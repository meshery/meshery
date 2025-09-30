---
name: 🖥 💡 Model Catalog
about: Create, export, and publish models in the Meshery Catalog
title: '[Model catalog]'
labels: 
- area/models
- hacktoberfest
- help wanted
assignees: ''
---
### Design description
<!-- A brief description of the design -->

### Pre-requisites
<!-- Add Pre-requisites -->
- 🛠️[Mesheryctl installation](https://docs.meshery.io/installation/mesheryctl)
- 📚[Publishing models](https://meshery.io/catalog/models#:~:text=%C3%97-,Publish%20Your%20Own%20Model,-Using%20Meshery%27s%20Registry)

### Implementation
<!-- Add your model publishing instructions -->
- Fork the meshery/meshery repository.
- Run `mesheryctl system start`
- Head to [meshmodel](https://github.com/meshery/meshery/tree/master/server/meshmodel) in your local meshery repository
- Import model definition file (JSON file) by using `mesheryctl model import [model-name/model-version]`
- Use `meshery model build [model-name/model-version]`, this build the model OCI image.

Now to publish your model to catalog:
- Fork the meshery/meshery.io repository.
- Create a new branch in your fork of the meshery/meshery.io repository.
- Add your model definition (see [template](https://github.com/meshery/meshery.io/blob/master/collections/_custom-models/_custom-model-definition-template.md) for reference) to the collections/_custom-models directory.
- Add your model oci artifact file to the assets/modelsFiles directory.
- Create a pull request to the meshery/meshery.io repository.
- Once your pull request is merged, your model will be available in the next Meshery release.

---
- 🙋🏾🙋🏼 Questions: [Discussion Forum](https://meshery.io/community#community-forums) and [Community Slack](https://slack.meshery.io)

### Additional resources

- 📚 [Models](https://docs.meshery.io/concepts/logical/models)
- 📚 [Creating models](https://docs.meshery.io/guides/configuration-management/creating-models#create-models)
- 📚 [Components](https://docs.meshery.io/concepts/logical/components)
- 📚 [Relationships](https://docs.meshery.io/concepts/logical/relationships)
- 👨‍💻 [Models Repository](https://github.com/meshery/meshery/tree/master/server/meshmodel)
- 📺 [Contributor training for Meshery Models](https://www.youtube.com/watch?v=K2gmdIlGXNo&t=82s)
