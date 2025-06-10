---
layout: default
title: Importing Models
abstract: Import existing Models and CRD-based Infrastructure Configurations into Meshery as Models.
permalink: guides/configuration-management/importing-models
category: configuration
type: guides
language: en
---

Import your existing Models and custom resource definitions (CRDs) into Meshery to leverage powerful service mesh configuration and lifecycle management capabilities. Meshery supports multiple application definition formats and offers flexible import options via both the CLI and UI.

**Note:** A valid [Model](/concepts/logical/models) must contain at least one valid [Component](/concepts/logical/components) or [Relationship](/concepts/logical/relationships) to be imported successfully.

{% tabs %}
{% tab CLI %}

## Import Models Using Meshery CLI

### Step 1: Install Meshery CLI

To begin, ensure the Meshery CLI (`mesheryctl`) is installed on your system. This CLI tool is the primary interface for managing Meshery resources from the command line. You can install it by following the instructions on the [Meshery installation page]({{site.baseurl}}/installation#install-mesheryctl).

### Step 2: Import the Model

You can import a model either from a local file or directly via a URL. The model file must be a Meshery-exported archive containing the appropriate model components.

Use the following command syntax:

```bash
mesheryctl model import -f [file-or-url]
