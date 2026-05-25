import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import AnimatedMeshery, { AnimatedMesheryDark } from './AnimatedMesheryCSS';

describe('AnimatedMesheryCSS', () => {
  it('renders SVG with the meshery title', () => {
    const { container } = render(<AnimatedMeshery />);
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(container.querySelector('title')).toHaveTextContent('meshery-logo-light-text');
  });

  it('uses default textFill #ccc when not provided', () => {
    const { container } = render(<AnimatedMeshery />);
    const filled = container.querySelectorAll('[fill="#ccc"]');
    expect(filled.length).toBeGreaterThan(0);
  });

  it('uses custom textFill when provided', () => {
    const { container } = render(<AnimatedMeshery textFill="#ff0000" />);
    const filled = container.querySelectorAll('[fill="#ff0000"]');
    expect(filled.length).toBeGreaterThan(0);
  });

  it('all polygons have the active class baked in', () => {
    const { container } = render(<AnimatedMeshery />);
    expect(container.querySelector('.svg-meshery-1.active')).toBeInTheDocument();
    expect(container.querySelector('.svg-meshery-25.active')).toBeInTheDocument();
  });
});

describe('AnimatedMesheryCSS Dark variant', () => {
  it('forwards through to the base component', () => {
    const { container } = render(<AnimatedMesheryDark textFill="#aaa" />);
    expect(container.querySelector('svg')).toBeInTheDocument();
    const filled = container.querySelectorAll('[fill="#aaa"]');
    expect(filled.length).toBeGreaterThan(0);
  });
});
