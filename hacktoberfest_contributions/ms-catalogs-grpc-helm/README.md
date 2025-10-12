# MS Catalogs gRPC Helm Chart Design

## 📦 Publishing Instructions

### Method 1: Via Meshery UI
1. Open Meshery at http://localhost:9081
2. Navigate to **Designs** section
3. Click **Import Design**
4. Select file: `design.yml`
5. Click **Publish to Catalog**

### Method 2: Via Meshery CLI
```bash
# Import the design
mesheryctl design import -f hacktoberfest_contributions/ms-catalogs-grpc-helm/design.yml

# Publish to catalog
mesheryctl design publish ms-catalogs-grpc-helm
```

## 🎯 Design Details
- **Name**: MS Catalogs gRPC Helm Chart
- **Components**: 8 Kubernetes resources
- **Status**: Production-ready

## ✅ Ready for Catalog Publication
This design is fully tested and ready for Meshery catalog publication.
