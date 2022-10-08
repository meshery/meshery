# Contributor Guide for UI component

This guide is specific to the Meshery UI component and involves steps/methods one need to follow while working on issues related to Meshery UI.

## How to run Meshery UI?
Meshery UI can be built and run in different ways. You will choose one of the two ways to build and run Meshery UI depending upon whether you are actively developing it (whether you are creating a new feature or fixing a bug in Meshery UI) or whether you simply need to use it as a user. Let's refer to these two methods as a _Development Build_ and _User Build._

#### 1. User Build:
For general usage, one can run Meshery UI using Meshery's command client `mesheryctl`, by simply running the `mesheryctl system start` command.
If you don't have the `mesheryctl` tool installed already, you can follow the [mesheryctl installation docs](https://docs.meshery.io/installation/mesheryctl) to install `mesheryctl` using various `package management` tools supported.

#### 2. Development Build:
For purposes of actively developing Meshery UI, you first need to ensure you have npm v7 installed (`npm -v`) and if not install it (`npm -g i npm@7`), then install the dependencies using `make setup-ui-libs` and then you can use either of the following approaches to build Meshery UI:
1. Follow the procedure mentioned in Step 1 (User build) above, and start Meshery UI sever on the 9081 port, and login to Meshery UI using either of the providers mentioned on the login page. Then, to run a development server of Meshery UI, install the dependencies using the command mentioned above, then execute `make ui` to run the livereload-nodemon server on port 3000.
    > **NOTE:** Please run the steps in order to avoid issues, as Meshery server should be running and logged-in before accessing the development server
    > on 3000 port.

    > **NOTE:** Its strongly recommended to use  either [Node Version Manager](https://github.com/nvm-sh/nvm#node-version-manager---)) in linux/mac os systems or [NVM for Windows](https://github.com/coreybutler/nvm-windows#nvm-for-windows) on Windows systems so single `nvm use` / `nvm install` simplifies installing and using correct node version locallly **(v14.8.0)**, see [NVM Intro](https://github.com/nvm-sh/nvm#intro) for details. Otherwise, you might experience issues during local `npm i` similar to [4674](https://github.com/meshery/meshery/issues/4674) due to how optional dependencies are resolved in npm v6.
    
1. **`make server`** - Alternatively, build all of Meshery UI's components upfront before serving the UI. Do this in two steps:
 - Execute `make ui-setup` to iniatilize your environment and then `make ui-build` to build and export all Meshery UI components.
 - Execute `make server` to serve the prebuilt components. 
This method doesn't provide a live reload server. You will have to build Meshery UI after making changes to the code and rerun these steps again in order to see those subsequent code changes reflected in the UI.
    > **NOTE:** If you are using this method, make sure you don't have Meshery already running on 9081 port, using `mesheryctl`.

## Tech stack used in Meshery-UI
- Meshery UI uses NextJs to do server side rendering of ReactJS components. The folder `ui/components` contains all the ReactJS components involved in
  building Meshery UI.
- MaterialUI is being used extensively for the visual components of Meshery UI.
- Billboard.js library is being used to display various charts, and comparison graphs in Meshery UI.


## Component naming convention
For reference and easy code search, the components are named accordingly following the rule 'Meshery<Part of UI it involves>', for example: components
involved in rendering the Results page of Meshery UI are named as 'MesheryResults.js', 'MesheryResultDialog.js', 'MesherySMIResults.js'. Please follow this convention if you are creating a new component.


<p style="text-align: center"><em>If you'll like to go to the main Meshery Contributor guide <a href="../CONTRIBUTING.md">click here</a></em></p>
