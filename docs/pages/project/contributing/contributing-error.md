---
layout: page
title: Creating Meshery Errors
permalink: project/contributing-error
description: How to declare errors in Meshery components.
language: en
type: project
category: contributing
---

Meshery pervasively uses MeshKit as a golang and service mesh management-specific library in all of its components. MeshKit helps populate error messages with a uniform and useful set of informative attributes. 
The below are the listed attributes:
- Code
- Short Description
- Long Description
- Probable Cause
- Suggested Remediation

Inorder to create a meshery error object, you will need to create a custom wrapper object for the native golang error. This can be done from the <a href="https://github.com/layer5io/meshkit/tree/master/errors">Meshkit Error</a> package. 

Use the `errors.New()` function to create a new instance of the error object and pass situation-specific attributes as function arguments. See the following example for reference.

```code
var (
    ErrExampleCode = "111"

    // Static errors
    ErrExample = errors.New(ErrExampleCode, errors.Alert, []string{"<short-description>"}, []string{"<long-description>"}, []string{"<probable-cause>"}, []string{"<suggested remediation>"})
)

// Dynamic errors
func ErrExample(err error) error {
    return errors.New(ErrExampleCode, errors.Alert, []string{"<short-description>"}, []string{"<long-description>"}, []string{"<probable-cause>"}, []string{"<suggested remediation>"})
}
```
