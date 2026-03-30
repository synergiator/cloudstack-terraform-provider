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

import "testing"

func TestResourceTypeMap_allTypesPresent(t *testing.T) {
	expectedTypes := []string{
		"instance", "ip", "volume", "snapshot", "template", "project",
		"network", "vpc", "cpu", "memory", "primarystorage", "secondarystorage",
	}
	for _, rt := range expectedTypes {
		if _, ok := resourceTypeMap[rt]; !ok {
			t.Fatalf("expected resource type %q in resourceTypeMap", rt)
		}
	}
}

func TestResourceTypeMap_correctValues(t *testing.T) {
	expected := map[string]int{
		"instance":         0,
		"ip":               1,
		"volume":           2,
		"snapshot":         3,
		"template":         4,
		"project":          5,
		"network":          6,
		"vpc":              7,
		"cpu":              8,
		"memory":           9,
		"primarystorage":   10,
		"secondarystorage": 11,
	}
	for k, v := range expected {
		if resourceTypeMap[k] != v {
			t.Fatalf("expected resourceTypeMap[%q] = %d, got %d", k, v, resourceTypeMap[k])
		}
	}
}

func TestResourceTypeMap_uniqueValues(t *testing.T) {
	seen := make(map[int]string)
	for k, v := range resourceTypeMap {
		if prev, ok := seen[v]; ok {
			t.Fatalf("duplicate value %d for keys %q and %q", v, prev, k)
		}
		seen[v] = k
	}
}

func TestResourceTypeMap_count(t *testing.T) {
	if len(resourceTypeMap) != 12 {
		t.Fatalf("expected 12 entries in resourceTypeMap, got %d", len(resourceTypeMap))
	}
}
