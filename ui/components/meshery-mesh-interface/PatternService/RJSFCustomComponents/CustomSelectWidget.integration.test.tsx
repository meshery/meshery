import React, { useState } from 'react';
import { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { SistentThemeProvider } from '@sistent/sistent';
import CustomSelectWidget from './CustomSelectWidget';

const enumOptions = [
  { value: 'k8s', label: 'Kubernetes' },
  { value: 'istio', label: 'Istio' },
  { value: 'long', label: 'Extremely long synthetic technology integration label' },
];

function setup(overrides = {}) {
  const onChange = vi.fn();
  const props = {
    id: 'technology',
    label: 'Technology',
    schema: { type: 'array' },
    options: { enumOptions },
    value: ['k8s'],
    onChange,
    onBlur: vi.fn(),
    onFocus: vi.fn(),
    ...overrides,
  };
  function Field() {
    const [value, setValue] = useState(props.value);
    return React.createElement(CustomSelectWidget as React.ComponentType<any>, {
      ...props,
      value,
      onChange: (next) => {
        setValue(next);
        onChange(next);
      },
    });
  }
  render(
    <SistentThemeProvider>
      <Field />
    </SistentThemeProvider>,
  );
  return { user: userEvent.setup(), onChange, combo: screen.getByRole('combobox') };
}

// Exercise the published Sistent/MUI controls: a TextField mock cannot detect
// ignored legacy props, nested checkbox focus, or broken select value rendering.
describe('CustomSelectWidget', () => {
  it('shows labels without checkboxes in the closed control and toggles multiple enum values', async () => {
    const { user, onChange, combo } = setup();
    expect(combo).toHaveTextContent('Kubernetes');
    expect(within(combo).queryByRole('checkbox', { hidden: true })).toBeNull();
    await user.click(combo);
    const list = screen.getByRole('listbox');
    expect(list).toHaveAttribute('aria-multiselectable', 'true');
    const options = within(list).getAllByRole('option');
    expect(options[0]).toHaveAttribute('aria-selected', 'true');
    expect(options[1]).toHaveAttribute('aria-selected', 'false');
    const checks = within(list).getAllByRole('checkbox', { hidden: true });
    expect(checks[0]).toBeChecked();
    expect(checks[1]).not.toBeChecked();
    expect(checks[0]).toHaveAttribute('tabindex', '-1');
    await user.click(options[1]);
    expect(onChange).toHaveBeenLastCalledWith(['k8s', 'istio']);
    await user.click(options[0]);
    expect(onChange).toHaveBeenLastCalledWith(['istio']);
    await user.keyboard('{Escape}');
    expect(combo).toHaveTextContent('Istio');
  });

  it('keeps single-select fields free of checkboxes and emits a scalar', async () => {
    const { user, onChange, combo } = setup({ schema: { type: 'string' }, value: 'k8s' });
    await user.click(combo);
    expect(screen.queryByRole('checkbox', { hidden: true })).toBeNull();
    await user.click(screen.getByRole('option', { name: 'Istio' }));
    expect(onChange).toHaveBeenLastCalledWith('istio');
    expect(combo).toHaveTextContent('Istio');
  });

  it('supports keyboard selection and skips disabled enum options', async () => {
    const { user, onChange, combo } = setup({ options: { enumOptions, enumDisabled: ['istio'] } });
    combo.focus();
    await user.keyboard('{ArrowDown}');
    expect(screen.getByRole('option', { name: 'Istio' })).toHaveAttribute('aria-disabled', 'true');
    await user.keyboard('{ArrowDown}{Enter}');
    expect(onChange).toHaveBeenLastCalledWith(['k8s', 'long']);
    await user.keyboard('{Escape}');
    expect(combo).toHaveFocus();
  });

  it('exposes the full option label in a hover tooltip', async () => {
    const { user, combo } = setup();
    await user.click(combo);
    await user.hover(within(screen.getByRole('listbox')).getByText(enumOptions[2].label));
    expect(await screen.findByRole('tooltip')).toHaveTextContent(enumOptions[2].label);
  });

  it('preserves rich labels in the comma-separated closed value', () => {
    const { combo } = setup({
      options: {
        enumOptions: [
          { value: 'k8s', label: <strong>Kubernetes</strong> },
          { value: 'istio', label: <em>Istio</em> },
        ],
      },
      value: ['k8s', 'istio'],
    });
    expect(combo).toHaveTextContent('Kubernetes, Istio');
    expect(combo.querySelector('strong')).toHaveTextContent('Kubernetes');
    expect(combo.querySelector('em')).toHaveTextContent('Istio');
    expect(combo).not.toHaveTextContent('[object Object]');
  });

  it.each([true, false])('honors hideError=%s without hiding field help', (hideError) => {
    setup({
      rawErrors: ['required'],
      hideError,
      schema: { type: 'array', description: 'Choose technologies' },
    });
    expect(screen.queryByLabelText('Validation errors') !== null).toBe(!hideError);
    expect(screen.getByLabelText('Field help')).toBeInTheDocument();
    const input = document.querySelector('input[name="technology"]');
    expect(input).toHaveAttribute('aria-invalid', String(!hideError));
  });

  it.each([{ disabled: true }, { readonly: true }])('disables the control for %j', (overrides) => {
    const { combo } = setup(overrides);
    expect(combo).toHaveAttribute('aria-disabled', 'true');
  });

  it('handles null and empty values and an external label', () => {
    const { combo } = setup({ value: null, schema: { type: 'array', 'x-rjsf-grid-area': 6 } });
    expect(combo).not.toHaveTextContent('Kubernetes');
    expect(document.querySelector('label[for="technology"]')).toHaveTextContent('Technology');
  });

  it('merges legacy and callback slot props, preserves adornments and menu defaults', async () => {
    const paper = vi.fn(() => ({
      style: { color: 'rgb(1, 2, 3)' },
      'data-testid': 'select-paper',
    }));
    const { user, combo } = setup({
      schema: { type: 'array', description: 'Choose technologies' },
      InputProps: { style: { borderTopWidth: 3 }, endAdornment: <span>Caller adornment</span> },
      InputLabelProps: { shrink: false },
      SelectProps: { MenuProps: { PaperProps: { style: { backgroundColor: 'rgb(4, 5, 6)' } } } },
      slotProps: {
        input: () => ({ style: { borderBottomWidth: 5 } }),
        inputLabel: () => ({ 'data-testid': 'select-label', shrink: false }),
        select: () => ({ MenuProps: { slotProps: { paper } } }),
      },
    });
    expect(screen.getByText('Caller adornment')).toBeInTheDocument();
    expect(screen.getByLabelText('Field help')).toBeInTheDocument();
    expect(combo.parentElement).toHaveStyle({ borderTopWidth: '3px', borderBottomWidth: '5px' });
    expect(screen.getByTestId('select-label')).toHaveAttribute('data-shrink', 'true');
    await user.click(combo);
    expect(paper).toHaveBeenCalled();
    expect(screen.getByTestId('select-paper')).not.toHaveClass('MuiInputBase-input');
    expect(screen.getByTestId('select-paper')).not.toHaveClass('MuiOutlinedInput-input');
    expect(screen.getByTestId('select-paper')).toHaveStyle({
      maxHeight: '400px',
      color: 'rgb(1, 2, 3)',
      backgroundColor: 'rgb(4, 5, 6)',
    });
  });
});
