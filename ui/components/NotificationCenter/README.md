# Bulk Operations in Notification Center

## Introduction

Bulk operations in the Notification Center allow users to perform actions like deleting multiple notifications or changing the status of multiple notifications in a batch. This documentation outlines the key features and functionality of bulk operations, including the restriction of performing only one bulk operation at a time, the disabling of buttons during ongoing operations, and the display of a loading icon to indicate ongoing activity.

## Features

### 1. Single Bulk Operation at a Time

To maintain system stability and prevent conflicts, the Notification Center enforces a rule that allows users to initiate only one bulk operation at a time. This means that if a bulk operation is already in progress, any attempt to start another operation will be blocked until the ongoing operation is completed or canceled.

### 2. Disabled Buttons During Ongoing Operations

When a bulk operation is initiated, such as deleting or changing the status of notifications, all relevant action buttons are automatically disabled. This prevents users from accidentally triggering additional operations that might interfere with the ongoing one. Disabled buttons ensure the user's actions are synchronized with the current state of notifications.

### 3. Loading Icon Indicator

While a bulk operation is in progress, a loading icon is displayed to provide feedback to the user. This visual indicator assures users that their request is being processed, and it also reinforces the idea that actions cannot be performed until the operation completes or is canceled.

## How It Works

### Initiating a Bulk Operation

1. Users select the notifications they want to include in the bulk operation. This is typically done by checking checkboxes next to each notification.

2. After selecting notifications, users trigger the desired bulk operation (e.g., delete or change status) by clicking the corresponding action button.

3. Once initiated, the bulk operation begins processing the selected notifications.

## Further Enchacements

To ensure a seamless user experience with bulk operations in the Notification Center, consider the following best practices:

1. Provide clear and concise messaging to inform users about the progress and outcome of the bulk operation.

2. Implement error handling to gracefully handle any issues that may arise during the operation and communicate errors to the user.

3. Consider offering a confirmation dialog before initiating a bulk operation to prevent accidental actions.

# Metadata Formatter Documentation

When an event is received from the server, it adheres to a fixed schema containing information that is valuable for presentation to the user. This information typically includes details such as the description, date, user_id, system_id, action, and acted-upon resources. Additionally, sometimes there may be a detailed traceback, a summary, or a comprehensive error log, all of which are dynamically generated data encapsulated within the metadata of the event. Presenting this structured data in a user-friendly manner is a crucial task because it contains valuable insights into ongoing operations.

To accomplish this task, we employ metadata formatters that transform structured data into visually appealing formats. There are currently two types of formatters in use:

1. **Metadata Specific Formatters**: These formatters are specifically designed for particular types of metadata, such as Error and DryRunResponse. Metadata Specific Formatters are implemented as React components that take the metadata as input and render it within the component.

2. **Dynamic Formatter**: Since metadata can vary significantly in structure, it is not practical to create a specific formatter for each kind. Dynamic formatters analyze the schema's structure and apply custom-defined rules for formatting:
   - Text strings are rendered using the _BodySectionRenderer_ (more on this later).
   - Arrays are rendered using the _ArrayRenderer_.
   - Key-value pairs are rendered using the _KeyValueRenderer_.
   - Nested objects are recursively rendered.

## BodySectionRenderer

The BodySectionRenderer is responsible for formatting and rendering raw text strings into React components. During this process, it parses the string to replace external links with `<Link>` components and checks if the link matches predefined sites to render the link accordingly.

## ArrayRenderer

The ArrayRenderer is responsible for rendering an array of items in a recursive manner, presenting them as a bulletized list using the _MetdataFormatter_.

## KeyValueRenderer

Object properties with string values are considered key-value pairs and are rendered as such.

## The Metadata Specific Formatter

Certain metadata, such as Kubernetes responses and Errors, hold high importance and have dedicated renderers. These dedicated renderers can still utilize the dynamic formatter to format specific parts of the response, such as _DryRunResponse_.

## Reusability

While this system was initially developed for our events and notification center, the components it comprises are highly reusable and can be employed in other contexts where dynamic formatting of structured data is required.
