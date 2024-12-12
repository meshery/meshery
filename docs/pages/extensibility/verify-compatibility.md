---
layout: default
title: "Verifying Compatibility With Golang Version Update"
permalink: extensibility/verify-compatibility
type: Extensibility
#redirect_from: extensibility/providers
abstract: "When Meshery is updated to use a new version of Golang, extension providers need to ensure their integrations remain compatible by updating their extension to align with Meshery."
language: en
list: include
---

When Meshery is updated to use a newer version of Golang, extension providers need to ensure their integrations remain compatible with Meshery. Follow these steps to verify compatibility of your integrations caused by Golang version changes:

{% include alert.html type="info" title="Note" content="The current Golang Version is mentioned in the <a href='https://docs.meshery.io/project/contributing#meshery-contribution-flow'>Contributing Guide</a>. Update to this version everywhere." %}

### Checkout the Latest Meshery Repository
Clone the latest version of the Meshery repository to ensure you are working with the most recent codebase. If you already have the repository, use git pull to fetch the latest changes.

### Update go.mod to Use the New Golang Version

Open the `go.mod` file located in the root of the Meshery repository. Update the go directive to the new Golang version.
<pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">  go 1.xx</div></div>
</pre>

### Update the Golang Version in the Makefile

Navigate to `/install/Makefile.core.mk` and locate the `$GOVERSION` variable. Update it to reflect the new Golang version.
<pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">  GOVERSION = 1.xx</div></div>
</pre>

### Build and Run the Updated Meshery Server

Compile the updated Meshery server using `make server` to ensure it works with the new Golang version. Once the server starts, check that it runs without errors.

### Update Your Golang-Based Plugin

Update your extension's `go.mod` file to use the same Golang version as the updated Meshery server. Then build the plugin.
<pre class="codeblock-pre">
  <div class="codeblock"><div class="clipboardjs">  go build -o my-plugin</div></div>
</pre>

### Run Meshery Server and Extension
Start your upgraded Meshery server and your updated extension together to verify compatibility. Ensure that both components run without any issues. Monitor the Meshery server logs for any errors or issues related to the extension.

### Validate Your Plugin’s Functionality
Finally, test your extension to ensure it operates as expected. Validate that all key functionalities are working, including:
- The interaction between Meshery and your extension.
- Any custom features specific to your extension.
- UI elements, if applicable.