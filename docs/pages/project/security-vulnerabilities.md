---
layout: default
title: Security Vulnerabilities
permalink: project/security-vulnerabilities
description: How the Meshery team handles security vulnerabilities.
language: en
type: project
category: none
list: include
---

## List of Announced Vulnerabilities


<table>
<tr>
  <th> DATE ANNOUNCED </th>
  <th> CVE ID </th>
  <th> DESCRIPTION </th>
  <th> AFFECTED COMPONENT </th>
  <th> VULNERABLE VERSION </th>
  <th> PATCHED VERSION </th>
  <th> FIX DETAILS </th>
  <th> LINKS </th>
</tr>
{% assign vulns = site.data.vulnerabilities.announce | sort: "Date-Announced" | reverse %}

{% for vuln in vulns %}

<tr>
  <td> {{vuln.DateAnnounced}} </td>
  <td> {{vuln.CVE}} </td>
  <td> {{vuln.Description}} </td>
  <td> {{vuln.AffectedComponent}} </td>
  <td> {{vuln.VulnerableVersion}} </td>
  <td> {{vuln.PatchedVersion}} </td>
  <td> {{vuln.FixDetails}} </td>
  <td> {{vuln.Links}} </td>
</tr>

{% endfor %}
</table>

## Reporting a vulnerability

We are very grateful to the security researchers and users that report
back Meshery security vulnerabilities. We investigate every report thoroughly.

To make a report, send an email to the private
[security@meshery.dev](mailto:security@meshery.dev)
mailing list with the vulnerability details. For normal product bugs
unrelated to latent security vulnerabilities, please head to
the appropriate repository and submit a [new issue](https://github.com/meshery/meshery/issues/new/choose).

### When to report a security vulnerability?

Send us a report whenever you:

- Think Meshery has a potential security vulnerability.
- Are unsure whether or how a vulnerability affects Meshery.
- Think a vulnerability is present in another project that Meshery
depends on (Docker for example).

### When not to report a security vulnerability?

Don't send a vulnerability report if:

- You need help tuning Meshery components for security.
- You need help applying security related updates.
- Your issue is not security related.

Instead, join the community [Slack](https://slack.meshery.io/) and ask questions.

### Evaluation

The Meshery team acknowledges and analyzes each vulnerability report within 10 working days.

Any vulnerability information you share with the Meshery team stays
within the Meshery project. We don't disseminate the information to other
projects. We only share the information as needed to fix the issue.

We keep the reporter updated as the status of the security issue is addressed.

### Fixing the issue

Once a security vulnerability has been fully characterized, a fix is developed by the Meshery team.
The development and testing for the fix happens in a private GitHub repository in order to prevent
premature disclosure of the vulnerability.

### Early disclosures

The Meshery project maintains a mailing list for private early disclosure of security vulnerabilities. 
The list is used to provide actionable information to close Meshery partners. The list is not intended 
for individuals to find out about security issues.

### Public disclosures

On the day chosen for public disclosure, a sequence of activities takes place as quickly as possible:

- Changes are merged from the private GitHub repository holding the fix into the appropriate set of public
branches.

- Meshery team ensures all necessary binaries are promptly built and published.

- Once the binaries are available, an announcement is sent out on the following channels:
  - The [Meshery blog](https://meshery.io/blog/)
  - The [Meshery Twitter feed](https://twitter.com/mesheryio)
  - The [#announcements](https://layer5io.slack.com/archives/CSF3PSZT9) channel on community [Slack](https://slack.meshery.io/)

As much as possible this announcement will be actionable, and include any mitigating steps customers can take prior to upgrading to a fixed version.


