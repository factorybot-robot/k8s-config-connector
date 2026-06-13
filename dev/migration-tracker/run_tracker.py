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
import re
import subprocess
import datetime

COORDINATOR_ISSUE_NUMBER = 10099

def run_command(cmd):
    """Helper to run shell commands safely and return stdout."""
    res = subprocess.run(cmd, capture_output=True, text=True)
    if res.returncode != 0:
        print(f"Error running command: {' '.join(cmd)}")
        print(f"stderr: {res.stderr}")
    return res.stdout, res.returncode

def check_direct_registered(kind):
    """Checks if the direct controller is registered in static_config.go."""
    static_config_path = 'pkg/controller/resourceconfig/static_config.go'
    if not os.path.exists(static_config_path):
        return False
    with open(static_config_path, 'r') as f:
        for line in f:
            if f'Kind: "{kind}"' in line or f'Kind:\"{kind}\"' in line:
                if 'ReconcilerTypeDirect' in line:
                    return True
    return False

def main():
    data_path = 'dev/migration-tracker/data.json'
    if not os.path.exists(data_path):
        print(f"Error: data.json not found at {data_path}")
        return

    # 1. Load data.json
    with open(data_path, 'r') as f:
        data = json.load(f)

    # Build maps for easy lookup
    by_kind = {r['kind']: r for r in data}
    all_kinds = list(by_kind.keys())

    # Get open issues
    print("Retrieving open issues from GitHub...")
    issues_stdout, _ = run_command([
        'gh', 'issue', 'list', '--state', 'open', '--limit', '1000',
        '--json', 'number,title,labels,assignees,createdAt,url'
    ])
    open_issues = json.loads(issues_stdout or '[]')

    # Get closed issues
    print("Retrieving closed issues from GitHub...")
    closed_stdout, _ = run_command([
        'gh', 'issue', 'list', '--state', 'closed', '--limit', '1000',
        '--json', 'number,title,labels,assignees,createdAt,url'
    ])
    closed_issues = json.loads(closed_stdout or '[]')

    # Get open PRs
    print("Retrieving open pull requests from GitHub...")
    prs_stdout, _ = run_command([
        'gh', 'pr', 'list', '--state', 'open', '--limit', '1000',
        '--json', 'number,title,url,author,labels,createdAt'
    ])
    open_prs = json.loads(prs_stdout or '[]')

    # Combine issues
    all_issues = open_issues + closed_issues

    # Keep track of transitions for journaling
    transitions = []
    new_external_tracking = []

    # Build active overseer issue maps
    open_overseer_issues = [iss for iss in open_issues if 'overseer' in [l['name'] for l in iss.get('labels', [])]]
    closed_overseer_issues = [iss for iss in closed_issues if 'overseer' in [l['name'] for l in iss.get('labels', [])]]

    # Step 2: Audit open (in-flight) migrations
    # We find open overseer issues matching unmigrated kinds using exact word boundaries
    for kind in all_kinds:
        r = by_kind[kind]
        if r['state'] == 'Completed':
            continue

        # Check if there is an active open overseer issue for this Kind
        pattern = re.compile(r'\b' + re.escape(kind) + r'\b', re.IGNORECASE)
        matching_open_issue = None
        for issue in open_overseer_issues:
            if pattern.search(issue['title']) and issue['number'] != COORDINATOR_ISSUE_NUMBER:
                matching_open_issue = issue
                break

        if matching_open_issue:
            issue_num = matching_open_issue['number']
            if r['state'] == 'Not Started':
                r['state'] = 'In Progress'
                r['trackingIssue'] = issue_num
                transitions.append(f"{kind} -> In Progress (#{issue_num}) [Discovered open overseer issue]")
                print(f"Discovered active open overseer issue #{issue_num} for {kind}. Updated to In Progress.")
            elif 'trackingIssue' not in r or r['trackingIssue'] is None:
                r['trackingIssue'] = issue_num
                print(f"Linked open overseer issue #{issue_num} to in-progress {kind}.")

    # Step 2.4: Audit closed migrations
    # Check if any in-progress migrations corresponding to closed overseer issues can be marked Completed
    for kind in all_kinds:
        r = by_kind[kind]
        if r['state'] == 'Completed':
            continue

        # Look for a closed overseer issue matching this kind
        pattern = re.compile(r'\b' + re.escape(kind) + r'\b', re.IGNORECASE)
        matching_closed_issue = None
        for issue in closed_overseer_issues:
            if pattern.search(issue['title']) and issue['number'] != COORDINATOR_ISSUE_NUMBER:
                matching_closed_issue = issue
                break

        if matching_closed_issue:
            # Check if direct controller is indeed registered in static_config.go
            if check_direct_registered(kind):
                r['state'] = 'Completed'
                r['trackingIssue'] = matching_closed_issue['number']
                # Reset/populate relevant step completion flags to true
                if 'steps' in r:
                    for step in r['steps']:
                        r['steps'][step] = True
                transitions.append(f"{kind} -> Completed (#{matching_closed_issue['number']}) [Direct controller registered & Overseer closed]")
                print(f"Closed overseer issue #{matching_closed_issue['number']} found and Direct controller registered for {kind}. Marked Completed.")

    # Step 3: Scan GitHub for Community and External Work
    # Find kinds with open PRs that aren't ours, matching direct controller keywords
    for kind in all_kinds:
        r = by_kind[kind]
        if r['state'] == 'Completed':
            continue

        pattern = re.compile(r'\b' + re.escape(kind) + r'\b', re.IGNORECASE)
        matching_prs = []
        for pr in open_prs:
            author = pr.get('author', {}).get('login', '')
            if author == 'factorybot-robot':
                continue
            title = pr['title']
            if pattern.search(title) and ('direct controller' in title.lower() or 'implement direct' in title.lower() or 'migrate' in title.lower()):
                matching_prs.append(pr)

        if matching_prs:
            # We found active external PR(s)! Record the first one's URL
            external_url = matching_prs[0]['url']
            r['notes'] = f"Community PR: {external_url}"
            new_external_tracking.append(f"Kind {kind} community contribution: {external_url}")
            print(f"Found community PR for {kind}: {external_url}")

            # Check if this resource has an active open overseer tracking issue
            has_open_overseer = False
            for issue in open_issues:
                labels = [l['name'] for l in issue.get('labels', [])]
                if 'overseer' in labels and pattern.search(issue['title']) and issue['number'] != COORDINATOR_ISSUE_NUMBER:
                    has_open_overseer = True
                    break

            if not has_open_overseer:
                # Create a new overseer tracking issue to link and monitor the work
                print(f"Creating overseer tracking issue for {kind} linking to community PR {external_url}...")
                title = f"Migrating {kind} to Direct controller overseer"
                body = f"Active in-flight work detected: {external_url}\n\nUsing the skill workflow, please complete creation of direct controller for {kind}\n\nhttps://raw.githubusercontent.com/gke-labs/gemini-for-kubernetes-development/refs/heads/main/.agents/workflows/kcc-example.txt"
                
                # Modifying the Github repository's issues
                issue_create_cmd = [
                    'gh', 'issue', 'create',
                    '--title', title,
                    '--body', body,
                    '--label', 'overseer',
                    '--assignee', 'factorybot-robot'
                ]
                issue_stdout, code = run_command(issue_create_cmd)
                
                if code == 0 and issue_stdout:
                    # Extract new issue number from URL
                    issue_url = issue_stdout.strip()
                    m = re.search(r'/issues/(\d+)', issue_url)
                    if m:
                        new_issue_num = int(m.group(1))
                        r['state'] = 'In Progress'
                        r['trackingIssue'] = new_issue_num
                        transitions.append(f"{kind} -> In Progress (#{new_issue_num}) [Linked external PR: {external_url}]")
                        print(f"Created overseer issue #{new_issue_num} for {kind}.")
                    else:
                        print(f"Failed to parse issue number from: {issue_url}")
                else:
                    print(f"Failed to create issue for {kind}.")

    # Step 4: Launch Next Migration Tasks
    # Select candidates: state "Not Started", defaultController is "Terraform" or "DCL", all dependencies Completed, no active community PR
    completed_kinds = {item['kind'] for item in data if item['state'] == 'Completed'}
    candidates = []
    for r in data:
        kind = r['kind']
        if r['state'] == 'Not Started' and r.get('defaultController') in ['Terraform', 'DCL']:
            # Must not have community PR recorded in notes
            if r.get('notes') and 'Community PR' in r['notes']:
                continue
            # Check if all dependencies are completed
            deps = r.get('dependencies', [])
            deps_met = all(dep in completed_kinds for dep in deps)
            if deps_met:
                candidates.append(r)

    # Sort candidates by sortOrder ascending
    candidates.sort(key=lambda x: x.get('sortOrder', 9999))

    # Limit new launches to a safe number, e.g., 3 new tasks per run
    MAX_LAUNCH = 3
    launched_count = 0
    for r in candidates:
        if launched_count >= MAX_LAUNCH:
            break

        kind = r['kind']
        print(f"Launching next migration task for {kind} (sortOrder={r.get('sortOrder')})...")
        title = f"Migrating {kind} to Direct controller overseer"
        body = f"Using the skill workflow, please complete creation of direct controller for {kind}\n\nhttps://raw.githubusercontent.com/gke-labs/gemini-for-kubernetes-development/refs/heads/main/.agents/workflows/kcc-example.txt"
        
        issue_create_cmd = [
            'gh', 'issue', 'create',
            '--title', title,
            '--body', body,
            '--label', 'overseer',
            '--assignee', 'factorybot-robot'
        ]
        issue_stdout, code = run_command(issue_create_cmd)
        
        if code == 0 and issue_stdout:
            # Extract new issue number from URL
            issue_url = issue_stdout.strip()
            m = re.search(r'/issues/(\d+)', issue_url)
            if m:
                new_issue_num = int(m.group(1))
                r['state'] = 'In Progress'
                r['trackingIssue'] = new_issue_num
                transitions.append(f"{kind} -> In Progress (#{new_issue_num}) [Launched new migration overseer]")
                print(f"Created overseer issue #{new_issue_num} for {kind}.")
                launched_count += 1
            else:
                print(f"Failed to parse issue number from: {issue_url}")
        else:
            print(f"Failed to create issue for {kind}.")

    # Step 6: Save Local Tracking Data
    print(f"Saving updated tracking data to {data_path}...")
    with open(data_path, 'w') as f:
        json.dump(data, f, indent=2)

    # Step 5: Journal Progress Locally
    journal_path = 'dev/migration-tracker/journal.md'
    now_str = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    # Calculate status counts
    completed_count = sum(1 for item in data if item['state'] == 'Completed')
    in_progress_count = sum(1 for item in data if item['state'] == 'In Progress')
    not_started_count = sum(1 for item in data if item['state'] == 'Not Started')
    total_count = len(data)

    journal_entry = f"""
## Run Timestamp: {now_str}

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

    if new_external_tracking:
        journal_entry += "\n### Newly Tracked External Contributions\n"
        for ext in new_external_tracking:
            journal_entry += f"- {ext}\n"

    print(f"Appending journal entry to {journal_path}...")
    # Read existing or create new journal
    existing_content = ""
    if os.path.exists(journal_path):
        with open(journal_path, 'r') as f:
            existing_content = f.read()

    with open(journal_path, 'w') as f:
        f.write("# Migration Tracker Journal\n" + journal_entry + "\n" + existing_content)

    # Step 7: Post or Update Summary Comment on Coordinator Issue
    # Build list of active in-progress migrations with details
    active_migrations_md = ""
    for r in data:
        if r['state'] == 'In Progress':
            kind = r['kind']
            issue_num = r.get('trackingIssue', '')
            notes = r.get('notes', '')
            # Find matching open issue to get date/assignee if possible
            assignee = 'factorybot-robot'
            date_started = 'N/A'
            for iss in open_issues:
                if iss['number'] == issue_num:
                    assignee_list = [a['login'] for a in iss.get('assignees', [])]
                    if assignee_list:
                        assignee = ', '.join(assignee_list)
                    # format date, e.g. "2026-06-12"
                    try:
                        dt = datetime.datetime.strptime(iss['createdAt'], "%Y-%m-%dT%H:%M:%SZ")
                        date_started = dt.strftime("%Y-%m-%d")
                    except Exception:
                        pass
                    break
            issue_str = f"#{issue_num}" if issue_num else ""
            active_migrations_md += f"| {kind} | {issue_str} | {assignee} | {date_started} | {notes} |\n"

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

    # Search for existing comments on the coordinator issue
    print(f"Retrieving comments on coordinator issue #{COORDINATOR_ISSUE_NUMBER}...")
    comments_stdout, _ = run_command([
        'gh', 'issue', 'view', str(COORDINATOR_ISSUE_NUMBER), '--json', 'comments'
    ])
    try:
        comments_data = json.loads(comments_stdout or '{}')
        comments_list = comments_data.get('comments', [])
    except Exception:
        comments_list = []

    existing_comment_id = None
    for comment in comments_list:
        if '### Migration Progress Tracker Summary' in comment.get('body', ''):
            # For GH API / gh cli, the id or databaseId can be used for editing
            # Let's see if we can edit via gh issue comment edit
            existing_comment_id = comment.get('id')
            break

    if existing_comment_id:
        print(f"Found existing summary comment ID: {existing_comment_id}. Editing comment...")
        # Modifying issue comment
        run_command([
            'gh', 'issue', 'comment', 'edit', str(existing_comment_id), '--body', summary_body
        ])
    else:
        print("No existing summary comment found. Posting new comment...")
        # Modifying issue comments
        run_command([
            'gh', 'issue', 'comment', str(COORDINATOR_ISSUE_NUMBER), '--body', summary_body
        ])

    print("Migration tracker execution completed successfully!")

if __name__ == '__main__':
    main()
