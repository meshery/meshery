import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children }: any) => <div data-testid="tooltip">{children}</div>,
  DeleteIcon: () => <svg data-testid="delete-icon" />,
  IconButton: ({ children, onClick, disabled }: any) => (
    <button type="button" onClick={onClick} disabled={disabled}>
      {children}
    </button>
  ),
}));

vi.mock('@/utils/can', () => ({
  default: () => true,
}));

import CustomToolbarSelect from './CustomToolbarSelect';

describe('CustomToolbarSelect', () => {
  it('renders a delete trigger button', () => {
    render(
      <CustomToolbarSelect
        selectedRows={{ data: [] }}
        patterns={[]}
        showModal={vi.fn()}
        deletePatterns={vi.fn()}
        setSelectedRows={vi.fn()}
      />,
    );
    expect(screen.getByTestId('delete-icon')).toBeInTheDocument();
  });

  it('asks for confirmation before delete and aborts on "no"', async () => {
    const showModal = vi.fn().mockResolvedValue('no');
    const deletePatterns = vi.fn();
    const setSelectedRows = vi.fn();

    render(
      <CustomToolbarSelect
        selectedRows={{ data: [{ index: 0 }] }}
        patterns={[{ id: 'a', name: 'Design A' }]}
        showModal={showModal}
        deletePatterns={deletePatterns}
        setSelectedRows={setSelectedRows}
      />,
    );

    fireEvent.click(screen.getByRole('button'));
    await waitFor(() => expect(showModal).toHaveBeenCalled());
    expect(deletePatterns).not.toHaveBeenCalled();
  });

  it('forwards delete to the supplied callback when confirmed', async () => {
    const showModal = vi.fn().mockResolvedValue('yes');
    const deletePatterns = vi.fn().mockReturnValue(Promise.resolve());
    const setSelectedRows = vi.fn();

    render(
      <CustomToolbarSelect
        selectedRows={{ data: [{ index: 0 }, { index: 1 }] }}
        patterns={[
          { id: 'a', name: 'Design A' },
          { id: 'b', name: 'Design B' },
        ]}
        showModal={showModal}
        deletePatterns={deletePatterns}
        setSelectedRows={setSelectedRows}
      />,
    );

    fireEvent.click(screen.getByRole('button'));
    await waitFor(() => expect(deletePatterns).toHaveBeenCalled());

    expect(deletePatterns).toHaveBeenCalledWith({
      patterns: [
        { id: 'a', name: 'Design A' },
        { id: 'b', name: 'Design B' },
      ],
    });

    await waitFor(() => expect(setSelectedRows).toHaveBeenCalledWith([]));
  });
});
