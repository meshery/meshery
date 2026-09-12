import { renderHook, act } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

// Permission outcome is a spy so the delayed-resolution case (permissions start
// denied while capabilities load, then flip) can be simulated between renders.
const canSpy = vi.fn((_key?: unknown) => true);

vi.mock('@sistent/sistent', () => ({
  useHasPermission: (key: unknown) => canSpy(key),
}));

// The real step registry pulls in the whole step UI; the wizard hook only needs
// a step list of a known length, so it is stubbed with two inert steps.
vi.mock('./registry', () => ({
  buildSteps: () => [
    { id: 'select', label: 'Choose Connection', Component: () => null },
    { id: 'details', label: 'Configure Connection', Component: () => null },
  ],
}));

vi.mock('@/rtk-query/connection', () => {
  // Declared inside the factory: `vi.mock` is hoisted above module-scope consts.
  const noopMutation = () => [vi.fn(() => ({ unwrap: () => Promise.resolve({}) }))];
  return {
    useAddKubernetesConfigMutation: noopMutation,
    useDiscoverKubernetesContextsMutation: noopMutation,
    useGetCredentialsQuery: () => ({ data: { credentials: [] } }),
    usePerformConnectionActionMutation: noopMutation,
    useProcessConnectionRegistrationMutation: noopMutation,
    useUpdateConnectionByIdMutation: noopMutation,
  };
});

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

import { Keys } from '@meshery/schemas/permissions';
import { useConnectionWizard } from './useConnectionWizard';
import type { ConnectionWizardKindConfig } from '../ConnectionWizard.helpers';

const kubernetesKind = {
  kind: 'kubernetes',
  type: 'platform',
  subType: 'orchestrator',
  label: 'Kubernetes',
  description: 'Kubernetes cluster',
  flow: 'kubernetes',
  docsUrl: 'https://docs.meshery.io',
  connectionSchema: null,
  credentialSchema: null,
  svgColor: null,
  svgWhite: null,
} as unknown as ConnectionWizardKindConfig;

const renderPresetWizard = () =>
  renderHook(() =>
    useConnectionWizard({
      mode: 'create',
      isOpen: true,
      availableKinds: [kubernetesKind],
      presetKind: 'kubernetes',
      skipKindSelection: true,
    }),
  );

const renderConfigureWizard = () =>
  renderHook(() =>
    useConnectionWizard({
      mode: 'configure',
      isOpen: true,
      availableKinds: [kubernetesKind],
      initialKindConfig: kubernetesKind,
    }),
  );

describe('useConnectionWizard permission gating', () => {
  beforeEach(() => {
    canSpy.mockClear();
    canSpy.mockReturnValue(true);
  });

  it('selects the preset kind and advances past selection when permitted', () => {
    const { result } = renderPresetWizard();

    // `flow: 'kubernetes'` maps to the add-cluster capability.
    expect(canSpy).toHaveBeenCalledWith(Keys.LifecycleManagementAddCluster);
    expect(result.current.ctx.data.kindConfig?.kind).toBe('kubernetes');
    expect(result.current.activeIndex).toBe(1);
  });

  it('leaves the preset unapplied when the kind permission is denied', () => {
    canSpy.mockReturnValue(false);
    const { result } = renderPresetWizard();

    // The kind chooser stays in charge (it disables the kinds the user cannot
    // use) rather than the wizard pre-selecting a kind it will never let them
    // leave.
    expect(result.current.ctx.data.kindConfig).toBeNull();
    expect(result.current.activeIndex).toBe(0);
  });

  // Regression: permissions resolve asynchronously. Recording the preset as
  // "applied" while access was still denied used to latch it in, so the effect
  // short-circuited once permission arrived and an authorized user was never
  // advanced past "Choose Connection".
  it('applies the preset once a delayed permission resolves to granted', () => {
    canSpy.mockReturnValue(false);
    const { result, rerender } = renderPresetWizard();

    expect(result.current.activeIndex).toBe(0);
    expect(result.current.ctx.data.kindConfig).toBeNull();

    act(() => {
      canSpy.mockReturnValue(true);
      rerender();
    });

    expect(result.current.ctx.data.kindConfig?.kind).toBe('kubernetes');
    expect(result.current.activeIndex).toBe(1);
  });

  it('blocks progression when the already-selected kind is not permitted', () => {
    canSpy.mockReturnValue(false);
    const { result } = renderConfigureWizard();

    // Seeded by the caller rather than by the preset effect, so the wizard-level
    // guard - not the preset guard - is what has to stop it.
    expect(result.current.ctx.data.kindConfig?.kind).toBe('kubernetes');
    expect(result.current.canProceed).toBe(false);
  });

  it('allows progression for a permitted selected kind', () => {
    const { result } = renderConfigureWizard();

    expect(result.current.canProceed).toBe(true);
  });
});
