import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  Dialog: ({ children, open, fullScreen, fullWidth }: any) =>
    open ? (
      <div
        data-testid="dialog"
        data-fullscreen={String(!!fullScreen)}
        data-fullwidth={String(!!fullWidth)}
      >
        {children}
      </div>
    ) : null,
  DialogActions: ({ children }: any) => <div data-testid="dialog-actions">{children}</div>,
  DialogContent: ({ children }: any) => <div data-testid="dialog-content">{children}</div>,
  Divider: () => <hr data-testid="divider" />,
  FullScreenExitIcon: () => <svg data-testid="exit-fullscreen-icon" />,
  FullScreenIcon: () => <svg data-testid="enter-fullscreen-icon" />,
  IconButton: ({ children, onClick, 'aria-label': ariaLabel }: any) => (
    <button data-testid={ariaLabel ? `icon-${ariaLabel}` : 'icon-button'} onClick={onClick}>
      {children}
    </button>
  ),
  SaveIcon: () => <svg data-testid="save-icon" />,
  Tooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('./CodeMirror', () => ({
  UnControlled: ({ value, onChange, options }: any) => (
    <textarea
      data-testid="codemirror"
      data-readonly={String(!!options?.readOnly)}
      value={value || ''}
      onChange={(e) => onChange?.(null, null, e.target.value)}
    />
  ),
}));

vi.mock('./YamlDialog.styles', () => ({
  YamlDialogTitleText: ({ children, variant }: any) => (
    <h6 data-testid="title-text" data-variant={variant}>
      {children}
    </h6>
  ),
  StyledDialog: ({ children, id }: any) => (
    <div data-testid="styled-dialog" id={id}>
      {children}
    </div>
  ),
}));

vi.mock('../designs/patterns/Cards.styles', () => ({
  StyledCodeMirrorWrapper: ({ children, fullScreen }: any) => (
    <div data-testid="codemirror-wrapper" data-fullscreen={String(!!fullScreen)}>
      {children}
    </div>
  ),
}));

import YAMLDialog from './YamlDialog';

describe('YamlDialog', () => {
  const baseProps = {
    fullScreen: false,
    name: 'my-pattern.yaml',
    toggleFullScreen: vi.fn(),
    config_file: 'kind: Deployment\n',
    setYaml: vi.fn(),
    deleteHandler: vi.fn(),
    updateHandler: vi.fn(),
  };

  it('renders the dialog with the provided name', () => {
    render(<YAMLDialog {...baseProps} />);
    expect(screen.getByTestId('dialog')).toBeInTheDocument();
    expect(screen.getByTestId('title-text')).toHaveTextContent('my-pattern.yaml');
  });

  it('renders the enter-fullscreen icon when not fullscreen', () => {
    render(<YAMLDialog {...baseProps} fullScreen={false} />);
    expect(screen.getByTestId('enter-fullscreen-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('exit-fullscreen-icon')).not.toBeInTheDocument();
  });

  it('renders the exit-fullscreen icon when fullscreen is true', () => {
    render(<YAMLDialog {...baseProps} fullScreen={true} />);
    expect(screen.getByTestId('exit-fullscreen-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('enter-fullscreen-icon')).not.toBeInTheDocument();
    expect(screen.getByTestId('dialog')).toHaveAttribute('data-fullscreen', 'true');
  });

  it('invokes toggleFullScreen when the fullscreen button is clicked', async () => {
    const toggleFullScreen = vi.fn();
    const user = userEvent.setup();
    render(<YAMLDialog {...baseProps} toggleFullScreen={toggleFullScreen} />);

    // The fullscreen IconButton has no aria-label, so find via the icon test id.
    const button = screen.getByTestId('enter-fullscreen-icon').closest('button');
    await user.click(button!);
    expect(toggleFullScreen).toHaveBeenCalledTimes(1);
  });

  it('renders save / delete actions when not read-only', async () => {
    const updateHandler = vi.fn();
    const deleteHandler = vi.fn();
    const user = userEvent.setup();

    render(
      <YAMLDialog
        {...baseProps}
        isReadOnly={false}
        updateHandler={updateHandler}
        deleteHandler={deleteHandler}
      />,
    );

    await user.click(screen.getByTestId('icon-Update'));
    await user.click(screen.getByTestId('icon-Delete'));

    expect(updateHandler).toHaveBeenCalledTimes(1);
    expect(deleteHandler).toHaveBeenCalledTimes(1);
  });

  it('omits action buttons when isReadOnly is true', () => {
    render(<YAMLDialog {...baseProps} isReadOnly={true} />);

    expect(screen.queryByTestId('icon-Update')).not.toBeInTheDocument();
    expect(screen.queryByTestId('icon-Delete')).not.toBeInTheDocument();
    expect(screen.queryByTestId('dialog-actions')).not.toBeInTheDocument();
    expect(screen.getByTestId('codemirror')).toHaveAttribute('data-readonly', 'true');
  });

  it('forwards CodeMirror changes through setYaml', async () => {
    const setYaml = vi.fn();
    const user = userEvent.setup();
    render(<YAMLDialog {...baseProps} setYaml={setYaml} config_file="" />);

    const editor = screen.getByTestId('codemirror');
    await user.type(editor, 'a');
    expect(setYaml).toHaveBeenLastCalledWith('a');
  });
});
