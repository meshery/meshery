#!/bin/bash

# Advanced troubleshooting script for Kanvas-Snapshot buffer allocation errors
# This script provides detailed analysis and potential workarounds

echo "=== Advanced Kanvas-Snapshot Diagnostic Tool ==="
echo "================================================="
echo ""

# Check system resources
echo "System Information:"
echo "  Node.js: $(node --version 2>/dev/null || echo 'Not found')"
echo "  Memory: $(free -h 2>/dev/null | grep 'Mem:' | awk '{print $2}' || echo 'Unknown')"
echo "  Disk: $(df -h . 2>/dev/null | tail -1 | awk '{print $4}' || echo 'Unknown')"
echo ""

# Analyze the specific error pattern
echo "Buffer Allocation Error Analysis:"
echo "  Error Pattern: 'size is invalid. Received -6667200'"
echo "  Calculation: -6667200 = likely negative width * height"
echo "  Screenshot size from logs: 5760x2823 = 16,273,920 pixels"
echo "  If each pixel = 4 bytes (RGBA), total = 65,095,680 bytes (~65MB)"
echo ""

# Mathematical analysis of the negative value
echo "Mathematical Analysis:"
echo "  -6667200 could be:"
echo "    • Signed 32-bit integer overflow"
echo "    • Incorrect crop coordinates (negative width/height)"
echo "    • Buffer size calculation error in image library"
echo "    • Memory alignment issue"
echo ""

# Potential solutions
echo "Recommended Solutions (in order of priority):"
echo ""
echo "1. IMMEDIATE FIX - Add error handling to workflow:"
echo "   - Add 'continue-on-error: true' to Kanvas step"
echo "   - Implement retry mechanism with smaller viewport"
echo ""
echo "2. CONFIGURATION FIXES:"
echo "   - Reduce browser viewport size (e.g., 1280x720)"
echo "   - Increase Node.js memory: NODE_OPTIONS='--max-old-space-size=8192'"
echo "   - Add Chrome flags for memory optimization"
echo ""
echo "3. WORKFLOW IMPROVEMENTS:"
echo "   - Add pre-flight memory checks"
echo "   - Implement fallback screenshot mechanism"
echo "   - Add timeout and retry logic"
echo ""

# Generate a proposed workflow fix
echo "=== PROPOSED WORKFLOW FIX ==="
cat << 'EOF'

# Add this step before Kanvas-Snapshot action:
- name: Configure environment for large screenshots
  run: |
    echo "NODE_OPTIONS=--max-old-space-size=8192" >> $GITHUB_ENV
    echo "CYPRESS_CACHE_FOLDER=/home/runner/.cache/Cypress" >> $GITHUB_ENV
    
# Modify the Kanvas step to handle errors gracefully:
- id: test_result
  continue-on-error: true  # Don't fail the entire workflow
  uses: layer5labs/Kanvas-Snapshot@v0.2.39
  with:
    githubToken: ${{ secrets.GITHUB_TOKEN }}
    mesheryToken: ${{ secrets.MESHERY_TOKEN }}
    prNumber: ${{ env.PULL_NO }}
    application_type: Kubernetes Manifest
    filePath: ${{ inputs.fileName == '' && 'install/deployment_yamls/k8s' || inputs.fileName }}

EOF

echo ""
echo "=== IMMEDIATE ACTION NEEDED ==="
echo "The buffer allocation error is preventing CI from completing."
echo "Recommend adding 'continue-on-error: true' to the Kanvas step"
echo "so that other CI checks can still pass while this issue is resolved."
echo ""
echo "Long-term solution requires coordination with Kanvas-Snapshot maintainers"
echo "to fix the underlying buffer allocation calculation bug."
echo ""
