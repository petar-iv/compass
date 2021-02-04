package k8s

import (
	"context"
	"fmt"
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/async"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/op-controller/api/v1beta1"
	"github.com/kyma-incubator/compass/components/op-controller/client"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type k8sClient interface {
	Create(ctx context.Context, operation *v1beta1.Operation) (*v1beta1.Operation, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Operation, error)
	Update(ctx context.Context, operation *v1beta1.Operation) (*v1beta1.Operation, error)
}

type Scheduler struct {
	kcli                         k8sClient
	lastProcessedResourceVersion string
	restartTimeout               time.Duration
}

func NewScheduler(kcli client.OperationsInterface, restartTimeout time.Duration) *Scheduler {
	return &Scheduler{
		kcli:                         kcli,
		restartTimeout:               restartTimeout,
		lastProcessedResourceVersion: "",
	}
}

func (s *Scheduler) Schedule(ctx context.Context, op *v1beta1.Operation) error {
	getOp, err := s.kcli.Get(ctx, op.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err = s.kcli.Create(ctx, op)
		return err
	}
	// TODO: Check if operation is in progress, if true: return an error that op is in progress; otherwise proceed with op update
	getOp.Spec = op.Spec
	_, err = s.kcli.Update(ctx, getOp)
	if errors.IsConflict(err) {
		return fmt.Errorf("another operation is in progress")
	}
	return err
}

func (s *Scheduler) Watch(ctx context.Context, processFunc async.OperationProcessor) error {
	for {
		select {
		case <-ctx.Done():
			log.C(ctx).Info("Context cancelled. Will not start scheduler watch...")
		default:
		}

		log.C(ctx).Infof("Starting scheduler watch..")
		w, err := s.kcli.Watch(ctx, metav1.ListOptions{
			ResourceVersion:      s.lastProcessedResourceVersion,
			ResourceVersionMatch: metav1.ResourceVersionMatchNotOlderThan,
		})
		if err != nil {
			// TODO: Track prometheus data that connection was not possible
			log.C(ctx).Errorf("Could not start scheduler watch: %s", err)
			time.Sleep(s.restartTimeout)
			continue
		}

		go func() {
			<-ctx.Done()
			log.C(ctx).Info("Context cancelled. Stopping scheduler watch...")
			w.Stop()
		}()

		ch := w.ResultChan()
		for ev := range ch {
			if err := s.processEvent(ctx, ev, processFunc); err != nil {
				log.C(ctx).Errorf("Could not process operation event: %s", err)
				break
			}
		}

		select {
		case <-ctx.Done():
		default:
			log.C(ctx).Error("Unexpected scheduler watch closed. Restarting scheduler watch...")
			w.Stop()
			time.Sleep(s.restartTimeout)
			continue
		}

		log.C(ctx).Info("Scheduler watcher has stopped")
		return nil
	}
}

func (s *Scheduler) processEvent(ctx context.Context, ev watch.Event, processFunc async.OperationProcessor) error {
	log.C(ctx).Infof("Event received %+v", ev.Type)
	switch op := ev.Object.(type) {
	case *v1beta1.Operation:
		// TODO: Use operation's correlationID here
		if err := processFunc(op); err != nil {
			return err
		}
		s.lastProcessedResourceVersion = op.ResourceVersion
	default:
		log.C(ctx).Errorf("Unexpected scheduler event received: %+v, %T", op, op)
	}
	return nil
}
