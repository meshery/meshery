# Event Parity Verification — View & Update

**Contributor:** [@KatalKavya96](https://github.com/KatalKavya96) 
**Goal:** Verify that “View” and “Update” events are generated consistently across Meshery clients.

---

## Results Table

| Event Type | Meshery Cloud Catalog | Meshery UI | Meshery.io Catalog | mesheryctl |
|-------------|----------------------|-------------|--------------------|-------------|
| View | ✅ | ✅ | ✅ | ✅ |
| Update | ✅ | ✅ | N/A | N/A |

✅ = Event generated ❌ = Not generated N/A = Not available

---

### Meshery Cloud Catalog
**Design tested:** `Kibana-Helm-Design-Fixed`  
- **View:** Opened design details page → page loaded correctly ✅  
- **Update:** Clicked *Edit*, renamed to `Kibana-Helm-Design-Fixed-test`, saved → toast appeared and “Updated At” changed ✅  

---

### Meshery UI (Local)
**Design tested:** `test-design-ui`  
- **View:** Clicking **Info** triggered `GET /api/meshmodels/models` with **200 OK** ✅  
- **Update:** Edited YAML (changed `replicas`) → save request **200 OK** ✅  

---

### Meshery.io Catalog
**Design tested:** Public design (e.g., `Kubernetes Deployment`)  
- **View:** Opened public design page on meshery.io → rendered successfully ✅  
- **Update:** Not supported for public users → N/A  

---

### mesheryctl
**Design tested:** `c3957bda` (`Online Boutique`)  
- **View:** `mesheryctl design view c3957bda` → details printed ✅  
- **Update:** N/A — no direct `design update` subcommand in this build (updates typically via apply/onboard).

