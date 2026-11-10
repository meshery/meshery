import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { CurvedArrowIcon } from './curvedArrowIcon';

describe('CurvedArrowIcon', () => {
  it('renders with default size and viewBox', () => {
    const { container } = render(<CurvedArrowIcon />);
    const svg = container.querySelector('svg');
    expect(svg).not.toBeNull();
    expect(svg).toHaveAttribute('width', '4.5rem');
    expect(svg).toHaveAttribute('height', '4.5rem');
    expect(svg).toHaveAttribute('viewBox', '0 0 50 73');
  });

  it('accepts explicit width and height props', () => {
    const { container } = render(<CurvedArrowIcon width="6rem" height={120} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '6rem');
    expect(svg).toHaveAttribute('height', '120');
  });

  it('forwards extra props to the svg element', () => {
    const { container } = render(<CurvedArrowIcon props={{ 'data-testid': 'arrow' } as any} />);
    expect(container.querySelector('[data-testid="arrow"]')).not.toBeNull();
  });
});
