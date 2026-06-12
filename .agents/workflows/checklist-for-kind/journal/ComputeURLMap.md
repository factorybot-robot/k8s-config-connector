# Migration Journal: ComputeURLMap

## Current Step
Step 1: Direct API Types & Step 2: Identity and Reference Types Pattern

## Migration Progress

| Step | Step Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | Direct API Types | [#6438](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/6438) | [#6854](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6854) | PR Created | 2026-05-31 | |
| 2 | Identity and Reference Types Pattern | [#6438](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/6438) | [#6854](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/6854) | PR Created | 2026-05-31 | |
| 3 | Create a Round-Trip KRM Fuzzer | | | Not Started | | |
| 4 | Implement Direct Controller & E2E Fixtures | | | Not Started | | |

## Status Updates
- **2026-06-12**: Initialized migration journal for ComputeURLMap. Identified existing open issue #6438 and open pull request #6854 covering both Step 1 (Direct API Types) and Step 2 (Identity and Reference Types Pattern).
- **2026-06-12**: Observed that PR #6854 currently has merge conflicts with the latest `master` and is failing several CI check-runs (`crd-equivalence-check`, `tests-e2e-fixtures-compute`, `validate-generated-files`).
- **2026-06-12**: Posted a comment on PR #6854 asking `factorybot-robot` to rebase on `master`, investigate and resolve the CI failures, and assigned `factorybot-robot` to continue the implementation.
