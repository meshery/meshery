#!/usr/bin/env python3
"""
Hacktoberfest Badge Automation Script for Meshery

This script automates the process of awarding Hacktoberfest badges to contributors
by monitoring GitHub issues, detecting merged PRs, extracting contributor emails,
and posting award commands to Slack.

Usage:
    python hacktoberfest-badge-automation.py --issues 16100,16101,16102
    python hacktoberfest-badge-automation.py --issues-file issues.txt
    python hacktoberfest-badge-automation.py --label hacktoberfest

Author: Meshery Community
Repository: https://github.com/meshery/meshery
"""

import os
import re
import sys
import json
import logging
import argparse
from typing import List, Set, Dict, Optional
from datetime import datetime
from pathlib import Path

# Check for required dependencies
try:
    import requests
except ImportError:
    print("ERROR: 'requests' module not found.")
    print("Install it with: pip install requests")
    sys.exit(1)

try:
    from dotenv import load_dotenv
except ImportError:
    print("WARNING: 'python-dotenv' not found. Environment variables must be set manually.")
    print("Install it with: pip install python-dotenv")
    load_dotenv = lambda: None


class Config:
    """Configuration management for the automation script"""
    
    def __init__(self):
        load_dotenv()
        
        # GitHub Configuration
        self.github_token = os.getenv('GITHUB_TOKEN', '')
        self.github_repo = os.getenv('GITHUB_REPO', 'meshery/meshery')
        self.github_api_base = 'https://api.github.com'
        
        # Slack Configuration
        self.slack_token = os.getenv('SLACK_BOT_TOKEN', '')
        self.slack_channel = os.getenv('SLACK_CHANNEL', 'event')
        self.badge_name = os.getenv('BADGE_NAME', 'hacktoberfest25')
        
        # Script Configuration
        self.dry_run = os.getenv('DRY_RUN', 'false').lower() == 'true'
        self.log_level = os.getenv('LOG_LEVEL', 'INFO').upper()
        self.awarded_emails_file = os.getenv('AWARDED_EMAILS_FILE', 'awarded_emails.json')
        
    def validate(self, dry_run: bool = False) -> List[str]:
        """Validate configuration and return list of missing items"""
        if dry_run:
            return []  # Don't require tokens in dry-run mode
        
        missing = []
        if not self.github_token:
            missing.append('GITHUB_TOKEN')
        if not self.slack_token:
            missing.append('SLACK_BOT_TOKEN')
        return missing


def setup_logging(log_level: str) -> logging.Logger:
    """Setup logging configuration"""
    log_format = '%(asctime)s - %(levelname)s - %(message)s'
    date_format = '%Y-%m-%d %H:%M:%S'
    
    # Create logger
    logger = logging.getLogger(__name__)
    logger.setLevel(getattr(logging, log_level))
    
    # Console handler
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setFormatter(logging.Formatter(log_format, date_format))
    logger.addHandler(console_handler)
    
    # File handler
    log_file = f'hacktoberfest_automation_{datetime.now().strftime("%Y%m%d_%H%M%S")}.log'
    file_handler = logging.FileHandler(log_file)
    file_handler.setFormatter(logging.Formatter(log_format, date_format))
    logger.addHandler(file_handler)
    
    return logger


