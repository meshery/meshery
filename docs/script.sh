#!/usr/bin/env bash

zip_file="site-dir.zip"
target_dir="./docs/_site"

# Check if the zip file exists
if [ -f "$zip_file" ]; then
    echo "Removing existing $zip_file..."
    rm "$zip_file"
fi

# Create a new zip file
echo "Creating $zip_file..."
zip -r "$zip_file" "$target_dir"

# Check the zip operation result
if [ $? -eq 0 ]; then
    echo "Zip file created successfully."
else
    echo "Failed to create zip file."
fi
