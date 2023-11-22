---
layout: default
title: Managing Events with Notification Center
description: A Comprehensive Guide to Managing Events with Meshery
permalink: guides/events-management
type: guides
language: en
abstract: "Meshery tracks operations that you perform on infrastructurees and their workloads. Meshery provides notification of environment issues, application conflicts with infrastructure configuration, and so on."
---

# Events Management in Meshery

Meshery Server, while running, generates a variety of events for operations happening inside your Kubernetes and Cloud-Native clusters and Meshery Server itself. To observe, monitor, and effectively manage these events, Meshery provides a powerful and user-friendly Notification Center.

## What is the Notification Center?

The Notification Center in Meshery's user interface is a dedicated feature that empowers you to monitor and manage the events occurring within your clusters. This guide will walk you through the various aspects of utilizing the Notification Center effectively.

## Key Features of the Notification Center

<a href="{{ site.baseurl }}/assets/img/notification-center/NotificationCenterOverview.png"><img alt="Notification Center Overview" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/notification-center/NotificationCenterOverview.png" /></a>

### 1. Views and Event Management

The Notification Center serves as your central hub for viewing and managing events within your environment. Whether you are monitoring a single cluster or an extensive mesh of interconnected services, Meshery's Notification Center provides you with a consolidated and comprehensive view of all events.

### 2. Event Search

To simplify the process of pinpointing specific events, Meshery offers a powerful search functionality. You can search for events based on various parameters, such as the action that occurred, the type of operation, the actor responsible for the event, and more. This feature is invaluable for quickly identifying and addressing critical events within your infrastructure.

### 3. Flexible Filtering Support

<a href="{{ site.baseurl }}/assets/img/notification-center/NotificationCenterFiltering.png"><img alt="Notification Center Filtering" style="width:500px;height:auto;" src="{{ site.baseurl }}/assets/img/notification-center/NotificationCenterFiltering.png" /></a>

Meshery's Notification Center offers robust filtering capabilities. You can customize the view by applying filters based on multiple criteria. For example, you can filter events based on their significance, source, or any other attributes that are important to you. These filters make it easy to focus on specific event categories and streamline your event management process.

### 4. Event Descriptions and Data

Each event in Meshery's Notification Center contains a detailed description to provide you with a clear understanding of what happened. This includes information about the action, the entities involved, and any relevant data. Additionally, events may include tracebacks and other contextual data to help you troubleshoot and diagnose issues effectively.

### 5. Data Sharing

Meshery allows you to easily share notifications with relevant stakeholders. Whether you need to collaborate with team members, share event details with developers, or involve other parties in the resolution process, Meshery simplifies the sharing of event information, facilitating efficient communication and issue resolution.

### 5. Event Status Tracking

Meshery's Notification Center provides the ability to track the status of events. You can mark events as read or unread, helping you to prioritize and manage your response to events efficiently. This status tracking feature ensures that no critical event goes unnoticed, and you can easily keep track of events that require immediate attention.

### 6. Event Persistence

Events in Meshery are always persisted, ensuring that you have a complete historical record of activities in your clusters and Meshery Server. This persistence is invaluable for compliance, auditing, and historical analysis, enabling you to track the evolution of your infrastructure and detect patterns over time.

## Summary

Meshery's Notification Center simplifies the task of managing events in your Kubernetes and Cloud-Native environments. This robust feature offers a central hub to view, search, and filter events, making it effortless to identify and address critical incidents. Each event comes with detailed descriptions and additional data for efficient troubleshooting.

You can easily collaborate with stakeholders by sharing event data, enhancing communication and issue resolution. What's more, Meshery enables you to track event statuses, ensuring no crucial event goes unnoticed. Plus, events are always persisted, providing a historical record for compliance and trend analysis. With Meshery's Notification Center, you can gain deep insights into your infrastructure, swiftly respond to critical events, and maintain optimal performance for your applications and services. Detailed instructions and best practices can be found in the official documentation and Meshery's user guide.

{% include suggested-reading.html language="en" %}

{% include related-discussions.html tag="meshery" %}
