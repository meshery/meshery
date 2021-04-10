# Contributor Guide for UI component

This guide is specific to the Meshery UI component and involves steps/methods one need to follow while working on issues related to Meshery UI.

## How to run Meshery UI?
Meshery UI can be used for different purposes, it can be for general use like testing out Meshery and it's functionalities or work on issues related to the Meshery UI.

#### 1. General Usage:
For general usage, one can run Meshery UI using Meshery's command client `mesheryctl`, by simply running the `mesheryctl system start` command.
If you don't have the `mesheryctl` tool installed already, you can follow the [mesheryctl installation docs](https://docs.meshery.io/installation/mesheryctl) to install `mesheryctl` using various `package management` tools supported.

#### 2. For development purporse:
For development purpose, either of the below mentioned approach could be followed
- Follow the procedure mentioned in Step 1 (General Usage) above, and start Meshery UI sever on the 9081 port, and login to Meshery UI using eihter of the providers mentioned on the login page. Then, to run a development server of Meshery UI, first install the dependencies using the `make setup-ui-libs` command, then build the UI package and export it using the command `make build-ui`, and finally execute `make run-ui-dev` to run the livereload-nodemon server on port 3000.
    > **NOTE:** Please run the steps in order to avoid issues, as Meshery server should be running and logged-in before accessing the development server
    > on 3000 port.
    
- Another way, to run Meshery UI for development purpose is by using the `make run-local` command. But, before that you will have to build the Meshery UI and export it using `make buil-ui` command. This method doesn't give a livereload server and one will have to build the Meshery-UI after making changes to the code and then again run the server.
    > **NOTE:** If you are using this method, make sure you don't have Meshery already running on 9081 port, using `mesheryctl`.

## Tech stack used in Meshery-UI
- Meshery UI uses NextJs to do server side rendering of ReactJS components. The folder `ui/components` contains all the ReactJS components involved in
  building Meshery UI.
- MaterialUI is being used extensively for the visual components of Meshery UI.
- Billboard.js library is being used to display various charts, and comparison graphs in Meshery UI.


## Component naming convention
For reference and easy code search, the components are named accordingly following the rule 'Meshery<Part of UI it involves>', for example: components
involved in rednering the Results page of Meshery UI are named as 'MesheryResults.js', 'MesheryResultDialog.js', 'MesherySMIResults.js'. Please follow this convention if you are creating a new component.


<p style="text-align: center"><em>If you'll like to go to the main Meshery Contributor guide <a href="../CONTRIBUTING.md">click here</a></em></p>