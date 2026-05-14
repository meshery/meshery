import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

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
    TextField: ({
      value,
      onChange,
      error,
      label,
      id,
      disabled,
      type,
      multiline,
      InputProps,
    }: any) => (
      <div data-testid="textfield-wrapper">
        <input
          data-testid="textfield"
          data-error={String(!!error)}
          data-label={String(label)}
          data-multiline={String(!!multiline)}
          type={type === 'file' ? 'text' : type}
          id={id}
          value={value == null ? '' : value}
          disabled={!!disabled}
          onChange={(e) => onChange(e)}
        />
        {InputProps?.endAdornment}
      </div>
    ),
    InputLabel: ({ children, htmlFor, required }: any) => (
      <label htmlFor={htmlFor} data-required={String(!!required)}>
        {children}
      </label>
    ),
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

vi.mock('../../../../css/icons.styles', () => ({
  iconSmall: {},
  iconMedium: {},
}));

vi.mock('../CustomTextTooltip', () => ({
  CustomTextTooltip: ({ children, title }: any) => (
    <div data-testid="tooltip" data-title={String(title)}>
      {children}
    </div>
  ),
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
}));

import BaseInput from './CustomBaseInput';

describe('CustomBaseInput', () => {
  it('renders the textfield with the supplied label and value', () => {
    render(
      <BaseInput
        id="input1"
        label="Name"
        value="alice"
        schema={{}}
        options={{}}
        onChange={vi.fn()}
      />,
    );
    const tf = screen.getByTestId('textfield');
    expect(tf).toHaveAttribute('data-label', 'Name');
    expect(tf).toHaveValue('alice');
  });

  it('emits empty string when the user clears the input', () => {
    const onChange = vi.fn();
    render(
      <BaseInput
        id="input1"
        label="Name"
        value="alice"
        schema={{}}
        options={{ emptyValue: '__EMPTY__' }}
        onChange={onChange}
      />,
    );
    const tf = screen.getByTestId('textfield');
    fireEvent.change(tf, { target: { value: '' } });
    expect(onChange).toHaveBeenCalledWith('__EMPTY__');
  });

  it('encodes values in URI when x-encode-in-uri is enabled', () => {
    const onChange = vi.fn();
    render(
      <BaseInput
        id="input1"
        label="Name"
        value=""
        schema={{ 'x-encode-in-uri': true }}
        options={{}}
        onChange={onChange}
      />,
    );
    const tf = screen.getByTestId('textfield');
    fireEvent.change(tf, { target: { value: 'foo bar' } });
    expect(onChange).toHaveBeenCalledWith(encodeURIComponent('foo bar'));
  });

  it('shows error tooltip when rawErrors are present', () => {
    render(
      <BaseInput
        id="input1"
        label="Name"
        value="alice"
        rawErrors={['required']}
        schema={{}}
        options={{}}
        onChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('textfield')).toHaveAttribute('data-error', 'true');
    expect(screen.getByTestId('error-icon')).toBeInTheDocument();
  });

  it('passes through file input changes raw without re-encoding', () => {
    const onChange = vi.fn();
    render(
      <BaseInput
        id="input1"
        label="File"
        value="data"
        schema={{}}
        options={{ inputType: 'file' }}
        onChange={onChange}
      />,
    );
    const tf = screen.getByTestId('textfield');
    fireEvent.change(tf, { target: { value: 'whatever' } });
    expect(onChange).toHaveBeenCalled();
  });

  it('renders an InputLabel when x-rjsf-grid-area is set', () => {
    render(
      <BaseInput
        id="input1"
        label="Name"
        value="x"
        schema={{ 'x-rjsf-grid-area': 6 }}
        options={{}}
        onChange={vi.fn()}
      />,
    );
    const label = document.querySelector('label[for="input1"]');
    expect(label).not.toBeNull();
  });

  it('uses "Value" as the label when the field is an additional property', () => {
    render(
      <BaseInput
        id="input1"
        label="newKey"
        value="New Value"
        schema={{ __additional_property: true }}
        options={{}}
        onChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('textfield')).toHaveAttribute('data-label', 'Value');
    // additional property "New Value" placeholder is replaced with empty string
    expect(screen.getByTestId('textfield')).toHaveValue('');
  });

  it('disables the input when disabled or readonly', () => {
    render(
      <BaseInput
        id="input1"
        label="Name"
        value="x"
        disabled
        schema={{}}
        options={{}}
        onChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('textfield')).toBeDisabled();
  });
});
