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
)

// ---------------------------------------------------------------------------
// latestVPC
// ---------------------------------------------------------------------------

func TestLatestVPC_single(t *testing.T) {
	vpcs := []*cloudstack.VPC{
		{Id: "vpc-1", Name: "first", Created: "2025-01-15T10:00:00+0000"},
	}
	got, err := latestVPC(vpcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vpc-1" {
		t.Fatalf("expected vpc-1, got %s", got.Id)
	}
}

func TestLatestVPC_picksMostRecent(t *testing.T) {
	vpcs := []*cloudstack.VPC{
		{Id: "vpc-old", Name: "old", Created: "2024-01-01T00:00:00+0000"},
		{Id: "vpc-new", Name: "new", Created: "2025-06-15T12:00:00+0000"},
		{Id: "vpc-mid", Name: "mid", Created: "2025-03-10T08:00:00+0000"},
	}
	got, err := latestVPC(vpcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vpc-new" {
		t.Fatalf("expected vpc-new, got %s", got.Id)
	}
}

func TestLatestVPC_badDateFormat(t *testing.T) {
	vpcs := []*cloudstack.VPC{
		{Id: "vpc-bad", Name: "bad", Created: "not-a-date"},
	}
	_, err := latestVPC(vpcs)
	if err == nil {
		t.Fatal("expected error for bad date format")
	}
}

func TestLatestVPC_sameTimestamp(t *testing.T) {
	vpcs := []*cloudstack.VPC{
		{Id: "vpc-a", Name: "a", Created: "2025-01-01T00:00:00+0000"},
		{Id: "vpc-b", Name: "b", Created: "2025-01-01T00:00:00+0000"},
	}
	got, err := latestVPC(vpcs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vpc-a" {
		t.Fatalf("expected vpc-a (first seen), got %s", got.Id)
	}
}

// ---------------------------------------------------------------------------
// applyVPCFilters
// ---------------------------------------------------------------------------

func newFilterSet(filters ...map[string]interface{}) *schema.Set {
	s := dataSourceFiltersSchema()
	filterSchema := s.Elem.(*schema.Resource)
	set := schema.NewSet(schema.HashResource(filterSchema), nil)
	for _, f := range filters {
		set.Add(f)
	}
	return set
}

func TestApplyVPCFilters_matchByName(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:   "vpc-1",
		Name: "my-vpc",
		Cidr: "10.0.0.0/16",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "my-vpc",
	})
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match")
	}
}

func TestApplyVPCFilters_noMatch(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:   "vpc-1",
		Name: "my-vpc",
		Cidr: "10.0.0.0/16",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "other-vpc",
	})
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected no match")
	}
}

func TestApplyVPCFilters_regexMatch(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:   "vpc-1",
		Name: "production-vpc-01",
		Cidr: "10.0.0.0/16",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "^production-.*",
	})
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected regex match")
	}
}

func TestApplyVPCFilters_invalidRegex(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:   "vpc-1",
		Name: "my-vpc",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "name",
		"value": "[invalid",
	})
	_, err := applyVPCFilters(vpc, filters)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApplyVPCFilters_underscoreFieldMapping(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:              "vpc-1",
		Name:            "test",
		Displaytext:     "Test VPC",
		Vpcofferingname: "Default VPC offering",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "display_text",
		"value": "Test VPC",
	})
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match on display_text -> displaytext")
	}
}

func TestApplyVPCFilters_multipleFilters(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:       "vpc-1",
		Name:     "my-vpc",
		Cidr:     "10.0.0.0/16",
		Zonename: "zone-1",
	}
	filters := newFilterSet(
		map[string]interface{}{"name": "name", "value": "my-vpc"},
		map[string]interface{}{"name": "zone_name", "value": "zone-1"},
	)
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match on all filters")
	}
}

func TestApplyVPCFilters_multipleFiltersPartialMismatch(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:       "vpc-1",
		Name:     "my-vpc",
		Cidr:     "10.0.0.0/16",
		Zonename: "zone-1",
	}
	filters := newFilterSet(
		map[string]interface{}{"name": "name", "value": "my-vpc"},
		map[string]interface{}{"name": "zone_name", "value": "zone-2"},
	)
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected no match when one filter doesn't match")
	}
}

