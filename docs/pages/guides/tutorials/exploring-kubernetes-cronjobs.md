---
layout: default
title: Exploring Kubernetes CronJobs
permalink: guides/tutorials/exploring-kubernetes-cronjobs
model: kubernetes
kind: cronjobs
type: guides
category: tutorials
language: en
list: include
abstract: "In this tutorial, we will explore how to use Meshery Playground, an interactive live cluster environment, to perform hands-on labs for managing Kubernetes CronJobs."
---

Introduction:
In this tutorial, we will explore Kubernetes CronJobs, a resource that allows you to run jobs periodically at specified intervals. We will use Meshery Playground, an interactive live cluster environment, to perform hands-on labs for working with CronJobs in Kubernetes.

Prerequisites:
- Basic understanding of Kubernetes concepts.
- Meshery Playground access. If you don't have an account, you can sign up at [Meshery Playground](https://meshery.layer5.io/play).

Lab Scenario: Scheduled Backups using CronJobs

Objective:
Learn how to use Kubernetes CronJobs to schedule and automate periodic backups of a sample application.

### Steps:

#### 1. **Accessing Meshery Playground:**
   - Log in to the [Meshery Playground](https://meshery.layer5.io/play) using your credentials.
   - Once logged in, navigate to the Meshery Playground dashboard.

#### 2. **Deploying an Application:**
   - Deploy a simple application that needs periodic backups. You can use a web application or any other application of your choice.


# Deploying a web application

To deploy a web application on Meshery, follow these steps:


1. Navigate to the Configuration menu in Meshery.
2. Select "Designs" from the menu.
3. Click on the "Import Design" button.
[![Configuration menu]({{site.baseurl}}/assets/img/meshery-design/configuration-menu-design-import.png)]({{site.baseurl}}/assets/img/meshery-design/configuration-menu-design-import.png)
4. Fill in the required details for your application or design, including the Design File Name, Design Type, and choose your preferred upload method:
   - For File Upload, click "Browse" to locate and select the design file.
   - For URL Import, paste the link to the design.
5. Click the "Import" button to initiate the upload process.
[![Click import button]({{site.baseurl}}/assets/img/meshery-design/click-import.png)]({{site.baseurl}}/assets/img/meshery-design/click-import.png)
6. Upon successful import, a pop-up message will confirm the completion, and it will also indicate that your design has been automatically saved in Meshery.
[![Design is auto-saved]({{site.baseurl}}/assets/img/meshery-design/design-auto-save.png)]({{site.baseurl}}/assets/img/meshery-design/design-auto-save.png)
7. Locate your design and click the deploy button to initiate the deployment process.
[![Located App]({{site.baseurl}}/assets/img/meshery-design/app-deploy.png)]({{site.baseurl}}/assets/img/meshery-design/app-deploy.png)
8. Once the popup window opens, wait for the dry run to complete, and then click the deploy button.
[![Dry Run]({{site.baseurl}}/assets/img/meshery-design/click-deploy.png)]({{site.baseurl}}/assets/img/meshery-design/click-deploy.png)
9. Upon successful deployment, a confirmation pop-up will appear, indicating that your app has been successfully deployed.
[![Deployment Success]({{site.baseurl}}/assets/img/meshery-design/deploy-success.png)]({{site.baseurl}}/assets/img/meshery-design/deploy-success.png)

<!-- 

Convey to user that Kanvas Designs are auto-saved. 

-->


#### 3. **Creating a CronJob for Backups:**

- Create a CronJob that runs a backup script at a specified interval to back up the application data.

<!--

Show user how to do this using Kanvas Designer to drag and drop components and configure them.

-->

1. Open the Kanvas tab located in the left panel.
2. Upon opening Kanvas, ensure that you are on the Design tab, which can be found at the top center of the canvas.
3. Navigate to the Design option located in the top menu of the left panel. Using the search bar, type in the name of your app, which in this instance is the Minecraft App.
4. Once your app appears in the list, click on it to upload the design file onto the canvas.
[![Navigate Kanvas]({{site.baseurl}}/assets/img/kanvas/navigate-kanvas.png)]({{site.baseurl}}/assets/img/kanvas/navigate-kanvas.png)
5. Locate the control panel at the bottom of the canvas and choose the Kubernetes option.
6. Using the search bar, enter "Cron Job" and click on the corresponding icon to display it on the canvas.
[![Select CronJob item]({{site.baseurl}}/assets/img/kanvas/select-cronjob.png)]({{site.baseurl}}/assets/img/kanvas/select-cronjob.png)
7. Once the cronjob component appears on the canvas, click on it to open the toolbar and begin configuration.
[![CronJob Toolbar]({{site.baseurl}}/assets/img/kanvas/toolbar-cronjob.png)]({{site.baseurl}}/assets/img/kanvas/toolbar-cronjob.png)
8. Fill out the details for the CronJob i.e Name field: "backup-cronjob". Under **Spec**, enter the specifications, i.e schedule: "0 * * * *" to run every hour.
[![CronJob Toolbar]({{site.baseurl}}/assets/img/kanvas/tool-bar.png)]({{site.baseurl}}/assets/img/kanvas/tool-bar.png)
9. Locate the **Save As** icon in the top right, once the pop up modal opens, give your design a name, then click save.
[![Save CronJob]({{site.baseurl}}/assets/img/kanvas/save.png)]({{site.baseurl}}/assets/img/kanvas/save.png)

10. Thereafter click on the **Action** drop down menu also located in the top right and click on the **Deploy** option.
[![Deploy CronJob]({{site.baseurl}}/assets/img/kanvas/deploy.png)]({{site.baseurl}}/assets/img/kanvas/deploy.png)
11. In the pop-up window, review and correct any errors as necessary. Then, click the deploy button.
[![Deploy CronJob]({{site.baseurl}}/assets/img/kanvas/deploy-app.png)]({{site.baseurl}}/assets/img/kanvas/deploy-app.png)
12. You'll receive a confirmation message indicating that your app has been successfully deployed.

#### 4. **Verifying CronJob Execution:**
   - Monitor the execution of the CronJob and verify that backups are created at the specified intervals.

To view the resources created for the CronJob, we will utilize the Visualize tab of the Kanvas. A view will be created with necessary filters to display the relevant resources.

   1. Ensure that you are on the Visualize tab, located at the top center of the canvas.
   2. Give the view a name.
   3. Click on the filter icon.
   4. Choose appropriate filters; for this lab, select the kind item that was set earlier during the CronJob creation, i.e., 'CronJob'.
   5. Click the filter icon again to close. This action should display a filtered view with only your resources, similar to the screenshot below:
   [![Vizualize CronJob]({{site.baseurl}}/assets/img/kanvas/view.png)]({{site.baseurl}}/assets/img/kanvas/view.png)


<!-- 

Show user how to use Views and filters in Kanvas Visualizer.

-->


#### 5. **Scaling and Updating CronJobs:**
   - Explore how to scale the CronJob or update its schedule by modifying the CronJob configuration.

   1. Return to the  Designer tab
   2. Select the CronJob Design:
        Choose the design that contains the CronJob you want to scale from the list of existing designs.
   3. Locate the CronJob Component:
        Within the design canvas, identify the representation of the CronJob you wish to scale. It should be labeled as "CronJob" or have a specific icon associated with CronJobs.
        [![Vizualize CronJob]({{site.baseurl}}/assets/img/kanvas/design-cronjob.png)]({{site.baseurl}}/assets/img/kanvas/design-cronjob.png)
   4.  Select the CronJob Component:
        Click on the CronJob component to open the toolbar..
   5. Continue to configure your cronjob with your preffered specifications that match your desired scaling options.
   6. Adjust the Number of Replicas or Parallelism:
        Within the toolbar, locate the field related to the number of replicas or parallelism for the CronJob. Adjust this value to scale the CronJob up or down according to your requirements.
        [![Vizualize CronJob]({{site.baseurl}}/assets/img/kanvas/scale.png)]({{site.baseurl}}/assets/img/kanvas/scale.png)

   7.  Save Changes:
        After verifying the adjustments, save the changes made to the CronJob settings within the Kanvas Designer interface to ensure they are retained for future reference.
        [![Save CronJob]({{site.baseurl}}/assets/img/kanvas/save.png)]({{site.baseurl}}/assets/img/kanvas/save.png)
<!-- 

Show user how to use Designs and components in Kanvas Designer.

-->

Use Meshery Playground to visualize the changes and observe the impact on the scheduled backups.

#### 6. **Clean-Up:**
   - Delete the CronJob and application resources after completing the lab.

1. Identify the CronJob Component:
    Within the design canvas, find the representation of the CronJob you wish to delete. It should be labeled as "CronJob" or have a specific icon associated with CronJobs.
2. Select the CronJob Component:
    Click on the CronJob component to open the tooltip. This action will enable access to the delete icon. Click to delete the CronJob.
    [![Save CronJob]({{site.baseurl}}/assets/img/kanvas/delete.png)]({{site.baseurl}}/assets/img/kanvas/delete.png)
3. Save Changes:
    After deleting the CronJob, save the changes made within the Kanvas Designer interface to reflect the cleanup.
    [![Save CronJob]({{site.baseurl}}/assets/img/kanvas/save-app.png)]({{site.baseurl}}/assets/img/kanvas/save-app.png)
<!-- 

Show user how to use Designs and components in Kanvas Designer.

-->



#### 7. **Saving and Sharing**
  Share your scenario with other Meshery users or the community for collaborative learning.

  
1. Save Your Scenario:
   - Click the save option in Kanvas Designer and give your scenario a descriptive name.

2. Make Design Public:
   - Toggle the visibility of your design to "Public" to allow others to view it.

3. Share Your Design:
   - Copy the shareable link or invite collaborators directly from Kanvas Designer.

4. Invite Friends to Collaborate:
   - Share the link with friends or collaborators to enable collaboration on your design.

5. Confirm Sharing Settings:
   - Ensure sharing settings align with your preferences, such as viewing, editing, or commenting permissions.

6. Save Changes:
   - Save any changes made to the sharing settings to apply them.

<!-- 

Show user how to make Design public and share with other users in Kanvas Designer.

-->

### Conclusion
Congratulations! You've successfully completed the lab on exploring Kubernetes CronJobs using Meshery Playground. This hands-on experience has provided valuable insights into scheduling and automating tasks in Kubernetes using CronJobs. Explore more scenarios in the Meshery Playground to enhance your skills in cloud-native technologies.
