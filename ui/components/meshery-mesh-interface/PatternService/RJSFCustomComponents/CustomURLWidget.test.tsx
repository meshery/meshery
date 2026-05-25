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

import CustomURLWidget from './CustomURLWidget';

describe('CustomURLWidget', () => {
  it('forwards options through to BaseInput', () => {
    baseInputMock.mockClear();
    render(<CustomURLWidget options={{ foo: 'bar' }} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({ options: expect.objectContaining({ foo: 'bar' }) }),
    );
  });
});
