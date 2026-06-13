# Migration Journal: ComputeFirewallPolicyAssociation

## Current Step
Step 4: Implement Direct Controller & E2E Fixtures (Completed)

## Progress Tracking

| Step | Name | Issue | Pull Request | Status | Date Started | Date Completed |
|------|------|-------|--------------|--------|--------------|----------------|
| 1    | Direct API Types | [#9976](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9976) | [#10068](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10068) | Completed | 2026-06-13 | 2026-06-13 |
| 2    | Identity & Reference Types | [#10110](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/10110) | N/A (Inlined) | Completed | 2026-06-13 | 2026-06-13 |
| 3    | Create Round-Trip KRM Fuzzer | [#9976](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/9976) | [#10068](https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10068) | Completed | 2026-06-13 | 2026-06-13 |
| 4    | Implement Direct Controller & E2E | [#10110](https://github.com/GoogleCloudPlatform/k8s-config-connector/issues/10110) | Local Branch | Completed | 2026-06-13 | 2026-06-13 |

## Learnings & Observations
- **Resource Scope:** `ComputeFirewallPolicyAssociation` is folder- or organization-level. Reference types like `attachmentTargetRef` were hand-coded locally to support referencing Folders and Organizations.
- **Mocking Support:** Added MockGCP handlers for `GetAssociation`, `AddAssociation`, and `RemoveAssociation` in `mockcompute/firewallpoliciesv1.go`.
- **E2E Golden Files:** Successfully generated and verified E2E golden files against MockGCP for both Folder and Organization scoped associations.
