# v1.0.72

LLMS index: [llms.txt](/pr-preview/pr-22662/llms.txt)

---

## What's New
## 🔤 General
- feat(models): add Infisical Kubernetes Operator model (secrets-operator v0.11.12) @yi-nuo426 (#22423)
- fix(server): prevent nil pointer panic on ServiceMeshFilter without type (#20988) @VimalN2005 (#22419)
- chore(deps): bump meshkit to v1.0.24 for GitHub connection changes @leecalcote (#22361)
- Bump github.com/moby/sys/user from 0.4.0 to 0.4.1 in /server/policies/wasm @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22318)

## ⌨️ Meshery CLI

- fix: \[mesheryctl\] Two --save bugs in view commands: unconditional UserHomeDir call and missing file extension @D-source1602 (#21333)
- \[Chore\] Remove leftover meshery-cpx references @hiyach28 (#22160)
- \[Chore\] Remove leftover meshery-octarine references @hiyach28 (#22161)
- \[Chore\] Remove the AWS App Mesh adapter @hiyach28 (#22355)
- \[mesheryctl\] Remove the NGINX Service Mesh adapter @hiyach28 (#22035)
- \[mesheryctl\] Fix docker not launching in macos because of hardcoded invalid path @syedmehdi03 (#22266)
- \[mesheryctl\] Fix TestReadValidationErrorMessages_Direct panic on Windows @abhijeetgaur-dev (#22334)
- \[mesheryctl\] Prevent index out of range panic in system status @rajpreetcodes (#22331)
- \[mesheryctl\] Remove dead code in healthcheck.go and registry.go @arnavmahajan630 (#22289)
- Fix/mesheryctl context create validate components @arnavmahajan630 (#22259)

## 🖥 Meshery UI

- fix(workspaces): resolve actions by workspace id @hiyach28 (#21266)
- \[Chore\] Remove leftover meshery-cpx references @hiyach28 (#22160)
- \[Chore\] Remove the AWS App Mesh adapter @hiyach28 (#22355)
- \[UI\] Remove the NGINX Service Mesh adapter @hiyach28 (#22013)
- \[UI\] Remove the AWS App Mesh adapter @hiyach28 (#22362)
- \[UI\] Fix malformed Performance Profile run URL when no Kubernetes contexts selected @VimalN2005 (#22261)
- \[UI\] Assert the delete icon uses the theme color @hiyach28 (#22262)
- \[UI\] Remove the NGINX Service Mesh adapter from the Meshery design e2e fixture @hiyach28 (#22325)
- \[UI\] Fix ping on Meshery Server connection chip @YashKumar-404 (#22332)
- \[UI\] Fix My Recent Designs error in Local Provider mode @25shivansh (#22274)
- \[UI\] Include invalid schema-validation payload type in error @pinankshah8-Aa (#22273)
- Bump sharp and next in /ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22322)
- Bump source-map-js from 1.2.1 to 1.2.2 in /ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22320)
- Bump joi from 18.2.8 to 18.2.9 in /ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22319)
- Bump next from 16.2.11 to 16.3.8 in /provider-ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22321)

## 🧰 Maintenance

- Update playground cleanup workflow to point to new infra (OCI) @sangramrath (#21548)
- Update GitHub token and add maturity level 2 @sangramrath (#22403)
- Replace expression mode with environment variable mode @sangramrath (#22389)
- \[Chore\] Remove the AWS App Mesh adapter @hiyach28 (#22355)
- Add GitHub Actions for publishing OSPS results to grc store @sangramrath (#22368)
- \[UI\] Remove the NGINX Service Mesh adapter from the Meshery design e2e fixture @hiyach28 (#22325)
- \[CI\] Publish Allure results after E2E failures @maverickaayush (#22328)
- Add OSPS Baseline Scanner tool to SECURITY-INSIGHTS.yml  @sangramrath (#22336)
- Bump sharp and next in /ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22322)
- Bump source-map-js from 1.2.1 to 1.2.2 in /ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22320)
- Bump joi from 18.2.8 to 18.2.9 in /ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22319)
- \[CI\] Fail docs workflows on npm ci errors @maverickaayush (#22326)
- Bump next from 16.2.11 to 16.3.8 in /provider-ui @\[dependabot\[bot\]\](https://github.com/apps/dependabot) (#22321)
- Update OSPS Baseline workflow configuration @sangramrath (#22301)

## 📖 Documentation

- \[Docs\] Remove nested links and redundant p tags inside summary elements @tejaswiverma121-byte (#21973)
- docs: add Favour Adeyemo introduction @BLVCK-MAMBA-6 (#22417)
- Add Yeisrael Dawit's meeting notes @YeisraelD (#22416)
- Add Om Parihar's introduction to meeting notes @salvatoreOm (#22415)
- Introduction - Okoko Jerry @OkokoJerry (#22414)
- Create parth-chandurkar introductory text.md @ParthrChandurkar (#22413)
- Add Kokole Wahid's introduction to the meetings doc @runningman145 (#22412)
- Add meeting notes for Alfred Dey @alfreddey (#22411)
- Add maanvi.md with introductory text @Maanvi212006 (#22407)
- Create adarsh-chauhan.md for personal introduction @MRIARC-08 (#22408)
- Create introduction for Amarachi Ezeonyekwere @Amarachi-Ezeonyekwere (#22410)
- \[Docs\] Remove tanzu-sm from the mesheryctl guide adapter list @hiyach28 (#22399)
- \[Newcomer's Meet\]: Add Maanvi's intro @Maanvi212006 (#22394)
- \[Chore\] Remove leftover meshery-cpx references @hiyach28 (#22160)
- \[Chore\] Remove leftover meshery-octarine references @hiyach28 (#22161)
- \[Docs\] Remove the archived meshery-tanzu-sm adapter page @hiyach28 (#22162)
- \[Chore\] Remove the AWS App Mesh adapter @hiyach28 (#22355)
- \[Docs\] Fix broken wasm-filters repository links in Architecture overview and Filter management guide @Aarushi0427 (#22316)
- \[CI\] Publish Allure results after E2E failures @maverickaayush (#22328)
- \[Docs\] Fix mesheryctl contributor path @fawcetteugene (#22299)
- \[Docs\] Fix cloned repository directory @fawcetteugene (#22297)
- \[Docs\] Remove Duplicated Content from the Catalog Page @Katotodan (#22340)
- \[Docs\] Document Kustomize package imports @BLVCK-MAMBA-6 (#22317)
- \[Docs\] Add meeting attendance for Rajpreet Singh Khurana @rajpreetcodes (#22330)

## 🔒 Security

- Update GitHub token and add maturity level 2 @sangramrath (#22403)
- Replace expression mode with environment variable mode @sangramrath (#22389)
- Add GitHub Actions for publishing OSPS results to grc store @sangramrath (#22368)
- Add OSPS Baseline Scanner tool to SECURITY-INSIGHTS.yml  @sangramrath (#22336)
- Update OSPS Baseline workflow configuration @sangramrath (#22301)

## 👨🏽‍💻 Contributors

Thank you to our contributors for making this release possible:
@25shivansh, @Aarushi0427, @Amarachi-Ezeonyekwere, @BLVCK-MAMBA-6, @D-source1602, @Katotodan, @MRIARC-08, @Maanvi212006, @OkokoJerry, @ParthrChandurkar, @VimalN2005, @YashKumar-404, @YeisraelD, @abhijeetgaur-dev, @alexquincy, @alfreddey, @arjunmehta-git, @arnavmahajan630, @dependabot\[bot\], @fawcetteugene, @fitzergerald, @hiyach28, @jamieplu, @jijillery, @l5io, @leecalcote, @lekaf974, @maverickaayush, @meshery-ci, @miacycle, @pinankshah8-Aa, @rajpreetcodes, @runningman145, @salvatoreOm, @sangramrath, @syedmehdi03, @tejaswiverma121-byte, @winkletinkle, @yi-nuo426 and \[dependabot\[bot\]\](https://github.com/apps/dependabot)
