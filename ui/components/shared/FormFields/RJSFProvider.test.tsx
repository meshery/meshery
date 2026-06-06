import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@rjsf/core', () => ({
  withTheme: (theme: any) => () => (
    <div data-testid="rjsf-themed-form" data-theme={theme?.name || 'material'} />
  ),
}));

vi.mock('@rjsf/mui', () => ({
  Theme: { name: 'material' },
}));

import RJSFProvider from './RJSFProvider';

describe('RJSFProvider', () => {
  it('exports a themed RJSF form factory', () => {
    expect(RJSFProvider).toBeDefined();
  });

  it('renders the themed form', () => {
    render(<RJSFProvider />);
    expect(screen.getByTestId('rjsf-themed-form')).toBeInTheDocument();
  });
});
