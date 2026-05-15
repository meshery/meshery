import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import EnvironmentFormModal from '../../environments/EnvironmentFormModal';

vi.mock('@sistent/sistent', () => ({
  EnvironmentIcon: ({ height, width, fill }: any) => (
    <svg data-testid="environment-icon" data-height={height} data-width={width} data-fill={fill} />
  ),
}));

vi.mock('@/theme', () => ({
  styled:
    (_Component: any) =>
    () =>
    ({ children, ...props }: any) => <div {...props}>{children}</div>,
  useTheme: () => ({
    palette: { background: { constant: { white: '#ffffff' } } },
  }),
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ isOpen, onClose, title, headerIcon, size, children }: any) =>
    isOpen ? (
      <div data-testid="modal" data-title={title} data-size={size}>
        <div data-testid="header-icon">{headerIcon}</div>
        <button onClick={onClose} aria-label="close-modal">
          Close
        </button>
        {children}
      </div>
    ) : null,
}));

vi.mock('../../lifecycle', () => ({
  EnvironmentComponent: () => <div data-testid="environment-component" />,
}));

describe('EnvironmentFormModal', () => {
  it('returns null when isOpenModal is false', () => {
    const { container } = render(
      <EnvironmentFormModal isOpenModal={false} setIsOpenModal={vi.fn()} />,
    );
    expect(container.textContent).toBe('');
  });

  it('renders the Environment modal when open', () => {
    render(<EnvironmentFormModal isOpenModal={true} setIsOpenModal={vi.fn()} />);

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Environment');
    expect(screen.getByTestId('modal')).toHaveAttribute('data-size', 'xl');
    expect(screen.getByTestId('environment-icon')).toHaveAttribute('data-fill', '#ffffff');
    expect(screen.getByTestId('environment-component')).toBeInTheDocument();
  });

  it('closes the modal when close button is clicked', async () => {
    const user = userEvent.setup();
    const setIsOpenModal = vi.fn();
    render(<EnvironmentFormModal isOpenModal={true} setIsOpenModal={setIsOpenModal} />);

    await user.click(screen.getByRole('button', { name: 'close-modal' }));
    expect(setIsOpenModal).toHaveBeenCalledWith(false);
  });
});
