#!/usr/bin/env python3
"""
Meshery API Schema Audit

Compares data sources within meshery/meshery and meshery/schemas to produce
a comprehensive audit of every API endpoint:

  1. server/router/server.go    → registered API endpoints
  2. docs/data/openapi.yml      → schema-backed check (from meshery/schemas)
  3. server/handlers/*.go       → schema-driven check (import analysis)

For each endpoint the script computes:
  - Coverage      — Overlap / Server Underlap / Schema Underlap
  - Status        — Active / Deprecated / Unimplemented / Cloud-only
  - Schema-Backed — Is the endpoint defined in the OpenAPI spec?
  - Schema Completeness — Full / Partial / Stub / N/A
  - Schema-Driven — Does the handler import+use meshery/schemas types?

Writes results to a Google Sheet. Credentials come from environment variables.

Intended to be run via Make targets from the repo root:

  make api-audit            # Dry-run audit — prints summary only
  make api-audit-update     # Audit and write results to Google Sheet
  make api-audit-setup      # Install Python dependencies only

For advanced usage or CI, the script can be invoked directly:

  python api_audit.py --repo . --dry-run --verbose   # Summary + per-endpoint details
  python api_audit.py --repo . --sheet-id $SHEET_ID  # Write to specific sheet
"""

# =====================================================================
# USER CONFIGURATION — edit these variables to match your environment
# =====================================================================

# Google Sheet ID to write audit results to.
# Find this in the sheet URL: https://docs.google.com/spreadsheets/d/<SHEET_ID>/edit
SHEET_ID = "1YYpMp0H1vlcdnWReaNfVzPcBHXAIlNh52tYBBaRBP7s"

# Absolute path to the meshery/meshery repository root.
# Leave empty to use the current working directory.
MESHERY_REPO_PATH = ""

# Path to Google service-account credentials JSON file (for local development).
# See: https://cloud.google.com/iam/docs/keys-create-delete
GOOGLE_APPLICATION_CREDENTIALS = ""

# Inline Google credentials JSON string (for CI / GitHub Actions).
# If set, this takes priority over GOOGLE_APPLICATION_CREDENTIALS.
GOOGLE_CREDENTIALS_JSON = ""

# Optional: path to an OpenAPI spec file (e.g. merged_openapi.yml from
# meshery/schemas). When set, overrides docs/data/openapi.yml in the repo.
OPENAPI_SPEC_PATH = ""

# =====================================================================

import argparse
import json
import os
import re
import subprocess
import sys
from collections import defaultdict
from datetime import datetime
from pathlib import Path
from typing import Any, Dict, List, Optional, Set, Tuple

try:
    import yaml
except ImportError:
    sys.exit("Missing dependency: pip install pyyaml")



# ---------------------------------------------------------------------------
# Paths relative to repo root
# ---------------------------------------------------------------------------
ROUTER_FILE = "server/router/server.go"
OPENAPI_FILE = "docs/data/openapi.yml"
HANDLERS_DIR = "server/handlers"
GO_MOD_FILE = "go.mod"

# ---------------------------------------------------------------------------
# Sheet configuration
# ---------------------------------------------------------------------------
SHEET_COLUMNS = [
    "Category",
    "Sub-Category",
    "Endpoints",
    "Methods",
    "Coverage",
    "Status",
    "Schema-Backed",
    "Schema Completeness",
    "Schema-Driven",
    "Notes",
    "Change Log",
]
COL_CATEGORY = 0
COL_SUBCATEGORY = 1
COL_ENDPOINTS = 2
COL_METHODS = 3
COL_COVERAGE = 4
COL_STATUS = 5
COL_BACKED = 6
COL_COMPLETENESS = 7
COL_DRIVEN = 8
COL_NOTES = 9
COL_CHANGELOG = 10

WORKSHEET_NAME = "Verification of Meshery Server API Endpoints"

HTTP_METHODS = frozenset({"get", "post", "put", "delete", "patch", "options", "head"})

# Methods that typically carry a request body
BODY_METHODS = frozenset({"POST", "PUT", "PATCH"})

# Sentinel value for handler response types that are []byte passthrough
# from the Provider interface (opaque JSON forwarded without typed struct).
_PROVIDER_BYTES_SENTINEL = "<provider-[]byte>"

MIDDLEWARE_NAMES = frozenset({
    "ProviderMiddleware", "AuthMiddleware", "SessionInjectorMiddleware",
    "KubernetesMiddleware", "K8sFSMMiddleware", "GraphqlMiddleware",
    "NoCacheMiddleware",
})

# Fallback category rules for router-only endpoints that have no OpenAPI tags.
# Spec-backed endpoints derive their category from x-tagGroups / tags instead.
CATEGORY_FALLBACK: List[Tuple[str, str, str]] = [
    ("/api/system/graphql", "Meshery Server and Components", "Meshery Operator"),
    ("/api/system/database", "Meshery Server and Components", "Database"),
    ("/api/system/adapter", "Meshery Server and Components", "Adapters"),
    ("/api/system/adapters", "Meshery Server and Components", "Adapters"),
    ("/api/system/availableAdapters", "Meshery Server and Components", "Adapters"),
    ("/api/system/meshsync", "Meshery Server and Components", "Meshsync"),
    ("/api/system", "Meshery Server and Components", "System"),
    ("/api/extension", "Meshery Server and Components", "System"),
    ("/api/meshmodels", "Capabilities Registry", "Entities"),
    ("/api/meshmodel", "Capabilities Registry", "Model Lifecycle"),
    ("/api/pattern", "Configuration", "Patterns"),
    ("/api/patterns", "Configuration", "Patterns"),
    ("/api/filter", "Configuration", "Filters"),
    ("/api/content/design", "Configuration", "Patterns"),
    ("/api/content/filter", "Configuration", "Filters"),
    ("/api/perf", "Benchmarking and Validation", "Performance (SMP)"),
    ("/api/mesh", "Benchmarking and Validation", "Performance (SMP)"),
    ("/api/smi", "Benchmarking and Validation", "Conformance (SMI)"),
    ("/api/user/performance", "Benchmarking and Validation", "Performance (SMP)"),
    ("/api/user/prefs/perf", "Benchmarking and Validation", "Performance (SMP)"),
    ("/api/telemetry/metrics/grafana", "Telemetry", "Grafana API"),
    ("/api/grafana", "Telemetry", "Grafana API"),
    ("/api/telemetry/metrics", "Telemetry", "Prometheus API"),
    ("/api/prometheus", "Telemetry", "Prometheus API"),
    ("/api/identity", "Identity", "User"),
    ("/api/user", "Identity", "User"),
    ("/api/token", "Identity", "User"),
    ("/api/provider", "Identity", "Providers, Extensions"),
    ("/api/providers", "Identity", "Providers, Extensions"),
    ("/api/schema", "Meshery Server and Components", "System"),
    ("/provider", "Identity", "Providers, Extensions"),
    ("/auth", "Identity", "User"),
    ("/user/login", "Identity", "User"),
    ("/user/logout", "Identity", "User"),
    ("/swagger.yaml", "Meshery Server and Components", "System"),
    ("/docs", "Meshery Server and Components", "System"),
    ("/healthz", "Meshery Server and Components", "System"),
    ("/error", "Other", ""),
]


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def normalize_path(path: str) -> str:
    """Replace {paramName} with positional {p1}, {p2}, ... for matching."""
    counter = [0]

    def _repl(_m):
        counter[0] += 1
        return f"{{p{counter[0]}}}"

    return re.sub(r"\{[^}]+\}", _repl, path)


def categorize(
    path: str,
    spec_categories: Optional[Dict[str, Tuple[str, str]]] = None,
) -> Tuple[str, str]:
    """Return (category, subcategory) for a given endpoint path.

    First checks *spec_categories* (derived from OpenAPI tags / x-tagGroups).
    Falls back to CATEGORY_FALLBACK prefix rules for router-only endpoints.
    """
    norm = normalize_path(path)
    if spec_categories and norm in spec_categories:
        return spec_categories[norm]
    for prefix, cat, sub in CATEGORY_FALLBACK:
        if path.startswith(prefix):
            return cat, sub
    return "Other", ""


def endpoint_sort_key(endpoint: Dict[str, Any]) -> Tuple[str, str, str, str]:
    """Return a deterministic sort key for sheet output."""
    return (
        endpoint["category"],
        endpoint["subcategory"],
        endpoint["path"],
        endpoint["methods"],
    )


# ---------------------------------------------------------------------------
# 1. Router parser — server/router/server.go
# ---------------------------------------------------------------------------

def parse_router(repo: Path) -> List[Dict[str, Any]]:
    """Parse route registrations from server.go."""
    router_file = repo / ROUTER_FILE
    if not router_file.exists():
        print(f"ERROR: {router_file} not found", file=sys.stderr)
        return []

    content = router_file.read_text(errors="replace")
    lines = content.splitlines()

    # Accumulate multi-line gMux statements
    statements: List[str] = []
    current = ""
    paren_depth = 0
    in_stmt = False
    current_commented = False

    for line in lines:
        stripped = line.strip()
        if re.match(r"^\s*(//\s*)?gMux\.(Handle|HandleFunc|PathPrefix)", line):
            if current and in_stmt:
                statements.append(current)
            current = stripped
            in_stmt = True
            current_commented = stripped.startswith("//")
            paren_depth = current.count("(") - current.count(")")
            if paren_depth <= 0 and not stripped.rstrip().endswith("."):
                statements.append(current)
                current, in_stmt, paren_depth, current_commented = "", False, 0, False
            continue
        if in_stmt:
            continuation = stripped
            if current_commented:
                continuation = re.sub(r"^//\s*", "", continuation)
            current += " " + continuation
            paren_depth += continuation.count("(") - continuation.count(")")
            if paren_depth <= 0 and not continuation.rstrip().endswith("."):
                statements.append(current)
                current, in_stmt, paren_depth, current_commented = "", False, 0, False

    if current:
        statements.append(current)

    routes = []
    for stmt in statements:
        route = _parse_route(stmt)
        if route:
            routes.append(route)
    return routes


