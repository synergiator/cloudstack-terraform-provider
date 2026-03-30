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
// resourceCloudStackVPNGateway schema
// ---------------------------------------------------------------------------

func TestResourceVPNGatewaySchema_requiredFields(t *testing.T) {
	r := resourceCloudStackVPNGateway()
	s, ok := r.Schema["vpc_id"]
	if !ok {
		t.Fatal("expected schema key vpc_id")
	}
	if !s.Required {
		t.Fatal("expected vpc_id to be Required")
	}
	if !s.ForceNew {
		t.Fatal("expected vpc_id to have ForceNew")
	}
}

func TestResourceVPNGatewaySchema_computedFields(t *testing.T) {
	r := resourceCloudStackVPNGateway()
	s, ok := r.Schema["public_ip"]
	if !ok {
		t.Fatal("expected schema key public_ip")
	}
	if !s.Computed {
		t.Fatal("expected public_ip to be Computed")
	}
	if s.Type != schema.TypeString {
		t.Fatalf("expected public_ip TypeString, got %v", s.Type)
	}
}

func TestResourceVPNGatewaySchema_hasCRD(t *testing.T) {
	r := resourceCloudStackVPNGateway()
	if r.Create == nil {
		t.Fatal("expected Create to be set")
	}
	if r.Read == nil {
		t.Fatal("expected Read to be set")
	}
	if r.Delete == nil {
		t.Fatal("expected Delete to be set")
	}
	if r.Update != nil {
		t.Fatal("expected Update to be nil (no updatable fields)")
	}
}

func TestResourceVPNGatewaySchema_hasImporter(t *testing.T) {
	r := resourceCloudStackVPNGateway()
	if r.Importer == nil {
		t.Fatal("expected Importer to be set")
	}
}

// ---------------------------------------------------------------------------
// resourceCloudStackVPNConnection schema
// ---------------------------------------------------------------------------

func TestResourceVPNConnectionSchema_requiredFields(t *testing.T) {
	r := resourceCloudStackVPNConnection()
	required := []string{"customer_gateway_id", "vpn_gateway_id"}
	for _, key := range required {
		s, ok := r.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.Required {
			t.Fatalf("expected %q to be Required", key)
		}
		if !s.ForceNew {
			t.Fatalf("expected %q to have ForceNew", key)
		}
	}
}

func TestResourceVPNConnectionSchema_hasCRD(t *testing.T) {
	r := resourceCloudStackVPNConnection()
	if r.Create == nil {
		t.Fatal("expected Create to be set")
	}
	if r.Read == nil {
		t.Fatal("expected Read to be set")
	}
	if r.Delete == nil {
		t.Fatal("expected Delete to be set")
	}
	if r.Update != nil {
		t.Fatal("expected Update to be nil (no updatable fields)")
	}
}

func TestResourceVPNConnectionSchema_noImporter(t *testing.T) {
	r := resourceCloudStackVPNConnection()
	if r.Importer != nil {
		t.Fatal("expected Importer to be nil (not importable)")
	}
}

// ---------------------------------------------------------------------------
// resourceCloudStackVPNCustomerGateway schema
// ---------------------------------------------------------------------------

func TestResourceVPNCustomerGatewaySchema_requiredFields(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
	required := []string{"name", "cidr", "esp_policy", "gateway", "ike_policy", "ipsec_psk"}
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

func TestResourceVPNCustomerGatewaySchema_optionalFields(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
	optional := []string{"dpd", "esp_lifetime", "ike_lifetime", "project"}
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

func TestResourceVPNCustomerGatewaySchema_hasCRUD(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
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

func TestResourceVPNCustomerGatewaySchema_hasImporter(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
	if r.Importer == nil {
		t.Fatal("expected Importer to be set")
	}
}

func TestResourceVPNCustomerGatewaySchema_projectForceNew(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
	s := r.Schema["project"]
	if !s.ForceNew {
		t.Fatal("expected project to have ForceNew")
	}
	if !s.Computed {
		t.Fatal("expected project to be Computed")
	}
}

func TestResourceVPNCustomerGatewaySchema_lifetimeTypes(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
	for _, key := range []string{"esp_lifetime", "ike_lifetime"} {
		s := r.Schema[key]
		if s.Type != schema.TypeInt {
			t.Fatalf("expected %q to be TypeInt, got %v", key, s.Type)
		}
		if !s.Computed {
			t.Fatalf("expected %q to be Computed", key)
		}
	}
}

func TestResourceVPNCustomerGatewaySchema_dpdType(t *testing.T) {
	r := resourceCloudStackVPNCustomerGateway()
	s := r.Schema["dpd"]
	if s.Type != schema.TypeBool {
		t.Fatalf("expected dpd to be TypeBool, got %v", s.Type)
	}
	if !s.Computed {
		t.Fatalf("expected dpd to be Computed")
	}
}
