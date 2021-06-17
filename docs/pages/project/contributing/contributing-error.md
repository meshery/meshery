---
layout: page
title: Using Meshery Error Utility
permalink: project/contributing-error
description: How to use error utility in meshery projects.
language: en
type: project
category: contributing
---

Meshery Error utility help populate error messages with certain set of meaningful attributes. 
The below are the listed attributes:
- Code
- Short Description
- Long Description
- Probable Cause
- Suggested Remediation

Inorder to create a meshery error object, you will need to create a custom wrapper object for the native golang error. This can be done from the <a href="https://github.com/layer5io/meshkit/tree/master/errors">Meshkit Error</a> package. 

Use `errors.New()` the above function to create a new instance of an error object and pass in the above listed attributes as function arguments. Refer below for an example.

```code
errors.New("111", errors.Alert, []string{"<short-description>"}, []string{"<long-description>"}, []string{"<probable-cause>"}, []string{"<suggested remediation>"})
```