def _parse_route(stmt: str) -> Optional[Dict[str, Any]]:
    """Parse a single gMux statement into a route dict."""
    commented = stmt.lstrip().startswith("//")
    clean = re.sub(r"^//\s*", "", stmt.strip()) if commented else stmt.strip()

    path_m = re.search(
        r'gMux\.(Handle|HandleFunc|PathPrefix)\s*\(\s*"([^"]+)"', clean
    )
    if not path_m:
        return None

    path = path_m.group(2)
    methods_m = re.search(r"\.\s*Methods\(\s*(.+?)\s*\)", clean)
    methods = re.findall(r'"([A-Z]+)"', methods_m.group(1)) if methods_m else ["ALL"]
    handler = _extract_handler(clean)

    return {
        "path": path,
        "methods": sorted(methods),
        "handler": handler,
        "commented": commented,
    }


def _extract_handler(line: str) -> str:
    """Extract handler function name from a route registration line."""
    # Exported methods on Handler receiver (h.FuncName)
    refs = re.findall(r"h\.([A-Z]\w+)", line)
    actual = [r for r in refs if r not in MIDDLEWARE_NAMES]
    if actual:
        return actual[-1]

    # Any h.funcName (including unexported)
    refs = re.findall(r"h\.([A-Za-z]\w+)", line)
    actual = [r for r in refs if r not in MIDDLEWARE_NAMES]
    if actual:
        return actual[-1]

    if "func(" in line or "func (" in line:
        return "<inline>"
    return "<unknown>"


# ---------------------------------------------------------------------------
# 2. OpenAPI parser — docs/data/openapi.yml
# ---------------------------------------------------------------------------

def _has_meaningful_schema(schema: Optional[dict]) -> bool:
    """Check if a schema object has real structure beyond a bare type."""
    if not schema or not isinstance(schema, dict):
        return False
    if "$ref" in schema:
        return True
    if any(k in schema for k in ("allOf", "oneOf", "anyOf")):
        return True
    if schema.get("type") == "array" and "items" in schema:
        return True
    if schema.get("type") == "object" and "properties" in schema:
        return True
    return False


def _get_content_schema(content: Any) -> Optional[dict]:
    """Extract schema from a content map (application/json, etc.)."""
    if not isinstance(content, dict):
        return None
    for _media_type, media_obj in content.items():
        if isinstance(media_obj, dict) and "schema" in media_obj:
            return media_obj["schema"]
    return None


def _describe_schema(schema: Optional[dict], label: str) -> List[str]:
    """Return human-readable findings about a schema object."""
    if not schema or not isinstance(schema, dict):
        return [f"{label}: no schema defined"]

    if "$ref" in schema:
        ref = schema["$ref"].rsplit("/", 1)[-1]
        return [f"{label}: references {ref}"]

    if any(k in schema for k in ("allOf", "oneOf", "anyOf")):
        combo = next(k for k in ("allOf", "oneOf", "anyOf") if k in schema)
        count = len(schema[combo]) if isinstance(schema[combo], list) else "?"
        return [f"{label}: {combo} with {count} sub-schemas"]

    s_type = schema.get("type", "untyped")

    if s_type == "array":
        items = schema.get("items", {})
        if isinstance(items, dict) and "$ref" in items:
            ref = items["$ref"].rsplit("/", 1)[-1]
            return [f"{label}: array of {ref}"]
        return [f"{label}: array (inline items)"]

    if s_type == "object":
        props = schema.get("properties", {})
        required = schema.get("required", [])
        if not props:
            return [f"{label}: object with no properties defined"]
        prop_names = sorted(props.keys())
        n = len(prop_names)
        preview = ", ".join(prop_names[:6])
        if n > 6:
            preview += f", ... ({n} total)"
        missing_desc = [
            k for k, v in props.items()
            if isinstance(v, dict) and not v.get("description")
        ]
        findings = [f"{label}: object with properties [{preview}]"]
        if not required:
            findings.append(f"{label}: no 'required' fields specified")
        if missing_desc:
            names = ", ".join(missing_desc[:5])
            if len(missing_desc) > 5:
                names += f", ... ({len(missing_desc)} total)"
            findings.append(
                f"{label}: properties missing description: {names}"
            )
        return findings

    # bare type (string, integer, boolean, etc.) with no structure
    return [f"{label}: bare type '{s_type}' with no properties or $ref"]


def _assess_completeness(operation: dict, method: str) -> Tuple[str, List[str]]:
    """Assess schema completeness for a single OpenAPI operation.

    Returns (completeness, detail_notes) where detail_notes lists every
    specific finding (missing fields, bare types, property gaps, etc.).
    """
    notes: List[str] = []
    expects_body = method.upper() in BODY_METHODS

    # --- Operation-level checks ---
    if not operation.get("operationId"):
        notes.append("missing operationId")
    if not operation.get("summary") and not operation.get("description"):
        notes.append("missing summary/description")

    # Parameters
    params = operation.get("parameters", [])
    if isinstance(params, list):
        no_desc = [
            p.get("name", "?") for p in params
            if isinstance(p, dict) and not p.get("description")
        ]
        if no_desc:
            notes.append(
                f"parameters missing description: {', '.join(no_desc[:5])}"
            )

    # --- Request side ---
    request_meaningful = False
    req_body = operation.get("requestBody", {})
    if isinstance(req_body, dict) and req_body:
        if "$ref" in req_body:
            request_meaningful = True
            ref = req_body["$ref"].rsplit("/", 1)[-1]
            notes.append(f"requestBody: references {ref}")
        else:
            req_content = req_body.get("content", {})
            req_schema = _get_content_schema(req_content)
            request_meaningful = _has_meaningful_schema(req_schema)
            notes.extend(_describe_schema(req_schema, "requestBody"))
    elif expects_body:
        notes.append("requestBody: not defined (method expects a body)")

    # --- Response side ---
    response_meaningful = False
    responses = operation.get("responses", {})
    defined_codes: Set[str] = set()

    if isinstance(responses, dict):
        for code, resp in responses.items():
            defined_codes.add(str(code))
            if not isinstance(resp, dict):
                continue

            if str(code).startswith("2"):
                if "$ref" in resp:
                    response_meaningful = True
                    ref = resp["$ref"].rsplit("/", 1)[-1]
                    notes.append(f"response {code}: references {ref}")
                else:
                    resp_content = resp.get("content", {})
                    resp_schema = _get_content_schema(resp_content)
                    if _has_meaningful_schema(resp_schema):
                        response_meaningful = True
                    notes.extend(
                        _describe_schema(resp_schema, f"response {code}")
                    )

    if not any(str(c).startswith("2") for c in defined_codes):
        notes.append("no 2xx success response defined")

    # Missing common error responses
    common_errors = {"400", "401", "404", "500"}
    missing_errors = sorted(common_errors - defined_codes)
    if missing_errors:
        notes.append(f"missing error responses: {', '.join(missing_errors)}")

    # --- Classify ---
    if expects_body:
        if request_meaningful and response_meaningful:
            return "Full", notes
        if request_meaningful or response_meaningful:
            return "Partial", notes
        return "Stub", notes
    else:
        # GET, DELETE, HEAD, OPTIONS — only response matters
        if response_meaningful:
            return "Full", notes
        return "Stub", notes


def _build_tag_category_map(doc: dict) -> Dict[str, Tuple[str, str]]:
    """Build a mapping from OpenAPI tag name to (category, subcategory).

    Uses x-tagGroups for the category (group name) and the tag's
    x-displayName for the subcategory.  Falls back to the raw tag name
    when x-displayName is absent.
    """
    # tag name → x-displayName
    tag_display: Dict[str, str] = {}
    for tag_def in doc.get("tags", []):
        if isinstance(tag_def, dict) and "name" in tag_def:
            display = tag_def.get("x-displayName", tag_def["name"])
            tag_display[tag_def["name"]] = display

    # tag name → (group_name, display_name)
    tag_to_category: Dict[str, Tuple[str, str]] = {}
    for group in doc.get("x-tagGroups", []):
        if not isinstance(group, dict):
            continue
        group_name = group.get("name", "Other")
        for tag_name in group.get("tags", []):
            display = tag_display.get(tag_name, tag_name)
            tag_to_category[tag_name] = (group_name, display)

    # Tags not listed in any group still get a category from their name
    for tag_name, display in tag_display.items():
        if tag_name not in tag_to_category:
            tag_to_category[tag_name] = (display, display)

    return tag_to_category


