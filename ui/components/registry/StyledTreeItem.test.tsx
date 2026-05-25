import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Typography: ({ children, ...props }: any) => <span {...props}>{children}</span>,
  Checkbox: ({ onClick, checked, ...props }: any) => (
    <input
      type="checkbox"
      checked={!!checked}
      onClick={onClick}
      data-testid="checkbox"
      onChange={() => {}}
      {...props}
    />
  ),
  useTheme: () => ({
    palette: {
      text: { default: 'black' },
    },
  }),
}));

vi.mock('@/utils/custom-search', () => ({
  default: ({ onSearch, expanded, setExpanded, placeholder }: any) => (
    <div data-testid="search-bar">
      <button data-testid="expand-search" onClick={() => setExpanded(!expanded)} type="button">
        toggle-{String(expanded)}
      </button>
      <input
        data-testid="search-input"
        placeholder={placeholder}
        onChange={(e) => onSearch(e.target.value)}
      />
    </div>
  ),
}));

vi.mock('@/utils/debounce', () => ({
  default: (fn: any) => fn,
}));

vi.mock('./MeshModel.style', () => ({
  StyledTreeItemRoot: ({ children, label, onMouseEnter, onMouseLeave, ...rest }: any) => (
    <div
      data-testid="styled-tree-item-root"
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
      data-item-id={rest.itemId}
      data-id={rest['data-id']}
    >
      <div data-testid="styled-tree-item-label">{label}</div>
      {children}
    </div>
  ),
}));

vi.mock('@/utils/dimension', () => ({
  useWindowDimensions: () => ({ width: 1920, height: 1080 }),
}));

import StyledTreeItem from './StyledTreeItem';

describe('StyledTreeItem', () => {
  it('renders the label text', () => {
    render(<StyledTreeItem labelText="Hello World" />);
    expect(screen.getByText('Hello World')).toBeInTheDocument();
  });

  it('renders the search bar when search is enabled', () => {
    render(<StyledTreeItem labelText="x" search setSearchText={vi.fn()} />);
    expect(screen.getByTestId('search-bar')).toBeInTheDocument();
  });

  it('does not render a search bar by default', () => {
    render(<StyledTreeItem labelText="x" />);
    expect(screen.queryByTestId('search-bar')).not.toBeInTheDocument();
  });

  it('renders a checkbox when check is enabled', () => {
    render(<StyledTreeItem labelText="x" check />);
    expect(screen.getByTestId('checkbox')).toBeInTheDocument();
  });

  it('does not render checkbox by default', () => {
    render(<StyledTreeItem labelText="x" />);
    expect(screen.queryByTestId('checkbox')).not.toBeInTheDocument();
  });

  it('toggles hover state on mouse events without crashing', () => {
    render(<StyledTreeItem labelText="hover" />);
    const root = screen.getByTestId('styled-tree-item-root');
    fireEvent.mouseEnter(root);
    fireEvent.mouseLeave(root);
    expect(screen.getByText('hover')).toBeInTheDocument();
  });

  it('toggles the checked state when the checkbox is clicked', async () => {
    const user = userEvent.setup();
    render(<StyledTreeItem labelText="x" check />);

    const cb = screen.getByTestId('checkbox') as HTMLInputElement;
    expect(cb.checked).toBe(false);
    await user.click(cb);
    expect(cb.checked).toBe(true);
  });

  it('forwards setSearchText to onSearch via the SearchBar', () => {
    const setSearchText = vi.fn();
    render(<StyledTreeItem labelText="x" search setSearchText={setSearchText} />);

    fireEvent.change(screen.getByTestId('search-input'), {
      target: { value: 'k8s' },
    });

    expect(setSearchText).toHaveBeenCalledWith('k8s');
  });

  it('passes extra props through (itemId, data-id)', () => {
    render(<StyledTreeItem labelText="hi" itemId="x.1" data-id="x.1" />);
    const root = screen.getByTestId('styled-tree-item-root');
    expect(root).toHaveAttribute('data-item-id', 'x.1');
    expect(root).toHaveAttribute('data-id', 'x.1');
  });
});
