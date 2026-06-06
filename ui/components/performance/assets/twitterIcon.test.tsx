import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../css/icons.styles', () => ({
  iconLarge: { width: 40 },
}));

import TwitterIcon from './twitterIcon';

describe('TwitterIcon', () => {
  it('renders with default dimensions when no props are supplied', () => {
    const { container } = render(<TwitterIcon />);
    const svg = container.querySelector('svg');
    expect(svg).not.toBeNull();
    expect(svg).toHaveAttribute('width', '40');
    expect(svg).toHaveAttribute('height', '40');
    expect(svg).toHaveAttribute('viewBox', '0 0 29 29');
  });

  it('honors width and height overrides', () => {
    const { container } = render(<TwitterIcon width={32} height={28} />);
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '32');
    expect(svg).toHaveAttribute('height', '28');
  });

  it('renders inner glyph paths', () => {
    const { container } = render(<TwitterIcon />);
    expect(container.querySelectorAll('path').length).toBeGreaterThan(0);
  });
});
