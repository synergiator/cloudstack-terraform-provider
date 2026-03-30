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
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ---------------------------------------------------------------------------
// getUserData
// ---------------------------------------------------------------------------

func TestGetUserData_plainText(t *testing.T) {
	input := "#!/bin/bash\necho hello"
	got, err := getUserData(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := base64.StdEncoding.EncodeToString([]byte(input))
	if got != expected {
		t.Fatalf("expected base64 encoded output %q, got %q", expected, got)
	}
}

func TestGetUserData_alreadyBase64(t *testing.T) {
	original := "#!/bin/bash\necho hello"
	encoded := base64.StdEncoding.EncodeToString([]byte(original))
	got, err := getUserData(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != encoded {
		t.Fatalf("expected passthrough of already-encoded data %q, got %q", encoded, got)
	}
}

func TestGetUserData_emptyString(t *testing.T) {
	got, err := getUserData("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestGetUserData_multilineScript(t *testing.T) {
	input := "#!/bin/bash\napt-get update\napt-get install -y nginx\nsystemctl start nginx"
	got, err := getUserData(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("result should be valid base64: %v", err)
	}
	if string(decoded) != input {
		t.Fatalf("decoded mismatch: expected %q, got %q", input, string(decoded))
	}
}

func TestGetUserData_cloudConfig(t *testing.T) {
	input := "#cloud-config\npackages:\n  - nginx\n  - git"
	got, err := getUserData(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(got)
	if err != nil {
		t.Fatalf("result should be valid base64: %v", err)
	}
	if string(decoded) != input {
		t.Fatalf("decoded mismatch: expected %q, got %q", input, string(decoded))
	}
}

// ---------------------------------------------------------------------------
// resourceCloudStackInstance schema
// ---------------------------------------------------------------------------

func TestResourceInstanceSchema_requiredFields(t *testing.T) {
	r := resourceCloudStackInstance()
	required := []string{"service_offering", "template", "zone"}
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

func TestResourceInstanceSchema_optionalFields(t *testing.T) {
	r := resourceCloudStackInstance()
	optional := []string{
		"name", "display_name", "disk_offering", "network_id", "ip_address",
		"root_disk_size", "group", "project", "keypair", "keypairs",
		"user_data", "userdata_id", "userdata_details", "details",
		"host_id", "cluster_id", "uefi", "boot_mode", "start_vm",
		"expunge", "pod_id", "nicnetworklist", "properties",
		"override_disk_offering",
	}
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

func TestResourceInstanceSchema_conflictsWith(t *testing.T) {
	r := resourceCloudStackInstance()
	cases := []struct {
		field     string
		conflicts string
	}{
		{"affinity_group_ids", "affinity_group_names"},
		{"affinity_group_names", "affinity_group_ids"},
		{"security_group_ids", "security_group_names"},
		{"security_group_names", "security_group_ids"},
		{"keypair", "keypairs"},
		{"keypairs", "keypair"},
	}
	for _, tc := range cases {
		s := r.Schema[tc.field]
		found := false
		for _, c := range s.ConflictsWith {
			if c == tc.conflicts {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("expected %q to conflict with %q", tc.field, tc.conflicts)
		}
	}
}

func TestResourceInstanceSchema_forceNewFields(t *testing.T) {
	r := resourceCloudStackInstance()
	forceNew := []string{
		"disk_offering", "network_id", "ip_address", "template",
		"root_disk_size", "project", "zone", "start_vm", "boot_mode",
		"security_group_ids", "security_group_names", "override_disk_offering",
	}
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

func TestResourceInstanceSchema_hasTags(t *testing.T) {
	r := resourceCloudStackInstance()
	_, ok := r.Schema["tags"]
	if !ok {
		t.Fatal("expected tags in schema")
	}
}

func TestResourceInstanceSchema_hasImporter(t *testing.T) {
	r := resourceCloudStackInstance()
	if r.Importer == nil {
		t.Fatal("expected Importer to be set")
	}
}

func TestResourceInstanceSchema_hasCRUD(t *testing.T) {
	r := resourceCloudStackInstance()
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

func TestResourceInstanceSchema_uefiDefaults(t *testing.T) {
	r := resourceCloudStackInstance()
	uefi := r.Schema["uefi"]
	if uefi.Default != false {
		t.Fatalf("expected uefi default false, got %v", uefi.Default)
	}
}

func TestResourceInstanceSchema_startVmDefaults(t *testing.T) {
	r := resourceCloudStackInstance()
	startVm := r.Schema["start_vm"]
	if startVm.Default != true {
		t.Fatalf("expected start_vm default true, got %v", startVm.Default)
	}
}

func TestResourceInstanceSchema_expungeDefaults(t *testing.T) {
	r := resourceCloudStackInstance()
	expunge := r.Schema["expunge"]
	if expunge.Default != false {
		t.Fatalf("expected expunge default false, got %v", expunge.Default)
	}
}

func TestResourceInstanceSchema_bootModeValidation(t *testing.T) {
	r := resourceCloudStackInstance()
	bootMode := r.Schema["boot_mode"]
	if bootMode.ValidateFunc == nil {
		t.Fatal("expected boot_mode to have ValidateFunc")
	}
	if bootMode.Type != schema.TypeString {
		t.Fatalf("expected boot_mode TypeString, got %v", bootMode.Type)
	}
}

func TestResourceInstanceSchema_userDataStateFunc(t *testing.T) {
	r := resourceCloudStackInstance()
	ud := r.Schema["user_data"]
	if ud.StateFunc == nil {
		t.Fatal("expected user_data to have StateFunc")
	}
	result := ud.StateFunc("test-data")
	if result == "" {
		t.Fatal("StateFunc should return non-empty hash")
	}
	if result == "test-data" {
		t.Fatal("StateFunc should hash the input, not return it unchanged")
	}
}
