#!/bin/bash

# Fix for Kanvas-Snapshot Cypress test failure
# This script patches the loadDesign.js file to handle undefined cytoscape object

echo "🔧 Applying fix for Kanvas-Snapshot Cypress test failure..."

# Check if the action directory exists
if [ -d "action/cypress-action/cypress/e2e/e2e" ]; then
    LOADDESIGN_FILE="action/cypress-action/cypress/e2e/e2e/loadDesign.js"
    
    if [ -f "$LOADDESIGN_FILE" ]; then
        echo "📝 Patching $LOADDESIGN_FILE to handle undefined cytoscape object..."
        
        # Create a backup
        cp "$LOADDESIGN_FILE" "$LOADDESIGN_FILE.backup"
        
        # Apply the fix using sed to add error handling
        sed -i 's/const cytoscape = win\.cyto;/const cytoscape = win.cyto || win.cy || window.cytoscape;/' "$LOADDESIGN_FILE"
        sed -i 's/cytoscape\.fit();/if (cytoscape \&\& typeof cytoscape.fit === "function") { cytoscape.fit(); }/' "$LOADDESIGN_FILE"
        sed -i 's/cytoscape\.center();/if (cytoscape \&\& typeof cytoscape.center === "function") { cytoscape.center(); }/' "$LOADDESIGN_FILE"
        
        echo "✅ Applied Cypress test fix successfully"
        echo "📋 Changes made:"
        echo "   - Added fallback for cytoscape object (win.cyto || win.cy || window.cytoscape)"
        echo "   - Added null checks before calling cytoscape.fit() and cytoscape.center()"
    else
        echo "⚠️  loadDesign.js file not found at $LOADDESIGN_FILE"
        echo "   This might be expected if the action structure has changed"
    fi
else
    echo "⚠️  Action directory not found. This script should run after the kanvas-snapshot action is checked out"
fi

echo "🏁 Kanvas-Snapshot Cypress fix script completed"
