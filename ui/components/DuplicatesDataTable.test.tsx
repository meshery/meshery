import React from 'react';
import { render, screen, act } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const dataTableMock = vi.fn();
const getDuplicateModelsMock = vi.fn();
const getDuplicateComponentsMock = vi.fn();

vi.mock('@sistent/mui-datatables', () => ({
  default: (props: any) => {
    dataTableMock(props);
    return <div data-testid="datatable" data-count={props.options?.count} />;
  },
}));

vi.mock('@sistent/sistent', () => ({
  Tooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  TableCell: ({ children }: any) => <th>{children}</th>,
  TableSortLabel: ({ children }: any) => <span>{children}</span>,
}));

vi.mock('../utils/debounce', () => ({
  default: (fn: any) => fn,
}));

vi.mock('../api/meshmodel', () => ({
  getDuplicateModels: (...args: any[]) => getDuplicateModelsMock(...args),
  getDuplicateComponents: (...args: any[]) => getDuplicateComponentsMock(...args),
}));

vi.mock('../constants/navigator', () => ({
  MODELS: 'MODELS',
  COMPONENTS: 'COMPONENTS',
}));

import DuplicatesDataTable from './DuplicatesDataTable';

describe('DuplicatesDataTable', () => {
  beforeEach(() => {
    dataTableMock.mockClear();
    getDuplicateModelsMock.mockReset();
    getDuplicateComponentsMock.mockReset();
  });

  it('fetches duplicate models when view is MODELS', async () => {
    getDuplicateModelsMock.mockResolvedValue({ totalCount: 3, models: [{ name: 'm1' }] });

    await act(async () => {
      render(
        <DuplicatesDataTable
          view="MODELS"
          rowData={{ model: 'kafka', version: '1.0' }}
          classes={{ tableHeader: 'header' }}
        />,
      );
    });

    expect(getDuplicateModelsMock).toHaveBeenCalledWith('kafka', '1.0');
    expect(screen.getByTestId('datatable')).toHaveAttribute('data-count', '3');
  });

  it('fetches duplicate components when view is COMPONENTS', async () => {
    getDuplicateComponentsMock.mockResolvedValue({
      totalCount: 2,
      components: [{ apiVersion: 'v1', kind: 'Pod' }],
    });

    await act(async () => {
      render(
        <DuplicatesDataTable
          view="COMPONENTS"
          rowData={{ kind: 'Pod', model: 'core', version: 'v1' }}
          classes={{ tableHeader: 'header' }}
        />,
      );
    });

    // getDuplicateComponents is invoked as (componentKind, modelName, apiVersion)
    // even though useEffect chains pass (kind, model, version)
    expect(getDuplicateComponentsMock).toHaveBeenCalledWith('Pod', 'v1', 'core');
  });

  it('passes serverSide options to MUIDataTable', async () => {
    getDuplicateModelsMock.mockResolvedValue({ totalCount: 0, models: [] });
    await act(async () => {
      render(
        <DuplicatesDataTable
          view="MODELS"
          rowData={{ model: 'foo', version: '1' }}
          classes={{ tableHeader: 'header' }}
        />,
      );
    });
    const props = dataTableMock.mock.calls.at(-1)![0];
    expect(props.options.serverSide).toBe(true);
    expect(props.options.rowsPerPage).toBe(10);
  });
});
