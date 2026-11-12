---
title: Command Palette
description: Fast keyboard navigation, design actions, and setting access across Meshery UI via global shortcut or header search.
categories: [configuration]
weight: 20
---

The Meshery Command Palette is a global, keyboard-first navigation and quick-action menu built directly into Meshery UI. It allows you to jump to any page, execute key workflows (like creating or importing designs), and manage settings from anywhere in the application.

> **Note:** You can open the command palette from anywhere in Meshery UI using <kbd>Ctrl</kbd> + <kbd>K</kbd> (Windows/Linux) or <kbd>⌘</kbd> + <kbd>K</kbd> (macOS), even while focus is in a text field.

## Opening and Closing the Palette

There are two ways to open the command palette:

1. **Keyboard Shortcut**: Press <kbd>Ctrl</kbd> + <kbd>K</kbd> on Windows and Linux, or <kbd>⌘</kbd> + <kbd>K</kbd> on macOS.
2. **Header Button**: Click the **Search...** button (displaying the platform shortcut hint) in the top application header.

To close the palette:
- Press <kbd>Esc</kbd> once.
- Press <kbd>Ctrl</kbd> + <kbd>K</kbd> or <kbd>⌘</kbd> + <kbd>K</kbd> while the palette is open.
- Click outside the dialog on the backdrop.

## Using the Command Palette

The palette is designed for rapid keyboard-driven operation:

- **Filter Commands**: Start typing to immediately filter commands. Matching evaluates command titles, the parent page name shown as subtitle, and keywords.
- **Navigate Items**: Press <kbd>↓</kbd> (Down Arrow) or <kbd>↑</kbd> (Up Arrow) to move the highlighted selection through the matching options.
- **Run a Command**: Press <kbd>Enter</kbd> to execute the currently highlighted command. The palette closes before navigation or action handling begins.

## Command Sections

Available commands are grouped into three contiguous sections:

### 1. Actions

- **Create new Design**: Navigates directly to the Design Configurator (`/configuration/designs/configurator`).
- **Import Design**: Navigates to the Designs page (`/configuration/designs`) and automatically opens the **Import Design** dialog.

### 2. Navigate

Direct navigation to pages across Meshery UI. Sub-pages display their parent section as a subtitle. Pages that are unavailable in your environment are hidden, and pages you lack permission for appear dimmed (see [Permissions and Access Control](#permissions-and-access-control)):

- **Dashboard**
- **Connections**, **Credentials**, **Environments**, **Workspaces**, **Adapters** (shown with the subtitle "Lifecycle")
- **Catalog**, **Designs** (shown with the subtitle "Configuration")
- **Telemetry**, **Charts**, **Metrics**
- **Performance**, **Profiles**
- **Extensions**

Parent sections like "Lifecycle" and "Configuration" duplicate the destinations of their first child pages (Connections and Designs) and are skipped as standalone commands.

### 3. Settings

- **Settings**
- **Preferences**

## Permissions and Access Control

Meshery evaluates commands against your current permissions each time the palette opens:

- **Disabled State with Hints**: Commands for which your account lacks permission stay visible but dimmed with a `Requires: <Permission>` hint and cannot be run.
- **Playground Restrictions**: In restricted playground mode, the Actions, Settings, and Preferences commands are omitted, and navigation follows the restricted menu.

For details on capability keys and permissions, see [Extensibility: Authorization]({{< ref "reference/extensibility/authorization/index.md" >}}).

## Accessibility

The command palette adheres to standard accessibility guidelines:

- **Keyboard-only Operation**: All actions, navigation, and dismissal can be performed entirely via keyboard.
- **Dialog Naming**: The dialog is labelled "Command palette".
- **Combobox Announcements**: The active item is announced through the WAI-ARIA combobox pattern (`aria-activedescendant`), and the result count is announced via an `aria-live="polite"` status region.
- **Theme Awareness**: Works in both light and dark themes.

## Adding Commands (for Contributors)

If you are developing or contributing new features to Meshery UI:

- **Navigation Commands**: Navigation commands are derived from the sidebar configuration in [`ui/components/layout/Navigator/Navigator.tsx`](https://github.com/meshery/meshery/blob/master/ui/components/layout/Navigator/Navigator.tsx) so new pages appear automatically.
- **Action Commands**: Action commands are defined in [`ui/components/command-palette/commands.ts`](https://github.com/meshery/meshery/blob/master/ui/components/command-palette/commands.ts) with a title, destination, and permission specification (no icon). Permission keys come from `Keys` in `@meshery/schemas/permissions`.
