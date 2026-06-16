import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import RelationshipTree from './RelationshipTree';

vi.mock('@sistent/sistent', () => ({
  CircularProgress: (props: any) => <div data-testid="loading-spinner" data-color={props.color} />,
}));

vi.mock('../shared/TreeView', () => ({
  SimpleTreeView: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('@/constants/navigator', () => ({
  RELATIONSHIPS: 'Relationships',
}));

vi.mock('../../assets/icons/MinusSquare', () => ({
  default: () => <svg data-testid="minus" />,
}));

vi.mock('../../assets/icons/PlusSquare', () => ({
  default: () => <svg data-testid="plus" />,
}));

vi.mock('../../assets/icons/DotSquare', () => ({
  default: () => <svg data-testid="dot" />,
}));

vi.mock('./StyledTreeItem', () => ({
  default: ({ itemId, labelText, onClick, children, 'data-id': dataId }: any) => (
    <div data-testid={`tree-item-${itemId}`} data-id={dataId} onClick={onClick}>
      <button data-testid={`tree-item-button-${itemId}`}>{labelText}</button>
      {children}
    </div>
  ),
}));

describe('RelationshipTree', () => {
  const baseProps = {
    expanded: [],
    selected: [],
    handleToggle: vi.fn(),
    handleSelect: vi.fn(),
    setShowDetailsData: vi.fn(),
  };

  const sampleData = [
    {
      kind: 'edge',
      relationships: [
        { id: 'r1', subType: 'parent', model: { name: 'kubernetes' } },
        { id: 'r2', subType: 'child', model: { name: 'kubernetes' } },
      ],
    },
    {
      kind: 'hierarchical',
      relationships: [{ id: 'r3', subType: 'inherits', model: { name: 'istio' } }],
    },
  ];

  it('renders a top-level node per kind with a count suffix', () => {
    render(<RelationshipTree {...baseProps} data={sampleData} />);

    // RELATIONSHIPS is the default view, so idForKind uses relationships[0].id
    expect(screen.getByTestId('tree-item-r1')).toBeInTheDocument();
    expect(screen.getByTestId('tree-item-r3')).toBeInTheDocument();
    expect(screen.getByTestId('tree-item-button-r1')).toHaveTextContent('edge (2)');
    expect(screen.getByTestId('tree-item-button-r3')).toHaveTextContent('hierarchical (1)');
  });

  it('renders child tree items labelled with subType (model.name)', () => {
    render(<RelationshipTree {...baseProps} data={sampleData} />);

    expect(screen.getByTestId('tree-item-r1.r2')).toBeInTheDocument();
    expect(screen.getByTestId('tree-item-button-r1.r2')).toHaveTextContent('child (kubernetes)');
  });

  it('uses idForKindAsProp prefix when view is not RELATIONSHIPS', () => {
    render(
      <RelationshipTree
        {...baseProps}
        data={sampleData.slice(0, 1)}
        view="Models"
        idForKindAsProp="parent-id"
      />,
    );

    expect(screen.getByTestId('tree-item-parent-id.r1')).toBeInTheDocument();
  });

  it('fires setShowDetailsData with the relationship payload on a leaf click', async () => {
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();

    render(
      <RelationshipTree
        {...baseProps}
        data={sampleData.slice(0, 1)}
        setShowDetailsData={setShowDetailsData}
      />,
    );

    await user.click(screen.getByTestId('tree-item-button-r1.r2'));

    expect(setShowDetailsData).toHaveBeenCalledTimes(1);
    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'Relationships',
      data: sampleData[0].relationships[1],
    });
  });

  it('fires setShowDetailsData with type: "none" on a kind-level click', async () => {
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();

    render(
      <RelationshipTree
        {...baseProps}
        data={sampleData.slice(0, 1)}
        setShowDetailsData={setShowDetailsData}
      />,
    );

    await user.click(screen.getByTestId('tree-item-button-r1'));

    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'none',
      data: { id: 'r1' },
    });
  });

  it('renders the loading spinner with color="inherit" when fetching', () => {
    render(<RelationshipTree {...baseProps} data={[]} isRelationshipFetching={true} />);

    expect(screen.getByTestId('loading-spinner')).toHaveAttribute('data-color', 'inherit');
  });

  it('does not render a loading spinner by default', () => {
    render(<RelationshipTree {...baseProps} data={[]} />);
    expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
  });
});
