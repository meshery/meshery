---
layout: default
title: Importing Models
abstract: Importing existing Model and CRD-based infrastructure configurations into Meshery
permalink: guides/configuration-management/importing-models
category: configuration
type: guides
language: en
---

Import your existing Models and CRD-based configurations into Meshery. The platform supports a variety of application definition formats. Models can be imported using either the Meshery CLI or Meshery UI.

> **Note:** A [Model](/concepts/logical/models) can only be imported if it contains at least one valid [Component](/concepts/logical/components) or [Relationship](/concepts/logical/relationships).

---

## 🧰 Import Models using Meshery CLI

### Step 1: Install Meshery CLI

Ensure you have Meshery CLI (`mesheryctl`) installed. Follow the [official installation guide]({{site.baseurl}}/installation#install-mesheryctl).

### Step 2: Import the Model

Meshery supports importing models via local files or URLs. The model must be a Meshery-exported model.

#### 📄 Import from Local File

```bash
mesheryctl model import -f ./istio-base.tar
