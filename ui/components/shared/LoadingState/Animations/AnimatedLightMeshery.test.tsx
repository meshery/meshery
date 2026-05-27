import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import AnimatedLightMeshery from './AnimatedLightMeshery';

describe('AnimatedLightMeshery', () => {
  it('renders the SVG with the meshery title', () => {
    const { container } = render(<AnimatedLightMeshery />);

    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(container.querySelector('title')).toHaveTextContent('meshery-logo-light-text');
  });

  it('uses white text fill', () => {
    const { container } = render(<AnimatedLightMeshery />);
    const whiteFilled = container.querySelectorAll('[fill="#FFF"]');
    expect(whiteFilled.length).toBeGreaterThan(0);
  });

  it('renders meshery polygon classes with active state', () => {
    const { container } = render(<AnimatedLightMeshery />);
    expect(container.querySelector('.svg-meshery-1.active')).toBeInTheDocument();
    expect(container.querySelector('.svg-meshery-25.active')).toBeInTheDocument();
  });

  it('forwards SVG props', () => {
    const { container } = render(<AnimatedLightMeshery className="custom-meshery" />);
    expect(container.querySelector('svg')).toHaveClass('custom-meshery');
  });
});
