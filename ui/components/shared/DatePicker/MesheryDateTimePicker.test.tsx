import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import moment from 'moment';
import MesheryDateTimePicker from './MesheryDateTimePicker';

let lastDateTimePickerProps: any = null;

vi.mock('@sistent/sistent', () => ({
  TextField: ({ label, value, onChange, type, disabled, fullWidth, InputLabelProps }: any) => (
    <input
      data-testid="text-field"
      aria-label={label}
      value={value}
      onChange={onChange}
      type={type}
      disabled={disabled}
      data-full-width={String(fullWidth)}
      data-shrink={String(InputLabelProps?.shrink ?? false)}
    />
  ),
}));

vi.mock('@mui/x-date-pickers/DateTimePicker', () => ({
  DateTimePicker: (props: any) => {
    lastDateTimePickerProps = props;
    return (
      <div data-testid="datetime-picker">
        <button aria-label="emit-error" onClick={() => props.onError(new Error('invalid'))}>
          err
        </button>
        <button
          aria-label="emit-valid"
          onClick={() => props.onChange(moment('2024-01-15T10:30:00'))}
        >
          valid
        </button>
        <button aria-label="emit-invalid" onClick={() => props.onChange(null)}>
          invalid
        </button>
        <span data-testid="picker-value">{String(props.value)}</span>
        <span data-testid="picker-disabled">{String(props.disabled)}</span>
        <span data-testid="picker-label">{props.label}</span>
        <span data-testid="picker-format">{props.format}</span>
      </div>
    );
  },
}));

describe('MesheryDateTimePicker', () => {
  beforeEach(() => {
    lastDateTimePickerProps = null;
  });

  it('renders the DateTimePicker with the provided label and format', () => {
    render(<MesheryDateTimePicker label="When?" onChange={vi.fn()} selectedDate={null} />);

    expect(screen.getByTestId('picker-label')).toHaveTextContent('When?');
    expect(screen.getByTestId('picker-format')).toHaveTextContent('YYYY-MM-DD, hh:mm:ss a');
  });

  it('passes the disabled prop through', () => {
    render(
      <MesheryDateTimePicker
        label="When?"
        onChange={vi.fn()}
        selectedDate={null}
        disabled={true}
      />,
    );

    expect(screen.getByTestId('picker-disabled')).toHaveTextContent('true');
  });

  it('invokes onChange when a valid date is selected', async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<MesheryDateTimePicker label="When?" onChange={onChange} selectedDate={null} />);

    await user.click(screen.getByRole('button', { name: 'emit-valid' }));

    expect(onChange).toHaveBeenCalled();
    expect(moment.isMoment(onChange.mock.calls[0][0])).toBe(true);
  });

  it('does not call onChange when an invalid value is selected', async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    render(<MesheryDateTimePicker label="When?" onChange={onChange} selectedDate={null} />);

    await user.click(screen.getByRole('button', { name: 'emit-invalid' }));
    expect(onChange).not.toHaveBeenCalled();
  });

  it('switches to native picker when the DateTimePicker reports an error', async () => {
    const user = userEvent.setup();
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    render(<MesheryDateTimePicker label="When?" onChange={vi.fn()} selectedDate={null} />);

    await user.click(screen.getByRole('button', { name: 'emit-error' }));

    expect(screen.getByTestId('text-field')).toBeInTheDocument();
    expect(screen.queryByTestId('datetime-picker')).not.toBeInTheDocument();
    errorSpy.mockRestore();
  });

  it('renders native input with formatted selectedDate after error', async () => {
    const user = userEvent.setup();
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const date = '2024-03-15T08:30:45';
    render(<MesheryDateTimePicker label="When?" onChange={vi.fn()} selectedDate={date} />);

    await user.click(screen.getByRole('button', { name: 'emit-error' }));

    const field = screen.getByTestId('text-field') as HTMLInputElement;
    expect(field.value).toContain('2024-03-15T08:30:45');
    errorSpy.mockRestore();
  });

  it('parses native datetime input and forwards to onChange', async () => {
    const user = userEvent.setup();
    const onChange = vi.fn();
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    render(
      <MesheryDateTimePicker
        label="When?"
        onChange={onChange}
        selectedDate="2024-03-15T08:30:45"
      />,
    );

    await user.click(screen.getByRole('button', { name: 'emit-error' }));

    fireEvent.change(screen.getByTestId('text-field'), {
      target: { value: '2024-04-01T10:00:00' },
    });

    expect(onChange).toHaveBeenCalled();
    errorSpy.mockRestore();
  });

  it('wraps the picker in a div with className', () => {
    const { container } = render(
      <MesheryDateTimePicker
        label="When?"
        onChange={vi.fn()}
        selectedDate={null}
        className="custom-date"
      />,
    );

    expect(container.querySelector('.custom-date')).toBeInTheDocument();
  });
});
