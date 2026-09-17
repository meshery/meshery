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
  MenuItem: ({ children, onClick, disabled, ...rest }: any) => (
    <li>
      <button type="button" onClick={onClick} disabled={disabled} {...rest}>
        {children}
      </button>
    </li>
  ),
  MenuList: ({ children }: any) => <ul>{children}</ul>,
  Paper: ({ children }: any) => <div>{children}</div>,
  Popper: ({ children, open }: any) => (open ? <div data-testid="popper">{children}</div> : null),
  CustomTooltip: ({ children, title }: any) => (
    <span data-testid="custom-tooltip" data-title={title}>
      {children}
    </span>
  ),
  Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
}));

import ActionButton from './ActionButton';

describe('ActionButton', () => {
  it('renders the Action button and dropdown toggle with tooltips', () => {
    render(<ActionButton options={[]} />);
    expect(screen.getByRole('button', { name: 'Action' })).toBeInTheDocument();
    expect(screen.getByTestId('action-btn-toggle')).toBeInTheDocument();

    const tooltips = screen.getAllByTestId('custom-tooltip');
    const titles = tooltips.map((t) => t.getAttribute('data-title'));
    expect(titles).toContain('Invoke actions interactively');
    expect(titles).toContain('Invoke actions in single click');
  });

  it('invokes defaultActionClick when provided', () => {
    const handler = vi.fn();
    render(<ActionButton defaultActionClick={handler} options={[]} />);
    fireEvent.click(screen.getByRole('button', { name: 'Action' }));
    expect(handler).toHaveBeenCalled();
  });

  it('toggles interactive mode when main Action button is clicked without defaultActionClick', () => {
    render(
      <ActionButton
        options={[
          { label: 'Validate', icon: <span data-testid="validate-icon" />, onClick: vi.fn() },
        ]}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Action' }));
    expect(screen.getByTestId('popper')).toBeInTheDocument();
    expect(screen.getByText('Validate')).toBeInTheDocument();
    expect(screen.getByTestId('validate-icon')).toBeInTheDocument();
  });

  it('toggles single-click mode when dropdown toggle is clicked', () => {
    render(
      <ActionButton
        options={[
          { label: 'Validate', icon: <span data-testid="validate-icon" />, onClick: vi.fn() },
        ]}
      />,
    );

    fireEvent.click(screen.getByTestId('action-btn-toggle'));
    expect(screen.getByTestId('popper')).toBeInTheDocument();
    expect(screen.getByTestId('validate-icon')).toBeInTheDocument();
    // In single-click mode, the label is inside the tooltip only, not rendered as direct text
    expect(screen.queryByText('Validate')).not.toBeInTheDocument();
  });

  it('invokes onClick in interactive mode', () => {
    const onClick = vi.fn();
    const onDirectClick = vi.fn();
    render(<ActionButton options={[{ label: 'Deploy', icon: <svg />, onClick, onDirectClick }]} />);

    fireEvent.click(screen.getByRole('button', { name: 'Action' }));
    fireEvent.click(screen.getByRole('button', { name: 'Deploy' }));

    expect(onClick).toHaveBeenCalled();
    expect(onDirectClick).not.toHaveBeenCalled();
    expect(screen.queryByTestId('popper')).not.toBeInTheDocument();
  });

  it('invokes onDirectClick in single-click mode when available', () => {
    const onClick = vi.fn();
    const onDirectClick = vi.fn();
    render(<ActionButton options={[{ label: 'Deploy', icon: <svg />, onClick, onDirectClick }]} />);

    fireEvent.click(screen.getByTestId('action-btn-toggle'));
    fireEvent.click(screen.getByTestId('action-btn-option-Deploy'));

    expect(onDirectClick).toHaveBeenCalled();
    expect(onClick).not.toHaveBeenCalled();
    expect(screen.queryByTestId('popper')).not.toBeInTheDocument();
  });

  it('falls back to onClick in single-click mode when onDirectClick is not provided', () => {
    const onClick = vi.fn();
    render(<ActionButton options={[{ label: 'Deploy', icon: <svg />, onClick }]} />);

    fireEvent.click(screen.getByTestId('action-btn-toggle'));
    fireEvent.click(screen.getByTestId('action-btn-option-Deploy'));

    expect(onClick).toHaveBeenCalled();
  });

  it('respects disabled options', () => {
    render(
      <ActionButton
        options={[{ label: 'Disabled', icon: <svg />, onClick: vi.fn(), disabled: true }]}
      />,
    );

    fireEvent.click(screen.getByRole('button', { name: 'Action' }));
    expect(screen.getByRole('button', { name: 'Disabled' })).toBeDisabled();
  });

  it('forwards permissionKey and permissionAction to the primary Action button', () => {
    const keyObj = { id: 'test_key', function: 'read' };
    render(<ActionButton options={[]} permissionKey={keyObj} permissionAction="hide" />);

    const mainBtn = screen.getByRole('button', { name: 'Action' });
    expect(mainBtn).toHaveAttribute('data-permission-key', JSON.stringify(keyObj));
    expect(mainBtn).toHaveAttribute('data-permission-action', 'hide');
  });
});
