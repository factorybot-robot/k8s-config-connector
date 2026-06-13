# Migration Journal: ComputeAutoscaler

## Current Step
**Step 1: Direct API Types**

The current focus is on implementing direct KRM types and ensuring schema compatibility with the existing CRD. PR #10046 is open with all checks passing, awaiting review and merger.

## Migration Progress Tracking

| Step | Step Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|------|-----------|--------------|---------------------|--------|--------------|----------------|
| 1 | Direct API Types | [#9956](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9956) | [#10046](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10046) | PR Created | 2026-06-13 | |
| 2 | Identity and Reference Types Pattern | | | | | |
| 3 | Create a Round-Trip KRM Fuzzer | | | | | |
| 4 | Implement Direct Controller & E2E Fixtures | | | | | |

## Status Updates
* **2026-06-13**: Monitored Step 1 PR #10046. All CI checks are green and passing. Awaiting merge before proceeding to Step 2.
