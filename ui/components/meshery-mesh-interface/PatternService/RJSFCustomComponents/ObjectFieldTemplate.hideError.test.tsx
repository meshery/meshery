import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import Form from '@rjsf/core';
import validator from '@rjsf/validator-ajv8';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div>{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  CssBaseline: () => null,
  IconButton: ({ children, onClick, disabled }: any) => (
    <button onClick={onClick} disabled={!!disabled}>
      {children}
    </button>
  ),
  Typography: ({ children }: any) => <span>{children}</span>,
  ExpandMoreIcon: () => <span data-testid="expand-more-icon" />,
  ExpandLessIcon: () => <span data-testid="expand-less-icon" />,
  useTheme: () => ({
    palette: { mode: 'light', error: { main: '#f00' } },
    typography: { pxToRem: (n: number) => `${n}px`, fontWeightRegular: 400 },
  }),
}));

vi.mock('../../../../assets/icons/AddIcon', () => ({
  default: () => <svg data-testid="add-icon" />,
}));

vi.mock('../CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('../../../../assets/icons/HelpOutlineIcon', () => ({
  default: () => <svg data-testid="help-icon" />,
}));

vi.mock('../../../../assets/icons/ErrorOutlineIcon', () => ({
  default: () => <svg data-testid="error-icon" />,
}));

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {}, iconMedium: {} }));

vi.mock('../helper', () => ({
  safeStringTitle: (v: any) => (v == null ? '' : String(v)),
  calculateGrid: () => ({ xs: 12 }),
}));

import ObjectFieldTemplate from './ObjectFieldTemplate';

// End-to-end check through a real RJSF Form: RJSF still passes errorSchema
// to the object template when hideError is set, so the template itself must
// hide the error icon. This is the same defect class as the array template's
// rawErrors block before the #22863 hideError fix.
describe('ObjectFieldTemplate hideError through a real RJSF Form', () => {
  const schema: any = {
    type: 'object',
    properties: {
      settings: {
        type: 'object',
        title: 'Settings',
        properties: { name: { type: 'string', title: 'Name' } },
        required: ['name'],
      },
    },
  };

  it('shows the error icon at the object title on a required-field violation', () => {
    const { container } = render(
      <Form
        schema={schema}
        formData={{ settings: {} }}
        validator={validator}
        templates={{ ObjectFieldTemplate }}
      />,
    );
    fireEvent.click(container.querySelector('button[type="submit"]')!);
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });

  it('hides the error icon when hideError is true', () => {
    const { container } = render(
      <Form
        schema={schema}
        uiSchema={{ settings: { 'ui:options': { hideError: true } } }}
        formData={{ settings: {} }}
        validator={validator}
        templates={{ ObjectFieldTemplate }}
      />,
    );
    fireEvent.click(container.querySelector('button[type="submit"]')!);
    expect(screen.queryByTestId('error-icon')).not.toBeInTheDocument();
  });
});
