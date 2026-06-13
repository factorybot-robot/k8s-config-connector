// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses///LICENSE-2.0
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

	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/common"

	"google.golang.org/api/option"

	gcp "cloud.google.com/go/compute/apiv1"
	computepb "cloud.google.com/go/compute/apiv1/computepb"
	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/compute/v1beta1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/config"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/directbase"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/registry"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	registry.RegisterModel(krm.ComputeFirewallPolicyAssociationGVK, NewFirewallPolicyAssociationModel)
}

func NewFirewallPolicyAssociationModel(ctx context.Context, config *config.ControllerConfig) (directbase.Model, error) {
	return &firewallPolicyAssociationModel{config: config}, nil
}

type firewallPolicyAssociationModel struct {
	config *config.ControllerConfig
}

// model implements the Model interface.
var _ directbase.Model = &firewallPolicyAssociationModel{}

type firewallPolicyAssociationAdapter struct {
	id                     *krm.FirewallPolicyAssociationIdentity
	firewallPoliciesClient *gcp.FirewallPoliciesClient
	desired                *krm.ComputeFirewallPolicyAssociation
	actual                 *computepb.FirewallPolicyAssociation
	reader                 client.Reader
}

var _ directbase.Adapter = &firewallPolicyAssociationAdapter{}

func (m *firewallPolicyAssociationModel) client(ctx context.Context) (*gcp.FirewallPoliciesClient, error) {
	var opts []option.ClientOption
	opts, err := m.config.RESTClientOptions()
	if err != nil {
		return nil, err
	}
	gcpClient, err := gcp.NewFirewallPoliciesRESTClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("building FirewallPolicy client: %w", err)
	}
	return gcpClient, err
}

func (m *firewallPolicyAssociationModel) AdapterForObject(ctx context.Context, op *directbase.AdapterForObjectOperation) (directbase.Adapter, error) {
	u := op.GetUnstructured()
	reader := op.Reader
	obj := &krm.ComputeFirewallPolicyAssociation{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &obj); err != nil {
		return nil, fmt.Errorf("error converting to %T: %w", obj, err)
	}

	id, err := krm.NewFirewallPolicyAssociationIdentity(ctx, reader, obj)
	if err != nil {
		return nil, err
	}

	adapter := &firewallPolicyAssociationAdapter{
		id:      id,
		desired: obj,
		reader:  reader,
	}

	// Get GCP client
	gcpClient, err := m.client(ctx)
	if err != nil {
		return nil, fmt.Errorf("building gcp client: %w", err)
	}
	adapter.firewallPoliciesClient = gcpClient

	return adapter, nil
}

func (m *firewallPolicyAssociationModel) AdapterForURL(ctx context.Context, url string) (directbase.Adapter, error) {
	// TODO: Support URLs
	return nil, nil
}

func (a *firewallPolicyAssociationAdapter) Find(ctx context.Context) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("getting ComputeFirewallPolicyAssociation", "name", a.id)

	association, err := a.get(ctx)
	if err != nil {
		if direct.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("getting ComputeFirewallPolicyAssociation %s: %w", a.id, err)
	}
	a.actual = association
	return true, nil
}

func (a *firewallPolicyAssociationAdapter) Create(ctx context.Context, createOp *directbase.CreateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("creating ComputeFirewallPolicyAssociation", "name", a.id)

	err := resolveFirewallPolicyAssociationRefs(ctx, a.reader, a.desired)
	if err != nil {
		return err
	}

	mapCtx := &direct.MapContext{}
	desired := a.desired.DeepCopy()
	firewallPolicyAssociation := ComputeFirewallPolicyAssociationSpec_v1beta1_ToProto(mapCtx, &desired.Spec)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	// Set name in association
	firewallPolicyAssociation.Name = direct.LazyPtr(a.id.Name())

	req := &computepb.AddAssociationFirewallPolicyRequest{
		FirewallPolicy:                    a.id.FirewallPolicy(),
		FirewallPolicyAssociationResource: firewallPolicyAssociation,
	}

	op, err := a.firewallPoliciesClient.AddAssociation(ctx, req)
	if err != nil {
		return fmt.Errorf("creating ComputeFirewallPolicyAssociation %s: %w", a.id, err)
	}

	if !op.Done() {
		err = op.Wait(ctx)
		if err != nil {
			return fmt.Errorf("waiting ComputeFirewallPolicyAssociation %s create failed: %w", a.id, err)
		}
	}
	log.V(2).Info("successfully created ComputeFirewallPolicyAssociation", "name", a.id)

	// Get the created resource
	created, err := a.get(ctx)
	if err != nil {
		return fmt.Errorf("getting ComputeFirewallPolicyAssociation %s: %w", a.id, err)
	}

	status := ComputeFirewallPolicyAssociationStatus_v1beta1_FromProto(mapCtx, created)
	status.ObservedGeneration = direct.LazyPtr(a.desired.GetGeneration())

	externalRef := a.id.String()
	status.ExternalRef = &externalRef

	return createOp.UpdateStatus(ctx, status, nil)
}

