import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  IconButton: ({ children, onClick, disabled }: any) => (
    <button data-testid="icon-button" onClick={onClick} disabled={!!disabled}>
      {children}
    </button>
  ),
  Accordion: ({ children }: any) => <div data-testid="accordion">{children}</div>,
  AccordionSummary: ({ children }: any) => <div data-testid="summary">{children}</div>,
  AccordionDetails: ({ children }: any) => <div data-testid="details">{children}</div>,
  Typography: ({ children }: any) => <span>{children}</span>,
  Box: ({ children }: any) => <div>{children}</div>,
  useTheme: () => ({
    typography: { pxToRem: (n: number) => `${n}px`, fontWeightRegular: 400 },
  }),
  ExpandMoreIcon: () => <svg data-testid="expand-icon" />,
  DeleteIcon: () => <svg data-testid="delete-icon" />,
}));

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {} }));

vi.mock('../helper', () => ({
  safeStringTitle: (v: any) => (v == null ? '' : String(v)),
}));

import SimpleAccordion from './Accordion';

describe('SimpleAccordion', () => {
  it('capitalizes the heading and renders children content', () => {
    render(
      <SimpleAccordion heading="hello world" childProps={{ hasRemove: false }}>
        <div data-testid="child">child content</div>
      </SimpleAccordion>,
    );
    expect(screen.getByText(/Hello world/)).toBeInTheDocument();
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('renders a delete button when hasRemove is true and invokes onRemoveItem on click', () => {
    const onRemoveItem = vi.fn();
    render(
      <SimpleAccordion heading="thing" childProps={{ hasRemove: true, onRemoveItem }}>
        <div>child</div>
      </SimpleAccordion>,
    );
    expect(screen.getByTestId('delete-icon')).toBeInTheDocument();
    // The old contract called onDropIndexClick(index) during render; the v6
    // contract passes onRemoveItem (already bound to this item's index by
    // RJSF) as the click handler, so it must fire only on click.
    expect(onRemoveItem).not.toHaveBeenCalled();
    fireEvent.click(screen.getByTestId('icon-button'));
    expect(onRemoveItem).toHaveBeenCalledTimes(1);
  });

  it('returns an empty heading when title is nullish', () => {
    render(
      <SimpleAccordion heading={null} childProps={{ hasRemove: false }}>
        <div>child</div>
      </SimpleAccordion>,
    );
    // No throw, just empty heading
    expect(screen.getByTestId('accordion')).toBeInTheDocument();
  });
});
