import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import MesheryIcon from './meshery-icon';

describe('MesheryIcon', () => {
  it('renders an svg with default props', () => {
    const { container } = render(<MesheryIcon />);
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(svg).toHaveAttribute('viewBox', '0 0 134.95 135.02');
    expect(svg).toHaveAttribute('xmlns', 'http://www.w3.org/2000/svg');
  });

  it('uses default primary and secondary fill colors when not overridden', () => {
    const { container } = render(<MesheryIcon />);
    const polygons = container.querySelectorAll('polygon');
    // Polygons should fill with default primary color or secondary color
    expect(polygons.length).toBeGreaterThan(0);
    const fillValues = Array.from(polygons).map((el) => el.getAttribute('fill'));
    expect(fillValues).toContain('#00b39f');
    expect(fillValues).toContain('#00d3a9');
  });

  it('uses custom primary and secondary fill colors when provided', () => {
    const { container } = render(<MesheryIcon primaryFill="#111111" secondaryFill="#222222" />);
    const polygons = container.querySelectorAll('polygon');
    const fillValues = Array.from(polygons).map((el) => el.getAttribute('fill'));
    expect(fillValues).toContain('#111111');
    expect(fillValues).toContain('#222222');
  });

  it('forwards width and height props', () => {
    const { container } = render(<MesheryIcon width={80} height={90} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '80');
    expect(svg).toHaveAttribute('height', '90');
  });

  it('applies passed style', () => {
    const { container } = render(<MesheryIcon style={{ padding: '4px' }} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveStyle({ padding: '4px' });
  });
});
