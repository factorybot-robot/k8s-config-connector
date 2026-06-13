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

package bigqueryreservation

import (
	"context"
	"fmt"

	gcp "cloud.google.com/go/bigquery/reservation/apiv1"
	pb "cloud.google.com/go/bigquery/reservation/apiv1/reservationpb"
	krm "github.com/GoogleCloudPlatform/k8s-config-connector/apis/bigqueryreservation/v1alpha1"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/config"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/directbase"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/controller/direct/registry"
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/structuredreporting"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
)

func init() {
	registry.RegisterModel(krm.BigQueryReservationCapacityCommitmentGVK, NewCapacityCommitmentModel)
}

func NewCapacityCommitmentModel(ctx context.Context, config *config.ControllerConfig) (directbase.Model, error) {
	return &modelCapacityCommitment{config: *config}, nil
}

var _ directbase.Model = &modelCapacityCommitment{}

type modelCapacityCommitment struct {
	config config.ControllerConfig
}

func (m *modelCapacityCommitment) client(ctx context.Context) (*gcp.Client, error) {
	var opts []option.ClientOption
	opts, err := m.config.RESTClientOptions()
	if err != nil {
		return nil, err
	}
	gcpClient, err := gcp.NewRESTClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("building CapacityCommitment client: %w", err)
	}
	return gcpClient, err
}

func (m *modelCapacityCommitment) AdapterForObject(ctx context.Context, op *directbase.AdapterForObjectOperation) (directbase.Adapter, error) {
	u := op.GetUnstructured()
	reader := op.Reader
	obj := &krm.BigQueryReservationCapacityCommitment{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &obj); err != nil {
		return nil, fmt.Errorf("error converting to %T: %w", obj, err)
	}

	res, err := obj.GetIdentity(ctx, reader)
	if err != nil {
		return nil, err
	}
	id := res.(*krm.BigQueryReservationCapacityCommitmentIdentity)

	// Get bigqueryreservation GCP client
	gcpClient, err := m.client(ctx)
	if err != nil {
		return nil, err
	}
	return &CapacityCommitmentAdapter{
		id:        id,
		gcpClient: gcpClient,
		desired:   obj,
	}, nil
}

func (m *modelCapacityCommitment) AdapterForURL(ctx context.Context, url string) (directbase.Adapter, error) {
	// TODO: Support URLs
	return nil, nil
}

type CapacityCommitmentAdapter struct {
	id        *krm.BigQueryReservationCapacityCommitmentIdentity
	gcpClient *gcp.Client
	desired   *krm.BigQueryReservationCapacityCommitment
	actual    *pb.CapacityCommitment
}

var _ directbase.Adapter = &CapacityCommitmentAdapter{}

func (a *CapacityCommitmentAdapter) Find(ctx context.Context) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("getting CapacityCommitment", "name", a.id.String())

	req := &pb.GetCapacityCommitmentRequest{Name: a.id.String()}
	capacityCommitmentpb, err := a.gcpClient.GetCapacityCommitment(ctx, req)
	if err != nil {
		if direct.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("getting CapacityCommitment %q: %w", a.id.String(), err)
	}

	a.actual = capacityCommitmentpb
	return true, nil
}

func parseBool(s *string) bool {
	if s == nil {
		return false
	}
	return *s == "true" || *s == "TRUE" || *s == "True"
}

func (a *CapacityCommitmentAdapter) Create(ctx context.Context, createOp *directbase.CreateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("creating CapacityCommitment", "name", a.id.String())
	mapCtx := &direct.MapContext{}

	desired := a.desired.DeepCopy()
	desiredPb := BigQueryReservationCapacityCommitmentSpec_v1alpha1_ToProto(mapCtx, &desired.Spec)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", a.id.Project, a.id.Location)

	req := &pb.CreateCapacityCommitmentRequest{
		Parent:                          parent,
		CapacityCommitmentId:            a.id.CapacityCommitment,
		CapacityCommitment:              desiredPb,
		EnforceSingleAdminProjectPerOrg: parseBool(desired.Spec.EnforceSingleAdminProjectPerOrg),
	}
	created, err := a.gcpClient.CreateCapacityCommitment(ctx, req)
	if err != nil {
		return fmt.Errorf("creating CapacityCommitment %s: %w", a.id.String(), err)
	}
	log.V(2).Info("successfully created CapacityCommitment", "name", created.Name)

	status := BigQueryReservationCapacityCommitmentStatus_v1alpha1_FromProto(mapCtx, created)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	return createOp.UpdateStatus(ctx, status, nil)
}

