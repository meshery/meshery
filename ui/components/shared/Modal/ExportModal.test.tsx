import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import ExportDesignModal from '../../designs/export/ExportDesignModal';

vi.mock('@sistent/sistent', () => ({
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
}));

vi.mock('@/theme', () => ({
  styled:
    (Component: any) =>
    () =>
    ({ children, ...props }: any) =>
      typeof Component === 'string' ? (
        React.createElement(Component, props, children)
      ) : (
        <Component {...props}>{children}</Component>
      ),
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
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ isOpen, onClose, title, children }: any) =>
    isOpen ? (
      <div data-testid="modal" data-title={title}>
        <button onClick={onClose} aria-label="close-modal">
          Close
        </button>
        {children}
      </div>
    ) : null,
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

describe('ExportDesignModal', () => {
  const baseProps = {
    downloadModal: { open: true, content: { id: 'design-1' } },
    handleDownloadDialogClose: vi.fn(),
    handleDesignDownload: vi.fn(),
  };

  it('renders the export modal with the 4 default options', () => {
    render(<ExportDesignModal {...baseProps} />);

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Export Design as...');
    expect(screen.getByText('Meshery Design (yaml)')).toBeInTheDocument();
    expect(screen.getByText('Meshery Design (OCI image)')).toBeInTheDocument();
    expect(screen.getByText('Kubernetes Manifest (yaml)')).toBeInTheDocument();
    expect(screen.getByText('Helm Chart (tar.gz)')).toBeInTheDocument();
    expect(screen.getAllByTestId('export-list-item')).toHaveLength(4);
  });

  it('returns null modal when open is false', () => {
    render(<ExportDesignModal {...baseProps} downloadModal={{ open: false, content: null }} />);
    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('triggers download with the content for yaml option', async () => {
    const user = userEvent.setup();
    const handleDesignDownload = vi.fn();

    render(<ExportDesignModal {...baseProps} handleDesignDownload={handleDesignDownload} />);

    await user.click(screen.getAllByTestId('icon-button')[0]);
    expect(handleDesignDownload).toHaveBeenCalledWith(expect.anything(), { id: 'design-1' });
  });

  it('appends extension export options', () => {
    const extraOption = {
      title: 'Extra Export',
      icon: <svg data-testid="extra-icon" />,
      onClick: vi.fn(),
      description: 'Extra',
    };

    render(<ExportDesignModal {...baseProps} extensionExportOptions={[extraOption]} />);

    expect(screen.getByText('Extra Export')).toBeInTheDocument();
    expect(screen.getAllByTestId('export-list-item')).toHaveLength(5);
  });
});
