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

vi.mock('@/theme', () => ({
  useTheme: () => ({ palette: { primary: { main: '#abc' } } }),
}));

vi.mock('../helper', () => ({
  safeDisplayValue: (v: any) => (v == null ? '' : String(v)),
}));

import CustomFileWidget from './CustomFileWidget';

describe('CustomFileWidget', () => {
  beforeEach(() => {
    baseInputMock.mockClear();
  });

  it('renders BaseInput with file input type', () => {
    render(<CustomFileWidget id="f1" options={{}} onChange={vi.fn()} />);
    expect(screen.getByTestId('base-input')).toHaveAttribute('data-input-type', 'file');
  });

  it('shows a default value mask label when options.default is provided', () => {
    render(
      <CustomFileWidget id="f1" options={{ default: 'placeholder.png' }} onChange={vi.fn()} />,
    );
    expect(screen.getByText('placeholder.png')).toBeInTheDocument();
  });
});
