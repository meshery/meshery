import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  ArrowDropDownIcon: () => <svg data-testid="dropdown-icon" />,
  Button: ({ children, onClick, permissionKey, permissionAction, ...rest }: any) => (
    <button
      type="button"
      onClick={onClick}
      data-permission-key={
        typeof permissionKey === 'object' ? JSON.stringify(permissionKey) : permissionKey
      }
      data-permission-action={permissionAction}
      {...rest}
    >
      {children}
    </button>
  ),
  ButtonGroup: ({ children }: any) => <div>{children}</div>,
  ClickAwayListener: ({ children, onClickAway }: any) => (
    <div>
      {children}
      <button
        data-testid="click-away"
        type="button"
        onClick={() => onClickAway?.({ target: null })}
      >
        away
      </button>
    </div>
  ),
  MenuItem: ({ children, onClick, disabled }: any) => (
    <li>
      <button type="button" onClick={onClick} disabled={disabled}>
        {children}
      </button>
    </li>
  ),
  MenuList: ({ children }: any) => <ul>{children}</ul>,
  Paper: ({ children }: any) => <div>{children}</div>,
  Popper: ({ children, open }: any) => (open ? <div data-testid="popper">{children}</div> : null),
}));

import ActionButton from './ActionButton';

describe('ActionButton', () => {
  it('renders the default Action button', () => {
    render(<ActionButton defaultActionClick={vi.fn()} options={[]} />);
    expect(screen.getByRole('button', { name: 'Action' })).toBeInTheDocument();
  });

  it('invokes the default click handler', () => {
    const handler = vi.fn();
    render(<ActionButton defaultActionClick={handler} options={[]} />);
    fireEvent.click(screen.getByRole('button', { name: 'Action' }));
    expect(handler).toHaveBeenCalled();
  });

  it('toggles the dropdown popper open and lists supplied options', () => {
    const onClick = vi.fn();
    render(
      <ActionButton
        defaultActionClick={vi.fn()}
        options={[{ label: 'Validate', icon: <svg />, onClick }]}
      />,
    );

    fireEvent.click(screen.getByTestId('action-btn-toggle'));
    expect(screen.getByTestId('popper')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Validate' })).toBeInTheDocument();
  });

  it('invokes the chosen option and closes the popper', () => {
    const onClick = vi.fn();
    render(
      <ActionButton
        defaultActionClick={vi.fn()}
        options={[{ label: 'Deploy', icon: <svg />, onClick }]}
      />,
    );

    fireEvent.click(screen.getByTestId('action-btn-toggle'));
    fireEvent.click(screen.getByRole('button', { name: 'Deploy' }));

    expect(onClick).toHaveBeenCalled();
    expect(screen.queryByTestId('popper')).not.toBeInTheDocument();
  });

  it('respects disabled options', () => {
    render(
      <ActionButton
        defaultActionClick={vi.fn()}
        options={[{ label: 'Disabled', icon: <svg />, onClick: vi.fn(), disabled: true }]}
      />,
    );

    fireEvent.click(screen.getByTestId('action-btn-toggle'));
    expect(screen.getByRole('button', { name: 'Disabled' })).toBeDisabled();
  });

  it('forwards permissionKey and permissionAction to the primary Action button', () => {
    const keyObj = { id: 'test_key', function: 'read' };
    render(
      <ActionButton
        defaultActionClick={vi.fn()}
        options={[]}
        permissionKey={keyObj}
        permissionAction="hide"
      />,
    );

    const mainBtn = screen.getByRole('button', { name: 'Action' });
    expect(mainBtn).toHaveAttribute('data-permission-key', JSON.stringify(keyObj));
    expect(mainBtn).toHaveAttribute('data-permission-action', 'hide');
  });
});
