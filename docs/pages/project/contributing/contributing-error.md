---
layout: page
title: Creating Meshery Errors
permalink: project/contributing-error
description: How to declare errors in Meshery components.
language: en
type: project
category: contributing
---

Meshery pervasively uses MeshKit as a golang and service mesh management-specific library in all of its components. MeshKit helps populate error messages with a uniform and useful set of informative attributes:

Meshery pervasively uses MeshKit as a golang and service mesh management-specific library in all of its components. MeshKit helps populate error messages with a uniform and useful set of informative attributes. 

This tool analyzes, verifies and updates error codes in Meshery source code trees. It extracts error details into a file that can be used for publishing all error code references on the Meshery website. The objective to create this was to avoid centralized handling of error codes, automating everything

Inorder to create a meshery error object, you will need to create a custom wrapper object for the native golang error. This can be done from the <a href="https://github.com/layer5io/meshkit/tree/master/errors">Meshkit Error</a> package. 

Use the `errors.New()` function to create a new instance of the error object and pass situation-specific attributes as function arguments. 
The below are the listed attributes:
- Code
- Short Description
- Long Description
- Probable Cause
- Suggested Remediation

### Syntax
     errors.New(ErrExampleCode, errors.Alert, []string{"<short-description>"}, []string{"<long-description>"}, []string{"<probable-cause>"}, []string{"<suggested remediation>"})


It is intended to be run locally and as part of a CI workflow.

- Errors names and codes are namespaced to components, i.e. they need to be unique within a component, which is verified by this tool.
- A component corresponds usually to a repository. Components have a type and a name. 
  They are also returned from the ComponentInfo endpoint, e.g. for adapters.
  Examples of a component types are 'adapter' and 'library', corresponding examples of names are 'istio' and 'meshkit'.
- There are no predefined error code ranges for components.
  Every component is free to use its own range, but it looks like the convention is to start at 1000.
- Errors are not to be reused across components and modules.
- Codes carry no meaning, as e.g. HTTP status codes do.
- In the code, create string var's or const's with names starting with Err[A-Z] and ending in Code, e.g. 'ErrApplyManifestCode'.
- Set the value to any string, like "replace_me" (no convention here), e.g. ErrApplyManifestCode = "test_code".
- If the value is a string, this tool will replace it with the next integer.
- If the value is an int, e.g. ErrGetName = "1000" the tool will not replace it unless it is forced (command line flag --force).
  If forced, all codes are renumbered. This can be useful to tidy up in earlier implementations of meshkit error codes.
- Setting an error code to a call expression like ErrNoneDatabase = errors.NewDefault(ErrNoneDatabaseCode, "No Database selected")
  is not allowed. This tool emits a warning if a call expression is detected.
- Using errors.NewDefault(...) is deprecated. This tool emits a warning if this is detected.
- Use errors.New(...) from meshkit to create actual errors with all the details.
  This is often done in a factory function. It is important that the error code variable is used here, not a literal.
  Specify detailed descriptions, probable causes, and remedies. They need to be string literals, call expressions are ignored.
  This tool extracts this information from the code and exports it.
- By convention, error codes and the factory functions live in files called error.go. The tool checks all files, but updates only error.go files.
- This tool will create a couple of files, one of them is designed to be used to generate the error reference on the meshery website.
  The file errorutil_analyze_summary.json contains a summary of the analysis, notably lists of duplicates etc.
- The tool requires a file called component_info.json. Its location can be customized, by default it is the root directory (-d flag). 
  This file has the following content, with concrete values specific for each component:
  {
    "name": "meshkit",
    "type": "library",
    "next_error_code": 11010
  }
- The tool updates next_error_code.


## Example

In this example we are trying to create an Error for not able to marhsal a JSON

```code
var (
    // Error code 
    ErrMarshalCode= "replace_me"

    //Static errors (for example)
    ErrExample = errors.New(ErrExampleCode, errors.Alert, []string{"<short-description>"}, []string{"<long-description>"}, []string{"<probable-cause>"}, []string{"<suggested remediation>"})
)

// Dynamic errors
//Error Name
func ErrMarshal(err error, obj string) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Unable to marshal the : ", obj}, []string{err.Error()}, []string{}, []string{})
}

```
### Replacing old Error Codes 

 Old
 ```Code 
    bd, err := json.Marshal(providers)
	if err != nil {
		http.Error(w, "unable to marshal the providers", http.StatusInternalServerError)
		return
	}
 ```
New
  ```Code 
    bd, err := json.Marshal(providers)
    if err != nil {
            obj := "provider"
            http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
            return
        }
 ```

 
## Replacing logrus 
 There already exists an interface for logger in <a href="https://github.com/layer5io/meshkit/blob/master/logger/logger.go">Meshkit </a>  

#### Defining a Logger 

   ```Code 
    type Logger struct {
        log   logger.Handler
    }
 ```

#### Debug
 
 Old

    `logrus.Debugf("meshLocationURL: %s", meshLocationURL)`
 New

    `l.log.Debug("meshLocationURL: ", meshLocationURL)`

#### Error

  Old

    `logrus.Errorf("error marshaling data: %v.", err)`
  New
  
    `l.log.Error(ErrMarshal(err, obj))`


# Suggested Reading

{% assign sorted_reading = site.pages | sort: page.title | reverse %}

<ul>
  {% for item in sorted_reading %}
  {% if item.type=="project" and item.category=="contributing" and item.list!="exclude" -%}
    <li><a href="{{ site.baseurl }}{{ item.url }}">{{ item.title }}</a>
    </li>
    {% endif %}
  {% endfor %}
</ul>