def parse_openapi(repo: Path, spec_override: Optional[Path] = None) -> dict:
    """Parse an OpenAPI spec and return structured data for all columns.

    When *spec_override* is given it is used directly; otherwise the script
    falls back to ``docs/data/openapi.yml`` inside *repo*.

    Returns a dict with:
      all_paths:        {norm_path: {METHOD, ...}}
      completeness:     {(norm_path, METHOD): "Full"/"Partial"/"Stub"}
      x_internal:       {(norm_path, METHOD): ["cloud"] or []}
      original_paths:   {norm_path: original_path}
      compl_notes:      {(norm_path, METHOD): [detail_notes]}
      path_categories:  {norm_path: (category, subcategory)}
      operations:       {(norm_path, METHOD): operation_dict}
    """
    spec_file = spec_override if spec_override else repo / OPENAPI_FILE
    empty = {
        "all_paths": {},
        "completeness": {},
        "x_internal": {},
        "original_paths": {},
        "compl_notes": {},
        "path_categories": {},
        "operations": {},
    }
    if not spec_file.exists():
        print(f"ERROR: {spec_file} not found", file=sys.stderr)
        return empty

    with open(spec_file, encoding="utf-8") as f:
        doc = yaml.safe_load(f)

    # Build tag → (category, subcategory) from x-tagGroups + tags
    tag_to_category = _build_tag_category_map(doc)

    all_paths: Dict[str, Set[str]] = {}
    completeness: Dict[Tuple[str, str], str] = {}
    x_internal: Dict[Tuple[str, str], List[str]] = {}
    original_paths: Dict[str, str] = {}
    compl_notes: Dict[Tuple[str, str], List[str]] = {}
    path_categories: Dict[str, Tuple[str, str]] = {}
    operations: Dict[Tuple[str, str], dict] = {}

    for path, methods_obj in doc.get("paths", {}).items():
        if not isinstance(methods_obj, dict):
            continue
        for method, details in methods_obj.items():
            if method.lower() not in HTTP_METHODS:
                continue
            if not isinstance(details, dict):
                continue

            norm = normalize_path(path)
            m_upper = method.upper()
            all_paths.setdefault(norm, set()).add(m_upper)

            # Track original path for spec-only endpoints
            if norm not in original_paths:
                original_paths[norm] = path

            # Derive category from the operation's first tag
            if norm not in path_categories:
                op_tags = details.get("tags", [])
                if isinstance(op_tags, list):
                    for tag_name in op_tags:
                        if tag_name in tag_to_category:
                            path_categories[norm] = tag_to_category[tag_name]
                            break

            # x-internal tag
            xi = details.get("x-internal", [])
            if not isinstance(xi, list):
                xi = [xi] if xi else []
            x_internal[(norm, m_upper)] = xi

            # Schema completeness (legacy — kept for fallback)
            comp, cnotes = _assess_completeness(details, method)
            completeness[(norm, m_upper)] = comp
            compl_notes[(norm, m_upper)] = cnotes

            # Store raw operation for cross-check
            operations[(norm, m_upper)] = details

    return {
        "all_paths": all_paths,
        "completeness": completeness,
        "x_internal": x_internal,
        "original_paths": original_paths,
        "compl_notes": compl_notes,
        "path_categories": path_categories,
        "operations": operations,
    }


# ---------------------------------------------------------------------------
# 3. Schema-driven detector — server/handlers/*.go
# ---------------------------------------------------------------------------

def _extract_function_body(text: str, func_name: str) -> Optional[str]:
    """Extract the body of a Go function using brace-depth counting.

    Skips braces inside string literals to avoid miscounting.
    """
    pat = re.compile(
        rf"func\s+(?:\([^)]*\)\s+)?{re.escape(func_name)}\s*\("
    )
    m = pat.search(text)
    if not m:
        return None

    brace_pos = text.find("{", m.end())
    if brace_pos == -1:
        return None

    depth = 1
    i = brace_pos + 1
    while i < len(text) and depth > 0:
        ch = text[i]
        if ch == "{":
            depth += 1
        elif ch == "}":
            depth -= 1
        elif ch == '"':
            i += 1
            while i < len(text) and text[i] != '"':
                if text[i] == "\\":
                    i += 1
                i += 1
        elif ch == "`":
            i += 1
            while i < len(text) and text[i] != "`":
                i += 1
        i += 1

    return text[brace_pos + 1 : i - 1] if depth == 0 else None


def build_schema_driven_map(repo: Path) -> Dict[str, Tuple[str, str]]:
    """Scan handler files for meshery/schemas imports at function level.

    For each handler function, extracts its body and checks which schema
    import aliases are actually used inside it — not just present in the file.

    Returns {handler_name: (status, reason)} where status is:
      TRUE    — function uses versioned schema types (models/v1beta1/*, etc.)
      Partial — function uses only models/core (utility types)
      FALSE   — function does not use any schema imports
    """
    handlers_dir = repo / HANDLERS_DIR
    if not handlers_dir.exists():
        print(f"WARNING: {handlers_dir} not found", file=sys.stderr)
        return {}

    # Read schema module path from go.mod
    schema_module = "github.com/meshery/schemas"
    go_mod = repo / GO_MOD_FILE
    if go_mod.exists():
        for line in go_mod.read_text().splitlines():
            m = re.match(r"\s*(github\.com/meshery/schemas)\s+v[\d.]+", line.strip())
            if m:
                schema_module = m.group(1)
                break

    escaped = re.escape(schema_module)
    alias_pat = re.compile(rf'(\w+)\s+"({escaped}[^"]*)"')
    bare_pat = re.compile(rf'"({escaped}[^"]*)"')

    # Per-file data
    handler_to_file: Dict[str, str] = {}
    file_texts: Dict[str, str] = {}
    file_aliases: Dict[str, Dict[str, str]] = {}  # fpath → {alias: import_path}

    for go_file in sorted(handlers_dir.glob("*.go")):
        if go_file.name.endswith("_test.go"):
            continue

        text = go_file.read_text(errors="replace")
        fpath = str(go_file)
        file_texts[fpath] = text

        # Map handler names → file
        for name in re.findall(
            r"func\s+\([^)]*\*?Handler[^)]*\)\s+(\w+)\s*\(", text
        ):
            handler_to_file[name] = fpath
        for name in re.findall(r"^func\s+(\w+)\s*\(", text, re.MULTILINE):
            if name not in handler_to_file:
                handler_to_file[name] = fpath

        # Build alias map: alias → full import path
        aliases: Dict[str, str] = {}
        for alias, imp_path in alias_pat.findall(text):
            aliases[alias] = imp_path
        seen_paths = set(aliases.values())
        for imp_path in bare_pat.findall(text):
            if imp_path not in seen_paths:
                last_seg = imp_path.rstrip("/").rsplit("/", 1)[-1]
                aliases[last_seg] = imp_path
                seen_paths.add(imp_path)
        file_aliases[fpath] = aliases

    # Classify each handler at function level
    result: Dict[str, Tuple[str, str]] = {}
    for name, fpath in handler_to_file.items():
        aliases = file_aliases.get(fpath, {})
        text = file_texts.get(fpath, "")

        # No schema imports in this file at all → fast path
        if not aliases:
            result[name] = ("FALSE", "no schema imports")
            continue

        # Try function-level analysis
        func_body = _extract_function_body(text, name)
        if func_body is not None:
            used: Set[str] = set()
            for alias, imp_path in aliases.items():
                if re.search(rf"\b{re.escape(alias)}\.", func_body):
                    used.add(imp_path)

            if used:
                versioned = {p for p in used if re.search(r"models/v\d+", p)}
                core_only = {p for p in used if "models/core" in p}
                if versioned:
                    pkgs = ", ".join(
                        sorted(p.replace(schema_module + "/", "") for p in versioned)
                    )
                    result[name] = ("TRUE", f"imports: {pkgs}")
                elif core_only:
                    result[name] = ("Partial", "imports: models/core only")
                else:
                    result[name] = ("FALSE", "schema dep but no model types")
            else:
                result[name] = ("FALSE", "no schema usage in function body")
        else:
            # Couldn't extract body — fall back to file-level
            all_imports = set(aliases.values())
            versioned = {p for p in all_imports if re.search(r"models/v\d+", p)}
            core_only = {p for p in all_imports if "models/core" in p}
            if versioned:
                pkgs = ", ".join(
                    sorted(p.replace(schema_module + "/", "") for p in versioned)
                )
                result[name] = ("TRUE", f"imports: {pkgs} (file-level)")
            elif core_only:
                result[name] = ("Partial", "imports: models/core only (file-level)")
            else:
                result[name] = ("FALSE", "schema dep but no model types")

    return result


# ---------------------------------------------------------------------------
# 4. Handler I/O Type Extractor
# ---------------------------------------------------------------------------

def _find_var_type(body: str, var_name: str) -> Optional[str]:
    """Find the Go type of a variable from its declaration in a function body.

    Handles:
      var x *pkg.Type             → pkg.Type
      var x pkg.Type              → pkg.Type
      x := &pkg.Type{}           → pkg.Type
      x := pkg.Type{field: val}  → pkg.Type
      x = &pkg.Type{}            → pkg.Type
      x := make([]pkg.Type, 0)   → []pkg.Type
    """
    esc = re.escape(var_name)

    # var x *Type  /  var x Type
    m = re.search(rf"\bvar\s+{esc}\s+\*?([\[\]]*[A-Za-z_][\w.]*)", body)
    if m:
        return m.group(1)

    # x := &Type{...}  /  x := Type{...}  /  x = &Type{...}
    # Must not match keywords (e.g. "if", "for", "return") or function calls
    m = re.search(
        rf"\b{esc}\s*:?=\s*&?([\[\]]*[A-Za-z_][\w.]*)\s*\{{", body
    )
    if m:
        candidate = m.group(1)
        # Filter out Go keywords and common non-type tokens
        if candidate not in (
            "if", "for", "range", "return", "func", "switch", "select",
            "map", "struct", "interface", "err", "nil",
        ):
            return candidate

    # x := make([]Type, ...)
    m = re.search(rf"\b{esc}\s*:?=\s*make\(([\[\]]+[A-Za-z_][\w.]*)", body)
    if m:
        return m.group(1)

    return None


