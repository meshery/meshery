import { describe, it } from 'vitest';

// finalizeDeployment.tsx uses CommonJS `require()` at module scope to pull
// in @sistent/sistent and ./common, which short-circuits vi.mock factories
// because the mocks are hoisted but require() ignores the ESM transformer
// at runtime. The component is also tightly coupled to redux + theme
// hooks (`selectSelectedEnvs`, useTheme) which are difficult to mock
// without a real store. Covered transitively by DeployStepper integration.
describe.skip('FinalizeDeployment', () => {
  it('skipped - component uses runtime require() that bypasses vi.mock', () => {});
});
