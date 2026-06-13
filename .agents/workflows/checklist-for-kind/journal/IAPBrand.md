# Migration Journal: IAPBrand

**Current Step**: Step 1: Direct API Types

## Progress Tracking

| Step Number and Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|---|---|---|---|---|---|
| Step 1: Direct API Types | [#9739](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9739) | [#9795](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9795), [#9809](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9809) | PR Created | 2026-06-12 | |
| Step 2: Identity and Reference Types Pattern | | | Pending | | |
| Step 3: Create a Round-Trip KRM Fuzzer | | | Pending | | |
| Step 4: Implement Direct Controller & E2E Fixtures | | | Pending | | |

## Status Update Notes

### 2026-06-13
- Monitored Step 1 progress: Checked the status of open pull requests [#9795](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9795) and [#9809](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9809). Confirmed that the automated CI failure resolver (`factorybot-robot`) has successfully investigated the failures (`validate-generated-files`, `unit-tests`, and `validations`), generated comprehensive root-cause analysis reports, and pushed/amended the required fixes (such as regenerated clients, resource docs, and API check exceptions) to both branches. Currently waiting for the clean set of CI checks to complete on both PRs.
- Monitored Step 1 progress: Checked CI check status for the new pull request [#9809](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9809) created by `factorybot-robot`. Confirmed that the new PR was successfully opened and all CI checks are currently queued and pending.
- Monitored Step 1 progress: Checked CI check status for pull request [#9795](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9795). Found that CI checks (`validate-generated-files`, `unit-tests`, and `validations`) are failing, and the bot previously gave up. Re-assigned the PR to `factorybot-robot` to request automated retry and resolution of CI failures.

### 2026-06-12
- Monitored Step 1 progress: Re-verified and successfully assigned PR [#9795](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9795) to `factorybot-robot` via the GitHub REST API to trigger the automated resolver for failing CI checks (`validate-generated-files`, `unit-tests`, and `validations`).
- Monitored Step 1 progress: Checked CI check status for pull request [#9795](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9795). Found that several checks (including `unit-tests`, `validate-generated-files`, and `validations`) failed. Assigned the PR to `factorybot-robot` to trigger automated resolution of these CI failures.
- Monitored Step 1 progress: Pull request [#9795](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9795) was successfully created by `factorybot-robot` to implement direct KRM types and `generate.sh` for IAPBrand. The PR is currently open and CI checks are running.
- Created Step 1 GitHub issue [#9739](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9739) for factorybot-robot to implement direct KRM types and `generate.sh` for IAPBrand.
- Checked progress of Step 1: issue #9739 is currently being worked on by `factorybot-robot` in a sandbox; awaiting pull request creation.
- Re-verified status: `factorybot-robot` has commented in issue #9739 that it started fixing the issue in a sandbox; awaiting pull request creation.
- Re-checked progress of Step 1: `factorybot-robot` is continuing to work on implementing the direct KRM types and `generate.sh` for IAPBrand in its sandbox; no Pull Request has been created yet.
- Confirmed that issue #9739 remains open with no Pull Request created yet. `factorybot-robot` is still working on implementing the direct KRM types and `generate.sh` in the sandbox.
- Monitored Step 1: Issue #9739 remains open and no pull request has been submitted by `factorybot-robot` yet.
- Checked progress of Step 1: Confirmed that GitHub issue #9739 remains open with no Pull Request submitted yet. `factorybot-robot` is still working on implementing the direct KRM types and `generate.sh` in its sandbox.
- Periodic check on Step 1: Verified that issue #9739 is still open and active. `factorybot-robot` continues working on the direct KRM types in the sandbox; no pull request has been opened yet.
- Checked progress of Step 1: Verified that Issue #9739 remains open with no Pull Request opened yet. `factorybot-robot` is still working on implementing direct KRM types in its sandbox.
- Periodic check on Step 1: Checked progress on GitHub issue #9739. Confirmed it is still open with no Pull Request submitted yet. `factorybot-robot` continues working on implementing the direct KRM types and `generate.sh` in its sandbox.
- Monitored Step 1 progress: Verified that GitHub issue #9739 is still open with no Pull Request created yet. `factorybot-robot` is still working on implementing the direct KRM types and `generate.sh` for IAPBrand in its sandbox.
- Monitored Step 1 progress: Checked GitHub issue #9739. Verified that the issue remains open with no active pull request opened yet, and that `factorybot-robot` is still working on implementing the direct KRM types and `generate.sh` in its sandbox.
- Verified Step 1 status: Checked GitHub issue #9739. Confirmed that no pull request has been submitted yet and `factorybot-robot` is continuing to work on implementing direct KRM types for IAPBrand in its sandbox.
- Monitored Step 1 progress: Checked GitHub issue #9739. Confirmed it is still open and `factorybot-robot` is working on direct KRM types in its sandbox; no Pull Request has been opened yet.
