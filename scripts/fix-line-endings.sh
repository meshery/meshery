#!/bin/bash

# Git Line Ending Fix Script for Meshery and other projects
# This script configures Git to handle line endings properly on Windows systems

echo "🔧 Fixing Git Line Ending Configuration..."

# Check if we're in a Git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "❌ Error: Not in a Git repository. Please run this script from within a Git repository."
    exit 1
fi

echo "📁 Current directory: $(pwd)"

# Configure Git settings for the current repository
echo "⚙️  Configuring local Git settings..."
git config core.autocrlf false
git config core.safecrlf warn
git config core.eol lf

echo "✅ Local Git configuration updated:"
echo "   - core.autocrlf: false"
echo "   - core.safecrlf: warn"
echo "   - core.eol: lf"

# Ask user if they want to apply global settings
read -p "🌍 Do you want to apply these settings globally for all Git repositories? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "⚙️  Configuring global Git settings..."
    git config --global core.autocrlf false
    git config --global core.safecrlf warn
    echo "✅ Global Git configuration updated"
else
    echo "ℹ️  Skipped global configuration"
fi

# Check if .gitattributes exists
if [ -f ".gitattributes" ]; then
    echo "📄 .gitattributes file already exists"
    read -p "🔄 Do you want to backup and update it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        cp .gitattributes .gitattributes.backup
        echo "💾 Backup created: .gitattributes.backup"
        echo "⚠️  Please manually merge the new .gitattributes content with your existing file"
    fi
else
    echo "📄 No .gitattributes file found. You should create one with proper line ending rules."
fi

# Show current Git configuration
echo ""
echo "📋 Current Git configuration:"
echo "   Local settings:"
git config --local --get-regexp "core\.(autocrlf|safecrlf|eol)" || echo "   No local core settings found"
echo "   Global settings:"
git config --global --get-regexp "core\.(autocrlf|safecrlf|eol)" || echo "   No global core settings found"

echo ""
echo "🎉 Git line ending configuration completed!"
echo ""
echo "📚 Next steps:"
echo "   1. Ensure you have a proper .gitattributes file"
echo "   2. Test by adding files: git add ."
echo "   3. Check for reduced line ending warnings"
echo ""
echo "💡 For more information, see GIT_LINE_ENDING_FIX.md"
