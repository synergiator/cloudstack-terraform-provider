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

func makeRule(protocol, trafficType, action string, extra map[string]interface{}) map[string]interface{} {
	r := map[string]interface{}{
		"protocol":     protocol,
		"traffic_type": trafficType,
		"action":       action,
		"cidr_list":    []interface{}{"10.0.0.0/8"},
	}
	for k, v := range extra {
		r[k] = v
	}
	return r
}

func TestRulesMatch_tcpSamePort(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	if !rulesMatch(old, new) {
		t.Fatal("expected rules to match")
	}
}

func TestRulesMatch_tcpDiffPort(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "443"})
	if rulesMatch(old, new) {
		t.Fatal("expected rules NOT to match")
	}
}

func TestRulesMatch_icmpMatch(t *testing.T) {
	old := makeRule("icmp", "ingress", "allow", map[string]interface{}{"icmp_type": 8, "icmp_code": 0})
	new := makeRule("icmp", "ingress", "allow", map[string]interface{}{"icmp_type": 8, "icmp_code": 0})
	if !rulesMatch(old, new) {
		t.Fatal("expected ICMP rules to match")
	}
}

func TestRulesMatch_icmpMismatch(t *testing.T) {
	old := makeRule("icmp", "ingress", "allow", map[string]interface{}{"icmp_type": 8, "icmp_code": 0})
	new := makeRule("icmp", "ingress", "allow", map[string]interface{}{"icmp_type": 0, "icmp_code": 0})
	if rulesMatch(old, new) {
		t.Fatal("expected ICMP rules NOT to match")
	}
}

func TestRulesMatch_protocolMismatch(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("udp", "ingress", "allow", map[string]interface{}{"port": "80"})
	if rulesMatch(old, new) {
		t.Fatal("expected rules NOT to match with different protocols")
	}
}

func TestRulesMatch_allProtocol(t *testing.T) {
	old := makeRule("all", "ingress", "allow", nil)
	new := makeRule("all", "ingress", "allow", nil)
	if !rulesMatch(old, new) {
		t.Fatal("expected 'all' protocol rules to match")
	}
}

func TestRulesMatch_actionMismatch(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("tcp", "ingress", "deny", map[string]interface{}{"port": "80"})
	if rulesMatch(old, new) {
		t.Fatal("expected rules NOT to match with different actions")
	}
}

func TestRulesMatch_tcpNoPort(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", nil)
	new := makeRule("tcp", "ingress", "allow", nil)
	if !rulesMatch(old, new) {
		t.Fatal("expected TCP rules without port to match")
	}
}

func TestRulesMatch_tcpOneHasPort(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("tcp", "ingress", "allow", nil)
	if rulesMatch(old, new) {
		t.Fatal("expected rules NOT to match when one has port and other doesn't")
	}
}

func TestRuleNeedsUpdate_noChange(t *testing.T) {
	rule := makeRule("tcp", "ingress", "allow", map[string]interface{}{
		"port":        "80",
		"rule_number": 1,
		"description": "web",
	})
	ruleCopy := makeRule("tcp", "ingress", "allow", map[string]interface{}{
		"port":        "80",
		"rule_number": 1,
		"description": "web",
	})
	if ruleNeedsUpdate(rule, ruleCopy) {
		t.Fatal("expected no update needed for identical rules")
	}
}

func TestRuleNeedsUpdate_actionChanged(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("tcp", "ingress", "deny", map[string]interface{}{"port": "80"})
	if !ruleNeedsUpdate(old, new) {
		t.Fatal("expected update needed when action changes")
	}
}

func TestRuleNeedsUpdate_portChanged(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "443"})
	if !ruleNeedsUpdate(old, new) {
		t.Fatal("expected update needed when port changes")
	}
}

func TestRuleNeedsUpdate_cidrChanged(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	old["cidr_list"] = []interface{}{"10.0.0.0/8"}
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new["cidr_list"] = []interface{}{"192.168.0.0/16"}
	if !ruleNeedsUpdate(old, new) {
		t.Fatal("expected update needed when CIDR changes")
	}
}

func TestRuleNeedsUpdate_descriptionChanged(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{
		"port":        "80",
		"description": "old desc",
	})
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{
		"port":        "80",
		"description": "new desc",
	})
	if !ruleNeedsUpdate(old, new) {
		t.Fatal("expected update needed when description changes")
	}
}

func TestRuleNeedsUpdate_ruleNumberChanged(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{
		"port":        "80",
		"rule_number": 1,
	})
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{
		"port":        "80",
		"rule_number": 2,
	})
	if !ruleNeedsUpdate(old, new) {
		t.Fatal("expected update needed when rule_number changes")
	}
}

func TestRuleNeedsUpdate_icmpChanged(t *testing.T) {
	old := makeRule("icmp", "ingress", "allow", map[string]interface{}{"icmp_type": 8, "icmp_code": 0})
	new := makeRule("icmp", "ingress", "allow", map[string]interface{}{"icmp_type": 0, "icmp_code": 0})
	if !ruleNeedsUpdate(old, new) {
		t.Fatal("expected update needed when ICMP type changes")
	}
}

func TestRuleNeedsUpdate_cidrOrderDoesNotMatter(t *testing.T) {
	old := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	old["cidr_list"] = []interface{}{"10.0.0.0/8", "192.168.0.0/16"}
	new := makeRule("tcp", "ingress", "allow", map[string]interface{}{"port": "80"})
	new["cidr_list"] = []interface{}{"192.168.0.0/16", "10.0.0.0/8"}
	if ruleNeedsUpdate(old, new) {
		t.Fatal("expected no update when CIDRs are the same but in different order")
	}
}
