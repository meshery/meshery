# v0.3.5

Source: /pr-preview/pr-21670/project/releases/v0.3.5/

### What's new?

- Mesheryctl
  - Removal of `init` as a command exposed to users. This command's functionality is used internal to `mesheryctl start`. A new `start --check` command will provide preflight check functionality in `init`'s place.