func (a *CapacityCommitmentAdapter) Update(ctx context.Context, updateOp *directbase.UpdateOperation) error {
	log := klog.FromContext(ctx)
	log.V(2).Info("updating CapacityCommitment", "name", a.id.String())

	mapCtx := &direct.MapContext{}
	desiredPb := BigQueryReservationCapacityCommitmentSpec_v1alpha1_ToProto(mapCtx, &a.desired.Spec)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	paths := []string{}
	if desiredPb.RenewalPlan != a.actual.RenewalPlan {
		paths = append(paths, "renewal_plan")
	}
	if desiredPb.SlotCount != a.actual.SlotCount {
		paths = append(paths, "slot_count")
	}

	if len(paths) == 0 {
		log.V(2).Info("no field needs update", "name", a.id.String())
		status := BigQueryReservationCapacityCommitmentStatus_v1alpha1_FromProto(mapCtx, a.actual)
		if mapCtx.Err() != nil {
			return mapCtx.Err()
		}
		return updateOp.UpdateStatus(ctx, status, nil)
	}

	report := &structuredreporting.Diff{Object: updateOp.GetUnstructured()}
	for _, path := range paths {
		report.AddField(path, nil, nil)
	}
	structuredreporting.ReportDiff(ctx, report)

	updateMask := &fieldmaskpb.FieldMask{Paths: paths}

	desiredPb.Name = a.id.String()

	req := &pb.UpdateCapacityCommitmentRequest{
		CapacityCommitment: desiredPb,
		UpdateMask:         updateMask,
	}
	updated, err := a.gcpClient.UpdateCapacityCommitment(ctx, req)
	if err != nil {
		return fmt.Errorf("updating CapacityCommitment %s: %w", a.id.String(), err)
	}

	log.V(2).Info("successfully updated CapacityCommitment", "name", updated.Name)

	status := BigQueryReservationCapacityCommitmentStatus_v1alpha1_FromProto(mapCtx, updated)
	if mapCtx.Err() != nil {
		return mapCtx.Err()
	}

	return updateOp.UpdateStatus(ctx, status, nil)
}

func (a *CapacityCommitmentAdapter) Export(ctx context.Context) (*unstructured.Unstructured, error) {
	if a.actual == nil {
		return nil, fmt.Errorf("Find() not called")
	}
	u := &unstructured.Unstructured{}

	obj := &krm.BigQueryReservationCapacityCommitment{}
	mapCtx := &direct.MapContext{}
	obj.Spec = direct.ValueOf(BigQueryReservationCapacityCommitmentSpec_v1alpha1_FromProto(mapCtx, a.actual))
	if mapCtx.Err() != nil {
		return nil, mapCtx.Err()
	}

	obj.Spec.Location = a.id.Location
	obj.Spec.ProjectRef.External = a.id.Project

	uObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}

	u.SetName(a.actual.Name)
	u.SetGroupVersionKind(krm.BigQueryReservationCapacityCommitmentGVK)

	u.Object = uObj
	return u, nil
}

func (a *CapacityCommitmentAdapter) Delete(ctx context.Context, deleteOp *directbase.DeleteOperation) (bool, error) {
	log := klog.FromContext(ctx)
	log.V(2).Info("deleting CapacityCommitment", "name", a.id.String())

	req := &pb.DeleteCapacityCommitmentRequest{
		Name:  a.id.String(),
		Force: true,
	}
	err := a.gcpClient.DeleteCapacityCommitment(ctx, req)
	if err != nil {
		if direct.IsNotFound(err) {
			log.V(2).Info("skipping delete for non-existent CapacityCommitment, assuming it was already deleted", "name", a.id.String())
			return true, nil
		}
		return false, fmt.Errorf("deleting CapacityCommitment %s: %w", a.id.String(), err)
	}
	log.V(2).Info("successfully deleted CapacityCommitment", "name", a.id.String())

	return true, nil
}
