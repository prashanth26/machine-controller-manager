/*
Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved.

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

// Package drain is used to drain nodes
package drain

import (
	"sync"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/klog"
)

type VolumeAttachmentHandler struct {
	sync.Mutex
	workers []chan storagev1.VolumeAttachment
}

func NewVolumeAttachmentHandler() *VolumeAttachmentHandler {
	return &VolumeAttachmentHandler{
		Mutex:   sync.Mutex{},
		workers: []chan storagev1.VolumeAttachment{},
	}
}

func (v *VolumeAttachmentHandler) dispatcher(obj interface{}) {
	volumeAttachment := obj.(*storagev1.VolumeAttachment)
	if volumeAttachment == nil {
		klog.Errorf("Couldn't convert to volumeAttachment from object %v", obj)
	}

	v.Lock()
	defer v.Unlock()

	for i, worker := range v.workers {
		klog.V(3).Info("Dispatching request for PV %s to worker %d", volumeAttachment.Spec.Source.PersistentVolumeName, i)
		worker <- *volumeAttachment
	}
}

func (v *VolumeAttachmentHandler) AddVolumeAttachment(obj interface{}) {
	klog.V(3).Infof("Adding volume attachment object")
	v.dispatcher(obj)
}

func (v *VolumeAttachmentHandler) UpdateVolumeAttachment(oldObj, newObj interface{}) {
	klog.V(3).Info("Updating volume attachment object")
	v.dispatcher(newObj)
}

func (v *VolumeAttachmentHandler) AddWorker(obj interface{}) *chan storagev1.VolumeAttachment {
	klog.V(3).Infof("Adding new worker. Current active workers %d", len(v.workers))

	v.Lock()
	defer v.Unlock()

	newWorker := make(chan storagev1.VolumeAttachment)
	v.workers = append(v.workers, newWorker)

	return &newWorker
}

func (v *VolumeAttachmentHandler) DeleteWorker(*chan storagev1.VolumeAttachment) {
	klog.V(3).Infof("Deleting an existing worker. Current active workers %d", len(v.workers))

	v.Lock()
	defer v.Unlock()

	for i, worker := range v.workers {

	}
}
