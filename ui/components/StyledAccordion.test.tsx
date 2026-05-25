import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  styled: () => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <div data-styled="true" {...props}>
        {children}
      </div>
    );
    StyledComponent.displayName = 'StyledMock';
    return StyledComponent;
  },
  Accordion: ({ children, ...props }: any) => (
    <div data-testid="accordion" {...props}>
      {children}
    </div>
  ),
  AccordionSummary: ({ children, ...props }: any) => (
    <div data-testid="accordion-summary" {...props}>
      {children}
    </div>
  ),
}));

import { StyledAccordion, StyledAccordionSummary } from './StyledAccordion';

describe('StyledAccordion', () => {
  it('exports StyledAccordion as a React component', () => {
    expect(StyledAccordion).toBeDefined();
    expect(typeof StyledAccordion).toBe('function');
  });

  it('renders children', () => {
    render(
      <StyledAccordion>
        <span data-testid="child">inner</span>
      </StyledAccordion>,
    );

    expect(screen.getByTestId('child')).toBeInTheDocument();
  });
});

describe('StyledAccordionSummary', () => {
  it('exports StyledAccordionSummary as a React component', () => {
    expect(StyledAccordionSummary).toBeDefined();
    expect(typeof StyledAccordionSummary).toBe('function');
  });

  it('renders children', () => {
    render(
      <StyledAccordionSummary>
        <span data-testid="summary-child">Summary</span>
      </StyledAccordionSummary>,
    );

    expect(screen.getByTestId('summary-child')).toBeInTheDocument();
  });

  it('renders both components together', () => {
    render(
      <StyledAccordion>
        <StyledAccordionSummary>
          <span data-testid="summary-text">Click me</span>
        </StyledAccordionSummary>
        <span data-testid="body-text">Body</span>
      </StyledAccordion>,
    );

    expect(screen.getByTestId('summary-text')).toBeInTheDocument();
    expect(screen.getByTestId('body-text')).toBeInTheDocument();
  });
});
