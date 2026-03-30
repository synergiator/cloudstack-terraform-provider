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

func TestValidateTrafficType_valid(t *testing.T) {
	validTypes := []string{"Public", "Guest", "Management", "Storage"}
	for _, tt := range validTypes {
		t.Run(tt, func(t *testing.T) {
			warnings, errs := validateTrafficType(tt, "traffic_type")
			if len(warnings) != 0 {
				t.Fatalf("expected no warnings, got %v", warnings)
			}
			if len(errs) != 0 {
				t.Fatalf("expected no errors for %q, got %v", tt, errs)
			}
		})
	}
}

func TestValidateTrafficType_invalid(t *testing.T) {
	warnings, errs := validateTrafficType("Invalid", "traffic_type")
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidateTrafficType_caseSensitive(t *testing.T) {
	_, errs := validateTrafficType("public", "traffic_type")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for lowercase 'public', got %d", len(errs))
	}
}

func TestValidateTrafficType_empty(t *testing.T) {
	_, errs := validateTrafficType("", "traffic_type")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error for empty string, got %d", len(errs))
	}
}
