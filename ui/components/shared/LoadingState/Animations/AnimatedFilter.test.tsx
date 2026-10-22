import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import AnimatedFilter from './AnimatedFilter';

describe('AnimatedFilter', () => {
  it('renders a 100x100 svg', () => {
    const { container } = render(<AnimatedFilter />);
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(svg).toHaveAttribute('width', '100');
    expect(svg).toHaveAttribute('height', '100');
  });

  it('renders the wa-1, wa-2 and wa-3 svg classes initially', () => {
    const { container } = render(<AnimatedFilter />);
    expect(container.querySelector('.svg-wa-1.active')).toBeInTheDocument();
    expect(container.querySelector('.svg-wa-2.active')).toBeInTheDocument();
    expect(container.querySelector('.svg-wa-3.active')).toBeInTheDocument();
  });

  it('accepts custom SVG props', () => {
    const { container } = render(<AnimatedFilter className="custom-filter" />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveClass('custom-filter');
  });
});
