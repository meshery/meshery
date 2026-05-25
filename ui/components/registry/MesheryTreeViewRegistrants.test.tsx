import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import MesheryTreeViewRegistrants from './MesheryTreeViewRegistrants';

vi.mock('@sistent/sistent', () => ({
  CircularProgress: () => <div data-testid="loading-spinner" />,
}));

vi.mock('../shared/TreeView', () => ({
  SimpleTreeView: ({ children }: any) => <div data-testid="tree-view">{children}</div>,
}));

vi.mock('@/constants/navigator', () => ({
  REGISTRANTS: 'Registrants',
}));

vi.mock('../../assets/icons/MinusSquare', () => ({ default: () => <svg /> }));
vi.mock('../../assets/icons/PlusSquare', () => ({ default: () => <svg /> }));
vi.mock('../../assets/icons/DotSquare', () => ({ default: () => <svg /> }));

vi.mock('./StyledTreeItem', () => ({
  default: ({ itemId, labelText, onClick, children }: any) => (
    <div data-testid={`styled-tree-item-${itemId}`}>
      <button onClick={onClick} data-testid={`styled-button-${itemId}`}>
        {labelText}
      </button>
      {children}
    </div>
  ),
}));

vi.mock('./MesheryTreeViewItem', () => ({
  default: ({ modelDef, registrantID }: any) => (
    <div data-testid={`model-${registrantID}-${modelDef.id}`}>{modelDef.name}</div>
  ),
}));

describe('MesheryTreeViewRegistrants', () => {
  const baseProps = {
    handleToggle: vi.fn(),
    handleSelect: vi.fn(),
    expanded: [],
    selected: [],
    setShowDetailsData: vi.fn(),
    lastRegistrantRef: { current: null } as React.MutableRefObject<any>,
    isRegistrantFetching: false,
    showDetailsData: { type: '', data: {} },
  };

  const sampleData = [
    {
      id: 'reg-1',
      name: 'Helm',
      kind: 'helm',
      summary: { models: 3 },
      models: [
        { id: 'm1', name: 'model-one', registrant: { kind: 'helm' } },
        { id: 'm2', name: 'model-two', registrant: { kind: 'helm' } },
        { id: 'm3', name: 'model-three', registrant: { kind: 'other' } },
      ],
    },
    {
      id: 'reg-2',
      name: 'GitHub',
      kind: 'github',
      models: [],
    },
    {
      id: 'reg-3',
      name: 'NoSummaryOrModels',
      kind: 'other',
    },
  ];

  it('renders registrants that have summary or models', () => {
    render(<MesheryTreeViewRegistrants {...baseProps} data={sampleData} />);

    expect(screen.getByTestId('styled-tree-item-reg-1')).toBeInTheDocument();
    expect(screen.getByTestId('styled-tree-item-reg-2')).toBeInTheDocument();
    expect(screen.queryByTestId('styled-tree-item-reg-3')).not.toBeInTheDocument();
  });

  it('shows "Models (N)" with the count from summary, otherwise "Models (0)"', () => {
    render(<MesheryTreeViewRegistrants {...baseProps} data={sampleData} />);

    expect(screen.getByTestId('styled-button-reg-1.1')).toHaveTextContent('Models (3)');
    expect(screen.getByTestId('styled-button-reg-2.1')).toHaveTextContent('Models (0)');
  });

  it('only renders nested models whose registrant.kind matches the registrant kind', () => {
    render(<MesheryTreeViewRegistrants {...baseProps} data={sampleData} />);

    expect(screen.getByTestId('model-reg-1-m1')).toBeInTheDocument();
    expect(screen.getByTestId('model-reg-1-m2')).toBeInTheDocument();
    expect(screen.queryByTestId('model-reg-1-m3')).not.toBeInTheDocument();
  });

  it('clicking a registrant calls setShowDetailsData with REGISTRANTS type', async () => {
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();

    render(
      <MesheryTreeViewRegistrants
        {...baseProps}
        data={sampleData.slice(0, 1)}
        setShowDetailsData={setShowDetailsData}
      />,
    );

    await user.click(screen.getByTestId('styled-button-reg-1'));

    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'Registrants',
      data: sampleData[0],
    });
  });

  it('renders loading spinner when fetching', () => {
    render(<MesheryTreeViewRegistrants {...baseProps} data={[]} isRegistrantFetching={true} />);
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('handles empty data array gracefully', () => {
    render(<MesheryTreeViewRegistrants {...baseProps} data={[]} />);
    expect(screen.getByTestId('tree-view')).toBeInTheDocument();
  });

  it('handles data with undefined values safely', () => {
    const data = [undefined as any, null as any];
    render(<MesheryTreeViewRegistrants {...baseProps} data={data} />);
    expect(screen.getByTestId('tree-view')).toBeInTheDocument();
  });
});
