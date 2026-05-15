import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import ExportModal from './ExportModal';

vi.mock('@sistent/sistent', () => ({
  Modal: ({ open, onClose, closeModal, title, children }: any) =>
    open ? (
      <div data-testid="modal" data-title={title}>
        <button onClick={onClose ?? closeModal} aria-label="close-modal">
          Close
        </button>
        {children}
      </div>
    ) : null,
  ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
  ListItem: ({ children, sx }: any) => (
    <li data-testid="export-list-item" data-sx={JSON.stringify(sx || {})}>
      {children}
    </li>
  ),
  ListItemIcon: ({ children }: any) => <span data-testid="list-icon">{children}</span>,
  ListItemText: ({ primary }: any) => <span data-testid="list-text">{primary}</span>,
  InfoTooltip: ({ title }: any) => <span data-testid="info-tooltip">{String(title ?? '')}</span>,
  IconButton: ({ children, onClick, disabled, sx }: any) => (
    <button
      onClick={onClick}
      disabled={disabled}
      data-testid="icon-button"
      data-sx={JSON.stringify(sx || {})}
    >
      {children}
    </button>
  ),
  DownloadIcon: () => <svg data-testid="download-icon" />,
  Box: ({ children, sx }: any) => <div data-sx={JSON.stringify(sx || {})}>{children}</div>,
  useTheme: () => ({
    palette: {
      primary: { main: '#1a73e8' },
      border: { normal: '#eee', brand: '#0aa' },
      icon: { default: '#444' },
      common: { white: '#fff' },
    },
    shadows: ['none', '0px 1px 1px rgba(0,0,0,0.1)'],
    spacing: (n: number) => `${n * 8}px`,
  }),
  useModal: (props: any) => ({ headerIcon: props?.headerIcon ?? 'icon' }),
}));

vi.mock('@/assets/icons/technology/kubernetes', () => ({
  default: (props: any) => <svg data-testid="k8s-icon" {...props} />,
}));

vi.mock('@/assets/icons/technology/HelmIcon', () => ({
  default: (props: any) => <svg data-testid="helm-icon" {...props} />,
}));

vi.mock('@/assets/icons/Pattern', () => ({
  default: (props: any) => <svg data-testid="pattern-icon" {...props} />,
}));

vi.mock('@/assets/icons/OciImage', () => ({
  OCIImageIcon: (props: any) => <svg data-testid="oci-icon" {...props} />,
}));

describe('ExportModal', () => {
  const baseProps = {
    downloadModal: { open: true, content: { id: 'design-1' } },
    handleDownloadDialogClose: vi.fn(),
    handleDesignDownload: vi.fn(),
  };

  it('renders the export modal with the 4 default options', () => {
    render(<ExportModal {...baseProps} />);

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Export Design as...');
    expect(screen.getByText('Meshery Design (yaml)')).toBeInTheDocument();
    expect(screen.getByText('Meshery Design (OCI image)')).toBeInTheDocument();
    expect(screen.getByText('Kubernetes Manifest (yaml)')).toBeInTheDocument();
    expect(screen.getByText('Helm Chart (tar.gz)')).toBeInTheDocument();
    expect(screen.getAllByTestId('export-list-item')).toHaveLength(4);
  });

  it('returns null modal when open is false', () => {
    render(<ExportModal {...baseProps} downloadModal={{ open: false, content: null }} />);

    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('triggers download with the content for yaml option', async () => {
    const user = userEvent.setup();
    const handleDesignDownload = vi.fn();

    render(<ExportModal {...baseProps} handleDesignDownload={handleDesignDownload} />);

    const downloadButtons = screen.getAllByTestId('icon-button');
    await user.click(downloadButtons[0]);

    expect(handleDesignDownload).toHaveBeenCalledWith(expect.anything(), { id: 'design-1' });
  });

  it('triggers oci=true for the OCI image option', async () => {
    const user = userEvent.setup();
    const handleDesignDownload = vi.fn();

    render(<ExportModal {...baseProps} handleDesignDownload={handleDesignDownload} />);

    const downloadButtons = screen.getAllByTestId('icon-button');
    await user.click(downloadButtons[1]);

    expect(handleDesignDownload).toHaveBeenLastCalledWith(
      expect.anything(),
      { id: 'design-1' },
      null,
      'oci=true',
    );
  });

  it('triggers Kubernetes Manifest export with the right query', async () => {
    const user = userEvent.setup();
    const handleDesignDownload = vi.fn();

    render(<ExportModal {...baseProps} handleDesignDownload={handleDesignDownload} />);

    const downloadButtons = screen.getAllByTestId('icon-button');
    await user.click(downloadButtons[2]);

    expect(handleDesignDownload).toHaveBeenLastCalledWith(
      expect.anything(),
      { id: 'design-1' },
      null,
      'export=Kubernetes Manifest',
    );
  });

  it('triggers helm chart export', async () => {
    const user = userEvent.setup();
    const handleDesignDownload = vi.fn();

    render(<ExportModal {...baseProps} handleDesignDownload={handleDesignDownload} />);

    const downloadButtons = screen.getAllByTestId('icon-button');
    await user.click(downloadButtons[3]);

    expect(handleDesignDownload).toHaveBeenLastCalledWith(
      expect.anything(),
      { id: 'design-1' },
      null,
      'export=helm-chart',
    );
  });

  it('appends extension export options', () => {
    const extraOption = {
      title: 'Extra Export',
      icon: <svg data-testid="extra-icon" />,
      onClick: vi.fn(),
      description: 'Extra',
    };

    render(<ExportModal {...baseProps} extensionExportOptions={[extraOption]} />);

    expect(screen.getByText('Extra Export')).toBeInTheDocument();
    expect(screen.getAllByTestId('export-list-item')).toHaveLength(5);
  });
});
