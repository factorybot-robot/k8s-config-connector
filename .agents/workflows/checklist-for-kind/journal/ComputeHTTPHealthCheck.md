# Migration Journal: ComputeHTTPHealthCheck

**Current Step:** `Step 1: Direct API Types` (CI failing, changes requested)

| Step | Step Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|---|---|---|---|---|---|---|
| 1 | Direct API Types | [#7480](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/7480) | [#9676](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9676) | CI Failing | 2026-06-10 | - |
| 2 | Identity & Reference Pattern | - | - | Pending | - | - |
| 3 | Round-Trip KRM Fuzzer | - | - | Pending | - | - |
| 4 | Direct Controller & E2E Fixtures | - | - | Pending | - | - |

## Notes & Status Updates

- **2026-06-10:** Step 1 initiated. PR #9676 created by `factorybot-robot`.
- **2026-06-11:** CI validation checks (`validations` and `validate-generated-files`) are failing on PR #9676. The resource Go clients are out-of-date and the branch needs to be rebased on `master` to resolve out-of-sync files. Reassigned PR back to `factorybot-robot` to apply fixes.
