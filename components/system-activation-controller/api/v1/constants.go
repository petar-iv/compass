/*
This file contains constants exposed by the operator
*/

package v1

const (
	// Group and Version for the operator exposed RTDs
	Group   = "ucl.orchestration.api.sap"
	Version = "v1"

	// Types exposed by the operator
	SystemActivationType     = "SystemActivation"
	SystemActivationSingular = "systemactivation"

	// Finalizers added and removed by the operator, used to block garbage collection until reconcile is finished
	Finalizer = "ucl.orchestration.api.sap/systemActivation-finalizer"

	// Conditions are used in the Status section
	ReadyCondition = "Ready" // Ready = True/False is the convention of whether the resource reconcile has finished

	ReadyMessageTrue          = "Resource is ready."
	ReadyMessageFalse         = "Resource is not ready."
	ReadyMessageFalseDeletion = "Resource is being deleted."
)
