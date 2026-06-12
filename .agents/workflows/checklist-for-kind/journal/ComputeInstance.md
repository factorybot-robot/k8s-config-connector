# ComputeInstance Migration Journal

## Current Step
Step 1: Direct API Types

## Progress Tracking
| Step Number and Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
| --- | --- | --- | --- | --- | --- |
| 1. Direct API Types | [#9735](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9735) | [#9741](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9741) | PR Created | 2026-06-12 | |
| 2. Identity and Reference Types Pattern | | | Pending | | |
| 3. Create a Round-Trip KRM Fuzzer | | | Pending | | |
| 4. Implement Direct Controller & E2E Fixtures | | | Pending | | |

## Recent Status Updates
* **2026-06-12**: Investigated failing CI check `validate-generated-files` on PR #9741. Identified that the generated resource reference documentation was out-of-date. Left detailed instructions on PR #9741 for `factorybot-robot` to run `make resource-docs` and re-assigned the PR back to them.
* **2026-06-12**: Identified CI failures in PR #9741. Analyzed root causes (pointer check on non-pointer slices in codegen generator, and schema mismatch for `Tags` message vs `[]string`). Commented on PR #9741 with detailed instructions on manual mapper overrides and reassigned it to `factorybot-robot`.
* **2026-06-12**: Detected and monitored open PR #9741 for Step 1. The pull request has been successfully created and CI checks are currently running.
* **2026-06-12**: Monitored Step 1 progress. Issue #9735 remains open; waiting for the AI Factory sandbox to generate and submit the pull request.
* **2026-06-12**: Initialized migration checklist. Created Step 1 issue #9735 to scaffold direct API types and configure `generate.sh` for `ComputeInstance`.
