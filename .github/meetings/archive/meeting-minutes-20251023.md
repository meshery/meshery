### April 2nd, 2025

**Attendees**: 

* Nikita B, Matthieu Evrin, Amit Amrutiya, Vivek VishalAabid  
* Varad Gupta , Tanishq Maheshwari, Zihan Kuang, Ijeoma Eti  
* Gourav Kumar Shaw, Pratik Chandra Pal  
* VInay Cheripally, Yash Sharma  
* Yogendar Singh, Prajwal KP  
* Oyefule Oluwatayo,Raghuraj Pratap Yadav  
* Faheem Mushtaq, Mahadevan KSFiona Murugi, Kapil Ramawani  
* Mercy Boma Naps-Nkari 

*Send regrets:* 

* Harsh Prakash(exams), Faizan Shaikh, Sangram Rath(kubecon), Yash Sharma(kubecon)

**Standing Agenda: Component Sync:**

| CLI | [development](https://docs.google.com/spreadsheets/d/1q63sIGAuCnIeDs8PeM-0BAkNj8BBgPUXhLbe1Y-318o/edit?gid=0#gid=0) (Matthieu) | [test plan](https://docs.google.com/spreadsheets/d/13Ir4gfaKoAX9r8qYjAFFl_U9ntke4X5ndREY1T7bnVs/edit?gid=0#gid=0) (??) |
| ----- | :---- | :---- |
| **UI** | [development](https://docs.google.com/spreadsheets/d/1RlUcr-iLtCFCsV1VZvWD0jBijPJAx-Rk-1GLLsL9sSI/edit?gid=0#gid=0) (Amit) | [test plan](https://docs.google.com/spreadsheets/d/1Ck_5q7U_vLSIDTtplugG3pCVC5zugXgTHtO7T7-yL8g/edit?gid=1336550606#gid=1336550606) (??) |
| *Release in-progress [v0.9.0](https://discuss.layer5.io/t/meshery-v0-9-roadmap/6296) [GitHub Milestone](https://github.com/meshery/meshery/milestone/6)* |  |  |

[**Contributor Training Schedule**](https://github.com/meshery/meshery/issues/14090)**:**

| \# | Topic and Component | Technology | Trainer | Date |
| :---: | ----- | ----- | ----- | ----- |
| 1 | [Contributing to Meshery Docs](https://docs.meshery.io/project/contributing/contributing-docs) | Hugo,,Jekyll, Markdown, HTML, CSS | Zihan Kuang | April 10th |
| 2 | [End-to-end Testing in Meshery CL](https://docs.meshery.io/project/contributing/contributing-cli-tests) | Bash (BATS | Riya Garg | April 17th |
| 3 | [Contributing to Meshery Build and Release](https://docs.meshery.io/project/contributing/build-and-release) | CI / CD, GitHub Workflows | Sangram Rath | April 24th |
| 4 | [Contributing to Meshery CLI](https://docs.meshery.io/project/contributing/contributing-cli) | Golang, Cobra | Matthieu Evrin,Aadhitya  | April 30th |
| 5 | [End-to-end Testing in Meshery UI](https://docs.meshery.io/project/contributing/contributing-ui-tests#:~:text=To%20automate%20functional%20integration%20and,not%20break%20the%20existing%20functionality.) | Playwright | Ian Whitney | May 8th |
| 6 | [Contributing to Meshery UI](https://docs.meshery.io/project/contributing/contributing-ui) | React, Sistent, CASL | Amit Amrutiya |  |
| 7 | [Contributing to Meshery Server](https://docs.meshery.io/project/contributing/contributing-server) | Golang, Gorilla, OpenAPI Schema |  |  |
| 8 | [Contributing to Meshery Models](https://docs.meshery.io/project/contributing/contributing-models) |  | Lee Calcote |  |

**Announcements:** 

* KubeCon EU 2025 Security Slam [preparation and participation](https://docs.google.com/document/d/15gFHw-A6uOchW-0U54-C6Miwv7sjvHFayynwboZwCbc/edit?tab=t.0). See [blog post](https://meshery.io/blog/meshery-security-slam-kubecon-eu-2025). Celebrate Security Sentinels next week.  
* ![][image2]![][image3]

* [Meshery Reaches 7,000 Stars](https://meshery.io/blog/meshery-reaches-7000-stars)  
    
* Meshery GitHub Org Split (draft post)   
  * Core: [https://github.com/meshery](https://github.com/meshery)    
  * Extensions: [https://github.com/meshery-extensions](https://github.com/meshery-extensions)

**Topics**

* \[Nikita B\] model init command  
* \[Ijeoma\] \- list Kanvas Snapshot Helm Plugin in extension page [\#14224](https://github.com/meshery/meshery/pull/14224)  
* \[Prajwal KP\]  Set appropriate permissions for GitHub workflows \#[14168](https://github.com/meshery/meshery/pull/14168)  
* \[Varad Gupta\]  Fix: Optimize Playwright Test Performance in CI/CD \#[14195](https://github.com/meshery/meshery/pull/14195)  
  	             Fix: Flaky Connection UI Tests \#[14233](https://github.com/meshery/meshery/pull/14233)  
* \[Baivab Mukhopdhyay\] Enhance error handling in login command [14165](https://github.com/meshery/meshery/pull/14165)  
* \[Fiona Murugi\] fix design apply status from 0 to non zero if a file path is invalid  \#[14189](https://github.com/meshery/meshery/pull/14189)  
* \[Yogendar Singh\]  docs: updated docs for docker-extension install \#[14176](https://github.com/meshery/meshery/pull/14176)  
* \[Yashwanth Reddy\] \-  fix : issue with sql injection in handlers/meshsync\_handler.go \#[14222](https://github.com/meshery/meshery/pull/14222)  
* \[Matthieu EVRIN\]   
  * \[server\] Enhancement request to add model delete support [\#14136](https://github.com/meshery/meshery/issues/14136)  
  * \[mesheryctl\]  
    * **model** command reuse common Fetch, List methods \+ Add tests [\#14171](https://github.com/meshery/meshery/pull/14171)  
    * **org** command renamed **organizations** [\#14169](https://github.com/meshery/meshery/pull/14169)  
      * Fix panic  
      * reuse common Fetch, List methods \+ Add tests  
* \[Mahadevan, Karan\] Feat: create programs layout [\#2085](https://github.com/meshery/meshery.io/pull/2085)  
* \[Tanishq Maheswari\]\[review\] fix : routing between meshsync and connection [\#14067](https://github.com/meshery/meshery/pull/14067)  
* \[Zihan Kuang\] docs: update account deletion and linking behavior [\#477](https://github.com/layer5io/docs/pull/477)  
* \[Gourav Kumar Shaw\] mesheryctl registry generate spreadsheet URL fix.merge \#[14087](https://github.com/meshery/meshery/pull/14087)  
* \[Amit Amrutiya\] \[call for volunteers\] migrating axios/fetch to rtk query  
* Faheem \- Fix all regression issues pending in meshery cloud  
* \[Vivek Vishal\] \[call for volunteers\] [https://github.com/meshery/meshery/issues/13546](https://github.com/meshery/meshery/issues/13546)  
* \[Aabid Sofi\]