func TestApplyVPCFilters_emptyFilters(t *testing.T) {
	vpc := &cloudstack.VPC{
		Id:   "vpc-1",
		Name: "my-vpc",
	}
	filters := newFilterSet()
	match, err := applyVPCFilters(vpc, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match with no filters")
	}
}

// ---------------------------------------------------------------------------
// vpcDescriptionAttributes
// ---------------------------------------------------------------------------

func TestVpcDescriptionAttributes(t *testing.T) {
	raw := dataSourceCloudstackVPC()
	d := raw.TestResourceData()

	vpc := &cloudstack.VPC{
		Id:              "vpc-123",
		Name:            "test-vpc",
		Displaytext:     "Test VPC Description",
		Cidr:            "10.0.0.0/16",
		Vpcofferingname: "Default VPC offering",
		Networkdomain:   "cs.internal",
		Project:         "terraform",
		Zonename:        "zone-1",
		Tags: []cloudstack.Tags{
			{Key: "env", Value: "prod"},
		},
	}

	err := vpcDescriptionAttributes(d, vpc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Id() != "vpc-123" {
		t.Fatalf("expected id vpc-123, got %s", d.Id())
	}
	if d.Get("name") != "test-vpc" {
		t.Fatalf("expected name test-vpc, got %s", d.Get("name"))
	}
	if d.Get("display_text") != "Test VPC Description" {
		t.Fatalf("expected display_text 'Test VPC Description', got %s", d.Get("display_text"))
	}
	if d.Get("cidr") != "10.0.0.0/16" {
		t.Fatalf("expected cidr 10.0.0.0/16, got %s", d.Get("cidr"))
	}
	if d.Get("vpc_offering_name") != "Default VPC offering" {
		t.Fatalf("expected vpc_offering_name, got %s", d.Get("vpc_offering_name"))
	}
	if d.Get("network_domain") != "cs.internal" {
		t.Fatalf("expected network_domain cs.internal, got %s", d.Get("network_domain"))
	}
	if d.Get("project") != "terraform" {
		t.Fatalf("expected project terraform, got %s", d.Get("project"))
	}
	if d.Get("zone_name") != "zone-1" {
		t.Fatalf("expected zone_name zone-1, got %s", d.Get("zone_name"))
	}
}

func TestVpcDescriptionAttributes_emptyTags(t *testing.T) {
	raw := dataSourceCloudstackVPC()
	d := raw.TestResourceData()

	vpc := &cloudstack.VPC{
		Id:   "vpc-empty",
		Name: "empty-tags",
		Tags: []cloudstack.Tags{},
	}

	err := vpcDescriptionAttributes(d, vpc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := d.Get("tags").(map[string]interface{})
	if len(tags) != 0 {
		t.Fatalf("expected empty tags, got %v", tags)
	}
}

// ---------------------------------------------------------------------------
// dataSourceCloudstackVPC schema
// ---------------------------------------------------------------------------

func TestDataSourceVPCSchema(t *testing.T) {
	ds := dataSourceCloudstackVPC()
	if ds.Read == nil {
		t.Fatal("expected Read function to be set")
	}
	expectedComputed := []string{"name", "display_text", "cidr", "vpc_offering_name", "network_domain", "zone_name"}
	for _, key := range expectedComputed {
		s, ok := ds.Schema[key]
		if !ok {
			t.Fatalf("expected schema key %q", key)
		}
		if !s.Computed {
			t.Fatalf("expected %q to be Computed", key)
		}
	}
	projectSchema := ds.Schema["project"]
	if !projectSchema.Computed || !projectSchema.Optional {
		t.Fatal("expected project to be Computed+Optional")
	}
	if _, ok := ds.Schema["filter"]; !ok {
		t.Fatal("expected filter in schema")
	}
	if _, ok := ds.Schema["tags"]; !ok {
		t.Fatal("expected tags in schema")
	}
}
