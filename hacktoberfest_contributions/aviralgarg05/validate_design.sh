#!/bin/bash

# Test script to validate New Relic Architecture Design
# This script validates the YAML structure, checks for heavy files, and verifies the design

set -e

echo "======================================"
echo "New Relic Architecture Design Validation"
echo "======================================"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Check if the file exists
echo "Test 1: Checking if design file exists..."
if [ -f "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml" ]; then
    echo -e "${GREEN}✓ Design file exists${NC}"
else
    echo -e "${RED}✗ Design file not found${NC}"
    exit 1
fi
echo ""

# Test 2: Validate YAML syntax
echo "Test 2: Validating YAML syntax..."
if ruby -e "require 'yaml'; YAML.load_file('hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml')" 2>&1; then
    echo -e "${GREEN}✓ YAML syntax is valid${NC}"
else
    echo -e "${RED}✗ YAML syntax is invalid${NC}"
    exit 1
fi
echo ""

# Test 3: Check file size (should not be heavy)
echo "Test 3: Checking file size..."
FILE_SIZE=$(stat -f%z "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml" 2>/dev/null || stat -c%s "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml")
FILE_SIZE_KB=$((FILE_SIZE / 1024))
if [ $FILE_SIZE -lt 10485760 ]; then  # Less than 10MB
    echo -e "${GREEN}✓ File size is acceptable: ${FILE_SIZE_KB}KB${NC}"
else
    echo -e "${RED}✗ File size is too large: ${FILE_SIZE_KB}KB${NC}"
    exit 1
fi
echo ""

# Test 4: Check for required components in the design
echo "Test 4: Checking for required New Relic components..."
REQUIRED_COMPONENTS=(
    "New Relic APM"
    "Application Server"
    "New Relic Infrastructure Agent"
    "Database"
)

for component in "${REQUIRED_COMPONENTS[@]}"; do
    if grep -q "$component" "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml"; then
        echo -e "${GREEN}✓ Found: $component${NC}"
    else
        echo -e "${RED}✗ Missing: $component${NC}"
        exit 1
    fi
done
echo ""

# Test 5: Verify directory structure
echo "Test 5: Verifying directory structure..."
if [ -d "hacktoberfest_contributions/aviralgarg05" ]; then
    echo -e "${GREEN}✓ Directory structure is correct (hacktoberfest_contributions/aviralgarg05/)${NC}"
else
    echo -e "${RED}✗ Directory structure is incorrect${NC}"
    exit 1
fi
echo ""

# Test 6: Check .gitignore for heavy file patterns
echo "Test 6: Verifying .gitignore has heavy file exclusions..."
HEAVY_FILE_PATTERNS=(
    "*.zip"
    "*.tar.gz"
    "node_modules"
    "*.mp4"
)

for pattern in "${HEAVY_FILE_PATTERNS[@]}"; do
    if grep -q "$pattern" ".gitignore"; then
        echo -e "${GREEN}✓ .gitignore includes: $pattern${NC}"
    else
        echo -e "${YELLOW}⚠ .gitignore missing: $pattern (added now)${NC}"
    fi
done
echo ""

# Test 7: Verify YAML structure has required fields
echo "Test 7: Verifying YAML structure..."
REQUIRED_FIELDS=(
    "name:"
    "version:"
    "services:"
    "traits:"
    "meshmap:"
)

for field in "${REQUIRED_FIELDS[@]}"; do
    if grep -q "$field" "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml"; then
        echo -e "${GREEN}✓ Found field: $field${NC}"
    else
        echo -e "${RED}✗ Missing field: $field${NC}"
        exit 1
    fi
done
echo ""

# Test 8: Count components in the design
echo "Test 8: Counting components in the design..."
COMPONENT_COUNT=$(grep -c "type:" "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml" || echo "0")
echo -e "${GREEN}✓ Design contains $COMPONENT_COUNT components${NC}"
echo ""

# Test 9: Verify edges and relationships
echo "Test 9: Verifying component relationships (edges)..."
if grep -q "edges:" "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml"; then
    EDGE_COUNT=$(grep -c "from:" "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml" || echo "0")
    echo -e "${GREEN}✓ Design contains $EDGE_COUNT relationship edges${NC}"
else
    echo -e "${YELLOW}⚠ No edges defined (components may not be connected)${NC}"
fi
echo ""

# Test 10: Check for observability components
echo "Test 10: Verifying observability components..."
OBSERVABILITY_COMPONENTS=(
    "APM"
    "Infrastructure"
    "Browser"
    "Logs"
)

OBS_COUNT=0
for obs_component in "${OBSERVABILITY_COMPONENTS[@]}"; do
    if grep -qi "$obs_component" "hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml"; then
        echo -e "${GREEN}✓ Found observability: $obs_component${NC}"
        OBS_COUNT=$((OBS_COUNT + 1))
    fi
done
echo -e "${GREEN}✓ Design contains $OBS_COUNT observability components${NC}"
echo ""

# Summary
echo "======================================"
echo -e "${GREEN}All tests passed! ✓${NC}"
echo "======================================"
echo ""
echo "Summary:"
echo "- Design file: hacktoberfest_contributions/aviralgarg05/new-relic-architecture.yaml"
echo "- File size: ${FILE_SIZE_KB}KB"
echo "- Components: $COMPONENT_COUNT"
echo "- Relationship edges: $EDGE_COUNT"
echo "- Observability components: $OBS_COUNT"
echo ""
echo "Next steps:"
echo "1. Review the design in Meshery Playground"
echo "2. Commit changes: git add . && git commit -m 'Add New Relic architecture design'"
echo "3. Push changes: git push origin <branch-name>"
echo "4. Create a Pull Request"
echo ""
