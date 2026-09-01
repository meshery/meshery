import React from 'react';
import { render, screen, fireEvent, within } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  IconButton: ({ children }: any) => <button>{children}</button>,
  InputAdornment: ({ children }: any) => <div data-testid="adornment">{children}</div>,
  ListItemText: ({ primary, primaryTypographyProps, style }: any) => (
    <span
      data-testid="list-item-text"
      data-nowrap={String(primaryTypographyProps?.noWrap)}
      style={style}
    >
      {primary}
    </span>
  ),
  MenuItem: ({ children, value, disabled, sx }: any) => (
    <li data-testid="menu-item" data-value={value} data-disabled={String(!!disabled)} style={sx}>
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
    slotProps,
    children,
  }: any) => {
    const endAdornment = slotProps?.input?.endAdornment || InputProps?.endAdornment;
    const selectConfig = slotProps?.select || SelectProps;
    const renderedDisplay = selectConfig?.renderValue ? selectConfig.renderValue(value) : null;
    return (
      <div
        data-testid="textfield-wrapper"
        data-id={id}
        data-multiple={String(selectConfig?.multiple)}
        data-rendered-value={typeof renderedDisplay === 'string' ? renderedDisplay : undefined}
      >
        <div data-testid="rendered-display">{renderedDisplay}</div>
        <input
          data-testid="textfield"
          data-error={String(!!error)}
          data-label={String(label)}
          value={value ?? ''}
          disabled={!!disabled}
          onChange={(e) => onChange(e)}
        />
        {endAdornment}
        <div data-testid="select-children">{children}</div>
      </div>
    );
  },
  InputLabel: ({ children, htmlFor, required }: any) => (
    <label htmlFor={htmlFor} data-required={String(!!required)}>
      {children}
    </label>
  ),
  Checkbox: ({ checked, sx, style }: any) => (
    <input
      type="checkbox"
      checked={!!checked}
      readOnly
      data-testid="checkbox"
      style={sx ?? style}
    />
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
  enumOptionsIndexForValue: (val: any, enumOpts: any[], multiple?: boolean) => {
    if (multiple) {
      if (!Array.isArray(val)) return [];
      return val
        .map((v) => {
          const idx = enumOpts.findIndex((o) => o.value === v);
          return idx === -1 ? null : String(idx);
        })
        .filter((i): i is string => i !== null);
    }
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

  it('renders checkboxes and tooltips when multiple is true', () => {
    render(
      <CustomSelectWidget
        id="s1"
        label="Technology"
        options={{ enumOptions, multiple: true }}
        schema={{ type: 'array' }}
        value={['a']}
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    const checkboxes = screen.getAllByTestId('checkbox');
    expect(checkboxes).toHaveLength(2);
    expect(checkboxes[0]).toBeChecked();
    expect(checkboxes[1]).not.toBeChecked();

    const tooltips = screen.getAllByTestId('tooltip');
    expect(tooltips).toHaveLength(2);
    expect(tooltips[0]).toHaveAttribute('data-title', 'Alpha');
    expect(tooltips[1]).toHaveAttribute('data-title', 'Beta');

    const textItems = screen.getAllByTestId('list-item-text');
    expect(textItems).toHaveLength(2);
    expect(textItems[0]).toHaveAttribute('data-nowrap', 'true');

    const wrapper = screen.getByTestId('textfield-wrapper');
    expect(wrapper).toHaveAttribute('data-multiple', 'true');
    expect(screen.getByTestId('rendered-display')).toHaveTextContent('Alpha');
  });

  it('renders without checkboxes for single-select fields', () => {
    render(
      <CustomSelectWidget
        id="s1"
        label="Type"
        options={{ enumOptions, multiple: false }}
        schema={{ type: 'string' }}
        value="a"
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    expect(screen.queryByTestId('checkbox')).toBeNull();
    const tooltips = screen.getAllByTestId('tooltip');
    expect(tooltips).toHaveLength(2);

    const wrapper = screen.getByTestId('textfield-wrapper');
    expect(wrapper).toHaveAttribute('data-multiple', 'false');
    expect(screen.getByTestId('rendered-display')).toHaveTextContent('Alpha');
  });

  it('preserves React-element labels when rendering multiple selected options', () => {
    const customEnumOptions = [
      { value: 'k8s', label: <span data-testid="k8s-badge">Kubernetes</span> },
      { value: 'istio', label: <span data-testid="istio-badge">Istio</span> },
    ];
    render(
      <CustomSelectWidget
        id="s1"
        label="Technology"
        options={{ enumOptions: customEnumOptions, multiple: true }}
        schema={{ type: 'array' }}
        value={['k8s', 'istio']}
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    const display = screen.getByTestId('rendered-display');
    expect(within(display).getByTestId('k8s-badge')).toBeInTheDocument();
    expect(within(display).getByTestId('istio-badge')).toBeInTheDocument();
    expect(display).toHaveTextContent('Kubernetes, Istio');
  });
});
