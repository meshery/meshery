import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import GenericModal from './GenericModal';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  return {
    Modal: ({ children, open, onClose, container, maxWidth }: any) =>
      open ? (
        <div
          data-testid="modal"
          data-open={String(open)}
          data-max-width={maxWidth}
          data-has-container={container ? 'true' : 'false'}
          onClick={onClose}
        >
          {children}
        </div>
      ) : null,
    Backdrop: ({ children }: any) => <div data-testid="backdrop">{children}</div>,
    Box: ({ children, sx }: any) => (
      <div data-testid="box" data-sx={JSON.stringify(sx || {})}>
        {children}
      </div>
    ),
    Fade: ({ children, in: inProp }: any) => (
      <div data-testid="fade" data-in={String(inProp)}>
        {children}
      </div>
    ),
    IconButton: ({ children, onClick }: any) => (
      <button data-testid="close-button" onClick={onClick}>
        {children}
      </button>
    ),
    CloseIcon: () => <span data-testid="close-icon">X</span>,
    styled,
  };
});

describe('GenericModal', () => {
  it('renders the content when open is true', () => {
    render(
      <GenericModal
        open={true}
        Content={<div data-testid="modal-content">Hello</div>}
        handleClose={vi.fn()}
        container={null}
      />,
    );

    expect(screen.getByTestId('modal')).toBeInTheDocument();
    expect(screen.getByTestId('modal-content')).toHaveTextContent('Hello');
    expect(screen.getByTestId('fade')).toHaveAttribute('data-in', 'true');
  });

  it('renders nothing when open is false', () => {
    render(
      <GenericModal
        open={false}
        Content={<div data-testid="modal-content">Hello</div>}
        handleClose={vi.fn()}
        container={null}
      />,
    );

    expect(screen.queryByTestId('modal')).not.toBeInTheDocument();
  });

  it('passes container reference through to the underlying modal', () => {
    const container = document.createElement('div');
    render(
      <GenericModal
        open={true}
        Content={<div>Body</div>}
        handleClose={vi.fn()}
        container={container}
      />,
    );

    expect(screen.getByTestId('modal')).toHaveAttribute('data-has-container', 'true');
  });

  it('sets max width to lg', () => {
    render(<GenericModal open={true} Content={<span>Body</span>} handleClose={vi.fn()} />);

    expect(screen.getByTestId('modal')).toHaveAttribute('data-max-width', 'lg');
  });

  it('renders the close button and calls handleClose when clicked', () => {
    const handleClose = vi.fn();
    render(<GenericModal open={true} Content={<span>Body</span>} handleClose={handleClose} />);

    const closeBtn = screen.getByTestId('close-button');
    expect(closeBtn).toBeInTheDocument();

    closeBtn.click();
    expect(handleClose).toHaveBeenCalledTimes(1);
  });
});
