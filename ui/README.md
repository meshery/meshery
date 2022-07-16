# Contributor Guide for Meshery UI

This guide is specific to the Meshery UI component and involves steps/methods one need to follow while working on issues related to Meshery UI.

## How to run Meshery UI?
Meshery UI can be built and run in different ways. You will choose one of the two ways to build and run Meshery UI depending upon whether you are actively developing it (whether you are creating a new feature or fixing a bug in Meshery UI) or whether you simply need to use it as a user. Let's refer to these two methods as a _Development Build_ and _User Build._

#### 1. User Build:
For general usage, one can run Meshery UI using Meshery's command client `mesheryctl`, by simply running the `mesheryctl system start` command.
If you don't have the `mesheryctl` tool installed already, you can follow the [mesheryctl installation docs](https://docs.meshery.io/installation/mesheryctl) to install `mesheryctl` using various `package management` tools supported.

#### 2. Development Build:
For purposes of actively developing Meshery UI, you first need to install the dependencies using `make ui-setup` and then you can use either of the following approaches to build Meshery UI:
1. Follow the procedure mentioned in Step 1 (User build) above, and start Meshery UI sever on the 9081 port, and login to Meshery UI using either of the providers mentioned on the login page. Then, to run a development server of Meshery UI, install the dependencies using the command mentioned above, and execute `make ui` to run the livereload-nodemon server on port 3000.
    > **NOTE:** Please run the steps in order to avoid issues, as Meshery server should be running and logged-in before accessing the development server
    > on 3000 port.
    
1. **`make run-local`** - Alternatively, build all of Meshery UI's components upfront before serving the UI. Do this in two steps:
 - Execute `make ui-build` to build and export all Meshery UI components.
 - Execute `make run-local` to serve the prebuilt components. 
This method doesn't provide a live reload server. You will have to build Meshery UI after making changes to the code and rerun these steps again in order to see those subsequent code changes reflected in the UI.
    > **NOTE:** If you are using this method, make sure you don't have Meshery already running on 9081 port, using `mesheryctl`.

## Development workflow

- Make sure you have a running instance of Meshery UI by following either one of the options provided above
- Go to `/ui`
- Make sure all the tests pass by running `npm run test`
- Create a new branch
- Make the changes
- Write tests if necessary
- Make sure all the tests pass by running `npm run test`
- Use `eslint` and `prettier` to lint and format the code respectively
- Commit your changes
- Create a PR

### Contribution workflow for UI restructuring

- Checkout to `meshery-ui-restructuring` branch
- Checkout to a new branch locally by `git checkout -b ui/restructuring/--name--` where `--name--` should be meaningful to the work that you are about to do
- Go to `/ui`
- Make sure all the tests pass by running `npm run test`
- Make the changes
- Write tests if necessary
- Make sure all the tests pass by running `npm run test`
- Use `eslint` and `prettier` to lint and format the code respectively
- Commit your changes
- Create a PR that targets `meshery-ui-restructuring` branch


## Tech stack used in Meshery UI
- Meshery UI uses NextJs to do server side rendering of ReactJS components. 
- Redux is being used for state management along with Redux-Thunk and Redux-Observable.
- MaterialUI is being used as the styling framework.
- Billboard.js library is being used to display various charts, and comparison graphs in Meshery UI.
- Relay is used as the GraphQL client

## Folder structure in Meshery UI
**App**

- Application configurations are put inside this directory

**Components**

- All the globally reusable react components are put inside this directory

**Features**
- This folder contains all the features of Meshery UI separated by concern
- If we think of Meshery UI as maintaining a huge tree of state, then we are slicing out the parts of state that are more connected into its own file so that it is easy to manage. 

**Lib**
- This folder contains the custom libraries that are used in this project

**Utils**
- Contains global utility functions

**Styles**
- Contains global theme declarations and global style declarations

> **NOTE:**  *Pages* and *Public* default NextJS directories


## Writing Tests

- We use Jest with Testing library/React for writing unit tests and integration tests
- Some things that we have to keep in mind while writing tests are, 
	- *Do not write tests for implementation details*
	- Write tests keeping in mind the perspective of a user and not a developer
	- Reading these blogs/articles before writing tests is highly recommended since it will help us in setting up some common principles amongst us contributors. - [Testing classifications](https://kentcdodds.com/blog/the-testing-trophy-and-testing-classifications) , [React docs about testing](https://reactjs.org/docs/testing.html), [JEST Docs](https://jestjs.io/docs/tutorial-react) , [RTL Docs](https://testing-library.com/docs/react-testing-library/intro/)

## Conventions

**Modular strucuture**

- Keeping the code modular is very important for it to be scalable and manageable
- One of the factors that affect the modularity of code is transparency in exports and imports
- By looking at just the exports and imports of a file, we can draw a dependency graph for the code
- So, it is important to make the imports and exports more transparent
- which is why, we have an `index.js` file for every module where all the exports will be written so that we have a single place to find what things are all exported by a module.

**Material UI**

- Always use MUI `styled` API
- Reuse Components as much as possible
- Never miss an opportunity to derive styling from the theme


## Philosophy

- Reuse code as much as possible
- “All code is guilty, until proven innocent.”
- Use comments, but don't overuse them

## JSDocs

- We use JSDocs to get type information about the data sets
- [JsDoc docs](https://jsdoc.app/)


## FAQ

1. How to change the target branch for a PR?

When opening a PR, at the top, we have to option to change the target branches as shown in the figure below. We can use it to create pull requests that can be merged into branches other than `master`

<img width="1431" alt="Screenshot 2021-09-25 at 12 26 31 AM" src="https://user-images.githubusercontent.com/75248557/134727998-5152a88d-0251-408c-97f3-c6ff0b8b389a.png">



<p style="text-align: center"><em>If you'll like to go to the main Meshery Contributor guide <a href="../CONTRIBUTING.md">click here</a></em></p>
