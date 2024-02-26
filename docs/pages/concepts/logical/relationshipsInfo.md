---
layout: enhanced
title: Relationships Information
permalink: concepts/logical/relationshipsInfo
abstract: "Information about relationships between different components."
language: en
list: include
redirect_from:
  - concepts/relationshipsinfo.html
---

**Relationships Information**

This page provides a concise overview of the interconnections between various components, along with essential caveats and considerations for users to bear in mind when defining these relationships.
<table class="table table-dark table-active">
    <tr>
        <th>From</th>
        <th>To</th>
        <th>Relationship Type</th>
        <th>Relationship SubType</th>
        <th>Evaluation Policy</th>
        <th>Considerations</th>
        <th>Cavets</th>
    </tr>
    <tr>
        <td>Pod</td>
        <td>Deployment, StatefulSet, DaemonSet, ReplicaSet</td>
        <td>Hierarchical</td>
        <td>Inventory</td>
        <td>heirarchical_inventory_relationship</td>
        <td>While defining a relationship between a pod and a deployment, any configurations added to the pod are automatically incorporated into the deployment manifest. Consequently, users don't need to separately specify configurations for the deployment.</td>
        <td>After the relationship is established. Unfortunately, there's no system to remove the extra pod configuration automatically in this situation. This can result in duplication, where both the pod and deployment have similar configurations. It's important to be aware of this possibility and manage configurations carefully to avoid unexpected issues during deployment.</td>
    </tr>
</table>