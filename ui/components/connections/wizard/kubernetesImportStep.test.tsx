import React from 'react';
import { act, render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

// ---------------------------------------------------------------------------
// @sistent/sistent — styled-component factory + the specific primitives used
// by KubernetesImportStep/KubeconfigStepBody.
//
// Mirrors the pattern in styles.test.tsx (styledFactory) and
// ConnectionStateTransitionModal.test.tsx (named-component stubs).
// ---------------------------------------------------------------------------

vi.mock('@sistent/sistent', () => {
  type AnyProps = React.HTMLAttributes<HTMLElement> & { children?: React.ReactNode };

  return {
    Alert: ({ children }: AnyProps) => <div role="alert">{children}</div>,
    Box: ({ children, ...props }: AnyProps) => <div {...props}>{children}</div>,
    CircularProgress: () => <span data-testid="preview-spinner" />,
    CloudUploadIcon: () => <svg data-testid="cloud-upload-icon" />,
    Typography: ({ children }: AnyProps) => <span>{children}</span>,
  };
});

// ---------------------------------------------------------------------------
// @/theme — alpha is only used inside styled-factory callbacks, which the mock
// above never invokes. Matches the minimal stub in styles.test.tsx line 55-57.
// ---------------------------------------------------------------------------

vi.mock('@/theme', () => {
  type AnyProps = React.HTMLAttributes<HTMLElement> & { children?: React.ReactNode };

  // styled(tag)(styleFactory) → passthrough rendered as the underlying element.
  // Same two-call curried shape as used in Filters.styled.test.tsx.
  const styledFactory = (Component: React.ComponentType<AnyProps> | string) => () => {
    const el = typeof Component === 'string' ? Component : 'div';
    const Passthrough = ({ children, ...props }: AnyProps) =>
      React.createElement(el as 'div', props as React.HTMLAttributes<HTMLElement>, children);
    Passthrough.displayName = 'StyledMock';
    return Passthrough;
  };

  return {
    styled: styledFactory,
    alpha: (color: string, value: number) => `alpha(${color},${value})`,
  };
});

// ---------------------------------------------------------------------------
// Transitive deps imported at the module level by kubernetesImportStep.tsx
// that are not exercised by the race-guard test but must resolve cleanly.
// ---------------------------------------------------------------------------

vi.mock('../ConnectionWizardStepContent', () => ({
  StepLayout: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  StepHeader: () => null,
}));

vi.mock('./kubernetesReviewStep', () => ({
  defaultChoice: (name: string) => ({ selected: true, name, meshsyncDeploymentMode: 'operator' }),
}));

vi.mock('./errors', () => ({
  formatWizardError: (err: unknown) => String(err),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'error', WARNING: 'warning' },
}));

// ---------------------------------------------------------------------------
// Import under test (must come after vi.mock calls — vitest hoists them).
// ---------------------------------------------------------------------------

import { kubernetesDetailsStep } from './kubernetesImportStep';
import type { WizardContext, DiscoveredKubeContext } from './types';

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

/** Returns a promise together with its resolve/reject handles so individual
 * test steps can release it at a controlled moment. */
function createDeferred<T>() {
  let resolve!: (value: T | PromiseLike<T>) => void;
  let reject!: (reason?: unknown) => void;
  const promise = new Promise<T>((res, rej) => {
    resolve = res;
    reject = rej;
  });
  return { promise, resolve, reject };
}

/** Builds the smallest WizardContext that KubeconfigStepBody actually reads.
 *
 * `data` is a plain object mutated by the `patch`/`patchPostConfig` stubs so
 * that subsequent reads inside the same async flow see the updated values —
 * the same pattern used by ConnectionStateTransitionModal.test.tsx's
 * `mockStoreState` object. Neither Redux nor React state is involved. */
