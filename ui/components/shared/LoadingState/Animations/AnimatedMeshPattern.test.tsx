import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import AnimatedMeshPattern from './AnimatedMeshPattern';

let mockMode: 'light' | 'dark' = 'light';

vi.mock('@sistent/sistent', () => ({
  useTheme: () => ({ palette: { mode: mockMode } }),
}));

describe('AnimatedMeshPattern', () => {
  it('renders the SVG with the service-mesh-pattern-text title', () => {
    mockMode = 'light';
    const { container } = render(<AnimatedMeshPattern />);

    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(container.querySelector('title')).toHaveTextContent('service-mesh-pattern-text');
  });

  it('uses light theme text color when in light mode', () => {
    mockMode = 'light';
    const { container } = render(<AnimatedMeshPattern />);
    const filled = container.querySelectorAll('[fill="#3c494f"]');
    expect(filled.length).toBeGreaterThan(0);
  });

  it('uses white text color when in dark mode', () => {
    mockMode = 'dark';
    const { container } = render(<AnimatedMeshPattern />);
    const filled = container.querySelectorAll('[fill="#FFF"]');
    expect(filled.length).toBeGreaterThan(0);
  });

  it('renders pattern paths with the active class', () => {
    const { container } = render(<AnimatedMeshPattern />);
    expect(container.querySelector('.svg-pattern-1.active')).toBeInTheDocument();
    expect(container.querySelector('.svg-pattern-8.active')).toBeInTheDocument();
  });

  it('forwards SVG attributes from props', () => {
    const { container } = render(<AnimatedMeshPattern className="custom-pattern" />);
    expect(container.querySelector('svg')).toHaveClass('custom-pattern');
  });
});
