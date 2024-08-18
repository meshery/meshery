package clomonitor

# Helper functions to check file existence and patterns
file_exists(filename) {
  input.files[_] == filename
}

file_contains(filename, pattern) {
  re_match(pattern, input.files[filename])
}

# Checks definitions
check("adopters") {
  file_exists("adopters") or
  file_exists("users") or
  file_contains("README.md", "(?i)adopters")
}

check("changelog") {
  file_exists("changelog") or
  file_contains("README.md", "(?i)changelog") or
  file_contains("latest_release", "(?i)changelog")
}

check("code_of_conduct") {
  file_exists("code_of_conduct") or
  file_contains("README.md", "(?i)code of conduct") or
  file_exists(".github/code_of_conduct") or
  file_exists("docs/code_of_conduct")
}

check("contributing") {
  file_exists("contributing") or
  file_contains("README.md", "(?i)contributing") or
  file_exists(".github/contributing") or
  file_exists("docs/contributing")
}

check("governance") {
  file_exists("governance") or
  file_contains("README.md", "(?i)governance") or
  file_exists("docs/governance")
}

check("maintainers") {
  file_exists("maintainers") or
  file_contains("README.md", "(?i)maintainers") or
  file_exists("owners") or
  file_exists("codeowners") or
  file_exists(".github/codeowners") or
  file_exists("docs/maintainers") or
  file_exists("docs/owners") or
  file_exists("docs/codeowners")
}

check("readme") {
  file_exists("README.md") or
  file_exists(".github/README.md") or
  file_exists("docs/README.md")
}

check("roadmap") {
  file_exists("roadmap") or
  file_contains("README.md", "(?i)roadmap")
}

check("summary_table") {
  file_contains("landscape.yaml", "extra:.*clomonitor_name")
}

check("website") {
  file_contains("repo_metadata", "(?i)website")
}

check("license_spdx_id") {
  file_exists("LICENSE") or
  file_exists("COPYING")
}

check("license_approved") {
  license := input.files["LICENSE"]
  re_match("Apache-2.0|BSD-2-Clause|BSD-3-Clause|ISC|MIT|PostgreSQL|Python-2.0|X11|Zlib", license)
}

check("license_scanning") {
  file_contains("README.md", "(https://app.fossa.(?:io|com)/projects/[^\"'\\)]+)") or
  file_contains("README.md", "(https://snyk.io/test/github/[^/]+/[^/\"]+)")
}

check("artifacthub_badge") {
  file_contains("README.md", "(https://artifacthub.io/packages/[^\"'\\)]+)")
}

check("cla") {
  file_contains("latest_merged_pr", "(?i)cncf-cla|(?i)cla/linuxfoundation|(?i)easycla|(?i)license/cla|(?i)cla/google")
}

check("community_meeting") {
  file_contains("README.md", "(?i)(community|developer|development|working group) \\[?(call|event|meeting|session)") or
  file_contains("README.md", "(?i)(weekly|biweekly|monthly) \\[?meeting") or
  file_contains("README.md", "(?i)meeting minutes")
}

check("dco") {
  file_contains("latest_commits", "(?i)dco") or
  file_contains("latest_merged_pr", "(?i)dco")
}

check("github_discussions") {
  file_contains("repo_metadata", "(?i)discussion")
}

check("openssf_badge") {
  file_contains("README.md", "(https://www.bestpractices.dev/projects/\\d+)|(https://bestpractices.coreinfrastructure.org/projects/\\d+)")
}

check("openssf_scorecard_badge") {
  file_contains("README.md", "(https://api.securityscorecards.dev/projects/github.com/[^/]+/[^/]+)/badge")
}

check("recent_release") {
  file_contains("repo_metadata", "(?i)release")
}

check("slack_presence") {
  file_contains("README.md", "(?i)https?://cloud-native.slack.com") or
  file_contains("README.md", "(?i)https?://slack.cncf.io") or
  file_contains("README.md", "(?i)https?://kubernetes.slack.com") or
  file_contains("README.md", "(?i)https?://slack.k8s.io")
}

check("security_policy") {
  file_exists("security") or
  file_contains("README.md", "(?i)security") or
  file_exists(".github/security") or
  file_exists("docs/security")
}

check("sbom") {
  file_contains("latest_release_assets", "(?i)sbom") or
  file_contains("README.md", "(?i)sbom")
}

check("trademark_disclaimer") {
  file_contains("website_content", "(?i)https://(?:w{3}\\.)?linuxfoundation.org/(?:legal/)?trademark-usage") or
  file_contains("website_content", "(?i)The Linux Foundation.* has registered trademarks and uses trademarks")
}

# Main rule to aggregate all checks
default allow = false

allow {
  not failed_checks[_]
}


failed_checks[check_name] {
  check_name := input.checks[_]
  not check(check_name)
}


