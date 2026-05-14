// Shim: master's ui/components/connections/ConnectionTable.tsx still imports
// LoadingScreen from '../LoadingComponents/LoadingComponent', a path that was
// renamed to ui/components/shared/LoadingState/LoadingComponent during the
// folder-casing normalization (#18745) but never updated in that one file.
// This re-export bridges the gap so vitest/vite/next can resolve the import.
// Remove this file once master's ConnectionTable.tsx is fixed to import from
// the canonical shared/LoadingState/ path.
export { default } from '../shared/LoadingState/LoadingComponent';
