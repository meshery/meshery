import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const baseInputMock = vi.fn();

vi.mock('./CustomBaseInput', () => ({
  default: (props: any) => {
    baseInputMock(props);
    return <div data-testid="base-input" />;
  },
}));

vi.mock('@rjsf/utils', () => ({
  utcToLocal: (v: string) => `local:${v}`,
  localToUTC: (v: string) => `utc:${v}`,
}));

import CustomDateTimeWidget from './CustomDateTimeWidget';

describe('CustomDateTimeWidget', () => {
  it('renders BaseInput with datetime-local input type', () => {
    baseInputMock.mockClear();
    render(<CustomDateTimeWidget value="2024-01-01" options={{}} onChange={vi.fn()} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({
        options: expect.objectContaining({ inputType: 'datetime-local', focused: true }),
        value: 'local:2024-01-01',
      }),
    );
  });

  it('converts the value back to UTC when onChange fires', () => {
    baseInputMock.mockClear();
    const onChange = vi.fn();
    render(<CustomDateTimeWidget value="2024-01-01" options={{}} onChange={onChange} />);
    const props = baseInputMock.mock.calls[0][0];
    props.onChange('2025-05-14');
    expect(onChange).toHaveBeenCalledWith('utc:2025-05-14');
  });
});
