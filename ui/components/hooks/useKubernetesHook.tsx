// Shim: master's ui/components/connections/ConnectionTable.tsx still imports
// useKubernetesHook from '../hooks/useKubernetesHook', a path that was
// moved to ui/utils/hooks/ during the hooks module migration. This
// re-export keeps the legacy import resolvable until the upstream file is
// updated to import from '@/utils/hooks/useKubernetesHook' directly.
export { default } from '../../utils/hooks/useKubernetesHook';
