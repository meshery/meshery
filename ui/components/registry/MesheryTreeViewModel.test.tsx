import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import MesheryTreeViewModel from './MesheryTreeViewModel';

vi.mock('@sistent/sistent', () => ({
  CircularProgress: () => <div data-testid="loading-spinner" />,
  Box: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('../shared/TreeView', () => ({
  SimpleTreeView: ({ children, ...props }: any) => (
    <div data-testid="simple-tree-view" data-multi-select={String(props.multiSelect)}>
      {children}
    </div>
  ),
}));

vi.mock('../../assets/icons/MinusSquare', () => ({ default: () => <svg /> }));
vi.mock('../../assets/icons/PlusSquare', () => ({ default: () => <svg /> }));
vi.mock('../../assets/icons/DotSquare', () => ({ default: () => <svg /> }));

vi.mock('./MesheryTreeViewItem', () => ({
  default: ({ modelDef }: any) => (
    <div data-testid={`tree-item-${modelDef.id}`}>{modelDef.name}</div>
  ),
}));

describe('MesheryTreeViewModel', () => {
  const baseProps = {
    handleToggle: vi.fn(),
    handleSelect: vi.fn(),
    expanded: [],
    selected: [],
    setShowDetailsData: vi.fn(),
    lastModelRef: { current: null } as React.MutableRefObject<any>,
    isModelFetching: false,
    showDetailsData: { type: '', data: {} },
  };

  it('renders one MesheryTreeViewItem per model', () => {
    const data = [
      { id: 'm1', name: 'kubernetes' },
      { id: 'm2', name: 'istio' },
    ];

    render(<MesheryTreeViewModel {...baseProps} data={data} />);

    expect(screen.getByTestId('tree-item-m1')).toHaveTextContent('kubernetes');
    expect(screen.getByTestId('tree-item-m2')).toHaveTextContent('istio');
  });

  it('shows the spinner when isModelFetching is true', () => {
    render(<MesheryTreeViewModel {...baseProps} data={[]} isModelFetching={true} />);
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('hides the spinner when not fetching', () => {
    render(<MesheryTreeViewModel {...baseProps} data={[]} />);
    expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
  });

  it('renders an empty tree when data is empty', () => {
    render(<MesheryTreeViewModel {...baseProps} data={[]} />);
    expect(screen.getByTestId('simple-tree-view')).toBeInTheDocument();
  });
});
