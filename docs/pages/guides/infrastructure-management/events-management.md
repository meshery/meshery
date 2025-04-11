---
layout: default
title: Managing Events with Notification Center
abstract: Meshery tracks operations performed on your infrastructure and workloads, and provides notification of environment issues, application conflicts with infrastructure configuration, policy violations, and so on.
permalink: guides/infrastructure-management/notification-management
redirect_from: guides/events-management
type: guides
category: infrastructure
language: en
---

Meshery continuously tracks operations performed on your infrastructure and workloads. While running, Meshery Server generates a variety of events reflecting activities inside your Kubernetes and Cloud-Native clusters, as well as within Meshery itself.

### What is the Notification Center?
The Notification Center is a dedicated panel in Meshery’s UI that helps you monitor, understand, and respond to events across your system. It acts as a central place where you can see important updates related to your infrastructure, workloads, and Meshery’s internal operations.

![Find Notification Center]({{ site.baseurl }}/assets/img/notification-center/find-notification-center.gif)

### Event Descriptions

Each notification in Meshery includes a clear summary of what occurred in your system. Notifications vary in format depending on the event type, but typically include:

- Action performed (e.g., saved a design, ran a dry run)
- Affected components (e.g., deployments, MeshSync, events)
- Validation results or errors (e.g., invalid values, missing fields)
- Relationship updates (e.g., how components are linked)
- Links to the related design or further details

You can mark notifications as read or unread to stay organized and focused. Meshery highlights critical, actionable events in red, helping you quickly spot and respond to urgent issues.

<a href="{{ site.baseurl }}/assets/img/notification-center/description.png">
  <img src="{{ site.baseurl }}/assets/img/notification-center/description.png" alt="Event Descriptions" style="width: 600px; height: auto; margin-top: 10px;" />
</a>

🔗 For more technical details, see the [Contributor Reference](#).

### Notification Timestamps
Each notification includes a timestamp showing when the event happened. The time is displayed based on your local device’s time zone, so it reflects your current time.

<a href="{{ site.baseurl }}/assets/img/notification-center/timestamps.png">
  <img src="{{ site.baseurl }}/assets/img/notification-center/timestamps.png" alt="Notification Timestamps" style="width: 600px; height: auto; margin-top: 10px;" />
</a>

### Data Sharing
Need to collaborate?
You can share notifications with teammates or stakeholders in just a few clicks — making it easier to communicate and resolve issues.

(???)

### Filtering and Searching

The Notification Center provides a powerful way to filter and search through events. You can narrow down results using filters such as **severity**, **status**, **action**, **category**, and **author**.

> Note: Some filter options such as `action`, `category`, and `author` are dynamically generated based on the notifications your Meshery instance has received. These values are retrieved from the `/api/system/events/types` endpoint.

(A GIF demonstrating the filter interaction will be added here after the UI update.)

#### Severity  

Filter notifications based on the level of severity, indicated by icon color and symbol. These levels are defined in the `SEVERITY` constant and styled using `SEVERITY_STYLE`.

| Level        | Code Value       | Icon       | Color (Light Mode)   | Description                     |
|--------------|------------------|------------|----------------------|---------------------------------|
| Info         | `informational`  | ℹ️ InfoIcon | Blue                 | General updates or logs         |
| Warning      | `warning`        | ⚠️ AlertIcon | Yellow               | Potential issues                |
| Error        | `error`          | ❌ ErrorIcon | Red                  | Failures or critical problems   |
| Success      | `success`        | ✅ InfoIcon  | Green                | Successfully completed actions  |

#### Status  
Filter notifications based on whether they have been read. These statuses are defined in `STATUS` and styled using `STATUS_STYLE`.

| Status       | Code Value | Icon          | Description                        |
|--------------|------------|---------------|------------------------------------|
| Read         | `read`     | ReadIcon      | Notifications that have been opened |
| Unread       | `unread`   | EnvelopeIcon  | New or untouched notifications      |

### Understanding Notification Logos and Icons

Meshery uses avatar icons to indicate who triggered a notification and what system was involved. These icons help users quickly understand the origin and nature of each event.

<a href="{{ site.baseurl }}/assets/img/notification-center/avatar-icons.png">
  <img src="{{ site.baseurl }}/assets/img/notification-center/avatar-icons.png" alt="Avatar Icons" style="width: 600px; height: auto; margin-top: 10px;" />
</a>

| Icon Type                       | Meaning                                                                 |
|----------------------------------|-------------------------------------------------------------------------|
| Meshery logo only               | System-triggered event – initiated automatically by Meshery (e.g., syncing errors, import failures). |
| User avatar + Meshery logo      | User-triggered event – the user performed an action, and Meshery processed it (e.g., registering a Kubernetes context). |
| User avatar only (rare)         | User-triggered event with no system action involved. |

These icons are generated dynamically using the event’s metadata:
 - If `user_id` is present → shows user avatar.
 - If `system_id` is present → shows Meshery logo.

<details close><summary>Visual Representation of System/User-triggered Notifications</summary>
  <br>

  <figure>
    <figcaption>
      1. 🟢 Meshery-only (System-triggered) Notification
      <a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=a7310bb4-e642-4e4e-807a-dbb602228f07">
        (open in playground)
      </a>
    </figcaption>
  </figure>
  <div id="embedded-design-a7310bb4-e642-4e4e-807a-dbb602228f07" style="height:30rem;width:100%;"></div>
  <script src="{{ site.baseurl }}/assets/img/notification-center/meshery-triggered.js" type="module"></script>

  <figure>
    <figcaption>
      2. 👤+🌐 User + Meshery (User-triggered) Notification
      <a target="_blank" href="https://playground.meshery.io/extension/meshmap?mode=design&design=a96a3008-9c36-4862-b1e3-20cc2c35ca89">
        (open in playground)
      </a>
    </figcaption>
  </figure>
  <div id="embedded-design-a96a3008-9c36-4862-b1e3-20cc2c35ca89" style="height:30rem;width:100%;"></div>
  <script src="{{ site.baseurl }}/assets/img/notification-center/user-triggered.js" type="module"></script>

</details>

### Notification Retention and Visibility
**How long are notifications stored?**
The duration for which notifications are retained is determined by the provider you are using (e.g., Meshery Cloud, local Meshery Server).

**What happens when retention ends?** 
In Meshery Cloud, notifications are removed once the provider is updated, helping ensure the event stream reflects the most recent and relevant information.

{% include related-discussions.html tag="meshery" %}
