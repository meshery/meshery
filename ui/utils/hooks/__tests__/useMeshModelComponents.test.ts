import { describe, expect, it, vi } from 'vitest';
import { act, renderHook } from '@testing-library/react';

vi.mock('../../../api/meshmodel', () => ({
  fetchCategories: vi.fn(() => Promise.resolve({ categories: [] })),
  getComponentFromModelApi: vi.fn(),
  getModelFromCategoryApi: vi.fn(),
  getVersionedComponentFromModel: vi.fn(),
}));

import { getComponentFromModelApi } from '../../../api/meshmodel';
import {
  removeDuplicatesAndSortByAlphabet,
  useMeshModelComponents,
} from '../useMeshModelComponents';

// The registry API returns every registration of a component, so the same kind
// can arrive more than once as separate objects (see FindDuplicateComponents).
const componentDef = (kind: string, version: string, modelVersion = 'v1.0.0') => ({
  component: { kind, version },
  model: { name: 'kubernetes', model: { version: modelVersion } },
});

describe('removeDuplicatesAndSortByAlphabet', () => {
  it('keeps one entry per component kind', () => {
    const result = removeDuplicatesAndSortByAlphabet([
      componentDef('Pod', 'v1'),
      componentDef('Pod', 'v1'),
      componentDef('Service', 'v1'),
    ]);

    expect(result.map((c) => c.component.kind)).toEqual(['Pod', 'Service']);
  });

  it('sets apiVersion to the most recent version of that kind', () => {
    const result = removeDuplicatesAndSortByAlphabet([
      componentDef('HorizontalPodAutoscaler', 'v1beta1'),
      componentDef('HorizontalPodAutoscaler', 'v1'),
    ]);

    expect(result).toHaveLength(1);
    expect(result[0].apiVersion).toBe('v1');
  });

  it('sorts unique components by kind', () => {
    const result = removeDuplicatesAndSortByAlphabet([
      componentDef('Service', 'v1'),
      componentDef('ConfigMap', 'v1'),
      componentDef('Pod', 'v1'),
    ]);

    expect(result.map((c) => c.component.kind)).toEqual(['ConfigMap', 'Pod', 'Service']);
  });
});

describe('useMeshModelComponents', () => {
  it('lists a component registered twice only once', async () => {
    vi.mocked(getComponentFromModelApi).mockResolvedValue({
      components: [componentDef('Pod', 'v1'), componentDef('Pod', 'v1')],
    });
    const { result } = renderHook(() => useMeshModelComponents());

    await act(async () => {
      await result.current.getComponentsFromModel('kubernetes');
    });

    const [versionGroup] = result.current.meshmodelComponents['kubernetes'];
    expect(versionGroup.components).toHaveLength(1);
  });
});
