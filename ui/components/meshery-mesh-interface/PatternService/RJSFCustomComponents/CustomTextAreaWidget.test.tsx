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

import CustomTextAreaWidget from './CustomTextAreaWidget';

describe('CustomTextAreaWidget', () => {
  it('defaults rows to 3 when options.rows is not provided', () => {
    baseInputMock.mockClear();
    render(<CustomTextAreaWidget options={{}} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({
        multiline: true,
        options: expect.objectContaining({ rows: 3 }),
      }),
    );
  });

  it('uses the provided row count from options', () => {
    baseInputMock.mockClear();
    render(<CustomTextAreaWidget options={{ rows: 7 }} />);
    expect(baseInputMock).toHaveBeenCalledWith(
      expect.objectContaining({
        options: expect.objectContaining({ rows: 7 }),
      }),
    );
  });
});
