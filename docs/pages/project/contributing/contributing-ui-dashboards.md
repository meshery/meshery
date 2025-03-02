---
layout: page
title: "Meshery Dashboard Contribution Guide"
permalink: project\contributing\contributing-ui-dashboards.md
abstract: "Guidelines for contributing to Meshery's dashboard framework, including customization, integration, and best practices."
language: en
type: project
category: contributing
list: include
---


# Contributing to Meshery Dashboards

Meshery's dashboard framework is designed to be highly extensible, allowing developers to create custom widgets and layouts that enhance the user experience. This guide will walk you through the process of extending Meshery's dashboards, including customization options, integration points, and best practices.

## 1. Dashboard Framework

### 1.1 Overview

Meshery's dashboard framework is built on a flexible grid system that allows widgets to be arranged in various layouts. The framework manages:

- **Layout**: Widget positioning and arrangement on the dashboard
- **Sizing**: Responsive resizing of widgets based on screen size and user preferences
- **Widget Thumbnails**: Preview representations of available widgets
- **User Preferences**: Persistent storage of dashboard configurations

The dashboard state is stored in user preferences, allowing configurations to persist between sessions and be synced across devices when using a remote provider.

### 1.2 Grid System

Dashboards use a 12-column grid layout that supports responsive design principles. Widgets can span multiple columns and rows, with the layout automatically adjusting based on screen size.

```javascript
// Example widget layout configuration
{
  id: "unique-widget-id",
  title: "My Custom Widget",
  x: 0,       // x-position on grid
  y: 0,       // y-position on grid
  w: 4,       // width (in columns)
  h: 2,       // height (in rows)
  component: MyWidgetComponent,
  minW: 2,    // minimum width
  minH: 1     // minimum height
}
```

### 1.3 Widget Registry

All widgets must be registered in the widget registry before they can be used in dashboards. The registry maintains metadata about each widget, including its type, display name, and component reference.

