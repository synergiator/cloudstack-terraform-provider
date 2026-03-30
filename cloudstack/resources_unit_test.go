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
	"errors"
	"strings"
	"testing"

	"github.com/apache/cloudstack-go/v2/cloudstack"
)

func TestSplitPorts_singlePort(t *testing.T) {
	m := splitPorts.FindStringSubmatch("80")
	if m == nil {
		t.Fatal("expected match for single port")
	}
	if m[1] != "80" {
		t.Fatalf("expected group 1 = '80', got '%s'", m[1])
	}
	if m[2] != "" {
		t.Fatalf("expected group 2 = '', got '%s'", m[2])
	}
}

func TestSplitPorts_portRange(t *testing.T) {
	m := splitPorts.FindStringSubmatch("80-443")
	if m == nil {
		t.Fatal("expected match for port range")
	}
	if m[1] != "80" {
		t.Fatalf("expected group 1 = '80', got '%s'", m[1])
	}
	if m[2] != "443" {
		t.Fatalf("expected group 2 = '443', got '%s'", m[2])
	}
}

func TestSplitPorts_invalid(t *testing.T) {
	invalids := []string{"abc", "", "80-", "-80", "80-443-900", "foo-bar"}
	for _, input := range invalids {
		m := splitPorts.FindStringSubmatch(input)
		if m != nil {
			t.Fatalf("expected no match for '%s', got %v", input, m)
		}
	}
}

func TestRetrieveError_Error(t *testing.T) {
	e := &retrieveError{
		name:  "zone",
		value: "us-east-1",
		err:   errors.New("not found"),
	}
	result := e.Error()
	if result == nil {
		t.Fatal("expected non-nil error")
	}
	msg := result.Error()
	if !strings.Contains(msg, "zone") {
		t.Fatalf("error should mention 'zone', got: %s", msg)
	}
	if !strings.Contains(msg, "us-east-1") {
		t.Fatalf("error should mention 'us-east-1', got: %s", msg)
	}
	if !strings.Contains(msg, "not found") {
		t.Fatalf("error should mention 'not found', got: %s", msg)
	}
}

func TestRetry_successFirst(t *testing.T) {
	calls := 0
	result, err := Retry(3, func() (interface{}, error) {
		calls++
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "ok" {
		t.Fatalf("expected 'ok', got %v", result)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetry_allFailures(t *testing.T) {
	calls := 0
	_, err := Retry(1, func() (interface{}, error) {
		calls++
		return nil, errors.New("fail")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "fail" {
		t.Fatalf("expected 'fail', got '%s'", err.Error())
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetry_asyncTimeout(t *testing.T) {
	calls := 0
	result, err := Retry(3, func() (interface{}, error) {
		calls++
		return "partial", cloudstack.AsyncTimeoutErr
	})
	if err != cloudstack.AsyncTimeoutErr {
		t.Fatalf("expected AsyncTimeoutErr, got %v", err)
	}
	if result != "partial" {
		t.Fatalf("expected 'partial', got %v", result)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call (no retry on async timeout), got %d", calls)
	}
}
