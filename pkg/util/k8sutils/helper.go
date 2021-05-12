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

// Package k8sutils is used to provider helper consts and functions for k8s operations
package k8sutils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

const (
	// VolumeAttachmentGroupName group name
	VolumeAttachmentGroupName = "storage.k8s.io"
	// VolumenAttachmentKind is the kind used for VolumeAttachment
	VolumeAttachmentResourceName = "volumeattachments"
	// VolumeAttachmentResource is the kind used for evicting pods
	VolumeAttachmentResourceKind = "VolumeAttachment"
)

// IsResourceSupported uses Discovery API to find out if the server supports
// the give resource of groupName, resourceName & resourceKind
// If support, it will return its groupVersion; Otherwise, it will return ""
func IsResourceSupported(
	clientset kubernetes.Interface,
	GroupName string,
	ResourceName string,
	ResourceKind string,
) bool {
	var (
		foundDesiredGroup   bool
		desiredGroupVersion string
	)

	discoveryClient := clientset.Discovery()
	groupList, err := discoveryClient.ServerGroups()
	if err != nil {
		return false
	}

	for _, group := range groupList.Groups {
		if group.Name == GroupName {
			foundDesiredGroup = true
			desiredGroupVersion = group.PreferredVersion.GroupVersion
			break
		}
	}
	if !foundDesiredGroup {
		return false
	}

	resourceList, err := discoveryClient.ServerResourcesForGroupVersion(desiredGroupVersion)
	if err != nil {
		return false
	}

	for _, resource := range resourceList.APIResources {
		if resource.Name == ResourceName && resource.Kind == ResourceKind {
			klog.V(4).Infof("Found Resource Name: %q, Resource Kind: %q", resource.Name, resource.Kind)
			return true
		}
	}
	return false
}