class GitHubAPI:
    """GitHub API interaction handler"""
    
    def __init__(self, token: str, repo: str, base_url: str, logger: logging.Logger):
        self.token = token
        self.repo = repo
        self.base_url = base_url
        self.logger = logger
        self.session = requests.Session()
        self.timeout = 30
        
        if token:
            self.session.headers.update({
                'Authorization': f'token {token}',
                'Accept': 'application/vnd.github.v3+json'
            })
    
    def get_issue(self, issue_number: int) -> Optional[Dict]:
        """Fetch issue details"""
        url = f'{self.base_url}/repos/{self.repo}/issues/{issue_number}'
        try:
            response = self.session.get(url, timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Failed to fetch issue #{issue_number}: {e}")
            return None
    
    def get_issues_by_label(self, label: str) -> List[Dict]:
        """Fetch all issues with a specific label"""
        url = f'{self.base_url}/repos/{self.repo}/issues'
        params = {'labels': label, 'state': 'all', 'per_page': 100}
        issues = []
        
        try:
            while url:
                response = self.session.get(url, params=params, timeout=self.timeout)
                response.raise_for_status()
                issues.extend(response.json())
                
                # Handle pagination
                url = response.links.get('next', {}).get('url')
                params = {}
                
            self.logger.info(f"Found {len(issues)} issues with label '{label}'")
            return issues
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Failed to fetch issues by label '{label}': {e}")
            return []
    
    def get_issue_timeline(self, issue_number: int) -> List[Dict]:
        """Fetch issue timeline to find linked PRs"""
        url = f'{self.base_url}/repos/{self.repo}/issues/{issue_number}/timeline'
        try:
            response = self.session.get(url, headers={'Accept': 'application/vnd.github.mockingbird-preview+json'}, timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Failed to fetch timeline for issue #{issue_number}: {e}")
            return []
    
    def get_pull_request(self, pr_number: int) -> Optional[Dict]:
        """Fetch PR details"""
        url = f'{self.base_url}/repos/{self.repo}/pulls/{pr_number}'
        try:
            response = self.session.get(url, timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Failed to fetch PR #{pr_number}: {e}")
            return None
    
    def get_pr_commits(self, pr_number: int) -> List[Dict]:
        """Fetch all commits from a PR"""
        url = f'{self.base_url}/repos/{self.repo}/pulls/{pr_number}/commits'
        try:
            response = self.session.get(url, timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Failed to fetch commits for PR #{pr_number}: {e}")
            return []
    
    def extract_pr_numbers_from_issue(self, issue_number: int) -> List[int]:
        """Extract PR numbers linked to an issue"""
        pr_numbers = []
        
        # Get issue timeline
        timeline = self.get_issue_timeline(issue_number)
        
        for event in timeline:
            if event.get('event') == 'cross-referenced':
                source = event.get('source', {})
                if source.get('type') == 'issue' and source.get('issue', {}).get('pull_request'):
                    pr_url = source['issue']['pull_request']['url']
                    pr_number = int(pr_url.split('/')[-1])
                    pr_numbers.append(pr_number)
            
            elif event.get('event') == 'referenced':
                commit_id = event.get('commit_id')
                if commit_id:
                    pr_num = self._find_pr_by_commit(commit_id)
                    if pr_num:
                        pr_numbers.append(pr_num)
        
        # Check issue body for PR references
        issue = self.get_issue(issue_number)
        if issue:
            body = issue.get('body', '') or ''
            pr_refs = re.findall(r'#(\d+)', body)
            pr_numbers.extend([int(num) for num in pr_refs])
        
        return list(set(pr_numbers))
    
    def _find_pr_by_commit(self, commit_sha: str) -> Optional[int]:
        """Find PR number associated with a commit"""
        url = f'{self.base_url}/repos/{self.repo}/commits/{commit_sha}/pulls'
        try:
            response = self.session.get(url, headers={'Accept': 'application/vnd.github.groot-preview+json'}, timeout=self.timeout)
            response.raise_for_status()
            prs = response.json()
            if prs:
                return prs[0]['number']
        except requests.exceptions.RequestException as e:
            self.logger.debug(f"Could not find PR for commit {commit_sha}: {e}")
        return None


class EmailExtractor:
    """Extract emails from commit messages"""
    
    SIGNOFF_PATTERN = re.compile(r'Signed-off-by:\s*(.+?)\s*<(.+?)>', re.IGNORECASE | re.MULTILINE)
    
    @staticmethod
    def extract_emails_from_commits(commits: List[Dict], logger: logging.Logger) -> Set[str]:
        """Extract all unique emails from commit Signed-off-by lines"""
        emails = set()
        
        for commit in commits:
            commit_data = commit.get('commit', {})
            message = commit_data.get('message', '')
            
            matches = EmailExtractor.SIGNOFF_PATTERN.findall(message)
            
            for name, email in matches:
                email = email.strip().lower()
                if EmailExtractor.is_valid_email(email):
                    emails.add(email)
                    logger.debug(f"Extracted email: {email} from commit {commit.get('sha', 'unknown')[:7]}")
                else:
                    logger.warning(f"Invalid email format: {email}")
        
        return emails
    
    @staticmethod
    def is_valid_email(email: str) -> bool:
        """Basic email validation"""
        pattern = re.compile(r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$')
        return bool(pattern.match(email))


class SlackAPI:
    """Slack API interaction handler"""
    
    def __init__(self, token: str, channel: str, badge_name: str, logger: logging.Logger, dry_run: bool = False):
        self.token = token
        self.channel = channel
        self.badge_name = badge_name
        self.logger = logger
        self.dry_run = dry_run
        self.api_url = 'https://slack.com/api/chat.postMessage'
    
    def award_badge(self, email: str) -> bool:
        """Post badge award command to Slack channel"""
        command = f"/award-badge {email} {self.badge_name}"
        
        if self.dry_run:
            self.logger.info(f"[DRY RUN] Would post to #{self.channel}: {command}")
            return True
        
        payload = {
            'channel': self.channel,
            'text': command
        }
        
        headers = {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
        
        try:
            response = requests.post(self.api_url, json=payload, headers=headers, timeout=30)
            response.raise_for_status()
            result = response.json()
            
            if result.get('ok'):
                self.logger.info(f"Successfully posted award command for {email}")
                return True
            else:
                self.logger.error(f"Slack API error: {result.get('error', 'Unknown error')}")
                return False
                
        except requests.exceptions.RequestException as e:
            self.logger.error(f"Failed to post to Slack: {e}")
            return False


class AwardTracker:
    """Track already awarded badges to prevent duplicates"""
    
    def __init__(self, file_path: str, logger: logging.Logger):
        self.file_path = Path(file_path)
        self.logger = logger
        self.awarded_emails: Set[str] = self._load()
    
    def _load(self) -> Set[str]:
        """Load previously awarded emails from file"""
        if not self.file_path.exists():
            return set()
        
        try:
            with open(self.file_path, 'r') as f:
                data = json.load(f)
                return set(data.get('awarded_emails', []))
        except (json.JSONDecodeError, IOError) as e:
            self.logger.warning(f"Could not load awarded emails file: {e}")
            return set()
    
    def _save(self):
        """Save awarded emails to file"""
        try:
            with open(self.file_path, 'w') as f:
                json.dump({
                    'awarded_emails': list(self.awarded_emails),
                    'last_updated': datetime.now().isoformat()
                }, f, indent=2)
        except IOError as e:
            self.logger.error(f"Failed to save awarded emails file: {e}")
    
    def is_awarded(self, email: str) -> bool:
        """Check if email has already been awarded"""
        return email.lower() in self.awarded_emails
    
    def mark_awarded(self, email: str):
        """Mark email as awarded"""
        self.awarded_emails.add(email.lower())
        self._save()


class HacktoberfestAutomation:
    """Main automation orchestrator"""
    
    def __init__(self, config: Config, logger: logging.Logger):
        self.config = config
        self.logger = logger
        self.github = GitHubAPI(config.github_token, config.github_repo, config.github_api_base, logger)
        self.slack = SlackAPI(config.slack_token, config.slack_channel, config.badge_name, logger, config.dry_run)
        self.tracker = AwardTracker(config.awarded_emails_file, logger)
        self.stats = {
            'issues_processed': 0,
            'prs_found': 0,
            'merged_prs': 0,
            'emails_extracted': 0,
            'badges_awarded': 0,
            'duplicates_skipped': 0,
            'errors': 0
        }
    
    def process_issue(self, issue_number: int) -> Set[str]:
        """Process a single issue and award badges"""
        self.logger.info(f"Processing issue #{issue_number}...")
        self.stats['issues_processed'] += 1
        awarded_emails = set()
        
        # Find linked PRs
        pr_numbers = self.github.extract_pr_numbers_from_issue(issue_number)
        
        if not pr_numbers:
            self.logger.warning(f"No PRs found for issue #{issue_number}")
            return awarded_emails
        
        self.logger.info(f"Found {len(pr_numbers)} PR(s) for issue #{issue_number}: {pr_numbers}")
        self.stats['prs_found'] += len(pr_numbers)
        
        # Process each PR
        for pr_number in pr_numbers:
            emails = self._process_pr(pr_number)
            awarded_emails.update(emails)
        
        return awarded_emails
    
    def _process_pr(self, pr_number: int) -> Set[str]:
        """Process a single PR"""
        awarded_emails = set()
        
        # Check if PR is merged
        pr = self.github.get_pull_request(pr_number)
        if not pr:
            self.stats['errors'] += 1
            return awarded_emails
        
        if not pr.get('merged', False):
            self.logger.info(f"PR #{pr_number} is not merged yet, skipping")
            return awarded_emails
        
        self.logger.info(f"PR #{pr_number} is merged [OK]")
        self.stats['merged_prs'] += 1
        
        # Get commits and extract emails
        commits = self.github.get_pr_commits(pr_number)
        if not commits:
            self.logger.warning(f"No commits found for PR #{pr_number}")
            return awarded_emails
        
        emails = EmailExtractor.extract_emails_from_commits(commits, self.logger)
        self.logger.info(f"Extracted {len(emails)} email(s) from PR #{pr_number}: {emails}")
        self.stats['emails_extracted'] += len(emails)
        
        # Award badges
        for email in emails:
            if self.tracker.is_awarded(email):
                self.logger.info(f"Email {email} already awarded, skipping")
                self.stats['duplicates_skipped'] += 1
                continue
            
            if self.slack.award_badge(email):
                self.tracker.mark_awarded(email)
                awarded_emails.add(email)
                self.stats['badges_awarded'] += 1
            else:
                self.stats['errors'] += 1
        
        return awarded_emails
    
    def run(self, issue_numbers: List[int]):
        """Run automation for list of issues"""
        self.logger.info(f"Starting Hacktoberfest badge automation for {len(issue_numbers)} issue(s)")
        self.logger.info(f"Configuration: Repo={self.config.github_repo}, Badge={self.config.badge_name}, DryRun={self.config.dry_run}")
        
        all_awarded_emails = set()
        
        for issue_number in issue_numbers:
            try:
                awarded = self.process_issue(issue_number)
                all_awarded_emails.update(awarded)
            except Exception as e:
                self.logger.error(f"Error processing issue #{issue_number}: {e}", exc_info=True)
                self.stats['errors'] += 1
        
        self._print_summary(all_awarded_emails)
    
    def run_by_label(self, label: str):
        """Run automation for all issues with a specific label"""
        self.logger.info(f"Fetching issues with label '{label}'...")
        issues = self.github.get_issues_by_label(label)
        
        if not issues:
            self.logger.warning(f"No issues found with label '{label}'")
            return
        
        issue_numbers = [issue['number'] for issue in issues]
        self.run(issue_numbers)
    
    def _print_summary(self, awarded_emails: Set[str]):
        """Print execution summary"""
        separator = "=" * 80
        self.logger.info(separator)
        self.logger.info("EXECUTION SUMMARY")
        self.logger.info(separator)
        self.logger.info(f"Issues processed: {self.stats['issues_processed']}")
        self.logger.info(f"PRs found: {self.stats['prs_found']}")
        self.logger.info(f"Merged PRs: {self.stats['merged_prs']}")
        self.logger.info(f"Emails extracted: {self.stats['emails_extracted']}")
        self.logger.info(f"Badges awarded: {self.stats['badges_awarded']}")
        self.logger.info(f"Duplicates skipped: {self.stats['duplicates_skipped']}")
        self.logger.info(f"Errors: {self.stats['errors']}")
        
        if awarded_emails:
            self.logger.info("\nEmails awarded in this run:")
            for email in sorted(awarded_emails):
                self.logger.info(f"  - {email}")
        
        self.logger.info(separator)


def parse_arguments() -> argparse.Namespace:
    """Parse command line arguments"""
    parser = argparse.ArgumentParser(
        description='Automate Hacktoberfest badge awards for merged contributions',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # Process specific issues
  python hacktoberfest-badge-automation.py --issues 16100,16101,16102
  
  # Process issues from file
  python hacktoberfest-badge-automation.py --issues-file issues.txt
  
  # Process all issues with a label
  python hacktoberfest-badge-automation.py --label hacktoberfest
  
  # Dry run (no actual Slack posts)
  python hacktoberfest-badge-automation.py --issues 16100 --dry-run
        """
    )
    
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument('--issues', help='Comma-separated list of issue numbers (e.g., 16100,16101)')
    group.add_argument('--issues-file', help='File containing issue numbers (one per line)')
    group.add_argument('--label', help='GitHub label to filter issues (e.g., hacktoberfest)')
    
    parser.add_argument('--dry-run', action='store_true', help='Run without posting to Slack')
    parser.add_argument('--log-level', choices=['DEBUG', 'INFO', 'WARNING', 'ERROR'], 
                       default='INFO', help='Set logging level')
    
    return parser.parse_args()


def load_issues_from_file(file_path: str, logger: logging.Logger) -> List[int]:
    """Load issue numbers from file"""
    try:
        with open(file_path, 'r') as f:
            lines = f.readlines()
        
        issues = []
        for line in lines:
            line = line.strip()
            if line and not line.startswith('#'):
                try:
                    issues.append(int(line))
                except ValueError:
                    logger.warning(f"Skipping invalid issue number: {line}")
        
        return issues
    except IOError as e:
        logger.error(f"Failed to read issues file: {e}")
        sys.exit(1)


def main():
    """Main entry point"""
    args = parse_arguments()
    
    # Setup configuration and logging
    config = Config()
    if args.dry_run:
        config.dry_run = True
    if args.log_level:
        config.log_level = args.log_level
    
    logger = setup_logging(config.log_level)
    
    logger.info("=" * 80)
    logger.info("Hacktoberfest Badge Automation Script")
    logger.info("=" * 80)
    
    # Validate configuration
    missing = config.validate(dry_run=config.dry_run)
    if missing:
        logger.error(f"Missing required environment variables: {', '.join(missing)}")
        logger.error("Please set them in .env file or environment")
        logger.info("Run with --dry-run to test without tokens")
        sys.exit(1)
    
    if config.dry_run:
        logger.warning("WARNING: Running in DRY RUN mode - no badges will be awarded")
    
    # Parse issue numbers
    issue_numbers = []
    if args.issues:
        try:
            issue_numbers = [int(num.strip()) for num in args.issues.split(',')]
        except ValueError:
            logger.error("Invalid issue numbers provided")
            sys.exit(1)
    elif args.issues_file:
        issue_numbers = load_issues_from_file(args.issues_file, logger)
    
    # Run automation
    automation = HacktoberfestAutomation(config, logger)
    
    if args.label:
        automation.run_by_label(args.label)
    else:
        if not issue_numbers:
            logger.error("No issue numbers provided")
            sys.exit(1)
        automation.run(issue_numbers)
    
    logger.info("[OK] Script execution completed")


if __name__ == '__main__':
    main()