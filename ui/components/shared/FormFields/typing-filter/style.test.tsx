import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = (props: any) => <Component {...props}>{props.children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    TextField: ({ children, ...rest }: any) => <input data-testid="text-field" {...rest} />,
    styled,
  };
});

import { Root, InputField, DropDown } from './style';

describe('typing-filter style', () => {
  it('renders Root as a div', () => {
    render(<Root data-testid="root">Hi</Root>);
    expect(screen.getByTestId('root')).toBeInTheDocument();
  });

  it('renders InputField using underlying TextField', () => {
    render(<InputField placeholder="search" />);
    expect(screen.getByTestId('text-field')).toBeInTheDocument();
    expect(screen.getByTestId('text-field')).toHaveAttribute('placeholder', 'search');
  });

  it('renders DropDown as a div', () => {
    render(<DropDown data-testid="dropdown">Options</DropDown>);
    expect(screen.getByTestId('dropdown')).toHaveTextContent('Options');
  });
});
