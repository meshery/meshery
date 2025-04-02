---
layout: default
title: Publishing a Design
abstract: Learn to publish and manage designs in Meshery with clear state transitions, role-based permissions, and common workflow FAQs.
permalink: extensions/publishing-a-design
language: en
display-title: "false"
list: include
type: extensions
category: kanvas
---

# Publishing a Design

Publishing a design makes it visible to all Meshery Cloud users and anonymous visitors. This guide explains the publishing workflow, state management, permissions, and operational constraints.

## Publishing Workflow

#### Step 1: Access Extensions UI
To begin, navigate to Meshery [Kanvas](https://playground.meshery.io/extension/meshmap) from the main dashboard.

#### Step 2: Select Design in Sidebar
Users can find designs through two different methods. Clicking the **info ("i")** button provides more details about the selected design.

- **Option 1:** Directly view designs in the sidebar under "My Designs."

  <a href="/assets/img/kanvas/designs-select-design.png">
    <img src="/assets/img/kanvas/designs-select-design.png" style="width:50%; height:auto;" alt="Design Details">
  </a> 

- **Option 2:** Go to Configuration → Designs to find the list of stored designs.

  <a href="/assets/img/kanvas/designs-select-sidebar.png">
    <img src="/assets/img/kanvas/designs-select-sidebar.png" style="width:50%; height:auto;" alt="Sidebar Design List">
  </a>

#### Step 3: Submit Design Metadata
After selecting your design, fill out the publication form:

- **Type:** Select the category that best represents your design.
- **Technology:** Specify the relevant technology associated with your design.
- **Description:** Briefly explain the purpose of the design, along with its intended and unintended usage.
- **Caveats and Considerations:** Add any specific stipulations, known behaviors, or important details that users should be aware of before using the design.

<a href="/assets/img/kanvas/publish-form.png">
  <img src="/assets/img/kanvas/publish-form.png" style="width:50%; height:auto;" alt="Publish Modal">
</a>

#### Step 4: Review Process
Once a design is submitted for publishing, it goes through an approval workflow based on the user's role:

- **Admin submissions:** Published immediately.  
- **User submissions:** Enter a "Pending Review" state.

Submissions in the “Pending Review” state will remain unpublished until they are manually reviewed and either approved or rejected by an Organization Admin or Provider Admin. This review process may take some time depending on the availability of reviewers. Regardless of approval or rejection, submitters will receive an email notification with the decision.

<a href="/assets/img/kanvas/approval-queue.png">
  <img src="/assets/img/kanvas/approval-queue.png" style="width:50%; height:auto;" alt="Approval Queue">
</a>

#### Step 5: Where to Find Your Published Design  
Once the review process is complete, you will receive an email notification informing you of the decision.  

- If your design is approved: 
  - It will no longer appear in "My Designs."  
  - Instead, it will be listed in the **[Catalog](https://cloud.layer5.io/catalog)** as a publicly available entry.  
  - This ensures that published designs remain accessible to all users while keeping personal design spaces uncluttered.  

- If your design is rejected:
  - You will receive an email notification with a rejection reason. 
  - Rejected designs **cannot be resubmitted** directly.  
  - If you want to revise and submit it again, you must clone the design, make changes, and submit it as a new entry.

## State Management

### 1. Design State Lifecycle
A design transitions through multiple states from creation to publication. The diagram below visually represents this process.

<a href="/assets/img/kanvas/Publishing-flow.svg">
  <img src="/assets/img/kanvas/Publishing-flow.svg" style="width:50%; height:auto;" alt="Publishing Flow">
</a>

#### Phases of Publishing Process
- **Pre-Publish:** Users freely create and edit designs.  
- **Pending Review:** Submitted designs undergo an approval process.  
- **Published:** Approved designs are locked and listed publicly.  
- **Withdrawn:** Unpublished designs return to private storage.

### 2. Design State Characteristics

| **State Stage**    | **Visibility**                    | **Operability**                                | **Key Restrictions**                                              |
|--------------------|----------------------------------|------------------------------------------------|--------------------------------------------------------------------|
| **Pre-Publish**    | Sidebar Design List & Configuration → Designs               | Free edit/delete/rename (editable by anyone)   |  No restrictions                                                                    |
| **Pending Review** | Sidebar Design List & Configuration → Designs                | Editable by all users (including guests)       | Deleting during this stage creates dead entries in the approval list |
| **Published**      | Category List (hidden in sidebar) | View-only copies (auto-appended "_copy")       | Original design permanently locked (edit in category only updates metadata) |
| **Withdrawn**      | Returns to sidebar as private     | Free edit/delete/rename                        | Original cannot be republished; must create a new canvas with identical content |

**Critical Rules:**
- Withdrawn designs require re-submission as new entries.
- Published designs are immutable.

### 3. Permission Controls

| **Operation**       | **Guest** | **Regular User** | **Owner** | **Admin** |
|---------------------|-----------|------------------|-----------|------------------------|
| Edit Pending Design | ✔         | ✔                | ✔         |  ✔                  |
| Submit for Review     | ✔         | ✔                | ✔         |✔                  |
| Unpublish           | ✘         | ✘                | ✔         | ✔                  |

**Key Notes:**
- Catalog approval queues are org-specific—only members of an organization can see its pending approvals, unless the user is a Admin.

## Frequently Asked Questions (FAQ) 

#### 1. When my design is in "Pending Review," can I still edit it? Will the changes take effect?
Yes, you can edit your design while it is in the "Pending Review" state, and any modifications will be automatically reflected in the submitted design. No need to resubmit the request.

#### 2. After my design is approved and published, can I modify it? Will the category be updated?
No, once a design is **published**, it becomes **immutable**. You cannot directly edit or modify the contents. The category and metadata remain locked to ensure version consistency. However, you can create a new version by cloning the design and making modifications.

#### 3. After my design is published, can I modify the category field?
No, the category field **cannot be changed** after the design is published. If a category update is required, you must clone the design, update the category, and submit it as a new entry.

#### 4. If my design is denied (rejected), can I submit it again?  
No, once a design is denied, it **cannot be resubmitted**. However, you can clone the design, make adjustments, and submit it as a new entry for review.

#### 5. Can I unpublish a design after it has been published?  
Only Admins and Owner can unpublish designs. Regular users and guests cannot perform this action.

#### 6. What is the difference between Public, Private, and Published?  

| **State**    | **Visibility** | **Editability** | **Notes** |
|-------------|---------------|----------------|-----------|
| **Public**   | Visible to all Meshery Cloud users | Fully editable | Available for all users to access |
| **Private**  | Only visible to the owner and organization members | Fully editable | Used for drafts and internal work |
| **Published** | Approved and locked for public access | Cannot be modified | Ensures design consistency and prevents unauthorized edits |  