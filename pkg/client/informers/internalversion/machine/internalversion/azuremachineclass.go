// This file was automatically generated by informer-gen

package internalversion

import (
	machine "github.com/gardener/node-controller-manager/pkg/apis/machine"
	internalinterfaces "github.com/gardener/node-controller-manager/pkg/client/informers/internalversion/internalinterfaces"
	internalclientset "github.com/gardener/node-controller-manager/pkg/client/internalclientset"
	internalversion "github.com/gardener/node-controller-manager/pkg/client/listers/machine/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// AzureMachineClassInformer provides access to a shared informer and lister for
// AzureMachineClasses.
type AzureMachineClassInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.AzureMachineClassLister
}

type azureMachineClassInformer struct {
	factory internalinterfaces.SharedInformerFactory
}

// NewAzureMachineClassInformer constructs a new informer for AzureMachineClass type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAzureMachineClassInformer(client internalclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				return client.Machine().AzureMachineClasses().List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				return client.Machine().AzureMachineClasses().Watch(options)
			},
		},
		&machine.AzureMachineClass{},
		resyncPeriod,
		indexers,
	)
}

func defaultAzureMachineClassInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewAzureMachineClassInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func (f *azureMachineClassInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&machine.AzureMachineClass{}, defaultAzureMachineClassInformer)
}

func (f *azureMachineClassInformer) Lister() internalversion.AzureMachineClassLister {
	return internalversion.NewAzureMachineClassLister(f.Informer().GetIndexer())
}
