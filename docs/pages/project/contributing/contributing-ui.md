---
layout: page
title: Contributing to Meshery UI
permalink: project/contributing-ui
description: How to contribute to Meshery UI (web-based user interface).
language: en
type: project
---

### <a name="contributing-ui">UI Contribution Flow</a>
Meshery is written in `Go` (Golang) and leverages Go Modules. UI is built on React and Next.js. To make building and packaging easier a `Makefile` is included in the main repository folder.

#### Install UI dependencies
To install/update the UI dependencies:
```
make setup-ui-libs
```

#### Build and export UI
To build and export the UI code:
```
make build-ui
```

Now that the UI code is built, Meshery UI will be available at `http://localhost:9081`.
Any time changes are made to the UI code, the above code will have to run to rebuild the UI.

#### UI Development Server
If you want to work on the UI, it will be a good idea to use the included UI development server. You can run the UI development server by running the following command:
```
make run-ui-dev
```

Make sure to have Meshery server configured, up and running on the default port `http://localhost:9081` before proceeding to access and work on the UI server at `http://localhost:3000`.
Any UI changes made now will automatically be recompiled and served in the browser.

#### Running Meshery from IDE
If you want to run Meshery from IDE like Goland, VSCode. set below environment variable
```
PROVIDER_BASE_URLS="https://meshery.layer5.io"
PORT=9081
DEBUG=true
ADAPTER_URLS=mesherylocal.layer5.io:10000 mesherylocal.layer5.io:10001 mesherylocal.layer5.io:10002 mesherylocal.layer5.io:10003 mesherylocal.layer5.io:10004 mesherylocal.layer5.io:10005 mesherylocal.layer5.io:10006 mesherylocal.layer5.io:10007 mesherylocal.layer5.io:10008 mesherylocal.layer5.io:10009
```
go tool argument
```shell
-tags draft
```
update /etc/hosts
```shell
127.0.0.1 mesherylocal.layer5.io
```
