#!/usr/bin/env bash

# This is a version archive script, executed when a new minor Meshery release is made.

if [ -f site-dir.zip ]; then
    rm -rf site-dir.zip
fi
zip -r site-dir.zip ./docs/_site
