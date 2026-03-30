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
	"context"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCommonRead_populatesAllFields(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{
		Id:                    "abc-123",
		Name:                  "test-offering",
		Displaytext:           "Test Offering",
		Deploymentplanner:     "ImplicitDedication",
		Diskofferingid:        "disk-456",
		Domainid:              "dom-1,dom-2",
		Hosttags:              "ssd",
		Networkrate:           100,
		Zoneid:                "zone-1,zone-2",
		Dynamicscalingenabled: true,
		Isvolatile:            false,
		Limitcpuuse:           true,
		Offerha:               true,
	}
	state := &serviceOfferingCommonResourceModel{}
	state.commonRead(ctx, so)

	if state.Id.ValueString() != "abc-123" {
		t.Fatalf("expected Id 'abc-123', got %q", state.Id.ValueString())
	}
	if state.Name.ValueString() != "test-offering" {
		t.Fatalf("expected Name 'test-offering', got %q", state.Name.ValueString())
	}
	if state.DisplayText.ValueString() != "Test Offering" {
		t.Fatalf("expected DisplayText 'Test Offering', got %q", state.DisplayText.ValueString())
	}
	if state.DeploymentPlanner.ValueString() != "ImplicitDedication" {
		t.Fatalf("expected DeploymentPlanner 'ImplicitDedication', got %q", state.DeploymentPlanner.ValueString())
	}
	if state.DiskOfferingId.ValueString() != "disk-456" {
		t.Fatalf("expected DiskOfferingId 'disk-456', got %q", state.DiskOfferingId.ValueString())
	}
	if state.HostTags.ValueString() != "ssd" {
		t.Fatalf("expected HostTags 'ssd', got %q", state.HostTags.ValueString())
	}
	if state.NetworkRate.ValueInt32() != 100 {
		t.Fatalf("expected NetworkRate 100, got %d", state.NetworkRate.ValueInt32())
	}
	if state.DynamicScalingEnabled.ValueBool() != true {
		t.Fatal("expected DynamicScalingEnabled true")
	}
	if state.LimitCpuUse.ValueBool() != true {
		t.Fatal("expected LimitCpuUse true")
	}
	if state.OfferHa.ValueBool() != true {
		t.Fatal("expected OfferHa true")
	}
	if state.IsVolatile.ValueBool() != false {
		t.Fatal("expected IsVolatile false")
	}
}

func TestCommonRead_emptyStringsSkipped(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{
		Id:   "abc-123",
		Name: "test",
	}
	state := &serviceOfferingCommonResourceModel{}
	state.commonRead(ctx, so)

	if state.Id.ValueString() != "abc-123" {
		t.Fatalf("expected Id set, got %q", state.Id.ValueString())
	}
	if !state.DeploymentPlanner.IsNull() && state.DeploymentPlanner.ValueString() != "" {
		t.Fatalf("expected DeploymentPlanner to remain zero, got %q", state.DeploymentPlanner.ValueString())
	}
	if !state.HostTags.IsNull() && state.HostTags.ValueString() != "" {
		t.Fatalf("expected HostTags to remain zero, got %q", state.HostTags.ValueString())
	}
}

func TestCommonRead_booleanFields(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{
		Id:                    "id-1",
		Dynamicscalingenabled: false,
		Isvolatile:            true,
		Limitcpuuse:           false,
		Offerha:               false,
	}
	state := &serviceOfferingCommonResourceModel{}
	state.commonRead(ctx, so)

	if state.DynamicScalingEnabled.ValueBool() != false {
		t.Fatal("expected DynamicScalingEnabled false")
	}
	if state.IsVolatile.ValueBool() != true {
		t.Fatal("expected IsVolatile true")
	}
	if state.LimitCpuUse.ValueBool() != false {
		t.Fatal("expected LimitCpuUse false")
	}
	if state.OfferHa.ValueBool() != false {
		t.Fatal("expected OfferHa false")
	}
}

func TestDiskQosHypervisorRead(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{
		DiskBytesReadRate:           1000,
		DiskBytesReadRateMax:        2000,
		DiskBytesReadRateMaxLength:  30,
		DiskBytesWriteRate:          500,
		DiskBytesWriteRateMax:       1000,
		DiskBytesWriteRateMaxLength: 60,
	}
	state := &ServiceOfferingDiskQosHypervisor{}
	state.commonRead(ctx, so)

	if state.DiskBytesReadRate.ValueInt64() != 1000 {
		t.Fatalf("expected DiskBytesReadRate 1000, got %d", state.DiskBytesReadRate.ValueInt64())
	}
	if state.DiskBytesWriteRateMax.ValueInt64() != 1000 {
		t.Fatalf("expected DiskBytesWriteRateMax 1000, got %d", state.DiskBytesWriteRateMax.ValueInt64())
	}
}

