import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

let mockMode: 'light' | 'dark' = 'light';

vi.mock('@/theme', () => ({
  useTheme: () => ({ palette: { mode: mockMode } }),
  SistentThemeProviderWithoutBaseLine: ({ children, initialMode }: any) => (
    <div data-testid="sistent-provider" data-initial-mode={initialMode}>
      {children}
    </div>
  ),
}));

import { UsesSistent } from './SistentWrapper';

describe('UsesSistent (SistentWrapper)', () => {
  it('wraps children with the Sistent theme provider', () => {
    mockMode = 'light';
    render(
      <UsesSistent>
        <div data-testid="content">hello</div>
      </UsesSistent>,
    );

    expect(screen.getByTestId('sistent-provider')).toBeInTheDocument();
    expect(screen.getByTestId('content')).toBeInTheDocument();
  });

  it('forwards the theme mode to the provider when light', () => {
    mockMode = 'light';
    render(
      <UsesSistent>
        <span>child</span>
      </UsesSistent>,
    );

    expect(screen.getByTestId('sistent-provider')).toHaveAttribute('data-initial-mode', 'light');
  });

  it('forwards the theme mode to the provider when dark', () => {
    mockMode = 'dark';
    render(
      <UsesSistent>
        <span>child</span>
      </UsesSistent>,
    );

    expect(screen.getByTestId('sistent-provider')).toHaveAttribute('data-initial-mode', 'dark');
  });

  it('renders multiple children', () => {
    mockMode = 'light';
    render(
      <UsesSistent>
        <span data-testid="a">A</span>
        <span data-testid="b">B</span>
      </UsesSistent>,
    );

    expect(screen.getByTestId('a')).toBeInTheDocument();
    expect(screen.getByTestId('b')).toBeInTheDocument();
  });
});
