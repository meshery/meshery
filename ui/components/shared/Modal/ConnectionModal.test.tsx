import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import ConnectionModal from './ConnectionModal';

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

  it('forwards the meshsync and connection metadata state to the table', () => {
    const controllerState = { status: 'running' };
    const metadataState = { foo: 'bar' };

    render(
      <ConnectionModal
        isOpenModal={true}
        setIsOpenModal={vi.fn()}
        meshsyncControllerState={controllerState}
        connectionMetadataState={metadataState}
      />,
    );

    expect(screen.getByTestId('connection-table')).toHaveAttribute(
      'data-controller',
      JSON.stringify(controllerState),
    );
    expect(screen.getByTestId('connection-table')).toHaveAttribute(
      'data-metadata',
      JSON.stringify(metadataState),
    );
  });
});
