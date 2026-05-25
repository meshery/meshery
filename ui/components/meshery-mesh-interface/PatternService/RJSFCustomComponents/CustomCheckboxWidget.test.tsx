import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  FormControlLabel: ({ control, label }: any) => (
    <label data-testid="form-control-label">
      {control}
      <span>{label}</span>
    </label>
  ),
  IconButton: ({ children }: any) => <button>{children}</button>,
  Checkbox: ({ checked, onChange, id, disabled, required }: any) => (
    <input
      data-testid="checkbox"
      type="checkbox"
      id={id}
      checked={!!checked}
      disabled={!!disabled}
      required={!!required}
      onChange={onChange}
    />
  ),
  useTheme: () => ({ palette: { mode: 'light' } }),
}));

vi.mock('@rjsf/utils', () => ({
  schemaRequiresTrueValue: (schema: any) => schema?.required === true,
  labelValue: (label: any, hideLabel: any) => (hideLabel ? '' : label),
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

vi.mock('../../../../css/icons.styles', () => ({ iconSmall: {} }));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
}));

import { CustomCheckboxWidget } from './CustomCheckboxWidget';

describe('CustomCheckboxWidget', () => {
  it('renders a checkbox reflecting the value', () => {
    render(
      <CustomCheckboxWidget id="c1" value={true} schema={{}} label="Enable" onChange={vi.fn()} />,
    );
    expect(screen.getByTestId('checkbox')).toBeChecked();
  });

  it('treats undefined value as unchecked', () => {
    render(<CustomCheckboxWidget id="c1" schema={{}} label="Enable" onChange={vi.fn()} />);
    expect(screen.getByTestId('checkbox')).not.toBeChecked();
  });

  it('calls onChange with checked state', () => {
    const onChange = vi.fn();
    render(
      <CustomCheckboxWidget id="c1" value={false} schema={{}} label="Enable" onChange={onChange} />,
    );
    fireEvent.click(screen.getByTestId('checkbox'));
    expect(onChange).toHaveBeenCalledWith(true);
  });

  it('renders a help tooltip when schema.description is present', () => {
    render(
      <CustomCheckboxWidget
        id="c1"
        schema={{ description: 'helpful text' }}
        label="Enable"
        onChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('tooltip')).toHaveAttribute('data-title', 'helpful text');
  });
});
