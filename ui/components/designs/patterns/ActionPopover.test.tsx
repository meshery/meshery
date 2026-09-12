import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  ClickAwayListener: ({ children, onClickAway }: any) => (
    <div>
      {children}
      <button type="button" data-testid="click-away" onClick={() => onClickAway?.({})}>
        away
      </button>
    </div>
  ),
  IconButton: ({ children, onClick, ...rest }: any) => (
    <button type="button" onClick={onClick} {...rest}>
      {children}
    </button>
  ),
  MenuItem: ({ children, onClick, disabled }: any) => (
    <li>
      <button type="button" onClick={onClick} disabled={disabled}>
        {children}
      </button>
    </li>
  ),
  MenuList: ({ children }: any) => <ul>{children}</ul>,
  MoreVertIcon: () => <svg data-testid="more-vert" />,
  Paper: ({ children }: any) => <div data-testid="paper">{children}</div>,
  Popper: ({ children, open }: any) => (open ? <div data-testid="popper">{children}</div> : null),
  Tooltip: ({ children }: any) => <>{children}</>,
}));

import ActionPopover from './ActionPopover';

describe('ActionPopover', () => {
  it('renders the trigger button by default and keeps the popper closed', () => {
    render(<ActionPopover actions={[]} />);
    expect(screen.getByRole('button', { name: 'Open actions menu' })).toBeInTheDocument();
    expect(screen.queryByTestId('popper')).not.toBeInTheDocument();
  });

  it('opens the popper on trigger click and renders the supplied actions', () => {
    const onClick = vi.fn();
    render(<ActionPopover actions={[{ label: 'Delete', icon: <svg />, onClick }]} />);

    fireEvent.click(screen.getByRole('button', { name: 'Open actions menu' }));
    expect(screen.getByTestId('popper')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeInTheDocument();
  });

  it('invokes the action callback and closes the popper after selection', () => {
    const onClick = vi.fn();
    render(<ActionPopover actions={[{ label: 'Edit', icon: <svg />, onClick }]} />);

    fireEvent.click(screen.getByRole('button', { name: 'Open actions menu' }));
    fireEvent.click(screen.getByRole('button', { name: 'Edit' }));

    expect(onClick).toHaveBeenCalledTimes(1);
    expect(screen.queryByTestId('popper')).not.toBeInTheDocument();
  });

  it('respects disabled actions', () => {
    const onClick = vi.fn();
    render(<ActionPopover actions={[{ label: 'Edit', icon: <svg />, onClick, disabled: true }]} />);

    fireEvent.click(screen.getByRole('button', { name: 'Open actions menu' }));
    const button = screen.getByRole('button', { name: 'Edit' });
    expect(button).toBeDisabled();
  });

  it('closes the popper on click-away', () => {
    render(<ActionPopover actions={[{ label: 'Edit', icon: <svg />, onClick: vi.fn() }]} />);
    fireEvent.click(screen.getByRole('button', { name: 'Open actions menu' }));
    expect(screen.getByTestId('popper')).toBeInTheDocument();

    fireEvent.click(screen.getByTestId('click-away'));
    expect(screen.queryByTestId('popper')).not.toBeInTheDocument();
  });
});
