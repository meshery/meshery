import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import ConnectionModal from '../../connections/ConnectionFormModal';

vi.mock('@/theme', () => ({
  styled: (Component: any) => () => {
    const Styled = ({ children, ...props }: any) =>
      typeof Component === 'string' ? (
        React.createElement(Component, props, children)
      ) : (
        <Component {...props}>{children}</Component>
      );
    return Styled;
  },
}));

vi.mock('@/components/shared/Modal', () => ({
  Modal: ({ isOpen, onClose, title, headerIcon, size, children }: any) =>
    isOpen ? (
      <div data-testid="modal" data-title={title} data-max-width={size}>
        <div data-testid="header-icon">{headerIcon}</div>
        <button onClick={onClose} aria-label="close-modal">
          Close
        </button>
        {children}
      </div>
    ) : null,
}));

vi.mock('@/assets/icons/Connection', () => ({
  default: ({ height, width }: any) => (
    <svg data-testid="connection-icon" data-height={height} data-width={width} />
  ),
}));

vi.mock('../../connections/ConnectionTable', () => ({
  default: ({ meshsyncControllerState, connectionMetadataState, selectedFilter }: any) => (
    <div
      data-testid="connection-table"
      data-selected-filter={selectedFilter}
      data-controller={JSON.stringify(meshsyncControllerState || {})}
      data-metadata={JSON.stringify(connectionMetadataState || {})}
    />
  ),
}));

describe('ConnectionModal', () => {
  it('renders nothing when isOpenModal is false', () => {
    const { container } = render(
      <ConnectionModal
        isOpenModal={false}
        setIsOpenModal={vi.fn()}
        meshsyncControllerState={{}}
        connectionMetadataState={{}}
      />,
    );

    expect(container.textContent).toBe('');
  });

  it('renders the modal with a connection table when isOpenModal is true', () => {
    render(
      <ConnectionModal
        isOpenModal={true}
        setIsOpenModal={vi.fn()}
        meshsyncControllerState={{ status: 'ok' }}
        connectionMetadataState={{ count: 1 }}
      />,
    );

    expect(screen.getByTestId('modal')).toHaveAttribute('data-title', 'Connections');
    expect(screen.getByTestId('modal')).toHaveAttribute('data-max-width', 'xl');
    expect(screen.getByTestId('connection-icon')).toBeInTheDocument();
    expect(screen.getByTestId('connection-table')).toHaveAttribute(
      'data-selected-filter',
      'kubernetes',
    );
  });

  it('invokes setIsOpenModal with false when close is triggered', async () => {
    const user = userEvent.setup();
    const setIsOpenModal = vi.fn();

    render(
      <ConnectionModal
        isOpenModal={true}
        setIsOpenModal={setIsOpenModal}
        meshsyncControllerState={{}}
        connectionMetadataState={{}}
      />,
    );

    await user.click(screen.getByRole('button', { name: 'close-modal' }));
    expect(setIsOpenModal).toHaveBeenCalledWith(false);
  });

  it('keeps rendering the connection table when compatibility props are passed', () => {
    render(
      <ConnectionModal
        isOpenModal={true}
        setIsOpenModal={vi.fn()}
        meshsyncControllerState={{ status: 'running' }}
        connectionMetadataState={{ foo: 'bar' }}
      />,
    );

    expect(screen.getByTestId('connection-table')).toHaveAttribute(
      'data-selected-filter',
      'kubernetes',
    );
  });
});
