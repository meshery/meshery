# v1.0.77

LLMS index: [llms.txt](/pr-preview/pr-23486/llms.txt)

---

## What's New
## 🔤 General
- \[Server\] Guard nil dereferences for user, SystemID, provider, and connection in telemetry handlers @Amansingh0807 (#23338)
- \[Server\] Guard against SSRF in remote URL imports @Supersaiyan101 (#23331)
- \[Server\] Guard nil MesheryInstanceID and typed-nil INSTANCE_ID in K8s context paths @piyxsh31 (#23305)
- \[Server\] Prevent nil pointer dereference on unauthenticated /user/logout @Supersaiyan101 (#23311)
- \[Server\] Evict state machine from tracker on DeleteContext error @Atharva-Mendhulkar (#23256)
- \[Helm\] Add env.MESHSYNC_DEFAULT_DEPLOYMENT_MODE with embedded default @ayushbuilds-code (#23291)
- \[Server\] Implement MapPreferencePersister unit tests and safe logger handling @Rohith-Saran (#23293)
- \[Docs\] Fix broken help-wanted search links in README.md @ayushbuilds-code (#23289)
- \[Server\] Probe NOTFOUND Kubernetes connections so their unreachable condition can resolve @dropps07 (#23235)

## ⌨️ Meshery CLI

- \[mesheryctl\] Fix nil dereference and slice bounds panic in `environment view` @abhinavkdeval08-design (#23365)
- fix(mesheryctl): sanitize spaces and path separators in relationship view --save filename @ThakurUjjwal1706 (#23359)
- \[mesheryctl\] Added Port Fallback To Unused Port @tanav29 (#22308)
- \[mesheryctl\] Enforce exact argument validation in system provider view @ayushbuilds-code (#23352)
- \[mesheryctl\] Use MeshKit errors for token assignment and rollback @piyxsh31 (#23340)
- \[mesheryctl\] Guard nil timestamp and UUID fields across design and filter commands @Swarup-Banik (#23325)
- \[mesheryctl\] Fix connection view empty-slice panic when TotalCount > 0 @Swarup-Banik (#23328)
- \[mesheryctl\] Respect -c/--context across system token create, set, and view @piyxsh31 (#23317)
- \[mesheryctl\] Validate positional arguments in design undeploy @RajJasani (#23309)
- \[mesheryctl\] Validate positional arguments in model export @RajJasani (#23304)
- fix(mesheryctl): escape model path and validate input in relationship view @ThakurUjjwal1706 (#23302)
- \[mesheryctl\] Replace path separators in file names saved by component view @RajJasani (#23300)
- \[mesheryctl\] Fix Kubernetes log stream descriptor leak in system logs @muskan-creates352 (#22471)
- \[mesheryctl\] Parse positional mesh name argument in adapter deploy and remove (#23275) @muskan-creates352 (#23276)
- \[mesheryctl\] Fix file descriptor leak in model init @muskan-creates352 (#22438)
- \[mesheryctl\] Honor the --context flag in system dashboard, login and logout @mayank609 (#23280)
- \[mesheryctl\] Return actionable error for unsupported meshery publish target @Rohith-Saran (#23270)

## 🖥 Meshery UI

- \[UI\] Handle empty component.schema in LazyComponentForm @abhinavkdeval08-design (#23357)
- \[UI\] Offer every credential in the connection wizard and MeshSync stepper @ChackoPAbraham19 (#23334)
- \[UI\] Paginate the Credentials table on the server @ChackoPAbraham19 (#23243)
- \[UI\] Guard missing cluster context in Dashboard View and enforce RBAC query skip checks @piyxsh31 (#23319)
- \[UI\] Keep failed notification page requests retryable @ayushgade06 (#22805)
- \[UI\] Enforce RBAC skip guards on connections, preferences, and registry @piyxsh31 (#23295)
- \[UI\] Disable LegacyRJSFModal submit while pending and add submitDisabled @piyxsh31 (#23299)
- \[UI\] Fix search paging and Count sort in the database summary table @mayank609 (#23282)
- \[UI\] Reset to the first page when searching performance results @mayank609 (#23284)
- \[UI\] Skip unauthorized queries in Dashboard charts when view permissions are denied @piyxsh31 (#23272)

## 🐛 Bug Fixes

- \[Server\] Dedup inventory parent additions across relationship copies @yi-nuo426 (#23366)

## 🧰 Maintenance

- \[Server\] Dedup inventory parent additions across relationship copies @yi-nuo426 (#23366)
- \[CI\] Restrict OSPS Baseline workflow to upstream repository @AdityaChaudhary5412 (#23308)
- \[CI\] Restrict Hugo Pages deployment to upstream repository @AdityaChaudhary5412 (#23307)
- \[CI\] Slim high-volume docs pages so GitHub Pages can publish again @yi-nuo426 (#23321)
- \[CI\] Restrict Trivy image scan workflow to upstream repository @AdityaChaudhary5412 (#23306)

## 📖 Documentation

- \[mesheryctl\] Added Port Fallback To Unused Port @tanav29 (#22308)
- \[mesheryctl\] Use MeshKit errors for token assignment and rollback @piyxsh31 (#23340)
- \[Docs\] KubeCon workshop: choose-your-own-adventure Kanvas or Meshery UI/mesheryctl paths @yi-nuo426 (#23350)
- docs: restructure KubeCon workshop to publish a design, then add a logo model @yi-nuo426 (#23344)
- \[Docs\] Collapse instructor guidance in KubeCon workshop page @leecalcote (#23339)
- \[mesheryctl\] Respect -c/--context across system token create, set, and view @piyxsh31 (#23317)
- docs: restructure KubeCon ContribFest workshop for 60 minutes @yi-nuo426 (#23337)
- docs(install): add playground sandbox-cluster Broker exposure, interim MeshSync RBAC and recovery runbook @yi-nuo426 (#23250)
- \[CI\] Slim high-volume docs pages so GitHub Pages can publish again @yi-nuo426 (#23321)
- \[Docs\] Stop telling Create Model users to turn on Immediate Registration @yi-nuo426 (#23323)
- \[Docs\] Hide generated error-code pages from the sidebar @ChackoPAbraham19 (#23222)
- \[Docs\] Add KubeCon Playground workshop: build your own model and design @ritzorama (#23267)

## 👨🏽‍💻 Contributors

Thank you to our contributors for making this release possible:
@AdityaChaudhary5412, @Amansingh0807, @Atharva-Mendhulkar, @ChackoPAbraham19, @Maanvi212006, @RajJasani, @Rohith-Saran, @Supersaiyan101, @Swarup-Banik, @ThakurUjjwal1706, @abhinavkdeval08-design, @alexquincy, @ayushbuilds-code, @ayushgade06, @carlosriosilva, @dropps07, @fitzergerald, @hortison, @leecalcote, @mayank609, @miacycle, @muskan-creates352, @piyxsh31, @pontusringblom, @ritzorama, @schizophrenicmaniac, @simihablo, @suttonskate, @tanav29, @winkletinkle, @yi-nuo426, Aditya Chaudhary and Ayush
