import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  // styled(Component) returns a function which returns a component that
  // forwards all props onto the wrapped Component — matching the real
  // emotion/MUI behaviour we care about here.
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => (
      <Component {...props}>{children}</Component>
    );
    StyledComponent.displayName = 'StyledMock';
    return StyledComponent;
  };

  return {
    NoSsr: ({ children }: any) => <>{children}</>,
    Typography: ({ children, variant }: any) => (
      <div data-testid={`typography-${variant}`}>{children}</div>
    ),
    Link: ({ children, href, target }: any) => (
      <a data-testid="link" href={href} target={target}>
        {children}
      </a>
    ),
    Box: ({ children, component, src, alt, ...rest }: any) => {
      if (component === 'img') {
        return <img data-testid="box-img" src={src} alt={alt} {...rest} />;
      }
      return (
        <div data-testid="box" {...rest}>
          {children}
        </div>
      );
    },
    styled,
  };
});

import CustomErrorMessage from './ErrorPage';

describe('CustomErrorMessage (ErrorPage)', () => {
  it('renders one of the random meshy messages and the default subheading', () => {
    render(<CustomErrorMessage />);

    expect(screen.getByTestId('typography-h1')).toBeInTheDocument();
    expect(screen.getByTestId('typography-h5')).toHaveTextContent('Page does not exist.');
  });

  it('uses a provided message in place of the default subheading', () => {
    render(<CustomErrorMessage message="Forbidden" />);
    expect(screen.getByTestId('typography-h5')).toHaveTextContent('Forbidden');
  });

  it('renders the discussion forum link', () => {
    render(<CustomErrorMessage />);
    const link = screen.getByTestId('link');
    expect(link).toHaveAttribute(
      'href',
      'https://meshery.io/community#community-forums/c/meshery/5',
    );
    expect(link).toHaveAttribute('target', '_blank');
    expect(link).toHaveTextContent('discussion forum');
  });

  it('renders the image by default', () => {
    render(<CustomErrorMessage />);
    const img = screen.getByTestId('box-img');
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', '/static/img/designs/service-mesh.svg');
    expect(img).toHaveAttribute('alt', 'service meshed');
  });

  it('omits the image when showImage is false', () => {
    render(<CustomErrorMessage showImage={false} />);
    expect(screen.queryByTestId('box-img')).not.toBeInTheDocument();
  });

  it('renders a non-empty random meshy heading', () => {
    render(<CustomErrorMessage />);
    expect(screen.getByTestId('typography-h1').textContent?.length).toBeGreaterThan(0);
  });
});
