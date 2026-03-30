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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDiffTags_noChanges(t *testing.T) {
	old := map[string]string{"foo": "bar", "baz": "qux"}
	new := map[string]string{"foo": "bar", "baz": "qux"}
	remove, create := diffTags(old, new)
	if len(remove) != 0 {
		t.Fatalf("expected no removes, got %v", remove)
	}
	if len(create) != 0 {
		t.Fatalf("expected no creates, got %v", create)
	}
}

func TestDiffTags_addOnly(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"foo": "bar", "baz": "qux"}
	remove, create := diffTags(old, new)
	if len(remove) != 0 {
		t.Fatalf("expected no removes, got %v", remove)
	}
	expected := map[string]string{"foo": "bar", "baz": "qux"}
	if !reflect.DeepEqual(create, expected) {
		t.Fatalf("expected creates %v, got %v", expected, create)
	}
}

func TestDiffTags_removeOnly(t *testing.T) {
	old := map[string]string{"foo": "bar", "baz": "qux"}
	new := map[string]string{}
	remove, create := diffTags(old, new)
	if len(create) != 0 {
		t.Fatalf("expected no creates, got %v", create)
	}
	expected := map[string]string{"foo": "bar", "baz": "qux"}
	if !reflect.DeepEqual(remove, expected) {
		t.Fatalf("expected removes %v, got %v", expected, remove)
	}
}

func TestDiffTags_mixedOperations(t *testing.T) {
	old := map[string]string{"keep": "same", "remove": "old", "change": "old"}
	new := map[string]string{"keep": "same", "add": "new", "change": "new"}
	remove, create := diffTags(old, new)
	expectedRemove := map[string]string{"remove": "old", "change": "old"}
	expectedCreate := map[string]string{"add": "new", "change": "new"}
	if !reflect.DeepEqual(remove, expectedRemove) {
		t.Fatalf("expected removes %v, got %v", expectedRemove, remove)
	}
	if !reflect.DeepEqual(create, expectedCreate) {
		t.Fatalf("expected creates %v, got %v", expectedCreate, create)
	}
}

func TestTagsFromSchema(t *testing.T) {
	cases := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]string
	}{
		{
			name:     "empty",
			input:    map[string]interface{}{},
			expected: map[string]string{},
		},
		{
			name:     "single entry",
			input:    map[string]interface{}{"key": "value"},
			expected: map[string]string{"key": "value"},
		},
		{
			name:     "multiple entries",
			input:    map[string]interface{}{"a": "1", "b": "2", "c": "3"},
			expected: map[string]string{"a": "1", "b": "2", "c": "3"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tagsFromSchema(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTagsToMap(t *testing.T) {
	cases := []struct {
		name     string
		input    []cloudstack.Tags
		expected map[string]string
	}{
		{
			name:     "empty",
			input:    []cloudstack.Tags{},
			expected: map[string]string{},
		},
		{
			name: "single tag",
			input: []cloudstack.Tags{
				{Key: "env", Value: "prod"},
			},
			expected: map[string]string{"env": "prod"},
		},
		{
			name: "multiple tags",
			input: []cloudstack.Tags{
				{Key: "env", Value: "prod"},
				{Key: "team", Value: "infra"},
			},
			expected: map[string]string{"env": "prod", "team": "infra"},
		},
		{
			name: "duplicate keys last wins",
			input: []cloudstack.Tags{
				{Key: "env", Value: "dev"},
				{Key: "env", Value: "prod"},
			},
			expected: map[string]string{"env": "prod"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tagsToMap(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestTagsSchema(t *testing.T) {
	s := tagsSchema()
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
