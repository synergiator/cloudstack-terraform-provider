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
	"reflect"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
)

func TestDsFlattenPodCapacity(t *testing.T) {
	input := []cloudstack.PodCapacity{
		{
			Capacityallocated: 100,
			Capacitytotal:     1000,
			Capacityused:      50,
			Clusterid:         "cluster-1",
			Clustername:       "Cluster One",
			Name:              "CPU",
			Percentused:       "5%",
			Podid:             "pod-1",
			Podname:           "Pod One",
			Type:              1,
			Zoneid:            "zone-1",
			Zonename:          "Zone One",
		},
	}
	result := dsFlattenPodCapacity(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0]["cluster_id"] != "cluster-1" {
		t.Fatalf("expected cluster_id 'cluster-1', got %v", result[0]["cluster_id"])
	}
	if result[0]["capacity_total"] != int64(1000) {
		t.Fatalf("expected capacity_total 1000, got %v", result[0]["capacity_total"])
	}
	if result[0]["zone_name"] != "Zone One" {
		t.Fatalf("expected zone_name 'Zone One', got %v", result[0]["zone_name"])
	}
}

func TestDsFlattenPodCapacity_empty(t *testing.T) {
	result := dsFlattenPodCapacity([]cloudstack.PodCapacity{})
	if len(result) != 0 {
		t.Fatalf("expected 0 items, got %d", len(result))
	}
}

func TestDsFlattenPodIpRanges(t *testing.T) {
	input := []cloudstack.PodIpranges{
		{
			Endip:        "10.0.0.100",
			Forsystemvms: "true",
			Startip:      "10.0.0.1",
			Vlanid:       "vlan-1",
		},
	}
	result := dsFlattenPodIpRanges(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	if result[0]["end_ip"] != "10.0.0.100" {
		t.Fatalf("expected end_ip '10.0.0.100', got %v", result[0]["end_ip"])
	}
	if result[0]["start_ip"] != "10.0.0.1" {
		t.Fatalf("expected start_ip '10.0.0.1', got %v", result[0]["start_ip"])
	}
	if result[0]["for_system_vms"] != "true" {
		t.Fatalf("expected for_system_vms 'true', got %v", result[0]["for_system_vms"])
	}
	if result[0]["vlan_id"] != "vlan-1" {
		t.Fatalf("expected vlan_id 'vlan-1', got %v", result[0]["vlan_id"])
	}
}

func TestDsFlattenPodIpRanges_empty(t *testing.T) {
	result := dsFlattenPodIpRanges([]cloudstack.PodIpranges{})
	if len(result) != 0 {
		t.Fatalf("expected 0 items, got %d", len(result))
	}
}

func TestDsFlattenClusterCapacity(t *testing.T) {
	input := []cloudstack.ClusterCapacity{
		{
			Capacityallocated: 200,
			Capacitytotal:     2000,
			Capacityused:      100,
			Clusterid:         "cluster-2",
			Clustername:       "Cluster Two",
			Name:              "Memory",
			Percentused:       "10%",
			Podid:             "pod-2",
			Podname:           "Pod Two",
			Type:              2,
			Zoneid:            "zone-2",
			Zonename:          "Zone Two",
		},
	}
	result := dsFlattenClusterCapacity(input)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	expected := map[string]interface{}{
		"capacity_allocated": int64(200),
		"capacity_total":     int64(2000),
		"capacity_used":      int64(100),
		"cluster_id":         "cluster-2",
		"cluster_name":       "Cluster Two",
		"name":               "Memory",
		"percent_used":       "10%",
		"pod_id":             "pod-2",
		"pod_name":           "Pod Two",
		"type":               2,
		"zone_id":            "zone-2",
		"zone_name":          "Zone Two",
	}
	if !reflect.DeepEqual(result[0], expected) {
		t.Fatalf("expected %v, got %v", expected, result[0])
	}
}

func TestDsFlattenClusterCapacity_empty(t *testing.T) {
	result := dsFlattenClusterCapacity([]cloudstack.ClusterCapacity{})
	if len(result) != 0 {
		t.Fatalf("expected 0 items, got %d", len(result))
	}
}
