import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import LoadingScreen from './LoadingComponent';

let mockMode: 'light' | 'dark' = 'light';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };

  return {
    Typography: ({ children, variant, component }: any) => (
      <span data-variant={variant} data-component={component}>
        {children}
      </span>
    ),
    useTheme: () => ({ palette: { mode: mockMode } }),
    styled,
  };
});

vi.mock('./Animations/AnimatedMeshSync', () => ({
  default: () => <div data-testid="animated-meshsync" />,
}));

vi.mock('./Animations/AnimatedMeshPattern', () => ({
  default: () => <div data-testid="animated-meshpattern" />,
}));

vi.mock('./Animations/AnimatedMeshery', () => ({
  default: () => <div data-testid="animated-meshery" />,
}));

vi.mock('./Animations/AnimatedFilter', () => ({
  default: () => <div data-testid="animated-filter" />,
}));

vi.mock('./Animations/AnimatedLightMeshery', () => ({
  default: () => <div data-testid="animated-light-meshery" />,
}));

describe('LoadingComponent (LoadingScreen)', () => {
  it('renders the page-loader container with the message', () => {
    mockMode = 'light';
    render(<LoadingScreen message="Just a moment" animatedIcon="AnimatedMeshPattern" />);

    expect(screen.getByTestId('page-loader')).toBeInTheDocument();
    expect(screen.getByText('Just a moment')).toBeInTheDocument();
  });

  it('uses the dark variant of meshery in dark mode', () => {
    mockMode = 'dark';
    render(<LoadingScreen message="Loading" animatedIcon="AnimatedMeshery" />);

    expect(screen.getByTestId('animated-light-meshery')).toBeInTheDocument();
    expect(screen.queryByTestId('animated-meshery')).not.toBeInTheDocument();
  });

  it('uses the light variant of meshery in light mode', () => {
    mockMode = 'light';
    render(<LoadingScreen message="Loading" animatedIcon="AnimatedMeshery" />);

    expect(screen.getByTestId('animated-meshery')).toBeInTheDocument();
    expect(screen.queryByTestId('animated-light-meshery')).not.toBeInTheDocument();
  });

  it('renders the AnimatedMeshSync icon with the corresponding logo image', () => {
    mockMode = 'light';
    render(<LoadingScreen message="Syncing" animatedIcon="AnimatedMeshSync" />);

    expect(screen.getByTestId('animated-meshsync')).toBeInTheDocument();
    expect(screen.getByAltText('mehsery-logo')).toHaveAttribute(
      'src',
      '/static/img/meshery-logo/meshery-black.svg',
    );
  });

  it('renders the AnimatedMeshSync white logo in dark mode', () => {
    mockMode = 'dark';
    render(<LoadingScreen message="Syncing" animatedIcon="AnimatedMeshSync" />);

    expect(screen.getByAltText('mehsery-logo')).toHaveAttribute(
      'src',
      '/static/img/meshery-logo/meshery-white.svg',
    );
  });

  it('renders AnimatedFilter', () => {
    mockMode = 'light';
    render(<LoadingScreen message="Filtering" animatedIcon="AnimatedFilter" />);

    expect(screen.getByTestId('animated-filter')).toBeInTheDocument();
  });

  it('passes className through to the loader container', () => {
    mockMode = 'light';
    render(
      <LoadingScreen
        message="Loading"
        animatedIcon="AnimatedMeshPattern"
        className="custom-loader-class"
      />,
    );

    expect(screen.getByTestId('page-loader')).toHaveAttribute('class', 'custom-loader-class');
  });
});
