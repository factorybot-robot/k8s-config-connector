# Migration Journal: ComputeSSLPolicy

This journal tracks the migration of `ComputeSSLPolicy` to a direct controller at `v1beta1`.

## Current Step
Step 1: Direct API Types - PR Created (Failing CI, Changes Requested)

## Migration Progress

| Step | Name | GitHub Issue | GitHub Pull Request | Status | Date Started | Date Completed |
|------|------|--------------|---------------------|--------|--------------|----------------|
| 1 | Direct API Types | [#9724](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9724) | [#9725](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9725) | PR Created | 2026-06-11 | - |
| 2 | Identity and Reference Types Pattern | - | - | Pending | - | - |
| 3 | Create a Round-Trip KRM Fuzzer | - | - | Pending | - | - |
| 4 | Implement Direct Controller & E2E Fixtures | - | - | Pending | - | - |

## Status Update Notes

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
