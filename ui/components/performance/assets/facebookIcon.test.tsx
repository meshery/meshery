import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../css/icons.styles', () => ({
  iconLarge: { width: 40 },
}));

import FacebookIcon from './facebookIcon';

describe('FacebookIcon', () => {
  it('renders with default dimensions when no props are supplied', () => {
    const { container } = render(<FacebookIcon />);
    const svg = container.querySelector('svg');
    expect(svg).not.toBeNull();
    expect(svg).toHaveAttribute('width', '40');
    expect(svg).toHaveAttribute('height', '40');
    expect(svg).toHaveAttribute('viewBox', '0 0 28 29');
  });

  it('honors width and height overrides', () => {
    const { container } = render(<FacebookIcon width={50} height={50} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '50');
    expect(svg).toHaveAttribute('height', '50');
  });

  it('renders inner glyph paths', () => {
    const { container } = render(<FacebookIcon />);
    expect(container.querySelectorAll('path').length).toBeGreaterThan(0);
  });
});