func TestDiskQosHypervisorRead_zeroValuesSkipped(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{}
	state := &ServiceOfferingDiskQosHypervisor{}
	state.commonRead(ctx, so)

	if !state.DiskBytesReadRate.IsNull() && state.DiskBytesReadRate.ValueInt64() != 0 {
		t.Fatal("expected DiskBytesReadRate to remain zero value")
	}
}

func TestDiskOfferingRead(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{
		CacheMode:              "writeback",
		Diskofferingstrictness: true,
		Provisioningtype:       "thin",
		Rootdisksize:           50,
		Storagetype:            "shared",
		Storagetags:            "fast",
	}
	state := &ServiceOfferingDiskOffering{}
	state.commonRead(ctx, so)

	if state.CacheMode.ValueString() != "writeback" {
		t.Fatalf("expected CacheMode 'writeback', got %q", state.CacheMode.ValueString())
	}
	if state.DiskOfferingStrictness.ValueBool() != true {
		t.Fatal("expected DiskOfferingStrictness true")
	}
	if state.ProvisionType.ValueString() != "thin" {
		t.Fatalf("expected ProvisionType 'thin', got %q", state.ProvisionType.ValueString())
	}
	if state.RootDiskSize.ValueInt64() != 50 {
		t.Fatalf("expected RootDiskSize 50, got %d", state.RootDiskSize.ValueInt64())
	}
	if state.StorageType.ValueString() != "shared" {
		t.Fatalf("expected StorageType 'shared', got %q", state.StorageType.ValueString())
	}
	if state.StorageTags.ValueString() != "fast" {
		t.Fatalf("expected StorageTags 'fast', got %q", state.StorageTags.ValueString())
	}
}

func TestDiskQosStorageRead(t *testing.T) {
	ctx := context.Background()
	so := &cloudstack.ServiceOffering{
		Iscustomizediops:          true,
		Hypervisorsnapshotreserve: 20,
		Maxiops:                   5000,
		Miniops:                   100,
	}
	state := &ServiceOfferingDiskQosStorage{}
	state.commonRead(ctx, so)

	if state.CustomizedIops.ValueBool() != true {
		t.Fatal("expected CustomizedIops true")
	}
	if state.HypervisorSnapshotReserve.ValueInt32() != 20 {
		t.Fatalf("expected HypervisorSnapshotReserve 20, got %d", state.HypervisorSnapshotReserve.ValueInt32())
	}
	if state.MaxIops.ValueInt64() != 5000 {
		t.Fatalf("expected MaxIops 5000, got %d", state.MaxIops.ValueInt64())
	}
	if state.MinIops.ValueInt64() != 100 {
		t.Fatalf("expected MinIops 100, got %d", state.MinIops.ValueInt64())
	}
}

func TestCommonUpdate_populatesFields(t *testing.T) {
	ctx := context.Background()
	resp := &cloudstack.UpdateServiceOfferingResponse{
		Displaytext: "updated text",
		Name:        "updated name",
		Hosttags:    "newtag",
		Domainid:    "dom-1",
		Zoneid:      "zone-1",
	}
	state := &serviceOfferingCommonResourceModel{
		DisplayText: types.StringValue("old text"),
		Name:        types.StringValue("old name"),
	}
	state.commonUpdate(ctx, resp)

	if state.DisplayText.ValueString() != "updated text" {
		t.Fatalf("expected DisplayText 'updated text', got %q", state.DisplayText.ValueString())
	}
	if state.Name.ValueString() != "updated name" {
		t.Fatalf("expected Name 'updated name', got %q", state.Name.ValueString())
	}
	if state.HostTags.ValueString() != "newtag" {
		t.Fatalf("expected HostTags 'newtag', got %q", state.HostTags.ValueString())
	}
}

func TestCommonUpdate_emptyStringsSkipped(t *testing.T) {
	ctx := context.Background()
	resp := &cloudstack.UpdateServiceOfferingResponse{
		Name: "updated",
	}
	state := &serviceOfferingCommonResourceModel{
		DisplayText: types.StringValue("keep this"),
		Name:        types.StringValue("old"),
	}
	state.commonUpdate(ctx, resp)

	if state.Name.ValueString() != "updated" {
		t.Fatalf("expected Name 'updated', got %q", state.Name.ValueString())
	}
	if state.DisplayText.ValueString() != "keep this" {
		t.Fatalf("expected DisplayText 'keep this' (unchanged), got %q", state.DisplayText.ValueString())
	}
}
