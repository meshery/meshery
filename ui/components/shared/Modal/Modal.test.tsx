import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import Modal, { RJSFModalWrapper } from './Modal';

const notify = vi.fn();

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

vi.mock('lib/event-types', () => ({
  EVENT_TYPES: { WARNING: 'warning' },
}));

vi.mock('@/assets/icons', () => ({
  ArrowDropDown: (props: any) => <svg data-testid="arrow-drop-down" {...props} />,
}));

vi.mock('../../meshery-mesh-interface/PatternService/RJSF_wrapper', () => ({
  default: ({ jsonSchema, formData, formRef, onChange }: any) => {
    if (formRef) {
      formRef.current = {
        state: { formData: { name: 'submission' } },
        validateForm: () => true,
      };
    }
    return (
      <div data-testid="rjsf-wrapper" data-schema={JSON.stringify(jsonSchema || {})}>
        <button onClick={() => onChange?.({ name: 'changed' })} aria-label="change-form">
          change
        </button>
      </div>
    );
  },
}));

vi.mock('../../meshery-mesh-interface/PatternService/helper', () => ({
  getSchema: (type: string) => ({ schemaType: type }),
}));

vi.mock('@sistent/sistent', () => ({
  IconButton: ({ children, onClick, ...props }: any) => (
    <button onClick={onClick} {...props}>
      {children}
    </button>
  ),
  Menu: ({ children, anchorEl, open }: any) =>
    open ? <div data-testid="menu">{children}</div> : null,
  MenuItem: ({ children, onClick, selected }: any) => (
    <li data-testid="menu-item" data-selected={String(selected)} onClick={onClick}>
      {children}
    </li>
  ),
  Tooltip: ({ children, title }: any) => (
    <span data-testid="tooltip" data-title={String(title)}>
      {children}
    </span>
  ),
  Typography: ({ children, variant }: any) => <p data-variant={variant}>{children}</p>,
  CircularProgress: () => <div data-testid="circular-progress" />,
  ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
  ModalFooter: ({ children, helpText, hasHelpText }: any) => (
    <div data-testid="modal-footer" data-help={helpText} data-has-help={String(hasHelpText)}>
      {children}
    </div>
  ),
  PrimaryActionButtons: ({
    primaryText,
    secondaryText,
    primaryButtonProps,
    secondaryButtonProps,
  }: any) => (
    <div data-testid="primary-action-buttons">
      <button
        data-testid="primary-button"
        onClick={primaryButtonProps.onClick}
        disabled={primaryButtonProps.disabled}
      >
        {primaryText}
      </button>
      <button data-testid="secondary-button" onClick={secondaryButtonProps.onClick}>
        {secondaryText}
      </button>
    </div>
  ),
  Modal: ({ children, open, closeModal, title }: any) =>
    open ? (
      <div data-testid="sistent-modal" data-title={title}>
        <button onClick={closeModal} aria-label="modal-close">
          close
        </button>
        {children}
      </div>
    ) : null,
  useTheme: () => ({ palette: { text: { primary: '#000' } } }),
}));

