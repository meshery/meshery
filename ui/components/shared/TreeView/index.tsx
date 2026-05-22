// Re-export the shared TreeView wrapper from its discoverable filename so
// `import { SimpleTreeView } from '../../shared/TreeView'` continues to resolve.
// The actual `@mui/x-tree-view` boundary lives in `./TreeView.tsx`.
export * from './TreeView';
