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

// ---------------------------------------------------------------------------
// resourceCloudStackVPC schema validation
// ---------------------------------------------------------------------------

func TestResourceVPCSchema_requiredFields(t *testing.T) {
	r := resourceCloudStackVPC()
	required := []string{"name", "cidr", "vpc_offering", "zone"}
	for _, key := range required {
		s, ok := r.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.Required {
			t.Fatalf("expected %q to be Required", key)
		}
	}
}

func TestResourceVPCSchema_optionalFields(t *testing.T) {
	r := resourceCloudStackVPC()
	optional := []string{"display_text", "network_domain", "project"}
	for _, key := range optional {
		s, ok := r.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.Optional {
			t.Fatalf("expected %q to be Optional", key)
		}
	}
}

func TestResourceVPCSchema_computedFields(t *testing.T) {
	r := resourceCloudStackVPC()
	computed := []string{"display_text", "network_domain", "project", "source_nat_ip"}
	for _, key := range computed {
		s, ok := r.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.Computed {
			t.Fatalf("expected %q to be Computed", key)
		}
	}
}

func TestResourceVPCSchema_forceNewFields(t *testing.T) {
	r := resourceCloudStackVPC()
	forceNew := []string{"cidr", "vpc_offering", "network_domain", "project", "zone"}
	for _, key := range forceNew {
		s, ok := r.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.ForceNew {
			t.Fatalf("expected %q to have ForceNew", key)
		}
	}
}

func TestResourceVPCSchema_hasTags(t *testing.T) {
	r := resourceCloudStackVPC()
	_, ok := r.Schema["tags"]
	if !ok {
		t.Fatal("expected tags in schema")
	}
}

func TestResourceVPCSchema_hasImporter(t *testing.T) {
	r := resourceCloudStackVPC()
	if r.Importer == nil {
		t.Fatal("expected Importer to be set")
	}
}

func TestResourceVPCSchema_hasCRUD(t *testing.T) {
	r := resourceCloudStackVPC()
	if r.Create == nil {
		t.Fatal("expected Create to be set")
	}
	if r.Read == nil {
		t.Fatal("expected Read to be set")
	}
	if r.Update == nil {
		t.Fatal("expected Update to be set")
	}
	if r.Delete == nil {
		t.Fatal("expected Delete to be set")
	}
}

func TestResourceVPC_sourceNatIPIsReadOnly(t *testing.T) {
	r := resourceCloudStackVPC()
	s := r.Schema["source_nat_ip"]
	if s.Type != schema.TypeString {
		t.Fatalf("expected source_nat_ip to be TypeString, got %v", s.Type)
	}
	if !s.Computed {
		t.Fatal("expected source_nat_ip to be Computed")
	}
	if s.Optional || s.Required {
		t.Fatal("expected source_nat_ip to be read-only (not Optional or Required)")
	}
}
