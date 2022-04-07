package controllers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	v1 "github.com/kyma-incubator/compass/components/system-activation-controller/api/v1"
	"github.com/pkg/errors"
	corev1 "github.tools.sap/unified-resource-manager/api/pkg/apis/core/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/meta"
	metav1 "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	"github.tools.sap/unified-resource-manager/api/pkg/client"
	controllers "github.tools.sap/unified-resource-manager/controller-utils/pkg/controller"
	"github.tools.sap/unified-resource-manager/controller-utils/pkg/manager"
	"net/http"
	"time"
)

// SystemActivationReconciler reconciles an SystemActivation object
type SystemActivationReconciler struct {
	HttpClient *http.Client
	client.Client
	Log logr.Logger
}

/*
	The Reconcile function is called each time an event for a SystemActivation is sent by the Unified Resource Manager.
	Parameter runtime.ResourceKey is an object that specifies the Name and Path of the SystemActivation.
	Returns:
	1. controllers.Result - empty if a requeue of the event is not required.
	2. error - if an error is returned, the event will be requeued.
*/
func (r *SystemActivationReconciler) Reconcile(ctx context.Context, key runtime.Key) (controllers.Result, error) {
	log := r.Log.WithValues("path", key.GetPath(), "name", key.GetName())
	log.Info("Starting SystemActivation reconcile")

	systemActivation := v1.NewSystemActivation()

	// Indicator for status changes - indicates whether resource should be updated or not.
	var isChanged = false

	if err := r.Get(ctx, key, systemActivation); err != nil {
		// If the object was not found, we have no SystemActivation to reconcile
		if IsNotFound(err) {
			log.Info("SystemActivation object not found.")
		}
		return controllers.Result{}, IgnoreNotFound(err)
	}

	// We want to let the consumer know the Desired State he requested has not reached yet.
	// It is important to update the Status as soon as possible, so the consumer will have an indication of the work being done.
	if len(systemActivation.Status.Conditions) == 0 {
		log.Info("Init conditions")
		systemActivation.Status.Conditions = []metav1.Condition{
			{
				Type:               v1.ReadyCondition,
				Status:             metav1.ConditionFalse,
				Reason:             "NotReady",
				Message:            v1.ReadyMessageFalse,
				ObservedGeneration: 1,
				LastTransitionTime: time.Now(),
			},
		}

		return controllers.Result{}, r.Update(ctx, systemActivation)
	}

	if IsObjectBeingDeleted(systemActivation) {
		systemActivation.Status.Conditions = []metav1.Condition{
			{
				Type:               v1.ReadyCondition,
				Status:             metav1.ConditionFalse,
				Reason:             "NotReady",
				Message:            v1.ReadyMessageFalseDeletion,
				ObservedGeneration: systemActivation.Generation,
				LastTransitionTime: time.Now(),
			},
		}

		return controllers.Result{}, r.handleDeletion(ctx, log, systemActivation)
	}

	if !meta.ContainsFinalizer(systemActivation, v1.Finalizer) {
		log.Info("Adding finalizer")
		meta.AddFinalizer(systemActivation, v1.Finalizer)
		return controllers.Result{}, r.Update(ctx, systemActivation)
	}

	if err := r.handleSecret(ctx, log, systemActivation); err != nil {
		return controllers.Result{}, err
	}

	if systemActivation.Status.SecretRef == nil {
		systemActivation.Status.SecretRef = &v1.SecretRefStatus{
			Name: systemActivation.Spec.Secret.Name,
			Path: systemActivation.Path,
		}
		isChanged = true
	}

	// We want to let the consumer know the Desired State he requested is now ready
	if !metav1.IsConditionTrue(systemActivation.Status.Conditions, v1.ReadyCondition) {
		readyCondition := metav1.Condition{
			Type:               v1.ReadyCondition,
			Status:             metav1.ConditionTrue,
			Reason:             "Ready",
			Message:            v1.ReadyMessageTrue,
			ObservedGeneration: systemActivation.Generation,
			LastTransitionTime: time.Now(),
		}
		metav1.SetCondition(&systemActivation.Status.Conditions, readyCondition)
		isChanged = true
	}

	// Update the resource if needed (causing re-triggering the reconciliation process with no changed)
	if isChanged == false {
		return controllers.Result{}, nil
	}

	if err := r.Update(ctx, systemActivation); err != nil {
		log.Error(err, "Update failed")
		return controllers.Result{}, err
	}
	return controllers.Result{}, nil
}

