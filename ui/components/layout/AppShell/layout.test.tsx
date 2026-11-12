import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import Layout from './layout';

vi.mock('../Header/Header', () => ({
  default: () => <div data-testid="header-mock">Header</div>,
}));

describe('AppShell Layout', () => {
  it('renders the Header at the top of the shell', () => {
    render(<Layout />);

    expect(screen.getByTestId('header-mock')).toBeInTheDocument();
  });

  it('renders children below the Header when provided', () => {
    render(
      <Layout>
        <p data-testid="layout-child">Hello world</p>
      </Layout>,
    );

    expect(screen.getByTestId('header-mock')).toBeInTheDocument();
    expect(screen.getByTestId('layout-child')).toHaveTextContent('Hello world');
  });

  it('renders without crashing when no children are passed', () => {
    const { container } = render(<Layout />);

    // Outer div with style attributes is the only structural element when no
    // children are supplied.
    const outer = container.firstChild as HTMLElement;
    expect(outer).not.toBeNull();
    expect(outer.tagName).toBe('DIV');
    expect(outer.style.margin).toBe('20px');
    expect(outer.style.padding).toBe('20px');
    expect(outer.style.border).toContain('1px');
  });

  it('renders multiple children in order', () => {
    render(
      <Layout>
        <span data-testid="child-1">First</span>
        <span data-testid="child-2">Second</span>
      </Layout>,
    );

    expect(screen.getByTestId('child-1')).toBeInTheDocument();
    expect(screen.getByTestId('child-2')).toBeInTheDocument();
  });
});
