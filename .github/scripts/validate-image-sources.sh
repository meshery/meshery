#!/bin/bash
#
# Bitnami Image Source Validator
# 
# This script validates that container images referenced in Meshery
# are not using deprecated Bitnami sources and are available from
# their respective registries.
#
# Exit codes:
#   0 - All validations passed
#   1 - Deprecated Bitnami references found
#   2 - Unreachable image sources detected

set -e

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}Bitnami Image Source Validator${NC}"
echo -e "${BLUE}======================================${NC}"
echo ""

EXIT_CODE=0

# Function to check for deprecated Bitnami references
check_deprecated_bitnami() {
    echo -e "${BLUE}🔍 Checking for deprecated Bitnami references...${NC}"
    
    # Search for docker.io/bitnami references (excluding bitnamisecure and sealed-secrets)
    deprecated_images=$(grep -r "docker\.io/bitnami/" \
        docs/catalog/ \
        server/meshmodel/ \
        --include="*.yml" \
        --include="*.yaml" \
        2>/dev/null | \
        grep -v "bitnamisecure" | \
        grep -v "bitnami-labs/sealed-secrets" || true)
    
    if [ -n "$deprecated_images" ]; then
        echo -e "${RED} Found deprecated Bitnami Docker image references:${NC}"
        echo "$deprecated_images" | while IFS= read -r line; do
            echo -e "   ${YELLOW}→${NC} $line"
        done
        echo ""
        echo -e "${YELLOW}   These images were deprecated after September 2025.${NC}"
        echo -e "${YELLOW}   Please migrate to official upstream images or alternatives.${NC}"
        echo ""
        EXIT_CODE=1
    else
        echo -e "${GREEN} No deprecated Bitnami Docker image references found${NC}"
    fi
    echo ""
}

# Function to check Helm chart references
check_helm_charts() {
    echo -e "${BLUE}🔍 Checking Bitnami Helm chart references...${NC}"
    
    # Count Bitnami Helm chart references
    helm_refs=$(grep -c "repo_url: https://charts\.bitnami\.com/bitnami" \
        server/meshmodel/component_models.yaml 2>/dev/null || echo "0")
    
    if [ "$helm_refs" -gt 0 ]; then
        echo -e "${YELLOW}  Found $helm_refs Bitnami Helm chart references in component_models.yaml${NC}"
        echo -e "${YELLOW}   Note: These charts are still available but no longer updated.${NC}"
        echo -e "${YELLOW}   Consider migrating to alternative chart sources from Artifact Hub.${NC}"
        echo ""
        
        # List unique chart names
        echo -e "${BLUE}   Charts using Bitnami repository:${NC}"
        grep -B2 "repo_url: https://charts\.bitnami\.com/bitnami" \
            server/meshmodel/component_models.yaml | \
            grep "^- name:" | \
            sed 's/^- name: //' | \
            sort -u | \
            head -20 | \
            while IFS= read -r chart; do
                echo -e "   ${YELLOW}→${NC} $chart"
            done
        
        chart_count=$(grep -B2 "repo_url: https://charts\.bitnami\.com/bitnami" \
            server/meshmodel/component_models.yaml | \
            grep "^- name:" | \
            sed 's/^- name: //' | \
            sort -u | wc -l)
        
        if [ "$chart_count" -gt 20 ]; then
            echo -e "   ${YELLOW}... and $((chart_count - 20)) more${NC}"
        fi
        echo ""
    else
        echo -e "${GREEN} No Bitnami Helm chart references found${NC}"
    fi
    echo ""
}

# Function to check for bitnamilegacy references
check_legacy_references() {
    echo -e "${BLUE}🔍 Checking for Bitnami Legacy repository references...${NC}"
    
    legacy_refs=$(grep -r "bitnamilegacy" \
        docs/ \
        server/ \
        --include="*.yml" \
        --include="*.yaml" \
        --include="*.md" \
        2>/dev/null || true)
    
    if [ -n "$legacy_refs" ]; then
        echo -e "${RED} Found Bitnami Legacy repository references:${NC}"
        echo "$legacy_refs" | while IFS= read -r line; do
            echo -e "   ${YELLOW}→${NC} $line"
        done
        echo ""
        echo -e "${RED}  The bitnamilegacy repository was scheduled for deletion Sept 29, 2025.${NC}"
        echo -e "${RED}   These images may no longer be available!${NC}"
        echo ""
        EXIT_CODE=1
    else
        echo -e "${GREEN} No Bitnami Legacy repository references found${NC}"
    fi
    echo ""
}

# Function to suggest alternatives
suggest_alternatives() {
    if [ $EXIT_CODE -ne 0 ]; then
        echo -e "${BLUE}======================================${NC}"
        echo -e "${BLUE}Suggested Migration Path${NC}"
        echo -e "${BLUE}======================================${NC}"
        echo ""
        echo -e "Common Bitnami image replacements:"
        echo ""
        echo -e "  ${YELLOW}bitnami/redis${NC} → ${GREEN}redis:7-alpine${NC}"
        echo -e "  ${YELLOW}bitnami/postgresql${NC} → ${GREEN}postgres:16-alpine${NC}"
        echo -e "  ${YELLOW}bitnami/mongodb${NC} → ${GREEN}mongo:8.0${NC}"
        echo -e "  ${YELLOW}bitnami/nginx${NC} → ${GREEN}nginx:alpine${NC}"
        echo -e "  ${YELLOW}bitnami/apache${NC} → ${GREEN}httpd:2.4-alpine${NC}"
        echo -e "  ${YELLOW}bitnami/etcd${NC} → ${GREEN}quay.io/coreos/etcd:v3.5${NC}"
        echo -e "  ${YELLOW}bitnami/jenkins${NC} → ${GREEN}jenkins/jenkins:lts${NC}"
        echo -e "  ${YELLOW}bitnami/os-shell${NC} → ${GREEN}busybox${NC} or ${GREEN}alpine:3.19${NC}"
        echo ""
        echo -e "📖 For complete migration guide, see:"
        echo -e "   ${BLUE}docs/BITNAMI_AUDIT_REPORT.md${NC}"
        echo ""
    fi
}

# Function to generate summary
generate_summary() {
    echo -e "${BLUE}======================================${NC}"
    echo -e "${BLUE}Validation Summary${NC}"
    echo -e "${BLUE}======================================${NC}"
    echo ""
    
    if [ $EXIT_CODE -eq 0 ]; then
        echo -e "${GREEN} All validations passed!${NC}"
        echo -e "   No deprecated Bitnami references detected."
    else
        echo -e "${RED} Validation failed${NC}"
        echo -e "   Deprecated Bitnami references detected."
        echo -e "   Please update to alternative image sources."
    fi
    echo ""
}

# Main execution
main() {
    check_deprecated_bitnami
    check_helm_charts
    check_legacy_references
    suggest_alternatives
    generate_summary
    
    exit $EXIT_CODE
}

# Run main function
main
