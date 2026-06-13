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

package v1beta1

import (
	"context"
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/k8s-config-connector/apis/common"
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FirewallPolicyAssociationIdentity struct {
	firewallPolicy string
	name           string
}

func (i *FirewallPolicyAssociationIdentity) String() string {
	return "locations/global/firewallPolicies/" + i.firewallPolicy + "/associations/" + i.name
}

func (i *FirewallPolicyAssociationIdentity) FirewallPolicy() string {
	return i.firewallPolicy
}

func (i *FirewallPolicyAssociationIdentity) Name() string {
	return i.name
}

func NewFirewallPolicyAssociationIdentity(ctx context.Context, reader client.Reader, obj *ComputeFirewallPolicyAssociation) (*FirewallPolicyAssociationIdentity, error) {
	// Resolve firewallPolicyRef
	firewallPolicyRef, err := refsv1beta1.ResolveComputeFirewallPolicy(ctx, reader, obj, &obj.Spec.FirewallPolicyRef)
	if err != nil {
		return nil, err
	}
	firewallPolicy := firewallPolicyRef.External
	if firewallPolicy == "" {
		return nil, fmt.Errorf("cannot resolve firewallPolicyRef")
	}

	// We only want the short ID of the firewall policy (e.g., if external is locations/global/firewallPolicies/12345, get 12345)
	if strings.HasPrefix(firewallPolicy, "locations/global/firewallPolicies/") {
		firewallPolicy = strings.TrimPrefix(firewallPolicy, "locations/global/firewallPolicies/")
	}

	// Get name
	name := common.ValueOf(obj.Spec.ResourceID)
	if name == "" {
		name = obj.GetName()
	}

	externalRef := common.ValueOf(obj.Status.ExternalRef)
	if externalRef != "" {
		actualIdentity, err := parseFirewallPolicyAssociationExternal(externalRef)
		if err != nil {
			return nil, err
		}
		if actualIdentity.firewallPolicy != firewallPolicy {
			return nil, fmt.Errorf("spec.firewallPolicyRef changed, expect %s, got %s", actualIdentity.firewallPolicy, firewallPolicy)
		}
		if actualIdentity.name != name {
			return nil, fmt.Errorf("cannot reset `spec.resourceID` to %s, since it has already assigned to %s",
				name, actualIdentity.name)
		}
	}

	return &FirewallPolicyAssociationIdentity{
		firewallPolicy: firewallPolicy,
		name:           name,
	}, nil
}

func parseFirewallPolicyAssociationExternal(externalRef string) (*FirewallPolicyAssociationIdentity, error) {
	// Format is: locations/global/firewallPolicies/{firewallPolicy}/associations/{name}
	tokens := strings.Split(externalRef, "/")
	if len(tokens) == 6 && tokens[0] == "locations" && tokens[1] == "global" && tokens[2] == "firewallPolicies" && tokens[4] == "associations" {
		return &FirewallPolicyAssociationIdentity{
			firewallPolicy: tokens[3],
			name:           tokens[5],
		}, nil
	}
	return nil, fmt.Errorf("format of ComputeFirewallPolicyAssociation externalRef %q is not valid", externalRef)
}
