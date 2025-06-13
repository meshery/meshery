---
title: Troubleshooting Model Generation
description: Common failures and fixes encountered during Meshery model generation.
---

# Troubleshooting Model Generation

This guide helps troubleshoot common errors encountered while generating Meshery models using `mesheryctl` or the model-generator.

---

## 🛠️ Common Issues & Suggested Fixes

### 1. **Invalid Spreadsheet ID**
- **Cause**: Incorrect GSheet ID passed.
- **Fix**: Double-check the ID in URL: it should appear after `/d/` and before `/edit`.
- **Also**: Ensure the spreadsheet is shared or published for reading.

---

### 2. **Invalid Spreadsheet Credentials**
- **Cause**: Misconfigured or expired service account credentials.
- **Fix**: Regenerate credentials, encode with `base64 -w 0`, export via:
  ```bash
  export SHEET_CRED='<your encoded JSON>'
