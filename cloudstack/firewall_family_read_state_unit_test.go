//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package cloudstack

import (
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.uber.org/mock/gomock"
)

func stringSet(values ...string) *schema.Set {
	s := &schema.Set{F: schema.HashString}
	for _, v := range values {
		s.Add(v)
	}
	return s
}

func TestFirewallRead_clearsIDWhenListReturnsNoMatchingRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := cloudstack.NewMockClient(ctrl)
	firewall := cs.Firewall.(*cloudstack.MockFirewallServiceIface)

	firewall.EXPECT().NewListFirewallRulesParams().Return(&cloudstack.ListFirewallRulesParams{})
	firewall.EXPECT().ListFirewallRules(gomock.Any()).DoAndReturn(func(p *cloudstack.ListFirewallRulesParams) (*cloudstack.ListFirewallRulesResponse, error) {
		if got, _ := p.GetIpaddressid(); got != "ip-1" {
			t.Fatalf("expected ipaddressid=ip-1, got %q", got)
		}
		return &cloudstack.ListFirewallRulesResponse{
			Count:         0,
			FirewallRules: []*cloudstack.FirewallRule{},
		}, nil
	})

	r := resourceCloudStackFirewall()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ip_address_id": "ip-1",
		"managed":       false,
	})
	d.SetId("ip-1")

	rules := r.Schema["rule"].ZeroValue().(*schema.Set)
	rules.Add(map[string]interface{}{
		"protocol":  "tcp",
		"cidr_list": stringSet("0.0.0.0/0"),
		"ports":     stringSet("80"),
		"uuids":     map[string]interface{}{"80": "fw-rule-1"},
	})
	if err := d.Set("rule", rules); err != nil {
		t.Fatalf("failed to seed rule set: %v", err)
	}

	if err := resourceCloudStackFirewallRead(d, cs); err != nil {
		t.Fatalf("read returned error: %v", err)
	}

	if d.Id() != "" {
		t.Fatalf("expected ID to be cleared when list returns no matching rules, got %q", d.Id())
	}
}

func TestEgressFirewallRead_clearsIDWhenListReturnsNoMatchingRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := cloudstack.NewMockClient(ctrl)
	firewall := cs.Firewall.(*cloudstack.MockFirewallServiceIface)

	firewall.EXPECT().NewListEgressFirewallRulesParams().Return(&cloudstack.ListEgressFirewallRulesParams{})
	firewall.EXPECT().ListEgressFirewallRules(gomock.Any()).DoAndReturn(func(p *cloudstack.ListEgressFirewallRulesParams) (*cloudstack.ListEgressFirewallRulesResponse, error) {
		if got, _ := p.GetNetworkid(); got != "net-1" {
			t.Fatalf("expected networkid=net-1, got %q", got)
		}
		return &cloudstack.ListEgressFirewallRulesResponse{
			Count:               0,
			EgressFirewallRules: []*cloudstack.EgressFirewallRule{},
		}, nil
	})

	r := resourceCloudStackEgressFirewall()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"network_id": "net-1",
		"managed":    false,
	})
	d.SetId("net-1")

	rules := r.Schema["rule"].ZeroValue().(*schema.Set)
	rules.Add(map[string]interface{}{
		"protocol":  "tcp",
		"cidr_list": stringSet("0.0.0.0/0"),
		"ports":     stringSet("443"),
		"uuids":     map[string]interface{}{"443": "eg-rule-1"},
	})
	if err := d.Set("rule", rules); err != nil {
		t.Fatalf("failed to seed rule set: %v", err)
	}

	if err := resourceCloudStackEgressFirewallRead(d, cs); err != nil {
		t.Fatalf("read returned error: %v", err)
	}

	if d.Id() != "" {
		t.Fatalf("expected ID to be cleared when list returns no matching rules, got %q", d.Id())
	}
}

func TestPortForwardRead_clearsIDEvenWhenIPAddressExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := cloudstack.NewMockClient(ctrl)
	address := cs.Address.(*cloudstack.MockAddressServiceIface)
	firewall := cs.Firewall.(*cloudstack.MockFirewallServiceIface)

	address.EXPECT().GetPublicIpAddressByID("ip-1", gomock.Any()).Return(&cloudstack.PublicIpAddress{Id: "ip-1"}, 1, nil)
	firewall.EXPECT().NewListPortForwardingRulesParams().Return(&cloudstack.ListPortForwardingRulesParams{})
	firewall.EXPECT().ListPortForwardingRules(gomock.Any()).DoAndReturn(func(p *cloudstack.ListPortForwardingRulesParams) (*cloudstack.ListPortForwardingRulesResponse, error) {
		if got, _ := p.GetIpaddressid(); got != "ip-1" {
			t.Fatalf("expected ipaddressid=ip-1, got %q", got)
		}
		return &cloudstack.ListPortForwardingRulesResponse{
			Count:               0,
			PortForwardingRules: []*cloudstack.PortForwardingRule{},
		}, nil
	})

	r := resourceCloudStackPortForward()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ip_address_id": "ip-1",
		"managed":       false,
	})
	d.SetId("ip-1")

	forwards := r.Schema["forward"].ZeroValue().(*schema.Set)
	forwards.Add(map[string]interface{}{
		"protocol":           "tcp",
		"private_port":       8080,
		"public_port":        8080,
		"virtual_machine_id": "vm-1",
		"uuid":               "pf-rule-1",
	})
	if err := d.Set("forward", forwards); err != nil {
		t.Fatalf("failed to seed forward set: %v", err)
	}

	if err := resourceCloudStackPortForwardRead(d, cs); err != nil {
		t.Fatalf("read returned error: %v", err)
	}

	if d.Id() != "" {
		t.Fatalf("expected ID to be cleared when forward list is empty, got %q", d.Id())
	}
}

