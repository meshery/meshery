import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import TroubleshootingModal from './TroubleshootingModalComponent';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };

  return {
    Typography: ({ children, variant, ...rest }: any) => (
      <span data-variant={variant} {...rest}>
        {children}
      </span>
    ),
    Accordion: ({ children, expanded, onChange }: any) => (
      <div data-testid="accordion" data-expanded={String(expanded)}>
        <button onClick={(e) => onChange?.(e, !expanded)} aria-label="accordion-toggle">
          toggle
        </button>
        {children}
      </div>
    ),
    AccordionDetails: ({ children }: any) => <div data-testid="accordion-details">{children}</div>,
    AccordionSummary: ({ children, expandIcon }: any) => (
      <div data-testid="accordion-summary">
        {expandIcon}
        {children}
      </div>
    ),
    Paper: ({ children, elevation, square }: any) => (
      <div data-testid="paper" data-elevation={elevation} data-square={String(square)}>
        {children}
      </div>
    ),
    IconButton: ({ children, onClick, ...rest }: any) => (
      <button onClick={onClick} {...rest}>
        {children}
      </button>
    ),
    InfoIcon: () => <svg data-testid="info-icon" />,
    LIGHT_TEAL: '#0aa',
    Modal: ({ open, onClose, children }: any) =>
      open ? (
        <div data-testid="ts-modal">
          <button onClick={onClose} aria-label="close-modal">
            close
          </button>
          {children}
        </div>
      ) : null,
    keyframes: () => 'mocked-keyframes',
    styled,
  };
});

vi.mock('@/assets/icons', () => ({
  ExpandMore: () => <svg data-testid="expand-more" />,
  Close: (props: any) => <svg data-testid="close-icon" {...props} />,
}));

describe('TroubleshootingModal', () => {
  it('renders nothing when not open', () => {
    render(<TroubleshootingModal open={false} setOpen={vi.fn()} />);
    expect(screen.queryByTestId('ts-modal')).not.toBeInTheDocument();
  });

  it('shows the modal title when open', () => {
    render(<TroubleshootingModal open={true} setOpen={vi.fn()} />);
    expect(screen.getByText('Extensions Troubleshooting Guide')).toBeInTheDocument();
  });

  it('calls setOpen(false) when the close button is clicked', async () => {
    const user = userEvent.setup();
    const setOpen = vi.fn();
    render(<TroubleshootingModal open={true} setOpen={setOpen} />);

    const closeBtn = screen.getByTestId('close-icon').closest('button');
    expect(closeBtn).toBeInTheDocument();
    if (closeBtn) await user.click(closeBtn);

    expect(setOpen).toHaveBeenCalledWith(false);
  });

  it('renders 3 accordions for Stale Data, Missing Data, Additional Resources', () => {
    render(<TroubleshootingModal open={true} setOpen={vi.fn()} />);

    expect(screen.getAllByTestId('accordion')).toHaveLength(3);
    expect(screen.getByText('Stale Data')).toBeInTheDocument();
    expect(screen.getByText('Missing Data')).toBeInTheDocument();
    expect(screen.getByText('Additional Resources')).toBeInTheDocument();
  });

  it('opens the missing data panel when error mentions data missing', () => {
    render(
      <TroubleshootingModal
        open={true}
        setOpen={vi.fn()}
        viewHeaderErrorMessage="data missing for namespace"
      />,
    );

    // The second accordion should be expanded by default
    const accordions = screen.getAllByTestId('accordion');
    expect(accordions[1]).toHaveAttribute('data-expanded', 'true');
  });

  it('expands a panel when toggled', async () => {
    const user = userEvent.setup();
    render(<TroubleshootingModal open={true} setOpen={vi.fn()} />);

    const accordions = screen.getAllByTestId('accordion');
    expect(accordions[0]).toHaveAttribute('data-expanded', 'false');

    const toggleButtons = screen.getAllByRole('button', { name: 'accordion-toggle' });
    await user.click(toggleButtons[0]);

    expect(screen.getAllByTestId('accordion')[0]).toHaveAttribute('data-expanded', 'true');
  });

  it('renders the footer help links', () => {
    render(<TroubleshootingModal open={true} setOpen={vi.fn()} />);
    expect(screen.getByText('email')).toBeInTheDocument();
    expect(screen.getByText('community forum')).toBeInTheDocument();
  });
});
