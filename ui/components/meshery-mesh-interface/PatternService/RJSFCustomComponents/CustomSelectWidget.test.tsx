import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  IconButton: ({ children }: any) => <button>{children}</button>,
  InputAdornment: ({ children }: any) => <div data-testid="adornment">{children}</div>,
  ListItemText: ({ primary }: any) => <span>{primary}</span>,
  MenuItem: ({ children, value, disabled }: any) => (
    <li data-testid="menu-item" data-value={value} data-disabled={String(!!disabled)}>
      {children}
    </li>
  ),
  TextField: ({
    value,
    onChange,
    label,
    id,
    disabled,
    error,
    InputProps,
    SelectProps,
    children,
  }: any) => (
    <div data-testid="textfield-wrapper" data-id={id}>
      <input
        data-testid="textfield"
        data-error={String(!!error)}
        data-label={String(label)}
        value={value ?? ''}
        disabled={!!disabled}
        onChange={(e) => onChange(e)}
      />
      {InputProps?.endAdornment}
      <div data-testid="select-children">{children}</div>
    </div>
  ),
  InputLabel: ({ children, htmlFor, required }: any) => (
    <label htmlFor={htmlFor} data-required={String(!!required)}>
      {children}
    </label>
  ),
  Checkbox: ({ checked }: any) => (
    <input type="checkbox" checked={!!checked} readOnly data-testid="checkbox" />
  ),
  useTheme: () => ({ palette: { error: { main: '#f00' }, mode: 'light' } }),
}));

vi.mock('../../../../assets/icons/HelpOutlineIcon', () => ({
  default: () => <svg data-testid="help-icon" />,
}));

vi.mock('../../../../assets/icons/ErrorOutlineIcon', () => ({
  default: () => <svg data-testid="error-icon" />,
}));

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {} }));

vi.mock('../CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('@rjsf/utils', () => ({
  ariaDescribedByIds: (id: string) => `desc-${id}`,
  enumOptionsIndexForValue: (val: any, enumOpts: any[]) => {
    const idx = enumOpts.findIndex((o) => o.value === val);
    return idx === -1 ? null : String(idx);
  },
  enumOptionsValueForIndex: (idx: any, enumOpts: any[]) => enumOpts[Number(idx)]?.value,
  labelValue: (label: any, hideLabel: any) => (hideLabel ? '' : label),
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
}));

import CustomSelectWidget from './CustomSelectWidget';

describe('CustomSelectWidget', () => {
  const enumOptions = [
    { value: 'a', label: 'Alpha' },
    { value: 'b', label: 'Beta' },
  ];

  it('renders MenuItems for every enum option', () => {
    render(
      <CustomSelectWidget
        id="s1"
        label="Pick"
        options={{ enumOptions }}
        schema={{}}
        value="a"
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    const items = screen.getAllByTestId('menu-item');
    expect(items).toHaveLength(2);
  });

  it('emits the underlying enum value on change', () => {
    const onChange = vi.fn();
    render(
      <CustomSelectWidget
        id="s1"
        label="Pick"
        options={{ enumOptions }}
        schema={{}}
        value="a"
        onChange={onChange}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    fireEvent.change(screen.getByTestId('textfield'), { target: { value: '1' } });
    expect(onChange).toHaveBeenCalledWith('b');
  });

  it('renders an error indicator when rawErrors are present', () => {
    render(
      <CustomSelectWidget
        id="s1"
        label="Pick"
        options={{ enumOptions }}
        schema={{}}
        rawErrors={['required']}
        value="a"
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    expect(screen.getByTestId('textfield')).toHaveAttribute('data-error', 'true');
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });

  it('shows a help tooltip when schema.description is provided', () => {
    render(
      <CustomSelectWidget
        id="s1"
        label="Pick"
        options={{ enumOptions }}
        schema={{ description: 'pick one' }}
        value="a"
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    expect(screen.getByTestId('help-icon')).toBeInTheDocument();
  });

  it('renders an external InputLabel when x-rjsf-grid-area is set', () => {
    render(
      <CustomSelectWidget
        id="s1"
        label="Pick"
        options={{ enumOptions }}
        schema={{ 'x-rjsf-grid-area': 6 }}
        value="a"
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    const label = document.querySelector('label[for="s1"]');
    expect(label).not.toBeNull();
  });
});
