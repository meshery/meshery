---
name: 🖥 💡 Meshery design
about: Create, export, and publish designs or models in the Meshery Catalog
title: '[Meshery Design]'
labels: 
- area/models
- hacktoberfest
assignees: ''
---
### Design description
<!-- A brief description of the design -->

### Pre-requisites
<!-- Add Pre-requisites -->
- 🛠️[Mesheryctl installation](https://docs.meshery.io/installation/mesheryctl)
- 📚[Create a Meshery Design](https://docs.meshery.io/guides/configuration-management/creating-a-meshery-design)

### Implementation
<!-- Add your design instructions -->
- Open [Meshery](https://play.meshery.io) in your browser.
- Navigate to components in the panel, where you can browse for components and define your design.


### Save and Publish 
<!-- Instructions to save and publish the design -->
- Once your design is completed, click on "Save as" button and name your design.
- Now head to [playground](https://playground.meshery.io/), in the side panel go to Configuration > Designs.
- Click on your saved design, select the "info" button, which should display the necessary fields to be filled.
- Once done with updating the fields, click on "Publish", this will send a request to the maintainers to review your design, once approved, your design will be listed at [Meshery Catalog](https://meshery.io/catalog/designs).
- Download the .yaml file for your design.

### Raise a Pull Request
- Fork the meshery/meshery repository.
- Create a new branch in your fork of the Meshery repository.
- Now head to [hacktoberfest-contributions](https://github.com/meshery/meshery/tree/master/hacktoberfest_contributions) directory
- Add your design.yaml file under the path structure:
meshery/hacktoberfest_contributions/<design-name>/<design.yaml>
- Create a pull request (PR) with your design to the Meshery repository.

---
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)

### Additonal resources

- 📚 [Components](https://docs.meshery.io/concepts/logical/components)
- 📚 [Relationships](https://docs.meshery.io/concepts/logical/relationships)
- 👨‍💻 [Models Repository](https://github.com/meshery/meshery/tree/master/models)
- 📺 [Self-paced Contributor Trainings](https://meshery.io/talks-and-trainings#trainings)
