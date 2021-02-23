package notifications

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/formation-watcher/script"
)

type RuntimeLabelNotificationHandler struct {
	RuntimeGetter  RuntimeGetter
	AppLister      ApplicationLister
	AppLabelGetter ApplicationLabelGetter
	Transact       persistence.Transactioner
	ScriptRunner   script.Runner
}

func (a *RuntimeLabelNotificationHandler) HandleCreate(ctx context.Context, label Label) error {
	// return a.handle(ctx, label)
	return nil
}

func (a *RuntimeLabelNotificationHandler) HandleUpdate(ctx context.Context, label Label) error {
	return nil
	// return a.handle(ctx, label)
}

func (a *RuntimeLabelNotificationHandler) HandleDelete(ctx context.Context, label Label) error {
	return nil

	// if label.Key != model.ScenariosKey {
	// 	log.C(ctx).Infof("label %v is not scenarios", label)
	// 	return nil
	// }

	// if len(label.RuntimeID) == 0 {
	// 	log.C(ctx).Infof("label %v is not for runtime", label)
	// 	return nil
	// }

	// tx, err := a.Transact.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer a.Transact.RollbackUnlessCommitted(ctx, tx)
	// ctx = persistence.SaveToContext(ctx, tx)
	// ctx = tenant.SaveToContext(ctx, label.TenantID, "")
	// runtime, err := a.RuntimeGetter.Get(ctx, label.RuntimeID)
	// if err != nil {
	// 	if apperrors.IsNotFoundError(err) {
	// 		log.C(ctx).Infof("runtime with id %s not found. Skipping label event", label.RuntimeID)
	// 		err = tx.Commit()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	}
	// 	return err
	// }

	// if runtime.Name != runtimeName {
	// 	log.C(ctx).Infof("event is not for the test runtime %s but for %s, skipping", runtimeName, runtime.Name)
	// 	return nil
	// }

	// if err := a.ScriptRunner.DeleteDependency(ctx, "dep-rt-"+label.RuntimeID, "admiral.yaml", "runtime.yaml"); err != nil {
	// 	return err
	// }

	// // TODO: this logic might need rethinking - if runtime is removed from scenario then this logic here might not get the right apps necessary to remove the service entries
	// parsedID, err := uuid.Parse(runtime.ID)
	// if err != nil {
	// 	return err
	// }

	// appsList, err := a.AppLister.ListByRuntimeID(ctx, parsedID, 100, "")
	// if err != nil {
	// 	return err
	// }

	// appNames := make([]string, 0, appsList.TotalCount)
	// for _, app := range appsList.Data {
	// 	scenarioLabel, err := a.AppLabelGetter.GetLabel(ctx, app.ID, "scenarios")
	// 	if err != nil {
	// 		if apperrors.IsNotFoundError(err) {
	// 			log.C(ctx).Warnf("app with id %s does not have scenarios label, skipping", label.AppID)
	// 			continue
	// 		}
	// 		return err
	// 	}
	// 	scenarioLabelSlice := scenarioLabel.Value.([]interface{})
	// 	if len(scenarioLabelSlice) == 1 && scenarioLabelSlice[0] == "DEFAULT" {
	// 		log.C(ctx).Warnf("app with id %s is only in the DEFAULT scenario, skipping", label.AppID)
	// 		continue
	// 	}
	// 	appNames = append(appNames, app.Name)
	// }

	// if err := cleanupServiceEntries(ctx, a.ScriptRunner); err != nil {
	// 	log.C(ctx).Errorf("unable to cleanup service entries for applications as part of runtime label event: %s", err)
	// 	return err
	// }

	// return nil
}
