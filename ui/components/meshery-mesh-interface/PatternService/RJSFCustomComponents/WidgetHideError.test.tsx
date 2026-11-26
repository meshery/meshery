import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import Form from '@rjsf/core';
import validator from '@rjsf/validator-ajv8';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledMock';
    return StyledComponent;
  };
  return {
    IconButton: ({ children, ...props }: any) => <button {...props}>{children}</button>,
    InputAdornment: ({ children }: any) => <div>{children}</div>,
    TextField: ({ value, onChange, error, id, disabled, type, InputProps, children }: any) => (
      <div data-testid="textfield-wrapper">
        <input
          data-testid="textfield"
          data-error={String(!!error)}
          type={type === 'file' ? 'text' : type}
          id={id}
          value={value == null ? '' : value}
          disabled={!!disabled}
          onChange={(e) => onChange(e)}
        />
        {InputProps?.endAdornment}
        {children}
      </div>
    ),
    InputLabel: ({ children, htmlFor, required }: any) => (
      <label htmlFor={htmlFor} data-required={String(!!required)}>
        {children}
      </label>
    ),
    ListItemText: ({ primary }: any) => <span>{primary}</span>,
    MenuItem: ({ children, value, disabled }: any) => (
      <option value={value} disabled={disabled}>
        {children}
      </option>
    ),
    Checkbox: ({ checked }: any) => <input type="checkbox" checked={!!checked} readOnly />,
    Grid2: ({ children }: any) => <div>{children}</div>,
    Box: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
    useTheme: () => ({
      palette: { error: { main: '#ff0000' }, mode: 'light' },
    }),
    styled,
  };
});

vi.mock('../../../../assets/icons/HelpOutlineIcon', () => ({
  default: () => <svg data-testid="help-icon" />,
}));

vi.mock('../../../../assets/icons/ErrorOutlineIcon', () => ({
  default: () => <svg data-testid="error-icon" />,
}));

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {}, iconMedium: {} }));

vi.mock('../CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
  safeStringTitle: (v: any) => (v == null ? '' : String(v)),
  calculateGrid: () => ({ xs: 12 }),
}));

import CustomSelectWidget from './CustomSelectWidget';
import CustomTextWidget from './CustomTextWidget';

// End-to-end coverage for #23108: RJSF's StringField passes hideError to the
// widget alongside rawErrors, so a real Form with hideError set must not show
// the widgets' error icon. Same defect class as #22954 (object template),
// one layer down in the widgets. Direct-render coverage lives in
// CustomBaseInput.test.tsx and CustomSelectWidget.test.tsx.
describe('Widget hideError in a real RJSF Form', () => {
  it('text field hides the error icon when hideError is true', () => {
    const schema: any = {
      type: 'object',
      properties: { name: { type: 'string', title: 'Name', minLength: 5 } },
      required: ['name'],
    };
    const { container } = render(
      <Form
        schema={schema}
        uiSchema={{ name: { 'ui:options': { hideError: true } } }}
        formData={{ name: 'ab' }}
        validator={validator}
        widgets={{ TextWidget: CustomTextWidget }}
      />,
    );
    fireEvent.click(container.querySelector('button[type="submit"]')!);
    expect(screen.queryByTestId('error-icon')).not.toBeInTheDocument();
  });

  it('enum field hides the error icon when hideError is true', () => {
    const schema: any = {
      type: 'object',
      properties: {
        color: { type: 'string', title: 'Color', enum: ['red', 'blue'] },
      },
      required: ['color'],
    };
    const { container } = render(
      <Form
        schema={schema}
        uiSchema={{ color: { 'ui:options': { hideError: true } } }}
        formData={{}}
        validator={validator}
        widgets={{ SelectWidget: CustomSelectWidget }}
      />,
    );
    fireEvent.click(container.querySelector('button[type="submit"]')!);
    expect(screen.queryByTestId('error-icon')).not.toBeInTheDocument();
  });

  it('text field still shows the error icon without hideError', () => {
    const schema: any = {
      type: 'object',
      properties: { name: { type: 'string', title: 'Name', minLength: 5 } },
      required: ['name'],
    };
    const { container } = render(
      <Form
        schema={schema}
        uiSchema={{}}
        formData={{ name: 'ab' }}
        validator={validator}
        widgets={{ TextWidget: CustomTextWidget }}
      />,
    );
    fireEvent.click(container.querySelector('button[type="submit"]')!);
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });

  it('enum field still shows the error icon without hideError', () => {
    const schema: any = {
      type: 'object',
      properties: {
        color: { type: 'string', title: 'Color', enum: ['red', 'blue'] },
      },
      required: ['color'],
    };
    const { container } = render(
      <Form
        schema={schema}
        uiSchema={{}}
        formData={{}}
        validator={validator}
        widgets={{ SelectWidget: CustomSelectWidget }}
      />,
    );
    fireEvent.click(container.querySelector('button[type="submit"]')!);
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });
});
