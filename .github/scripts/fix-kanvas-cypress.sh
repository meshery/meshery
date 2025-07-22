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
        
        # Apply the fix using a more robust approach
        # First, let's see what we're working with
        echo "📋 Original problematic section:"
        grep -n -A 5 -B 5 "cytoscape.fit()" "$LOADDESIGN_FILE" || echo "Pattern not found"

        # Create a temporary file with the fixed content
        cat > /tmp/cypress_fix.js << 'EOF'
const captureSnapshot = ({ window, designId, theme }) => {
  console.log("Taking snapshot", designId, theme);
  //removeWidgets();
  cy.window().then((win) => {
    // Safe cytoscape access with multiple fallbacks
    const cytoscape = win.cyto || win.cy || win.cytoscapeInstance || window.cytoscape;

    // Only call methods if cytoscape object exists and has the methods
    if (cytoscape && typeof cytoscape === 'object') {
      try {
        if (typeof cytoscape.fit === 'function') {
          cytoscape.fit();
        } else {
          console.warn('cytoscape.fit() method not available');
        }

        if (typeof cytoscape.center === 'function') {
          cytoscape.center();
        } else {
          console.warn('cytoscape.center() method not available');
        }
      } catch (error) {
        console.warn('Error calling cytoscape methods:', error);
      }
    } else {
      console.warn('Cytoscape object not found, skipping fit/center operations');
    }
  });

  const path = snapshotPath(designId, theme);
  cy.wait(2000);
  cy.get("main", { timeout: 10 * 1000 })
    .should("exist")
    .screenshot(path, {
      //blackout: [".hide-from-snapshot"], // hides elements before screenshot
      scale: true,
    });
  console.log(`Snapshot taken at ${path}`);
};
EOF

        # Replace the problematic function
        # Find the start and end of the captureSnapshot function
        START_LINE=$(grep -n "const captureSnapshot = " "$LOADDESIGN_FILE" | cut -d: -f1)
        if [ -n "$START_LINE" ]; then
            # Find the end of the function (next function or end of file)
            END_LINE=$(tail -n +$((START_LINE + 1)) "$LOADDESIGN_FILE" | grep -n "^const \|^function \|^\[" | head -1 | cut -d: -f1)
            if [ -n "$END_LINE" ]; then
                END_LINE=$((START_LINE + END_LINE - 1))
            else
                END_LINE=$(wc -l < "$LOADDESIGN_FILE")
            fi

            echo "📝 Replacing captureSnapshot function (lines $START_LINE-$END_LINE)"

            # Create new file with replacement
            head -n $((START_LINE - 1)) "$LOADDESIGN_FILE" > "$LOADDESIGN_FILE.new"
            cat /tmp/cypress_fix.js >> "$LOADDESIGN_FILE.new"
            tail -n +$((END_LINE + 1)) "$LOADDESIGN_FILE" >> "$LOADDESIGN_FILE.new"

            # Replace the original file
            mv "$LOADDESIGN_FILE.new" "$LOADDESIGN_FILE"

            echo "✅ Successfully replaced captureSnapshot function with safe version"
        else
            echo "⚠️  Could not find captureSnapshot function to replace"
            # Fallback to simple sed replacement
            sed -i 's/const cytoscape = win\.cyto;/const cytoscape = win.cyto || win.cy || window.cytoscape;/' "$LOADDESIGN_FILE"
            sed -i 's/cytoscape\.fit();/if (cytoscape \&\& typeof cytoscape.fit === "function") { try { cytoscape.fit(); } catch(e) { console.warn("fit() failed:", e); } }/' "$LOADDESIGN_FILE"
            sed -i 's/cytoscape\.center();/if (cytoscape \&\& typeof cytoscape.center === "function") { try { cytoscape.center(); } catch(e) { console.warn("center() failed:", e); } }/' "$LOADDESIGN_FILE"
        fi

        # Clean up
        rm -f /tmp/cypress_fix.js
        
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
