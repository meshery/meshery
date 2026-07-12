import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (_tag: any) => () => {
    const Styled = ({ children, ...rest }: any) => <div {...rest}>{children}</div>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    styled,
    CustomTooltip: ({ children, title, onClick }: any) => (
      <div data-testid={`tooltip-${title}`} onClick={onClick}>
        {children}
      </div>
    ),
    Dialog: ({ children, open, onClose }: any) =>
      open ? (
        <div data-testid="dialog">
          <button type="button" onClick={onClose} aria-label="dialog-close">
            close
          </button>
          {children}
        </div>
      ) : null,
    DialogActions: ({ children }: any) => <div>{children}</div>,
    DialogContent: ({ children }: any) => <div>{children}</div>,
    Divider: () => <hr />,
    IconButton: ({ children, onClick, disabled, 'aria-label': ariaLabel }: any) => (
      <button type="button" onClick={onClick} disabled={disabled} aria-label={ariaLabel}>
        {children}
      </button>
    ),
    FullScreenIcon: () => <svg data-testid="fullscreen-icon" />,
    FullScreenExitIcon: () => <svg data-testid="fullscreen-exit-icon" />,
  };
});

vi.mock('@/assets/icons', () => ({
  Close: () => <svg data-testid="close-icon" />,
  Delete: () => <svg data-testid="delete-icon" />,
  Save: () => <svg data-testid="save-icon" />,
}));

vi.mock('../../CodeMirror', () => ({
  UnControlled: ({ value, onChange }: any) => (
    <textarea
      data-testid="codemirror"
      defaultValue={value}
      onChange={(e) => onChange?.('editor', { data: 1 }, e.target.value)}
    />
  ),
}));

vi.mock('../../../utils/Enum', () => ({
  FILE_OPS: { UPDATE: 'UPDATE', DELETE: 'DELETE' },
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

vi.mock('./MesheryPatterns.styled', () => ({
  YamlDialogTitle: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
  YamlDialogTitleText: ({ children }: any) => <span>{children}</span>,
}));

import YAMLEditor from './YAMLEditor';

describe('YAMLEditor', () => {
  const pattern = {
    id: 'p1',
    name: 'pattern-1',
    patternFile: 'foo: bar',
    catalogData: { tag: 'demo' },
  };

  it('renders the pattern name in the dialog title', () => {
    render(<YAMLEditor pattern={pattern} onClose={vi.fn()} onSubmit={vi.fn()} />);
    expect(screen.getByText('pattern-1')).toBeInTheDocument();
    expect(screen.getByTestId('dialog')).toBeInTheDocument();
  });

  it('invokes onSubmit with UPDATE payload when the save button is clicked', () => {
    const onSubmit = vi.fn();
    render(<YAMLEditor pattern={pattern} onClose={vi.fn()} onSubmit={onSubmit} />);

    fireEvent.click(screen.getByRole('button', { name: 'Update' }));
    expect(onSubmit).toHaveBeenCalledWith({
      data: 'foo: bar',
      id: 'p1',
      name: 'pattern-1',
      type: 'UPDATE',
      catalogData: { tag: 'demo' },
    });
  });

  it('invokes onSubmit with DELETE payload when the delete button is clicked', () => {
    const onSubmit = vi.fn();
    render(<YAMLEditor pattern={pattern} onClose={vi.fn()} onSubmit={onSubmit} />);

    fireEvent.click(screen.getByRole('button', { name: 'Delete' }));
    expect(onSubmit).toHaveBeenCalledWith({
      data: 'foo: bar',
      id: 'p1',
      name: 'pattern-1',
      type: 'DELETE',
      catalogData: { tag: 'demo' },
    });
  });

  it('hides Save and Delete buttons in read-only mode', () => {
    render(<YAMLEditor pattern={pattern} onClose={vi.fn()} onSubmit={vi.fn()} isReadOnly />);
    expect(screen.queryByRole('button', { name: 'Update' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Delete' })).not.toBeInTheDocument();
  });
});
