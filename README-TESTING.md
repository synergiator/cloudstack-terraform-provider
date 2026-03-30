# Testing Guide

This repository has two main test paths:

- `make test`: local validation (format check + Go test run)
- `make testacc`: acceptance tests against a real CloudStack API endpoint (typically simulator)

## 1) Local Validation

Run:

```sh
make test
```

What this does:

- Runs `fmtcheck`
- Runs `go test` across packages with a short timeout

Use this as your fast default before commits.

## 2) Acceptance Tests

Run:

```sh
make testacc
```

Acceptance tests require these environment variables:

```sh
export CLOUDSTACK_API_URL=http://localhost:8080/client/api
export CLOUDSTACK_API_KEY=<your-api-key>
export CLOUDSTACK_SECRET_KEY=<your-secret-key>
```

If any are missing, tests fail during pre-check.

### Run one acceptance test

```sh
make testacc TESTARGS='-run ^TestAccCloudStackNetworkACLRule_update$'
```

## 3) CloudStack Simulator Setup (local)

Example flow:

```sh
docker pull apache/cloudstack-simulator:4.20.1.0
docker run --name simulator -p 8080:5050 -d apache/cloudstack-simulator:4.20.1.0
docker exec -it simulator python /root/tools/marvin/marvin/deployDataCenter.py -i /root/setup/dev/advanced.cfg
```

Then retrieve or create API credentials in CloudStack and export them as shown above.

Note: acceptance tests also expect an empty project named `terraform`.

## 4) Do I Need a `.cmk` Profile?

Short answer:

- Not required for `make test`
- Not required for `make testacc` if env vars are already exported
- Required only when you use `cmk -p <profile> ...` commands

In CI, a `.cmk` profile is used to bootstrap API keys and create setup resources. Locally, you can skip `.cmk` entirely if you already have API URL/key/secret in environment variables.

## 5) Running Acceptance Tests Against a Real CloudStack

Acceptance tests work against **any** CloudStack deployment — not just the simulator. Point the three environment variables at your real API:

```sh
export CLOUDSTACK_API_URL=https://your-cloudstack-host:8080/client/api
export CLOUDSTACK_API_KEY=<your-api-key>
export CLOUDSTACK_SECRET_KEY=<your-secret-key>
```

Then run `make testacc` as usual.

### Hardcoded test assumptions

The test suite has simulator-specific names hardcoded in HCL config literals. Your environment must either provide matching resources or you must patch the test files.

| Hardcoded Value | Occurrences | Type |
|---|---|---|
| `Sandbox-simulator` | ~93 | Zone name |
| `DefaultIsolatedNetworkOfferingWithSourceNatService` | ~39 | Network offering |
| `terraform` | ~27 | Project name |
| `Small Instance` | ~26 | Service offering |
| `Default VPC offering` | ~22 | VPC offering |
| `Medium Instance` | ~11 | Service offering |

Some template tests also require `CLOUDSTACK_TEMPLATE_URL` to be set.

### Approaches

1. **Pre-create matching resources** in your CloudStack so the hardcoded names resolve. This is the least-invasive option.

2. **Find-and-replace** the hardcoded values in test files to match your environment:

   ```sh
   sed -i 's/Sandbox-simulator/YourZoneName/g' cloudstack/*_test.go
   sed -i 's/Small Instance/YourSmallOffering/g' cloudstack/*_test.go
   sed -i 's/Medium Instance/YourMediumOffering/g' cloudstack/*_test.go
   sed -i 's/Default VPC offering/YourVPCOffering/g' cloudstack/*_test.go
   sed -i 's/DefaultIsolatedNetworkOfferingWithSourceNatService/YourNetworkOffering/g' cloudstack/*_test.go
   ```

3. **Run a subset** of tests you have prepared resources for:

   ```sh
   make testacc TESTARGS='-run ^TestAccCloudStackVPC_basic$'
   ```

### Important warnings

- Tests **create and destroy real resources** (VMs, networks, VPCs, volumes, etc.). Always run against a non-production environment.
- There is no env-var or config mechanism to override these values without editing the test files — they are all hardcoded string literals in the HCL test configs.

## 6) What Test Types Exist?

Most tests in this repo are acceptance-style (`TestAcc...`), plus a small number of non-acceptance unit-style tests.

