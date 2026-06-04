import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import LoadingScreen, { PureHtmlLoadingScreen } from './LoadingComponentServer';

vi.mock('../../../ui.config', () => ({
  default: {
    AnimatedLogoDark: () => <div data-testid="animated-logo-dark" />,
  },
}));

describe('LoadingScreen', () => {
  it('renders the loading screen when isLoading is true', () => {
    render(
      <LoadingScreen message="Loading data..." isLoading={true} id="loader-1">
        <span data-testid="children">should not show</span>
      </LoadingScreen>,
    );

    expect(screen.getByTestId('animated-logo-dark')).toBeInTheDocument();
    expect(screen.getByText('Loading data...')).toBeInTheDocument();
    expect(screen.queryByTestId('children')).not.toBeInTheDocument();
  });

  it('renders children when isLoading is false', () => {
    render(
      <LoadingScreen message="ignored" isLoading={false} id="loader-2">
        <span data-testid="children">Children rendered</span>
      </LoadingScreen>,
    );

    expect(screen.queryByTestId('animated-logo-dark')).not.toBeInTheDocument();
    expect(screen.getByTestId('children')).toHaveTextContent('Children rendered');
  });

  it('forwards id prop to the message heading suffix', () => {
    render(<LoadingScreen message="Loading..." isLoading={true} id="my-id" />);

    const heading = screen.getByText('Loading...');
    expect(heading).toHaveAttribute('id', 'my-id-text-message');
  });
});

describe('PureHtmlLoadingScreen', () => {
  it('renders the message in a heading with the appended -text-message id', () => {
    render(<PureHtmlLoadingScreen message="Working" id="abc" />);

    const heading = screen.getByText('Working');
    expect(heading.tagName.toLowerCase()).toBe('h1');
    expect(heading).toHaveAttribute('id', 'abc-text-message');
  });

  it('merges custom style overrides with the defaults', () => {
    render(
      <PureHtmlLoadingScreen
        message="Loading"
        id="merge"
        style={{ color: 'red', backgroundColor: '#ffffff' }}
      />,
    );

    const outer = screen.getByText('Loading').parentElement?.parentElement;
    expect(outer).toHaveStyle({ color: 'rgb(255, 0, 0)' });
  });
});
