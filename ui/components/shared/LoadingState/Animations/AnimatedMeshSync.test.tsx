import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import AnimatedMeshSync from './AnimatedMeshSync';

describe('AnimatedMeshSync', () => {
  it('renders the SVG with the meshsync title', () => {
    const { container } = render(<AnimatedMeshSync />);

    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(container.querySelector('title')).toHaveTextContent('meshsync');
  });

  it('renders sync paths with the active class', () => {
    const { container } = render(<AnimatedMeshSync />);

    expect(container.querySelector('.svg-sync-1.active')).toBeInTheDocument();
    expect(container.querySelector('.svg-sync-5.active')).toBeInTheDocument();
  });

  it('forwards SVG attributes from props', () => {
    const { container } = render(<AnimatedMeshSync className="sync-class" />);
    expect(container.querySelector('svg')).toHaveClass('sync-class');
  });
});
