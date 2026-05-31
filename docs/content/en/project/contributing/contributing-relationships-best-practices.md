---
title: Contributing Relationships - Testing and Best Practices
description: How to test, validate, and debug relationship definitions before contribution
categories: [contributing]
---

# Relationship Testing and Best Practices

This guide complements the [Contributing to Model Relationships](/project/contributing/contributing-relationships) documentation by focusing on **practical testing and debugging techniques** for relationship definitions.

## Before You Start

Make sure you've read the [Contributing to Model Relationships](/project/contributing/contributing-relationships) guide to understand relationship structure, schema, and components.

## 1. Validate JSON Syntax

Before opening a PR, validate that your relationship JSON is well-formed.

```bash
python -m json.tool relationship.json > /dev/null
```

If the command completes without errors, the JSON syntax is valid.

This simple check catches malformed JSON and is sufficient for most relationship contributions.

## 2. Check for Common Relationship Mistakes

### ❌ Mistake 1: Mismatched `from`/`to` for Hierarchical Relationships

**Wrong:**
```json
{
  "kind": "hierarchical",
  "selectors": [{
    "allow": {
      "from": [{"kind": "Parent"}],
      "to": [{"kind": "Child"}]
    }
  }]
}
```

**Correct:**
```json
{
  "kind": "hierarchical",
  "selectors": [{
    "allow": {
      "from": [{"kind": "Child"}],
      "to": [{"kind": "Parent"}]
    }
  }]
}
```

**Why:** In Meshery, `from` is the child and `to` is the parent for hierarchical relationships. This is counterintuitive but important.

---

### ❌ Mistake 2: Cross-model relationships without proper scoping

**Problem:** Relating components across models without specifying both models.

**Wrong:**
```json
{
  "selectors": [{
    "allow": {
      "from": [{"kind": "ServiceAccount"}],
      "to": [{"kind": "Role"}]
    }
  }]
}
```

**Correct:**
```json
{
  "selectors": [{
    "allow": {
      "from": [
        {"kind": "ServiceAccount", "model": "kubernetes"}
      ],
      "to": [
        {"kind": "Role", "model": "aws-iam-controller"}
      ]
    }
  }]
}
```

**Why:** Without explicit model specification, Meshery may match components from unintended models.

---

### ❌ Mistake 3: Overlapping allow and deny selectors

**Problem:** Your `allow` and `deny` selectors contradict each other.

**Wrong:**
```json
{
  "selectors": [{
    "allow": {
      "from": [{"kind": "ServiceAccount"}],
      "to": [{"kind": "Role"}]
    },
    "deny": {
      "from": [{"kind": "ServiceAccount"}],
      "to": [{"kind": "Role"}]
    }
  }]
}
```

**Correct:** Use separate selector objects for different logic:
```json
{
  "selectors": [
    {
      "allow": {
        "from": [{"kind": "ServiceAccount", "metadata": {"labels": {"type": "system"}}}],
        "to": [{"kind": "Role"}]
      }
    },
    {
      "deny": {
        "from": [{"kind": "ServiceAccount", "metadata": {"labels": {"type": "test"}}}],
        "to": [{"kind": "Role"}]
      }
    }
  ]
}
```

**Why:** Conflicting selectors can make relationship behavior difficult to reason about and should generally be avoided.

---

## 3. Validate the Relationship Definition

Before submitting a relationship:

### Step 1: Verify the referenced field exists

Inspect the source component schema and confirm that the field referenced in the relationship actually exists.

For example:

```json
"delegatedSubnetResourceReference"
```

should exist in the source component schema before being used in a relationship definition.

### Step 2: Check for existing relationships

Search the repository for similar relationships before creating a new one.

This helps avoid duplicate contributions and ensures consistency with existing relationship patterns.

### Step 3: Compare with an existing relationship

Find a similar accepted relationship and follow the same structure whenever possible.

Using an existing relationship as a reference is often the fastest way to create a correct relationship definition.

### Step 4: Validate JSON syntax

Run:

```bash
python -m json.tool relationship.json > /dev/null
```

to ensure the JSON is valid.


## 4. Capture and Share Relationship Screenshots

Relationship screenshots provide visual proof that your definition works correctly and helps maintainers quickly validate the functionality.

### How to Take a Kanvas Screenshot

