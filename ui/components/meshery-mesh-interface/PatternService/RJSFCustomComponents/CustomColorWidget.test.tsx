import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const baseInputMock = vi.fn();

vi.mock('./CustomBaseInput', () => ({
  default: (props: any) => {
    baseInputMock(props);
    return <div data-testid="base-input" data-input-type={props.options?.inputType} />;
  },
}));

import CustomColorWidget from './CustomColorWidget';

describe('CustomColorWidget', () => {
  it('passes inputType=color to BaseInput', () => {
    baseInputMock.mockClear();
    render(<CustomColorWidget id="id1" options={{ foo: 'bar' }} value="#fff" onChange={vi.fn()} />);
    expect(screen.getByTestId('base-input')).toHaveAttribute('data-input-type', 'color');
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({
        options: expect.objectContaining({ inputType: 'color', foo: 'bar' }),
      }),
    );
  });

  it('still works without options', () => {
    baseInputMock.mockClear();
    render(<CustomColorWidget id="id1" value="#fff" onChange={vi.fn()} />);
    expect(baseInputMock).toHaveBeenCalled();
  });
});
