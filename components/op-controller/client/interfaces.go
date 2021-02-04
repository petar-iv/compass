package client

import (
	"context"

	"github.com/kyma-incubator/compass/components/op-controller/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type OperationsInterface interface {
	List(ctx context.Context, opts metav1.ListOptions) (*v1beta1.OperationList, error)
	Get(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Operation, error)
	Create(ctx context.Context, operation *v1beta1.Operation) (*v1beta1.Operation, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Update(ctx context.Context, operation *v1beta1.Operation) (*v1beta1.Operation, error)
}
