import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  useTheme: () => ({
    palette: {
      text: { error: '#b00020' },
      background: { error: { default: '#fdecea' } },
    },
  }),
  ErrorIcon: ({ fill }: { fill?: string }) => <svg data-testid="error-icon" data-fill={fill} />,
  Typography: ({ children, variant }: { children?: React.ReactNode; variant?: string }) => (
    <p data-variant={variant}>{children}</p>
  ),
  Button: ({ children, onClick }: { children?: React.ReactNode; onClick?: () => void }) => (
    <button type="button" onClick={onClick}>
      {children}
    </button>
  ),
}));

vi.mock('../style', () => ({
  ErrorContainer: ({ children }: { children?: React.ReactNode }) => (
    <div data-testid="error-container">{children}</div>
  ),
}));

import WidgetErrorFallback from './WidgetErrorFallback';

describe('WidgetErrorFallback', () => {
  it('renders the widget title in the message', () => {
    render(<WidgetErrorFallback widgetTitle="Getting Started" />);
    expect(screen.getByText('Unable to load Getting Started')).toBeInTheDocument();
  });

  it('renders a default message when none is provided', () => {
    render(<WidgetErrorFallback widgetTitle="My Designs" />);
    expect(
      screen.getByText(
        'Unable to load. Confirm your organization selection and role assignments are appropriate and try again.',
      ),
    ).toBeInTheDocument();
  });

  it('renders a custom message when provided', () => {
    render(<WidgetErrorFallback widgetTitle="Latest Blogs" message="Unable to reach the feed." />);
    expect(screen.getByText('Unable to reach the feed.')).toBeInTheDocument();
  });

  it('does not render a retry button when resetErrorBoundary is not provided', () => {
    render(<WidgetErrorFallback widgetTitle="My Designs" />);
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
  });

  it('calls resetErrorBoundary when Try Again is clicked', async () => {
    const resetErrorBoundary = vi.fn();
    render(
      <WidgetErrorFallback widgetTitle="My Designs" resetErrorBoundary={resetErrorBoundary} />,
    );
    await userEvent.click(screen.getByRole('button', { name: 'Try Again' }));
    expect(resetErrorBoundary).toHaveBeenCalledTimes(1);
  });
});
