package notifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal2/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal2/labelfilter"
	"github.com/kyma-incubator/compass/components/director/internal2/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/formation-watcher/pkg/destination"
	"github.com/kyma-incubator/compass/components/formation-watcher/pkg/log"
	"github.com/kyma-incubator/compass/components/formation-watcher/script"
)

const runtimeName = "sap-graph"
const commerceMockName = "commerce-mock"

type AppData struct {
	ID   string
	Name string
}

type AppLabelNotificationHandler struct {
	RuntimeLister      RuntimeLister
	AppLister          ApplicationLister
	AppLabelGetter     ApplicationLabelGetter
	RuntimeLabelGetter RuntimeLabelGetter
	BundleGetter       BundleGetter
	Transact           persistence.Transactioner
	DestinationCient   DestinationCient
	DestinationsData   []destination.Destination
}

func (a *AppLabelNotificationHandler) HandleCreate(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *AppLabelNotificationHandler) HandleUpdate(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *AppLabelNotificationHandler) HandleDelete(ctx context.Context, label Label) error {
	return a.handle(ctx, label)
}

func (a *AppLabelNotificationHandler) handle(ctx context.Context, label Label) error {
	if label.Key != model.ScenariosKey {
		log.C(ctx).Infof("label %v is not scenarios", label)
		return nil
	}

	if len(label.AppID) == 0 {
		log.C(ctx).Infof("label %v is not for apps", label)
		return nil
	}

	tx, err := a.Transact.Begin()
	if err != nil {
		return err
	}
	defer a.Transact.RollbackUnlessCommitted(ctx, tx)

	ctx = persistence.SaveToContext(ctx, tx)
	ctx = tenant.SaveToContext(ctx, label.TenantID, "")
	query := `$[*] ? ( `
	queryEnd := ` )`
	queries := make([]string, 0, len(label.Value))
	for _, val := range label.Value {
		queries = append(queries, fmt.Sprintf("@ == \"%s\"", val))
	}
	query = query + strings.Join(queries, "||") + queryEnd
	runtimesList, err := a.RuntimeLister.List(ctx, []*labelfilter.LabelFilter{
		labelfilter.NewForKeyWithQuery(model.ScenariosKey, query),
	}, 100, "")
	if err != nil {
		return err
	}
	for _, runtime := range runtimesList.Data {
		if runtime.Name != runtimeName {
			log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", runtimeName, runtime.Name)
			continue
		}

		scenarioLabel, err := a.RuntimeLabelGetter.GetLabel(ctx, runtime.ID, "scenarios")
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				log.C(ctx).Warnf("runtime with id %s does not have scenarios label, skipping", runtime.ID)
				continue
			}
			return err
		}
		scenarioLabelSlice := scenarioLabel.Value.([]interface{})
		if len(scenarioLabelSlice) == 1 && scenarioLabelSlice[0] == "DEFAULT" {
			log.C(ctx).Warnf("runtime with id %s is only in the DEFAULT scenario, skipping", runtime.ID)
			continue
		}

		parsedID, err := uuid.Parse(runtime.ID)
		if err != nil {
			return err
		}

		appsList, err := a.AppLister.ListByRuntimeID(ctx, parsedID, 100, "")
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				log.C(ctx).Warnf("app with id %s not found during handling of label event", label.AppID)
				err = tx.Commit()
				if err != nil {
					return err
				}
				return nil
			}
			return err
		}

		appToSync := make([]AppData, 0, appsList.TotalCount)
		// appNames := make([]string, 0, appsList.TotalCount)
		for _, app := range appsList.Data {
			// if app.Status.Condition != model.ApplicationStatusConditionConnected {
			// 	log.C(ctx).Infof("app %s is not connected but is in status %s", app.Name, app.Status.Condition)
			// 	continue
			// }
			scenarioLabel, err := a.AppLabelGetter.GetLabel(ctx, app.ID, "scenarios")
			if err != nil {
				if apperrors.IsNotFoundError(err) {
					log.C(ctx).Warnf("app with id %s does not have scenarios label, skipping", label.AppID)
					continue
				}
				return err
			}
			scenarioLabelSlice := scenarioLabel.Value.([]interface{})
			if len(scenarioLabelSlice) == 1 && scenarioLabelSlice[0] == "DEFAULT" {
				log.C(ctx).Warnf("app with id %s is only in the DEFAULT scenario, skipping", label.AppID)
				continue
			}

			appToSync = append(appToSync, AppData{
				ID:   app.ID,
				Name: app.Name,
			})
		}

		log.C(ctx).Infof("Number of applications in scenario with test runtime: %d", len(appToSync))

		if err := a.syncDestinations(ctx, appToSync); err != nil {
			log.C(ctx).Errorf("unable to sync destinations for applications as part of application label event: %s", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (a *AppLabelNotificationHandler) syncDestinations(ctx context.Context, appData []AppData) error {
	// Get all Destinations
	currentDests, err := a.DestinationCient.GetAllDestinations(ctx)
	if err != nil {
		return err
	}
	destsMap := destsToMap(currentDests)
	expectedBundles := make([]*model.Bundle, 0)
	fmt.Printf(">>>>>>\n%+v\n>>>>>>>>\n", destsMap)
	for _, app := range appData {
		bundles, err := a.BundleGetter.ListByApplicationID(ctx, app.ID, 100, "")
		if err != nil {
			return err
		}
		// bundleName := bundles.Data[0].Name
		if len(bundles.Data) > 0 {
			expectedBundles = append(expectedBundles, bundles.Data[0])
		}
		// TODO: Get destination by bundle Name
	}

	destsToCreate := make([]string, 0)
	for _, expectedBndl := range expectedBundles {
		if _, found := destsMap[expectedBndl.Name]; !found {
			destsToCreate = append(destsToCreate, expectedBndl.Name)
		}
		delete(destsMap, expectedBndl.Name)
	}

	for _, destName := range destsToCreate {
		fmt.Printf("Destination to create: %s\n", destName)
		destToCreate, found := a.findDestinationToCreate(destName)
		if !found {
			return fmt.Errorf("Could not find data for destination with name %s", destName)
		}
		if err := a.DestinationCient.CreateDestination(ctx, destToCreate); err != nil {
			return err
		}
	}

	for name := range destsMap {
		fmt.Printf("Destination to delete: %s\n", name)
		if err := a.DestinationCient.DeleteDestination(ctx, name); err != nil {
			return err
		}
	}

	return nil
}

func (a *AppLabelNotificationHandler) findDestinationToCreate(name string) (destination.Destination, bool) {
	for _, destData := range a.DestinationsData {
		if destData["Name"].(string) == name {
			return destData, true
		}
	}
	return nil, false
}

func destsToMap(dests []destination.Destination) map[string]destination.Destination {
	result := make(map[string]destination.Destination)
	for _, dest := range dests {
		result[dest["Name"].(string)] = dest
	}
	return result
}

func syncServiceEntries(ctx context.Context, scriptRunner script.Runner, expectedAppNames []string) error {
	existingAppNames, err := scriptRunner.GetExistingServices(ctx)
	if err != nil {
		return err
	}

	appsToDelete := make([]string, 0)
	for _, app := range existingAppNames {
		if stringsAnyEquals(expectedAppNames, app) {
			continue
		} else {
			appsToDelete = append(appsToDelete, app)
		}
	}

	appsToCreate := make([]string, 0)
	for _, app := range expectedAppNames {
		if stringsAnyEquals(existingAppNames, app) {
			continue
		} else {
			appsToCreate = append(appsToCreate, app)
		}
	}

	// delete all resources for systems removed from scenario
	for _, appName := range appsToDelete {
		if appName == commerceMockName {
			log.C(ctx).Infof("service resources won't be edited for the %q application", appName)
			continue
		}

		if err := scriptRunner.DeleteResource(ctx, fmt.Sprintf("service-entries/%s.yaml", appName)); err != nil {
			return err
		}
	}

	// apply resources for systems added to scenario
	for _, appName := range appsToCreate {
		if appName == commerceMockName {
			log.C(ctx).Infof("service resources won't be edited for the %q application", appName)
			continue
		}

		if err := scriptRunner.ApplyResource(ctx, fmt.Sprintf("service-entries/%s.yaml", appName)); err != nil {
			return err
		}
	}

	return nil
}

func cleanupServiceEntries(ctx context.Context, scriptRunner script.Runner) error {
	return scriptRunner.DeleteResource(ctx, "service-entries/")
}

// stringsAnyEquals returns true if any of the strings in the slice equal the given string.
func stringsAnyEquals(stringSlice []string, str string) bool {
	for _, v := range stringSlice {
		if v == str {
			return true
		}
	}
	return false
}
