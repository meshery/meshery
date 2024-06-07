---
layout: default
title: Understanding How Meshery Generates Models
abstract: Models are generated for capabilities defined in the Meshery Registry using a combination of manual entry and dynamic generation techniques.
permalink: guides/operating/model-generation
type: guides
category: operating
language: en
---

Meshery uses a combination of techniques to generate models for capabilities defined in its Registry. The following are the primary techniques used:

1) **Manual Entry**: The Meshery team manually enters the capabilities of a infrastructure or adapter into the Registry. This is the most common method used to add new capabilities to the Registry.
2) **Model Generation**: Meshery uses a model generation tool to generate models for capabilities. This tool is used to generate models for capabilities that are not manually entered into the Registry.

<h4>Importing Models into the Registry using Meshery CLI</h4>
<p>To register a model using the Meshery CLI, you can use the mesheryctl command to import a model from a specified path:</p>

<pre><code>mesheryctl model import -f &lt;path-to-model&gt; </code></pre>
<h4>Using Meshery UI</h4>
<p>You can also register a model through the Meshery UI:</p>
<ul>
    <li>Navigate to the Settings → Registry page.</li>
    <li>Click the "Import" button.</li>
    <li>Select the model you want to import.</li>
</ul>