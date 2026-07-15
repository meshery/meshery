import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Typography: ({ children, variant }: any) => <div data-variant={variant}>{children}</div>,
  InfoCircleIcon: () => <svg data-testid="info-icon" />,
  CustomTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
  useTheme: () => ({
    palette: {
      mode: 'light',
      text: { default: '#222' },
    },
  }),
}));

vi.mock('@/constants/common', () => ({
  ErrorTypes: { UNKNOWN: 'UNKNOWN', PERMISSION: 'PERMISSION' },
}));

vi.mock('./OrgSwitcher', () => ({
  default: () => <div data-testid="org-switcher" />,
}));

vi.mock('./CurrentSession', () => ({
  default: () => <div data-testid="current-session" />,
}));

vi.mock('./styles', () => ({
  ErrorMain: ({ children }: any) => <main>{children}</main>,
  ErrorContainer: ({ children }: any) => <section>{children}</section>,
  ErrorSection: ({ children }: any) => <div data-testid="error-section">{children}</div>,
  ErrorSectionContainer: ({ children }: any) => <div>{children}</div>,
  ErrorSectionContent: ({ children }: any) => <div>{children}</div>,
  ErrorContentContainer: ({ children }: any) => (
    <div data-testid="content-container">{children}</div>
  ),
  StyledButton: ({ children, href }: any) => <a href={href}>{children}</a>,
  IconWrapper: ({ children }: any) => <div>{children}</div>,
  StyledDivider: () => <hr />,
  ErrorLink: ({ children, href }: any) => <a href={href}>{children}</a>,
}));

import DefaultError from './index';

describe('DefaultError', () => {
  it('renders the default permissions message when no errorTitle is provided', () => {
    render(<DefaultError />);
    expect(screen.getByText(/you don't have the required permissions/i)).toBeInTheDocument();
    expect(screen.getByTestId('org-switcher')).toBeInTheDocument();
    expect(screen.getByTestId('current-session')).toBeInTheDocument();
  });

  it('renders custom error title when provided', () => {
    render(<DefaultError errorTitle="Something broke" />);
    expect(screen.getByText('Something broke')).toBeInTheDocument();
  });

  it('renders the unknown error content when errorType is UNKNOWN', () => {
    render(<DefaultError errorType="UNKNOWN" errorContent="server is down" />);
    expect(screen.getByTestId('content-container')).toHaveTextContent('server is down');
  });

  it('does not render the unknown error content for other error types', () => {
    render(<DefaultError errorType="PERMISSION" />);
    expect(screen.queryByTestId('content-container')).not.toBeInTheDocument();
  });

  it('renders a link back to the dashboard and community resources', () => {
    render(<DefaultError />);
    const dashboard = screen.getByText(/Return to Dashboard/i).closest('a');
    expect(dashboard).toHaveAttribute('href', '/');

    const forumLink = screen.getByText('discussion forum').closest('a');
    expect(forumLink).toHaveAttribute('href', expect.stringContaining('meshery.io/community'));
    const slackLink = screen.getByText('Slack workspace').closest('a');
    expect(slackLink).toHaveAttribute('href', 'https://slack.meshery.io');
  });
});
