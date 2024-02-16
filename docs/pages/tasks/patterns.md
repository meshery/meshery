---
layout: default
title: Cloud Native Design Patterns
permalink: tasks/pattern-management
# redirect_from: tasks/pattern-management/
type: tasks
language: en
list: include
abstract: 'Meshery applies DRY principle when managing the configuration of cloud native infrastructure .'
---

{% include alert.html
    type="warning"
    title="What are Cloud Native Patterns?"
    content='This feature area is under active development. Patterns will be importable and exportable as OCI images in v0.8.0.' %}

Patterns are essentially atomic designs with one or more components made in composed into an atomic, reusable design. Patterns are a way to apply the DRY principle when managing the configuration of cloud native infrastructure.  

## Use Meshery to Deploy a Cloud Native Pattern

Use Meshery to deploy a cloud native pattern. Cloud native patterns are detailed as a YAML file. See the [Configuration Management]({{ site.baseurl }}/guides/configuration-management) guide.

### Pattern Repository

_See [Meshery Catalog](https://meshery.io/catalog)_.

You may bring your own patterns or find them available through your chosen provider. Each cloud native pattern carries a unique identifier for reference. The patterns in this repository serve in an educational capacity, facilitating learning, and also serve in an operational capacity, facilitating implementation and validation of your cloud native deployment’s adherence to a pattern.

{% include alert.html
    type="info"
    title="Repository of Cloud Native Patterns"
    content="A central set of design patterns is available in Meshery Catalog (<a href='/concepts/catalog'>concept</a>, <a href='https://meshery.io/catalog'>site</a>)." %}


## Importing Patterns

You can import a pattern using `mesheryctl` or Meshery UI. Patterns can be imported from your local filesystem, an HTTP/S endpoint, or from GitHub. When provided a GitHub location (org/repo), Meshery will recursively search the given directory (or the entire repository) for existing pattern files.

### Using Meshery CLI

#### To import a pattern, execute this command:

```
$ mesheryctl pattern import -f <path to the pattern>
Importing pattern…
Pattern successfully imported.
```

### Using Meshery UI:

From Meshery UI, patterns can be imported from your local filesystem or imported from a remote URL.

<a href="{{ site.baseurl }}/assets/img/patterns/ImportPatternUI.png">
    <img src="{{ site.baseurl }}/assets/img/patterns/ImportPatternUI.png" style="width: 60%" />
</a>

_To upload from url click the link icon_

<a href="{{ site.baseurl }}/assets/img/patterns/UrlImport.png">
    <img src="{{ site.baseurl }}/assets/img/patterns/UrlImport.png" style="width: 60%" />
</a>


#### To deploy a pattern, execute this command:

```
$ mesheryctl pattern apply -f <path to the pattern>
Deploying pattern…
Pattern successfully applied.
```

From here, output and behavior will vary based upon the specific pattern you’re deploying. Should you find that your pattern is not successfully deployed, refer to the pattern troubleshooting guide in the Meshery documentation.

As you step through each pattern, you might choose to modify the pattern to suit your specific environment and workload, exploring in-context of your specific situation. Optionally, you may choose to use a plugin like MeshMap, or another visual cloud native topology tool, to facilitate your comprehension of the patterns and to literally see the patterns in-action.

Take time to explore. Try deploying one cloud native’s sample application onto a different cloud native and compare differences in behavior and each cloud native’s ability to manage it. If using Meshery to do so, execute the following commands as an example to deploy the sample application, BookCatalog, onto Open cloud native:

```
$ mesheryctl pattern apply -f book-catalog
Deploying application “BookCatalog”...
Deployed. Endpoint(s) available at: http://localhost:8000/catalog
```

## Related Reading

- [`mesheryctl pattern`]({{ site.baseurl }}/reference/mesheryctl/pattern)
