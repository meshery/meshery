# Security Policy
We are very grateful to the security researchers and users that report
back Meshery security vulnerabilities. We investigate every report thoroughly.

## Reporting a vulnerability
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

## Evaluation

The Meshery team acknowledges and analyzes each vulnerability report within 10 working days.

Any vulnerability information you share with the Meshery team stays
within the Meshery project. We don't disseminate the information to other
projects. We only share the information as needed to fix the issue.

We keep the reporter updated as the status of the security issue is addressed.

## Fixing the issue

Once a security vulnerability has been fully characterized, a fix is developed by the Meshery team.
The development and testing for the fix happens in a private GitHub repository in order to prevent
premature disclosure of the vulnerability.

## Early disclosure

The Meshery project maintains a mailing list for private early disclosure of security vulnerabilities. 
The list is used to provide actionable information to close Meshery partners. The list is not intended 
for individuals to find out about security issues.

## Public disclosure

On the day chosen for public disclosure, a sequence of activities takes place as quickly as possible:

- Changes are merged from the private GitHub repository holding the fix into the appropriate set of public
branches.
- Meshery team ensures all necessary binaries are promptly built and published.
- Once the binaries are available, an announcement is sent out on the following channels:
  - The [Meshery blog](https://meshery.io/blog/)
  - The [Meshery Twitter feed](https://twitter.com/mesheryio)
  - The #announcements channel on Slack

As much as possible this announcement will be actionable, and include any mitigating steps customers can take prior to
upgrading to a fixed version. 