For more information on the core dashboard architecture, see the [Dashboard Architecture Documentation](https://docs.meshery.io/concepts/architectural/dashboard).

## 2. Setting Up the Development Environment

Before you can start extending Meshery's dashboards, you need to set up your development environment. This section provides a step-by-step guide to get you started.

### 2.1 Prerequisites

Make sure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.18 or later)
- [Node.js](https://nodejs.org/) (version 16 or later)
- [npm](https://www.npmjs.com/) (version 8 or later)
- [Docker](https://docs.docker.com/get-docker/) (for running dependencies)
- [Git](https://git-scm.com/downloads)

### 2.2 Forking and Cloning the Repository

1. Fork the Meshery repository on GitHub by visiting [https://github.com/meshery/meshery](https://github.com/meshery/meshery) and clicking the "Fork" button.

2. Clone your forked repository:

```bash
git clone https://github.com/YOUR-USERNAME/meshery.git
cd meshery
```

3. Add the upstream repository as a remote:

```bash
git remote add upstream https://github.com/meshery/meshery.git
```

### 2.3 Setting Up the UI Development Environment

The Meshery UI is built with Next.js and React. To set up the UI development environment:

1. Navigate to the UI directory:

```bash
cd ui
```

2. Install dependencies:

```bash
npm install
```

3. Create a `.env.local` file in the `ui` directory with the following content:

```
NEXT_PUBLIC_MESHERY_API=http://localhost:9081
```

4. Start the UI development server:

```bash
npm run dev
```

This will start the UI on `http://localhost:3000`.

### 2.4 Setting Up the Server

In a separate terminal, you'll need to run the Meshery server:

1. Navigate to the root directory of the Meshery repository.

2. Build and run the server:

```bash
make run-backend
```

This will start the Meshery server on `http://localhost:9081`.

### 2.5 Development Workflow

The recommended workflow for dashboard development is:

1. Run the Meshery server (`make run-backend`)
2. Run the UI development server (`cd ui && npm run dev`)
3. Access the Meshery UI at `http://localhost:3000`
4. Make changes to the UI code in the `ui/components` and `ui/pages` directories
5. The UI will automatically reload when you save changes

### 2.6 Useful Development Commands

- `npm run lint` - Run linting checks
- `npm run test` - Run UI tests
- `npm run build` - Build the UI for production
- `make build-ui` - Build the UI from the root directory
- `make docker` - Build a Docker image with your changes

For more detailed setup instructions, refer to the [Meshery Development Guide](https://docs.meshery.io/project/contributing/contributing-server).

## 3. Local vs. Remote Provider Considerations

When developing dashboard extensions, it's important to understand the differences between local and remote providers, as they affect how dashboard configurations are stored and synchronized.

### 3.1 Local Provider

With a local provider:
- Dashboard configurations are stored in the browser's local storage
- Configurations do not persist across different devices
- No user authentication is required
- Limited to basic dashboard functionality

When developing widgets for local provider use, ensure they can function without requiring authenticated API access.

### 3.2 Remote Provider

With a remote provider (such as Meshery Cloud):
- Dashboard configurations are stored on the server and synced across devices
- User authentication is required
- Enhanced persistence and sharing capabilities
- Support for team collaboration features

When developing widgets for remote provider use, you can take advantage of authenticated APIs and user-specific data.

For more information on providers, refer to the [Meshery Providers Documentation](https://docs.meshery.io/concepts/logical/providers).

## 4. Creating a Custom Widget

### 4.1 Widget Structure

A Meshery widget consists of:
1. A React component that renders the widget content
2. Configuration metadata for the dashboard framework
3. Optional API integrations to fetch and display data

### 4.2 Widget File Structure

Organize your widget files following this recommended structure:

```
ui/
├── components/
│   ├── DashboardComponents/
│   │   ├── Widgets/
│   │   │   ├── MyCustomWidget/
│   │   │   │   ├── index.js           # Main widget component
│   │   │   │   ├── styles.js          # Component styles
│   │   │   │   ├── helpers.js         # Helper functions
│   │   │   │   ├── config.js          # Widget configuration
│   │   │   │   └── README.md          # Documentation
```

### 4.3 Example Widget Implementation

Below is an example of a simple metrics widget:

```jsx
import React, { useState, useEffect } from 'react';
import { useQuery } from '@apollo/client';
import { GET_METRICS_DATA } from '../graphql/queries';
import { Card, CardContent, Typography } from '@material-ui/core';
import { LineChart, Line, XAxis, YAxis, Tooltip } from 'recharts';

const MetricsWidget = ({ widgetId, refreshInterval = 30000 }) => {
  const { loading, error, data, refetch } = useQuery(GET_METRICS_DATA, {
    variables: { widgetId },
    fetchPolicy: 'network-only',
  });

  // Refresh data at specified interval
  useEffect(() => {
    const interval = setInterval(() => {
      refetch();
    }, refreshInterval);
    
    return () => clearInterval(interval);
  }, [refetch, refreshInterval]);

  if (loading) return <Typography>Loading metrics...</Typography>;
  if (error) return <Typography color="error">Error loading metrics: {error.message}</Typography>;

  return (
    <Card>
      <CardContent>
        <Typography variant="h6">Service Metrics</Typography>
        <LineChart width={300} height={200} data={data.metrics}>
          <XAxis dataKey="timestamp" />
          <YAxis />
          <Tooltip />
          <Line type="monotone" dataKey="value" stroke="#8884d8" />
        </LineChart>
      </CardContent>
    </Card>
  );
};

export default MetricsWidget;
```

### 4.4 Registering the Widget

To make your widget available in the dashboard, you need to register it in the widget registry:

```javascript
import { registerWidget } from '../utils/widgetUtils';
import MetricsWidget from '../components/MetricsWidget';

registerWidget({
  id: 'metrics-widget',
  name: 'Service Metrics',
  description: 'Displays real-time service metrics',
  component: MetricsWidget,
  category: 'metrics',
  thumbnail: '/images/widgets/metrics-thumbnail.png',
  defaultSize: {
    w: 4,
    h: 2,
    minW: 2,
    minH: 1,
  }
});
```

For a complete example of widget implementation, see the [Example Widgets Directory](https://docs.meshery.io/guides/sample-apps/widgets-examples).

## 5. API Integration

Meshery provides both REST and GraphQL APIs that can be used to fetch data for your widgets.

### 5.1 REST API

The REST API is ideal for simple data retrieval and actions. When integrating with the REST API:

- Use the Meshery API client library when available
- Include proper authentication headers for authenticated endpoints
- Implement error handling and loading states
- Consider caching strategies for performance optimization

Example REST API integration:

```javascript
import { mesheryApiClient } from '../utils/apiUtils';

async function fetchPerformanceData(serviceId) {
  try {
    const response = await mesheryApiClient.get(`/api/performance/results/${serviceId}`);
    return response.data;
  } catch (error) {
    console.error('Error fetching performance data:', error);
    throw error;
  }
}
```

For more information, refer to the [Meshery REST API Documentation](https://docs.meshery.io/guides/meshery-api).

### 5.2 GraphQL API

The GraphQL API provides more flexibility and is ideal for complex data requirements. Benefits include:

- Fetching only the required data fields
- Combining multiple data requirements in a single request
- Built-in subscription support for real-time updates

Example GraphQL integration:

```javascript
import { useQuery, gql } from '@apollo/client';

const GET_MESH_DATA = gql`
  query GetMeshData($meshId: ID!) {
    mesh(id: $meshId) {
      id
      name
      services {
        id
        name
        status
        metrics {
          cpu
          memory
          latency
        }
      }
    }
  }
`;

function MeshOverviewWidget({ meshId }) {
  const { loading, error, data } = useQuery(GET_MESH_DATA, {
    variables: { meshId },
    pollInterval: 10000,
  });
  
  // Component rendering logic
}
```

For complete API reference, see the [Meshery GraphQL API Documentation](https://docs.meshery.io/guides/meshery-graphql).

## 6. Permissions

Meshery's permission system is a crucial aspect of dashboard and widget development. Understanding and properly implementing permissions ensures that users only access functionality and data appropriate to their role and authorization level.

### 6.1 Permission Architecture

Meshery uses a role-based access control (RBAC) system with the following key components:

- **Roles**: Collections of permissions assigned to users (e.g., Admin, Operator, Viewer)
- **Resources**: Items that can be accessed or manipulated (e.g., dashboards, services, policies)
- **Actions**: Operations that can be performed on resources (e.g., view, create, edit, delete)
- **Permission Keys**: Unique identifiers that represent the combination of actions and resources

### 6.2 Permission Keys

Permission keys follow the format: `{action}:{resource}[:{subresource}]`

Examples:
- `view:dashboards` - Permission to view all dashboards
- `edit:dashboards:system` - Permission to edit system dashboards
- `create:widgets:performance` - Permission to create performance widgets

### 6.3 Implementing Permission-Aware Widgets

When developing widgets, follow these best practices for permissions:

#### 6.3.1 Check Permissions Before Rendering Sensitive UI Elements

```jsx
import { usePermissions } from '../hooks/usePermissions';

function ConfigurableWidget({ resourceId }) {
  const { checkPermission } = usePermissions();
  
  const canView = checkPermission('view', `widgets:${resourceId}`);
  const canEdit = checkPermission('edit', `widgets:${resourceId}`);
  const canDelete = checkPermission('delete', `widgets:${resourceId}`);
  
  if (!canView) {
    return <AccessDeniedPlaceholder />;
  }
  
  return (
    <div className="widget-container">
      <div className="widget-content">
        {/* Widget main content here */}
      </div>
      
      <div className="widget-actions">
        {canEdit && <EditButton onClick={handleEdit} />}
        {canDelete && <DeleteButton onClick={handleDelete} />}
      </div>
    </div>
  );
}
```

#### 6.3.2 Permission-Based API Requests

Filter API requests based on user permissions to prevent unnecessary or unauthorized calls:

```jsx
import { useEffect, useState } from 'react';
import { usePermissions } from '../hooks/usePermissions';
import { fetchDetailedMetrics, fetchBasicMetrics } from '../api/metrics';

function MetricsWidget() {
  const [metrics, setMetrics] = useState(null);
  const { checkPermission } = usePermissions();
  const hasAdvancedAccess = checkPermission('view', 'metrics:detailed');
  
  useEffect(() => {
    async function loadData() {
      try {
        if (hasAdvancedAccess) {
          const detailedData = await fetchDetailedMetrics();
          setMetrics(detailedData);
        } else {
          const basicData = await fetchBasicMetrics();
          setMetrics(basicData);
        }
      } catch (error) {
        console.error('Failed to load metrics:', error);
      }
    }
    
    loadData();
  }, [hasAdvancedAccess]);
  
  // Render logic
}
```

#### 6.3.3 Graceful Permission Handling

Instead of showing error messages or blank screens, implement graceful fallbacks:

```jsx
function DashboardWidget({ widgetType }) {
  const { checkPermission } = usePermissions();
  
  // Different widget types require different permissions
  const permissionMap = {
    'performance': 'view:performance',
    'security': 'view:security',
    'config': 'view:configuration'
  };
  
  const requiredPermission = permissionMap[widgetType];
  const hasPermission = checkPermission(requiredPermission);
  
  if (!hasPermission) {
    return (
      <div className="limited-access-widget">
        <h3>Limited Access</h3>
        <p>This widget shows {widgetType} metrics that require additional permissions.</p>
        <a href="/docs/permissions">Learn more about permissions</a>
      </div>
    );
  }
  
  // Full widget implementation
}
```

### 6.4 Declaring Required Permissions

When registering widgets, declare required permissions to help the dashboard framework make intelligent decisions about widget visibility and behavior:

```javascript
registerWidget({
  id: 'advanced-metrics-widget',
  name: 'Advanced Service Metrics',
  description: 'Displays detailed service metrics with advanced controls',
  component: AdvancedMetricsWidget,
  category: 'metrics',
  requiredPermissions: [
    'view:metrics:detailed',
    'view:services'
  ],
  // Optional permissions that enable additional features
  optionalPermissions: {
    'edit:metrics:thresholds': 'Enables threshold configuration',
    'create:alerts': 'Enables alert creation from the widget'
  }
});
```

### 6.5 Widget Adaptation Based on Permissions

Design widgets to adapt their functionality based on available permissions:

```jsx
function AdaptiveWidget() {
  const { checkPermission } = usePermissions();
  
  // Check various permission levels
  const permissionLevels = {
    basic: checkPermission('view:basic'),
    intermediate: checkPermission('view:intermediate'),
    advanced: checkPermission('view:advanced'),
    admin: checkPermission('view:admin')
  };
  
  // Determine the highest permission level the user has
  let activeLevel = 'basic';
  if (permissionLevels.admin) activeLevel = 'admin';
  else if (permissionLevels.advanced) activeLevel = 'advanced';
  else if (permissionLevels.intermediate) activeLevel = 'intermediate';
  
  // Render appropriate content based on permission level
  return (
    <div className="adaptive-widget">
      <WidgetHeader level={activeLevel} />
      
      {/* Base content visible to all users */}
      <BasicContent />
      
      {/* Progressive enhancement based on permissions */}
      {permissionLevels.intermediate && <IntermediateContent />}
      {permissionLevels.advanced && <AdvancedContent />}
      {permissionLevels.admin && <AdminContent />}
    </div>
  );
}
```

### 6.6 Testing with Different Permission Levels

When developing widgets, test with various permission configurations:

1. Create test user accounts with different roles
2. Implement a permission override mode during development
3. Use permission mocking in unit and integration tests

Example test helper:

```javascript
// permissionTestHelper.js
export function mockPermissions(permissions) {
  // Store original implementation
  const original = window.meshery.checkPermission;
  
  // Override with mock
  window.meshery.checkPermission = (action, resource) => {
    const key = `${action}:${resource}`;
    return permissions.includes(key);
  };
  
  // Return cleanup function
  return () => {
    window.meshery.checkPermission = original;
  };
}

// In tests
import { mockPermissions } from './permissionTestHelper';

describe('AdminWidget', () => {
  test('shows all controls with admin permissions', () => {
    const cleanup = mockPermissions(['view:admin', 'edit:settings', 'delete:resources']);
    // Test widget with admin permissions
    cleanup();
  });
  
  test('hides sensitive controls without admin permissions', () => {
    const cleanup = mockPermissions(['view:basic']);
    // Test widget with basic permissions
    cleanup();
  });
});
```

For more information on the permission system, see the [Meshery RBAC Documentation](https://docs.meshery.io/concepts/logical/identity).

## 7. Additional Examples

Here are other examples of dashboard extensions to inspire your development:

### 7.1 Performance Testing Widget

A widget that allows users to run performance tests and view results directly from the dashboard.

See the [Performance Widget Example](https://docs.meshery.io/guides/performance-management/widgets).

### 7.2 Service Mesh Topology Viewer

A widget that visualizes the topology of a service mesh, showing connections between services.

See the [Topology Widget Example](https://docs.meshery.io/guides/service-mesh-management/topology-widget).

### 7.3 Policy Compliance Dashboard

A widget that displays compliance status for various policies applied to the mesh.

See the [Policy Widget Example](https://docs.meshery.io/guides/policies/policy-widgets).

## 8. Publishing Extensions

Once you've developed and tested your dashboard extension, you can publish it to the Meshery Extensions Marketplace.

### 8.1 Requirements

To publish your extension:

1. Create a GitHub repository for your extension
2. Include comprehensive documentation
3. Provide usage examples
4. Ensure code quality and test coverage
5. Create a release with semantic versioning

### 8.2 Submission Process

1. Fork the Meshery Extensions repository
2. Add your extension metadata to the catalog
3. Submit a pull request
4. Address any review comments
5. Once approved, your extension will be listed on the marketplace

### 8.3 Extension Metadata

```yaml
name: "My Custom Dashboard Widget"
version: "1.0.0"
description: "A widget that displays custom metrics"
icon: "https://path-to-your-icon.png"
category: "Metrics"
maintainer: "Your Name <your.email@example.com>"
repo: "https://github.com/yourusername/meshery-widget"
requiredPermissions:
  - "view:metrics"
  - "view:services"
optionalPermissions:
  - "edit:metrics:thresholds"
  - "create:alerts"
```

For detailed publishing guidelines, see the [Meshery Extensions Marketplace Documentation](https://docs.meshery.io/extensions).

## 9. Best Practices

### 9.1 Performance Considerations

- Implement efficient data fetching and caching strategies
- Use pagination for large datasets
- Implement virtualization for long lists
- Optimize rendering cycles

### 9.2 User Experience

- Provide loading and error states
- Use consistent styling with Meshery's design system
- Implement responsive layouts
- Provide tooltips and help text for complex functionality

### 9.3 Code Quality

- Write unit and integration tests
- Document your code thoroughly
- Follow Meshery's coding conventions
- Use TypeScript for type safety

### 9.4 Accessibility

- Ensure proper color contrast
- Add ARIA attributes where needed
- Support keyboard navigation
- Test with screen readers

### 9.5 Permission-Related Best Practices

- **Granular Permission Checks**: Check specific permissions rather than general role access
- **Progressive Enhancement**: Build UIs that enhance functionality as more permissions are available
- **Helpful Access Denied Messages**: Explain why content is unavailable and how to request access
- **Performance Optimization**: Cache permission check results to avoid redundant checks
- **Documentation**: Document all required and optional permissions for your widget
- **Testing**: Test widgets with various permission combinations to ensure proper behavior
- **Security First**: Always validate permissions on both client and server sides
- **Error Handling**: Implement graceful error handling for permission-related API failures

## 10. Troubleshooting

Common issues and solutions when developing dashboard extensions:

### 10.1 Development Environment Issues

- **UI not connecting to backend**: Verify that your `.env.local` file has the correct API URL and that the backend server is running
- **Hot reloading not working**: Restart the development server with `npm run dev`
- **Module not found errors**: Run `npm install` to ensure all dependencies are installed

### 10.2 Widget Development Issues

1. **Widget not appearing in dashboard**: Ensure it's properly registered in the widget registry
2. **Data not loading**: Check API endpoints and authentication
3. **Layout issues**: Verify the grid configuration and responsive behavior
4. **Performance problems**: Look for unnecessary re-renders or inefficient data fetching
5. **Permission issues**: Verify that permission keys match the expected format and that permission checks are implemented correctly

### 10.3 Common Error Messages

- **"Failed to fetch"**: Check that your backend server is running and accessible
- **"Widget registration failed"**: Ensure your widget ID is unique and that the component is properly exported
- **"Permission denied"**: Check that you have the required permissions for the action you're attempting

For more troubleshooting guidance, see the [Meshery Troubleshooting Guide](https://docs.meshery.io/reference/troubleshooting).

## 11. Contributing Back

If you've developed a useful dashboard extension, consider contributing it back to the Meshery project:

1. Open an issue describing your extension
2. Submit a pull request with your code
3. Work with the maintainers to integrate it into the Meshery codebase

## 12. Resources

- [Meshery UI Architecture](https://docs.meshery.io/project/contributing/contributing-ui)
- [Dashboard API Reference](https://docs.meshery.io/reference/api/dashboard)
- [Widget Component Library](https://docs.meshery.io/reference/ui-components)
- [Extension Examples Repository](https://github.com/meshery/meshery-extensions-examples)
- [Meshery Community Slack](https://slack.meshery.io) - Join the #ui-extensions channel