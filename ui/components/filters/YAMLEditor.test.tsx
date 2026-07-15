import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const can = vi.fn(() => true);

vi.mock('@/utils/can', () => ({
  default: (...args: unknown[]) => can(...args),
}));

vi.mock('@/assets/icons', () => ({
  Close: () => <svg data-testid="close-icon" />,
  Delete: () => <svg data-testid="delete-icon" />,
  Save: () => <svg data-testid="save-icon" />,
}));

vi.mock('@sistent/sistent', () => {
  const styled = (_Component: any) => (_factory?: any) => {
    const Styled = ({ children }: any) => <div>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    CustomTooltip: ({ children, title }: any) => (
      <div data-testid={`tooltip-${title}`}>{children}</div>
    ),
    Dialog: ({ children, open }: any) => (open ? <div data-testid="dialog">{children}</div> : null),
    DialogActions: ({ children }: any) => <div data-testid="actions">{children}</div>,
    DialogTitle: ({ children }: any) => <div data-testid="title">{children}</div>,
    Divider: () => <hr data-testid="divider" />,
    IconButton: ({ children, onClick, disabled, ...props }: any) => (
      <button onClick={onClick} disabled={disabled} {...props}>
        {children}
      </button>
    ),
    styled,
    FullScreenIcon: () => <svg data-testid="fullscreen-icon" />,
    FullScreenExitIcon: () => <svg data-testid="fullscreen-exit-icon" />,
  };
});

vi.mock('../CodeMirror', () => ({
  UnControlled: ({ value, onChange }: any) => (
    <textarea
      data-testid="codemirror"
      value={value || ''}
      onChange={(e) => onChange?.({}, {}, e.target.value)}
    />
  ),
}));

vi.mock('../../utils/Enum', () => ({
  FILE_OPS: { UPDATE: 'UPDATE', DELETE: 'DELETE' },
}));

vi.mock('../../css/icons.styles', () => ({
  iconMedium: {},
}));

vi.mock('./TooltipIcon', () => ({
  default: ({ children, onClick, title }: any) => (
    <button data-testid={`tip-${title}`} onClick={onClick}>
      {children}
    </button>
  ),
}));

vi.mock('./Filters.styled', () => ({
  YmlDialogTitle: ({ children }: any) => <div>{children}</div>,
  YmlDialogTitleText: ({ children }: any) => <h2>{children}</h2>,
}));

import YAMLEditor from './YAMLEditor';

describe('YAMLEditor', () => {
  beforeEach(() => {
    can.mockReset();
    can.mockReturnValue(true);
  });

  const filter = {
    id: 'f-1',
    name: 'my filter',
    catalogData: { foo: 'bar' },
    filter_resource: JSON.stringify({ settings: { config: 'apiVersion: v1' } }),
  };

  it('renders the filter name as the dialog title', () => {
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={() => {}} />);
    expect(screen.getByText('my filter')).toBeInTheDocument();
  });

  it('initialises CodeMirror with the YAML from filter_resource.settings.config', () => {
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={() => {}} />);
    expect(screen.getByTestId('codemirror')).toHaveValue('apiVersion: v1');
  });

  it('falls back to empty config when filter_resource is invalid JSON', () => {
    const badFilter = { ...filter, filter_resource: '{invalid json' };
    render(<YAMLEditor filter={badFilter} onClose={() => {}} onSubmit={() => {}} />);
    expect(screen.getByTestId('codemirror')).toHaveValue('');
  });

  it('calls onSubmit with type=UPDATE when the save button is clicked', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={onSubmit} />);

    await user.click(screen.getByRole('button', { name: 'Update' }));
    expect(onSubmit).toHaveBeenCalledWith({
      data: 'apiVersion: v1',
      id: 'f-1',
      name: 'my filter',
      type: 'UPDATE',
      catalogData: { foo: 'bar' },
    });
  });

  it('calls onSubmit with type=DELETE when the delete button is clicked', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn();
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={onSubmit} />);

    await user.click(screen.getByRole('button', { name: 'Delete' }));
    expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({ type: 'DELETE', id: 'f-1' }));
  });

  it('toggles fullscreen on the fullscreen icon click', async () => {
    const user = userEvent.setup();
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={() => {}} />);

    await user.click(screen.getByTestId('tip-Enter Fullscreen'));
    // Once toggled, the title text should switch to "Exit Fullscreen"
    expect(screen.getByTestId('tip-Exit Fullscreen')).toBeInTheDocument();
  });

  it('calls onClose when the exit button is clicked', async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    render(<YAMLEditor filter={filter} onClose={onClose} onSubmit={() => {}} />);
    await user.click(screen.getByTestId('tip-Exit'));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('disables save/delete buttons when permission CAN returns false', () => {
    can.mockReturnValue(false);
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={() => {}} />);
    expect(screen.getByRole('button', { name: 'Update' })).toBeDisabled();
    expect(screen.getByRole('button', { name: 'Delete' })).toBeDisabled();
  });

  it('updates the yaml state when CodeMirror onChange fires', async () => {
    const onSubmit = vi.fn();
    render(<YAMLEditor filter={filter} onClose={() => {}} onSubmit={onSubmit} />);

    const editor = screen.getByTestId('codemirror') as HTMLTextAreaElement;
    // Drive a single change event - the mocked editor doesn't preserve value
    // across multiple character events, but a fireEvent.change captures the
    // final value verbatim.
    const { fireEvent } = await import('@testing-library/react');
    fireEvent.change(editor, { target: { value: 'foo: 1' } });

    await import('@testing-library/user-event').then(({ default: userEvent2 }) =>
      userEvent2.setup().click(screen.getByRole('button', { name: 'Update' })),
    );
    expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({ data: 'foo: 1' }));
  });
});
