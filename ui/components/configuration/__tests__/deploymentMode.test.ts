import { describe, expect, it } from 'vitest';
import type { ControllersConfigDoc } from '../ControllersConfigForm';
import {
  connectionDeploymentMode,
  serverDefaultDeploymentMode,
  takesEffectIn,
  SECTION_PATHS,
} from '../deploymentMode';

const modeDoc = (mode: 'operator' | 'embedded'): ControllersConfigDoc => ({
  operator: { deploymentMode: mode },
});

describe('takesEffectIn', () => {
  it('applies every setting in operator mode', () => {
    const everyPath = [
      ...SECTION_PATHS.operator,
      ...SECTION_PATHS.meshsync,
      ...SECTION_PATHS.broker,
    ];
    expect(everyPath.filter((path) => !takesEffectIn('operator', path))).toEqual([]);
  });

  it('applies only the mode itself and the output filters in embedded mode', () => {
    const everyPath = [
      ...SECTION_PATHS.operator,
      ...SECTION_PATHS.meshsync,
      ...SECTION_PATHS.broker,
    ];
    expect(everyPath.filter((path) => takesEffectIn('embedded', path))).toEqual([
      ['operator', 'deploymentMode'],
      ['meshsync', 'outputNamespaces'],
      ['meshsync', 'outputResources'],
    ]);
  });

  it('treats the whole broker section as inert in embedded mode', () => {
    // Meshery Broker exists only when Meshery Operator deploys it.
    expect(SECTION_PATHS.broker.some((path) => takesEffectIn('embedded', path))).toBe(false);
  });
});

describe('connectionDeploymentMode', () => {
  it('resolves from the connection override when the draft sets one', () => {
    const draft = modeDoc('operator');
    expect(
      connectionDeploymentMode({
        draft,
        persistedOverride: draft,
        serverDefault: modeDoc('embedded'),
        serverEffective: modeDoc('operator'),
      }),
    ).toEqual({
      mode: 'operator',
      sourceLabel: "this connection's override",
      unsaved: false,
      scope: 'connection',
    });
  });

  it('marks a draft mode that differs from the persisted override as unsaved', () => {
    const governance = connectionDeploymentMode({
      draft: modeDoc('operator'),
      persistedOverride: {},
      serverDefault: {},
      serverEffective: modeDoc('embedded'),
    });
    expect(governance.mode).toBe('operator');
    expect(governance.unsaved).toBe(true);
  });

  it('falls back to the server-wide default when the draft inherits', () => {
    expect(
      connectionDeploymentMode({
        draft: {},
        persistedOverride: {},
        serverDefault: modeDoc('operator'),
        serverEffective: modeDoc('operator'),
      }),
    ).toEqual({
      mode: 'operator',
      sourceLabel: 'the server-wide default',
      unsaved: false,
      scope: 'connection',
    });
  });

  it("reports the server's effective mode when neither editable layer sets one", () => {
    // The mode came from the materialized meshsync_deployment_mode cache or
    // MESHSYNC_DEFAULT_DEPLOYMENT_MODE; only the server can see those, and
    // guessing the built-in default here would declare a live cluster inert.
    expect(
      connectionDeploymentMode({
        draft: {},
        persistedOverride: {},
        serverDefault: {},
        serverEffective: modeDoc('operator'),
      }),
    ).toEqual({
      mode: 'operator',
      sourceLabel: 'the Meshery Server default',
      unsaved: false,
      scope: 'connection',
    });
  });

  it('says the mode is resolved on save when the user has just cleared an override', () => {
    const governance = connectionDeploymentMode({
      draft: {},
      persistedOverride: modeDoc('operator'),
      serverDefault: {},
      serverEffective: modeDoc('operator'),
    });
    expect(governance.mode).toBe('operator');
    expect(governance.unsaved).toBe(true);
    expect(governance.sourceLabel).toBe('the Meshery Server default, resolved when you save');
  });

  it('falls back to the built-in default when the server reports nothing', () => {
    expect(
      connectionDeploymentMode({ draft: {}, persistedOverride: {}, serverDefault: {} }).mode,
    ).toBe('embedded');
  });
});

describe('serverDefaultDeploymentMode', () => {
  it('reports the built-in default when the server-wide document sets no mode', () => {
    expect(serverDefaultDeploymentMode({})).toEqual({
      mode: 'embedded',
      sourceLabel: 'the built-in default',
      unsaved: false,
      scope: 'serverDefault',
    });
  });

  it('never scopes to a connection, so no field is rendered inert on this layer', () => {
    // A connection can override the mode to Operator, so a value stored here
    // is not dead - it is conditional, and the form annotates it instead.
    expect(serverDefaultDeploymentMode(modeDoc('embedded')).scope).toBe('serverDefault');
  });
});
