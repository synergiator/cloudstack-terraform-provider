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
// latestVPNConnection
// ---------------------------------------------------------------------------

func TestLatestVPNConnection_single(t *testing.T) {
	conns := []*cloudstack.VpnConnection{
		{Id: "vpn-1", Created: "2025-01-15T10:00:00+0000"},
	}
	got, err := latestVPNConnection(conns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vpn-1" {
		t.Fatalf("expected vpn-1, got %s", got.Id)
	}
}

func TestLatestVPNConnection_picksMostRecent(t *testing.T) {
	conns := []*cloudstack.VpnConnection{
		{Id: "vpn-old", Created: "2024-01-01T00:00:00+0000"},
		{Id: "vpn-new", Created: "2025-06-15T12:00:00+0000"},
		{Id: "vpn-mid", Created: "2025-03-10T08:00:00+0000"},
	}
	got, err := latestVPNConnection(conns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vpn-new" {
		t.Fatalf("expected vpn-new, got %s", got.Id)
	}
}

func TestLatestVPNConnection_badDateFormat(t *testing.T) {
	conns := []*cloudstack.VpnConnection{
		{Id: "vpn-bad", Created: "garbage"},
	}
	_, err := latestVPNConnection(conns)
	if err == nil {
		t.Fatal("expected error for bad date format")
	}
}

func TestLatestVPNConnection_sameTimestamp(t *testing.T) {
	conns := []*cloudstack.VpnConnection{
		{Id: "vpn-a", Created: "2025-01-01T00:00:00+0000"},
		{Id: "vpn-b", Created: "2025-01-01T00:00:00+0000"},
	}
	got, err := latestVPNConnection(conns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Id != "vpn-a" {
		t.Fatalf("expected vpn-a (first seen), got %s", got.Id)
	}
}

// ---------------------------------------------------------------------------
// applyVPNConnectionFilters
// ---------------------------------------------------------------------------

func TestApplyVPNConnectionFilters_matchByID(t *testing.T) {
	conn := &cloudstack.VpnConnection{
		Id:                   "vpn-1",
		S2scustomergatewayid: "cg-1",
		S2svpngatewayid:      "gw-1",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "id",
		"value": "vpn-1",
	})
	match, err := applyVPNConnectionFilters(conn, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match")
	}
}

func TestApplyVPNConnectionFilters_noMatch(t *testing.T) {
	conn := &cloudstack.VpnConnection{
		Id:                   "vpn-1",
		S2scustomergatewayid: "cg-1",
		S2svpngatewayid:      "gw-1",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "id",
		"value": "vpn-other",
	})
	match, err := applyVPNConnectionFilters(conn, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected no match")
	}
}

func TestApplyVPNConnectionFilters_regexMatch(t *testing.T) {
	conn := &cloudstack.VpnConnection{
		Id:      "vpn-123-abc",
		Gateway: "192.168.1.1",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "gateway",
		"value": "^192\\.168\\..*",
	})
	match, err := applyVPNConnectionFilters(conn, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected regex match")
	}
}

func TestApplyVPNConnectionFilters_invalidRegex(t *testing.T) {
	conn := &cloudstack.VpnConnection{
		Id: "vpn-1",
	}
	filters := newFilterSet(map[string]interface{}{
		"name":  "id",
		"value": "[bad",
	})
	_, err := applyVPNConnectionFilters(conn, filters)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestApplyVPNConnectionFilters_emptyFilters(t *testing.T) {
	conn := &cloudstack.VpnConnection{
		Id: "vpn-1",
	}
	filters := newFilterSet()
	match, err := applyVPNConnectionFilters(conn, filters)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected match with no filters")
	}
}

// ---------------------------------------------------------------------------
// vpnConnectionDescriptionAttributes
// ---------------------------------------------------------------------------

func TestVpnConnectionDescriptionAttributes(t *testing.T) {
	raw := dataSourceCloudstackVPNConnection()
	d := raw.TestResourceData()

	conn := &cloudstack.VpnConnection{
		Id:                   "vpn-123",
		S2scustomergatewayid: "cg-456",
		S2svpngatewayid:      "gw-789",
	}

	err := vpnConnectionDescriptionAttributes(d, conn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Id() != "vpn-123" {
		t.Fatalf("expected id vpn-123, got %s", d.Id())
	}
	if d.Get("s2s_customer_gateway_id") != "cg-456" {
		t.Fatalf("expected s2s_customer_gateway_id cg-456, got %s", d.Get("s2s_customer_gateway_id"))
	}
	if d.Get("s2s_vpn_gateway_id") != "gw-789" {
		t.Fatalf("expected s2s_vpn_gateway_id gw-789, got %s", d.Get("s2s_vpn_gateway_id"))
	}
}

// ---------------------------------------------------------------------------
// dataSourceCloudstackVPNConnection schema
// ---------------------------------------------------------------------------

func TestDataSourceVPNConnectionSchema(t *testing.T) {
	ds := dataSourceCloudstackVPNConnection()
	if ds.Read == nil {
		t.Fatal("expected Read function to be set")
	}
	expectedComputed := []string{"s2s_customer_gateway_id", "s2s_vpn_gateway_id"}
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
}
