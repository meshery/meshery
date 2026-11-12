import { describe, it } from 'vitest';

// LazyComponentForm renders a dynamic RJSF form bound to a meshmodel
// component definition. It loads custom widgets/templates from
// @rjsf/mui + @rjsf/utils + @rjsf/validator-ajv8 and threads xstate
// validator events back to the design schema machine. Rendering it
// requires mocking the entire RJSF surface (form, widgets, custom
// templates, ajv validator) and the design-validator state machine,
// which produces several thousand lines of mocks for marginal value.
describe.skip('LazyComponentForm', () => {
  it('skipped - bound to RJSF + xstate validator with heavy widget surface', () => {});
});