def _build_provider_return_types(repo: Path) -> Dict[str, str]:
    """Parse the Provider interface to build {MethodName: return_type}.

    The Provider interface in server/models/providers.go defines the
    return types for all provider methods.  Many return ([]byte, error)
    which is opaque, but others return concrete typed values like
    (*connections.ConnectionPage, error).

    Returns {"GetConnections": "*connections.ConnectionPage", ...}.
    Only includes methods that return a typed (non-[]byte) first value.
    """
    providers_file = repo / "server" / "models" / "providers.go"
    if not providers_file.exists():
        return {}

    text = providers_file.read_text(errors="replace")

    # Find the Provider interface block
    m = re.search(r"type\s+Provider\s+interface\s*\{", text)
    if not m:
        return {}

    # Walk braces to extract interface body
    start = m.end()
    depth = 1
    i = start
    while i < len(text) and depth > 0:
        if text[i] == "{":
            depth += 1
        elif text[i] == "}":
            depth -= 1
        i += 1
    iface_body = text[start : i - 1]

    result: Dict[str, str] = {}

    # Match method signatures:  MethodName(...) (RetType, error)
    # or MethodName(...) (RetType, int, error) etc.
    for m in re.finditer(
        r"^\s*(\w+)\s*\([^)]*\)\s*\(([^)]+)\)",
        iface_body,
        re.MULTILINE,
    ):
        method_name = m.group(1)
        returns = m.group(2).strip()

        # Parse return types — split by comma
        ret_parts = [r.strip() for r in returns.split(",")]
        if len(ret_parts) < 2:
            continue

        first_ret = ret_parts[0]

        # Skip []byte returns — opaque passthrough
        if first_ret == "[]byte":
            continue
        # Skip string, int, bool, error returns — not typed structs
        if first_ret in ("string", "int", "bool", "error"):
            continue

        result[method_name] = first_ret

    return result


def _resolve_resp_type(
    body: str,
    resp_var: str,
    provider_types: Dict[str, str],
) -> Optional[str]:
    """Try to resolve the type of a response variable.

    Checks in order:
      1. Local type declaration (var, :=, make)
      2. Provider method return type (x, err := provider.Method(...))
      3. Config/field method return type (h.config.Client.Method(...))
    """
    # 1. Local declaration
    t = _find_var_type(body, resp_var)
    if t:
        return t

    esc = re.escape(resp_var)

    # 2. x, err := provider.MethodName(...)
    m = re.search(
        rf"\b{esc}\s*(?:,\s*\w+)*\s*:?=\s*(?:\w+\.)*provider\.(\w+)\s*\(",
        body,
    )
    if m:
        method = m.group(1)
        ptype = provider_types.get(method)
        if ptype:
            return ptype

    # 3. x, err := provider.MethodName(...)  (provider as p, prov, etc.)
    m = re.search(
        rf"\b{esc}\s*(?:,\s*\w+)*\s*:?=\s*(?:p|prov|provider)\.(\w+)\s*\(",
        body,
    )
    if m:
        method = m.group(1)
        ptype = provider_types.get(method)
        if ptype:
            return ptype

    # 4. x, err := h.config.SomeClient.Method(...) — can't resolve generically
    #    but at least return None cleanly

    return None


def extract_handler_io_types(
    repo: Path,
) -> Dict[str, Dict[str, Optional[str]]]:
    """For each handler function, extract request-body and response Go types.

    Scans server/handlers/*.go for json.Decode / json.Unmarshal (request)
    and json.Encode / json.Marshal (response) calls, then resolves the
    variable's declared type using local declarations and Provider interface
    return types.

    Returns {handler_name: {"request_type": "pkg.Type" | None,
                            "response_type": "pkg.Type" | None}}.
    """
    handlers_dir = repo / HANDLERS_DIR
    if not handlers_dir.exists():
        return {}

    provider_types = _build_provider_return_types(repo)

    result: Dict[str, Dict[str, Optional[str]]] = {}

    for go_file in sorted(handlers_dir.glob("*.go")):
        if go_file.name.endswith("_test.go"):
            continue

        text = go_file.read_text(errors="replace")

        # Collect handler function names (methods on *Handler receiver)
        handler_names = re.findall(
            r"func\s+\([^)]*\*?Handler[^)]*\)\s+(\w+)\s*\(", text
        )

        for name in handler_names:
            func_body = _extract_function_body(text, name)
            if not func_body:
                continue

            io: Dict[str, Optional[str]] = {
                "request_type": None,
                "response_type": None,
            }

            # --- Request type ---
            # json.NewDecoder(r.Body).Decode(&var)
            req_vars = re.findall(
                r"json\.NewDecoder\(\w+\.Body\)\.Decode\(&(\w+)\)",
                func_body,
            )
            # json.Unmarshal(bytes, &var)  (only if preceded by io.ReadAll on Body)
            if not req_vars and re.search(
                r"io\.ReadAll\(\w+\.Body\)", func_body
            ):
                req_vars = re.findall(
                    r"json\.Unmarshal\(\w+,\s*&(\w+)\)", func_body
                )

            if req_vars:
                io["request_type"] = _find_var_type(func_body, req_vars[0])
                # Fallback: check function parameters for request var type
                if not io["request_type"]:
                    sig_m = re.search(
                        rf"func\s+\([^)]*\)\s+{re.escape(name)}"
                        rf"\s*\(([^)]*)\)",
                        text,
                    )
                    if sig_m:
                        params = sig_m.group(1)
                        pm = re.search(
                            rf"\b{re.escape(req_vars[0])}\s+"
                            rf"\*?([\[\]]*[A-Za-z_][\w.]*)",
                            params,
                        )
                        if pm:
                            io["request_type"] = pm.group(1)

            # --- Response type ---
            # json.NewEncoder(w).Encode(expr)  (chained)
            enc_matches = re.findall(
                r"json\.NewEncoder\(\w+\)\.Encode\(([^)]+)\)", func_body
            )
            # enc := json.NewEncoder(w) ... enc.Encode(expr) (stored encoder)
            if not enc_matches:
                enc_var_m = re.search(
                    r"(\w+)\s*:?=\s*json\.NewEncoder\(\w+\)", func_body
                )
                if enc_var_m:
                    enc_var = enc_var_m.group(1)
                    enc_matches = re.findall(
                        rf"{re.escape(enc_var)}\.Encode\(([^)]+)\)",
                        func_body,
                    )
            # json.Marshal(expr)
            marsh_matches = re.findall(
                r"json\.Marshal\(([^)]+)\)", func_body
            )

            resp_expr = None
            if enc_matches:
                resp_expr = enc_matches[0].strip().lstrip("&*")
            elif marsh_matches:
                resp_expr = marsh_matches[0].strip().lstrip("&*")

            if resp_expr:
                # Direct struct literal: pkg.Type{...}
                m = re.match(r"([\[\]]*[A-Za-z_][\w.]*)\s*\{", resp_expr)
                if m:
                    candidate = m.group(1)
                    if candidate not in (
                        "if", "for", "range", "return", "func",
                        "map", "struct", "interface",
                    ):
                        io["response_type"] = candidate
                # Package-level variable: pkg.VarName (no parens/braces)
                elif re.match(r"[A-Za-z_]\w*\.[A-Z]\w*$", resp_expr):
                    io["response_type"] = resp_expr
                # Simple variable name — trace via local decls + provider types
                elif re.match(r"[A-Za-z_]\w*$", resp_expr):
                    io["response_type"] = _resolve_resp_type(
                        func_body, resp_expr, provider_types
                    )
                    # Fallback: check function parameters
                    if not io["response_type"]:
                        sig_m = re.search(
                            rf"func\s+\([^)]*\)\s+{re.escape(name)}"
                            rf"\s*\(([^)]*)\)",
                            text,
                        )
                        if sig_m:
                            params = sig_m.group(1)
                            pm = re.search(
                                rf"\b{re.escape(resp_expr)}\s+"
                                rf"\*?([\[\]]*[A-Za-z_][\w.]*)",
                                params,
                            )
                            if pm:
                                io["response_type"] = pm.group(1)

            # --- Fallback: fmt.Fprint with provider response ---
            # Pattern: resp, err := provider.X(...) then fmt.Fprint(w, string(resp))
            if not io["response_type"]:
                fprint_m = re.search(
                    r"fmt\.Fprint\w*\(\w+,\s*string\((\w+)\)\)", func_body
                )
                if fprint_m:
                    resolved = _resolve_resp_type(
                        func_body, fprint_m.group(1), provider_types
                    )
                    if resolved:
                        io["response_type"] = resolved
                    else:
                        # Provider returns []byte — opaque passthrough
                        io["response_type"] = _PROVIDER_BYTES_SENTINEL

            result[name] = io

    return result


# ---------------------------------------------------------------------------
# 5. Go Struct Field Parser — json tags → field name set
# ---------------------------------------------------------------------------

def _gomodcache() -> Optional[Path]:
    """Return GOMODCACHE path, or None."""
    try:
        out = subprocess.check_output(
            ["go", "env", "GOMODCACHE"],
            text=True,
            stderr=subprocess.DEVNULL,
        ).strip()
        if out:
            return Path(out)
    except (FileNotFoundError, subprocess.CalledProcessError):
        pass
    fallback = Path.home() / "go" / "pkg" / "mod"
    return fallback if fallback.exists() else None


def _parse_struct_json_fields(text: str, type_name: str) -> Optional[Set[str]]:
    """Extract the set of JSON field names from a Go struct definition.

    Looks for ``type <type_name> struct { ... }`` and pulls out every
    ``json:"<name>,..."`` tag.  Returns None if the struct is not found.
    """
    # Match: type TypeName struct {
    pat = re.compile(
        rf"\btype\s+{re.escape(type_name)}\s+struct\s*\{{", re.MULTILINE
    )
    m = pat.search(text)
    if not m:
        return None

    # Walk braces to find end of struct
    start = m.end()
    depth = 1
    i = start
    while i < len(text) and depth > 0:
        if text[i] == "{":
            depth += 1
        elif text[i] == "}":
            depth -= 1
        i += 1

    struct_body = text[start : i - 1]
    fields: Set[str] = set()
    for tag_m in re.finditer(r'json:"([^"]+)"', struct_body):
        json_name = tag_m.group(1).split(",")[0]
        if json_name and json_name != "-":
            fields.add(json_name)
    return fields if fields else None


def _resolve_type_alias(text: str, type_name: str) -> Optional[str]:
    """If ``type TypeName = pkg.Actual``, return the aliased type."""
    m = re.search(
        rf"\btype\s+{re.escape(type_name)}\s*=\s*([\w./]+\.(\w+))",
        text,
    )
    if m:
        return m.group(1)  # e.g. schemasEnvironment.EnvironmentPayload
    return None


