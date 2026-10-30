# v0.3.5

LLMS index: [llms.txt](/pr-preview/pr-22671/llms.txt)

---

### What's new?

- Mesheryctl
  - Removal of `init` as a command exposed to users. This command's functionality is used internal to `mesheryctl start`. A new `start --check` command will provide preflight check functionality in `init`'s place.