func (a *firewallPolicyAssociationAdapter) Update(ctx context.Context, updateOp *directbase.UpdateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("updating ComputeFirewallPolicyAssociation", "name", a.id)

	err := resolveFirewallPolicyAssociationRefs(ctx, a.reader, a.desired)
	if err != nil {
		return err
	}

	mapCtx := &direct.MapContext{}
	desired := a.desired.DeepCopy()
	firewallPolicyAssociation := ComputeFirewallPolicyAssociationSpec_v1beta1_ToProto(mapCtx, &desired.Spec)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	// Since fields are immutable, we don't expect changes. Let's run CompareProtoMessage.
	paths, err := common.CompareProtoMessage(firewallPolicyAssociation, a.actual, common.BasicDiff)
	if err != nil {
		return err
	}

	if len(paths) > 0 {
		return fmt.Errorf("cannot update ComputeFirewallPolicyAssociation %s: all spec fields are immutable", a.id)
	}

	status := ComputeFirewallPolicyAssociationStatus_v1beta1_FromProto(mapCtx, a.actual)
	status.ObservedGeneration = direct.LazyPtr(a.desired.GetGeneration())
	externalRef := a.id.String()
	status.ExternalRef = &externalRef

	return updateOp.UpdateStatus(ctx, status, nil)
}

func (a *firewallPolicyAssociationAdapter) Export(ctx context.Context) (*unstructured.Unstructured, error) {
	if a.actual == nil {
		return nil, fmt.Errorf("firewallPolicyAssociation %s not found", a.id)
	}

	mc := &direct.MapContext{}
	spec := ComputeFirewallPolicyAssociationSpec_v1beta1_FromProto(mc, a.actual)
	specObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(spec)
	if err != nil {
		return nil, fmt.Errorf("error converting firewallPolicyAssociation spec to unstructured: %w", err)
	}

	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(krm.ComputeFirewallPolicyAssociationGVK)

	if err := unstructured.SetNestedField(u.Object, specObj, "spec"); err != nil {
		return nil, fmt.Errorf("setting spec: %w", err)
	}

	return u, nil
}

// Delete implements the Adapter interface.
func (a *firewallPolicyAssociationAdapter) Delete(ctx context.Context, deleteOp *directbase.DeleteOperation) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("deleting ComputeFirewallPolicyAssociation", "name", a.id)

	req := &computepb.RemoveAssociationFirewallPolicyRequest{
		FirewallPolicy: a.id.FirewallPolicy(),
		Name:           direct.LazyPtr(a.id.Name()),
	}

	op, err := a.firewallPoliciesClient.RemoveAssociation(ctx, req)
	if err != nil {
		if direct.IsNotFound(err) {
			return true, nil
		}
		return false, fmt.Errorf("deleting ComputeFirewallPolicyAssociation %s: %w", a.id, err)
	}

	if !op.Done() {
		err = op.Wait(ctx)
		if err != nil {
			return false, fmt.Errorf("waiting ComputeFirewallPolicyAssociation %s delete failed: %w", a.id, err)
		}
	}
	log.V(2).Info("successfully deleted ComputeFirewallPolicyAssociation", "name", a.id)

	return true, nil
}

func (a *firewallPolicyAssociationAdapter) get(ctx context.Context) (*computepb.FirewallPolicyAssociation, error) {
	getReq := &computepb.GetAssociationFirewallPolicyRequest{
		FirewallPolicy: a.id.FirewallPolicy(),
		Name:           direct.LazyPtr(a.id.Name()),
	}
	return a.firewallPoliciesClient.GetAssociation(ctx, getReq)
}
