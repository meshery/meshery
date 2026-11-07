import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Box: ({ children }: any) => <div>{children}</div>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  IconButton: ({ children, onClick, className, disabled }: any) => (
    <button data-testid={className || 'icon-button'} onClick={onClick} disabled={!!disabled}>
      {children}
    </button>
  ),
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({ palette: { mode: 'light', error: { main: '#f00' } } }),
  CssBaseline: () => null,
  ExpandMoreIcon: () => <svg data-testid="expand-more" />,
}));

vi.mock('@rjsf/utils', () => ({
  canExpand: (schema: any) => !!schema?.canExpand,
  // Faithful port of RJSF's real getUiOptions: the registry's global UI
  // options are the base, then every ui:* key is applied in uiSchema key
  // order, with ui:options spread in and bare ui:* keys losing their prefix.
  // (The real one also logs when ui:widget is an object; that path is
  // skipped here because no test sets it.)
  getUiOptions: (uiSchema: any = {}, globalOptions: any = {}) => {
    if (!uiSchema) return { ...globalOptions };
    return Object.keys(uiSchema)
      .filter((key: string) => key.startsWith('ui:'))
      .reduce(
        (options: any, key: string) => {
          const value = uiSchema[key];
          if (key === 'ui:widget' && typeof value === 'object' && value !== null) {
            return options;
          }
          if (key === 'ui:options' && typeof value === 'object' && value !== null) {
            return { ...options, ...value };
          }
          return { ...options, [key.substring(3)]: value };
        },
        { ...globalOptions },
      );
  },
}));

vi.mock('../../../../assets/icons/AddIcon', () => ({
  default: () => <svg data-testid="add-icon" />,
}));

vi.mock('../../../../assets/icons/ExpandLessIcon', () => ({
  default: () => <svg data-testid="expand-less" />,
}));

vi.mock('../../../../assets/icons/ErrorOutlineIcon', () => ({
  default: () => <svg data-testid="error-icon" />,
}));

vi.mock('../../../../assets/icons/HelpOutlineIcon', () => ({
  default: () => <svg data-testid="help-icon" />,
}));

vi.mock('../CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {}, iconMedium: {} }));

vi.mock('../helper', () => ({
  calculateGrid: () => ({ xs: 12, md: 12, lg: 6 }),
  safeStringTitle: (v: any) => (v == null ? '' : String(v)),
}));

import ObjectFieldTemplate from './ObjectFieldTemplate';

describe('ObjectFieldTemplate', () => {
  const baseProps = {
    description: 'an object',
    title: 'Settings',
    properties: [{ name: 'name', content: <div data-testid="prop-content">name field</div> }],
    disabled: false,
    readonly: false,
    uiSchema: {},
    idSchema: { $id: 'root' },
    schema: { properties: { name: { type: 'string' } } },
    formData: {},
    errorSchema: {},
  };

  it('renders the title and a toggle button when properties exist', () => {
    render(<ObjectFieldTemplate {...(baseProps as any)} />);
    expect(screen.getByText('Settings')).toBeInTheDocument();
    expect(screen.getByTestId('object-property-expand')).toBeInTheDocument();
  });

  it('initially keeps properties collapsed and expands them on toggle', () => {
    render(<ObjectFieldTemplate {...(baseProps as any)} />);
    expect(screen.queryByTestId('prop-content')).not.toBeInTheDocument();
    fireEvent.click(screen.getByTestId('object-property-expand'));
    expect(screen.getByTestId('prop-content')).toBeInTheDocument();
  });

  it('renders an Add button when canExpand is true', () => {
    const onAddProperty = vi.fn(() => vi.fn());
    render(
      <ObjectFieldTemplate
        {...(baseProps as any)}
        schema={{ canExpand: true, properties: { name: { type: 'string' } } }}
        onAddProperty={onAddProperty}
      />,
    );
    expect(screen.getByTestId('object-property-expand')).toBeInTheDocument();
    expect(screen.getByTestId('add-icon')).toBeInTheDocument();
  });

  it('renders an error tooltip when errorSchema contains errors', () => {
    render(
      <ObjectFieldTemplate {...(baseProps as any)} errorSchema={{ name: { __errors: ['bad'] } }} />,
    );
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });

  it('hides the error icon when ui:options.hideError is true', () => {
    // ui:options.hideError must hide the indicator even though RJSF still
    // passes errorSchema to the template: the caller asked for no error to
    // be shown. RJSF's ObjectField does not forward hideError to the
    // template, so the template reads it from the effective UI options.
    render(
      <ObjectFieldTemplate
        {...(baseProps as any)}
        errorSchema={{ name: { __errors: ['bad'] } }}
        uiSchema={{ 'ui:options': { hideError: true } }}
      />,
    );
    expect(screen.queryByTestId('error-icon')).not.toBeInTheDocument();
  });

  it('hides the error icon when the bare ui:hideError option is set', () => {
    render(
      <ObjectFieldTemplate
        {...(baseProps as any)}
        errorSchema={{ name: { __errors: ['bad'] } }}
        uiSchema={{ 'ui:hideError': true }}
      />,
    );
    expect(screen.queryByTestId('error-icon')).not.toBeInTheDocument();
  });

  it('hides the error icon when hideError comes from the registry global UI options', () => {
    render(
      <ObjectFieldTemplate
        {...(baseProps as any)}
        errorSchema={{ name: { __errors: ['bad'] } }}
        registry={{ globalUiOptions: { hideError: true } }}
      />,
    );
    expect(screen.queryByTestId('error-icon')).not.toBeInTheDocument();
  });

  it('renders only the properties grid when no title is provided', () => {
    render(<ObjectFieldTemplate {...(baseProps as any)} title="" uiSchema={{}} />);
    // No fieldTitle means CustomTitleField is skipped — content renders directly
    expect(screen.getByTestId('prop-content')).toBeInTheDocument();
  });

  it('uses ui:options.expand to initialize expansion', () => {
    render(
      <ObjectFieldTemplate {...(baseProps as any)} uiSchema={{ 'ui:options': { expand: true } }} />,
    );
    expect(screen.getByTestId('prop-content')).toBeInTheDocument();
  });

  it('uses "Value" as title for additional properties', () => {
    render(
      <ObjectFieldTemplate
        {...(baseProps as any)}
        schema={{ ...baseProps.schema, __additional_property: true }}
      />,
    );
    expect(screen.getByText('Value')).toBeInTheDocument();
  });
});
