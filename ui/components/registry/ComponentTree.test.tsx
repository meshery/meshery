import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import ComponentTree from './ComponentTree';

vi.mock('@sistent/sistent', () => ({
  CircularProgress: () => <div data-testid="loading-spinner" />,
}));

vi.mock('../shared/TreeView', () => ({
  SimpleTreeView: ({ children, ...props }: any) => (
    <div data-testid="simple-tree-view" data-multi-select={String(props.multiSelect)}>
      {children}
    </div>
  ),
}));

vi.mock('../../constants/navigator', () => ({
  COMPONENTS: 'Components',
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
  default: ({ itemId, labelText, onClick }: any) => (
    <button data-testid={`tree-item-${itemId}`} onClick={onClick}>
      {labelText}
    </button>
  ),
}));

describe('ComponentTree', () => {
  const baseProps = {
    expanded: [],
    selected: [],
    handleToggle: vi.fn(),
    handleSelect: vi.fn(),
    setShowDetailsData: vi.fn(),
    lastComponentRef: { current: null } as React.MutableRefObject<any>,
    isComponentFetching: false,
  };

  it('renders one tree item per component in data', () => {
    const data = [
      { id: 'c1', displayName: 'Component One' },
      { id: 'c2', displayName: 'Component Two' },
    ];

    render(<ComponentTree {...baseProps} data={data} />);

    expect(screen.getByTestId('tree-item-c1')).toHaveTextContent('Component One');
    expect(screen.getByTestId('tree-item-c2')).toHaveTextContent('Component Two');
  });

  it('passes setShowDetailsData with the component data when a tree item is clicked', async () => {
    const setShowDetailsData = vi.fn();
    const user = userEvent.setup();
    const data = [{ id: 'c1', displayName: 'One' }];

    render(<ComponentTree {...baseProps} data={data} setShowDetailsData={setShowDetailsData} />);

    await user.click(screen.getByTestId('tree-item-c1'));

    expect(setShowDetailsData).toHaveBeenCalledWith({
      type: 'Components',
      data: data[0],
    });
  });

  it('renders a loading spinner when isComponentFetching is true', () => {
    render(<ComponentTree {...baseProps} data={[]} isComponentFetching={true} />);
    expect(screen.getByTestId('loading-spinner')).toBeInTheDocument();
  });

  it('does not render a loading spinner when not fetching', () => {
    render(<ComponentTree {...baseProps} data={[]} />);
    expect(screen.queryByTestId('loading-spinner')).not.toBeInTheDocument();
  });
});
