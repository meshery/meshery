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

{{< code code="kubectl apply -f meshery.yaml" >}}

### Multi-line Code

{{< code code=`curl -L https://meshery.io/install | PLATFORM=docker bash -

mesheryctl system start
mesheryctl system status` >}}

### Code with Special Characters

{{< code code=`echo "Hello <world> & 'universe'"` >}}

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

## 4. Alert Shortcode

### Info (default)

{{< alert title="Information" color="info" >}}
This is an informational alert. Supports **markdown** content.
{{< /alert >}}

### Warning

{{< alert title="Warning" color="warning" >}}
This is a warning alert.
{{< /alert >}}

### Danger

{{< alert title="Danger" color="danger" >}}
This is a danger alert.
{{< /alert >}}

### No title

{{< alert color="success" >}}
This alert has no title, just content.
{{< /alert >}}

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
