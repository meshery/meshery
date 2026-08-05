import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import MesheryVersionCompatibilityNotice from './MesheryVersionCompatibilityNotice';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Typography: ({ children, ...props }: any) => <p {...props}>{children}</p>,
  Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
  Chip: ({ label, ...props }: any) => <span {...props}>{label}</span>,
  IconButton: ({ children, ...props }: any) => <button {...props}>{children}</button>,
  CustomTooltip: ({ children }: any) => <div>{children}</div>,
  useTheme: () => ({
    spacing: (val: number) => `${val * 8}px`,
    palette: {
      mode: 'light',
      warning: { main: '#ed6c02' },
      success: { main: '#2e7d32' },
      text: { primary: '#000' },
    },
  }),
}));

describe('MesheryVersionCompatibilityNotice', () => {
  it('renders title and default version details', () => {
    render(<MesheryVersionCompatibilityNotice />);

    expect(screen.getByText(/Version Compatibility Notice/i)).toBeInTheDocument();
    expect(screen.getByText(/Current: v0.7.0/i)).toBeInTheDocument();
    expect(screen.getByText(/Required: v0.7.1\+/i)).toBeInTheDocument();
    expect(screen.getByText(/mesheryctl system update/i)).toBeInTheDocument();
  });

  it('renders custom prop values correctly', () => {
    render(
      <MesheryVersionCompatibilityNotice
        currentVersion="v0.6.9"
        requiredVersion="v0.7.0"
        componentName="Meshery Adapter"
        upgradeCommand="mesheryctl system restart"
      />,
    );

    expect(screen.getByText(/Meshery Adapter/i)).toBeInTheDocument();
    expect(screen.getByText(/Current: v0.6.9/i)).toBeInTheDocument();
    expect(screen.getByText(/Required: v0.7.0/i)).toBeInTheDocument();
    expect(screen.getByText(/mesheryctl system restart/i)).toBeInTheDocument();
  });

  it('copies upgrade command when copy button is clicked', async () => {
    const user = userEvent.setup();
    const writeTextMock = vi.fn();
    Object.assign(navigator, {
      clipboard: {
        writeText: writeTextMock,
      },
    });

    render(<MesheryVersionCompatibilityNotice upgradeCommand="mesheryctl system update" />);

    const copyBtn = screen.getByRole('button', { name: /copy upgrade command/i });
    await user.click(copyBtn);

    expect(writeTextMock).toHaveBeenCalledWith('mesheryctl system update');
  });
});
