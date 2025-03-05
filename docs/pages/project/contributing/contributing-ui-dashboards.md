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


## 2. Local vs. Remote Provider Considerations

When developing dashboard extensions, it's important to understand the differences between local and remote providers, as they affect how dashboard configurations are stored and synchronized.

### 2.1 Local Provider

With a local provider:
- Dashboard configurations are stored in the browser's local storage
- Configurations do not persist across different devices
- No user authentication is required
- Limited to basic dashboard functionality

When developing widgets for local provider use, ensure they can function without requiring authenticated API access.

### 2.2 Remote Provider

With a remote provider (such as Meshery Cloud):
- Dashboard configurations are stored on the server and synced across devices
- User authentication is required
- Enhanced persistence and sharing capabilities
- Support for team collaboration features

When developing widgets for remote provider use, you can take advantage of authenticated APIs and user-specific data.

For more information on providers, refer to the [Meshery Providers Documentation](https://docs.meshery.io/concepts/logical/providers).

## 3. Widget Thumbnails and User Preferences

### 3.1 Widget Thumbnails

Widget thumbnails are preview images that represent available widgets in the dashboard. They help users quickly identify and select the widgets they want to add to their dashboards.

To add a thumbnail to your widget:

1. Create a PNG image with a resolution of 200x150 pixels.
2. Save the image in the ui/public/images/widgets directory.
3. Reference the image in your widget configuration.

Example:

```javascript
// Example widget layout configuration with thumbnail
{
  id: "my-custom-widget",
  title: "My Custom Widget",
  component: MyCustomWidgetComponent,
  thumbnail: "/images/widgets/my-custom-widget.png"
}
```

## 3.2 User Preferences
User preferences allow users to customize their dashboard experience. These preferences are stored and persisted between sessions, providing a consistent experience across devices.

The dashboard framework automatically manages user preferences, including:

Widget positions and sizes
Dashboard layout
Theme settings
Other user-specific settings
To access and modify user preferences:

1. Use the useDashboardPreferences hook.
2. Get or set preferences using the provided API.

Example:
```javascript
import { useDashboardPreferences } from "../hooks/useDashboardPreferences";

function MyWidget() {
  const { preferences, setPreferences } = useDashboardPreferences();

  const handleThemeChange = (theme) => {
    setPreferences({ ...preferences, theme: theme });
  };

  return (
    <div>
      <h2>My Widget</h2>
      <button onClick={() => handleThemeChange("light")}>Light Theme</button>
      <button onClick={() => handleThemeChange("dark")}>Dark Theme</button>
    </div>
  );
}
```
## 4. API Integrations

### 4.1 REST API
Meshery provides a REST API for interacting with various resources. When developing widgets, you can use the REST API to fetch data and perform actions.
Meshery provides a REST API available through the default port of 9081/tcp at <hostname>:<port>/api/. 

```javascript
async function fetchPerformanceData() {
  try {
    const response = await fetch('/api/performance');
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error('Error fetching performance data:', error);
    throw error;
  }
}
```
For more information, refer to the [Meshery REST API Documentation](https://docs.meshery.io/extensibility/api).

### 4.2 GraphQL API
Meshery provides a GraphQL API for more flexible data queries. The GraphQL API is available through the default port of 9081/tcp at `<hostname>:<port>/api/graphql/query`. Relay is the client used for interacting with the GraphQL API.

Example GraphQL query:

```javascript
import { graphql } from 'react-relay';

const query = graphql`
  query MyWidgetQuery {
    performance {
      id
      name
      value
    }
  }
`;

async function fetchGraphQLData() {
  const response = await fetch('/api/graphql/query', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: query.loc.source.body,
    }),
  });
  const result = await response.json();
  return result.data;
}
```
For more information, refer to the [Meshery GraphQL API Documentation](https://docs.meshery.io/extensibility/api).

## 5. Permissions
Meshery's permission system is a crucial aspect of dashboard and widget development. Understanding and properly implementing permissions ensures that users only access functionality and data appropriate to their role and authorization level.

### 5.1 Permission Architecture
Meshery uses a role-based access control (RBAC) system with the following key components:

Roles: Collections of permissions assigned to users (e.g., Admin, Operator, Viewer)
Resources: Items that can be accessed or manipulated (e.g., dashboards, services, policies)
Actions: Operations that can be performed on resources (e.g., view, create, edit, delete)
Permission Keys: Unique identifiers that represent the combination of actions and resources

### 5.2 Permission Keys
Permission keys follow the format: {action}:{resource}[:{subresource}]

Examples:

1. view:dashboards - Permission to view all dashboards
2. edit:dashboards:system - Permission to edit system dashboards
3. create:widgets:performance - Permission to create performance widgets

## 6. Additional Examples
### Widget Implementation Examples.

The `MyDesignsWidget` is an excellent reference implementation of a dashboard widget. This widget manages designs, allowing users to create, choose templates, or import from GitHub.

```javascript
import React, { useEffect, useState } from 'react';
import { createTheme, darken, lighten, useMediaQuery } from '@material-ui/core';
import { useTheme } from '@material-ui/styles';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { Box, Button, CircularProgress, Grid, Paper, Typography } from '@material-ui/core';
import { updateProgress } from '../../../lib/store';
import { DASHBOARD_DESIGN_LINK } from '../../../constants/navigator';
import useStyles from '../../../assets/styles/general/tool.styles';
import DesignConfigurator from '../../configuratorComponents/MeshModel';
import dataFetch from '../../../lib/data-fetch';
import { iconMedium } from '../../../css/icons.styles';

const MyDesignsWidget = ({ updateProgress }) => {
  const router = useRouter();
  const classes = useStyles();

  // Widget states
  const [designs, setDesigns] = useState([]);
  const [selectedDesign, setSelectedDesign] = useState(null);
  const [page, setPage] = useState(0);
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [importConfig, setImportConfig] = useState({
    showImportModal: false,
    selectedProject: '',
    githubURL: '',
    isLocalFile: false,
  });
  const [loadingDesigns, setLoadingDesigns] = useState(false);

  // Fetch designs on component mount
  useEffect(() => {
    fetchDesigns();
  }, [page, pageSize]);

  // Fetch designs from API
  const fetchDesigns = () => {
    setLoadingDesigns(true);
    updateProgress({ showProgress: true });

    dataFetch(
      `/api/designs?page=${page}&pagesize=${pageSize}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (result) => {
        updateProgress({ showProgress: false });
        setLoadingDesigns(false);
        if (result) {
          setDesigns(result.designs || []);
          setCount(result.total_count || 0);
          setPageSize(result.page_size || 10);
        }
      },
      handleError,
    );
  };

  // Error handling
  const handleError = (error) => {
    updateProgress({ showProgress: false });
    setLoadingDesigns(false);
    console.error('Error fetching designs:', error);
  };

  // Navigate to create new design
  const handleNewDesign = () => {
    router.push(DASHBOARD_DESIGN_LINK);
  };

  // Render design cards
  const renderDesigns = () => {
    if (loadingDesigns) {
      return (
        <Box display="flex" justifyContent="center" alignItems="center" p={2}>
          <CircularProgress />
        </Box>
      );
    }

    if (designs.length === 0) {
      return (
        <Box p={2}>
          <Typography variant="body1" align="center">
            No designs found. Create your first design!
          </Typography>
        </Box>
      );
    }

    return (
      <Grid container spacing={2}>
        {designs.slice(0, 3).map((design) => (
          <Grid item xs={12} md={4} key={design.id}>
            <Paper 
              className={classes.designCard}
              onClick={() => router.push(`/designs/${design.id}`)}
            >
              <Typography variant="h6" noWrap title={design.name}>
                {design.name}
              </Typography>
              <Typography variant="body2" color="textSecondary">
                {design.pattern_file ? 'Pattern' : 'Component'}
              </Typography>
              <Typography variant="caption" color="textSecondary">
                Last updated: {new Date(design.updated_at).toLocaleDateString()}
              </Typography>
            </Paper>
          </Grid>
        ))}
      </Grid>
    );
  };

  return (
    <div className={classes.widgetContainer}>
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb={2}
      >
        <Typography variant="h6">My Designs</Typography>
        <Button
          variant="contained"
          color="primary"
          size="small"
          onClick={handleNewDesign}
        >
          Create
        </Button>
      </Box>
      {renderDesigns()}
      {designs.length > 0 && (
        <Box mt={2} display="flex" justifyContent="flex-end">
          <Link href="/designs">
            <Button color="primary">View All</Button>
          </Link>
        </Box>
      )}
    </div>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(null, mapDispatchToProps)(MyDesignsWidget);

// Widget configuration for dashboard
export const MyDesignsWidgetConfig = {
  id: 'my-designs',
  title: 'My Designs',
  component: MyDesignsWidget,
  thumbnail: '/static/img/designs-light.svg',
  minW: 4,
  minH: 2,
  defaultW: 6,
  defaultH: 3
};
```

### Adding Your Widget to the Dashboard

To make your widget available in the dashboard, you need to include it in the `getWidgets.js` file. Here's how widgets are typically structured in this file:

```javascript
import MyDesignsWidget from './MyDesignsWidget';
import WorkspaceActivityWidget from './WorkspaceActivityWidget';
import HelpCenterWidget from './HelpCenterWidget';
// Other widget imports...

export function getWidgets() {
  return [
    {
      id: "my-designs",
      title: "My Designs",
      x: 0,
      y: 0,
      w: 4,
      h: 2,
      component: MyDesignsWidget,
      minW: 2,
      minH: 1
    },
    {
      id: "workspace-activity",
      title: "Workspace Activity",
      x: 4,
      y: 0,
      w: 4,
      h: 2,
      component: WorkspaceActivityWidget,
      minW: 2,
      minH: 1
    },
    {
      id: "help-center",
      title: "Help Center",
      x: 8,
      y: 0,
      w: 4,
      h: 2,
      component: HelpCenterWidget,
      minW: 2,
      minH: 1
    },
    // Other widget configurations...
  ];
}
```

## More Widget Examples

For additional examples of different widget implementations, Meshery offers several reference widgets:

1. **WorkspaceActivityWidget**: Displays recent activities in the workspace.

2. **HelpCenterWidget**: Provides links to help resources like documentation, support, and forums.

3. **LatestBlogs**: Displays the latest blog posts from the Meshery community.

4. **ClusterStatus**: Shows the status of Kubernetes clusters.

For a complete list of available widgets and their configurations, you can refer to the `getWidgets.js` file in the Meshery codebase.

## 7. Publishing Extensions
Once you've developed and tested your dashboard extension, you can publish it to the Meshery Extensions Marketplace. The Meshery Extension Hub at [meshery.io/extensions](https://meshery.io/extensions) serves as the central repository for all Meshery extensions, including dashboard widgets and components.

### 7.1 Requirements
To publish your extension:
1. Create a GitHub repository for your extension
2. Include comprehensive documentation
3. Provide usage examples
4. Ensure code quality and test coverage
5. Create a release with semantic versioning

### 7.2 Submission Process

1. Fork the Meshery Extensions repository
2. Add your extension metadata to the catalog
3. Submit a pull request
4. Address any review comments
5. Once approved, your extension will be listed on the marketplace

### 7.3 Extension Metadata

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


## 8. Best Practices
### 8.1 Performance Considerations
- Implement efficient data fetching and caching strategies
- Use pagination for large datasets
- Implement virtualization for long lists
- Optimize rendering cycles

### 8.2 User Experience

- Provide loading and error states
- Use consistent styling with Meshery's design system
- Implement responsive layouts
- Provide tooltips and help text for complex functionality

### 8.3 Code Quality
- Write unit and integration tests
- Document your code thoroughly
- Follow Meshery's coding conventions
- Use TypeScript for type safety

### 8.4 Accessibility

- Ensure proper color contrast
- Add ARIA attributes where needed
- Support keyboard navigation
- Test with screen readers

### 8.5 Permission-Related Best Practices

- **Granular Permission Checks**: Check specific permissions rather than general role access
- **Progressive Enhancement**: Build UIs that enhance functionality as more permissions are available
- **Helpful Access Denied Messages**: Explain why content is unavailable and how to request access
- **Performance Optimization**: Cache permission check results to avoid redundant checks
- **Data not loading**: Check API endpoints and authentication
- **Layout issues**: Verify the grid configuration and responsive behavior
- **Performance problems**: Look for unnecessary re-renders or inefficient data fetching
- **Permission issues**: Verify that permission keys match the expected format and that permission checks are implemented correctly.

For more troubleshooting guidance, see the [Meshery Troubleshooting Guide](https://docs.meshery.io/reference/error-codes).

## 10. Contributing Back
If you've developed a useful dashboard extension, consider contributing it back to the Meshery project:
1. Open an issue describing your extension
2. Submit a pull request with your code
3. Work with the maintainers to integrate it into the Meshery codebase
For more information on contributing to Meshery, see the [Contributing Guidelines](https://docs.meshery.io/project/contributing).