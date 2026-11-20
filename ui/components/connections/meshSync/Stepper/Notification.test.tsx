import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import Notification from './Notification';

vi.mock('@/theme', () => ({
  useTheme: () => ({
    palette: {
      mode: 'light',
      success: { main: '#00B39F' },
      error: { main: '#F91484' },
      text: { secondary: 'rgba(0, 0, 0, 0.6)' },
    },
  }),
  alpha: (color: string, opacity: number) => `rgba(${color}, ${opacity})`,
}));

describe('MeshSync Stepper Notification', () => {
  it('renders success notification with accessible text and no retry button', () => {
    render(
      <Notification type="success" message="Credential for Kubernetes created" retry={false} />,
    );

    expect(screen.getByText('Success')).toBeDefined();
    expect(screen.getByText('Credential for Kubernetes created')).toBeDefined();
    expect(screen.queryByRole('button', { name: /retry/i })).toBeNull();
  });

  it('renders verification failed notification with retry button and triggers callback', () => {
    const handleRetry = vi.fn();
    render(
      <Notification
        type="error"
        message="Credential for Kubernetes verification failed"
        retry={true}
        onRetry={handleRetry}
      />,
    );

    expect(screen.getByText('Verification Failed')).toBeDefined();
    expect(screen.getByText('Credential for Kubernetes verification failed')).toBeDefined();

    const retryBtn = screen.getByRole('button', { name: /retry/i });
    expect(retryBtn).toBeDefined();

    fireEvent.click(retryBtn);
    expect(handleRetry).toHaveBeenCalledTimes(1);
  });

  it('omits retry button when retry prop is false in error state', () => {
    render(<Notification type="error" message="Verification failed" retry={false} />);

    expect(screen.getByText('Verification Failed')).toBeDefined();
    expect(screen.queryByRole('button', { name: /retry/i })).toBeNull();
  });
});
