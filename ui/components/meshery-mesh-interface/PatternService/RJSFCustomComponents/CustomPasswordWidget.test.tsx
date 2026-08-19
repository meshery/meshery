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

import CustomPasswordWidget from './CustomPasswordWidget';

describe('CustomPasswordWidget', () => {
  it('passes inputType=password to BaseInput', () => {
    baseInputMock.mockClear();
    render(
      <CustomPasswordWidget
        id="secret"
        options={{ foo: 'bar' }}
        value="token"
        onChange={vi.fn()}
      />,
    );
    expect(screen.getByTestId('base-input')).toHaveAttribute('data-input-type', 'password');
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({
        options: expect.objectContaining({ inputType: 'password', foo: 'bar' }),
      }),
    );
  });

  it('still works without options', () => {
    baseInputMock.mockClear();
    render(<CustomPasswordWidget id="secret" value="" onChange={vi.fn()} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({
        options: expect.objectContaining({ inputType: 'password' }),
      }),
    );
  });
});
