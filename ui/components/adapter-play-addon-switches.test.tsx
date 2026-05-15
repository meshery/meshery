import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  FormControl: ({ children, component }: any) => (
    <div data-testid="form-control" data-component={component}>
      {children}
    </div>
  ),
  FormControlLabel: ({ control, label }: any) => (
    <label data-testid="form-control-label">
      {control}
      <span>{label}</span>
    </label>
  ),
  FormGroup: ({ children }: any) => <div data-testid="form-group">{children}</div>,
  FormLabel: ({ children }: any) => <div data-testid="form-label">{children}</div>,
  Switch: ({ checked, onChange, name }: any) => (
    <input
      data-testid={`switch-${name}`}
      type="checkbox"
      name={name}
      checked={!!checked}
      onChange={onChange}
    />
  ),
}));

import AdapterAddonSwitches from './adapter-play-addon-switches';

describe('AdapterAddonSwitches', () => {
  it('returns null when given no ops', () => {
    const { container } = render(
      <AdapterAddonSwitches
        selectedAdapterOps={[]}
        addonSwitchGroup={{}}
        onSwitchChange={vi.fn()}
      />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('renders a Customize Addons section with one switch per op', () => {
    render(
      <AdapterAddonSwitches
        selectedAdapterOps={[
          { key: 'op1', value: 'Add-on:Tracing' },
          { key: 'op2', value: 'Add-on:Metrics' },
        ]}
        addonSwitchGroup={{ op1: true, op2: false }}
        onSwitchChange={vi.fn()}
      />,
    );

    expect(screen.getByTestId('form-label')).toHaveTextContent('Customize Addons');
    expect(screen.getByTestId('switch-op1')).toBeChecked();
    expect(screen.getByTestId('switch-op2')).not.toBeChecked();
  });

  it('strips the Add-on: prefix from the rendered label', () => {
    render(
      <AdapterAddonSwitches
        selectedAdapterOps={[{ key: 'op1', value: 'Add-on:Tracing' }]}
        addonSwitchGroup={{}}
        onSwitchChange={vi.fn()}
      />,
    );

    expect(screen.getByTestId('form-control-label')).toHaveTextContent('Tracing');
    expect(screen.getByTestId('form-control-label')).not.toHaveTextContent('Add-on:');
  });

  it('sorts switches alphabetically by their cleaned label', () => {
    render(
      <AdapterAddonSwitches
        selectedAdapterOps={[
          { key: 'opb', value: 'Add-on:Bravo' },
          { key: 'opa', value: 'Add-on:Alpha' },
          { key: 'opc', value: 'Add-on:Charlie' },
        ]}
        addonSwitchGroup={{}}
        onSwitchChange={vi.fn()}
      />,
    );

    const labels = screen.getAllByTestId('form-control-label');
    expect(labels[0]).toHaveTextContent('Alpha');
    expect(labels[1]).toHaveTextContent('Bravo');
    expect(labels[2]).toHaveTextContent('Charlie');
  });

  it('calls onSwitchChange with name, checked, and the original op', () => {
    const onSwitchChange = vi.fn();
    const ops = [{ key: 'op1', value: 'Add-on:Tracing' }];
    render(
      <AdapterAddonSwitches
        selectedAdapterOps={ops}
        addonSwitchGroup={{}}
        onSwitchChange={onSwitchChange}
      />,
    );

    const sw = screen.getByTestId('switch-op1');
    fireEvent.click(sw);

    expect(onSwitchChange).toHaveBeenCalledTimes(1);
    const [name, checked, op] = onSwitchChange.mock.calls[0];
    expect(name).toBe('op1');
    expect(checked).toBe(true);
    // The op forwarded is the post-map version (Add-on: prefix stripped).
    expect(op.key).toBe('op1');
  });

  it('treats undefined entries in addonSwitchGroup as unchecked', () => {
    render(
      <AdapterAddonSwitches
        selectedAdapterOps={[{ key: 'op1', value: 'Add-on:Tracing' }]}
        addonSwitchGroup={{}}
        onSwitchChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('switch-op1')).not.toBeChecked();
  });
});
