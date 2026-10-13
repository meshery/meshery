#!/usr/bin/env python3
"""Check Meshery docs URLs for alias problems and broken redirects.

Uses only the Python standard library so it can run in CI without new
dependencies.

Usage (run from the docs/ directory):

  # Source check: fail on duplicate aliases, or an alias that equals the
  # URL of a real (published) page.
  python3 scripts/check_docs_urls.py aliases

  # Build check: every URL that existed in an earlier build must still
  # resolve in the new build, either as a page or as an alias redirect.
  python3 scripts/check_docs_urls.py redirects --old old-urls.txt --public public

  # Write the URL list of a build, for use as --old later.
  python3 scripts/check_docs_urls.py urls --public public > old-urls.txt
"""
import argparse
import os
import re
import sys

CONTENT = os.path.join("content", "en")
REFRESH = re.compile(r'http-equiv="?refresh"?', re.I)


def norm(url):
    url = url.strip().strip("'\"")
    if not url.startswith("/"):
        url = "/" + url
    if not url.endswith("/") and "." not in url.rsplit("/", 1)[-1]:
        url += "/"
    return url.lower()


def front_matter(path):
    with open(path, encoding="utf-8") as fh:
        text = fh.read()
    if not text.startswith("---"):
        return {}
    end = text.find("\n---", 3)
    if end == -1:
        return {}
    data, key = {}, None
    for line in text[3:end].splitlines():
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        m = re.match(r"^([A-Za-z_][\w-]*):\s*(.*)$", line)
        if m:
            key, val = m.group(1), m.group(2).strip()
            if val.startswith("[") and val.endswith("]"):
                data[key] = [v.strip().strip("'\"") for v in val[1:-1].split(",") if v.strip()]
            elif val:
                data[key] = val.strip("'\"")
            else:
                data[key] = []
            continue
        m = re.match(r"^\s*-\s*(.+)$", line)
        if m and key and isinstance(data.get(key), list):
            data[key].append(m.group(1).strip().strip("'\""))
    return data


def page_url(rel, fm):
    if isinstance(fm.get("url"), str):
        return norm(fm["url"])
    parts = rel.replace(os.sep, "/").split("/")
    name = parts[-1]
    if name in ("_index.md", "index.md"):
        parts = parts[:-1]
    else:
        parts[-1] = fm.get("slug") if isinstance(fm.get("slug"), str) else name[:-3]
    return norm("/".join(parts))


def is_published(fm):
    return str(fm.get("draft", "false")).lower() != "true" and \
        str(fm.get("publish", "true")).lower() != "false"


def check_aliases():
    pages, aliases = {}, {}
    for root, dirs, files in os.walk(CONTENT):
        for name in files:
            if not name.endswith(".md"):
                continue
            path = os.path.join(root, name)
            fm = front_matter(path)
            if not is_published(fm):
                continue
            rel = os.path.relpath(path, CONTENT)
            pages.setdefault(page_url(rel, fm), []).append(rel)
            al = fm.get("aliases", [])
            for a in al if isinstance(al, list) else [al]:
                aliases.setdefault(norm(a), []).append(rel)
    errors = []
    for url, srcs in sorted(aliases.items()):
        if len(set(srcs)) > 1:
            errors.append("duplicate alias %s in %s" % (url, ", ".join(sorted(set(srcs)))))
        if url in pages:
            owners = set(pages[url])
            if owners != set(srcs):
                errors.append("alias %s (in %s) collides with page %s" % (
                    url, ", ".join(sorted(set(srcs))), ", ".join(sorted(owners))))
    for e in errors:
        print("ERROR:", e)
    print("checked %d pages, %d aliases, %d problems" % (len(pages), len(aliases), len(errors)))
    return 1 if errors else 0


def build_urls(public):
    urls = set()
    for root, dirs, files in os.walk(public):
        if "index.html" in files:
            rel = os.path.relpath(root, public).replace(os.sep, "/")
            urls.add("/" if rel == "." else "/" + rel + "/")
    return urls


def check_redirects(old_file, public):
    with open(old_file, encoding="utf-8") as fh:
        old = {norm(l) for l in fh if l.strip()}
    new = build_urls(public)
    missing, redirected = [], 0
    for url in sorted(old - {u.lower() for u in new}):
        missing.append(url)
    for url in sorted(old):
        f = os.path.join(public, url.strip("/"), "index.html")
        if os.path.isfile(f):
            with open(f, encoding="utf-8", errors="replace") as fh:
                if REFRESH.search(fh.read(4096)):
                    redirected += 1
    for u in missing:
        print("MISSING:", u)
    print("%d old URLs, %d now redirects, %d missing" % (len(old), redirected, len(missing)))
    return 1 if missing else 0


def main():
    p = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    sub = p.add_subparsers(dest="cmd", required=True)
    sub.add_parser("aliases")
    r = sub.add_parser("redirects")
    r.add_argument("--old", required=True)
    r.add_argument("--public", default="public")
    u = sub.add_parser("urls")
    u.add_argument("--public", default="public")
    args = p.parse_args()
    if args.cmd == "aliases":
        return check_aliases()
    if args.cmd == "redirects":
        return check_redirects(args.old, args.public)
    for url in sorted(build_urls(args.public)):
        print(url)
    return 0


if __name__ == "__main__":
    sys.exit(main())