def build_go_struct_fields(
    repo: Path,
) -> Dict[str, Set[str]]:
    """Build a map of Go type → {json field names} for types used in handlers.

    Searches:
      1. server/handlers/*.go  (inline struct types)
      2. server/models/**/*.go (models package + sub-packages)
      3. GOMODCACHE/github.com/meshery/schemas/... (generated schema types)

    Handles one level of type-alias indirection (common in meshery — e.g.
    ``type EnvironmentPayload = schemasEnv.EnvironmentPayload``).

    Returns {"TypeName": {"field_a", "field_b", ...}, ...}.
    The key is the short type name (without package qualifier).
    """
    result: Dict[str, Set[str]] = {}
    alias_map: Dict[str, str] = {}  # short_name → aliased_full_import

    # --- Scan local Go files ---
    search_dirs = [
        repo / HANDLERS_DIR,
        repo / "server" / "models",
    ]

    for search_dir in search_dirs:
        if not search_dir.exists():
            continue
        for go_file in sorted(search_dir.rglob("*.go")):
            if go_file.name.endswith("_test.go"):
                continue
            text = go_file.read_text(errors="replace")

            # Find all struct definitions
            for name in re.findall(
                r"\btype\s+(\w+)\s+struct\s*\{", text
            ):
                fields = _parse_struct_json_fields(text, name)
                if fields:
                    result[name] = fields

            # Find type aliases  (type X = pkg.Y)
            for m in re.finditer(
                r"\btype\s+(\w+)\s*=\s*(\w+)\.(\w+)", text
            ):
                local_name, _alias_pkg, remote_name = (
                    m.group(1),
                    m.group(2),
                    m.group(3),
                )
                alias_map[local_name] = remote_name
                # Also try to find the import path for the alias package
                imp_match = re.search(
                    rf'{re.escape(_alias_pkg)}\s+"([^"]+)"', text
                )
                if imp_match:
                    alias_map[f"_imp_{local_name}"] = imp_match.group(1)

    # --- Resolve aliases from GOMODCACHE ---
    modcache = _gomodcache()
    schema_version = None
    go_mod = repo / GO_MOD_FILE
    if go_mod.exists():
        for line in go_mod.read_text().splitlines():
            vm = re.match(
                r"\s*github\.com/meshery/schemas\s+(v[\d.]+)", line.strip()
            )
            if vm:
                schema_version = vm.group(1)
                break

    if modcache and schema_version:
        schema_root = (
            modcache
            / f"github.com/meshery/schemas@{schema_version}"
        )
        if schema_root.exists():
            for go_file in sorted(schema_root.rglob("*.go")):
                if go_file.name.endswith("_test.go"):
                    continue
                text = go_file.read_text(errors="replace")
                for name in re.findall(
                    r"\btype\s+(\w+)\s+struct\s*\{", text
                ):
                    fields = _parse_struct_json_fields(text, name)
                    if fields and name not in result:
                        result[name] = fields

    # --- Resolve aliases: if local name not found, use remote name ---
    for local_name, remote_name in alias_map.items():
        if local_name.startswith("_imp_"):
            continue
        if local_name not in result and remote_name in result:
            result[local_name] = result[remote_name]

    return result


# ---------------------------------------------------------------------------
# 6. Spec Schema Field Extractor
# ---------------------------------------------------------------------------

def _collect_property_names(schema: Any) -> Set[str]:
    """Recursively collect property names from an OpenAPI schema."""
    if not isinstance(schema, dict):
        return set()

    props = set()

    if "properties" in schema and isinstance(schema["properties"], dict):
        props.update(schema["properties"].keys())

    # Walk allOf / oneOf / anyOf
    for combo in ("allOf", "oneOf", "anyOf"):
        if combo in schema and isinstance(schema[combo], list):
            for sub in schema[combo]:
                props.update(_collect_property_names(sub))

    # Array items
    if schema.get("type") == "array" and isinstance(
        schema.get("items"), dict
    ):
        props.update(_collect_property_names(schema["items"]))

    return props


def extract_spec_schema_fields(
    operation: dict, method: str
) -> Dict[str, Set[str]]:
    """Extract property name sets from an OpenAPI operation.

    Returns {"request_fields": set, "response_fields": set}.
    """
    req_fields: Set[str] = set()
    resp_fields: Set[str] = set()

    # --- Request body ---
    req_body = operation.get("requestBody", {})
    if isinstance(req_body, dict) and req_body:
        content = req_body.get("content", {})
        if isinstance(content, dict):
            for _mt, media_obj in content.items():
                if isinstance(media_obj, dict) and "schema" in media_obj:
                    req_fields = _collect_property_names(media_obj["schema"])
                    break

    # --- Response (first 2xx) ---
    responses = operation.get("responses", {})
    if isinstance(responses, dict):
        for code, resp in responses.items():
            if not str(code).startswith("2") or not isinstance(resp, dict):
                continue
            content = resp.get("content", {})
            if isinstance(content, dict):
                for _mt, media_obj in content.items():
                    if isinstance(media_obj, dict) and "schema" in media_obj:
                        resp_fields = _collect_property_names(
                            media_obj["schema"]
                        )
                        break
            break  # use first 2xx only

    return {"request_fields": req_fields, "response_fields": resp_fields}


# ---------------------------------------------------------------------------
# 7. Cross-Check Completeness
# ---------------------------------------------------------------------------

def _field_overlap(
    go_fields: Optional[Set[str]], spec_fields: Set[str]
) -> Tuple[Optional[float], Set[str], Set[str], Set[str]]:
    """Compute overlap between Go struct fields and spec properties.

    Returns (ratio_or_None, common, only_go, only_spec).
    """
    if go_fields is None or not spec_fields or not go_fields:
        return None, set(), set(), set()

    common = go_fields & spec_fields
    only_go = go_fields - spec_fields
    only_spec = spec_fields - go_fields
    ratio = len(common) / max(len(go_fields), len(spec_fields))
    return ratio, common, only_go, only_spec


def _format_field_set(fields: Set[str], limit: int = 6) -> str:
    """Format a set of field names as a compact string."""
    if not fields:
        return ""
    items = sorted(fields)
    preview = ", ".join(items[:limit])
    if len(items) > limit:
        preview += f", ... ({len(items)} total)"
    return preview


def cross_check_completeness(
    handler_name: str,
    handler_io: Dict[str, Optional[str]],
    go_fields_map: Dict[str, Set[str]],
    spec_fields: Dict[str, Set[str]],
    method: str,
) -> Tuple[str, List[str]]:
    """Cross-check handler's Go types against spec schemas.

    Returns (completeness, structured_notes) where structured_notes is a
    list of lines with clear section markers for later formatting.

    Lines are prefixed with markers:
      [REQ]  — request-side finding
      [RESP] — response-side finding
      [INFO] — general info
    """
    notes: List[str] = []
    expects_body = method.upper() in BODY_METHODS

    req_type = handler_io.get("request_type")
    resp_type = handler_io.get("response_type")
    is_resp_passthrough = resp_type == _PROVIDER_BYTES_SENTINEL

    # Strip pointer/slice prefix for type lookup
    def _bare_type(t: Optional[str]) -> Optional[str]:
        if t is None or t == _PROVIDER_BYTES_SENTINEL:
            return None
        return t.lstrip("*[]").split(".")[-1]

    req_type_short = _bare_type(req_type)
    resp_type_short = _bare_type(resp_type)

    req_go_fields = (
        go_fields_map.get(req_type_short) if req_type_short else None
    )
    resp_go_fields = (
        go_fields_map.get(resp_type_short) if resp_type_short else None
    )

    spec_req = spec_fields.get("request_fields", set())
    spec_resp = spec_fields.get("response_fields", set())
    has_spec = bool(spec_req or spec_resp)

    if not has_spec:
        return "N/A", ["[INFO] No spec schema to compare against"]

    # --- Request side ---
    req_ratio, req_common, req_only_go, req_only_spec = _field_overlap(
        req_go_fields, spec_req
    )

    if expects_body:
        if req_type:
            notes.append(f"[REQ] Handler type: {req_type}")
            if req_ratio is not None:
                pct = int(req_ratio * 100)
                notes.append(
                    f"[REQ] Field match: {pct}% "
                    f"({len(req_common)}/{max(len(req_go_fields or set()), len(spec_req))})"
                )
                if req_common:
                    notes.append(
                        f"[REQ] Matching: {_format_field_set(req_common)}"
                    )
                if req_only_go:
                    notes.append(
                        f"[REQ] In handler only: "
                        f"{_format_field_set(req_only_go, 5)}"
                    )
                if req_only_spec:
                    notes.append(
                        f"[REQ] In spec only: "
                        f"{_format_field_set(req_only_spec, 5)}"
                    )
            else:
                notes.append(
                    "[REQ] Struct fields not found for comparison"
                )
        else:
            notes.append("[REQ] Could not extract request type from handler")

    # --- Response side ---
    resp_ratio, resp_common, resp_only_go, resp_only_spec = _field_overlap(
        resp_go_fields, spec_resp
    )

    if is_resp_passthrough:
        notes.append(
            "[RESP] Provider returns raw []byte — "
            "response type is opaque, cannot cross-check"
        )
    elif resp_type:
        notes.append(f"[RESP] Handler type: {resp_type}")
        if resp_ratio is not None:
            pct = int(resp_ratio * 100)
            notes.append(
                f"[RESP] Field match: {pct}% "
                f"({len(resp_common)}/{max(len(resp_go_fields or set()), len(spec_resp))})"
            )
            if resp_common:
                notes.append(
                    f"[RESP] Matching: {_format_field_set(resp_common)}"
                )
            if resp_only_go:
                notes.append(
                    f"[RESP] In handler only: "
                    f"{_format_field_set(resp_only_go, 5)}"
                )
            if resp_only_spec:
                notes.append(
                    f"[RESP] In spec only: "
                    f"{_format_field_set(resp_only_spec, 5)}"
                )
        else:
            notes.append(
                "[RESP] Struct fields not found for comparison"
            )
    else:
        if spec_resp:
            notes.append(
                "[RESP] Could not extract response type from handler"
            )
        else:
            notes.append("[RESP] No response body in spec (expected)")

    # --- Classify ---
    GOOD_THRESHOLD = 0.70

    if expects_body:
        req_good = req_ratio is not None and req_ratio >= GOOD_THRESHOLD
        resp_good = resp_ratio is not None and resp_ratio >= GOOD_THRESHOLD
        req_any = req_ratio is not None and req_ratio > 0
        resp_any = resp_ratio is not None and resp_ratio > 0

        if req_ratio is None and resp_ratio is None:
            if req_type or resp_type:
                return "Stub", notes
            return "N/A", notes

        if req_good and resp_good:
            return "Full", notes
        if req_good or resp_good or (req_any and resp_any):
            return "Partial", notes
        return "Stub", notes
    else:
        if resp_ratio is None:
            if is_resp_passthrough or (resp_type and resp_type != _PROVIDER_BYTES_SENTINEL):
                return "Stub", notes
            if spec_resp:
                return "N/A", notes
            return "N/A", notes

        if resp_ratio >= GOOD_THRESHOLD:
            return "Full", notes
        if resp_ratio > 0:
            return "Partial", notes
        return "Stub", notes


