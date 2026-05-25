import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const muiDataTableProps: any[] = [];

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    if (typeof Component === 'string') {
      const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
      Styled.displayName = 'StyledHostMock';
      return Styled;
    }
    const Styled = ({ children, ...rest }: any) => <Component {...rest}>{children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    NoSsr: ({ children }: any) => <>{children}</>,
    Chip: ({ label }: any) => <span data-testid="chip">{label}</span>,
    Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  };
});

vi.mock('@sistent/mui-datatables', () => ({
  default: (props: any) => {
    muiDataTableProps.push(props);
    return (
      <div data-testid="mui-datatable" data-rows={props.data.length}>
        {props.data.map((row: any, idx: number) => (
          <div key={idx} data-testid="row">
            <span>{row.board}</span>
            {row.panels}
            {row.template_variables}
          </div>
        ))}
      </div>
    );
  },
}));

import GrafanaDisplaySelection from './GrafanaDisplaySelection';

describe('GrafanaDisplaySelection', () => {
  it('renders a data table row per board config', () => {
    muiDataTableProps.length = 0;
    render(
      <GrafanaDisplaySelection
        classes={{}}
        boardPanelConfigs={[
          { board: { title: 'Board A' }, panels: [{ id: 1, title: 'cpu' }], templateVars: ['k=v'] },
          { board: { title: 'Board B' }, panels: [{ id: 2, title: 'mem' }], templateVars: [] },
        ]}
        deleteSelectedBoardPanelConfig={vi.fn()}
      />,
    );
    expect(screen.getByTestId('mui-datatable')).toHaveAttribute('data-rows', '2');
    expect(screen.getAllByTestId('row')).toHaveLength(2);
    const labels = screen.getAllByTestId('chip').map((c) => c.textContent);
    expect(labels).toContain('cpu');
    expect(labels).toContain('mem');
    expect(labels).toContain('k=v');
  });

  it('passes pagination=false / download=false / print=false through the options', () => {
    muiDataTableProps.length = 0;
    render(
      <GrafanaDisplaySelection
        classes={{}}
        boardPanelConfigs={[]}
        deleteSelectedBoardPanelConfig={vi.fn()}
      />,
    );
    const lastProps = muiDataTableProps.at(-1);
    expect(lastProps.options.pagination).toBe(false);
    expect(lastProps.options.print).toBe(false);
    expect(lastProps.options.download).toBe(false);
    expect(lastProps.options.viewColumns).toBe(false);
  });

  it('delegates row deletion to the deleteSelectedBoardPanelConfig callback', () => {
    muiDataTableProps.length = 0;
    const del = vi.fn();
    render(
      <GrafanaDisplaySelection
        classes={{}}
        boardPanelConfigs={[]}
        deleteSelectedBoardPanelConfig={del}
      />,
    );
    const lastProps = muiDataTableProps.at(-1);
    const result = lastProps.options.onRowsDelete({ data: [{ dataIndex: 0 }, { dataIndex: 2 }] });
    expect(del).toHaveBeenCalledWith([0, 2]);
    expect(result).toBe(false);
  });
});
