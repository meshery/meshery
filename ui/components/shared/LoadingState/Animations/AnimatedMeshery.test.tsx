import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import AnimatedMeshery, { AnimatedMesheryDark } from './AnimatedMeshery';

vi.mock('@/utils/hooks', () => ({
  useTimeout: vi.fn(),
}));

describe('AnimatedMeshery', () => {
  it('renders the SVG with the default text fill color', () => {
    const { container } = render(<AnimatedMeshery />);

    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    const filled = container.querySelectorAll('[fill="#000000"]');
    expect(filled.length).toBeGreaterThan(0);
  });

  it('respects the textFill prop', () => {
    const { container } = render(<AnimatedMeshery textFill="#abcdef" />);
    const filled = container.querySelectorAll('[fill="#abcdef"]');
    expect(filled.length).toBeGreaterThan(0);
  });

  it('renders an svg with title meshery-logo-light-text', () => {
    const { container } = render(<AnimatedMeshery />);
    expect(container.querySelector('title')).toHaveTextContent('meshery-logo-light-text');
  });

  it('has the polygon points', () => {
    const { container } = render(<AnimatedMeshery />);
    const polygons = container.querySelectorAll('polygon');
    expect(polygons.length).toBeGreaterThan(0);
  });
});

describe('AnimatedMesheryDark', () => {
  it('forces the text fill to white', () => {
    const { container } = render(<AnimatedMesheryDark />);
    const filled = container.querySelectorAll('[fill="#fff"]');
    expect(filled.length).toBeGreaterThan(0);
  });
});