# ---------------------------------------------------------------------------
# 8. Actionable Notes Builder
# ---------------------------------------------------------------------------

def _build_actionable_notes(
    *,
    coverage: str,
    status: str,
    is_commented: bool,
    completeness: str,
    compl_notes: List[str],
    driven: str,
    handler: str,
    cloud_methods: List[str],
    spec_methods: Set[str],
) -> str:
    """Build a structured, readable summary for the Notes column.

    Organises findings into clearly labelled sections separated by
    newlines so that the sheet cell (or verbose CLI output) is scannable.

    Sections (only included when relevant):
      COVERAGE / STATUS  — high-level action for the endpoint
      CLOUD              — x-internal annotation info
      CROSS-CHECK        — handler ↔ spec field comparison (request, response)
      SPEC QUALITY       — legacy completeness notes (fallback path)
      SCHEMA-DRIVEN      — whether handler imports meshery/schemas types
    """
    sections: List[str] = []

    # ── COVERAGE / STATUS ──────────────────────────────────────────────
    if is_commented:
        sections.append(
            "[COVERAGE] Route is commented out — "
            "consider removal from router and spec"
        )

    if coverage == "Server Underlap":
        sections.append(
            "[COVERAGE] Not in OpenAPI spec — add spec definition"
        )
    elif coverage == "Schema Underlap":
        if status == "Unimplemented":
            sections.append(
                "[COVERAGE] In spec but no server route — "
                "implement handler or remove from spec"
            )

    # ── CLOUD ──────────────────────────────────────────────────────────
    if coverage == "Overlap" and cloud_methods:
        if len(cloud_methods) == len(spec_methods):
            sections.append(
                "[CLOUD] All methods marked x-internal: cloud in spec, "
                "but an equivalent route exists in the server"
            )
        else:
            sections.append(
                f"[CLOUD] Partially cloud-annotated: "
                f"{', '.join(cloud_methods)} marked x-internal in spec"
            )

    # ── CROSS-CHECK ────────────────────────────────────────────────────
    # compl_notes from cross_check_completeness carry [REQ]/[RESP]/[INFO]
    # prefixes.  Group them for readability.
    req_lines = [n for n in compl_notes if n.startswith("[REQ]")]
    resp_lines = [n for n in compl_notes if n.startswith("[RESP]")]
    info_lines = [n for n in compl_notes if n.startswith("[INFO]")]
    legacy_lines = [
        n for n in compl_notes
        if not n.startswith(("[REQ]", "[RESP]", "[INFO]"))
    ]

    if req_lines or resp_lines or info_lines:
        cross_parts: List[str] = []
        if info_lines:
            for line in info_lines:
                cross_parts.append(line.replace("[INFO] ", ""))
        if req_lines:
            cross_parts.append("Request:")
            for line in req_lines:
                cross_parts.append("  " + line.replace("[REQ] ", ""))
        if resp_lines:
            cross_parts.append("Response:")
            for line in resp_lines:
                cross_parts.append("  " + line.replace("[RESP] ", ""))
        sections.append(
            "[CROSS-CHECK]\n" + "\n".join(cross_parts)
        )

    # ── SPEC QUALITY (legacy fallback) ─────────────────────────────────
    if legacy_lines:
        sections.append(
            "[SPEC QUALITY]\n" + "\n".join(legacy_lines)
        )

    # ── SCHEMA-DRIVEN ──────────────────────────────────────────────────
    if driven == "FALSE" and coverage != "Schema Underlap":
        if handler in ("<inline>", "<unknown>"):
            sections.append(
                f"[SCHEMA-DRIVEN] Handler is {handler} — "
                "extract to a named function and adopt schema types"
            )
        else:
            sections.append(
                "[SCHEMA-DRIVEN] Handler does not import meshery/schemas "
                "types — migrate to schema-driven"
            )
    elif driven == "Partial":
        sections.append(
            "[SCHEMA-DRIVEN] Uses only core schema types — "
            "adopt versioned model types (v1beta1, etc.)"
        )

    return "\n\n".join(sections) if sections else ""


# ---------------------------------------------------------------------------
# 5. Classification — bidirectional walk (Router ∪ Spec)
# ---------------------------------------------------------------------------

def _aggregate_completeness(
    norm: str,
    methods: List[str],
    spec_data: dict,
) -> Tuple[str, List[str]]:
    """Aggregate completeness across methods for a single endpoint."""
    comp_map = spec_data["completeness"]
    cnotes_map = spec_data["compl_notes"]
    all_paths = spec_data["all_paths"]
    spec_methods = all_paths.get(norm, set())

    if not spec_methods:
        return "N/A", []

    method_comps = []
    agg_notes: List[str] = []
    check_methods = spec_methods if methods == ["ALL"] else [
        m for m in methods if m in spec_methods
    ]

    for m in check_methods:
        c = comp_map.get((norm, m), "Stub")
        method_comps.append(c)
        agg_notes.extend(cnotes_map.get((norm, m), []))

    if not method_comps:
        return "Stub", ["spec path exists but no method match"]

    # Deduplicate notes
    seen: Set[str] = set()
    unique_notes = []
    for n in agg_notes:
        if n not in seen:
            seen.add(n)
            unique_notes.append(n)

    if all(c == "Full" for c in method_comps):
        return "Full", unique_notes
    if any(c == "Full" for c in method_comps) or any(c == "Partial" for c in method_comps):
        return "Partial", unique_notes
    return "Stub", unique_notes


