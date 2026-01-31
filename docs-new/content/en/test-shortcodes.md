---
title: Shortcode Test Page
description: Test page for validating Hugo shortcodes
draft: true
---

# Shortcode Test Page

This page tests the Hugo shortcodes created for the Jekyll to Hugo migration.

---

## 1. Code Shortcode

### Basic Usage

{{< code >}}
kubectl apply -f meshery.yaml
{{< /code >}}

### Multi-line Code

{{< code >}}
curl -L https://meshery.io/install | PLATFORM=docker bash -

mesheryctl system start
mesheryctl system status
{{< /code >}}

### Code with Special Characters

{{< code >}}
echo "Hello <world> & 'universe'"
{{< /code >}}

---

## 2. Discuss Shortcode

{{< discuss >}}

---

## 3. Related Discussions Shortcode

### Meshery Tag

{{< related-discussions tag="meshery" >}}

### Mesheryctl Tag

{{< related-discussions tag="mesheryctl" >}}

---

## Validation Checklist

- [ ] Code blocks render with dark background
- [ ] Code blocks have proper padding and styling
- [ ] Special characters are escaped correctly
- [ ] Discuss alert shows with dark theme
- [ ] Discussion forum link works
- [ ] Related discussions show 10 topics
- [ ] Topic dates format correctly
- [ ] Topic links work
- [ ] Author names display correctly
- [ ] Dark mode compatible
