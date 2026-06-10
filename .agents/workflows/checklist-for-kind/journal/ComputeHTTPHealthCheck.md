# Migration Journal: ComputeHTTPHealthCheck

## Current Step
**Step 1: Direct API Types**

The current step is to ensure that the direct API types are generated and matching the CRD. There is an existing open Pull Request #7491 addressing this step, but it currently has merge conflicts and is in a `dirty` state. We have requested `factorybot-robot` to rebase and resolve the conflicts.

## Progress Table

| Step | Step Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|------|-----------|--------------|---------------------|--------|--------------|----------------|
| 1 | Direct API Types | [#7480](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/7480) | [#7491](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7491) | PR Created (Conflict - Rebase Requested) | 2026-06-10 | - |
| 2 | Identity & Reference Pattern | - | - | - | - | - |
| 3 | Round-Trip KRM Fuzzer | - | - | - | - | - |
| 4 | Direct Controller & E2E Fixtures | - | - | - | - | - |

## Notes & Observations
* **Step 1 Status:** PR #7491 is open and its checks have passed previously. However, it currently has merge conflicts. We have triggered a rebase/conflict resolution by commenting on the PR and assigning `factorybot-robot`.
* **Step 2 Status:** Pending merge of Step 1.
* **Step 3 Status:** Pending completion of Step 2.
* **Step 4 Status:** Pending completion of Step 3.
