import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Troubleshoot from './TroubleshootingComponent';

const mockNotify = vi.fn();
const mockTriggerWebhook = vi.fn();

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: mockNotify }),
}));

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => ({
    data: { firstName: 'Test', lastName: 'User', email: 'test@example.com' },
  }),
}));

vi.mock('@/rtk-query/webhook', () => ({
  useSupportWebHookMutation: () => [mockTriggerWebhook, { isLoading: false }],
}));

vi.mock('@sistent/sistent', () => ({
  Button: ({ children, variant, color, 'aria-label': ariaLabel, onClick, ...props }: any) => (
    <button
      data-testid="sistent-button"
      aria-label={ariaLabel}
      data-variant={variant}
      data-color={color}
      onClick={onClick}
      {...props}
    >
      {children}
    </button>
  ),
  Modal: ({ children, open, title }: any) =>
    open ? (
      <div data-testid="support-modal" data-title={title}>
        {children}
      </div>
    ) : null,
  helpAndSupportModalSchema: {},
  helpAndSupportModalUiSchema: {},
}));

vi.mock('@/components/shared/Troubleshooting/TroubleshootingModal', () => ({
  default: ({ open }: any) => (open ? <div data-testid="troubleshooting-modal">Modal</div> : null),
}));

vi.mock('./shared/Modal/Modal', () => ({
  RJSFModalWrapper: () => <div data-testid="rjsf-modal-wrapper" />,
}));

vi.mock('./general/feedback', () => ({
  StickyFeedbackButton: () => <div data-testid="sticky-feedback-button" />,
}));

vi.mock('@/assets/icons/support', () => ({
  default: () => <span data-testid="support-icon" />,
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

describe('TroubleshootingComponent', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders troubleshooting guide and get help buttons with accessible ARIA labels', () => {
    render(<Troubleshoot errorMessage="Test error message" />);

    const guideButton = screen.getByRole('button', { name: /Open Troubleshooting Guide/i });
    const helpButton = screen.getByRole('button', { name: /Open Help and Support Form/i });

    expect(guideButton).toBeInTheDocument();
    expect(guideButton).toHaveTextContent('Troubleshooting Guide');
    expect(helpButton).toBeInTheDocument();
    expect(helpButton).toHaveTextContent('Get Help');
  });

  it('opens support form modal when Get Help button is clicked', async () => {
    const user = userEvent.setup();
    render(<Troubleshoot errorMessage="Test error message" />);

    expect(screen.queryByTestId('support-modal')).not.toBeInTheDocument();

    const helpButton = screen.getByRole('button', { name: /Open Help and Support Form/i });
    await user.click(helpButton);

    expect(screen.getByTestId('support-modal')).toBeInTheDocument();
  });
});
