---
layout: default
title: Exploring Kubernetes CronJobs
abstract: Learn Kubernetes CronJobs with Meshery Playground
permalink: guides/tutorials/exploring-kubernetes-cronjobs
# redirect_from: guides/tutorials/
type: guides
category: tutorials
language: en
abstract: "In this tutorial, we will explore how to use Meshery Playground, an interactive live cluster environment, to perform hands-on labs for managing Kubernetes CronJobs."
list: include
abstract: Meshery is the self-service engineering platform, enabling collaborative design and operation of cloud and cloud native infrastructure.
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

Convey to user that MeshMap Designs are auto-saved. 

-->


#### 3. **Creating a CronJob for Backups:**

- Create a CronJob that runs a backup script at a specified interval to back up the application data.

<!--

Show user how to do this using MeshMap Designer to drag and drop components and configure them.

-->

1. Open the MeshMap tab located in the left panel.
2. Upon opening MeshMap, ensure that you are on the Design tab, which can be found at the top center of the canvas.
3. Navigate to the Design option located in the top menu of the left panel. Using the search bar, type in the name of your app, which in this instance is the Minecraft App.
4. Once your app appears in the list, click on it to upload the design file onto the canvas.
[![Navigate MeshMap]({{site.baseurl}}/assets/img/meshmap/navigate-meshmap.png)]({{site.baseurl}}/assets/img/meshmap/navigate-meshmap.png)
5. Locate the control panel at the bottom of the canvas and choose the Kubernetes option.
6. Using the search bar, enter "Cron Job" and click on the corresponding icon to display it on the canvas.
[![Navigate MeshMap]({{site.baseurl}}/assets/img/meshmap/select-cronjob.png)]({{site.baseurl}}/assets/img/meshmap/select-cronjob.png)
7. Once the cronjob component appears on the canvas, click on it to open the toolbar and begin configuration.
[![CronJob Toolbar]({{site.baseurl}}/assets/img/meshmap/tool-bar.png)]({{site.baseurl}}/assets/img/meshmap/tool-bar.png)
8. Fill out the details for the CronJob i.e Name field: "backup-cronjob". Under **Spec**, enter the specifications, i.e schedule: "0 * * * *" to run every hour.
    
   
<details>
<summary>
Alternate: Command line instructions
</summary>

\```yaml
# backup-cronjob.yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: backup-cronjob
spec:
  schedule: "0 * * * *"  # Run every hour
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: backup-container
            image: backup-image:latest  # Replace with your backup script image
            # Add volume mounts and commands as needed for backup
          restartPolicy: OnFailure
\```

</details>

9. Click on the deploy icon located in the top right corner.
10. In the pop-up window, review and correct any errors as necessary. Then, click the deploy button. You'll receive a confirmation message indicating that your app has been successfully deployed.
[![Deploy CronJob]({{site.baseurl}}/assets/img/meshmap/deploy.png)]({{site.baseurl}}/assets/img/meshmap/deploy.png)

Apply the CronJob configuration:

\```bash
mesheryctl pattern apply -f backup-cronjob.yaml
\```

#### 4. **Verifying CronJob Execution:**
   - Monitor the execution of the CronJob and verify that backups are created at the specified intervals.


<!-- 

Show user how to use Views and filters in MeshMap Visualizer.

-->


#### 5. **Scaling and Updating CronJobs:**
   - Explore how to scale the CronJob or update its schedule by modifying the CronJob configuration.


<!-- 

Show user how to use Designs and components in MeshMap Designer.

-->


\```bash
kubectl edit cronjob backup-cronjob
\```

Use Meshery Playground to visualize the changes and observe the impact on the scheduled backups.

#### 6. **Clean-Up:**
   - Delete the CronJob and application resources after completing the lab.


<!-- 

Show user how to use Designs and components in MeshMap Designer.

-->



#### 7. **Saving and Sharing**
   - Save your scenario in Meshery Playground for future reference.
   - Invite friends to collaborate in your design.

<!-- 

Show user how to make Design public and share with other users in MeshMap Designer.

-->

Share your scenario with other Meshery users or the community for collaborative learning.

### Conclusion
Congratulations! You've successfully completed the lab on exploring Kubernetes CronJobs using Meshery Playground. This hands-on experience has provided valuable insights into scheduling and automating tasks in Kubernetes using CronJobs. Explore more scenarios in the Meshery Playground to enhance your skills in cloud-native technologies.

{% include related-content.html type="tutorials" %}