1. **Open Kanvas in Meshery Playground**
   - Visit [playground.meshery.io](https://playground.meshery.io/)
   - Import your custom model (or use existing models if testing standard relationships)

2. **Prepare the canvas view**
   - Drag both components involved in the relationship onto the canvas
   - Position them clearly with space between them
   - Zoom to ~75-100% so components are clearly visible

3. **Draw the relationship**
   - Hover over the first component (the "from" component)
   - Click and drag to the second component (the "to" component)
   - If your relationship is valid, the edge will appear with the correct visual style

4. **Take the screenshot**
   - Include:
     - Both components on the canvas
     - The relationship edge/line connecting them
     - Clear view of the relationship's visual style (color, line type, direction)
   - **On Windows:** Press `Win + Shift + S` and select the area
   - **On macOS:** Press `Cmd + Shift + 4` and select the area
   - **On Linux:** Use PrintScreen or your screenshot tool

### What Makes a Good Relationship Screenshot

✅ **Good screenshot includes:**
- Both components clearly visible and labeled
- The relationship edge connecting them (not blurry or cut off)
- Relationship direction is obvious (especially for directional relationships)
- Canvas is clean (no extra components cluttering the view)
- Image is at least 800x600 pixels (readable resolution)

❌ **Avoid:**
- Screenshot so zoomed in that context is lost
- Multiple unrelated components on canvas
- Cut-off components or relationship edges
- Very small or pixelated images

### Embedding Screenshots in Your PR

1. **Save your screenshot** (e.g., `relationship-iam-role-to-sa.png`)

2. **In your PR description**, add:
   ```markdown
   ## Relationship Visualization

   ![IAM Role to ServiceAccount relationship in Kanvas](/uploads/relationship-iam-role-to-sa.png)
   ```

3. **Alternatively, drag and drop directly:**
   - Open your PR description editor
   - Click the attachment icon or drag-and-drop the image
   - GitHub will automatically create the reference

4. **Add context near the screenshot:**
   ```markdown
   ## Relationship Visualization
   
   The screenshot below shows the IAM Role to Kubernetes ServiceAccount relationship 
   successfully creating an edge in Meshery Kanvas, demonstrating IRSA pattern support.
   
   ![Relationship screenshot](image-link)
   ```

### Example PR with Screenshots

When submitting multiple relationships, organize your PR description like:

```markdown
## Relationships Added

1. **IAM Role ↔ ServiceAccount (IRSA)**
   ![Screenshot: IRSA relationship](image-url)
   - Pattern: IRSA pod authentication
   - Components: AWS IAM Role + Kubernetes ServiceAccount

2. **SecurityGroup ↔ NetworkPolicy**
   ![Screenshot: SG to NP relationship](image-url)
   - Pattern: AWS SecurityGroup to K8s NetworkPolicy
   - Components: AWS SecurityGroup + Kubernetes NetworkPolicy
```

---

## 5. Pre-submission Checklist

Before opening a PR, verify:

* [ ] Relationship JSON is valid
* [ ] Referenced fields exist in the source component schema
* [ ] Similar relationships were reviewed for consistency
* [ ] Existing relationships were checked to avoid duplicates
* [ ] Relationship file is placed in the correct model directory
* [ ] Components referenced by the relationship exist
* [ ] Screenshot captured showing the relationship in Kanvas
* [ ] Screenshot included in the PR description

---

## 6. Debugging Tips

### Issue: Relationship does not appear as expected

Possible causes:

1. The referenced field does not exist in the source component schema.
2. The source or target component kind is incorrect.
3. The relationship already exists and conflicts with your definition.
4. The model name does not match the component's model.

### Debug Checklist

* Verify the referenced field exists in the component schema.
* Compare your relationship against a similar accepted relationship.
* Confirm component kinds and model names match existing definitions.
* Check for duplicate relationships before creating a new one.

---

## 7. Relationship File Naming

Relationship filenames vary across models and providers.

Before creating a new relationship file:

1. Look at existing relationship files in the same model directory.
2. Follow the naming pattern already used there.
3. Keep the filename consistent with nearby relationship definitions.

Using the existing naming style helps maintain consistency across models and avoids unnecessary review comments.

---

## 8. Contributing Your Relationship

Once validated and tested:

1. **Create a new branch:**
   ```bash
   git checkout -b relationships/add-{from}-{to}
   ```

2. **Commit with descriptive message:**
   ```bash
   git commit -s -m "Add {from} to {to} relationship (#ISSUE_NUMBER)"
   ```

3. **Include in your PR description:**

   - What relationship was added
   - Which models and components are involved
   - Why the relationship is useful
   - Screenshot showing the relationship in Kanvas
   - Link to any related issue or discussion

---

## Resources

- [Contributing to Model Relationships](/project/contributing/contributing-relationships)
- [Meshery Relationship Schema](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1alpha3/relationship)
- [Existing Relationships](https://github.com/meshery/meshery/tree/master/server/meshmodel)
- [OPA Policies](https://github.com/meshery/meshery/tree/master/server/meshmodel/meshery-core)
- [Issue #17096: AWS Relationships Epic](https://github.com/meshery/meshery/issues/17096)
