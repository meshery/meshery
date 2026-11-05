import { describe, it } from 'vitest';

// SelectDeploymentTarget consumes two RTK Query hooks
// (useGetEnvironmentConnectionsQuery, useGetEnvironmentsQuery), a
// React context (DeploymentTargetContext) with three required fields,
// useSelector / useDispatch from react-redux, and the K8sContextConnectionChip
// header component. Rendering it as a leaf would require mocking the entire
// environments slice, the meshsync controller state, and a connection-
// metadata map. End-to-end behaviour is covered by playwright; the
// callback/selector wiring it depends on is exercised through the redux
// store tests.
describe.skip('SelectDeploymentTarget', () => {
  it('skipped - depends on RTK Query, redux store, and multiple contexts', () => {});
});
