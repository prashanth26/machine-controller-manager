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
	if len(v.workers) == 0 {
		// As no workers are registered, nothing to do here.
		return
	}

	volumeAttachment := obj.(*storagev1.VolumeAttachment)
	if volumeAttachment == nil {
		klog.Errorf("Couldn't convert to volumeAttachment from object %v", obj)
	}

	klog.V(3).Infof("Dispatching request for PV %s", *volumeAttachment.Spec.Source.PersistentVolumeName)

	v.Lock()
	defer v.Unlock()

	for i, worker := range v.workers {
		klog.V(3).Infof("Dispatching request for PV %s to worker %d", *volumeAttachment.Spec.Source.PersistentVolumeName, i)
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

func (v *VolumeAttachmentHandler) AddWorker() chan storagev1.VolumeAttachment {
	klog.V(3).Infof("Adding new worker. Current active workers %d", len(v.workers))

	v.Lock()
	defer v.Unlock()

	newWorker := make(chan storagev1.VolumeAttachment, 20)
	v.workers = append(v.workers, newWorker)

	klog.V(3).Infof("Successfully added new worker %v. Current active workers %d", newWorker, len(v.workers))
	return newWorker
}

func (v *VolumeAttachmentHandler) DeleteWorker(desiredWorker chan storagev1.VolumeAttachment) {
	klog.V(3).Infof("Deleting an existing worker %v. Current active workers %d", desiredWorker, len(v.workers))

	v.Lock()
	defer v.Unlock()

	finalWorkers := []chan storagev1.VolumeAttachment{}

	for _, worker := range v.workers {
		if worker == desiredWorker {
			klog.V(3).Infof("Found worker to be deleted. Worker %s", worker)
		} else {
			finalWorkers = append(finalWorkers, worker)
		}

	}

	v.workers = finalWorkers
	klog.V(3).Infof("Successfully removed worker. Current active workers %d", len(v.workers))
}
