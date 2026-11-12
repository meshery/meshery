import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
let queryReturn: any = {
  data: { relationships: [] },
  isLoading: false,
  isError: false,
  error: null,
};
let routerSelectedItemUUID = '';

vi.mock('@sistent/sistent', () => ({
  CircularProgress: () => <div data-testid="loading-spinner" />,
}));

vi.mock('../../constants/navigator', () => ({
  MODELS: 'Models',
  RELATIONSHIPS: 'Relationships',
}));

vi.mock('./StyledTreeItem', () => ({
  default: ({ itemId, labelText, children }: any) => (
    <div data-testid={`tree-item-${itemId}`}>
      <span data-testid={`label-${itemId}`}>{labelText}</span>
      {children}
    </div>
  ),
}));

vi.mock('./helper', () => ({
  groupRelationshipsByKind: (rels: any[]) => {
    const map: Record<string, any> = {};
    for (const r of rels) {
      if (!map[r.kind]) {
        map[r.kind] = { kind: r.kind, relationships: [] };
      }
      map[r.kind].relationships.push(r);
    }
    return Object.values(map);
  },
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useGetRelationshipsFromModalQuery: () => queryReturn,
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { ERROR: 'error' },
}));

vi.mock('./hooks', () => ({
  useRegistryRouter: () => ({ selectedItemUUID: routerSelectedItemUUID }),
}));

vi.mock('./RelationshipTree', () => ({
  default: ({ data, view, idForKindAsProp }: any) => (
    <div data-testid="relationship-tree" data-view={view} data-id-prop={idForKindAsProp}>
      {data.length} groups
    </div>
  ),
}));

import VersionedModelRelationshipTree from './VersionedModelRelationshipTree';

describe('VersionedModelRelationshipTree', () => {
  const baseProps = {
    modelDef: { id: 'm1' },
    versionedModelDef: { id: 'v1', name: 'mname', model: { version: '1.0.0' } },
    setShowDetailsData: vi.fn(),
    showDetailsData: { type: '', data: {} },
    handleToggle: vi.fn(),
    handleSelect: vi.fn(),
    selected: [],
    expanded: [],
  };

  beforeEach(() => {
    queryReturn = {
      data: { relationships: [] },
      isLoading: false,
      isError: false,
      error: null,
    };
    routerSelectedItemUUID = '';
    notify.mockReset();
  });

  it('renders the loading spinner while loading', () => {
    queryReturn = {
      data: { relationships: [] },
      isLoading: true,
      isError: false,
      error: null,
    };
    render(<VersionedModelRelationshipTree {...baseProps} />);
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('shows "Relationships (0)" when no relationships exist', () => {
    render(<VersionedModelRelationshipTree {...baseProps} />);

    expect(screen.getByTestId('label-m1.v1.2')).toHaveTextContent('Relationships (0)');
    expect(screen.queryByTestId('relationship-tree')).not.toBeInTheDocument();
  });

  it('renders the RelationshipTree when there are relationships', () => {
    queryReturn = {
      data: {
        relationships: [
          { id: 'r1', kind: 'edge', subType: 's1', model: { name: 'k' } },
          { id: 'r2', kind: 'hierarchical', subType: 's2', model: { name: 'k' } },
        ],
      },
      isLoading: false,
      isError: false,
      error: null,
    };

    render(<VersionedModelRelationshipTree {...baseProps} />);

    const tree = screen.getByTestId('relationship-tree');
    expect(tree).toHaveAttribute('data-view', 'Models');
    expect(tree).toHaveAttribute('data-id-prop', 'm1.v1.2');
    expect(tree).toHaveTextContent('2 groups');
  });

  it('shows the right label count', () => {
    queryReturn = {
      data: {
        relationships: [{ id: 'r1', kind: 'edge', subType: 's', model: { name: 'k' } }],
      },
      isLoading: false,
      isError: false,
      error: null,
    };
    render(<VersionedModelRelationshipTree {...baseProps} />);
    expect(screen.getByTestId('label-m1.v1.2')).toHaveTextContent('Relationships (1)');
  });

  it('notifies when isError is true', () => {
    queryReturn = {
      data: { relationships: [] },
      isLoading: false,
      isError: true,
      error: { data: 'boom' },
    };
    render(<VersionedModelRelationshipTree {...baseProps} />);
    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ event_type: 'error' }));
  });

  it('prefixes the registrantID into tree item itemId', () => {
    queryReturn = {
      data: { relationships: [] },
      isLoading: false,
      isError: false,
      error: null,
    };
    render(<VersionedModelRelationshipTree {...baseProps} registrantID="reg-1" />);
    expect(screen.getByTestId('tree-item-reg-1.1.m1.v1.2')).toBeInTheDocument();
  });
});
