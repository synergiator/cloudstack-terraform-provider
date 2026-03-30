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

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestServiceOfferingMergeCommonSchema_includesCommonFields(t *testing.T) {
	result := serviceOfferingMergeCommonSchema(map[string]schema.Attribute{})
	commonKeys := []string{
		"name", "display_text", "id", "offer_ha", "is_volatile",
		"limit_cpu_use", "dynamic_scaling_enabled", "deployment_planner",
		"host_tags", "network_rate", "zone_ids", "domain_ids",
		"disk_offering_id", "disk_offering", "disk_hypervisor", "disk_storage",
	}
	for _, key := range commonKeys {
		if _, ok := result[key]; !ok {
			t.Fatalf("expected common key %q in merged schema", key)
		}
	}
}

func TestServiceOfferingMergeCommonSchema_includesCustomFields(t *testing.T) {
	custom := map[string]schema.Attribute{
		"cpu_speed": schema.Int32Attribute{
			Description: "CPU speed",
			Required:    true,
		},
	}
	result := serviceOfferingMergeCommonSchema(custom)
	if _, ok := result["cpu_speed"]; !ok {
		t.Fatal("expected custom key 'cpu_speed' in merged schema")
	}
	if _, ok := result["name"]; !ok {
		t.Fatal("expected common key 'name' to still be present")
	}
}

func TestServiceOfferingMergeCommonSchema_customOverridesCommon(t *testing.T) {
	custom := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description: "overridden",
			Optional:    true,
		},
	}
	result := serviceOfferingMergeCommonSchema(custom)
	attr, ok := result["name"]
	if !ok {
		t.Fatal("expected key 'name' in merged schema")
	}
	strAttr, ok := attr.(schema.StringAttribute)
	if !ok {
		t.Fatal("expected StringAttribute for 'name'")
	}
	if strAttr.Description != "overridden" {
		t.Fatalf("expected overridden description, got %q", strAttr.Description)
	}
}

func TestServiceOfferingMergeCommonSchema_emptyCustom(t *testing.T) {
	result := serviceOfferingMergeCommonSchema(map[string]schema.Attribute{})
	if _, ok := result["name"]; !ok {
		t.Fatal("expected common key 'name'")
	}
	if _, ok := result["display_text"]; !ok {
		t.Fatal("expected common key 'display_text'")
	}
}
