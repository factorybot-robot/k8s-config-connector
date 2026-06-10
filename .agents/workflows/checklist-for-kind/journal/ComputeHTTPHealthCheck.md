# Migration Journal: ComputeHTTPHealthCheck

## Current Step
**Step 1: Direct API Types**

The current step is to ensure that the direct API types are generated and matching the CRD. There is an existing open Pull Request #7491 addressing this step, but it currently has merge conflicts and is in a `dirty` state. We have requested `factorybot-robot` to rebase and resolve the conflicts.

## Progress Table

| Step | Step Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|------|-----------|--------------|---------------------|--------|--------------|----------------|
| 1 | Direct API Types | [#7480](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/7480) | [#9676](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9676) | PR Created (Rebase Completed - CI Running) | 2026-06-10 | - |
| 2 | Identity & Reference Pattern | - | - | Pending | - | - |
| 3 | Round-Trip KRM Fuzzer | - | - | Pending | - | - |
| 4 | Direct Controller & E2E Fixtures | - | - | Pending | - | - |

## Notes & Observations
* **Step 1 Status:** PR #9676 is open and rebased successfully. CI is running.
* **Step 2 Status:** Pending merge of Step 1.
* **Step 3 Status:** Pending completion of Step 2.
* **Step 4 Status:** Pending completion of Step 3.
