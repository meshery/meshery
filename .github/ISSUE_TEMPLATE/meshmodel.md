---
name: 🖥 💡 MeshModel Update
about: An update, addition, or removal of one or more of the models, components, relationships, or policies within MeshModel.
title: '[MeshModel]: '
labels: area/meshmodel
assignees: ''
---
### Current Situation
<!-- A brief description of the current state of MeshModel -->

### Proposed Change
<!-- A brief description of the change. -->

---

### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://layer5.io/community/handbook)

- [MeshModel Repository](https://github.com/meshery/meshery/tree/master/server/meshmodel)
- 🖥 [MeshModel Overview](https://docs.google.com/presentation/d/1SQMfyu5shjpGKlYONdVzOtd7UYTgLWBcgUvHMLCZ2tY/edit#slide=id.g226f5de5180_19_259)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://discuss.layer5.io) and [Community Slack](http://slack.layer5.io)

## How to make updates to MeshModel

### Instructions for Models
1. _Forthcoming_

### Instructions for Components
While the default shape for new components is a circle, each component should be considered for its best-fit shape.
1. Review and familiarize with the available set of predefined relationship types. Refer the Cytoscape [node types](https://js.cytoscape.org/demos/node-types/) for a list of possible shapes. 
1. Propose a specific shape, best-suited to visually represent the Component. Example - Deployment as a pentagon.
1. Proposee a specific icon, best-suited to visually represent the Component. Example - DaemonSet as a skull icon.

### Instructions for Relationships
1. Identify the relationship and any specific constraints to be enforced between the two specific components, their models or potentially other components, models, or environmental considerations.
1. Propose a specific visual representation for the relationship. Visual representation examples:
    - [Hierarchical](https://github.com/meshery/meshery/blob/master/.github/assets/images/hierarchical_relationship.png)
    - [Sibling](https://github.com/meshery/meshery/blob/master/.github/assets/images/sibling_relationship.png)
    - [Binding](https://github.com/meshery/meshery/blob/master/.github/assets/images/binding_relationship.png)
1. Prospose the appropriate relationship type, using one of the predefined set of relationship types or suggest a new relationship where an existing type does not fit.
1. Create a Relationship Definition (yaml). See the [Relationship Schema](https://github.com/meshery/meshery/tree/master/server/meshmodel/schemas) in MeshModel and [examples]([url](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships)).
1. Create a policy for evaluation of the relationship (rego). See examples.

### Instructions for Policies
1. _Forthcoming_

