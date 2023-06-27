---
layout: page
title: Contributing to Meshery Server
permalink: project/contributing/contributing-server
redirect_from: project/contributing/contributing-server/
description: How to contribute to Meshery Server
language: en
type: project
category: contributing
---

As a new contributor, you’re going to want to familiarize with the project in order to resolve the issues in the best way. Installing and playing around with Meshery will give you context for the any issues that you might work on.

Once an issue has been addressed, you’ll need to test it as well. Ideally, these tests are run from the user’s perspective (someone running Meshery in a container), not from a contributor’s perspective (someone running Meshery as a locally-compiled service).

## Compiling and Running Meshery server

To build and run Meshery server from source:

1. Build the static assets for the UI by running

{% capture code_content %}make ui-setup
make ui-build{% endcapture %}
{% include code.html code=code_content %}


2. Build & run the server code by running

{% capture code_content %}make server{% endcapture %}
{% include code.html code=code_content %}

Any time changes are made to the GO code, you will have to stop the server and run the above command again.
Once the Meshery server is up and running, you should be able to access Meshery on your `localhost` on port `9081` at `http://localhost:9081`. One thing to note, you might NOT see the [Meshery UI](#contributing-ui) until the UI code is built as well.
After running Meshery server, you will need to select your **Cloud Provider** by navigating to `localhost:9081`. Only then you will be able to use the Meshery UI on port `3000`.

**Please note**: If you get error while starting the server as **"Meshery Development Incompatible"** then follow the below guideline 👇

<div style="text-align: center;">
<a href="{{ site.baseurl }}/assets/img/meshery-development-incompatible-error.png">
  <img style= "width: 600px;" src="{{ site.baseurl }}/assets/img/meshery-development-incompatible-error.png" />
</a>
</div>

Potential Solution: 

-  Go to your meshery folder in your local-system where you’ve cloned it.
Execute: 


- Restart the meshery server.
- Make sure that `BUILD="v0.X.X"` is logged in some starting line while starting the server.

- Addionally, before restarting the server, if you like to pull the latest changes, you can do: `git pull upstream master`

#### Building Docker image

To build a Docker image of Meshery, please ensure you have `Docker` installed to be able to build the image. Now, run the following command to build the Docker image:

{% capture code_content %}make docker{% endcapture %}
{% include code.html code=code_content %}

#### Define and validate errors

Every Golang-based project that belongs to Layer5 incorporates a utility to define and manage error messages for every error instance. This is internally done with several make commands, but one can explicitely validate with the help of the following make command. This checks and validates the errors that are present in the particular project.

{% capture code_content %}make error{% endcapture %}
{% include code.html code=code_content %}

For more details, <a href="{{ site.baseurl }}/project/contributing/contributing-error">Error Utility</a>

{% include suggested-reading.html %}
