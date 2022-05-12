package reconcilers

import (
	"context"
	"strings"

	"github.com/kyma-incubator/compass/components/atom-tenants-controller/pkg"

	"github.com/kyma-incubator/compass/components/atom-tenants-controller/internal/model"
	rmerrors "github.tools.sap/unified-resource-manager/api/pkg/apis/errors"
	rmlogger "github.tools.sap/unified-resource-manager/api/pkg/apis/logger"
	metav1 "github.tools.sap/unified-resource-manager/api/pkg/apis/meta/v1"
	"github.tools.sap/unified-resource-manager/api/pkg/apis/runtime"
	rmclient "github.tools.sap/unified-resource-manager/api/pkg/client"
	rmcontroller "github.tools.sap/unified-resource-manager/controller-utils/pkg/controller"
	rmmanager "github.tools.sap/unified-resource-manager/controller-utils/pkg/manager"
)

type ResourceGroupController struct {
	rmclient.Client
	Creator pkg.TenantCreator
	Log     rmlogger.Logger
}

func (r *ResourceGroupController) Reconcile(ctx context.Context, resourceKey runtime.Key) (rmcontroller.Result, error) {
	r.Log.Info("begin ResourceGroup reconciliation")
	resourceGroup := model.NewResourceGroup()
	err := r.Get(ctx, resourceKey, resourceGroup)
	if err != nil {
		if rmerrors.IsNotFound(err) {
			r.Log.Info("ResourceGroup was not found. Probably already deleted.")
			return rmcontroller.Result{}, nil
		}
		r.Log.Error(err, "Unable to fetch the ResourceGroup")
		return rmcontroller.Result{}, err
	}

	tenantHierarchy := strings.Split(resourceGroup.Path, PathDelimiter)
	organizationName := tenantHierarchy[1]

	var folders []pkg.Tenant
	for i := 2; i < len(tenantHierarchy); i++ {
		currFolder := tenantHierarchy[i]
		pathToCurrentFolder := strings.Join(tenantHierarchy[:i], PathDelimiter)
		folders = append(folders, pkg.Tenant{
			Name: currFolder,
			Path: pathToCurrentFolder + PathDelimiter + currFolder,
		})

	}

	crmID, err := getCustomerIDForOrganization(ctx, r.Client, organizationName, PathDelimiter)
	if err != nil {
		return rmcontroller.Result{}, err
	}

	payload := pkg.RequestPayload{
		Customer: crmID,
		Organization: pkg.Tenant{
			Name: organizationName,
			Path: PathDelimiter + organizationName,
		},
		Folders: folders,
		ResourceGroup: &pkg.Tenant{
			Name: resourceGroup.Name,
			Path: resourceGroup.Path + PathDelimiter + resourceGroup.Name,
		},
	}

	if err = r.Creator.StoreTenants(ctx, payload); err != nil {
		r.Log.Error(err, "while storing tenants")
		return rmcontroller.Result{}, err
	}

	return rmcontroller.Result{}, nil
}

func (r *ResourceGroupController) RunWithController(ctx context.Context) error {
	r.initFolderController(1).StartAsync(ctx)
	return nil
}

func (r *ResourceGroupController) ControllerWithManager(controllerMng *rmmanager.ControllerManager, maxConcurrentThreads int) error {
	r.Client = controllerMng.Client
	controller := r.initFolderController(maxConcurrentThreads)
	return controllerMng.AddController(controller)
}

func (r *ResourceGroupController) initFolderController(maxConcurrentThreads int) *rmcontroller.Controller {
	rootPath := PathDelimiter
	options := []metav1.ListOption{metav1.InPath(rootPath)}
	changeTypes := []runtime.ChangeType{runtime.ResourceCreate}

	return rmcontroller.NewController(
		"resource_group_controller",
		r.Client,
		r.Client,
		maxConcurrentThreads,
		r,
		model.NewResourceGroup(),
		rmcontroller.WatchParams{Options: options, ChangeTypes: changeTypes})
}
