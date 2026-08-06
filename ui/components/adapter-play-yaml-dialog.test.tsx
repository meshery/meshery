import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  Dialog: ({ children, open, onClose, fullWidth, maxWidth }: any) =>
    open ? (
      <div data-testid="dialog" data-fullwidth={String(!!fullWidth)} data-maxwidth={maxWidth}>
        <button onClick={onClose} type="button">
          backdrop-close
        </button>
        {children}
      </div>
    ) : null,
  DialogActions: ({ children }: any) => <div data-testid="dialog-actions">{children}</div>,
  DialogContent: ({ children }: any) => <div data-testid="dialog-content">{children}</div>,
  DialogTitle: ({ children, onClose, id }: any) => (
    <div data-testid="dialog-title" id={id}>
      <button onClick={onClose} type="button">
        close-from-title
      </button>
      {children}
    </div>
  ),
  Divider: () => <hr data-testid="divider" />,
  Grid: ({ children, item, container, xs, spacing }: any) => (
    <div
      data-testid="grid"
      data-item={String(!!item)}
      data-container={String(!!container)}
      data-xs={xs}
      data-spacing={spacing}
    >
      {children}
    </div>
  ),
  IconButton: ({ children, onClick, 'aria-label': ariaLabel }: any) => (
    <button data-testid={ariaLabel ? `icon-${ariaLabel}` : 'icon-button'} onClick={onClick}>
      {children}
    </button>
  ),
  PlayArrowIcon: () => <svg data-testid="play-icon" />,
}));

vi.mock('./general/CodeMirror', () => ({
  Controlled: ({ value, onBeforeChange, options }: any) => (
    <textarea
      data-testid="codemirror"
      data-mode={options?.mode}
      value={value || ''}
      onChange={(e) => onBeforeChange?.(null, null, e.target.value)}
    />
  ),
}));

vi.mock('./general/ReactSelectWrapper', () => ({
  default: ({ label, value, error, options, onChange }: any) => (
    <div data-testid={`select-${label}`} data-error={String(!!error)}>
      <span data-testid={`select-${label}-value`}>{value?.label}</span>
      <button
        type="button"
        onClick={() => onChange?.(options[0])}
        data-testid={`select-${label}-choose-first`}
      >
        Choose first
      </button>
    </div>
  ),
}));

vi.mock('../css/icons.styles', () => ({
  iconMedium: { width: 24 },
}));

import AdapterYamlDialog from './adapter-play-yaml-dialog';

const baseProps = {
  open: true,
  isDelete: false,
  adapterName: 'Istio',
  namespace: { value: 'default', label: 'default' },
  namespaceError: false,
  namespaceList: [
    { value: 'ns-a', label: 'ns-a' },
    { value: 'ns-b', label: 'ns-b' },
  ],
  onNamespaceChange: vi.fn(),
  version: { value: '1.18', label: '1.18' },
  versionError: false,
  versionList: [
    { value: '1.18', label: '1.18' },
    { value: '1.19', label: '1.19' },
  ],
  onVersionChange: vi.fn(),
  value: 'kind: Deployment\n',
  onBeforeChange: vi.fn(),
  onClose: vi.fn(),
  onApply: vi.fn(),
};

describe('AdapterYamlDialog', () => {
  it('renders nothing when open is false', () => {
    render(<AdapterYamlDialog {...baseProps} open={false} />);
    expect(screen.queryByTestId('dialog')).not.toBeInTheDocument();
  });

  it('renders the apply title in non-delete mode with the PlayIcon', () => {
    render(<AdapterYamlDialog {...baseProps} isDelete={false} />);
    expect(screen.getByTestId('dialog-title')).toHaveTextContent('Istio Adapter - Custom YAML');
    expect(screen.getByTestId('dialog-title')).not.toHaveTextContent('(delete)');
    expect(screen.getByTestId('play-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('delete-icon')).not.toBeInTheDocument();
  });

  it('renders the (delete) suffix and DeleteIcon in delete mode', () => {
    render(<AdapterYamlDialog {...baseProps} isDelete={true} />);
    expect(screen.getByTestId('dialog-title')).toHaveTextContent(
      'Istio Adapter - Custom YAML(delete)',
    );
    expect(screen.getByTestId('delete-icon')).toBeInTheDocument();
    expect(screen.queryByTestId('play-icon')).not.toBeInTheDocument();
  });

  it('renders namespace and version selects with the given values', () => {
    render(<AdapterYamlDialog {...baseProps} />);
    expect(screen.getByTestId('select-Namespace-value')).toHaveTextContent('default');
    expect(screen.getByTestId('select-Version-value')).toHaveTextContent('1.18');
  });

  it('forwards namespace error and version error flags', () => {
    render(<AdapterYamlDialog {...baseProps} namespaceError={true} versionError={true} />);
    expect(screen.getByTestId('select-Namespace')).toHaveAttribute('data-error', 'true');
    expect(screen.getByTestId('select-Version')).toHaveAttribute('data-error', 'true');
  });

  it('fires onNamespaceChange / onVersionChange when a selection is made', async () => {
    const user = userEvent.setup();
    const onNamespaceChange = vi.fn();
    const onVersionChange = vi.fn();
    render(
      <AdapterYamlDialog
        {...baseProps}
        onNamespaceChange={onNamespaceChange}
        onVersionChange={onVersionChange}
      />,
    );

    await user.click(screen.getByTestId('select-Namespace-choose-first'));
    expect(onNamespaceChange).toHaveBeenCalledWith(baseProps.namespaceList[0]);

    await user.click(screen.getByTestId('select-Version-choose-first'));
    expect(onVersionChange).toHaveBeenCalledWith(baseProps.versionList[0]);
  });

  it('fires onClose from the dialog title close handle', async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();
    render(<AdapterYamlDialog {...baseProps} onClose={onClose} />);

    await user.click(screen.getByText('close-from-title'));
    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('fires onApply when the apply IconButton is clicked', async () => {
    const user = userEvent.setup();
    const onApply = vi.fn();
    render(<AdapterYamlDialog {...baseProps} onApply={onApply} />);
    await user.click(screen.getByTestId('icon-Apply'));
    expect(onApply).toHaveBeenCalledTimes(1);
  });

  it('renders the CodeMirror with yaml mode and the provided value', () => {
    render(<AdapterYamlDialog {...baseProps} />);
    const editor = screen.getByTestId('codemirror') as HTMLTextAreaElement;
    expect(editor).toHaveAttribute('data-mode', 'text/x-yaml');
    expect(editor.value).toBe('kind: Deployment\n');
  });
});