// ControllerWithManager registers the SystemActivationController as a controller in the controllers manager
func (r *SystemActivationReconciler) ControllerWithManager(mgr *manager.ControllerManager) error {
	watchPath := fmt.Sprintf("/service-orchestration/managed-service-workspaces/%s-%s", v1.Group, v1.SystemActivationSingular)
	controller := controllers.NewController("SystemActivationController", r.Client, r.Client, 1, r, v1.NewSystemActivation(),
		controllers.WatchParams{Options: []metav1.ListOption{metav1.InPath(watchPath)}})
	return mgr.AddController(controller)
}

// handleDeletion is called when a delete process is required
func (r *SystemActivationReconciler) handleDeletion(ctx context.Context, log logr.Logger, systemActivation *v1.SystemActivation) error {
	log.Info("Handling deletion")

	// Remove the secret
	if systemActivation.Status.SecretRef != nil {
		if err := r.deleteSecret(ctx, log, runtime.ResourceKey{Path: systemActivation.Path, Name: systemActivation.Status.SecretRef.Name}); err != nil {
			log.Error(err, "Secret deletion error")
			return err
		}
	}

	// After we are done and made sure the object can be safely removed from Unified Resource Manager, we can remove the finalizer
	if meta.ContainsFinalizer(systemActivation, v1.Finalizer) {
		log.Info("Removing finalizer")
		meta.RemoveFinalizer(systemActivation, v1.Finalizer)
		return r.Update(ctx, systemActivation)
	}

	return nil
}

// handleSecret checks if a secret creation/maintenance is required, and executes
func (r *SystemActivationReconciler) handleSecret(ctx context.Context, log logr.Logger, systemActivation *v1.SystemActivation) error {
	log.Info("Handling Secret")
	if systemActivation.Status.SecretRef != nil {
		log.Info("Secret already present")
		return nil
	}

	externalAppURL := systemActivation.Spec.URL
	req, err := http.NewRequest("GET", externalAppURL, nil)
	if err != nil {
		log.Error(err, "New request error")
		return err
	}

	resp, err := r.HttpClient.Do(req)
	if err != nil {
		log.Error(err, "HttpClient request error")
		return err
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error(err, "Got error on closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("wrong status code, got [%d], expected [%d]", resp.StatusCode, http.StatusOK)
		log.Error(err, "Wrong status code")
		return err
	}

	responseBody := SystemActivationCredentials{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		log.Error(err, "While decoding response")
		return errors.Wrap(err, "while decoding response")
	}

	secretName := systemActivation.Spec.Secret.Name

	if err := r.createSecret(ctx, responseBody.Secret, log, runtime.ResourceKey{Path: systemActivation.Path, Name: secretName}); err != nil {
		return err
	}

	return nil
}

func (r *SystemActivationReconciler) deleteSecret(ctx context.Context, log logr.Logger, resourceKey runtime.ResourceKey) error {
	secret := corev1.NewSecret()
	if err := r.Get(ctx, resourceKey, secret); err != nil {
		if !IsNotFound(err) {
			log.Error(err, FormatMessageResource("Couldn't delete secret.", resourceKey))
			return err
		}
		log.Info(FormatMessageResource("Secret does not exist.", resourceKey))
		return nil
	} else {
		log.Info(FormatMessageResource("Deleting secret.", resourceKey))
		err := r.Delete(ctx, secret)
		return err
	}
}

func (r *SystemActivationReconciler) createSecret(ctx context.Context, secretValue string, log logr.Logger, resourceKey runtime.ResourceKey) error {
	secret := corev1.NewSecret()
	secret.Path = resourceKey.Path
	secret.Name = resourceKey.Name
	encoded := base64.StdEncoding.EncodeToString([]byte(secretValue))
	secret.Data = map[string]string{
		"secret": encoded,
	}

	log.Info(FormatMessageResource("Creating Secret.", resourceKey))
	if err := r.Create(ctx, secret); err != nil {
		if IsResourceAlreadyExists(err) {
			log.Info(FormatMessageResource("Secret already exists.", resourceKey))
			return nil
		}
		log.Error(err, FormatMessageResource("Cannot create Secret.", resourceKey))
		return err
	}
	return nil
}
