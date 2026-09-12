import { describe, expect, it, vi } from 'vitest';

// designValidator pulls in `@/rtk-query/meshModel`, `@/rtk-query/design`, and
// `@/rtk-query/utils` which transitively boot the redux store + RTK query API
// (and a relay environment). For unit tests we mock those out so the machine
// graph and the pure helpers (formatDryRunResponse, selectValidator,
// selectComponentValidationResults, selectComponentDryRunResults) are
// exercisable in isolation.

vi.mock('@/rtk-query/meshModel', () => ({
  getComponentDefinition: vi.fn(),
}));
vi.mock('@/rtk-query/design', () => ({
  designsApi: {
    endpoints: { deployPattern: { initiate: vi.fn() }, undeployPattern: { initiate: vi.fn() } },
  },
}));
vi.mock('@/rtk-query/utils', () => ({
  initiateQuery: vi.fn(),
}));
vi.mock('@/utils/utils', () => ({
  encodeDesignFile: vi.fn(),
  processDesign: vi.fn(() => ({ components: [] })),
}));

import {
  designValidatorCommands,
  designValidatorEvents,
  formatDryRunResponse,
  selectComponentDryRunResults,
  selectComponentValidationResults,
  selectValidator,
  DESIGN_VALIDATOR_EVENTS,
} from '../validator/designValidator';

describe('designValidatorCommands', () => {
  it('builds a validate-design-schema command', () => {
    const cmd = designValidatorCommands.validateDesignSchema({
      design: { name: 'D1' },
      returnAddress: 'addr-1',
    });
    expect(cmd).toEqual({
      type: 'VALIDATE_DESIGN_SCHEMA',
      returnAddress: 'addr-1',
      data: { design: { name: 'D1' }, validationPayloadType: 'design' },
    });
  });

  it('builds a validate-design-component command', () => {
    const cmd = designValidatorCommands.validateDesignComponent({
      component: { id: 'c1' },
      returnAddress: 'addr-2',
    });
    expect(cmd).toEqual({
      type: 'VALIDATE_DESING_COMPONENT',
      returnAddress: 'addr-2',
      data: { component: { id: 'c1' }, validationPayloadType: 'component' },
    });
  });

  it('builds a dry-run deploy command', () => {
    const cmd = designValidatorCommands.dryRunDesignDeployment({
      design: { name: 'D1' },
      k8sContexts: ['ctx-1'],
      includeDependencies: true,
      returnAddress: 'addr-3',
    });
    expect(cmd.type).toBe('DRY_RUN_DESIGN');
    expect(cmd.returnAddress).toBe('addr-3');
    expect(cmd.data.dryRunType).toBe('Deploy');
    expect(cmd.data.includeDependencies).toBe(true);
  });

  it('builds a dry-run undeploy command', () => {
    const cmd = designValidatorCommands.dryRunDesignUnDeployment({
      design: { name: 'D1' },
      k8sContexts: ['ctx-1'],
      includeDependencies: false,
      returnAddress: 'addr-4',
    });
    expect(cmd.data.dryRunType).toBe('Undeploy');
    expect(cmd.data.includeDependencies).toBe(false);
  });
});

describe('designValidatorEvents', () => {
  it('builds a design-schema-validated event', () => {
    const ev = designValidatorEvents.designSchemaValidated({
      design: { id: 'd-1' },
      validationResults: { ok: true },
    });
    expect(ev.type).toBe(DESIGN_VALIDATOR_EVENTS.DESIGN_SCHEMA_VALIDATION_DONE);
    expect(ev.data).toEqual({ design: { id: 'd-1' }, validationResults: { ok: true } });
  });

  it('builds a tap-on-error event', () => {
    const ev = designValidatorEvents.tapOnError({
      error: 'boom',
      type: 'schema',
      component: { id: 'c-1' },
    });
    expect(ev.type).toBe(DESIGN_VALIDATOR_EVENTS.TAP_ON_ERROR);
    expect(ev.data).toEqual({ error: 'boom', type: 'schema', component: { id: 'c-1' } });
  });
});

