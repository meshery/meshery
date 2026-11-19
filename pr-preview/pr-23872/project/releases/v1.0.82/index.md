# v1.0.82

LLMS index: [llms.txt](/pr-preview/pr-23872/llms.txt)

---

## What's New
## 🔤 General
- Revert "\[Server\] Prevent MeshSync cache loss when the resync handler is unavailable" @miacycle (#23896)
- \[Server\] Prevent MeshSync cache loss when the resync handler is unavailable @Nixantsingh943 (#23878)
- \[Server\] Guard nil provider, validate path UUIDs, and propagate status codes in pattern handlers @patelshivani2283-lab (#23891)
- \[Server\] Bound request body size and guard against DoS in MeshModelGenerationHandler @Supersaiyan101 (#23864)
- \[Server\] Fix MeshSync flush handling for disconnected connections and nil handler @Nixantsingh943 (#23846)
- \[Server\] Introduce DataLoaders in GraphQL resolvers to resolve N+1 query bottleneck @Sarthak-Pandey (#22851)
- \[Server\] Import a raw .wasm filter from a URL @Asheesh7298 (#23842)
- \[Server\] Resolve array wildcard paths during policy mutation and deletion cleanup @nishant-devv (#23670)

## ⌨️ Meshery CLI

- \[mesheryctl\] Validate model search positional arguments @Karman580 (#22565)
- \[mesheryctl\] Fix profile lookup when paginating performance results @FayezBast (#22775)
- \[mesheryctl\] Do not parse a registry summary when model generate skips registration @Eternity0207 (#23868)
- \[mesheryctl\] Fail system check --adapter for an adapter that is not configured @Eternity0207 (#23849)
- \[mesheryctl\] Keep the current context when system reset uses --context @Eternity0207 (#23653)
- \[mesheryctl\] Show each design's file in the design deploy selection prompt @mayank609 (#23286)
- \[mesheryctl\] Include container names in Kubernetes logs @itsCodeTide (#23847)
- \[mesheryctl\] Keep the callback URL when system update upgrades the Helm charts @Eternity0207 (#23845)

## 🖥 Meshery UI

- fix(server):Reject oversized telemetry responses @FayezBast (#22762)
- \[UI\] Guard urlEncodeArrayParam against undefined input and prevent malformed query strings @madhavzanwar (#23674)
- \[UI\] Turn analytics preference switches back when saving fails @mayank609 (#23640)
- \[UI\] Shallow-copy object schema and guard array properties in LazyComponentForm @abhinavkdeval08-design (#23363)
- \[UI\] Show the real error when fetching Kubernetes resources fails @mayank609 (#23856)

## 🧰 Maintenance

- \[CI\] Replace actions/first-interaction with github-script in welcome workflow @Junnygram (#23803)

## 📖 Documentation

- \[Server\] Guard unchecked nil provider dereferences, nil logger panics, and unbounded body reads in UserTestPreference handlers @Amansingh0807 (#23861)
- fix(server): keep Meshery Operator installed across shared-cluster connection lifecycle @yi-nuo426 (#23261)
- fix(server):Reject oversized telemetry responses @FayezBast (#22762)
- \[mesheryctl\] Fix profile lookup when paginating performance results @FayezBast (#22775)
- \[Server\] Restrict filter URL imports to public destinations @itsCodeTide (#23874)
- \[Docs\] Adjust mobile content padding to prevent sidebar toggle button overlap @itsrockmukesh (#23859)

## 👨🏽‍💻 Contributors

Thank you to our contributors for making this release possible:
@Amansingh0807, @Asheesh7298, @Eternity0207, @FayezBast, @Junnygram, @Karman580, @Maanvi212006, @Nixantsingh943, @Sarthak-Pandey, @Supersaiyan101, @abhinavkdeval08-design, @alexquincy, @carlosriosilva, @hortison, @itsCodeTide, @itsrockmukesh, @leecalcote, @madhavzanwar, @marblom007, @mayank609, @meshery-ci, @miacycle, @nishant-devv, @patelshivani2283-lab, @pontusringblom, @ritzorama, @simihablo, @suttonskate, @yi-nuo426 and Sarthak Pandey
