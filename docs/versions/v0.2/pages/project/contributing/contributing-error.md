---
layout: page
title: How to write MeshKit compatible errors
permalink: /v0.2/project/contributing/contributing-error
description: How to declare errors in Meshery components.
language: en
type: project
category: contributing
---

Meshery pervasively uses MeshKit as a golang and service mesh management-specific library in all of its components. MeshKit helps populate error messages with a uniform and useful set of informative attributes.

To help with creating error codes, MeshKit contains a tool that analyzes, verifies and updates error codes in Meshery source code trees. It extracts error details into a file that can be used for publishing all error code references on the Meshery [error codes reference page](https://docs.meshery.io/reference/error-codes). The objective to create this was to avoid centralized handling of error codes and automating everything

In order to create a Meshery error object, you will need to create a custom wrapper object for the native golang error. This can be done from the <a href="https://github.com/layer5io/meshkit/tree/master/errors">MeshKit Error</a> package.

## Some rules for making Errors codes

- Errors names and codes are namespaced to components, i.e. they need to be unique within a component, which is verified by this tool.

- Errors are not to be reused across components and modules.

- Codes carry no meaning, as e.g. HTTP status codes do.

- For the Code argument in the errors.New use the same Error name and append a "Code" after it. e.g error name : ErrApplyManifest then the error code is ErrApplyManifestCode

- Set the value to any string, like "replace_me" (no convention here), e.g. ErrApplyManifestCode = "replace_me".

- Error codes are not to be set as integer

- CI will take care of updating Error codes from a string to an integer.

- Using errors.NewDefault(...) is deprecated. This tool emits a warning if this is detected.

- Use errors.New(...) from MeshKit to create actual errors with all the details.
  This is often done in a factory function. It is important that the error code variable is used here, not a literal.
  Specify detailed descriptions, probable causes, and remedies. They need to be string literals, call expressions are ignored.
  This tool extracts this information from the code and exports it.

- First letter of the every error description need to capitalized

- By convention, error codes and the factory functions live in files called error.go. The tool checks all files, but updates only error.go files.

- This tool will create a couple of files, one of them is designed to be used to generate the error reference on the meshery website.
  The file errorutil_analyze_summary.json contains a summary of the analysis, notably lists of duplicates etc.

- Running `make error` would analyze the code and return you with a warning.

Use the `errors.New()` function to create a new instance of the error object and pass situation-specific attributes as function arguments.
These attributes are:

- Code
- Short Description
- Long Description
- Probable Cause
- Suggested Remediation

### Syntax

     errors.New(ErrExampleCode, errors.Alert, []string{"<short-description>"}, []string{"<long-description>"}, []string{"<probable-cause>"}, []string{"<suggested remediation>"})

## Example

In this example we are creating an Error for being unable to marshal JSON

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

There already exists an [interface for logger](https://github.com/layer5io/meshkit/blob/master/logger/logger.go) in MeshKit.

#### Defining a Logger

```Code
 type Logger struct {
     log   logger.Handler
 }
```

#### Debug

##### Old

`logrus.Debugf("meshLocationURL: %s", meshLocationURL)`

##### New

`l.log.Debug("meshLocationURL: ", meshLocationURL)`

#### Error

##### Old

`logrus.Errorf("error marshaling data: %v.", err)`

##### New

`l.log.Error(ErrMarshal(err, obj))`

{% include suggested-reading.html %}
