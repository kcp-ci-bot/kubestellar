//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The KubeStellar Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by kcp code-generator. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	kcpcache "github.com/kcp-dev/apimachinery/v2/pkg/cache"
	kcpinformers "github.com/kcp-dev/apimachinery/v2/third_party/informers"
	"github.com/kcp-dev/logicalcluster/v3"

	edgev1alpha1 "github.com/kcp-dev/edge-mc/pkg/apis/edge/v1alpha1"
	scopedclientset "github.com/kcp-dev/edge-mc/pkg/client/clientset/versioned"
	clientset "github.com/kcp-dev/edge-mc/pkg/client/clientset/versioned/cluster"
	"github.com/kcp-dev/edge-mc/pkg/client/informers/externalversions/internalinterfaces"
	edgev1alpha1listers "github.com/kcp-dev/edge-mc/pkg/client/listers/edge/v1alpha1"
)

// SyncerConfigClusterInformer provides access to a shared informer and lister for
// SyncerConfigs.
type SyncerConfigClusterInformer interface {
	Cluster(logicalcluster.Name) SyncerConfigInformer
	Informer() kcpcache.ScopeableSharedIndexInformer
	Lister() edgev1alpha1listers.SyncerConfigClusterLister
}

type syncerConfigClusterInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewSyncerConfigClusterInformer constructs a new informer for SyncerConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSyncerConfigClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredSyncerConfigClusterInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredSyncerConfigClusterInformer constructs a new informer for SyncerConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSyncerConfigClusterInformer(client clientset.ClusterInterface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) kcpcache.ScopeableSharedIndexInformer {
	return kcpinformers.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EdgeV1alpha1().SyncerConfigs().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EdgeV1alpha1().SyncerConfigs().Watch(context.TODO(), options)
			},
		},
		&edgev1alpha1.SyncerConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *syncerConfigClusterInformer) defaultInformer(client clientset.ClusterInterface, resyncPeriod time.Duration) kcpcache.ScopeableSharedIndexInformer {
	return NewFilteredSyncerConfigClusterInformer(client, resyncPeriod, cache.Indexers{
		kcpcache.ClusterIndexName: kcpcache.ClusterIndexFunc,
	},
		f.tweakListOptions,
	)
}

func (f *syncerConfigClusterInformer) Informer() kcpcache.ScopeableSharedIndexInformer {
	return f.factory.InformerFor(&edgev1alpha1.SyncerConfig{}, f.defaultInformer)
}

func (f *syncerConfigClusterInformer) Lister() edgev1alpha1listers.SyncerConfigClusterLister {
	return edgev1alpha1listers.NewSyncerConfigClusterLister(f.Informer().GetIndexer())
}

// SyncerConfigInformer provides access to a shared informer and lister for
// SyncerConfigs.
type SyncerConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() edgev1alpha1listers.SyncerConfigLister
}

func (f *syncerConfigClusterInformer) Cluster(clusterName logicalcluster.Name) SyncerConfigInformer {
	return &syncerConfigInformer{
		informer: f.Informer().Cluster(clusterName),
		lister:   f.Lister().Cluster(clusterName),
	}
}

type syncerConfigInformer struct {
	informer cache.SharedIndexInformer
	lister   edgev1alpha1listers.SyncerConfigLister
}

func (f *syncerConfigInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

func (f *syncerConfigInformer) Lister() edgev1alpha1listers.SyncerConfigLister {
	return f.lister
}

type syncerConfigScopedInformer struct {
	factory          internalinterfaces.SharedScopedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

func (f *syncerConfigScopedInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&edgev1alpha1.SyncerConfig{}, f.defaultInformer)
}

func (f *syncerConfigScopedInformer) Lister() edgev1alpha1listers.SyncerConfigLister {
	return edgev1alpha1listers.NewSyncerConfigLister(f.Informer().GetIndexer())
}

// NewSyncerConfigInformer constructs a new informer for SyncerConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSyncerConfigInformer(client scopedclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSyncerConfigInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredSyncerConfigInformer constructs a new informer for SyncerConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSyncerConfigInformer(client scopedclientset.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EdgeV1alpha1().SyncerConfigs().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.EdgeV1alpha1().SyncerConfigs().Watch(context.TODO(), options)
			},
		},
		&edgev1alpha1.SyncerConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *syncerConfigScopedInformer) defaultInformer(client scopedclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSyncerConfigInformer(client, resyncPeriod, cache.Indexers{}, f.tweakListOptions)
}
