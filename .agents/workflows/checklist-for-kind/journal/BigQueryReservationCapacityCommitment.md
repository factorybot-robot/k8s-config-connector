# Migration Journal: BigQueryReservationCapacityCommitment

This journal tracks the progress of migrating `BigQueryReservationCapacityCommitment` to a production-ready direct controller.

## Current Status
* **Current Step**: Step 4: Implement Direct Controller & E2E Fixtures
* **Overall Progress**: 4 / 4 steps in progress (Step 4 is currently open as PR #9577 and all CI checks are green/passing)

## Progress Tracking

| Step | Step Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
| :---: | :--- | :---: | :---: | :---: | :---: | :---: |
| 1 | Direct API Types | [#9425](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9425) | [#9431](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9431) | `Completed` | 2026-06-06 | 2026-06-07 |
| 2 | Identity & Reference Types | [#9507](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9507) | [#9510](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9510) | `Completed` | 2026-06-07 | 2026-06-07 |
| 3 | Round-Trip KRM Fuzzer | [#9524](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9524) | [#9528](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9528) | `Completed` | 2026-06-07 | 2026-06-07 |
| 4 | Direct Controller & E2E Fixtures | [#9562](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9562) | [#9577](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9577) | `PR Created` | 2026-06-08 | *In progress* |

## Status Updates
- **2026-06-13**: Verified that Step 4 PR (#9577) is currently open and has successfully passed all CI check-runs. The PR is currently waiting for human/OWNER review. Created local journal and updated migration progress on the parent tracking issue #10100.
