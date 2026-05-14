import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import EnvironmentModal from './EnvironmentModal';

vi.mock('@sistent/sistent', () => ({
  Modal: ({ open, closeModal, title, headerIcon, maxWidth, children }: any) =>
    open ? (
      <div data-testid="modal" data-title={title} data-max-width={maxWidth}>
        <div data-testid="header-icon">{headerIcon}</div>
        <button onClick={closeModal} aria-label="close-modal">
          Close
        </button>
        {children}
      </div>
    ) : null,
  ModalBody: ({ children }: any) => <div data-testid="modal-body">{children}</div>,
  EnvironmentIcon: ({ height, width, fill }: any) => (
    <svg data-testid="environment-icon" data-height={height} data-width={width} data-fill={fill} />
  ),
  useTheme: () => ({
    palette: { background: { constant: { white: '#ffffff' } } },
  }),
  Box: ({ children, maxHeight }: any) => (
    <div data-testid="box" data-max-height={maxHeight}>
      {children}
    </div>
  ),
}));

vi.mock('../../lifecycle', () => ({
  EnvironmentComponent: () => <div data-testid="environment-component" />,
}));

describe('EnvironmentModal', () => {
  it('returns null when isOpenModal is false', () => {
    const { container } = render(<EnvironmentModal isOpenModal={false} setIsOpenModal={vi.fn()} />);

    expect(container.textContent).toBe('');
  });

  it('renders the Environment modal when open', () => {
    render(<EnvironmentModal isOpenModal={true} setIsOpenModal={vi.fn()} />);

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Environment');
    expect(screen.getByTestId('modal')).toHaveAttribute('data-max-width', 'xl');
    expect(screen.getByTestId('environment-icon')).toHaveAttribute('data-fill', '#ffffff');
    expect(screen.getByTestId('environment-component')).toBeInTheDocument();
  });

  it('closes the modal when close button is clicked', async () => {
    const user = userEvent.setup();
    const setIsOpenModal = vi.fn();
    render(<EnvironmentModal isOpenModal={true} setIsOpenModal={setIsOpenModal} />);

    await user.click(screen.getByRole('button', { name: 'close-modal' }));
    expect(setIsOpenModal).toHaveBeenCalledWith(false);
  });

  it('wraps the EnvironmentComponent in a Box with max height', () => {
    render(<EnvironmentModal isOpenModal={true} setIsOpenModal={vi.fn()} />);

    expect(screen.getByTestId('box')).toHaveAttribute('data-max-height', '65vh');
  });
});
