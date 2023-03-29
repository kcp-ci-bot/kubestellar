/*
Copyright 2023 The KCP Authors.

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

package placement

import (
	"sync"

	"k8s.io/klog/v2"

	"github.com/kcp-dev/logicalcluster/v3"

	edgeapi "github.com/kcp-dev/edge-mc/pkg/apis/edge/v1alpha1"
)

type setBinder struct {
	logger klog.Logger
	sync.Mutex
	resolvedWhatDifferencerConstructor  ResolvedWhatDifferencerConstructor
	resolvedWhereDifferencerConstructor ResolvedWhereDifferencerConstructor
	perCluster                          map[logicalcluster.Name]*setBindingForCluster
	singleBinder                        SingleBinder
	ProjectionMapProvider
}

type setBindingForCluster struct {
	*setBinder
	cluster          logicalcluster.Name
	perPlacement     map[string]*setBindingForPlacement
	joinXY           SetChangeReceiver[Pair[WorkloadPart, string /*epName*/]]
	joinYZ           SetChangeReceiver[Pair[string /*epName*/, edgeapi.SinglePlacement]]
	singleBindingOps SingleBindingOps
}

type setBindingForPlacement struct {
	*setBindingForCluster
	resolvedWhatReceiver  Receiver[WorkloadParts]
	resolvedWhereReceiver Receiver[ResolvedWhere]
}

func NewSetBinder(
	logger klog.Logger,
	resolvedWhatDifferencerConstructor ResolvedWhatDifferencerConstructor,
	resolvedWhereDifferencerConstructor ResolvedWhereDifferencerConstructor,
	bindingOrganizer BindingOrganizer,
	discovery APIMapProvider,
	resourceModes ResourceModes,
	eventHandler EventHandler,
) SetBinder {
	singleBinder, projectionMapProvider := bindingOrganizer(discovery, resourceModes, eventHandler)
	ans := &setBinder{
		logger:                              logger,
		resolvedWhatDifferencerConstructor:  resolvedWhatDifferencerConstructor,
		resolvedWhereDifferencerConstructor: resolvedWhereDifferencerConstructor,
		singleBinder:                        singleBinder,
		ProjectionMapProvider:               projectionMapProvider,
	}
	return ans
}

func (sb *setBinder) AsWhatReceiver() MappingReceiver[ExternalName, WorkloadParts] {
	return sbAsResolvedWhatReceiver{sb}
}

func (sb *setBinder) AsWhereReceiver() MappingReceiver[ExternalName, ResolvedWhere] {
	return sbAsResolvedWhereReceiver{sb}
}

type sbAsResolvedWhatReceiver struct{ *setBinder }

func (sb sbAsResolvedWhatReceiver) Receive(epName ExternalName, resolvedWhat WorkloadParts) {
	sb.Lock()
	defer sb.Unlock()
	sbc := sb.ensureCluster(epName.Cluster)
	sbc.singleBinder.Transact(func(sbo SingleBindingOps) {
		sbp := sbc.ensurePlacement(epName.Name)
		sbc.singleBindingOps = sbo
		sbp.resolvedWhatReceiver.Receive(resolvedWhat)
		sbc.singleBindingOps = nil
	})
}

type sbAsResolvedWhereReceiver struct{ *setBinder }

func (sb sbAsResolvedWhereReceiver) Receive(epName ExternalName, resolvedWhere ResolvedWhere) {
	sb.Lock()
	defer sb.Unlock()
	sbc := sb.ensureCluster(epName.Cluster)
	sbc.singleBinder.Transact(func(sbo SingleBindingOps) {
		sbp := sbc.ensurePlacement(epName.Name)
		sbc.singleBindingOps = sbo
		sbp.resolvedWhereReceiver.Receive(resolvedWhere)
		sbc.singleBindingOps = nil
	})
}

func (sb *setBinder) ensureCluster(cluster logicalcluster.Name) *setBindingForCluster {
	sbc := sb.perCluster[cluster]
	if sbc == nil {
		sbc = &setBindingForCluster{
			setBinder:    sb,
			cluster:      cluster,
			perPlacement: map[string]*setBindingForPlacement{},
		}
		sbc.joinXY, sbc.joinYZ = NewDynamicJoin[WorkloadPart, string, edgeapi.SinglePlacement](sb.logger, sbc)
		sb.perCluster[cluster] = sbc
	}
	return sbc
}

func (sbc *setBindingForCluster) Add(pair Pair[WorkloadPart, edgeapi.SinglePlacement]) bool {
	sbc.logger.V(4).Info("Adding joined pair", "cluster", sbc.cluster, "part", pair.First, "where", pair.Second)
	return sbc.singleBindingOps.AddBinding(pair.First, pair.Second)
}

func (sbc *setBindingForCluster) Remove(pair Pair[WorkloadPart, edgeapi.SinglePlacement]) bool {
	sbc.logger.V(4).Info("Removing joined pair", "cluster", sbc.cluster, "part", pair.First, "where", pair.Second)
	return sbc.singleBindingOps.RemoveBinding(pair.First, pair.Second)
}

func (sbc *setBindingForCluster) ensurePlacement(epName string) *setBindingForPlacement {
	sbp := sbc.perPlacement[epName]
	if sbp == nil {
		sbp = &setBindingForPlacement{
			setBindingForCluster: sbc,
		}
		sbp.resolvedWhatReceiver = sbc.resolvedWhatDifferencerConstructor(&TransformSetChangeReceiver[WorkloadPart, Pair[WorkloadPart, string]]{
			Transform: AddSecondFunc[WorkloadPart, string](epName),
			Inner:     sbp.joinXY,
		})
		sbp.resolvedWhereReceiver = sbc.resolvedWhereDifferencerConstructor(&TransformSetChangeReceiver[edgeapi.SinglePlacement, Pair[string, edgeapi.SinglePlacement]]{
			Transform: AddFirstFunc[string, edgeapi.SinglePlacement](epName),
			Inner:     sbp.joinYZ,
		})
		sbc.perPlacement[epName] = sbp
	}
	return sbp
}
