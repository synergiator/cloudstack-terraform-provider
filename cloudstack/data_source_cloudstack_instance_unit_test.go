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
)

// ---------------------------------------------------------------------------
// latestInstance
// ---------------------------------------------------------------------------

func TestLatestInstance_single(t *testing.T) {
	instances := []*cloudstack.VirtualMachine{
		{Id: "vm-1", Name: "first", Created: "2025-01-15T10:00:00+0000"},
	}
	got, err := latestInstance(instances)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vm-1" {
		t.Fatalf("expected vm-1, got %s", got.Id)
	}
}

func TestLatestInstance_picksMostRecent(t *testing.T) {
	instances := []*cloudstack.VirtualMachine{
		{Id: "vm-old", Name: "old", Created: "2024-01-01T00:00:00+0000"},
		{Id: "vm-new", Name: "new", Created: "2025-06-15T12:00:00+0000"},
		{Id: "vm-mid", Name: "mid", Created: "2025-03-10T08:00:00+0000"},
	}
	got, err := latestInstance(instances)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vm-new" {
		t.Fatalf("expected vm-new, got %s", got.Id)
	}
}

func TestLatestInstance_badDateFormat(t *testing.T) {
	instances := []*cloudstack.VirtualMachine{
		{Id: "vm-bad", Name: "bad", Created: "not-a-date"},
	}
	_, err := latestInstance(instances)
	if err == nil {
		t.Fatal("expected error for bad date format")
	}
}

func TestLatestInstance_sameTimestamp(t *testing.T) {
	instances := []*cloudstack.VirtualMachine{
		{Id: "vm-a", Name: "a", Created: "2025-01-01T00:00:00+0000"},
		{Id: "vm-b", Name: "b", Created: "2025-01-01T00:00:00+0000"},
	}
	got, err := latestInstance(instances)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vm-a" {
		t.Fatalf("expected vm-a (first seen), got %s", got.Id)
	}
}

// ---------------------------------------------------------------------------
// applyInstanceFilters
// ---------------------------------------------------------------------------

func TestApplyInstanceFilters_matchByName(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:   "vm-1",
		Name: "web-server-01",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "web-server-01",
	})
	match, err := applyInstanceFilters(vm, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match")
	}
}

func TestApplyInstanceFilters_noMatch(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:   "vm-1",
		Name: "web-server-01",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "db-server",
	})
	match, err := applyInstanceFilters(vm, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected no match")
	}
}

func TestApplyInstanceFilters_regexMatch(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:    "vm-1",
		Name:  "web-server-01",
		State: "Running",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "state",
		"value": "^Run.*",
	})
	match, err := applyInstanceFilters(vm, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected regex match")
	}
}

func TestApplyInstanceFilters_invalidRegex(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:   "vm-1",
		Name: "test",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "[invalid",
	})
	_, err := applyInstanceFilters(vm, filters)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApplyInstanceFilters_underscoreFieldMapping(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:          "vm-1",
		Name:        "test",
		Displayname: "My Display Name",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "display_name",
		"value": "My Display Name",
	})
	match, err := applyInstanceFilters(vm, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match on display_name -> displayname")
	}
}

func TestApplyInstanceFilters_emptyFilters(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:   "vm-1",
		Name: "test",
	}
	filters := newFilterSet()
	match, err := applyInstanceFilters(vm, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match with no filters")
	}
}

func TestApplyInstanceFilters_multipleFilters(t *testing.T) {
	vm := &cloudstack.VirtualMachine{
		Id:    "vm-1",
		Name:  "web-01",
		State: "Running",
	}
	filters := newFilterSet(
		map[string]interface{}{"name": "name", "value": "web-01"},
		map[string]interface{}{"name": "state", "value": "Running"},
	)
	match, err := applyInstanceFilters(vm, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match on all filters")
	}
}

// ---------------------------------------------------------------------------
// instanceDescriptionAttributes
// ---------------------------------------------------------------------------

func TestInstanceDescriptionAttributes(t *testing.T) {
	raw := dataSourceCloudstackInstance()
	d := raw.TestResourceData()

	vm := &cloudstack.VirtualMachine{
		Id:          "vm-123",
		Account:     "admin",
		Created:     "2025-01-15T10:00:00+0000",
		Displayname: "My VM",
		State:       "Running",
		Hostid:      "host-1",
		Zoneid:      "zone-1",
		Tags: []cloudstack.Tags{
			{Key: "env", Value: "staging"},
		},
		Nic: []cloudstack.Nic{
			{Ipaddress: "10.0.0.5"},
		},
	}

	err := instanceDescriptionAttributes(d, vm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Id() != "vm-123" {
		t.Fatalf("expected id vm-123, got %s", d.Id())
	}
	if d.Get("instance_id") != "vm-123" {
		t.Fatalf("expected instance_id vm-123, got %s", d.Get("instance_id"))
	}
	if d.Get("account") != "admin" {
		t.Fatalf("expected account admin, got %s", d.Get("account"))
	}
	if d.Get("display_name") != "My VM" {
		t.Fatalf("expected display_name 'My VM', got %s", d.Get("display_name"))
	}
	if d.Get("state") != "Running" {
		t.Fatalf("expected state Running, got %s", d.Get("state"))
	}
	if d.Get("host_id") != "host-1" {
		t.Fatalf("expected host_id host-1, got %s", d.Get("host_id"))
	}
	if d.Get("zone_id") != "zone-1" {
		t.Fatalf("expected zone_id zone-1, got %s", d.Get("zone_id"))
	}
}

func TestInstanceDescriptionAttributes_emptyTags(t *testing.T) {
	raw := dataSourceCloudstackInstance()
	d := raw.TestResourceData()

	vm := &cloudstack.VirtualMachine{
		Id:   "vm-empty",
		Tags: []cloudstack.Tags{},
		Nic: []cloudstack.Nic{
			{Ipaddress: "10.0.0.1"},
		},
	}

	err := instanceDescriptionAttributes(d, vm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := d.Get("tags").(map[string]interface{})
	if len(tags) != 0 {
		t.Fatalf("expected empty tags, got %v", tags)
	}
}

// ---------------------------------------------------------------------------
// dataSourceCloudstackInstance schema
// ---------------------------------------------------------------------------

func TestDataSourceInstanceSchema(t *testing.T) {
	ds := dataSourceCloudstackInstance()
	if ds.Read == nil {
		t.Fatal("expected Read function to be set")
	}
	expectedComputed := []string{"instance_id", "account", "display_name", "state", "host_id", "zone_id", "created"}
	for _, key := range expectedComputed {
		s, ok := ds.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.Computed {
			t.Fatalf("expected %q to be Computed", key)
		}
	}
	if _, ok := ds.Schema["filter"]; !ok {
		t.Fatal("expected filter in schema")
	}
	if _, ok := ds.Schema["tags"]; !ok {
		t.Fatal("expected tags in schema")
	}
	nicSchema, ok := ds.Schema["nic"]
	if !ok {
		t.Fatal("expected nic in schema")
	}
	if nicSchema.Type != 5 {
		t.Fatal("expected nic to be TypeList")
	}
}
