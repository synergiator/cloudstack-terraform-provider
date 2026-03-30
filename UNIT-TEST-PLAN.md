# Unit Test Plan — cloudstack-terraform-provider

## Status Quo

Only **3 true unit tests** exist today (`TestProvider`, `TestProvider_impl`, `TestDiffTags`). Everything else is an acceptance test requiring a running CloudStack simulator. This plan adds **3–4 unit tests per component** for all pure-logic functions that can be tested without any infrastructure.

---

## Master Task List

### 1. Tags — `tags.go` / `tags_test.go`

**Existing**: `TestDiffTags` (2 cases only).

| # | Test | What it validates |
|---|------|-------------------|
| 1.1 | `TestDiffTags_noChanges` | Identical old/new maps return empty remove and create sets |
| 1.2 | `TestDiffTags_addOnly` | Empty old map, populated new map → only creates, no removes |
| 1.3 | `TestDiffTags_removeOnly` | Populated old map, empty new map → only removes, no creates |
| 1.4 | `TestDiffTags_mixedOperations` | Disjoint add + remove + unchanged keys in a single diff |
| 1.5 | `TestTagsFromSchema` | Converts `map[string]interface{}` to `map[string]string`; empty input; multiple entries |
| 1.6 | `TestTagsToMap` | Converts `[]cloudstack.Tags` to `map[string]string`; empty slice; duplicate keys (last wins) |
| 1.7 | `TestTagsSchema` | `tagsSchema()` returns correct Type, Optional, Computed flags |

---

### 2. Resources Helpers — `resources.go` / `resources_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 2.1 | `TestSplitPorts_singlePort` | `splitPorts` regex matches `"80"` → groups `["80", ""]` |
| 2.2 | `TestSplitPorts_portRange` | `splitPorts` regex matches `"80-443"` → groups `["80", "443"]` |
| 2.3 | `TestSplitPorts_invalid` | Non-numeric strings like `"abc"` or `"80-"` don't match |
| 2.4 | `TestRetrieveError_Error` | `retrieveError{name:"zone", value:"foo", err:…}.Error()` produces expected message format |
| 2.5 | `TestRetry_successFirst` | `Retry(3, fn)` returns immediately when fn succeeds on first call |
| 2.6 | `TestRetry_successAfterFailures` | `Retry(3, fn)` succeeds on attempt 2 after 1 failure |
| 2.7 | `TestRetry_allFailures` | `Retry(2, fn)` returns last error when all attempts fail |
| 2.8 | `TestRetry_asyncTimeout` | `Retry(3, fn)` returns `cloudstack.AsyncTimeoutErr` without retrying |

> **Note**: `Retry` calls `time.Sleep(30s)`. Tests should use a counter-based function that succeeds on attempt N to avoid actual sleep (it only sleeps on error, so success-on-first-call is instant; for failure paths, keep retry count = 1 to minimize sleep or consider refactoring `Retry` to accept a sleep function).

---

### 3. Snapshot Policy Helpers — `resource_cloudstack_snapshot_policy.go` / `resource_cloudstack_snapshot_policy_unit_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 3.1 | `TestIntervalTypeToString_hourly` | `intervalTypeToString(0)` → `"HOURLY"` |
| 3.2 | `TestIntervalTypeToString_daily` | `intervalTypeToString(1)` → `"DAILY"` |
| 3.3 | `TestIntervalTypeToString_weekly` | `intervalTypeToString(2)` → `"WEEKLY"` |
| 3.4 | `TestIntervalTypeToString_monthly` | `intervalTypeToString(3)` → `"MONTHLY"` |
| 3.5 | `TestIntervalTypeToString_unknown` | `intervalTypeToString(99)` → `"99"` (string of the int) |

---

### 4. Traffic Type Validation — `resource_cloudstack_traffic_type.go` / `resource_cloudstack_traffic_type_unit_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 4.1 | `TestValidateTrafficType_valid` | Each of `"Public"`, `"Guest"`, `"Management"`, `"Storage"` returns no errors |
| 4.2 | `TestValidateTrafficType_invalid` | `"Invalid"` returns an error mentioning the valid set |
| 4.3 | `TestValidateTrafficType_caseSensitive` | `"public"` (lowercase) is rejected — validator is case-sensitive |
| 4.4 | `TestValidateTrafficType_empty` | `""` is rejected |

