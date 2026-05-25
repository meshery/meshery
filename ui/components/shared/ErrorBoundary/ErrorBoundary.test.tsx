import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import CustomErrorFallback from './ErrorBoundary';

const notify = vi.fn();
const supportWebhook = vi.fn().mockReturnValue({
  unwrap: () => Promise.resolve(),
});

let providerData: any = { providerType: 'remote' };
let userData: any = { firstName: 'A', lastName: 'B', email: 'a@b.c' };

vi.mock('@sistent/sistent', () => ({
  Fallback: ({ children, showPackageInfo }: any) => (
    <div data-testid="fallback" data-show-package-info={String(showPackageInfo)}>
      {children}
    </div>
  ),
  Modal: ({ open, closeModal, title, children }: any) =>
    open ? (
      <div data-testid="support-modal" data-title={title}>
        <button onClick={closeModal} aria-label="modal-close">
          close
        </button>
        {children}
      </div>
    ) : null,
  helpAndSupportModalSchema: { type: 'object' },
  helpAndSupportModalUiSchema: {},
  useTheme: () => ({ palette: { text: { default: '#000' } } }),
}));

vi.mock('../Modal/Modal', () => ({
  RJSFModalWrapper: ({ handleSubmit }: any) => (
    <button data-testid="rjsf-modal-wrapper" onClick={() => handleSubmit({ message: 'hello' })}>
      submit-help
    </button>
  ),
}));

vi.mock('@/assets/icons/support', () => ({
  default: (props: any) => <svg data-testid="support-icon" {...props} />,
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('@/rtk-query/webhook', () => ({
  useSupportWebHookMutation: () => [supportWebhook],
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { SUCCESS: 'success', ERROR: 'error' },
}));

vi.mock('@/rtk-query/user', () => ({
  useGetLoggedInUserQuery: () => ({ data: userData }),
  useGetProviderCapabilitiesQuery: () => ({ data: providerData }),
}));

vi.mock('../../general/style', () => ({
  EditButton: ({ children, onClick, ...rest }: any) => (
    <button onClick={onClick} {...rest} data-testid="edit-btn">
      {children}
    </button>
  ),
  FallbackWrapper: ({ children }: any) => <div data-testid="fallback-wrapper">{children}</div>,
  TextButton: ({ children }: any) => <span>{children}</span>,
  ToolBarButtonContainer: ({ children }: any) => <div>{children}</div>,
  TryAgainButton: ({ children, onClick, color }: any) => (
    <button onClick={onClick} data-color={color} data-testid="try-again-btn">
      {children}
    </button>
  ),
}));

vi.mock('../../general/feedback', () => ({
  StickyFeedbackButton: ({ defaultMessage, defaultOpen }: any) => (
    <div
      data-testid="sticky-feedback"
      data-message={defaultMessage}
      data-open={String(defaultOpen)}
    />
  ),
}));

describe('CustomErrorFallback', () => {
  beforeEach(() => {
    notify.mockReset();
    supportWebhook.mockClear();
    providerData = { providerType: 'remote' };
    userData = { firstName: 'A', lastName: 'B', email: 'a@b.c' };
  });

  it('renders the fallback and try again button', () => {
    const resetErrorBoundary = vi.fn();
    render(
      <CustomErrorFallback error={{ message: 'Boom!' }} resetErrorBoundary={resetErrorBoundary} />,
    );

    expect(screen.getByTestId('fallback')).toBeInTheDocument();
    expect(screen.getByTestId('try-again-btn')).toHaveTextContent('Try Again');
  });

  it('calls resetErrorBoundary when Try Again is clicked', async () => {
    const user = userEvent.setup();
    const resetErrorBoundary = vi.fn();

    render(
      <CustomErrorFallback error={{ message: 'Boom!' }} resetErrorBoundary={resetErrorBoundary} />,
    );

    await user.click(screen.getByTestId('try-again-btn'));
    expect(resetErrorBoundary).toHaveBeenCalledTimes(1);
  });

  it('shows Get Help button when provider is remote', () => {
    render(<CustomErrorFallback error={{ message: 'x' }} resetErrorBoundary={vi.fn()} />);
    expect(screen.getByTestId('edit-btn')).toHaveTextContent('Get Help');
  });

  it('hides Get Help button and sticky feedback when provider is not remote', () => {
    providerData = { providerType: 'local' };
    render(<CustomErrorFallback error={{ message: 'x' }} resetErrorBoundary={vi.fn()} />);
    expect(screen.queryByTestId('edit-btn')).not.toBeInTheDocument();
    expect(screen.queryByTestId('sticky-feedback')).not.toBeInTheDocument();
  });

  it('opens the support modal when Get Help is clicked', async () => {
    const user = userEvent.setup();
    render(<CustomErrorFallback error={{ message: 'x' }} resetErrorBoundary={vi.fn()} />);

    await user.click(screen.getByTestId('edit-btn'));
    expect(screen.getByTestId('support-modal')).toBeInTheDocument();
  });

  it('submits the support webhook with the user metadata', async () => {
    const user = userEvent.setup();
    render(<CustomErrorFallback error={{ message: 'x' }} resetErrorBoundary={vi.fn()} />);

    await user.click(screen.getByTestId('edit-btn'));
    await user.click(screen.getByTestId('rjsf-modal-wrapper'));

    expect(supportWebhook).toHaveBeenCalledWith({
      body: {
        memberFormOne: {
          message: 'hello',
          firstname: 'A',
          lastname: 'B',
          email: 'a@b.c',
        },
      },
      type: 'support',
    });
  });

  it('renders the sticky feedback button when remote', () => {
    Object.defineProperty(window, 'location', {
      value: { href: 'http://test/path' },
      writable: true,
    });
    render(<CustomErrorFallback error={{ message: 'b' }} resetErrorBoundary={vi.fn()} />);
    const sticky = screen.getByTestId('sticky-feedback');
    expect(sticky).toHaveAttribute('data-open', 'true');
    expect(sticky.getAttribute('data-message')).toContain('http://test/path');
  });

  it('forwards showPackageInfo true to the Fallback', () => {
    render(<CustomErrorFallback error={{ message: 'b' }} resetErrorBoundary={vi.fn()} />);
    expect(screen.getByTestId('fallback')).toHaveAttribute('data-show-package-info', 'true');
  });
});
