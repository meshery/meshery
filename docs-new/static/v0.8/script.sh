#!/usr/bin/env bash

# This is a version archive script, executed when a new minor Meshery release is made.
# A "minor" release like v0.7.0 or v0.8.0
# (semver)

if [ -f site-dir.zip ]; then
    rm -rf site-dir.zip
fi
zip -r site-dir.zip ./docs/_site