---

### 5. Network ACL Rule Helpers — `resource_cloudstack_network_acl_rule.go` / `resource_cloudstack_network_acl_rule_unit_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 5.1 | `TestRulesMatch_sameProtocolTcpSamePort` | Two TCP rules with same protocol/action/traffic_type/port → `true` |
| 5.2 | `TestRulesMatch_sameProtocolTcpDiffPort` | Two TCP rules with different ports → `false` |
| 5.3 | `TestRulesMatch_icmpMatch` | Two ICMP rules with same type/code → `true` |
| 5.4 | `TestRulesMatch_icmpMismatch` | Two ICMP rules with different type or code → `false` |
| 5.5 | `TestRulesMatch_protocolMismatch` | Rules with different protocols → `false` |
| 5.6 | `TestRulesMatch_allProtocol` | `protocol: "all"` rules always match |
| 5.7 | `TestRuleNeedsUpdate_noChange` | Identical rules → `false` |
| 5.8 | `TestRuleNeedsUpdate_actionChanged` | Same rule but different action → `true` |
| 5.9 | `TestRuleNeedsUpdate_portChanged` | Same TCP rule but different port → `true` |
| 5.10 | `TestRuleNeedsUpdate_cidrChanged` | Same rule but different cidr_list → `true` |
| 5.11 | `TestRuleNeedsUpdate_descriptionChanged` | Same rule but different description → `true` |
| 5.12 | `TestRuleNeedsUpdate_ruleNumberChanged` | Same rule but different rule_number → `true` |

---

### 6. Service Offering Schema — `service_offering_schema.go` / `service_offering_schema_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 6.1 | `TestServiceOfferingMergeCommonSchema_includesCommonFields` | Result contains all common keys (`name`, `display_text`, `id`, `offer_ha`, etc.) |
| 6.2 | `TestServiceOfferingMergeCommonSchema_includesCustomFields` | Custom fields passed in appear in the merged result |
| 6.3 | `TestServiceOfferingMergeCommonSchema_customOverridesCommon` | If custom map has same key as common, custom wins |
| 6.4 | `TestServiceOfferingMergeCommonSchema_emptyCustom` | Empty custom map → result equals just the common schema |

---

### 7. Service Offering Utils — `service_offering_util.go` / `service_offering_util_test.go` *(new file)*

Tests for the `commonRead` methods that transform CloudStack API response structs into Terraform model structs.

| # | Test | What it validates |
|---|------|-------------------|
| 7.1 | `TestCommonRead_populatesAllFields` | Given a fully populated `cloudstack.ServiceOffering`, all model fields are set correctly |
| 7.2 | `TestCommonRead_emptyStringsSkipped` | Empty-string fields in the API response leave model fields at their zero value |
| 7.3 | `TestCommonRead_booleanFields` | `Dynamicscalingenabled`, `Isvolatile`, `Limitcpuuse`, `Offerha` map correctly |
| 7.4 | `TestDiskQosHypervisorRead` | Given populated disk QoS fields, model is set correctly; zero values are skipped |
| 7.5 | `TestDiskOfferingRead` | `CacheMode`, `StorageType`, `ProvisionType`, `RootDiskSize` mapped correctly |
| 7.6 | `TestDiskQosStorageRead` | `CustomizedIops`, `MaxIops`, `MinIops`, `HypervisorSnapshotReserve` mapped correctly |
| 7.7 | `TestCommonUpdate_populatesFields` | Given an `UpdateServiceOfferingResponse`, model fields are updated |
| 7.8 | `TestCommonUpdate_emptyStringsSkipped` | Empty response fields don't overwrite existing model values |

---

