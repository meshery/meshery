import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../css/icons.styles', () => ({
  iconLarge: { width: 40 },
}));

import LinkedinIcon from './linkedinIcon';

describe('LinkedinIcon', () => {
  it('renders with default dimensions when no props are supplied', () => {
    const { container } = render(<LinkedinIcon />);
    const svg = container.querySelector('svg');
    expect(svg).not.toBeNull();
    expect(svg).toHaveAttribute('width', '40');
    expect(svg).toHaveAttribute('height', '40');
    expect(svg).toHaveAttribute('viewBox', '0 0 29 29');
  });

  it('honors width and height overrides', () => {
    const { container } = render(<LinkedinIcon width={64} height={48} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '64');
    expect(svg).toHaveAttribute('height', '48');
  });

  it('renders inner LinkedIn glyph paths', () => {
    const { container } = render(<LinkedinIcon />);
    expect(container.querySelectorAll('path').length).toBeGreaterThan(0);
    expect(container.querySelector('filter')).not.toBeNull();
  });
});
