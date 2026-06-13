# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import json
import os
import datetime

COORDINATOR_ISSUE_NUMBER = 10099

def main():
    data_path = 'dev/migration-tracker/data.json'
    if not os.path.exists(data_path):
        print("data.json not found!")
        return

    with open(data_path, 'r') as f:
        data = json.load(f)

    by_kind = {r['kind']: r for r in data}

    transitions = []

    # 1. Update Open Overseer Issue Matches
    open_overseer_matches = {
        'APIGatewayAPIConfig': 7286,
        'AccessContextManagerAccessLevel': 6814,
        'AccessContextManagerAccessLevelCondition': 7343,
        'AccessContextManagerAccessPolicy': 6821,
        'AccessContextManagerServicePerimeter': 8842,
        'AccessContextManagerServicePerimeterResource': 8812,
        'AlloyDBBackup': 8983,
        'BigQueryDatasetAccess': 7359,
        'BinaryAuthorizationAttestor': 7022,
        'BinaryAuthorizationPolicy': 7054,
        'CloudFunctionsFunction': 7061,
        'ComputeAddress': 9730,
        'ComputeBackendBucket': 8880,
        'ComputeBackendService': 8856,
        'ComputeHTTPHealthCheck': 9656,
        'ComputeHealthCheck': 8813,
        'ComputeInstance': 9745,
        'ComputeInstanceGroup': 8798,
        'ComputeInstanceTemplate': 8966,
        'ComputeNetworkEndpointGroup': 7511,
        'ComputeNodeGroup': 8982,
        'ComputeNodeTemplate': 7176,
        'ComputeSSLPolicy': 9724,
        'ComputeSecurityPolicy': 9785,
        'ComputeURLMap': 9740,
        'ContainerCluster': 5186,
        'ContainerNodePool': 9393,
        'DLPJobTrigger': 7047,
        'DataprocAutoscalingPolicy': 9383,
        'DataprocCluster': 9000,
        'Folder': 7374,
        'IAMServiceAccount': 7421,
        'IAPBrand': 9739,
        'LoggingLogBucket': 7058,
        'NetworkConnectivitySpoke': 7020,
        'NetworkServicesGRPCRoute': 7028,
        'PrivateCACertificateTemplate': 9788,
        'Project': 8816,
        'RecaptchaEnterpriseKey': 7026,
        'RedisInstance': 6635,
        'Service': 8819
    }

    for kind, issue_num in open_overseer_matches.items():
        if kind in by_kind:
            r = by_kind[kind]
            if r['state'] != 'Completed':
                if r['state'] == 'Not Started':
                    r['state'] = 'In Progress'
                    transitions.append(f"{kind} -> In Progress (#{issue_num}) [Discovered open overseer issue]")
                r['trackingIssue'] = issue_num

    # 2. Update Completed Matches
    completed_matches = {
        'KMSCryptoKey': 7389,
        'LoggingLogBucket': 9045,
        'PubSubSchema': 9031,
        'KMSKeyRing': 5357
    }

    for kind, issue_num in completed_matches.items():
        if kind in by_kind:
            r = by_kind[kind]
            if r['state'] != 'Completed':
                r['state'] = 'Completed'
                r['trackingIssue'] = issue_num
                if 'steps' in r:
                    for s in r['steps']:
                        r['steps'][s] = True
                transitions.append(f"{kind} -> Completed (#{issue_num}) [Direct controller registered & Overseer closed]")

    # 3. Update Community PR Tracking Issues Created
    new_overseer_creations = {
        'BigQueryReservationCapacityCommitment': (10100, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9577"),
        'BillingBudgetsBudget': (10101, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9176"),
        'ComputeAutoscaler': (10102, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10046"),
        'ComputeBackendBucketSignedURLKey': (10103, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10001"),
        'ComputeBackendServiceSignedURLKey': (10104, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10005"),
        'ComputeDisk': (10105, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10045"),
        'ComputeDiskResourcePolicyAttachment': (10106, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10021"),
        'ComputeExternalVPNGateway': (10107, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10032"),
        'ComputeFirewall': (10108, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10031"),
        'ComputeFirewallPolicy': (10109, "https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10064")
    }

    for kind, (issue_num, pr_url) in new_overseer_creations.items():
        if kind in by_kind:
            r = by_kind[kind]
            r['state'] = 'In Progress'
            r['trackingIssue'] = issue_num
            r['notes'] = f"Community PR: {pr_url}"
            transitions.append(f"{kind} -> In Progress (#{issue_num}) [Linked external PR: {pr_url}]")

    # 4. Record other Community PR URLs in Notes
    all_community_prs = {
        'AlloyDBBackup': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/8987',
        'CloudFunctionsFunction': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/8649',
        'ComputeAddress': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10024',
        'ComputeBackendBucket': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10026',
        'ComputeBackendService': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10049',
        'ContainerNodePool': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9800',
        'DNSRecordSet': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9783',
        'DNSResponsePolicy': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9773',
        'DataflowJob': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9420',
        'FirestoreIndex': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/5570',
        'Folder': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9495',
        'IdentityPlatformConfig': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9421',
        'KMSCryptoKey': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7377',
        'LoggingLogBucket': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9424',
        'MonitoringAlertPolicy': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9692',
        'MonitoringUptimeCheckConfig': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9580',
        'Project': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10079',
        'PubSubSubscription': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9801',
        'SQLUser': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7148',
        'Service': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/10057',
        'StorageBucket': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9784',
        'VertexAIDataset': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/9787',
        'VertexAIEndpoint': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/7994',
        'VertexAIIndex': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/8422',
        'VertexAIIndexEndpoint': 'https://github.com/GoogleCloudPlatform/k8s-config-connector/pull/8422'
    }

    for kind, pr_url in all_community_prs.items():
        if kind in by_kind:
            r = by_kind[kind]
            if not r.get('notes'):
                r['notes'] = f"Community PR: {pr_url}"

    # 5. Save updated data.json
    print(f"Saving updated tracking data to {data_path}...")
    with open(data_path, 'w') as f:
        json.dump(data, f, indent=2)

    # 6. Journal Progress Locally
    journal_path = 'dev/migration-tracker/journal.md'
    now_str = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    completed_count = sum(1 for item in data if item['state'] == 'Completed')
    in_progress_count = sum(1 for item in data if item['state'] == 'In Progress')
    not_started_count = sum(1 for item in data if item['state'] == 'Not Started')
    total_count = len(data)

    journal_entry = f"""
## Run Timestamp: {now_str} (Offline-Updated)

### High-level Status Breakdown
- **Completed**: {completed_count}
- **In Progress**: {in_progress_count}
- **Not Started**: {not_started_count}
- **Total**: {total_count}

### Resource Status Transitions
"""
    if transitions:
        for t in transitions:
            journal_entry += f"- {t}\n"
    else:
        journal_entry += "- No status transitions occurred this run.\n"

    print(f"Appending journal entry to {journal_path}...")
    existing_content = ""
    if os.path.exists(journal_path):
        with open(journal_path, 'r') as f:
            existing_content = f.read()

    with open(journal_path, 'w') as f:
        f.write("# Migration Tracker Journal\n" + journal_entry + "\n" + existing_content)

    # 7. Print the Summary Table to stdout
    active_migrations_md = ""
    for r in data:
        if r['state'] == 'In Progress':
            kind = r['kind']
            issue_num = r.get('trackingIssue', '')
            notes = r.get('notes', '')
            issue_str = f"#{issue_num}" if issue_num else ""
            active_migrations_md += f"| {kind} | {issue_str} | factorybot-robot | {now_str[:10]} | {notes} |\n"

    summary_body = f"""### Migration Progress Tracker Summary

## High-Level Status
| State | Count |
|-------|-------|
| Completed | {completed_count} |
| In Progress | {in_progress_count} |
| Not Started | {not_started_count} |
| Total | {total_count} |

## Active Migrations (In Flight)
| Kind | Overseer Issue | Assignee | Date Started | External Notes |
|------|----------------|----------|--------------|----------------|
{active_migrations_md}"""

    print("\n================ PROGRESS SUMMARY ================")
    print(summary_body)
    print("==================================================\n")
    print("Offline update completed successfully!")

if __name__ == '__main__':
    main()
