# Security Vulnerabilities

> How the Meshery team handles security vulnerabilities.

Source: /pr-preview/pr-21670/project/security-vulnerabilities/

## List of Announced Vulnerabilities

<table>
<tr>
  <th> DATE ANNOUNCED </th>
  <th> CVE ID </th>
  <th> SEVERITY (CVSS v3) </th>
  <th> DESCRIPTION </th>
  <th> AFFECTED COMPONENT </th>
  <th> VULNERABLE VERSION </th>
  <th> PATCHED VERSION </th>
  <th> FIX DETAILS </th>
  <th> LINKS </th>
</tr>


<tr>
  <td> 2024-08-05 </td>
  <td> CVE-2024-35182 </td>
  <td> Medium (5.9) </td>
  <td> A SQL injection vulnerability in Meshery prior to v0.7.22 in the events API (/api/v2/events, GetAllEvents) allows a remote attacker to execute arbitrary SQL via the sort parameter, including stacked queries and arbitrary file writes via ATTACH DATABASE. </td>
  <td> Events API (GetAllEvents) </td>
  <td> &lt; v0.7.22 </td>
  <td> v0.7.22 </td><td> <a href="https://github.com/meshery/meshery/pull/10280">PR #10280</a> </td>
  <td> <a href="https://nvd.nist.gov/vuln/detail/CVE-2024-35182">NVD</a>,<br><a href="https://github.com/advisories/GHSA-h7cm-jvpp-69xf">GHSA-h7cm-jvpp-69xf</a>,<br><a href="https://securitylab.github.com/advisories/GHSL-2024-013_GHSL-2024-014_Meshery/">GitHub Security Lab</a> </td>
</tr>

<tr>
  <td> 2024-08-05 </td>
  <td> CVE-2024-35181 </td>
  <td> Medium (5.9) </td>
  <td> A SQL injection vulnerability in Meshery prior to v0.7.22 in the GetMeshSyncResourcesKinds handler (/api/system/meshsync/resources/kinds) allows a remote attacker to execute arbitrary SQL via the order parameter, including stacked queries and arbitrary file writes via ATTACH DATABASE. </td>
  <td> MeshSync resources API (GetMeshSyncResourcesKinds) </td>
  <td> &lt; v0.7.22 </td>
  <td> v0.7.22 </td><td> <a href="https://github.com/meshery/meshery/pull/10280">PR #10280</a> </td>
  <td> <a href="https://nvd.nist.gov/vuln/detail/CVE-2024-35181">NVD</a>,<br><a href="https://github.com/advisories/GHSA-9f24-jrv4-f8g5">GHSA-9f24-jrv4-f8g5</a>,<br><a href="https://securitylab.github.com/advisories/GHSL-2024-013_GHSL-2024-014_Meshery/">GitHub Security Lab</a> </td>
</tr>

<tr>
  <td> 2024-08-05 </td>
  <td> CVE-2024-29031 </td>
  <td> High (7.5) </td>
  <td> A SQL injection vulnerability in Meshery prior to v0.7.17 allows a remote attacker to obtain sensitive information via the order parameter of the GetMeshSyncResources function. </td>
  <td> MeshSync resources API (GetMeshSyncResources) </td>
  <td> &lt; v0.7.17 </td>
  <td> v0.7.17 </td><td> <a href="https://github.com/meshery/meshery/pull/10207">PR #10207</a> </td>
  <td> <a href="https://nvd.nist.gov/vuln/detail/CVE-2024-29031">NVD</a>,<br><a href="https://github.com/advisories/GHSA-652r-q29p-m25h">GHSA-652r-q29p-m25h</a>,<br><a href="https://securitylab.github.com/advisories/GHSL-2023-249_Meshery/">GitHub Security Lab</a> </td>
</tr>

<tr>
  <td> 2024-07-24 </td>
  <td> CVE-2024-36535 </td>
  <td> Critical (9.8) </td>
  <td> Insecure default permissions in Meshery v0.7.51 allow an attacker to access sensitive data and escalate privileges by obtaining the Meshery ServiceAccount token. With the default Helm installation the interface is exposed on an external IP and permits open self-registration, so anyone who can reach it can sign up, gain broad visibility into cluster activity, deploy pods, and execute arbitrary code unless Meshery is secured or restricted to internal networks. </td>
  <td> Default deployment / RBAC (ServiceAccount token) </td>
  <td> v0.7.51 </td>
  <td> Mitigation only </td><td> Mitigated by configuration &mdash; see <a href="/pr-preview/pr-21670/installation/production/security-hardening">Security Hardening</a> and <a href="/pr-preview/pr-21670/installation/production/authentication-and-identity">Authentication &amp; Identity</a> </td>
  <td> <a href="https://nvd.nist.gov/vuln/detail/CVE-2024-36535">NVD</a>,<br><a href="https://www.cve.org/CVERecord?id=CVE-2024-36535">cve.org</a>,<br><a href="https://gist.github.com/HouqiyuA/2950c3993cdeff23afcbd73ba7a33879">advisory</a> </td>
</tr>

<tr>
  <td> 2023-11-24 </td>
  <td> CVE-2023-46575 </td>
  <td> Critical (9.1) </td>
  <td> A SQL injection vulnerability in Meshery prior to v0.6.179 enables a remote attacker to retrieve sensitive information and execute arbitrary code via the order parameter. </td>
  <td> REST API (order parameter) </td>
  <td> &lt; v0.6.179 </td>
  <td> v0.6.179 </td><td> <a href="https://github.com/meshery/meshery/pull/9372">PR #9372</a> </td>
  <td> <a href="https://nvd.nist.gov/vuln/detail/CVE-2023-46575">NVD</a>,<br><a href="https://github.com/advisories/GHSA-9jjc-grg5-67gj">GHSA-9jjc-grg5-67gj</a> </td>
</tr>

<tr>
  <td> 2021-04-28 </td>
  <td> CVE-2021-31856 </td>
  <td> Critical (9.8) </td>
  <td> A SQL Injection vulnerability in the REST API in Meshery 0.5.2 allows an attacker to execute arbitrary SQL commands via the /experimental/patternfiles endpoint (order parameter in GetMesheryPatterns in models/meshery_pattern_persister.go). </td>
  <td> REST API </td>
  <td> v0.5.2 </td>
  <td> v0.5.3 </td><td> <a href="https://github.com/meshery/meshery/pull/2745">PR #2745</a> </td>
  <td> <a href="https://nvd.nist.gov/vuln/detail/CVE-2021-31856">NVD</a>,<br><a href="https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-31856">mitre</a>,<br><a href="https://github.com/ssst0n3/CVE-2021-31856">details</a> </td>
</tr>

</table>


## Reporting a vulnerability

We are very grateful to the security researchers and users that report
back Meshery security vulnerabilities. We investigate every report thoroughly.

To make a report, send an email to the private
[security@meshery.io](mailto:security@meshery.io)
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
  - The [Meshery blog](https://meshery.io/blog)
  - The [Meshery X feed](https://x.com/mesheryio)
  - The [#announcements](https://mesheryio.slack.com/archives/CSF3PSZT9) channel on community [Slack](https://slack.meshery.io/)

As much as possible this announcement will be actionable, and include any mitigating steps customers can take prior to upgrading to a fixed version.
