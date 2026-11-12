import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ title, children }: any) => (
    <div data-testid="custom-tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

import { CustomTextTooltip } from './CustomTextTooltip';

describe('CustomTextTooltip', () => {
  it('forwards props to the underlying tooltip', () => {
    render(
      <CustomTextTooltip title="hello" interactive>
        <span>child</span>
      </CustomTextTooltip>,
    );

    const tooltip = screen.getByTestId('custom-tooltip');
    expect(tooltip).toHaveAttribute('data-title', 'hello');
    expect(screen.getByText('child')).toBeInTheDocument();
  });
});
