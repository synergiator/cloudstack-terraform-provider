# AGENTS.md

Guide for AI agents working in the Apache CloudStack Terraform Provider codebase.

## Project Overview

This is an **Apache-licensed Terraform provider** for [Apache CloudStack](https://cloudstack.apache.org/), written in Go. It enables managing CloudStack infrastructure (VMs, networks, VPCs, storage, etc.) via Terraform and OpenTofu.

- **Module**: `github.com/terraform-providers/terraform-provider-cloudstack`
- **Go version**: 1.23+ (specified in `go.mod`)
- **CloudStack Go SDK**: `github.com/apache/cloudstack-go/v2`
- **Registry**: `registry.terraform.io/cloudstack/cloudstack`

## Commands

| Task | Command |
|---|---|
| **Build** | `make build` (runs `fmtcheck` then `go install`) |
| **Unit tests** | `make test` (runs `fmtcheck` + `go test` with 30s timeout) |
| **Acceptance tests** | `make testacc` (requires running CloudStack; 30m timeout) |
| **Single acceptance test** | `make testacc TESTARGS='-run ^TestAccCloudStackVPC_basic$'` |
| **Format code** | `make fmt` (runs `gofmt -w`) |
| **Format check** | `make fmtcheck` |
| **Vet** | `make vet` |
| **Error check** | `make errcheck` (requires `errcheck` tool) |
| **Compile tests only** | `make test-compile TEST=./cloudstack` |

### Acceptance Test Environment Variables

Acceptance tests (`TestAcc*`) require a running CloudStack instance (typically the simulator):

```sh
export CLOUDSTACK_API_URL=http://localhost:8080/client/api
export CLOUDSTACK_API_KEY=<your-api-key>
export CLOUDSTACK_SECRET_KEY=<your-secret-key>
```

Optional: `CLOUDSTACK_TEMPLATE_URL` (used by some template tests).

### Local CloudStack Simulator

```sh
docker pull apache/cloudstack-simulator:4.20.1.0
docker run --name simulator -p 8080:5050 -d apache/cloudstack-simulator:4.20.1.0
docker exec -it simulator python /root/tools/marvin/marvin/deployDataCenter.py -i /root/setup/dev/advanced.cfg
```

Acceptance tests also expect a project named `terraform` to exist in CloudStack.

## Code Organization

```
.
├── main.go                       # Entrypoint - mux server combining SDK v2 + Framework providers
├── cloudstack/
│   ├── provider.go               # SDK v2 provider (most resources/data sources)
│   ├── provider_v6.go            # Plugin Framework provider (newer resources)
│   ├── provider_test.go          # Provider tests + mux server test + testAccPreCheck
│   ├── config.go                 # Config struct, CloudStack client creation
│   ├── resources.go              # Shared helpers: retrieveID, setValueOrID, Retry, ResourceWithConfigure
│   ├── tags.go                   # Tag CRUD helpers (setTags, updateTags, diffTags)
│   ├── metadata.go               # Metadata CRUD helpers (setMetadata, updateMetadata)
│   ├── data_source_cloudstack_common_schema.go  # Shared filter schema for data sources
│   │
│   ├── resource_cloudstack_*.go          # SDK v2 resource implementations
│   ├── resource_cloudstack_*_test.go     # Acceptance tests for resources
│   ├── data_source_cloudstack_*.go       # SDK v2 data source implementations
│   ├── data_source_cloudstack_*_test.go  # Acceptance tests for data sources
│   │
│   ├── service_offering_*.go             # Framework-based service offering resources
│   │   ├── service_offering_schema.go        # Shared schema via serviceOfferingMergeCommonSchema()
│   │   ├── service_offering_models.go        # Model structs with tfsdk tags
│   │   ├── service_offering_util.go          # Shared CRUD helper methods on model structs
│   │   ├── service_offering_constrained_resource.go
│   │   ├── service_offering_unconstrained_resource.go
│   │   └── service_offering_fixed_resource.go
│   └── ...
├── website/docs/
│   ├── r/          # Resource documentation (HTML markdown)
│   └── d/          # Data source documentation (HTML markdown)
├── scripts/        # Helper scripts (gofmtcheck, errcheck, gogetcookie, changelog-links)
├── GNUmakefile     # Build/test/lint targets
└── .github/
    ├── workflows/
    │   ├── build.yml       # CI: make build + make test
    │   ├── acceptance.yml  # CI: acceptance tests matrix (Terraform + OpenTofu × CloudStack versions)
    │   └── rat.yaml        # Apache RAT license header check
    └── actions/
        └── setup-cloudstack/action.yml  # Composite action for CloudStack simulator setup
```

## Two Provider Frameworks (Mux Server)

This provider uses **both** the Terraform Plugin SDK v2 and the Terraform Plugin Framework, combined via `tf6muxserver`:

1. **SDK v2** (`provider.go`): Contains the majority of resources and all data sources. Uses `schema.Resource` with `Create`/`Read`/`Update`/`Delete` function fields.
2. **Plugin Framework** (`provider_v6.go`): Contains only the service offering resources (`cloudstack_service_offering_constrained`, `cloudstack_service_offering_unconstrained`, `cloudstack_service_offering_fixed`). Uses `resource.Resource` interface with `Schema`/`Create`/`Read`/`Update`/`Delete` methods.

The `main.go` upgrades the SDK v2 provider to protocol 6 via `tf5to6server`, then muxes it with the Framework provider.

**When adding new resources**: New resources should use the Plugin Framework pattern (see `service_offering_*` files). Register them in `provider_v6.go` → `Resources()`.

**When modifying existing resources**: Most existing resources use SDK v2. Follow the existing patterns in that resource file.

## Resource Patterns (SDK v2)

Most resources follow this pattern:

### File Naming
- Resource: `resource_cloudstack_<name>.go`
- Data source: `data_source_cloudstack_<name>.go`
- Tests: append `_test.go` to either

### Resource Structure
```go
func resourceCloudStack<Name>() *schema.Resource {
    return &schema.Resource{
        Create: resourceCloudStack<Name>Create,
        Read:   resourceCloudStack<Name>Read,
        Update: resourceCloudStack<Name>Update,
        Delete: resourceCloudStack<Name>Delete,
        Importer: &schema.ResourceImporter{
            State: importStatePassthrough,  // supports project/ID import format
        },
        Schema: map[string]*schema.Schema{...},
    }
}
```

### Key Patterns
- **Client access**: Cast meta to `*cloudstack.CloudStackClient` → `cs := meta.(*cloudstack.CloudStackClient)`
- **ID resolution**: Use `retrieveID(cs, "type", value)` to resolve names to IDs (supports zone, domain, project, service_offering, etc.)
- **Template ID**: Use `retrieveTemplateID(cs, zoneid, value)` specifically for templates
- **Project support**: Use `setProjectid(p, cs, d)` to set project ID on API params
- **Tags**: Use `tagsSchema()` for schema, `setTags(cs, d, "ResourceType")` on create, `updateTags(cs, d, "ResourceType")` on update
- **Metadata**: Use `metadataSchema()` for schema, `setMetadata()`/`updateMetadata()`/`getMetadata()` for CRUD
- **Name-or-ID fields**: Use `setValueOrID(d, key, name, id)` to preserve whether user specified name or ID
- **Read not-found**: When `count == 0`, clear ID with `d.SetId("")` and return nil (don't error)
- **Delete not-found**: Check for "entity does not exist" in error message, return nil (idempotent delete)
- **Import**: Use `importStatePassthrough` which supports `project/id` format

### Registration
Resources are registered in `provider.go` → `ResourcesMap` and `DataSourcesMap`.

## Resource Patterns (Plugin Framework)

Used only for `service_offering_*` resources:

- Struct implements `resource.Resource` and `resource.ResourceWithConfigure`
- Configuration via `ResourceWithConfigure` base struct in `resources.go`
- Models use `tfsdk` struct tags
- Schema merging via `serviceOfferingMergeCommonSchema()` for shared attributes
- Shared CRUD logic on model structs (e.g., `commonCreateParams`, `commonRead`, `commonUpdate`)
- Interface compliance via `var _ resource.Resource = &myResource{}`

## Data Source Patterns

- All data sources use SDK v2
- Use `dataSourceFiltersSchema()` for the standard `filter` block (name/value regex matching)
- Filter matching: JSON-marshal the API response, then regex-match against fields
- When multiple results match, return the most recently created one
- Function naming: `dataSourceCloudstack<Name>()` (lowercase 's' in 'stack' — inconsistent with resources)

## Testing Patterns

### Test Structure
```go
func TestAccCloudStack<Name>_basic(t *testing.T) {
    var obj cloudstack.<Type>
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,       // SDK v2 resources
        // OR:
        ProtoV6ProviderFactories: testAccMuxProvider,  // Framework resources
        CheckDestroy: testAccCheckCloudStack<Name>Destroy,
        Steps: []resource.TestStep{
            {
                Config: testAccCloudStack<Name>_basic,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckCloudStack<Name>Exists("cloudstack_<name>.foo", &obj),
                    testAccCheckCloudStack<Name>Attributes(&obj),
                ),
            },
        },
    })
}
```

### Test Conventions
- Test configs are `const` string literals at the bottom of the test file (raw HCL)
- Test resource addresses use `.foo` as the instance name
- `testAccPreCheck(t)` validates required env vars are set
- Exists checks: look up resource in state, fetch from API, verify IDs match
- Destroy checks: iterate all resources of the type, verify API returns not-found
- Import tests: separate `TestAccCloudStack<Name>_import` functions using `ImportState`/`ImportStateVerify`
- CloudStack simulator zone is `"Sandbox-simulator"`

## Naming Conventions

| Item | Convention | Example |
|---|---|---|
| Resource constructor | `resourceCloudStack<PascalName>()` | `resourceCloudStackVPC()` |
| Data source constructor | `dataSourceCloudstack<PascalName>()` | `dataSourceCloudstackVPC()` |
| CRUD functions (SDK v2) | `resourceCloudStack<Name><CRUD>` | `resourceCloudStackVPCCreate` |
| Framework resource struct | `serviceOffering<Variant>Resource` | `serviceOfferingConstrainedResource` |
| Framework constructor | `Newservice<Variant>Resource` | `NewserviceOfferingConstrainedResource` |
| Model structs | `<name>ResourceModel` or `<Name>Model` | `serviceOfferingConstrainedResourceModel` |
| Test functions | `TestAccCloudStack<Name>_<variant>` | `TestAccCloudStackVPC_basic` |
| Terraform resource type | `cloudstack_<snake_name>` | `cloudstack_vpc` |

**Note**: There is inconsistency in capitalization of "CloudStack" vs "Cloudstack" vs "cloudstack" across the codebase. Resource constructors use `CloudStack` (capital S), data source constructors sometimes use `Cloudstack` (lowercase s). Follow the pattern in the specific file you're editing.

## CI/CD

### GitHub Actions Workflows

1. **Build-Check** (`build.yml`): Runs on every push/PR — `make build` + `make test`
2. **Acceptance Test** (`acceptance.yml`): Runs on every push/PR — matrix of:
   - Terraform versions: 1.11.x, 1.12.x
   - OpenTofu versions: 1.8.x, 1.9.x
   - CloudStack versions: 4.19.0.1, 4.19.1.3, 4.19.2.0, 4.19.3.0, 4.20.1.0
   - Uses `apache/cloudstack-simulator` Docker image as a service container
3. **RAT Check** (`rat.yaml`): Apache Release Audit Tool — checks license headers on all files

### Release
- Uses GoReleaser (`.goreleaser.yml`) for cross-platform binary builds
- Apache release process via `performrelease.sh` (source tarball, GPG signing, SHA-512 checksums)

## License Headers

**Every source file must have the Apache License 2.0 header.** The RAT CI check enforces this. Files excluded from RAT checking are listed in `.rat-excludes`.

Standard Go file header:
```go
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
```

## Documentation

- Provider docs live in `website/docs/`
- Resource docs: `website/docs/r/<name>.html.markdown`
- Data source docs: `website/docs/d/<name>.html.markdown`
- Format: YAML frontmatter + markdown with HCL example usage, argument reference, attributes reference, and import section
- When adding/modifying a resource, update the corresponding doc file

## Important Gotchas

1. **Two frameworks in one provider**: Read `main.go` to understand the mux setup. SDK v2 and Framework resources coexist. Tests for SDK v2 resources use `Providers: testAccProviders`; tests for Framework resources use `ProtoV6ProviderFactories: testAccMuxProvider`.

2. **`go test -i` is deprecated**: The `make test` target uses `go test -i` which is a no-op in modern Go but still works. Don't be confused by it.

3. **Acceptance tests need real infrastructure**: All `TestAcc*` tests talk to a CloudStack API. They cannot run without `CLOUDSTACK_API_URL`, `CLOUDSTACK_API_KEY`, and `CLOUDSTACK_SECRET_KEY` set. Use the simulator for local testing.

4. **Import format with projects**: The `importStatePassthrough` helper splits on `/` to extract optional project name. Import format is `project_name/resource_id` or just `resource_id`.

5. **Name-or-ID pattern**: Many resources accept either a name or UUID for reference fields (zone, vpc_offering, etc.). The `retrieveID` helper resolves names to IDs, and `setValueOrID` preserves the user's original format on read-back.

6. **CloudStack SDK field naming**: The `cloudstack-go` SDK uses lowercase concatenated field names (e.g., `Displaytext`, `Vpcofferingid`, `Networkdomain`). These don't follow Go conventions but match CloudStack API response fields.

7. **Error handling in deletes**: Delete functions check for "entity does not exist" error strings to handle already-deleted resources gracefully. This is a workaround for the CloudStack API's error reporting.

8. **`fmtcheck` runs before build**: The `build` target depends on `fmtcheck`. Always run `make fmt` before building if you've made changes.

9. **Package is `cloudstack`**: All provider code is in the single `cloudstack` package. There are no sub-packages.

10. **Service offering resources are special**: They are the only resources using the Plugin Framework and have a different code organization with shared schema/models/utils split across multiple files.