describe('formatDryRunResponse', () => {
  it('returns an empty list when nothing is provided', () => {
    expect(formatDryRunResponse(null)).toEqual([]);
    expect(formatDryRunResponse(undefined)).toEqual([]);
  });

  it('returns an empty list when there are no failed contexts', () => {
    const response = {
      compA: { ctx1: { success: true } },
    };
    expect(formatDryRunResponse(response)).toEqual([]);
  });

  it('reports component errors with the Causes list when present', () => {
    const response = {
      compA: {
        ctx1: {
          success: false,
          component: { name: 'compA' },
          error: { Causes: [{ message: 'missing field' }] },
        },
      },
    };
    const errors = formatDryRunResponse(response);
    expect(errors).toHaveLength(1);
    expect(errors[0]).toEqual({
      type: 'ComponentError',
      compName: 'compA',
      component: { name: 'compA' },
      contextId: 'ctx1',
      errors: [{ message: 'missing field' }],
    });
  });

  it('falls back to Status when Causes is missing', () => {
    const response = {
      compA: {
        ctx1: {
          success: false,
          component: { name: 'compA' },
          error: { Status: 'failed' },
        },
      },
    };
    const errors = formatDryRunResponse(response);
    expect(errors[0].errors).toEqual(['failed']);
  });

  it('returns an empty errors array when error has neither Causes nor Status', () => {
    const response = {
      compA: { ctx1: { success: false, component: { name: 'compA' }, error: {} } },
    };
    const errors = formatDryRunResponse(response);
    expect(errors[0].errors).toEqual([]);
  });

  it('skips components whose context list is null', () => {
    const response = {
      compA: null,
      compB: { ctx1: { success: true } },
    };
    expect(formatDryRunResponse(response)).toEqual([]);
  });

  it('produces one error entry per failed context', () => {
    const response = {
      compA: {
        ctx1: { success: false, error: { Causes: [{ m: 1 }] } },
        ctx2: { success: false, error: { Causes: [{ m: 2 }] } },
        ctx3: { success: true },
      },
    };
    const errors = formatDryRunResponse(response);
    expect(errors).toHaveLength(2);
    expect(errors.map((e) => e.contextId).sort()).toEqual(['ctx1', 'ctx2']);
  });
});

describe('selectValidator', () => {
  it('returns the validator slot from context when present', () => {
    const state = { context: { schemaValidator: { id: 'sv' } } };
    expect(selectValidator(state, 'schemaValidator')).toEqual({ id: 'sv' });
  });

  it('returns null when missing', () => {
    expect(selectValidator({ context: {} }, 'schemaValidator')).toBeNull();
  });

  it('returns null when accessing throws and logs a warning', () => {
    const warn = vi.spyOn(console, 'warn').mockImplementation(() => {});
    // simulate context being a getter that throws
    const trap = {
      get context() {
        throw new Error('boom');
      },
    };
    expect(selectValidator(trap, 'schemaValidator')).toBeNull();
    expect(warn).toHaveBeenCalled();
    warn.mockRestore();
  });
});

describe('selectComponentValidationResults', () => {
  it('returns null when there is no schema validator', () => {
    expect(selectComponentValidationResults({ context: {} }, 'c-1')).toBeNull();
  });

  it('returns null when validation results are missing', () => {
    const state = {
      context: {
        schemaValidator: {
          getSnapshot: () => ({ context: { validationResults: null } }),
        },
      },
    };
    expect(selectComponentValidationResults(state, 'c-1')).toBeNull();
  });

  it('finds the validation result by component traits.meshmap.id', () => {
    const result1 = { component: { traits: { meshmap: { id: 'c-1' } } }, errors: ['e'] };
    const result2 = { component: { traits: { meshmap: { id: 'c-2' } } }, errors: [] };
    const state = {
      context: {
        schemaValidator: {
          getSnapshot: () => ({
            context: { validationResults: { a: result1, b: result2 } },
          }),
        },
      },
    };
    expect(selectComponentValidationResults(state, 'c-2')).toBe(result2);
  });
});

describe('selectComponentDryRunResults', () => {
  it('returns null when there is no dry-run validator', () => {
    expect(selectComponentDryRunResults({ context: {} }, 'compA')).toBeNull();
  });

  it('finds the dry-run result by compName', () => {
    const state = {
      context: {
        dryRunValidator: {
          getSnapshot: () => ({
            context: {
              validationResults: [
                { compName: 'compA', errors: ['e'] },
                { compName: 'compB', errors: [] },
              ],
            },
          }),
        },
      },
    };
    expect(selectComponentDryRunResults(state, 'compA')?.compName).toBe('compA');
  });
});
