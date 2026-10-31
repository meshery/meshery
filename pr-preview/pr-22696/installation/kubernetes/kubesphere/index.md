# KubeSphere

> Install Meshery on KubeSphere

---

LLMS index: [llms.txt](/pr-preview/pr-22696/llms.txt)

---

<h1>Quick Start with KubeSphere <img src="/pr-preview/pr-22696/installation/kubernetes/kubesphere/images/kubesphere.png" style="width:35px;height:35px;" /></h1>

<div class="prereqs"><h4>Prerequisites</h4>

1. Install the Meshery command line client, 
<a href="/pr-preview/pr-22696/installation/mesheryctl/" class="meshery-light">
    mesheryctl
</a>.
</div>


[Meshery](https://meshery.io/) is the open source, cloud native manager for Kubernetes, all kinds of cloud native infrastructure, and their workloads.

This guide walks you through deploying Meshery from the KubeSphere App Store.


## Prerequisites

- Enable the OpenPitrix system (the KubeSphere App Store).
- Create a workspace, a project, and a user account. The account must be a platform regular user, invited to the project as an operator with the `operator` role. This guide uses the user `project-regular`, the project `demo-project`, and the workspace `demo-workspace`. For details, see Create Workspaces, Projects, Users and Roles in the KubeSphere documentation.


## Hands-on Lab

Follow these steps in order.

### 1. <b>Deploy Meshery from the App Store</b>


1. On the **Overview** page of the project `demo-project`, click **App Store** in the upper-left corner.
2. Search for **Meshery** in the App Store, and click the search result to open the app.

    ![meshery-app](images/meshery-app.png)
3. On the **App Information** page, click **Install** in the upper-right corner.

    ![meshery-install](images/Meshery-install.png)

4. On the **App Settings** page, set the application **Name**, **Location** (your namespace), and **App Version**, then click **Next** in the upper-right corner.

    ![meshery-info](images/Meshery-info.png)

5. Edit the **values.yaml** file as needed, or click **Install** to use the default configuration. The chart creates the Meshery service as a `LoadBalancer` by default. If your cluster has no load balancer, set `service.type` to `NodePort` here.

    ![meshery-yaml](images/Meshery-yaml.png)

6. Wait for the deployment to finish. When it completes, KubeSphere shows **Meshery** as **Running**.

    ![meshery-app-running](images/Meshery-app-running.png)



### 2. <b>Access the Meshery Dashboard</b>


1. Go to **Services** and click the Meshery service.
2. In the **Resource Status** page, copy the external IP address of the Meshery service or, if you set the service type to `NodePort`, its **NodePort**.

    ![meshery-service](images/Meshery-service.png)

3. Open Meshery UI in your browser at **${EXTERNAL_IP}:9081**, or at **${NodeIP}:${NODEPORT}** for a `NodePort` service.

    ![meshery-dashboard](images/meshery-dashboard.png)

<div class="related-discussions">
  <h3>Recent Discussions with "meshery" Tag</h3><ul><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-seeking-an-experienced-meshery-contributor-as-a-guide/8166" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate: Seeking an Experienced Meshery Contributor as a Guide
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-to-reactivate-my-slack-account/8168" target="_blank" rel="noopener noreferrer">
            Request to Reactivate My Slack Account
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/bug-new-catalog-publish-request-email-renders-with-oversized-social-icons-and-squashed-content-on-mobile/8138" target="_blank" rel="noopener noreferrer">
            [Bug] &#34;New Catalog publish request&#34; email renders with oversized social icons and squashed content on mobile
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/new-contributor-intro-vedit-agrawal-bughunterx2101-looking-to-pair-with-a-meshmate-for-meshery-ui-ai-adapter/8172" target="_blank" rel="noopener noreferrer">
            New Contributor Intro: Vedit Agrawal (@BugHunterX2101) | Looking to pair with a MeshMate for Meshery UI &amp; AI Adapter
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/looking-for-a-meshmate-frontend-backend-cloud-native-development/8161" target="_blank" rel="noopener noreferrer">
            Looking for a MeshMate — Frontend, Backend &amp; Cloud-Native Development
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/request-for-a-meshmate-active-meshery-contributor/8158" target="_blank" rel="noopener noreferrer">
            Request for a MeshMate — Active Meshery Contributor
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/how-to-present-an-issue-pr-in-a-development-meeting-feedback-wanted/8107" target="_blank" rel="noopener noreferrer">
            How to Present an Issue/PR in a Development Meeting: feedback wanted
          </a>
        </li><li>
          <a href="https://discuss.meshery.io/t/newbie-here-excited-to-join/8052" target="_blank" rel="noopener noreferrer">
            Newbie here -- Excited to join
          </a>
        </li></ul><p>
    <a href="https://discuss.meshery.io/tag/meshery" target="_blank" rel="noopener noreferrer">
      View all discussions tagged with <code>meshery</code>
    </a>
  </p>
</div>
