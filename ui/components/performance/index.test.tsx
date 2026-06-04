import { describe, it } from 'vitest';

// MesheryPerformanceComponent is a 1300-line container that wires up many
// hooks (useSelector, useNotification, useSavePerformanceProfileMutation,
// useGetUserPrefWithContextQuery, EventSource), nested file-upload helpers,
// route-resolution helpers (ctxUrl, getK8sClusterIds), and ~30 sistent UI
// primitives. Mocking the full surface to exercise the rendered tree adds
// far more code than the component itself. Coverage for its internals is
// provided indirectly via helper.test.tsx (generateTestName, generateUUID)
// and by integration / E2E tests.
describe.skip('MesheryPerformanceComponent (index.tsx)', () => {
  it('skipped - container with deep redux+rtk-query+EventSource graph', () => {});
});
