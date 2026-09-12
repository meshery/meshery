import React from 'react';
import { render, screen } from '@testing-library/react';
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

  it('renders a delete button when hasRemove is true and invokes onDropIndexClick', () => {
    const onDropIndexClick = vi.fn(() => vi.fn());
    render(
      <SimpleAccordion heading="thing" childProps={{ hasRemove: true, index: 2, onDropIndexClick }}>
        <div>child</div>
      </SimpleAccordion>,
    );
    expect(screen.getByTestId('delete-icon')).toBeInTheDocument();
    expect(onDropIndexClick).toHaveBeenCalledWith(2);
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
