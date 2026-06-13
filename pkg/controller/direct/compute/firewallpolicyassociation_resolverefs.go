// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compute

import (
	"context"
	"fmt"

	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/compute/v1beta1"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResolveAttachmentTarget(ctx context.Context, reader client.Reader, src client.Object, ref *krm.ComputeFirewallPolicyAssociationAttachmentTargetRef) (string, error) {
	if ref == nil {
		return "", fmt.Errorf("attachmentTargetRef must be specified")
	}

	if ref.External != "" {
		if ref.Name != "" {
			return "", fmt.Errorf("cannot specify both name and external on attachmentTargetRef")
		}
		return ref.External, nil
	}

	if ref.Name == "" {
		return "", fmt.Errorf("must specify either name or external on attachmentTargetRef")
	}

	// Organizations are not yet supported dynamically as KCC resources, so Name can only point to a Folder.
	if ref.Kind != "" && ref.Kind != "Folder" {
		return "", fmt.Errorf("unsupported kind on attachmentTargetRef: %s, only Folder is supported", ref.Kind)
	}

	folderRef := &refsv1beta1.FolderRef{
		Name:      ref.Name,
		Namespace: ref.Namespace,
	}
	folder, err := refsv1beta1.ResolveFolder(ctx, reader, src, folderRef)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("folders/%s", folder.FolderID), nil
}

func resolveFirewallPolicyAssociationRefs(ctx context.Context, reader client.Reader, obj *krm.ComputeFirewallPolicyAssociation) error {
	// Resolve attachmentTargetRef
	resolvedTarget, err := ResolveAttachmentTarget(ctx, reader, obj, &obj.Spec.AttachmentTargetRef)
	if err != nil {
		return fmt.Errorf("resolving attachmentTargetRef: %w", err)
	}
	obj.Spec.AttachmentTargetRef.External = resolvedTarget

	// Resolve firewallPolicyRef
	resolvedPolicy, err := refsv1beta1.ResolveComputeFirewallPolicy(ctx, reader, obj, &obj.Spec.FirewallPolicyRef)
	if err != nil {
		return fmt.Errorf("resolving firewallPolicyRef: %w", err)
	}
	obj.Spec.FirewallPolicyRef.External = resolvedPolicy.External
	return nil
}
