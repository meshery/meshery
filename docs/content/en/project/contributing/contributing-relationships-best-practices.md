---
title: Contributing Relationships - Testing and Best Practices
description: How to test, validate, and debug relationship definitions before contribution
categories: [contributing]
---

# Relationship Testing and Best Practices

This guide complements the [Contributing to Model Relationships](/project/contributing/contributing-relationships) documentation by focusing on **practical testing and debugging techniques** for relationship definitions.

## Before You Start

Make sure you've read the [Contributing to Model Relationships](/project/contributing/contributing-relationships) guide to understand relationship structure, schema, and components.

## 1. Validate Your Relationship JSON Schema

Before testing in Meshery, validate your relationship definition against the schema.

### Using JSON Schema Validator Online

1. Go to [jsonschemavalidator.io](https://www.jsonschemavalidator.io/)
2. In the **JSON Schema** section, paste the schema from [Meshery schemas repository](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1alpha3/relationship)
3. In the **JSON Input** section, paste your relationship definition
4. Click **Validate** - errors will appear if your JSON doesn't match the schema

### Using `jsonschema` (CLI)

If you have Python installed:

```bash
pip install jsonschema
jsonschema -i your-relationship.json schema.json
```

**Common validation errors:**
- Missing required fields: `kind`, `type`, `selectors`
- Incorrect data types: `selectors` must be an array
- Missing or malformed `from`/`to` in selector entries
- Typos in enum values (e.g., `kind: "edge"` vs `kind: "Edge"`)

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

### ❌ Mistake 3: Not specifying version for model-specific relationships

**Problem:** Your relationship only applies to specific API versions, but you use `*`.

**Wrong:**
```json
{
  "model": "kubernetes",
  "version": "*",
  "selectors": [{
    "allow": {
      "from": [{"kind": "Pod", "apiVersion": "v1"}],
      "to": [{"kind": "PersistentVolume", "apiVersion": "v1"}]
    }
  }]
}
```

**Correct (if only v1):**
```json
{
  "model": "kubernetes",
  "version": "1.28.0",
  "selectors": [{
    "allow": {
      "from": [{"kind": "Pod"}],
      "to": [{"kind": "PersistentVolume"}]
    }
  }]
}
```

**Why:** Specifying exact versions prevents unexpected behavior when older API versions are used.

---

### ❌ Mistake 4: Missing `evaluationQuery` for custom relationship types

**Problem:** Your relationship won't evaluate correctly without the right OPA policy.

**Wrong:**
```json
{
  "kind": "edge",
  "type": "custom-type",
  "subType": "custom-subtype"
  // No evaluationQuery!
}
```

**Correct:**
```json
{
  "kind": "edge",
  "type": "network",
  "subType": "custom-subtype",
  "evaluationQuery": "edge_network_relationship"
}
```

**Why:** The `evaluationQuery` tells Meshery's policy engine which OPA rules to use. If it doesn't match an existing policy, your relationship won't be evaluated.

---

### ❌ Mistake 5: Overlapping allow and deny selectors

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

**Why:** Conflicting selectors in the same object cause undefined behavior.

---

## 3. Test Locally Before Submitting

### Step 1: Create your relationship JSON file

Place it in the correct model directory:
```
server/meshmodel/{model}/{version}/v1.0.0/relationships/{relationship-name}.json
```

Example:
```
server/meshmodel/aws-iam-controller/v1.0.0/v1.0.0/relationships/edge-binding-role-to-serviceaccount.json
```

### Step 2: Run Meshery locally and import the model

```bash
# From Meshery repo root
make server
```

Then in the Meshery UI:
1. Go to **Settings** → **Models**
2. Click **Import Custom Model**
3. Point to your model folder containing the relationships
4. Verify relationships appear without errors

### Step 3: Verify in Kanvas

1. Open **Kanvas** (visual designer)
2. Search for and add both components (e.g., IAM Role + ServiceAccount)
3. Try to draw a relationship edge between them
4. If your relationship is valid, the edge should:
   - Allow connection
   - Show correct visual style (color, line type)
   - Not produce error messages in browser console

### Step 4: Check browser console for errors

Open DevTools (**F12** → **Console**):
- Look for any `relationship` or `schema` validation errors
- Errors like `"kind not recognized"` indicate your relationship definition has issues

---

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

- [ ] Relationship JSON passes schema validation (section 1)
- [ ] No common mistakes from section 2
- [ ] File placed in correct directory (`server/meshmodel/{model}/*/v1.0.0/relationships/`)
- [ ] File name follows pattern: `{type}-{subtype}-{from}-to-{to}.json`
- [ ] Tested locally in Meshery (section 3)
- [ ] No browser console errors when viewing in Kanvas
- [ ] Screenshots captured showing relationship working in Kanvas
- [ ] Screenshots included in PR description
- [ ] `from` and `to` components exist in their respective models
- [ ] `evaluationQuery` matches an existing OPA policy for your `kind`/`type`/`subType`
- [ ] For cross-model relationships: both `model` values are correct
- [ ] Selectors have meaningful conditions (not too broad, not contradictory)
- [ ] Relationship has descriptive `description` field
- [ ] Components can actually be related in real-world scenario

---

## 6. Debugging Tips

### Issue: Relationship appears but won't connect

**Possible causes:**
1. Selectors too restrictive (component doesn't match all conditions)
2. Wrong `model` name for a component
3. `evaluationQuery` invalid or missing

**Debug:**
- Add more logging to your selector conditions
- Verify component names match exactly (case-sensitive)
- Check that both components have `status: enabled` in their model

### Issue: Relationship connects but doesn't evaluate

**Possible causes:**
1. OPA policy (`evaluationQuery`) doesn't exist
2. Policy name has typo

**Debug:**
- Verify `evaluationQuery` matches a policy in `server/meshmodel/meshery-core/*/policies/`
- Check policy file names: `{kind}_{type}_{subtype}_relationship.rego`

### Issue: Relationship appears in wrong place visually

**Possible causes:**
1. Incorrect `type` or `subType` for visual representation
2. Missing OPA policy for your relationship kind

**Debug:**
- Cross-reference your `type`/`subType` against the [Relationship Visualizations](contributing-relationships#relationship-visualizations)
- Ensure a corresponding OPA policy exists

---

## 7. Relationship Naming Convention

Follow this pattern for filenames and metadata:

```text
{kind}-{type}-{from-component}-to-{to-component}.json
```

### Examples

✅ Good names:
- `edge-binding-role-to-serviceaccount.json` (clear, descriptive)
- `hierarchical-parent-namespace-to-pod.json` (direction is clear)
- `edge-network-service-to-pod.json` (specific and brief)

❌ Avoid:
- `relationship-1.json` (too generic)
- `aws-k8s-relationship.json` (ambiguous)
- `k8s_Service_Pod_Network.json` (inconsistent formatting)

---

## 8. Contributing Your Relationship

Once validated and tested:

1. **Create a new branch:**
   ```bash
   git checkout -b relationships/add-{from}-{to}
   ```

2. **Commit with descriptive message:**
   ```bash
   git commit -m "Add {from} to {to} relationship (#ISSUE_NUMBER)"
   ```

3. **Include in PR description:**
   - Which components relate
   - Real-world use case (e.g., "IRSA pattern for EKS pods")
   - Screenshots from Kanvas showing relationship working
   - Reference relevant issue (e.g., #17096)

4. **PR description example:**
   ```
   ## Description
   Adds relationship between AWS IAM Role and Kubernetes ServiceAccount for IRSA pattern support.
   
   ## Changes
   - Added edge-binding-role-to-serviceaccount.json
   - Supports cross-model relationships between aws-iam-controller and kubernetes
   - Uses annotation-based selector for pod authentication
   
   ## Relationship Visualization
   ![IRSA Relationship](image-url)
   
   ## Testing
   - Validated against schema ✓
   - Tested in Kanvas ✓
   - Screenshots included ✓
   - No console errors ✓
   
   Fixes #17096
   ```

---

## Resources

- [Contributing to Model Relationships](/project/contributing/contributing-relationships)
- [Meshery Relationship Schema](https://github.com/meshery/schemas/tree/master/schemas/constructs/v1alpha3/relationship)
- [Existing Relationships](https://github.com/meshery/meshery/tree/master/server/meshmodel)
- [OPA Policies](https://github.com/meshery/meshery/tree/master/server/meshmodel/meshery-core)
- [Issue #17096: AWS Relationships Epic](https://github.com/meshery/meshery/issues/17096)
