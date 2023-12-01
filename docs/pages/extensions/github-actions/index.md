---
layout: default
title: GitHub Actions
permalink: extensions/github-actions
language: en
abstract: GitHub Actions with Meshery streamlines development by automating Git repository synchronization and visual validation in pull requests.
display-title: "false"
list: include
type: extensions
---

# <img style="height: 3rem; width: 3rem;" src="{{site.baseurl}}/assets/img/meshmap-icon-color.svg" /><img style="height: 3rem; width: 3rem;" src="{{site.baseurl}}/assets/img/github.svg" />  GitHub Actions: MeshMap Snapshot

GitHub Actions with Meshery streamlines and enhances your development workflow by integrating Git version control with Meshery's robust features. This methodology automates the synchronization of your actual infrastructure state with the desired state defined in your Git repository. With Meshery. Meshery enables GitOps by offering a comprehensive view of your infrastructure, allowing you to verify workload designs and Kubernetes cluster configurations before merging pull requests.

## MeshMap Snapshot GitHub Actions

Meshery introduces the MeshMap Snapshot GitHub Actions, a tool that allows you to visualize changes in your infrastructure directly in your pull requests. With MeshMap Snapshot, you can:

- See your deployment changes pull request-to-pull request.
- Get snapshots of your infrastructure directly in your pull requests.

### Getting ready for GitOps

**Step 1: Get access to MeshMap**

To initiate the GitOps workflow with Meshery, it is essential to **[sign up](https://layer5.io/cloud-native-management/meshmap#:~:text=Signup%20for,MeshMap%20Beta)** for MeshMap. MeshMap serves as a fundamental component, providing visual insights into infrastructure changes and facilitating the synchronization between the actual and desired states. Ensure that you have proper access to MeshMap to leverage its capabilities for effective GitOps implementation.

**Step 2: Enable GitHub Actions: MeshMap Snapshot Extension in Meshery Playground**

In the Meshery Playground environment, navigate to the [Extensions](https://playground.meshery.io/extensions) section and specifically enable the **GitHub Actions: MeshMap Snapshot**. By enabling this extension, you can seamlessly integrate MeshMap Snapshot into your GitOps workflow, providing a visual layer to validate and manage modifications. This streamlined process within Meshery Playground ensures that you can efficiently leverage the capabilities of MeshMap Snapshot for enhanced GitOps practices in a controlled and experimental environment.

**Step 3: Next**

Upon clicking the `Enable` button, you will be seamlessly redirected to [Meshery Cloud](https://meshery.layer5.io/dashboard). Once there, a **"Where do you want to start?"** pop-up will appear. Choose the option that says **"Visualize your code on GitHub"** from the available selections. This step marks the initiation of setting up Meshery Cloud with GitHub.

## Steps for Setting up Meshery Cloud with GitHub

**Step 1: Connect App:** 

Begin by connecting your application to Meshery. Ensure Meshery has access to your application for effective synchronization.

**Step 2: Select Repository:** 

Choose the Git repository where your application's desired state is defined. Meshery will use this repository for synchronization.

**Step 3: Configure Application Path:** 

Specify the path within your Git repository where the application configurations are stored. This ensures Meshery knows where to find and apply the desired configurations.

**Step 4: Configure Secrets:** 

Provide any necessary secrets or authentication tokens required for Meshery to access your Git repository securely.

**Step 5: Add Workflow:** 

Define the workflow within your Git repository. This includes specifying the triggers, jobs, and steps Meshery should follow when synchronizing with the desired state.

**Step 6: Finish:** 

Once configurations are set up, finalize the process. Meshery will now automatically manage the synchronization between the actual and desired states based on your defined workflow.

## Customizing Snapshot Workflow Triggers in MeshMap Snapshot
You can configure your workflows to run when specific activity on GitHub happens, at a scheduled time, or when an event outside of GitHub occurs.

### About events that trigger workflows
GitHub Actions provides a variety of events that can trigger workflows, allowing you to automate your software development process. Each event corresponds to a specific activity, such as creating a pull request, pushing code to a repository, or releasing a new version.
