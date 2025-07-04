#!/bin/bash

# Script to help diagnose and fix Kanvas snapshot test issues
# This script provides debugging information and potential fixes for common Cypress/screenshot issues

echo "Kanvas Snapshot Test Diagnostics"
echo "================================="
echo ""

# Check Node.js version (important for buffer allocation)
echo "Node.js version:"
node --version
echo ""

# Check available memory
echo "Available memory:"
free -h 2>/dev/null || echo "Memory info not available"
echo ""

# Check for common issues
echo "Checking for common issues:"
echo ""

# 1. Check if there are any large files that might cause memory issues
echo "Looking for large files in the project..."
find . -type f -size +100M 2>/dev/null | head -5
echo ""

# 2. Check browser/Cypress related packages if they exist
if [ -f "package.json" ]; then
    echo "Checking for Cypress/browser dependencies..."
    grep -E "(cypress|puppeteer|playwright|chrome|browser)" package.json 2>/dev/null || echo "No browser testing packages found in main package.json"
    echo ""
fi

# 3. Display potential fixes
echo "Potential fixes for buffer allocation errors:"
echo "1. Update Kanvas-Snapshot to latest version (already done: v0.2.37)"
echo "2. Increase Node.js heap size: NODE_OPTIONS='--max-old-space-size=8192'"
echo "3. Reduce screenshot resolution or canvas size"
echo "4. Check for memory leaks in browser automation"
echo "5. Use headless browser with smaller viewport"
echo ""

echo "The error 'size is invalid. Received -6667200' suggests:"
echo "- Image crop dimensions are calculated incorrectly"
echo "- Possible negative width/height in crop operation"
echo "- Canvas size calculation overflow"
echo ""

echo "Next steps:"
echo "1. Monitor the CI run with updated Kanvas-Snapshot v0.2.37"
echo "2. Check if the external action has memory limits"
echo "3. Consider reducing test complexity if issues persist"
echo ""

echo "Script completed."
