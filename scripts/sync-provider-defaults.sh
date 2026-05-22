#!/usr/bin/env bash
# Sync default remote provider URLs from install/Makefile.core.mk into every
# consumer file in the repo. Edit the canonical block in Makefile.core.mk
# (between the AUTO-SYNC SOURCE / END AUTO-SYNC SOURCE markers) and run this
# script. Do NOT hand-edit any sentinel-managed region.

set -euo pipefail

cd "$(dirname "$0")/.."

MK="install/Makefile.core.mk"
MARK="AUTO-GENERATED-FROM-MAKEFILE"

read_mk() {
	make -s -f "$MK" "print-$1" | tr -d '"'
}

URLS=$(read_mk REMOTE_PROVIDER_URLS)
PRIMARY=$(read_mk PRIMARY_PROVIDER_URL)

# Parallel name/url arrays for the docker-extension chooser. Each entry is
# "<NAME_VAR>|<URL_VAR>"; the order here is the order rendered into the
# REMOTE_PROVIDERS JS array. Keep parallel with REMOTE_PROVIDER_URLS in
# install/Makefile.core.mk - only list the active providers.
PAIRS=(
	"MESHERY_NAME|MESHERY_CLOUD_PROD"
	"LAYER5_NAME|LAYER5_CLOUD_PROD"
)

# replace_single_line <file> <perl-regex-matching-old-line> <new-line>
# The regex must match the ENTIRE line. The new line should not include a
# trailing newline.
#
# Pre-check: a managed file must already contain the AUTO-GENERATED-FROM-MAKEFILE
# sentinel. If it does not, the regex would silently no-op and the file would
# stay out of sync without CI noticing. Fail loudly instead so the person
# adding a new managed file remembers to seed the sentinel comment first.
replace_single_line() {
	local file="$1" pattern="$2" replacement="$3"
	if [ ! -f "$file" ]; then
		echo "sync-provider-defaults: missing file $file" >&2
		exit 1
	fi
	if ! grep -qF "$MARK" "$file"; then
		echo "sync-provider-defaults: sentinel \"$MARK\" not seeded in $file." >&2
		echo "  Add the sentinel as a comment on the line you want managed, then re-run." >&2
		exit 1
	fi
	R="$replacement" perl -i -pe 'BEGIN{$r=$ENV{R}} s~'"$pattern"'~$r~' "$file"
	if ! grep -qF "$MARK" "$file"; then
		echo "sync-provider-defaults: sentinel $MARK not found in $file after rewrite" >&2
		exit 1
	fi
}

# replace_block <file> <begin-marker> <end-marker> <new-block-stdin>
# Replaces every line between the begin and end marker lines (exclusive).
# The marker lines themselves are preserved.
replace_block() {
	local file="$1" begin="$2" end="$3"
	local tmp
	tmp=$(mktemp)
	awk -v begin="$begin" -v end="$end" '
		BEGIN { state = 0 }
		state == 0 && index($0, begin) { print; state = 1; next }
		state == 1 && index($0, end)   { while ((getline line < "/dev/stdin") > 0) print line; print; state = 2; next }
		state != 1 { print }
	' "$file" > "$tmp"
	mv "$tmp" "$file"
	if ! grep -qF "$begin" "$file" || ! grep -qF "$end" "$file"; then
		echo "sync-provider-defaults: block markers not found in $file after rewrite" >&2
		exit 1
	fi
}

# ---------------------------------------------------------------------------
# install/kubernetes/helm/meshery/values.yaml
# ---------------------------------------------------------------------------
replace_single_line \
	"install/kubernetes/helm/meshery/values.yaml" \
	'^  PROVIDER_BASE_URLS:.*# '"$MARK"'.*$' \
	"  PROVIDER_BASE_URLS: \"$URLS\"  # $MARK"

# ---------------------------------------------------------------------------
# install/docker/docker-compose.yaml
# ---------------------------------------------------------------------------
replace_single_line \
	"install/docker/docker-compose.yaml" \
	'^      - "PROVIDER_BASE_URLS=.*# '"$MARK"'.*$' \
	"      - \"PROVIDER_BASE_URLS=$URLS\"  # $MARK"

# ---------------------------------------------------------------------------
# install/mesheryapp.dockerapp/docker-compose.yml
# ---------------------------------------------------------------------------
replace_single_line \
	"install/mesheryapp.dockerapp/docker-compose.yml" \
	'^      - "PROVIDER_BASE_URLS=.*# '"$MARK"'.*$' \
	"      - \"PROVIDER_BASE_URLS=$URLS\"  # $MARK"

# ---------------------------------------------------------------------------
# install/deployment_yamls/k8s/meshery-deployment.yaml
# (k8s manifest uses two lines: `- name:` then `  value:`. Patch the value.)
# ---------------------------------------------------------------------------
replace_single_line \
	"install/deployment_yamls/k8s/meshery-deployment.yaml" \
	'^          value: .*# '"$MARK"'.*$' \
	"          value: $URLS  # $MARK"

