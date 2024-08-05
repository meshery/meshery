---
layout: page
title: Contributing to Meshery Server
permalink: project/contributing/contributing-Server
redirect_from: project/contributing/contributing-Server/
abstract: How to contribute to Meshery Server
language: en
type: project
category: contributing
list: include
---

As a new contributor, you’re going to want to familiarize with the project in order to resolve the issues in the best way. Installing and playing around with Meshery will give you context for any issues that you might work on.

Once an issue has been addressed, you’ll need to test it as well. Ideally, these tests are run from the user’s perspective (someone running Meshery in a container), not from a contributor’s perspective (someone running Meshery as a locally-compiled service).

## Compiling and Running Meshery Server

To build and run Meshery Server from source:

1. Build the static assets for the UI by running

{% capture code_content %}make ui-setup
make ui-build{% endcapture %}
{% include code.html code=code_content %}

2. Build & run the Server code by running

{% capture code_content %}make server{% endcapture %}
{% include code.html code=code_content %}

Any time changes are made to the Go code, you will have to stop the Server and run the above command again.
Once the Meshery Server is up and running, you should be able to access Meshery on your `localhost` on port `9081` at `http://localhost:9081`. One thing to note, you might NOT see the [Meshery UI](#contributing-ui) until the UI code is built as well.
After running Meshery Server, you will need to select your **Cloud Provider** by navigating to `localhost:9081`. Only then you will be able to use the Meshery UI on port `3000`.

**Please note**: If you get error while starting the Server as **"Meshery Development Incompatible"** then follow the below guideline 👇

<a href="{{ site.baseurl }}/assets/img/meshery-development-incompatible-error.png">
  <img style= "max-width: 450px;" src="{{ site.baseurl }}/assets/img/meshery-development-incompatible-error.png" />
</a>

Potential Solution:

- Go to your meshery folder in your local-system where you’ve cloned it.
  Execute:

- `git remote add upstream https://github.com/meshery/meshery`
- `git fetch upstream`
- Restart the meshery Server
- Additionally, before restarting the Server, if you like to pull the latest changes, you can do: `git pull upstream master`

### Build the Docker image

To build a Docker image of Meshery, please ensure you have `Docker` installed to be able to build the image. Now, run the following command to build the Docker image:

{% capture code_content %}make docker{% endcapture %}
{% include code.html code=code_content %}

### Define and validate errors

Every Golang-based component within the Meshery ecosystem incorporates a utility to define and manage error messages for every error instance. This is internally done with several make commands, but one can explicitly validate with the help of the following make command. This checks and validates the errors that are present in the particular project.

{% capture code_content %}make error{% endcapture %}
{% include code.html code=code_content %}

For more details, <a href="{{ site.baseurl }}/project/contributing/contributing-error">Error Utility</a>

### Configuring Log levels at Runtime

Meshery Server log levels can be configured at runtime by changing the env variable `LOG_LEVEL` defined in file [`Server-config.env`](https://github.com/meshery/meshery/blob/master/Server/cmd/Server-config.env). The configuration library (`viper`) watches for the env file, any change in the file content results in the `file_system` event to be emitted and the log level is updated accordingly.

**_Should there be any alterations to the location or name of the environment file, it will result in the inability to configure log levels during runtime. In the event of such modifications, it is essential to update the Server to preserve proper functionality._**

Available Meshery Server log levels are:

```
 - Panic - 0
 - Fatal - 1
 - Error - 2
 - Warn  - 3
 - Info  - 4
 - Debug - 5
 - Trace level - 6
```

The default setting for the `LOG_LEVEL` is `4` (Info). However, if the `DEBUG` environmental variable is configured as `TRUE`, it supersedes the value set in the `LOG_LEVEL` environmental variable, and the logging level is then adjusted to `5`(Debug).

### Using custom MeshKit code for Meshery Server development

<ol>
  <li>
    <p>Checkout <strong>meshery</strong> and <strong>meshkit</strong> repository in adjacent directories.</p>
    {% capture code_content %}
$ git clone https://github.com/meshery/meshery.git
$ git clone https://github.com/meshery/meshkit.git
    {% endcapture %}
    {% include code.html code=code_content %}
  </li>
  <li>
    <p>In your <code>meshery</code> go.mod, update the meshkit package.</p>
    {% capture code_content %}
github.com/layer5io/meshkit => ../meshkit
    {% endcapture %}
    {% include code.html code=code_content %}
    <p>Remember to remove this go.mod change when creating pull requests.</p>
  </li>
</ol>

## Practicing Schema-driven Development

Meshery has many logical constructs, each of which are represented by a JSON Schema in the [meshery/schemas repo](https://github.com/meshery/schemas). This section of the contributing guide provides instructions for generating Go code from these JSON Schemas / OpenAPI specifications, best practices and coding conventions while developing within any of Meshery's Golang-based components.

### Install Required Tools

1. `oapi-codegen`: https://github.com/oapi-codegen/oapi-codegen?tab=readme-ov-file#install
2. `redocly`: `npm i -g @redocly/cli`

### Prerequisite

1. Ensure you have forked the [meshery/schemas](https://github.com/meshery/schemas) and that your current working directory is the root of your `schemas` clone.
2. Resolving References in JSON Schemas:
    The OpenAPI spec references other JSON schemas, which in turn has references to other schemas, the `oapi-codegen` tool fails to resolve such nested references, hence the references needs to be resolved before generating the code.
    Use the `ref-resolver.js` script provided in the `schemas` (scripts/ref-resover.js) repository to perform the resolution:

_If the openapi spec points to a ref, eg: ../a.json and the schema inside a.json doesn't have any other references then the `oapi-codegen` works fine, the task of the `ref-resovler.js` is to just ensure the nested external references are resolved to the point `oapi-codegen` can work._

`Run: make resolve-ref path=<file/directory containing the schemas whose references needs to be resolved (the path should be relative from the root of the repository)> eg: make resolve-ref path=schemas/constructs/v1beta1.`

### Generating constructs

```bash
oapi-codegen -config config.yml schemas/constructs/openapi/patterns.yml
```

In this command:

- `-config config.yml`: Specifies the configuration file.
- `schemas/constructs/openapi/patterns.yml`: Location of the OpenAPI specification file.

The tool supports passing the configuration via flags or via `config.yaml`. The supported options are described [here](https://github.com/oapi-codegen/oapi-codegen/blob/main/
configuration-schema.json).
Refer [scripts/config.yml](github.com/meshery/schemas) , for mostly used options:

Below is an example for `config.yml` file.

```yaml
package: v1beta1 # package name under which the genreate code should reside.

generate:
  models: true # only generate models i.e. golang structs, skips the genenration Server code i.e. API Handlers, middleware etc.

# A schema can refer to other schema/sub-schema, and the golang code might have already been generated, so to ensure code is not repeated but uses the already exisitng package, define the pacakge name and location.
# For eg: For most of the commonly used data-types (uuid, email, URL) the schema is defined in the `core.yml` file (inside schemas repo) and the corresponding golang code is already generated in the `meshery/schemas/models/core` package. Hence at the time of generation, in order to prevent code duplication, we need to map the schema to the already generated package. "./common/core.yml" it informs the tool to replace any usages with the existing golang package, also adding necessary imports in the code to be generated.

_NOTE: Import mappings cannot be used for references within the same file, for eg: "#/components/schemas/dummy-schema": "github.com/../../dummay-schema", such mappings are not handled by the tool._

import-mapping:
  "./core.yml": ""
  "./common/core.yml": "github.com/meshery/schemas/models/core"
  "../core.json": "github.com/meshery/schemas/models/core"
  "./common/parameters.yml": "github.com/meshery/schemas/models/core"
  "../v1beta1/model.json": "github.com/meshery/schemas/models/v1beta1"
  "../v1beta1/component.json": "github.com/meshery/schemas/models/v1beta1"

# Location of the generated code, this can be overriden by passing the `-o` flag in the CLI command.
output: models/v1beta1/patterns.go 

output-options:
  # name-normalizer: ["ToCamelCase" | "ToCamelCaseWithDigits" | "ToCamelCaseWithInitialisms" ] 
  skip-prune: false # Remove structs/model which are not referenced from any part of the generated code.
  
  # An OpenAPI spec may contain schema for many resources, and grouped together using tags, to generate code only for schemas with specific tags, use the `include-tags` flag, `exclude-tags` to excude.
  include-tags: 
  - patterns
```

## Generating REST API docs
Run `make docs-build` to generate statid REST API docs.
Run `make preview-docs` to see live changes for the REST API docs.

## Best Practices for Schema Design

### 1. Specifying Additional Tags

To include additional tags in the generated Go structs, use the `x-oapi-codegen-extra-tags` extension in your OpenAPI specification. For example, to add a `yaml` or `gorm` tag to a struct field, update your schema as follows:

The tool adds the json tags automatically.

```yaml
pattern:
  properties:
    Id:
      $ref: '#/components/id'
      x-oapi-codegen-extra-tags:
        gorm: id
```

### 2. Casing Conventions

Ensure that the casing/naming schema of your Go structs matches your requirements. 
Use `x-go-name`

```yaml
properties:
  name:
    type: string
    x-go-name: PatternName
```

```gotype Pattern struct {
   // Unique identifier
   PatternName string `json:"name,omitempty"`
}
```
### 3. Skip pointer type for optional fields in structs	

Add `x-go-type-skip-optional-pointer: true` tag in the schema to prevent pointer fields from being generated for fields which are not marked as required in the schema.


### 4. Overriding/Specifying data type and package to use.

To use specific imports for certain fields, use `x-go-type` and `x-go-type-import`. This helps in managing dependencies and ensures that your generated code uses the correct imports.

Specify the Go types to be used for defined attributes in your OpenAPI schema. For example:

```yaml
components:
  schemas:
    SomeSchema:
      properties:
        id:
          type: string
          format: uuid
            x-go-type: gofrsuuid.UUID
            x-go-type-import:
              path: github.com/gofrs/uuid
              name: gofrsuuid
```

## Further Reading

- [oapi-codegen Documentation](https://github.com/deepmap/oapi-codegen)
- [Meshery Schemas](https://github.com/meshery/schemas)
- [Reference Resolver Script](https://github.com/meshery/schemas/blob/cba1c8d13a4eb801b4c0642cf7af4fc4c1fe617c/scripts/ref-resolver.js)

Thank you for contributing to the project! For any questions or issues, please open an issue or pull request on the respective GitHub repository.