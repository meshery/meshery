import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import InfoIcon from './info-icon';

describe('InfoIcon', () => {
  it('renders an svg with default props', () => {
    const { container } = render(<InfoIcon />);
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(svg).toHaveAttribute('viewBox', '0 0 21 22');
    expect(svg).toHaveAttribute('xmlns', 'http://www.w3.org/2000/svg');
  });

  it('forwards width, height, and fill props', () => {
    const { container } = render(<InfoIcon width={32} height={48} fill="red" />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '32');
    expect(svg).toHaveAttribute('height', '48');
    expect(svg).toHaveAttribute('fill', 'red');
  });

  it('applies passed style', () => {
    const { container } = render(<InfoIcon style={{ marginLeft: '4px' }} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveStyle({ marginLeft: '4px' });
  });

  it('renders both path elements (icon contents)', () => {
    const { container } = render(<InfoIcon />);
    const paths = container.querySelectorAll('path');
    expect(paths.length).toBe(2);
  });
});