# ---------------------------------------------------------------------------
# install/docker-extension/docker-compose.yaml - single-URL product flow
# ---------------------------------------------------------------------------
replace_single_line \
	"install/docker-extension/docker-compose.yaml" \
	'^      - "PROVIDER_BASE_URLS=.*# '"$MARK"'.*$' \
	"      - \"PROVIDER_BASE_URLS=$PRIMARY\"  # $MARK"

# ---------------------------------------------------------------------------
# install/playground/docker/docker-compose.yaml - single-URL product flow
# ---------------------------------------------------------------------------
replace_single_line \
	"install/playground/docker/docker-compose.yaml" \
	'^            - PROVIDER_BASE_URLS=.*# '"$MARK"'.*$' \
	"            - PROVIDER_BASE_URLS=$PRIMARY  # $MARK"

# ---------------------------------------------------------------------------
# install/playground/docker/Makefile - single-URL product flow
# ---------------------------------------------------------------------------
replace_single_line \
	"install/playground/docker/Makefile" \
	'^PROVIDER_BASE_URLS=.*# '"$MARK"'.*$' \
	"PROVIDER_BASE_URLS=$PRIMARY  # $MARK"

# ---------------------------------------------------------------------------
# mesheryctl/pkg/utils/helpers.go - entry in the Services map
# ---------------------------------------------------------------------------
replace_single_line \
	"mesheryctl/pkg/utils/helpers.go" \
	'^			"PROVIDER_BASE_URLS=.*// '"$MARK"'.*$' \
	"			\"PROVIDER_BASE_URLS=$URLS\", // $MARK"

# ---------------------------------------------------------------------------
# ui/constants/endpoints.ts - SaaS deep-link host
# ---------------------------------------------------------------------------
replace_single_line \
	"ui/constants/endpoints.ts" \
	"^export const MESHERY_CLOUD_PROD = .*// $MARK.*\$" \
	"export const MESHERY_CLOUD_PROD = '$PRIMARY'; // $MARK"

# ---------------------------------------------------------------------------
# provider-ui/lib/data-fetch.js - return-to fallback URL
# ---------------------------------------------------------------------------
replace_single_line \
	"provider-ui/lib/data-fetch.js" \
	"^export const PROVIDER_URL = .*// $MARK.*\$" \
	"export const PROVIDER_URL = \"$PRIMARY\"; // $MARK"

# ---------------------------------------------------------------------------
# ui/tests/e2e/env.js - e2e default fallback
# ---------------------------------------------------------------------------
replace_single_line \
	"ui/tests/e2e/env.js" \
	"^const REMOTE_PROVIDER_URL = process.env.REMOTE_PROVIDER_URL \|\| .*// $MARK.*\$" \
	"const REMOTE_PROVIDER_URL = process.env.REMOTE_PROVIDER_URL || '$PRIMARY'; // $MARK"

# ---------------------------------------------------------------------------
# .github/workflows/build-ui-and-server.yml - e2e env var
# ---------------------------------------------------------------------------
replace_single_line \
	".github/workflows/build-ui-and-server.yml" \
	'^          REMOTE_PROVIDER_URL: ".*" # '"$MARK"'.*$' \
	"          REMOTE_PROVIDER_URL: \"$PRIMARY\" # $MARK"

# ---------------------------------------------------------------------------
# server/models/default_remote_providers.go - generated Go fallback
# ---------------------------------------------------------------------------
{
	echo "// Code generated by scripts/sync-provider-defaults.sh; DO NOT EDIT."
	echo "// $MARK"
	echo
	echo "package models"
	echo
	echo "// DefaultRemoteProviderURLs is the comma-separated canonical list of"
	echo "// default remote provider URLs. The server seeds this into viper via"
	echo "// SetDefault(\"PROVIDER_BASE_URLS\", ...) so operators who do not set the"
	echo "// env var still register the canonical providers."
	echo "const DefaultRemoteProviderURLs = \"$URLS\""
	echo
	echo "// PrimaryProviderURL is the single canonical provider host used by"
	echo "// single-URL consumers (SaaS deep-link constants, provider-ui return-to"
	echo "// fallback, e2e test defaults)."
	echo "const PrimaryProviderURL = \"$PRIMARY\""
} > "server/models/default_remote_providers.go"

# ---------------------------------------------------------------------------
# install/docker-extension/ui/src/components/utils/constants.js
# Multi-line REMOTE_PROVIDERS array. The block between the BEGIN/END markers
# is fully regenerated from PAIRS.
# ---------------------------------------------------------------------------
{
	for pair in "${PAIRS[@]}"; do
		name_var="${pair%%|*}"
		url_var="${pair##*|}"
		name=$(read_mk "$name_var")
		url=$(read_mk "$url_var")
		echo "    {"
		echo "        name: \"$name\","
		echo "        url: \"$url\","
		echo "    },"
	done
} | replace_block \
	"install/docker-extension/ui/src/components/utils/constants.js" \
	"// BEGIN $MARK" \
	"// END $MARK"

echo "sync-provider-defaults: ok ($(date '+%Y-%m-%d %H:%M:%S'))"