describe('Modal (RJSF version)', () => {
  beforeEach(() => {
    notify.mockReset();
  });

  it('renders a modal with title and schema when open', () => {
    render(
      <Modal
        open={true}
        title="My Modal"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('sistent-modal')).toHaveAttribute('data-title', 'My Modal');
    expect(screen.getByTestId('rjsf-wrapper')).toBeInTheDocument();
  });

  it('shows the loading spinner before schema is provided', () => {
    render(
      <Modal
        open={true}
        title="Loading Modal"
        schema={null}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('circular-progress')).toBeInTheDocument();
    expect(screen.queryByTestId('rjsf-wrapper')).not.toBeInTheDocument();
  });

  it('warns when the title contains forbidden words', () => {
    render(
      <Modal
        open={true}
        title="Untitled Design"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(notify).toHaveBeenCalledWith({
      event_type: 'warning',
      message: 'Design name should not contain Untitled Design, Untitled, LFX',
    });
  });

  it('disables the primary submit button when forbidden title is set', () => {
    render(
      <Modal
        open={true}
        title="lfx pattern"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('primary-button')).toBeDisabled();
  });

  it('submits and closes when primary button is clicked', async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    const handleClose = vi.fn();

    render(
      <Modal
        open={true}
        title="Normal Title"
        schema={{ type: 'object' }}
        handleClose={handleClose}
        handleSubmit={handleSubmit}
      />,
    );

    await user.click(screen.getByTestId('primary-button'));

    expect(handleClose).toHaveBeenCalled();
    expect(handleSubmit).toHaveBeenCalledWith({ name: 'submission' });
  });

  it('uses the schema_array dropdown when schema_array has length less than 1', () => {
    render(
      <Modal
        open={true}
        title="Schemed"
        schema={{ type: 'object' }}
        schema_array={[]}
        type="vN"
        schemaChangeHandler={vi.fn()}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('arrow-drop-down')).toBeInTheDocument();
  });

  it('uses custom submitBtnText if provided', () => {
    render(
      <Modal
        open={true}
        title="Normal"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
        submitBtnText="Send"
      />,
    );

    expect(screen.getByTestId('primary-button')).toHaveTextContent('Send');
  });

  it('falls back to Submit text', () => {
    render(
      <Modal
        open={true}
        title="Normal"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('primary-button')).toHaveTextContent('Submit');
  });

  it('renders help text in the footer', () => {
    render(
      <Modal
        open={true}
        title="HelpfulModal"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
        helpText="Need help?"
      />,
    );

    expect(screen.getByTestId('modal-footer')).toHaveAttribute('data-help', 'Need help?');
    expect(screen.getByTestId('modal-footer')).toHaveAttribute('data-has-help', 'true');
  });
});

describe('RJSFModalWrapper', () => {
  beforeEach(() => {
    notify.mockReset();
  });

  it('renders the wrapper with the RJSF form when schema is provided', () => {
    render(
      <RJSFModalWrapper
        title="Wrapper"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('rjsf-wrapper')).toBeInTheDocument();
  });

  it('shows a loader before schema is provided', () => {
    render(
      <RJSFModalWrapper
        title="Wrapper"
        schema={null}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(screen.getByTestId('circular-progress')).toBeInTheDocument();
  });

  it('submits and calls handleNext when both are provided', async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    const handleNext = vi.fn();

    render(
      <RJSFModalWrapper
        title="Wrapper"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={handleSubmit}
        handleNext={handleNext}
      />,
    );

    await user.click(screen.getByTestId('primary-button'));

    expect(handleSubmit).toHaveBeenCalledWith({ name: 'submission' });
    expect(handleNext).toHaveBeenCalled();
  });

  it('warns when the title contains forbidden words', () => {
    render(
      <RJSFModalWrapper
        title="My Untitled Design"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
      />,
    );

    expect(notify).toHaveBeenCalled();
  });

  it('closes when secondary button is clicked', async () => {
    const user = userEvent.setup();
    const handleClose = vi.fn();

    render(
      <RJSFModalWrapper
        title="Wrapper"
        schema={{ type: 'object' }}
        handleClose={handleClose}
        handleSubmit={vi.fn()}
      />,
    );

    await user.click(screen.getByTestId('secondary-button'));
    expect(handleClose).toHaveBeenCalled();
  });

  it('uses custom submitBtnText', () => {
    render(
      <RJSFModalWrapper
        title="Wrapper"
        schema={{ type: 'object' }}
        handleClose={vi.fn()}
        handleSubmit={vi.fn()}
        submitBtnText="Save Now"
      />,
    );

    expect(screen.getByTestId('primary-button')).toHaveTextContent('Save Now');
  });
});
