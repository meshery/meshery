---
name: 🖥 💡 Models Update
about: An update, addition, or removal of one or more of the models, components, relationships, workflows, or policies within Meshery Models.
title: '[Models] '
labels: area/models
assignees: ''
---
### Current Situation
<!-- A brief description of the current state of Models -->

### Proposed Change
<!-- A brief description of the change. -->

---

### General Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://layer5.io/community/handbook)

- 🙋🏾🙋🏼 Questions: [Discussion Forum](http://discuss.meshery.io) and [Community Slack](https://slack.meshery.io)

### Understanding Meshery Models

- 🖥 [Models Overview](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g226f5de5180_19_259)
- 👨‍💻[Models Repository](https://github.com/meshery/meshery/tree/master/server/meshmodel)

### Contributing to Meshery Models

- https://docs.meshery.io/project/contributing/contributing-models
- https://docs.meshery.io/project/contributing/contributing-relationships

### Instructions for Models

To add or update a model, follow these steps:
1. Create a Model Definition (yaml).

### Instructions for Components

While the default shape for new components is a circle, each component should be considered for its best-fit shape.
1. Review and familiarize with the available set of predefined relationship types. Refer the Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for a list of possible shapes. 
1. Propose a specific shape, best-suited to visually represent the Component. _Example - Deployment as a pentagon._
1. Proposee a specific icon, best-suited to visually represent the Component. _Example - DaemonSet as a skull icon._

### Instructions for Relationships

1. Identify the relationship and any specific constraints to be enforced between one or more specific components within the same or different models.
1. Propose a specific visual representation for the relationship. See list of visualizations on [Visualizing Relationships](https://docs.meshery.io/project/contributing/contributing-relationships#relationship-visualizations)
1. Prospose the appropriate relationship type, using one of the predefined set of relationship types or suggest a new relationship where an existing type does not fit.
1. Create a Relationship Definition (yaml).
1. (Typically not necessary) create a policy for evaluation of the relationship (rego).
1. Review a prior pull request as an example of how to define a Relationships: https://github.com/meshery/meshery/pull/9880/files

<!-- ### Instructions for Policies
1. _Forthcoming_
 -->
