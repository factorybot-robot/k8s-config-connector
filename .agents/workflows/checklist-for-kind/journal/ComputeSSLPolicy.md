# Migration Journal: ComputeSSLPolicy

This journal tracks the migration of `ComputeSSLPolicy` to a direct controller at `v1beta1`.

## Current Step
Step 1: Direct API Types - Approved & Queued for Merge

## Migration Progress

| Step | Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|------|------|--------------|---------------------|--------|--------------|----------------|
| 1 | Direct API Types | [#9724](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9724) | [#9725](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9725) | Rebase Required | 2026-06-11 | - |
| 2 | Identity and Reference Types Pattern | - | - | Pending | - | - |
| 3 | Create a Round-Trip KRM Fuzzer | - | - | Pending | - | - |
| 4 | Implement Direct Controller & E2E Fixtures | - | - | Pending | - | - |

## Status Update Notes

### 2026-06-12
- Monitored merge queue progress and identified that both merge queue workflow runs `27392322062` and `27391517383` failed due to an unrelated, flaky infrastructure test timeout in `fields/management/gkehub/featuremembership/set_unset`.
- Consequently, PR #9725 was kicked out of the merge queue. Although auto-merge was re-enabled by `barney-s`, the PR remains behind the newly merged upstream `master` commit `1cdcf7237e` (#9312) and is not re-queueing automatically.
- Commented on PR #9725 requesting `factorybot-robot` to rebase the branch on latest upstream `master` to resolve the stale state, trigger fresh green checks, and re-enqueue the PR into the merge queue. Reassigned the PR back to `factorybot-robot`.
- Re-verified the merge queue workflow runs `27392322062` and `27391517383` on June 12, 2026. Verified that all completed check-runs (such as `test-mockgcp`, `unit-tests`, `validate-generated-files`) have passed with 100% success (0 failures), and the remaining fixtures/fuzzing tests are progressing cleanly in the queue.
- Since PR #9725 is not yet merged, we continue to wait for the merge queue to complete and merge the PR into upstream `master` before starting Step 2 (Identity and Reference Types Pattern).
- Monitored the merge queue and check-runs for PR #9725. Verified that the PR is actively being validated inside the GitHub Actions native merge queue on merge commits `4caf565444` (workflow run `27391517383`) and `1cdcf7237e` (workflow run `27392322062`). All active workflows (including `Presubmit` and `ci-presubmit`) are running and progressing cleanly with 0 failures.
- Checked the PR merge status on GitHub. PR #9725 remains in `open` state, fully approved with `/lgtm` and `/approve` from maintainers, and its own presubmits are completely green.
- We continue to wait for the GitHub Actions native merge queue to successfully complete and merge PR #9725 into upstream `master` before starting Step 2 (Identity and Reference Types Pattern).
- Checked the GitHub merge queue status. Verified that PR #9725 is actively being validated inside the GitHub Actions native merge queue (workflow run `27388212871`) on the merge commit `4caf565444e183731a7ba752862d7e79c17c54d8`. Almost all core test suites have successfully completed with 0 failures, and only `tests-e2e-fixtures-bigquery` and `tests-e2e-fixtures-sql` are currently in-progress. We continue to wait for the merge queue to successfully complete and merge PR #9725 into upstream `master` before starting Step 2 (Identity and Reference Types Pattern).
- Monitored the merge queue and check-runs for PR #9725 again. Verified that the head commit `e9ab8142fe` has all check-runs successfully completed/skipped, with no failing checks on the branch itself.
- Verified that the merge queue branch commit `1cdcf7237e` for PR #9725 is progressing with 0 failures across all completed check-runs (including E2E, unit tests, and validations), while some remaining tasks are still queued.
- Confirmed that the preceding PR #9720 in the merge queue (commit `4caf565444`) is also progressing cleanly with no failures.
- PR #9725 is in a fully green, approved, and healthy state within the native merge queue. We continue to wait for the merge queue to complete the merge of PR #9725 into upstream `master` before we can start Step 2 (Identity and Reference Types Pattern).
- Monitored PR #9725 checks. The PR has been approved with `/lgtm` and `/approve` by the owners/approvers (maqiuyujoyce and barney-s). All 21 core CI check-runs have passed successfully.
- The PR is now queued to merge by the Prow/Tide bot. We are waiting for the merge to complete before proceeding to Step 2 (Identity and Reference Types Pattern).
- Verified that `factorybot-robot` ran `make resource-docs` and committed/pushed the updated documentation in commit `1dc8ae1ef475ae80798dab61b6443910ef4e23d0`. All core CI checks are in a fully green state.
- Re-verified that the PR remains healthy and is actively in the merge queue (behind PR #9720) with no blocks, and we are continuing to monitor it until it merges.
- Confirmed that the GitHub Actions Presubmit run for the blocking PR #9720 is actively running and making progress. PR #9725 remains cleanly approved, fully green, and queued directly behind #9720 in the merge queue.
- Monitored merge queue progress for PR #9725. Verified that the PR is currently inside the GitHub Actions native merge queue on the merge commit `4caf565444`. 17 core test suites have successfully completed on this merge commit with 0 failures, and the remaining tests are healthy and progressing. We are waiting for the merge queue to complete and merge PR #9725 into upstream `master` before starting Step 2.
- Checked the GitHub merge queue status again. Verified that PR #9725 is at the front of the queue under the merge commit '4caf565444', and the blocking PR #9720 has been queued on top of it. All completed tests on both merge queue commits continue to pass cleanly, and the remaining exhaustive test suites are queued. We continue to wait for the merge queue to complete and merge PR #9725 before initiating Step 2 (Identity and Reference Types Pattern).

### 2026-06-11
- Initialized the ComputeSSLPolicy migration journal.
- Opened GitHub Issue #9724 for Step 1 (Direct API Types) and assigned it to `factorybot-robot`.
- Checked progress; issue #9724 is active and we are waiting for `factorybot-robot` to create the implementation PR.
- Re-verified status: GitHub Issue #9724 is still open and assigned to `factorybot-robot`. The implementation PR for Step 1 has not yet been submitted. We are continuing to monitor progress on the active sandbox work.
- Detected that `factorybot-robot` submitted PR #9725. However, the `validate-generated-files` check failed.
- Diagnosed the failure locally and found that the resource reference documentation (`scripts/generate-google3-docs/resource-reference/generated/resource-docs/compute/computesslpolicy.md`) was out of date.
- Commented on PR #9725 requesting `factorybot-robot` to run `make resource-docs` to regenerate the documentation and update the PR, then reassigned the PR back to them.
- Monitored PR #9725 progress. Discovered that the latest commit `6e3e981253f4da13160740644c820b63f7b9cb35` is failing all CI check-runs due to a compilation error.
- Diagnosed the failure locally by checking out the PR branch. Found that `6e3e981253` incorrectly added `DeepCopy` methods back to `apis/networkconnectivity/v1alpha1/zz_generated.deepcopy.go` for commented-out/unreachable types (`AllocationOptions`, `Any`, and `AutoCreatedSubnetworkInfo`), causing `undefined: AllocationOptions` and other compilation errors.
- Commented on PR #9725 requesting `factorybot-robot` to discard these changes or regenerate cleanly and re-assigned it back to them.
- Detected that `factorybot-robot` pushed a new commit `eb2e4b4a24` which successfully pruned the undefined deepcopy methods and fixed the compilation failures (unit-tests passed).
- However, the `validate-generated-files` check continues to fail for PR #9725 because they did not commit/stage the regenerated documentation `scripts/generate-google3-docs/resource-reference/generated/resource-docs/compute/computesslpolicy.md` after updating the schema.
- Commented on PR #9725 requesting them to run `make resource-docs` locally, commit/push the changes to `computesslpolicy.md`, and re-assigned to `factorybot-robot`.
