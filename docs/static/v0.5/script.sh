#! /usr/bin/env bash

if [ -f site-dir.zip ]; then
    rm -rf site-dir.zip
fi
zip -r site-dir.zip ./docs/_site