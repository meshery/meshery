// Shim: master's ui/components/connections/ConnectionTable.tsx still imports
// resetDatabase from '../graphql/queries/ResetDatabaseQuery', a path that
// was moved out of components/ during the GraphQL module migration. This
// re-export keeps the legacy import resolvable until the upstream file is
// updated to import from '@/graphql/queries/ResetDatabaseQuery' directly.
export { default } from '../../../graphql/queries/ResetDatabaseQuery';
