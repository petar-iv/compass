package k8s

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/op-controller/api/v1beta1"
	"github.com/kyma-incubator/compass/components/op-controller/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Scheduler struct {
	kcli client.OperationsInterface
}

func NewScheduler(kcli client.OperationsInterface) *Scheduler {
	return &Scheduler{
		kcli: kcli,
	}
}

func (s *Scheduler) Schedule(ctx context.Context, op *v1beta1.Operation) error {
	result, err := s.kcli.Create(ctx, op)
	fmt.Println(">>>>>>>>>", result.CreationTimestamp)
	return err
}

func (s *Scheduler) Watch(ctx context.Context) error {
	w, err := s.kcli.Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}
	ch := w.ResultChan()
	for ev := range ch {
		op := ev.Object.(*v1beta1.Operation)
		if op.Status.Ready {
			fmt.Printf("op with name %s is ready\n", op.Name)
		} else {
			fmt.Printf("op with name %s is NOT ready\n", op.Name)
		}
	}
	return nil
}