### 8. Data Source Flatten Helpers — `data_source_cloudstack_pod.go` + `data_source_cloudstack_cluster.go` / `data_source_flatten_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 8.1 | `TestDsFlattenPodCapacity` | Converts `[]cloudstack.PodCapacity` → `[]map[string]interface{}` with correct keys |
| 8.2 | `TestDsFlattenPodCapacity_empty` | Empty input → empty output slice |
| 8.3 | `TestDsFlattenPodIpRanges` | Converts `[]cloudstack.PodIpranges` → correct map with `end_ip`, `start_ip`, `for_system_vms`, `vlan_id` |
| 8.4 | `TestDsFlattenPodIpRanges_empty` | Empty input → empty output slice |
| 8.5 | `TestDsFlattenClusterCapacity` | Converts `[]cloudstack.ClusterCapacity` → correct map structure |
| 8.6 | `TestDsFlattenClusterCapacity_empty` | Empty input → empty output slice |

---

### 9. Limits Resource Type Map — `resource_cloudstack_limits.go` / `resource_cloudstack_limits_unit_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 9.1 | `TestResourceTypeMap_allTypesPresent` | All 12 documented types exist in the map |
| 9.2 | `TestResourceTypeMap_correctValues` | Each type maps to its correct integer (instance=0, ip=1, …, secondarystorage=11) |
| 9.3 | `TestResourceTypeMap_uniqueValues` | No two types map to the same integer |

---

### 10. Provider & Schemas — `provider_test.go` (extend existing) + `schema_test.go` *(new file)*

| # | Test | What it validates |
|---|------|-------------------|
| 10.1 | `TestDataSourceFiltersSchema` | `dataSourceFiltersSchema()` returns TypeSet, Required, ForceNew with correct nested schema |
| 10.2 | `TestMetadataSchema` | `metadataSchema()` returns TypeMap, Optional, Computed |
| 10.3 | `TestProviderSchema_hasRequiredAttributes` | Provider schema contains `api_url`, `api_key`, `secret_key`, `config`, `profile`, `http_get_only`, `timeout` |
| 10.4 | `TestProviderSchema_sensitiveFields` | `api_key` and `secret_key` are marked Sensitive |

---

## Execution Order (recommended)

These tasks are **independent** and can be worked in any order or in parallel. The suggested order balances impact and complexity:

1. **Task 1** (Tags) — extend existing file, highest coverage gain per effort
2. **Task 2** (Resources helpers) — foundational helpers used everywhere
3. **Task 5** (ACL rule helpers) — richest pure logic, most edge cases
4. **Task 3** (Snapshot interval) — trivial, quick win
5. **Task 4** (Traffic type validation) — trivial, quick win
6. **Task 7** (Service offering utils) — tests the newer Framework pattern
7. **Task 6** (Service offering schema) — schema merge logic
8. **Task 8** (Flatten helpers) — data transformation correctness
9. **Task 9** (Limits type map) — data integrity guard
10. **Task 10** (Provider & schemas) — structural regression tests

## Conventions

- **File naming**: unit test files use `_test.go` suffix. For files that already have acceptance tests (e.g., `tags_test.go`), add unit tests to the same file. For files with no test file or only acceptance tests, create a new `*_unit_test.go` or add to the existing `*_test.go`.
- **Test naming**: `Test<FunctionName>_<scenario>` — no `Acc` prefix (that's reserved for acceptance tests).
- **No infrastructure**: Every test in this plan runs with `go test ./cloudstack/ -run 'Test[^A]'` — no CloudStack simulator needed.
- **Table-driven tests**: Preferred where there are multiple input/output cases for the same function (follow the existing `TestDiffTags` pattern).
- **Package**: All tests are in `package cloudstack` (same-package testing to access unexported functions).

## Expected Outcome

| Metric | Before | After |
|--------|--------|-------|
| Unit tests | 3 | ~55 |
| Components with unit coverage | 1 (tags) | 10 |
| `go test` time (no simulator) | 0.01s | < 1s |

All tests runnable via:
```sh
export PATH="/usr/local/go/bin:$PATH"
go test ./cloudstack/ -run 'Test[^A]' -v -count=1
```
