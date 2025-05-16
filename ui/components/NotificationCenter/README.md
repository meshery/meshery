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

Here's a refactored and enhanced version of your documentation with improved clarity, structure, and readability:  

---

### Event Formatting Component  

This component serves as the entry point for formatting events. It provides the logic for mapping custom formatters to different event types and properties, ensuring consistent and dynamic formatting.  

### Formatter Hierarchy  

Formatting follows a structured hierarchy to determine the most appropriate transformation for event data:  

1. **Custom Event Formatter**  
   - Defines specialized formatting for specific event types.  
   - Applied when formatting depends on the event's type, action, or category.  
   - Should only be used for events that require highly specific or unique formatting.  

2. **Custom Property Formatter**  
   - Automatically formats specific properties based on their names.  
   - Useful for generic properties that appear across multiple events (e.g., `id`, `data`, `doc`, `description`).  
   - Can also be defined for complex properties like `error` or `trace` to ensure structured formatting.  

### When to Use Which Formatter?  

- Use a **Property Formatter** if the formatting is generic and applies to a specific property across multiple event types (e.g., `description`, `error`).  
- Use an **Event Formatter** if the formatting is unique to a particular event type and does not apply elsewhere.  

### Recursive Dynamic Formatting  

Both **event-based** and **property-based** formatters can recursively utilize the **Dynamic Formatter**, allowing for structured and consistent data transformation across nested properties.  


## Reusability

While this system was initially developed for our events and notification center, the components it comprises are highly reusable and can be employed in other contexts where dynamic formatting of structured data is required.
