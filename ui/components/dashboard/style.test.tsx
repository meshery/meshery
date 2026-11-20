import React from 'react';
import { render, screen } from '@testing-library/react';
import { SistentThemeProvider } from '@sistent/sistent';
import { describe, expect, it } from 'vitest';
import { DashboardActionsContainer } from './style';

describe('DashboardActionsContainer', () => {
  it('renders with responsive wrapping and full width styles in the DOM', () => {
    render(
      <SistentThemeProvider initialMode="light">
        <DashboardActionsContainer data-testid="dashboard-actions-container">
          <div>Action 1</div>
          <div>Action 2</div>
        </DashboardActionsContainer>
      </SistentThemeProvider>,
    );

    const container = screen.getByTestId('dashboard-actions-container');
    expect(container).toBeInTheDocument();

    const computedStyle = window.getComputedStyle(container);
    expect(computedStyle.display).toBe('flex');
    expect(computedStyle.flexWrap).toBe('wrap');
    expect(computedStyle.width).toBe('100%');
    expect(computedStyle.maxWidth).toBe('100%');
    expect(computedStyle.justifyContent).toBe('flex-end');
  });
});
