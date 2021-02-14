require "json"

$common_ignore_paths = [
  "CHANGELOG.md",
  "README.md",
  "concourse/**",
].to_json
