import { act, renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@/rtk-query/design', () => ({
  useUpdatePatternFileMutation: () => [vi.fn()],
  useDeletePatternFileMutation: () => [vi.fn()],
  useLazyGetDesignQuery: () => [vi.fn()],
}));

vi.mock('../../../../../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));

vi.mock('../../../../../lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

import useDesignLifecycle from './useDesignLifecycle';

const makeDefinition = (kind: string) => ({
  component: { kind, version: 'v1' },
  schemaVersion: 'components.meshery.io/v1beta2',
  version: 'v1.0.0',
  model: {
    name: 'kubernetes',
    version: 'v1.0.0',
    category: { name: 'Orchestration & Management' },
    registrant: { kind: 'github' },
  },
});

// Mirrors the reference PatternServiceFormCore fills in for each mounted form.
const makeReference = (referKey: string) => ({ current: { referKey } });

describe('useDesignLifecycle onSettingsChange', () => {
  it('writes edits from two form references to two separate components', () => {
    const { result } = renderHook(() => useDesignLifecycle());
    const onNamespaceChange = result.current.onSettingsChange(
      makeDefinition('Namespace'),
      makeReference('ns-id'),
    );
    const onServiceChange = result.current.onSettingsChange(
      makeDefinition('Service'),
      makeReference('svc-id'),
    );

    act(() => onNamespaceChange({ name: 'web-ns', namespace: 'default' }));
    act(() => onServiceChange({ name: 'web-svc', namespace: 'default' }));

    const { components } = result.current.designJson;
    expect(components).toHaveLength(2);
    expect(components.map((c) => [c.id, c.component.kind, c.displayName])).toEqual([
      ['ns-id', 'Namespace', 'web-ns'],
      ['svc-id', 'Service', 'web-svc'],
    ]);
  });

  it('updates the existing component in place when the same reference edits again', () => {
    const { result } = renderHook(() => useDesignLifecycle());
    const onNamespaceChange = result.current.onSettingsChange(
      makeDefinition('Namespace'),
      makeReference('ns-id'),
    );

    act(() => onNamespaceChange({ name: 'web-ns', namespace: 'default' }));
    act(() => onNamespaceChange({ name: 'web-ns-edited', namespace: 'apps' }));

    const { components } = result.current.designJson;
    expect(components).toHaveLength(1);
    expect(components[0].id).toBe('ns-id');
    expect(components[0].displayName).toBe('web-ns-edited');
    expect(components[0].configuration.metadata.namespace).toBe('apps');
  });
});