def classify_endpoints(
    routes: List[Dict[str, Any]],
    spec_data: dict,
    schema_map: Dict[str, Tuple[str, str]],
    handler_io_map: Optional[Dict[str, Dict[str, Optional[str]]]] = None,
    go_fields_map: Optional[Dict[str, Set[str]]] = None,
) -> List[Dict[str, Any]]:
    """Classify endpoints from both router and spec (bidirectional walk).

    When *handler_io_map* and *go_fields_map* are provided, schema
    completeness is determined by cross-checking handler Go types against
    spec schemas (field-level comparison).  Otherwise falls back to the
    legacy structural assessment.
    """
    all_paths = spec_data["all_paths"]
    x_internal_map = spec_data["x_internal"]
    original_paths = spec_data["original_paths"]
    path_categories = spec_data.get("path_categories", {})
    operations_map = spec_data.get("operations", {})

    endpoints: List[Dict[str, Any]] = []
    router_norm_paths: Set[str] = set()

    # ------------------------------------------------------------------
    # Pass 1: Router-sourced endpoints
    # ------------------------------------------------------------------
    grouped_routes: Dict[Tuple[str, str], List[Dict[str, Any]]] = defaultdict(list)
    for route in routes:
        methods_str = ", ".join(route["methods"])
        grouped_routes[(route["path"], methods_str)].append(route)

    for path, methods_str in sorted(grouped_routes):
        route_group = grouped_routes[(path, methods_str)]
        route = next((r for r in route_group if not r["commented"]), route_group[0])
        methods = route["methods"]
        is_commented = all(r["commented"] for r in route_group)

        norm = normalize_path(path)
        category, subcategory = categorize(path, path_categories)
        router_norm_paths.add(norm)

        spec_methods = all_paths.get(norm, set())

        # --- Coverage ---
        coverage = "Overlap" if spec_methods else "Server Underlap"

        # --- Cloud methods (needed for Status and Notes) ---
        cloud_methods: List[str] = []
        if coverage != "Server Underlap":
            check_m = (
                sorted(spec_methods)
                if methods == ["ALL"]
                else [m for m in methods if m in spec_methods]
            )
            for m in check_m:
                xi = x_internal_map.get((norm, m), [])
                if "cloud" in xi:
                    cloud_methods.append(m)

        # --- Status ---
        if is_commented:
            status = "Deprecated"
        elif cloud_methods and len(cloud_methods) == len(spec_methods):
            status = "Active (Cloud-annotated)"
        else:
            status = "Active"

        # --- Schema-Backed ---
        backed = "TRUE" if spec_methods else "FALSE"

        # --- Schema Completeness (cross-check or legacy fallback) ---
        handler = route["handler"]
        use_crosscheck = (
            handler_io_map is not None
            and go_fields_map is not None
            and handler in handler_io_map
            and spec_methods
        )

        if use_crosscheck:
            # Cross-check: pick the first matching method for comparison
            check_methods_list = (
                sorted(spec_methods)
                if methods == ["ALL"]
                else [m for m in methods if m in spec_methods]
            )
            if check_methods_list:
                # Aggregate cross-check across all matching methods
                method_comps: List[str] = []
                agg_notes: List[str] = []
                for m in check_methods_list:
                    op = operations_map.get((norm, m))
                    if op:
                        sf = extract_spec_schema_fields(op, m)
                        comp, cnotes = cross_check_completeness(
                            handler, handler_io_map[handler],
                            go_fields_map, sf, m,
                        )
                        method_comps.append(comp)
                        agg_notes.extend(cnotes)
                    else:
                        method_comps.append("N/A")

                # Aggregate: best-of across methods
                if all(c == "Full" for c in method_comps):
                    completeness = "Full"
                elif any(c == "Full" for c in method_comps) or any(
                    c == "Partial" for c in method_comps
                ):
                    completeness = "Partial"
                elif all(c == "N/A" for c in method_comps):
                    completeness = "N/A"
                else:
                    completeness = "Stub"
                # Deduplicate notes
                seen_n: Set[str] = set()
                compl_notes = []
                for n in agg_notes:
                    if n not in seen_n:
                        seen_n.add(n)
                        compl_notes.append(n)
            else:
                completeness, compl_notes = _aggregate_completeness(
                    norm, methods, spec_data
                )
        else:
            completeness, compl_notes = _aggregate_completeness(
                norm, methods, spec_data
            )

        # --- Schema-Driven ---
        if handler in ("<inline>", "<unknown>"):
            driven, driven_reason = "FALSE", f"handler: {handler}"
        else:
            driven, driven_reason = schema_map.get(
                handler, ("FALSE", "handler not mapped")
            )

        # --- Notes (actionable summary) ---
        notes = _build_actionable_notes(
            coverage=coverage,
            status=status,
            is_commented=is_commented,
            completeness=completeness,
            compl_notes=compl_notes,
            driven=driven,
            handler=handler,
            cloud_methods=cloud_methods,
            spec_methods=spec_methods,
        )

        endpoints.append({
            "category": category,
            "subcategory": subcategory,
            "path": path,
            "methods": methods_str,
            "coverage": coverage,
            "status": status,
            "backed": backed,
            "completeness": completeness,
            "driven": driven,
            "notes": notes,
        })

    # ------------------------------------------------------------------
    # Pass 2: Spec-only endpoints (Schema Underlap)
    # ------------------------------------------------------------------
    for norm_path, spec_methods in sorted(all_paths.items()):
        if norm_path in router_norm_paths:
            continue

        original = original_paths.get(norm_path, norm_path)
        methods_sorted = sorted(spec_methods)
        category, subcategory = categorize(original, path_categories)

        # Determine x-internal across all methods for this path
        all_cloud = True
        any_cloud = False
        for m in methods_sorted:
            xi = x_internal_map.get((norm_path, m), [])
            if "cloud" in xi:
                any_cloud = True
            else:
                all_cloud = False

        # --- Coverage ---
        coverage = "Schema Underlap"

        # --- Status ---
        if all_cloud:
            status = "Cloud-only"
        else:
            status = "Unimplemented"

        # --- Schema-Backed ---
        backed = "TRUE"

        # --- Schema-Driven ---
        driven = "N/A"

        # --- Schema Completeness & Notes ---
        if status == "Cloud-only":
            # No router equivalent — spec completeness checks don't apply
            completeness = "N/A"
            notes = "No equivalent route in Meshery server; defined in spec as x-internal: cloud"
        else:
            completeness, compl_notes = _aggregate_completeness(
                norm_path, methods_sorted, spec_data
            )
            cloud_methods_list = [
                m for m in methods_sorted
                if "cloud" in x_internal_map.get((norm_path, m), [])
            ]
            notes = _build_actionable_notes(
                coverage=coverage,
                status=status,
                is_commented=False,
                completeness=completeness,
                compl_notes=compl_notes,
                driven=driven,
                handler="",
                cloud_methods=cloud_methods_list,
                spec_methods=set(methods_sorted),
            )

        endpoints.append({
            "category": category,
            "subcategory": subcategory,
            "path": original,
            "methods": ", ".join(methods_sorted),
            "coverage": coverage,
            "status": status,
            "backed": backed,
            "completeness": completeness,
            "driven": driven,
            "notes": notes,
        })

    return sorted(endpoints, key=endpoint_sort_key)


# ---------------------------------------------------------------------------
# Google Sheet — credentials from environment, never hardcoded
# ---------------------------------------------------------------------------

def _get_sheet_client():
    """Authenticate with Google Sheets using env-var credentials."""
    try:
        import gspread
        from google.oauth2.service_account import Credentials
    except ImportError:
        sys.exit(
            "Missing packages. Run: pip install gspread google-auth"
        )

    scopes = [
        "https://www.googleapis.com/auth/spreadsheets",
        "https://www.googleapis.com/auth/drive",
    ]

    # Option 1: inline JSON (config variable or env var, for CI / GitHub Actions)
    creds_json = GOOGLE_CREDENTIALS_JSON or os.environ.get("GOOGLE_CREDENTIALS_JSON")
    if creds_json:
        info = json.loads(creds_json)
        creds = Credentials.from_service_account_info(info, scopes=scopes)
        return gspread.authorize(creds)

    # Option 2: file path (config variable or env var, for local development)
    creds_file = GOOGLE_APPLICATION_CREDENTIALS or os.environ.get("GOOGLE_APPLICATION_CREDENTIALS")
    if creds_file and os.path.exists(creds_file):
        creds = Credentials.from_service_account_file(creds_file, scopes=scopes)
        return gspread.authorize(creds)

    return None


def _col_letter(idx: int) -> str:
    """Convert 0-based column index to sheet column letter (A, B, ... Z, AA, ...)."""
    result = ""
    while True:
        result = chr(65 + idx % 26) + result
        idx = idx // 26 - 1
        if idx < 0:
            break
    return result


# Columns the script compares and updates on matched rows.
# (column_index, endpoint_dict_key, human_label)
_UPDATABLE_COLUMNS = [
    (COL_COVERAGE, "coverage", "coverage"),
    (COL_STATUS, "status", "status"),
    (COL_BACKED, "backed", "backed"),
    (COL_COMPLETENESS, "completeness", "completeness"),
    (COL_DRIVEN, "driven", "driven"),
    (COL_NOTES, "notes", "notes"),
]


def update_sheet(
    endpoints: List[Dict[str, Any]],
    sheet_id: str,
    dry_run: bool = False,
) -> List[str]:
    """Diff computed endpoints against the sheet and apply updates.

    - Matches rows by normalized endpoint path + method overlap.
    - Updates Coverage, Status, Schema-Backed, Schema Completeness,
      Schema-Driven, and Notes columns when they differ.
    - Inserts new rows into matching category groups when possible.
    - Stamps the Change Log column on modified rows.
    """
    gc = _get_sheet_client()
    if not gc:
        print(
            "ERROR: No credentials found.\n"
            "  Set GOOGLE_CREDENTIALS_JSON (inline JSON for CI) or\n"
            "  GOOGLE_APPLICATION_CREDENTIALS (file path for local dev).",
            file=sys.stderr,
        )
        sys.exit(1)

    sheet = gc.open_by_key(sheet_id)
    ws = sheet.get_worksheet(4)

    print(f"Connected to worksheet: {ws.title}")
    current_rows = ws.get_all_values()

    # Index sheet rows by normalized path
    sheet_index: Dict[str, List[Tuple[int, Set[str]]]] = defaultdict(list)
    for idx, row in enumerate(current_rows):
        if idx == 0:
            continue
        ep = row[COL_ENDPOINTS].strip() if len(row) > COL_ENDPOINTS else ""
        if not ep:
            continue
        if not ep.startswith("/"):
            ep = "/" + ep
        norm = normalize_path(ep)
        raw_methods = row[COL_METHODS].strip() if len(row) > COL_METHODS else ""
        mset = {
            m.strip().upper()
            for m in raw_methods.replace(";", ",").split(",")
            if m.strip()
        }
        sheet_index[norm].append((idx, mset))

    changes: List[str] = []
    batch_updates: List[Dict[str, Any]] = []
    new_rows_info: List[Tuple[List[str], str, str]] = []
    matched_rows: Set[int] = set()
    today = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    for ep in endpoints:
        norm = normalize_path(ep["path"])
        ep_mset = {m.strip() for m in ep["methods"].split(",")}
        candidates = sheet_index.get(norm, [])

        # Find matching sheet row
        matched_idx = None
        for idx, sheet_mset in candidates:
            if idx in matched_rows:
                continue
            if (
                "ALL" in ep_mset
                or "ALL" in sheet_mset
                or ep_mset & sheet_mset
                or not sheet_mset
                or not ep_mset
            ):
                matched_idx = idx
                break

        if matched_idx is not None:
            matched_rows.add(matched_idx)
            row = current_rows[matched_idx]
            while len(row) < len(SHEET_COLUMNS):
                row.append("")

            row_changed = False
            for col_idx, field, label in _UPDATABLE_COLUMNS:
                old_val = row[col_idx].strip() if len(row) > col_idx else ""
                new_val = ep[field]
                if old_val != new_val:
                    cl = _col_letter(col_idx)
                    changes.append(
                        f"UPDATE row {matched_idx + 1} [{ep['path']}] "
                        f"{label}: '{old_val}' -> '{new_val}'"
                    )
                    batch_updates.append({
                        "range": f"{cl}{matched_idx + 1}",
                        "values": [[new_val]],
                    })
                    row_changed = True

            if row_changed:
                cl = _col_letter(COL_CHANGELOG)
                batch_updates.append({
                    "range": f"{cl}{matched_idx + 1}",
                    "values": [[today]],
                })
        else:
            new_row = [
                ep["category"],
                ep["subcategory"],
                ep["path"],
                ep["methods"],
                ep["coverage"],
                ep["status"],
                ep["backed"],
                ep["completeness"],
                ep["driven"],
                ep["notes"],
                today,
            ]
            changes.append(
                f"NEW ROW: {ep['path']} [{ep['methods']}] "
                f"coverage={ep['coverage']} status={ep['status']} "
                f"backed={ep['backed']} completeness={ep['completeness']} "
                f"driven={ep['driven']}"
            )
            new_rows_info.append((new_row, ep["category"], ep["subcategory"]))

    new_rows_info.sort(
        key=lambda item: endpoint_sort_key({
            "category": item[1],
            "subcategory": item[2],
            "path": item[0][COL_ENDPOINTS],
            "methods": item[0][COL_METHODS],
        })
    )

    # --- Apply batch cell updates ---
    if not dry_run and batch_updates:
        try:
            ws.batch_update(batch_updates, value_input_option="RAW")
            print(f"Batch updated {len(batch_updates)} cells")
        except Exception as exc:
            changes.append(f"BATCH UPDATE ERROR: {exc}")

    # --- Insert new rows ---
    if not dry_run and new_rows_info:
        _insert_rows_by_group(ws, new_rows_info, changes)

    return changes


