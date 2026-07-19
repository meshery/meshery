import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const CAN_mock = vi.fn(() => true);

vi.mock('@sistent/sistent', () => ({
  AddIcon: () => <svg data-testid="add-icon" />,
  Box: ({ children, ...props }: any) => (
    <div data-testid="box" {...props}>
      {children}
    </div>
  ),
  CardActions: ({ children, disableSpacing }: any) => (
    <div data-testid="card-actions" data-disable-spacing={String(!!disableSpacing)}>
      {children}
    </div>
  ),
  CardHeader: ({ title, subheader }: any) => (
    <div data-testid="card-header">
      <span data-testid="card-title">{title}</span>
      <span data-testid="card-subheader">{subheader}</span>
    </div>
  ),
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  IconButton: ({ children, onClick, disabled, 'aria-label': ariaLabel, ref }: any) => (
    <button
      data-testid={ariaLabel ? `icon-button-${ariaLabel}` : 'icon-button'}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  ),
  Menu: ({ children, open, onClose }: any) =>
    open ? (
      <div data-testid="menu">
        <button type="button" onClick={onClose}>
          close-menu
        </button>
        {children}
      </div>
    ) : null,
  MenuItem: ({ children, onClick }: any) => (
    <button data-testid="menu-item" type="button" onClick={onClick}>
      {children}
    </button>
  ),
  PlayArrowIcon: () => <svg data-testid="play-icon" />,
}));

vi.mock('../css/icons.styles', () => ({
  iconMedium: { width: 24 },
}));

vi.mock('@/utils/can', () => ({
  default: (...args: any[]) => CAN_mock(...args),
}));

vi.mock('./adapter-play-styled', () => ({
  AdapterCard: ({ children }: any) => <div data-testid="adapter-card">{children}</div>,
}));

import AdapterCategoryCard from './adapter-play-category-card';

const buildProps = (overrides: Partial<any> = {}) => {
  const addIconEles = { current: {} as Record<number, HTMLElement | null> };
  const delIconEles = { current: {} as Record<number, HTMLElement | null> };
  return {
    cat: 0,
    adapterOps: [
      { key: 'op1', value: 'Bravo' },
      { key: 'op2', value: 'Alpha' },
      { key: 'op3', value: 'Charlie', category: 1 },
    ],
    menuState: {
      0: { add: false, delete: false },
      1: { add: false, delete: false },
      2: { add: false, delete: false },
      3: { add: false, delete: false },
      4: { add: false, delete: false },
    },
    addIconEles,
    delIconEles,
    onMenuToggle: vi.fn(),
    onMenuItemClick: vi.fn(() => vi.fn()),
    renderYamlDialog: vi.fn((cat: number, isDelete: boolean) => (
      <div data-testid={`yaml-dialog-${cat}-${String(isDelete)}`} />
    )),
    ...overrides,
  };
};

describe('AdapterCategoryCard', () => {
  beforeEach(() => {
    CAN_mock.mockClear();
    CAN_mock.mockReturnValue(true);
  });

  it('renders category 0 (Lifecycle) with install + delete buttons and the AddIcon', () => {
    render(<AdapterCategoryCard {...buildProps({ cat: 0 })} />);

    expect(screen.getByTestId('card-title')).toHaveTextContent(
      'Manage Cloud Native Infrastructure Lifecycle',
    );
    expect(screen.getByTestId('icon-button-install')).toBeInTheDocument();
    expect(screen.getByTestId('icon-button-delete')).toBeInTheDocument();
    expect(screen.getByTestId('add-icon')).toBeInTheDocument();
  });

  it('renders the PlayIcon for category 4 (custom configuration)', () => {
    render(<AdapterCategoryCard {...buildProps({ cat: 4 })} />);
    expect(screen.getByTestId('play-icon')).toBeInTheDocument();
  });

  it('omits the delete button for category 3 (validation)', () => {
    render(<AdapterCategoryCard {...buildProps({ cat: 3 })} />);
    expect(screen.queryByTestId('icon-button-delete')).not.toBeInTheDocument();
  });

  it('disables buttons when CAN returns false', () => {
    CAN_mock.mockReturnValue(false);
    render(<AdapterCategoryCard {...buildProps({ cat: 0 })} />);
    expect(screen.getByTestId('icon-button-install')).toBeDisabled();
    expect(screen.getByTestId('icon-button-delete')).toBeDisabled();
  });

  it('invokes onMenuToggle when install/delete are clicked', async () => {
    const user = userEvent.setup();
    const onMenuToggle = vi.fn();
    render(<AdapterCategoryCard {...buildProps({ cat: 0, onMenuToggle })} />);

    await user.click(screen.getByTestId('icon-button-install'));
    expect(onMenuToggle).toHaveBeenCalledWith(0, false);

    await user.click(screen.getByTestId('icon-button-delete'));
    expect(onMenuToggle).toHaveBeenCalledWith(0, true);
  });

  it('renders menu items for matching adapter ops when menu is open', () => {
    const props = buildProps({
      cat: 0,
      menuState: {
        0: { add: true, delete: false },
        1: { add: false, delete: false },
        2: { add: false, delete: false },
        3: { add: false, delete: false },
        4: { add: false, delete: false },
      },
    });
    render(<AdapterCategoryCard {...props} />);

    const items = screen.getAllByTestId('menu-item');
    // Adapter ops with no category default to cat 0 (Bravo & Alpha).
    expect(items).toHaveLength(2);
    // Items should be sorted alphabetically: Alpha, then Bravo.
    expect(items[0]).toHaveTextContent('Alpha');
    expect(items[1]).toHaveTextContent('Bravo');
  });

  it('renders YAML dialog for category 4 instead of a menu', () => {
    render(<AdapterCategoryCard {...buildProps({ cat: 4 })} />);
    expect(screen.getByTestId('yaml-dialog-4-false')).toBeInTheDocument();
    expect(screen.getByTestId('yaml-dialog-4-true')).toBeInTheDocument();
  });

  it('filters out Add-on: prefixed ops for category 2', () => {
    const adapterOps = [
      { key: 'op1', value: 'Visible', category: 2 },
      { key: 'op2', value: 'Add-on:Hidden', category: 2 },
    ];
    const menuState = {
      0: { add: false, delete: false },
      1: { add: false, delete: false },
      2: { add: true, delete: false },
      3: { add: false, delete: false },
      4: { add: false, delete: false },
    };
    render(<AdapterCategoryCard {...buildProps({ cat: 2, adapterOps, menuState })} />);

    const items = screen.getAllByTestId('menu-item');
    expect(items).toHaveLength(1);
    expect(items[0]).toHaveTextContent('Visible');
  });

  it('handles missing adapterOps gracefully', () => {
    expect(() =>
      render(<AdapterCategoryCard {...buildProps({ cat: 0, adapterOps: undefined as any })} />),
    ).not.toThrow();
  });
});