function makeMockCtx(discoverKubeContexts: WizardContext['services']['discoverKubeContexts']): {
  ctx: WizardContext;
  patchPostConfig: ReturnType<typeof vi.fn>;
} {
  const data = {
    kubeconfigFile: null as File | null,
    postConfig: {} as Record<string, unknown>,
  } as unknown as WizardContext['data'];

  const patchPostConfig = vi.fn((partial: Record<string, unknown>) => {
    Object.assign(data.postConfig, partial);
  });

  const ctx: WizardContext = {
    mode: 'create',
    data,
    patch: vi.fn((partial) => Object.assign(data, partial)),
    patchPostConfig,
    services: {
      discoverKubeContexts,
      notify: vi.fn(),
      registerConnection: vi.fn(),
      connectConnection: vi.fn(),
      uploadKubeconfig: vi.fn(),
      updateConnectionById: vi.fn(),
      setMeshsyncMode: vi.fn(),
      flushMeshsync: vi.fn(),
      credentials: [],
    },
    formRefs: {
      connection: { current: null },
      credential: { current: null },
    },
    advance: vi.fn(),
  } as unknown as WizardContext;

  return { ctx, patchPostConfig };
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

const { Component: KubeconfigStepBody } = kubernetesDetailsStep;

describe('KubeconfigStepBody preview race-guard (requestVersionRef)', () => {
  // Regression guard for the CodeRabbit fix: `requestVersion` must be
  // incremented before the early return on null-file so that any in-flight
  // discovery promise from a previous pick cannot write its stale result.
  it('discards a stale preview result when the file is cleared before discovery resolves', async () => {
    const fileA = new File(['apiVersion: v1\nclusters: []'], 'cluster-a.yaml', {
      type: 'application/yaml',
    });

    const staleResult: DiscoveredKubeContext[] = [
      { id: 'ctx-a', name: 'cluster-a', server: 'https://a.example.com', reachable: true },
    ];

    // Deferred discovery: the test controls when file A's request settles.
    const deferred = createDeferred<DiscoveredKubeContext[]>();
    const discoverKubeContexts = vi.fn(() => deferred.promise);

    const { ctx, patchPostConfig } = makeMockCtx(
      discoverKubeContexts as unknown as WizardContext['services']['discoverKubeContexts'],
    );

    const { container } = render(<KubeconfigStepBody ctx={ctx} />);

    const fileInput = container.querySelector(
      '#connection-wizard-kubeconfig-input',
    ) as HTMLInputElement;
    expect(fileInput).toBeTruthy();

    // Step 1: pick file A — requestVersion becomes 1, discovery starts in the background.
    await act(async () => {
      Object.defineProperty(fileInput, 'files', { value: [fileA], configurable: true });
      fileInput.dispatchEvent(new Event('change', { bubbles: true }));
    });

    expect(discoverKubeContexts).toHaveBeenCalledTimes(1);
    expect(discoverKubeContexts).toHaveBeenCalledWith(fileA);

    // Step 2: clear the file before the discovery resolves — requestVersion
    // becomes 2. The early-return path executes without awaiting anything.
    await act(async () => {
      Object.defineProperty(fileInput, 'files', { value: [], configurable: true });
      fileInput.dispatchEvent(new Event('change', { bubbles: true }));
    });

    // Step 3: resolve the now-stale discovery with real-looking data.
    await act(async () => {
      deferred.resolve(staleResult);
    });

    // The version mismatch (1 !== 2) must have caused handlePickFile to exit
    // before calling patchPostConfig with the preview contexts.
    const callsWithPreviewData = patchPostConfig.mock.calls.filter(
      ([partial]) => partial.previewContexts !== undefined,
    );
    expect(callsWithPreviewData).toHaveLength(0);
  });

  it('writes preview contexts when discovery resolves before the file is changed again', async () => {
    const fileA = new File(['apiVersion: v1\nclusters: []'], 'cluster-a.yaml', {
      type: 'application/yaml',
    });

    const validResult: DiscoveredKubeContext[] = [
      { id: 'ctx-a', name: 'cluster-a', server: 'https://a.example.com', reachable: true },
    ];

    const discoverKubeContexts = vi.fn(() => Promise.resolve(validResult));

    const { ctx, patchPostConfig } = makeMockCtx(
      discoverKubeContexts as unknown as WizardContext['services']['discoverKubeContexts'],
    );

    const { container } = render(<KubeconfigStepBody ctx={ctx} />);

    const fileInput = container.querySelector(
      '#connection-wizard-kubeconfig-input',
    ) as HTMLInputElement;

    // Pick a file and let discovery resolve immediately.
    await act(async () => {
      Object.defineProperty(fileInput, 'files', { value: [fileA], configurable: true });
      fileInput.dispatchEvent(new Event('change', { bubbles: true }));
    });

    // The request version still matches — preview must have been written.
    const callsWithPreviewData = patchPostConfig.mock.calls.filter(
      ([partial]) => partial.previewContexts !== undefined,
    );
    expect(callsWithPreviewData).toHaveLength(1);
    expect(callsWithPreviewData[0][0].previewContexts).toEqual(validResult);
  });
});
