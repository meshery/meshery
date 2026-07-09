import React from 'react';
import { render, screen, act, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';

// Shared mock state. Declared via vi.hoisted so the (hoisted) vi.mock factories
// below can reference it without tripping the "used before initialization" rule.
const h = vi.hoisted(() => ({
  notifyMock: vi.fn(),
  unwrapMock: vi.fn(),
  addTriggerMock: vi.fn(),
  // Resolver for the first PromptComponent.show() call (the upload modal), kept
  // pending so the test can pick a file before "IMPORT" is confirmed — mirroring
  // the real user flow.
  resolveUploadModal: { current: undefined as undefined | ((value: string) => void) },
  // Counts show() calls so only the first (upload) modal stays pending. Held in
  // hoisted state (not a factory-local) so beforeEach can reset it per test.
  showCount: { current: 0 },
}));

vi.mock('../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: h.notifyMock }),
}));
vi.mock('../rtk-query/connection', () => ({
  useAddKubernetesConfigMutation: () => [h.addTriggerMock],
}));
vi.mock('@/utils/can', () => ({ default: () => true }));

vi.mock('@/utils/hooks/useKubernetesHook', () => ({ default: () => vi.fn() }));
vi.mock('@/utils/hooks/useTestIDs', () => ({ default: () => () => 'test-id' }));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));
vi.mock('../assets/icons/AddIconCircleBorder', () => ({ default: () => null }));
vi.mock('../lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'error', WARNING: 'warning', SUCCESS: 'success', INFO: 'info' },
}));
vi.mock('./connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: () => null,
  ConnectionStateChip: () => null,
}));

// PromptComponent is imperative (ref.show()). The first call (the kubeconfig
// upload prompt) stays pending until the test resolves it; it also renders the
// passed subtitle so the hidden <input id="k8sfile"> mounts and can be driven.
vi.mock('./PromptComponent', async () => {
  const { forwardRef, useState, useImperativeHandle } = await import('react');
  const PromptComponent = forwardRef((_props: unknown, ref: React.Ref<unknown>) => {
    const [content, setContent] = useState<React.ReactNode>(null);
    useImperativeHandle(ref, () => ({
      show: (args: { subtitle?: React.ReactNode }) => {
        setContent(args?.subtitle ?? null);
        h.showCount.current += 1;
        if (h.showCount.current === 1) {
          return new Promise<string>((resolve) => {
            h.resolveUploadModal.current = resolve;
          });
        }
        return Promise.resolve('OK');
      },
    }));
    return <div data-testid="prompt">{content}</div>;
  });
  return { __esModule: true, default: PromptComponent };
});

vi.mock('@sistent/sistent', () => ({
  Button: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
  CloudUploadIcon: () => null,
  Typography: ({ children }: any) => <span>{children}</span>,
  FormGroup: ({ children }: any) => <div>{children}</div>,
  TextField: ({ id }: any) => <input id={id} readOnly />,
  InputAdornment: ({ children }: any) => <span>{children}</span>,
  Tooltip: ({ children }: any) => <>{children}</>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  Box: ({ children }: any) => <div>{children}</div>,
  styled: () => () => ({}),
  PROMPT_VARIANTS: { SUCCESS: 'success', WARNING: 'warning', ERROR: 'error' },
}));

import MesherySettingsEnvButtons from './MesherySettingsEnvButtons';

const selectKubeconfig = () => {
  const fileInput = document.getElementById('k8sfile') as HTMLInputElement;
  expect(fileInput).toBeTruthy();
  const file = new File(['apiVersion: v1'], 'config.yaml', { type: 'text/yaml' });
  Object.defineProperty(fileInput, 'files', { value: [file], configurable: true });
  fireEvent.change(fileInput);
};

describe('MesherySettingsEnvButtons – kubeconfig upload', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    h.resolveUploadModal.current = undefined;
    h.showCount.current = 0;
    h.addTriggerMock.mockReturnValue({ unwrap: h.unwrapMock });
  });

  it('surfaces errored contexts from a camelCase response without crashing', async () => {
    // Mirrors the server wire format: camelCase buckets, and errored contexts are
    // plain K8sContext objects with NO per-context `error` field (the failure is
    // recorded server-side in event metadata). The previous implementation read
    // snake_case keys and called `ctx.error.toString()`, throwing
    // "Cannot read properties of undefined (reading 'toString')", which surfaced
    // as a misleading "failed to upload kubernetes config" toast.
    h.unwrapMock.mockResolvedValue({
      registeredContexts: [],
      connectedContexts: [],
      ignoredContexts: [],
      erroredContexts: [{ name: 'broken-ctx', server: 'https://broken.example' }],
    });

    render(<MesherySettingsEnvButtons />);

    await act(async () => {
      fireEvent.click(screen.getByText('Add Cluster'));
    });

    selectKubeconfig();

    await act(async () => {
      h.resolveUploadModal.current?.('IMPORT');
    });

    await waitFor(() =>
      expect(h.notifyMock).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Failed to add cluster "broken-ctx" at https://broken.example',
          event_type: 'error',
        }),
      ),
    );

    const messages = h.notifyMock.mock.calls.map((call) => call[0]?.message);
    // The catch-all upload-failure toast (the symptom of the old crash) must not appear.
    expect(
      messages.some(
        (message) =>
          typeof message === 'string' && message.includes('failed to upload kubernetes config'),
      ),
    ).toBe(false);
  });
});
