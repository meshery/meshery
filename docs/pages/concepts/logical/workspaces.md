---
layout: enhanced
title: Workspaces
permalink: concepts/logical/workspaces
type: concepts
abstract: "Meshery Workspaces act as central collaboration point for teams."
language: en
list: include
---

Meshery Workspaces serve as a virtual space for your team-based work. Create a Workspace to organize your work and to serve as the central point of collaboration for you and your teams and a central point of access control to Environments and their resources.

You may create Workspaces to organize project-based work or to create domains of responsibility for your teams or segregate Designs and Environments and track team activity.

## Summary

Workspaces facilitate collaboration between you and your teams, allow you to control access to resources, and track activity and report on related events.

## Key Features

- **Resource Sharing** Workspaces allow for seamless resource sharing among team members, fostering collaboration.
- **Logical Grouping** Within Workspaces, you can group related components such as environments and infrastructure designs.
- **Flexibility**: Workspaces support various resources like Kubernetes, Prometheus, Jaeger, Nginx, and more.
- **Simplified Management** Managing and deploying resources is made easy within Workspaces.
- **Access Control** Workspaces allow you to control access to resources by granting permissions to users and teams.

### Workspace Relationships and Restrictions

- Access to Workspaces may be granted to one or more teams.
- As a point of collaboration to facilitate work, Workspaces may have zero or more Environments associated.

After creating a Workspace, of your next steps is to resource that Workspace. Like a shared drive (or or shared collection of files). Workspaces are your Google Drive, while Meshery Designs are your Google Docs.

## Key Components

### Environments

- Environments are a central part of a workspace. They serve as a logical grouping for managing connections. A connection, in this context, can be either managed or discovered by Meshery. Examples of connections include Kubernetes clusters, Prometheus instances, Jaeger traces, Nginx servers, and more.
- One or more environments can be assigned to a workspace.
- Same environment can be assigned to multiple workspaces.

See "[Environments](/concepts/logical/environments)" section for more information.

### Designs

- Infrastructure Designs are essential for creating reusable deployment templates. Users belonging to teams with access to a workspace can utilize these designs to deploy resources in the Kubernetes clusters associated with that workspace.
- Like a shared drive (or or shared collection of files), Workspaces are your Google Drive, while Meshery Designs are your Google Docs.
- One ore more designs can be assigned to a workspace.
- Same design can be assigned to multiple workspaces.

See "[Meshery Designs](/concepts/logical/designs)" section for more information.

### Organizations

- Organizationa are the unit of tenancy in Meshery. Organizations group users together.
- Organizations own all resources created by users, like Workspaces, Designs, Environments, and so on.
- Remote Providers can extend Meshery to include additional identity and user management features like heirarchical organizations, or teams (as user groups), and so on.
- Remote Providers can extend Meshery to offer fine-grained permissions and access control to resources like Workspaces, Designs, Environments, and so on.

Learn more about [extensible authorization](/extensibility/authorization).
 <!-- "[Organizations](/extensions/team-management)" section for more information. -->

<!-- ### Teams

- A Workspace is closely associated with Teams in Meshery. Teams are groups of users with varying permissions, and they are at the center of resource access and management within a Workspace.
- One ore more teams can be assigned to a workspace.
- Same team can be assigned to multiple workspaces.

See "[Teams](/extensions/team-management)" section for more information. -->

### Connections

- Connections in Meshery refer to various resources that can be either managed or unmanaged, but are discovered and made accessible for use. Examples of Connections include Kubernetes clusters, Prometheus instances, Jaeger services, Nginx deployments, and more. These Connections serve as a fundamental part of Workspaces, as users can deploy infrastructure designs within the context of these Connections.
- Connections can be assigned to one or more environments.
- Same connection can be assigned to multiple environments.

See "[Connections](/concepts/logical/connections)" section for more information.

## Best Practices

To make the most of Meshery Workspaces, here are some best practices:

- Clearly define permissions in the form of team assignement to ensure proper access control.
- Use Infrastructure Designs to standardize resource deployments.
- Regularly review and update your Workspace's resources and configurations.

Meshery Workspaces enhance collaboration within your teams, providing a structured environment for sharing and managing resources. By following best practices and understanding the core components of Workspaces, you can maximize the benefits of this feature in Meshery.