func TestFirewallRead_preservesIDWhenListReturnsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := cloudstack.NewMockClient(ctrl)
	firewall := cs.Firewall.(*cloudstack.MockFirewallServiceIface)

	firewall.EXPECT().NewListFirewallRulesParams().Return(&cloudstack.ListFirewallRulesParams{})
	firewall.EXPECT().ListFirewallRules(gomock.Any()).Return(&cloudstack.ListFirewallRulesResponse{
		Count:         0,
		FirewallRules: []*cloudstack.FirewallRule{},
	}, nil)

	r := resourceCloudStackFirewall()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ip_address_id": "ip-1",
		"managed":       false,
	})
	d.SetId("ip-1")

	rules := r.Schema["rule"].ZeroValue().(*schema.Set)
	rules.Add(map[string]interface{}{
		"protocol":  "tcp",
		"cidr_list": stringSet("0.0.0.0/0"),
		"ports":     stringSet("80"),
		"uuids":     map[string]interface{}{"80": "fw-rule-1"},
	})
	if err := d.Set("rule", rules); err != nil {
		t.Fatalf("failed to seed rule set: %v", err)
	}

	if err := resourceCloudStackFirewallRead(d, cs); err != nil {
		t.Fatalf("read returned error: %v", err)
	}

	if d.Id() != "ip-1" {
		t.Fatalf("expected ID to stay set when parent still exists, got %q", d.Id())
	}
}

func TestEgressFirewallRead_preservesIDWhenListReturnsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := cloudstack.NewMockClient(ctrl)
	firewall := cs.Firewall.(*cloudstack.MockFirewallServiceIface)

	firewall.EXPECT().NewListEgressFirewallRulesParams().Return(&cloudstack.ListEgressFirewallRulesParams{})
	firewall.EXPECT().ListEgressFirewallRules(gomock.Any()).Return(&cloudstack.ListEgressFirewallRulesResponse{
		Count:               0,
		EgressFirewallRules: []*cloudstack.EgressFirewallRule{},
	}, nil)

	r := resourceCloudStackEgressFirewall()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"network_id": "net-1",
		"managed":    false,
	})
	d.SetId("net-1")

	rules := r.Schema["rule"].ZeroValue().(*schema.Set)
	rules.Add(map[string]interface{}{
		"protocol":  "tcp",
		"cidr_list": stringSet("0.0.0.0/0"),
		"ports":     stringSet("443"),
		"uuids":     map[string]interface{}{"443": "eg-rule-1"},
	})
	if err := d.Set("rule", rules); err != nil {
		t.Fatalf("failed to seed rule set: %v", err)
	}

	if err := resourceCloudStackEgressFirewallRead(d, cs); err != nil {
		t.Fatalf("read returned error: %v", err)
	}

	if d.Id() != "net-1" {
		t.Fatalf("expected ID to stay set when parent still exists, got %q", d.Id())
	}
}

func TestPortForwardRead_preservesIDWhenParentIPAddressExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := cloudstack.NewMockClient(ctrl)
	address := cs.Address.(*cloudstack.MockAddressServiceIface)
	firewall := cs.Firewall.(*cloudstack.MockFirewallServiceIface)

	address.EXPECT().GetPublicIpAddressByID("ip-1", gomock.Any()).Return(&cloudstack.PublicIpAddress{Id: "ip-1"}, 1, nil)
	firewall.EXPECT().NewListPortForwardingRulesParams().Return(&cloudstack.ListPortForwardingRulesParams{})
	firewall.EXPECT().ListPortForwardingRules(gomock.Any()).Return(&cloudstack.ListPortForwardingRulesResponse{
		Count:               0,
		PortForwardingRules: []*cloudstack.PortForwardingRule{},
	}, nil)

	r := resourceCloudStackPortForward()
	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"ip_address_id": "ip-1",
		"managed":       false,
	})
	d.SetId("ip-1")

	forwards := r.Schema["forward"].ZeroValue().(*schema.Set)
	forwards.Add(map[string]interface{}{
		"protocol":           "tcp",
		"private_port":       8080,
		"public_port":        8080,
		"virtual_machine_id": "vm-1",
		"uuid":               "pf-rule-1",
	})
	if err := d.Set("forward", forwards); err != nil {
		t.Fatalf("failed to seed forward set: %v", err)
	}

	if err := resourceCloudStackPortForwardRead(d, cs); err != nil {
		t.Fatalf("read returned error: %v", err)
	}

	if d.Id() != "ip-1" {
		t.Fatalf("expected ID to stay set when parent IP exists, got %q", d.Id())
	}
}
