### April 9th, 2025

**Attendees**: 

* Vivek Vishala, Lee Calcote, Amit Amrutiya, Zihan Kuang, Abhinav Sharma, Faheem Mushtaq  
* Kavya samudrala, Mahadevan KS, Yash Sharma, Tanishq Maheshwari, Vinay Cheripally  
* Varad Gupta, Pratik Chandra Pal , Riya Garg, Nikita B,, Ijeoma Eti, Rishikesh Maddhesiya,Tamaghna, Mercy Boma Naps-Nkari

*Send regrets:* 

* 

**Standing Agenda: Component Sync:**

| CLI | [development](https://docs.google.com/spreadsheets/d/1q63sIGAuCnIeDs8PeM-0BAkNj8BBgPUXhLbe1Y-318o/edit?gid=0#gid=0) (Matthieu) | [test plan](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit?gid=0#gid=0) (??) |
| ----- | :---- | :---- |
| **UI** | [development](https://docs.google.com/spreadsheets/d/1RlUcr-iLtCFCsV1VZvWD0jBijPJAx-Rk-1GLLsL9sSI/edit?gid=0#gid=0) (Amit) | [test plan](https://docs.google.com/spreadsheets/d/1Ck_5q7U_vLSIDTtplugG3pCVC5zugXgTHtO7T7-yL8g/edit?gid=1336550606#gid=1336550606) (??) |
| *Release in-progress [v0.9.0](https://discuss.layer5.io/t/meshery-v0-9-roadmap/6296) [GitHub Milestone](https://github.com/meshery/meshery/milestone/6)* |  |  |

[**Contributor Training Schedule**](https://github.com/meshery/meshery/issues/14090)**:**

| \# | Topic and Component | Technology | Trainer | Date |
| :---: | ----- | ----- | ----- | ----- |
| 1 | [Contributing to Meshery Docs](https://docs.meshery.io/project/contributing/contributing-docs) | Hugo, Jekyll, Markdown, HTML, CSS | Zihan Kuang | April 10th |
| 2 | [End-to-end Testing in Meshery CLI](https://docs.meshery.io/project/contributing/contributing-cli-tests) | Bash (BATS | Riya Garg | April 17th |
| 3 | [Contributing to Meshery Build and Release](https://docs.meshery.io/project/contributing/build-and-release) | CI / CD, GitHub Workflows | Sangram Rath | \! April 24th |
| 4 | [Contributing to Meshery CLI](https://docs.meshery.io/project/contributing/contributing-cli) | Golang, Cobra | Matthieu Evrin,Aadhitya  | April 30th |
| 5 | [End-to-end Testing in Meshery UI](https://docs.meshery.io/project/contributing/contributing-ui-tests#:~:text=To%20automate%20functional%20integration%20and,not%20break%20the%20existing%20functionality.) | Playwright | Ian Whitney | \! May 8th |
| 6 | [Contributing to Meshery UI](https://docs.meshery.io/project/contributing/contributing-ui) | React, Sistent, CASL | Amit Amrutiya |  |
| 7 | [Contributing to Meshery Server](https://docs.meshery.io/project/contributing/contributing-server) | Golang, Gorilla, OpenAPI Schema |  |  |
| 8 | [Contributing to Meshery Models](https://docs.meshery.io/project/contributing/contributing-models) | All cloud native tech | Lee Calcote |  |

**Announcements:** 

* \[Recognition\]\[Security Sentinels\] \- [celebration](https://discuss.layer5.io/t/meshery-at-security-slam-kubecon-eu-2025/6840/2?u=vishalvivekm) 🎉  
* \[Meet GSoc’25 Mentors\] [https://x.com/mesheryio/status/19097...](https://x.com/mesheryio/status/1909740016407814553)  
  


  
**Topics**

* \[Nikita B\] \- “generate yaml from json schema” \-  Let’s be schema-driven. 🙂  
* \[Abhinav Sharma\] \[Draft Tutorial\] \- Kafka Deployment on Kubernetes Using Meshery Kanvas  
* \[Sumit Shandilya\]  Fixed Cross-site Scripting (XSS) in prometheus\_handlers.go \#[14241](https://github.com/meshery/meshery/pull/14241)  
* \[Riya Garg\]    
  * \[mesheryctl\] Move registry logger setup from initialisation function to PersistentPreRun of Registry command \#[14341](https://github.com/meshery/meshery/pull/14341)  
  * \[Server\] Resolve $ref schema  [\#713](https://github.com/meshery/meshkit/pull/713)  
  * Improve Error logging for Source URL related issues in Model Generation [\#709](https://github.com/meshery/meshkit/pull/709)  
* \[Alberto\] Chore\] Update the go and alpine Docker images \#[14249](https://github.com/meshery/meshery/pull/14249)  
* \[Yash Sharma\]  Improved mesheryctl design view command \#[14340](https://github.com/meshery/meshery/pull/14340)  
* \[Varad Gupta\]  
*  Fix: Extension Spec \#[14333](https://github.com/meshery/meshery/pull/14333)  
*  Fix: Connection Spec \-\> Delete Kubernetes cluster connections \#[14316](https://github.com/meshery/meshery/pull/14316)  
*  Fix: Connection Spec \-\> Transition tests \#[14313](https://github.com/meshery/meshery/pull/14313)  
* \[Arnav Kapoor\]  Added ResourceGraphDefinition model via mesheryctl \#[14304](https://github.com/meshery/meshery/pull/14304)  
* \[Anuj Kumar Sharma\]  fixed failing performance tests \#[14252](https://github.com/meshery/meshery/pull/14252)  
* \[Kapil Ramwani\] Issues currently faced in Prometheus and Grafana Support. [\#13847](https://github.com/meshery/meshery/issues/13847#issuecomment-2787052758) . [One Pager on Issues and Ask to maintainers](https://docs.google.com/document/d/1Rk8SotJT3ZXAfMVPhqh5iTHbOivSUzAFXuKB4AmtL48/edit?tab=t.0#heading=h.pcmxjsi73rtw)  
* \[Gourav Kumar Shaw\]  
* \[Harsh Prakash\]  
* \[Tanishq Maheswari\]  
        \-    Update on kanvas snapshot \[workflow being success but email not sent\]  
        \-    fix : routing between connection and meshsync [\#14336](https://github.com/meshery/meshery/pull/14336)

                    \-   add e2e test for import,list and design commands [\#14273](https://github.com/meshery/meshery/pull/14273)  

* \[Mahadevan\] Add Switch Component to Sistent Components page [\#6380](https://github.com/layer5io/layer5/pull/6380)  
* \[Faizan Shaik\] Update on Helm export functionality  
* \[Zihan Kuang\] fix: restyle blockquote for better visual emphasis [\#14329](https://github.com/meshery/meshery/pull/14329)  
* \[Aabid Sofi\]   
*          \-  Fix component updater   
*  \[Amit Amrutiya\] \- update typing filter component         
* \[Faheem Mushtaq\]  
  	\- Working on RTK Query Migration  
* \[“your name”\] \- “your topic”  
* \[“your name”\] \- “your topic”