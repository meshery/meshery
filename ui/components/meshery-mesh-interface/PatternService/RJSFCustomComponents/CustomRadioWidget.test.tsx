import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  FormControlLabel: ({ control, label, value, disabled }: any) => (
    <label data-testid="form-control-label" data-value={value} data-disabled={String(!!disabled)}>
      {control}
      <span>{label}</span>
    </label>
  ),
  FormLabel: ({ children, required, htmlFor }: any) => (
    <span data-testid="form-label" data-required={String(!!required)} data-html-for={htmlFor}>
      {children}
    </span>
  ),
  Radio: ({ name, id }: any) => (
    <input data-testid={`radio-${id}`} type="radio" name={name} id={id} />
  ),
  RadioGroup: ({ children, value, onChange, name }: any) => (
    <div data-testid="radio-group" data-name={name} data-value={String(value)}>
      <button
        type="button"
        data-testid="trigger-onchange"
        onClick={() => onChange?.({ target: { value: '1' } }, '1')}
      >
        change
      </button>
      {children}
    </div>
  ),
}));

vi.mock('@rjsf/utils', () => ({
  ariaDescribedByIds: (id: string) => `desc-${id}`,
  enumOptionsIndexForValue: (val: any, enumOpts: any[]) =>
    enumOpts.findIndex((o) => o.value === val)?.toString?.() ?? null,
  enumOptionsValueForIndex: (index: any, enumOpts: any[]) => enumOpts[Number(index)]?.value,
  optionId: (id: string, idx: number) => `${id}-option-${idx}`,
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
}));

import CustomRadioWidget from './CustomRadioWidget';

describe('CustomRadioWidget', () => {
  const enumOptions = [
    { value: 'a', label: 'Option A' },
    { value: 'b', label: 'Option B' },
  ];

  it('renders enumOptions as radios with labels', () => {
    render(
      <CustomRadioWidget
        id="r1"
        label="Pick"
        value="a"
        options={{ enumOptions }}
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    expect(screen.getByText('Option A')).toBeInTheDocument();
    expect(screen.getByText('Option B')).toBeInTheDocument();
  });

  it('hides the form label when hideLabel is true', () => {
    render(
      <CustomRadioWidget
        id="r1"
        label="Pick"
        hideLabel
        value="a"
        options={{ enumOptions }}
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    expect(screen.queryByTestId('form-label')).not.toBeInTheDocument();
  });

  it('disables individual options when enumDisabled is set', () => {
    render(
      <CustomRadioWidget
        id="r1"
        label="Pick"
        value="a"
        options={{ enumOptions, enumDisabled: ['b'] }}
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    const labels = screen.getAllByTestId('form-control-label');
    // Option B should be marked disabled
    expect(labels[1]).toHaveAttribute('data-disabled', 'true');
  });

  it('converts the radio change to the underlying enum value', () => {
    const onChange = vi.fn();
    render(
      <CustomRadioWidget
        id="r1"
        label="Pick"
        value="a"
        options={{ enumOptions }}
        onChange={onChange}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    fireEvent.click(screen.getByTestId('trigger-onchange'));
    expect(onChange).toHaveBeenCalledWith('b');
  });

  it('coerces non-string non-React label values to strings', () => {
    render(
      <CustomRadioWidget
        id="r1"
        label={{ complex: 'label' } as any}
        value="a"
        options={{ enumOptions }}
        onChange={vi.fn()}
        onBlur={vi.fn()}
        onFocus={vi.fn()}
      />,
    );
    expect(screen.getByTestId('form-label')).toHaveTextContent('[object Object]');
  });
});
