package reconcilers

import (
	"context"

	"github.com/kyma-incubator/compass/components/atom-tenants-controller/internal/model"
	"github.com/kyma-incubator/compass/components/atom-tenants-controller/pkg"
	"github.com/pkg/errors"
	rmerrors "github.tools.sap/unified-resource-manager/api/pkg/apis/errors"
	rmlogger "github.tools.sap/unified-resource-manager/api/pkg/apis/logger"
	metav1 "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	rmclient "github.tools.sap/unified-resource-manager/api/pkg/client"
	rmcontroller "github.tools.sap/unified-resource-manager/controller-utils/pkg/controller"
	rmmanager "github.tools.sap/unified-resource-manager/controller-utils/pkg/manager"
)

type OrganizationController struct {
	rmclient.Client
	Creator pkg.TenantCreator
	Log     rmlogger.Logger
}

func (r *OrganizationController) Reconcile(ctx context.Context, resourceKey runtime.Key) (rmcontroller.Result, error) {
	r.Log.Info("begin Organization reconciliation")
	org := model.NewOrganization()
	if err := r.Get(ctx, resourceKey, org); err != nil {
		if rmerrors.IsNotFound(err) {
			r.Log.Info("Organization was not found. Probably already deleted.")
			return rmcontroller.Result{}, nil
		}
		r.Log.Error(err, "Unable to fetch the Organization")
		return rmcontroller.Result{}, err
	}

	if !metav1.IsConditionTrue(org.Status.Conditions, "Ready") {
		return rmcontroller.Result{}, errors.New("Organization status was not Ready")
	}

	crmID, err := getCustomerIDForOrganization(ctx, r.Client, org.Name, org.Path)
	if err != nil {
		return rmcontroller.Result{}, err
	}

	payload := pkg.RequestPayload{
		Customer: crmID,
		Organization: pkg.Tenant{
			Name: org.Name,
			Path: PathDelimiter + org.Name,
		},
	}

	if err = r.Creator.StoreTenants(ctx, payload); err != nil {
		r.Log.Error(err, "while storing tenants")
		return rmcontroller.Result{}, err
	}

	return rmcontroller.Result{}, nil
}

func (r *OrganizationController) RunWithController(ctx context.Context) error {
	r.initOrganizationController(1).StartAsync(ctx)
	return nil
}

func (r *OrganizationController) ControllerWithManager(controllerMng *rmmanager.ControllerManager, maxConcurrentThreads int) error {
	r.Client = controllerMng.Client
	controller := r.initOrganizationController(maxConcurrentThreads)
	return controllerMng.AddController(controller)
}

func (r *OrganizationController) initOrganizationController(maxConcurrentThreads int) *rmcontroller.Controller {
	rootPath := "/"
	options := []metav1.ListOption{metav1.InPath(rootPath)}
	changeTypes := []runtime.ChangeType{runtime.ResourceCreate}

	return rmcontroller.NewController(
		"organization_tenant_controller",
		r.Client,
		r.Client,
		maxConcurrentThreads,
		r,
		model.NewOrganization(),
		rmcontroller.WatchParams{Options: options, ChangeTypes: changeTypes})
}

func getCustomerIDForOrganization(ctx context.Context, client rmclient.Client, orgName string, orgPath string) (string, error) {
	orgBaseKey := runtime.FullResourceKey{
		Path:             orgPath,
		Name:             orgName,
		GroupVersionType: model.NewOrganizationGVT(),
	}
	orgBase := model.NewOrganizationBase()
	if err := client.Get(ctx, orgBaseKey, orgBase); err != nil {
		if rmerrors.IsNotFound(err) { // referred organization is internal
			return "", nil
		}
		return "", err
	}
	return orgBase.Labels["accounts.commercial.resource.api.sap/crm-id"], nil

}
