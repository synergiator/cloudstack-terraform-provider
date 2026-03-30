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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceFiltersSchema(t *testing.T) {
	s := dataSourceFiltersSchema()
	if s.Type != schema.TypeSet {
		t.Fatalf("expected TypeSet, got %v", s.Type)
	}
	if !s.Required {
		t.Fatal("expected Required to be true")
	}
	if !s.ForceNew {
		t.Fatal("expected ForceNew to be true")
	}
	elem, ok := s.Elem.(*schema.Resource)
	if !ok {
		t.Fatal("expected Elem to be *schema.Resource")
	}
	if _, ok := elem.Schema["name"]; !ok {
		t.Fatal("expected 'name' in nested schema")
	}
	if _, ok := elem.Schema["value"]; !ok {
		t.Fatal("expected 'value' in nested schema")
	}
}

func TestMetadataSchema(t *testing.T) {
	s := metadataSchema()
	if s.Type != schema.TypeMap {
		t.Fatalf("expected TypeMap, got %v", s.Type)
	}
	if !s.Optional {
		t.Fatal("expected Optional to be true")
	}
	if !s.Computed {
		t.Fatal("expected Computed to be true")
	}
}

func TestProviderSchema_hasRequiredAttributes(t *testing.T) {
	p := Provider()
	requiredKeys := []string{"api_url", "api_key", "secret_key", "config", "profile", "http_get_only", "timeout"}
	for _, key := range requiredKeys {
		if _, ok := p.Schema[key]; !ok {
			t.Fatalf("expected provider schema to contain %q", key)
		}
	}
}

func TestProviderSchema_sensitiveFields(t *testing.T) {
	p := Provider()
	sensitiveKeys := []string{"api_key", "secret_key"}
	for _, key := range sensitiveKeys {
		s, ok := p.Schema[key]
		if !ok {
			t.Fatalf("expected provider schema to contain %q", key)
		}
		if !s.Sensitive {
			t.Fatalf("expected %q to be marked Sensitive", key)
		}
	}
}

func TestProviderSchema_conflictsWith(t *testing.T) {
	p := Provider()

	apiUrlSchema := p.Schema["api_url"]
	if len(apiUrlSchema.ConflictsWith) != 2 {
		t.Fatalf("expected api_url to conflict with 2 fields, got %d", len(apiUrlSchema.ConflictsWith))
	}

	configSchema := p.Schema["config"]
	if len(configSchema.ConflictsWith) != 3 {
		t.Fatalf("expected config to conflict with 3 fields, got %d", len(configSchema.ConflictsWith))
	}
}

func TestProviderSchema_registersResources(t *testing.T) {
	p := Provider()
	sampleResources := []string{
		"cloudstack_vpc", "cloudstack_instance", "cloudstack_network",
		"cloudstack_firewall", "cloudstack_template", "cloudstack_zone",
	}
	for _, name := range sampleResources {
		if _, ok := p.ResourcesMap[name]; !ok {
			t.Fatalf("expected resource %q in ResourcesMap", name)
		}
	}
}

func TestProviderSchema_registersDataSources(t *testing.T) {
	p := Provider()
	sampleDataSources := []string{
		"cloudstack_vpc", "cloudstack_zone", "cloudstack_template",
		"cloudstack_instance", "cloudstack_ssh_keypair",
	}
	for _, name := range sampleDataSources {
		if _, ok := p.DataSourcesMap[name]; !ok {
			t.Fatalf("expected data source %q in DataSourcesMap", name)
		}
	}
}
