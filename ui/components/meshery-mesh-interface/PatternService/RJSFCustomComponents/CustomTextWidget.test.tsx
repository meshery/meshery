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

import CustomTextWidget from './CustomTextWidget';

describe('CustomTextWidget', () => {
  it('defaults to inputType=string when schema type is unset', () => {
    baseInputMock.mockClear();
    render(<CustomTextWidget id="t1" schema={{}} options={{}} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({ options: expect.objectContaining({ inputType: 'string' }) }),
    );
  });

  it('uses inputType=number when schema type is number', () => {
    baseInputMock.mockClear();
    render(<CustomTextWidget id="t1" schema={{ type: 'number' }} options={{}} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({ options: expect.objectContaining({ inputType: 'number' }) }),
    );
  });

  it('uses inputType=number when schema type is integer', () => {
    baseInputMock.mockClear();
    render(<CustomTextWidget id="t1" schema={{ type: 'integer' }} options={{}} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({ options: expect.objectContaining({ inputType: 'number' }) }),
    );
  });
});
