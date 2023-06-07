---
name: 🖥 💡 MeshModal Update
about: An update, addition, or removal of one or more of the model, component, relationship, policy within MeshModel
title: '[MeshModel]'
labels: component/meshmodel
assignees: ''
---

#### Issue Type
- [ ] Update
- [ ] Addition
- [ ] Removal

#### Desired Behavior
<!-- A brief description of the enhancement. -->

#### Implementation
<!-- Specifics on the approach to fulfilling the feature request. -->

#### Instructions for Defining a new Relationship
1. Identify the relationship and any specific constraints to be enforced between the two specific components, their models or potentially other components, models, or environmental considerations.
2. Propose a specific visual representation for the relationship. Review and familiarize with the available set of predefined relationship types. Refer the cytoscape [style guide](https://js.cytoscape.org/#style). Virtual representation for existing context-aware designs. 
i. [Hierarchical](https://github.com/meshery/meshery/blob/master/.github/assets/images/hierarchical_relationship.png)
ii. [Sibling](https://github.com/meshery/meshery/blob/master/.github/assets/images/sibling_relationship.png)
iii. [Binding](https://github.com/meshery/meshery/blob/master/.github/assets/images/binding_realtionship.png)
3. Prospose the appropriate relationship type, using one of the predefined set of relationship types or suggest a new relationship where an existing type does not fit.
4. Create a Relationship Definition (yaml). See the [Relationship Schema](https://github.com/meshery/meshery/tree/master/server/meshmodel/schemas) in MeshModel and [examples]([url](https://github.com/meshery/meshery/tree/master/server/meshmodel/relationships)).
5. Create a policy for evaluation of the relationship (rego). See examples.

---

#### Contributor [Guides](https://docs.meshery.io/project/contributing) and [Handbook](https://layer5.io/community/handbook)

- [MeshModel Repository](https://github.com/meshery/meshery/tree/master/server/meshmodel)
- 🖥 [MeshModel Overview](#)
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://discuss.layer5.io) and [Community Slack](http://slack.layer5.io)
