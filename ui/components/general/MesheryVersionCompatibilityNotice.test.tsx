import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useGetSystemVersionQuery } from '@/rtk-query/user';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children, ...props }: any) => <div {...props}>{children}</div>,
  Typography: ({ children, ...props }: any) => <p {...props}>{children}</p>,
  Button: ({ children, ...props }: any) => <button {...props}>{children}</button>,
  Chip: ({ label, ...props }: any) => <span {...props}>{label}</span>,
  IconButton: ({ children, ...props }: any) => <button {...props}>{children}</button>,
  CustomTooltip: ({ children }: any) => <div>{children}</div>,
  CheckIcon: (props: any) => <svg data-testid="check-icon" {...props} />,
  CopyIcon: (props: any) => <svg data-testid="copy-icon" {...props} />,
  WarningIcon: (props: any) => <svg data-testid="warning-icon" {...props} />,
  alpha: (color: string) => color,
  useTheme: () => ({
    spacing: (val: number) => `${val * 8}px`,
    palette: {
      mode: 'light',
      common: { white: '#fff', black: '#000' },
      warning: { main: '#ed6c02' },
      success: { main: '#2e7d32' },
      text: { primary: '#000' },
      divider: '#e0e0e0',
      grey: { 100: '#f5f5f5', 900: '#1e1e1e' },
    },
  }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetSystemVersionQuery: vi.fn(),
}));

import MesheryVersionCompatibilityNotice from './MesheryVersionCompatibilityNotice';

const mockedUseGetSystemVersionQuery = vi.mocked(useGetSystemVersionQuery);

describe('MesheryVersionCompatibilityNotice', () => {
  beforeEach(() => {
    mockedUseGetSystemVersionQuery.mockReset();
    mockedUseGetSystemVersionQuery.mockReturnValue({ data: undefined } as never);
  });
  it('renders fallback version when no prop and no system data are present', () => {
    render(<MesheryVersionCompatibilityNotice />);

    expect(screen.getByText(/Version Compatibility Notice/i)).toBeInTheDocument();
    // Neither currentVersion prop nor system data -> hardcoded fallback
    expect(screen.getByText(/Current: v0.7.0/i)).toBeInTheDocument();
    expect(screen.getByText(/Required: v0.7.1\+/i)).toBeInTheDocument();
    expect(screen.getByText(/mesheryctl system update/i)).toBeInTheDocument();
  });

  it('uses the runtime Meshery version when no prop is passed but system data is available', () => {
    mockedUseGetSystemVersionQuery.mockReturnValue({ data: { build: 'v1.2.3' } } as never);

    render(<MesheryVersionCompatibilityNotice />);

    // No explicit prop -> falls back to system data
    expect(screen.getByText(/Current: v1.2.3/i)).toBeInTheDocument();
  });

  it('prefers explicit currentVersion prop over system data', () => {
    mockedUseGetSystemVersionQuery.mockReturnValue({ data: { build: 'v9.9.9' } } as never);

    render(<MesheryVersionCompatibilityNotice currentVersion="v0.6.9" />);

    // Explicit prop wins; system version should NOT appear
    expect(screen.getByText(/Current: v0.6.9/i)).toBeInTheDocument();
    expect(screen.queryByText(/v9.9.9/)).not.toBeInTheDocument();
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

  it('copies upgrade command when copy button is clicked', () => {
    const writeTextMock = vi.fn().mockResolvedValue(undefined);
    Object.assign(navigator, {
      clipboard: {
        writeText: writeTextMock,
      },
    });

    render(<MesheryVersionCompatibilityNotice upgradeCommand="mesheryctl system update" />);

    const copyBtn = screen.getByRole('button', { name: /copy upgrade command/i });
    fireEvent.click(copyBtn);

    return waitFor(() => {
      expect(writeTextMock).toHaveBeenCalledWith('mesheryctl system update');
    });
  });

  it('does not show copied state when clipboard write fails', async () => {
    const writeTextMock = vi.fn().mockRejectedValue(new Error('denied'));
    Object.assign(navigator, {
      clipboard: {
        writeText: writeTextMock,
      },
    });

    render(<MesheryVersionCompatibilityNotice upgradeCommand="mesheryctl system update" />);

    const copyBtn = screen.getByRole('button', { name: /copy upgrade command/i });
    fireEvent.click(copyBtn);

    await waitFor(() => {
      expect(writeTextMock).toHaveBeenCalledWith('mesheryctl system update');
    });

    await waitFor(() => {
      expect(screen.queryByText(/Copied!/i)).not.toBeInTheDocument();
    });
  });
});
