import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const notify = vi.fn();
let queryReturn: any = {
  data: { components: [] },
  isLoading: false,
  isError: false,
  error: null,
};
let routerSelectedItemUUID = '';

vi.mock('@sistent/sistent', () => ({
  CircularProgress: () => <div data-testid="loading-spinner" />,
}));

vi.mock('../../constants/navigator', () => ({
  COMPONENTS: 'Components',
}));

vi.mock('./StyledTreeItem', () => ({
  default: ({ itemId, labelText, onClick, children }: any) => (
    <div data-testid={`tree-item-${itemId}`}>
      <button onClick={onClick} data-testid={`button-${itemId}`}>
        {labelText}
      </button>
      {children}
    </div>
  ),
}));

vi.mock('./helper', () => ({
  getFilteredDataForDetailsComponent: (data: any[], id: string) => ({
    type: 'Components',
    data: data.find((d) => d.id === id) || {},
  }),
}));

vi.mock('@/rtk-query/meshModel', () => ({
  useGetComponentsFromModalQuery: () => queryReturn,
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

import VersionedModelComponentTree from './VersionedModelComponentTree';

describe('VersionedModelComponentTree', () => {
  const baseProps = {
    modelDef: { id: 'm1' },
    versionedModelDef: { id: 'v1', name: 'mname', model: { version: '1.0.0' } },
    setShowDetailsData: vi.fn(),
    showDetailsData: { type: '', data: {} },
  };

  beforeEach(() => {
    queryReturn = {
      data: { components: [] },
      isLoading: false,
      isError: false,
      error: null,
    };
    routerSelectedItemUUID = '';
    notify.mockReset();
    baseProps.setShowDetailsData = vi.fn();
  });

  it('renders the loading spinner when isLoading is true', () => {
    queryReturn = {
      data: { components: [] },
      isLoading: true,
      isError: false,
      error: null,
    };
    render(<VersionedModelComponentTree {...baseProps} />);
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('renders a Components header with the count', () => {
    queryReturn = {
      data: {
        components: [
          { id: 'c1', displayName: 'Comp One' },
          { id: 'c2', displayName: 'Comp Two' },
        ],
      },
      isLoading: false,
      isError: false,
      error: null,
    };

    render(<VersionedModelComponentTree {...baseProps} />);

    expect(screen.getByTestId('button-m1.v1.1')).toHaveTextContent('Components (2)');
  });

  it('shows "Components (0)" when no components are returned', () => {
    queryReturn = {
      data: { components: [] },
      isLoading: false,
      isError: false,
      error: null,
    };

    render(<VersionedModelComponentTree {...baseProps} />);

    expect(screen.getByTestId('button-m1.v1.1')).toHaveTextContent('Components (0)');
  });

  it('renders one tree item per component and fires setShowDetailsData on click', async () => {
    queryReturn = {
      data: {
        components: [{ id: 'c1', displayName: 'Comp One' }],
      },
      isLoading: false,
      isError: false,
      error: null,
    };
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();

    render(<VersionedModelComponentTree {...baseProps} setShowDetailsData={setShowDetailsData} />);

    expect(screen.getByText('Comp One')).toBeInTheDocument();
    await user.click(screen.getByTestId('button-m1.v1.1.c1'));
    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'Components',
      data: queryReturn.data.components[0],
    });
  });

  it('notifies on isError', () => {
    queryReturn = {
      data: { components: [] },
      isLoading: false,
      isError: true,
      error: { data: 'boom' },
    };

    render(<VersionedModelComponentTree {...baseProps} />);

    expect(notify).toHaveBeenCalledWith(expect.objectContaining({ event_type: 'error' }));
  });

  it('syncs showDetailsData when route UUID matches model.version path', () => {
    queryReturn = {
      data: { components: [{ id: 'c1', displayName: 'X' }] },
      isLoading: false,
      isError: false,
      error: null,
    };
    routerSelectedItemUUID = 'm1.v1.1.c1';
    const setShowDetailsData = vi.fn();

    render(
      <VersionedModelComponentTree
        {...baseProps}
        setShowDetailsData={setShowDetailsData}
        showDetailsData={{ type: '', data: {} }}
      />,
    );

    expect(setShowDetailsData).toHaveBeenCalled();
  });

  it('includes the registrantID prefix in the tree item itemId', () => {
    queryReturn = {
      data: { components: [{ id: 'c1', displayName: 'Comp' }] },
      isLoading: false,
      isError: false,
      error: null,
    };

    render(<VersionedModelComponentTree {...baseProps} registrantID="reg-1" />);

    expect(screen.getByTestId('tree-item-reg-1.1.m1.v1.1')).toBeInTheDocument();
  });
});
