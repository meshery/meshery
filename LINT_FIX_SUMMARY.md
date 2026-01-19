# ESLint Fixes for PR #17020

## Summary
Fixed 96 ESLint errors in the TypeScript migration PR by prefixing unused parameters with underscore (_).

## Changes Made

### Files Modified (27 files):
1. `ui/components/MesheryAdapterPlayComponent.tsx`
2. `ui/components/MesheryFilters/Filters.tsx`
3. `ui/components/MesheryFilters/FiltersCard.tsx`
4. `ui/components/MesheryFilters/FiltersGrid.tsx`
5. `ui/components/PromptComponent.tsx`
6. `ui/components/Settings/MesherySettings.tsx`
7. `ui/components/Settings/Registry/ComponentTree.tsx`
8. `ui/components/Settings/Registry/CreateModelModal.tsx`
9. `ui/components/Settings/Registry/ImportModel.tsx`
10. `ui/components/Settings/Registry/ImportModelModal.tsx`
11. `ui/components/Settings/Registry/MeshModelComponent.tsx`
12. `ui/components/Settings/Registry/MeshModelDetails.tsx`
13. `ui/components/Settings/Registry/MesheryTreeView.tsx`
14. `ui/components/Settings/Registry/MesheryTreeViewItem.tsx`
15. `ui/components/Settings/Registry/MesheryTreeViewModel.tsx`
16. `ui/components/Settings/Registry/MesheryTreeViewRegistrants.tsx`
17. `ui/components/Settings/Registry/RelationshipTree.tsx`
18. `ui/components/Settings/Registry/StyledTreeItem.tsx`
19. `ui/components/Settings/Registry/VersionedModelComponentTree.tsx`
20. `ui/components/Settings/Registry/VersionedModelRelationshipTree.tsx`
21. `ui/components/UserPreferences/index.tsx`
22. `ui/components/graphql/mutations/AdapterStatusMutation.tsx`
23. `ui/components/graphql/mutations/OperatorStatusMutation.tsx`
24. `ui/components/graphql/subscriptions/ClusterResourcesSubscription.tsx`
25. `ui/components/graphql/subscriptions/MeshModelSummarySubscription.tsx`
26. `ui/components/graphql/subscriptions/MesheryControllersStatusSubscription.tsx`
27. `ui/components/hooks/useTestIDs.tsx`

## Fix Pattern

All fixes follow the ESLint rule that requires unused parameters to start with underscore (_).

### Examples:
- `onClick: (event: React.MouseEvent) => void` → `onClick: (_event: React.MouseEvent) => void`
- `handleSubmit: (data: any) => void` → `handleSubmit: (_data: any) => void`
- `interface UserData` → `interface _UserData` (for unused type definitions)
- Fixed JSX namespace issue by using `React.JSX.Element` instead of `JSX.Element`
- Removed unused import `ChangeEvent` from MesheryAdapterPlayComponent.tsx

## How to Apply

Apply the generated patch file to PR #17020:
```bash
git apply /tmp/0001-UI-Fix-ESLint-errors-prefix-unused-params-with-under.patch
```

## Verification

After applying fixes:
```bash
cd ui && npx eslint .
```

Expected result: 0 errors, 1 warning (about synchronous script in _document.js - non-blocking)

## Original Errors

- 96 ESLint errors about unused variables/parameters
- 1 warning about synchronous scripts (kept as-is, non-blocking)

Total: 97 problems → 1 warning after fixes
