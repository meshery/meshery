import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  Dialog: ({ open, onClose, children, 'aria-labelledby': labelledBy }: any) =>
    open ? (
      <div role="dialog" aria-labelledby={labelledBy} data-testid="dialog">
        {children}
        <button data-testid="backdrop-close" onClick={onClose}>
          backdrop
        </button>
      </div>
    ) : null,
  DialogActions: ({ children }: any) => <div data-testid="actions">{children}</div>,
  DialogContent: ({ children }: any) => <div data-testid="content">{children}</div>,
  DialogContentText: ({ children }: any) => <div data-testid="content-text">{children}</div>,
  DialogTitle: ({ children, id, 'data-testid': testId }: any) => (
    <h2 id={id} data-testid={testId}>
      {children}
    </h2>
  ),
}));

import MesheryChartDialog from './MesheryChartDialog';

describe('MesheryChartDialog', () => {
  it('returns nothing when closed', () => {
    const { container } = render(
      <MesheryChartDialog
        open={false}
        title="Hello"
        content={<div>chart</div>}
        handleClose={vi.fn()}
      />,
    );
    expect(container.querySelector('[role="dialog"]')).toBeNull();
  });

  it('renders the provided title and content when open', () => {
    render(
      <MesheryChartDialog
        open={true}
        title="My Comparison"
        content={<div data-testid="content-body">chart content</div>}
        handleClose={vi.fn()}
      />,
    );
    expect(screen.getByTestId('chart-dialog-title')).toHaveTextContent('My Comparison');
    expect(screen.getByTestId('content-body')).toBeInTheDocument();
  });

  it('falls back to the default "Comparison" title when none is supplied', () => {
    render(
      <MesheryChartDialog open={true} title="" content={<div>x</div>} handleClose={vi.fn()} />,
    );
    expect(screen.getByTestId('chart-dialog-title')).toHaveTextContent('Comparison');
  });

  it('invokes handleClose when Close is clicked', () => {
    const handleClose = vi.fn();
    render(
      <MesheryChartDialog open={true} title="t" content={<div>x</div>} handleClose={handleClose} />,
    );
    fireEvent.click(screen.getByText('Close'));
    expect(handleClose).toHaveBeenCalled();
  });
});