def _insert_rows_by_group(
    ws,
    new_rows_info: List[Tuple[List[str], str, str]],
    changes: List[str],
) -> None:
    """Insert new rows into the correct category/sub-category block.

    Groups insertions by target position and processes from bottom to top
    so that earlier inserts don't shift indices for later ones.
    """
    try:
        all_rows = ws.get_all_values()
    except Exception as exc:
        changes.append(f"INSERT ERROR (read failed): {exc}")
        return

    # Build index: last row for each (category, subcategory) and category
    group_last_row: Dict[Tuple[str, str], int] = {}
    cat_last_row: Dict[str, int] = {}
    last_cat = ""
    last_sub = ""

    for idx, row in enumerate(all_rows):
        if idx == 0:
            continue

        cat = row[COL_CATEGORY].strip() if len(row) > COL_CATEGORY else ""
        sub = row[COL_SUBCATEGORY].strip() if len(row) > COL_SUBCATEGORY else ""

        if cat:
            last_cat = cat
        else:
            cat = last_cat

        if sub:
            last_sub = sub
        else:
            sub = last_sub

        if cat:
            group_last_row[(cat, sub)] = idx
            cat_last_row[cat] = idx

    # Classify each new row: targeted insert or append
    inserts: List[Tuple[int, List[str]]] = []
    append_rows: List[List[str]] = []

    for row_data, cat, sub in new_rows_info:
        target = group_last_row.get((cat, sub))
        if target is None:
            target = cat_last_row.get(cat)

        if target is not None:
            inserts.append((target, row_data))
            # Advance the group pointer so subsequent rows in the same
            # group land after this one rather than on top of it.
            group_last_row[(cat, sub)] = target + 1
            cat_last_row[cat] = target + 1
        else:
            append_rows.append(row_data)

    # Insert from bottom to top to preserve indices
    inserts.sort(key=lambda item: item[0], reverse=True)

    for insert_after, row_data in inserts:
        try:
            ws.insert_row(row_data, insert_after + 2, value_input_option="RAW")
        except Exception as exc:
            changes.append(f"INSERT ERROR at row {insert_after + 2}: {exc}")

    if append_rows:
        try:
            ws.append_rows(append_rows, value_input_option="RAW")
        except Exception as exc:
            changes.append(f"APPEND ROWS ERROR: {exc}")


# ---------------------------------------------------------------------------
# CLI
# ---------------------------------------------------------------------------

def main():
    parser = argparse.ArgumentParser(
        description=(
            "Audit Meshery API endpoints for schema coverage, completeness, "
            "and schema-driven status. Writes results to Google Sheet."
        )
    )
    parser.add_argument(
        "--repo",
        default=MESHERY_REPO_PATH or os.environ.get("MESHERY_REPO_PATH", "."),
        help=(
            "Path to the meshery/meshery repo root "
            "(default: MESHERY_REPO_PATH config at top of file, "
            "or $MESHERY_REPO_PATH env var, or cwd)"
        ),
    )
    parser.add_argument(
        "--spec",
        default=OPENAPI_SPEC_PATH or os.environ.get("OPENAPI_SPEC_PATH"),
        help=(
            "Path to an OpenAPI spec file (e.g. merged_openapi.yml from "
            "meshery/schemas). Overrides docs/data/openapi.yml in the repo."
        ),
    )
    parser.add_argument(
        "--sheet-id",
        default=SHEET_ID or os.environ.get("SHEET_ID"),
        help="Google Sheet ID (default: SHEET_ID config at top of file, or $SHEET_ID env var)",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Print diff without writing to the sheet",
    )
    parser.add_argument(
        "--verbose", "-v",
        action="store_true",
        help="Print per-endpoint details",
    )
    args = parser.parse_args()
    repo = Path(args.repo).resolve()

    # Validate
    if not (repo / ROUTER_FILE).exists():
        print(
            f"ERROR: {ROUTER_FILE} not found in {repo}\n"
            "Use --repo to point to the meshery/meshery repo root.",
            file=sys.stderr,
        )
        sys.exit(1)

    # --- Phase 1: Parse ---
    routes = parse_router(repo)
    spec_override = Path(args.spec).resolve() if args.spec else None
    spec_data = parse_openapi(repo, spec_override=spec_override)
    n_spec = len(spec_data["all_paths"])
    schema_map = build_schema_driven_map(repo)
    n_driven = sum(1 for s, _ in schema_map.values() if s == "TRUE")

    # Handler I/O types and Go struct fields for cross-check completeness
    handler_io_map = extract_handler_io_types(repo)
    go_fields_map = build_go_struct_fields(repo)
    n_io_extracted = sum(
        1 for io in handler_io_map.values()
        if io.get("request_type") or io.get("response_type")
    )
    print(f"Cross-check: {n_io_extracted}/{len(handler_io_map)} handlers "
          f"with extractable I/O types, {len(go_fields_map)} Go struct definitions")

    # --- Phase 2: Classify ---
    endpoints = classify_endpoints(
        routes, spec_data, schema_map,
        handler_io_map=handler_io_map,
        go_fields_map=go_fields_map,
    )
    total = len(endpoints)

    # Counts
    n_srv_under = sum(1 for e in endpoints if e["coverage"] == "Server Underlap")
    n_active = sum(1 for e in endpoints if e["status"] == "Active")
    n_cloud_compat = sum(1 for e in endpoints if e["status"] == "Active (Cloud-annotated)")
    n_deprecated = sum(1 for e in endpoints if e["status"] == "Deprecated")
    n_unimpl = sum(1 for e in endpoints if e["status"] == "Unimplemented")
    n_cloud = sum(1 for e in endpoints if e["status"] == "Cloud-only")
    b_true = sum(1 for e in endpoints if e["backed"] == "TRUE")
    comp_stub = sum(1 for e in endpoints if e["completeness"] == "Stub")
    d_true = sum(1 for e in endpoints if e["driven"] == "TRUE")
    d_part = sum(1 for e in endpoints if e["driven"] == "Partial")
    n_na_driven = sum(1 for e in endpoints if e["driven"] == "N/A")

    # --- Summary ---
    print("\nMeshery API Audit Summary")
    print("=" * 25)
    print(f"Sources: {len(routes)} routes | {n_spec} spec paths | "
          f"{len(schema_map)} handlers ({n_driven} schema-driven)")

    status_parts = []
    for label, count in [
        ("Active", n_active + n_cloud_compat),
        ("Deprecated", n_deprecated),
        ("Unimplemented", n_unimpl),
        ("Cloud-only", n_cloud),
    ]:
        if count:
            status_parts.append(f"{count} {label}")
    print(f"\n{total} endpoints (router + spec combined): {', '.join(status_parts)}")

    # Actionable gaps — only show non-zero items
    gaps = []
    if n_srv_under:
        gaps.append(f"  {n_srv_under} routes have no OpenAPI spec definition")
    if n_unimpl:
        gaps.append(f"  {n_unimpl} spec-defined endpoints have no router registration")
    not_backed = total - b_true
    if not_backed:
        gaps.append(f"  {not_backed}/{total} endpoints are not schema-backed")
    if comp_stub:
        gaps.append(f"  {comp_stub} endpoint schemas are stubs")
    router_total = total - n_na_driven
    not_driven = router_total - d_true - d_part
    if not_driven > 0:
        gaps.append(f"  {not_driven}/{router_total} handlers are not schema-driven")

    if gaps:
        print("\nNeeds attention:")
        for g in gaps:
            print(g)
    else:
        print("\nNo gaps found — all endpoints are fully covered.")

    print("\nRun with --verbose for per-endpoint details.")

    if args.verbose:
        print()
        for ep in endpoints:
            print(
                f"  {ep['path']:55s} [{ep['methods']:20s}] "
                f"cov={ep['coverage']:16s} st={ep['status']:14s} "
                f"bk={ep['backed']:5s} comp={ep['completeness']:7s} "
                f"drv={ep['driven']:7s}"
            )

    # --- Phase 3: Update sheet ---
    if not args.sheet_id:
        print(
            "\nNo --sheet-id provided. Set $SHEET_ID or pass --sheet-id "
            "to write results to Google Sheet."
        )
        sys.exit(0)

    label = "DRY RUN — previewing" if args.dry_run else "Updating"
    print(f"\n{label} Google Sheet...")

    changes = update_sheet(endpoints, args.sheet_id, args.dry_run)

    if not changes:
        print("Sheet is up to date.")
    else:
        n_updates = sum(1 for c in changes if c.startswith("UPDATE"))
        n_new = sum(1 for c in changes if c.startswith("NEW ROW"))
        n_errors = sum(1 for c in changes if "ERROR" in c)
        print(f"\n{len(changes)} change(s): {n_updates} updated, {n_new} new, {n_errors} error(s)")
        if args.verbose:
            for ch in changes:
                print(f"  {ch}")


if __name__ == "__main__":
    main()
