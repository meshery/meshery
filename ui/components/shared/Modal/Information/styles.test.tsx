import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CreatAtContainer, ActionContainer, CopyLinkButton, ResourceName } from './styles';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Styled = ({ children, ...props }: any) => <Component {...props}>{children}</Component>;
    Styled.displayName = 'StyledMock';
    return Styled;
  };
  return {
    Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
    Button: ({ children, ...rest }: any) => <button {...rest}>{children}</button>,
    Typography: ({ children, ...rest }: any) => <span {...rest}>{children}</span>,
    styled,
  };
});

describe('Info Modal styles', () => {
  it('CreatAtContainer renders children with isBold true', () => {
    render(<CreatAtContainer isBold={true}>Bold</CreatAtContainer>);
    expect(screen.getByText('Bold')).toBeInTheDocument();
  });

  it('CreatAtContainer renders children with isBold false', () => {
    render(<CreatAtContainer isBold={false}>Light</CreatAtContainer>);
    expect(screen.getByText('Light')).toBeInTheDocument();
  });

  it('ActionContainer renders children', () => {
    render(<ActionContainer>action</ActionContainer>);
    expect(screen.getByText('action')).toBeInTheDocument();
  });

  it('CopyLinkButton renders children', () => {
    render(<CopyLinkButton>Copy</CopyLinkButton>);
    expect(screen.getByText('Copy')).toBeInTheDocument();
  });

  it('ResourceName renders children', () => {
    render(<ResourceName>Resource</ResourceName>);
    expect(screen.getByText('Resource')).toBeInTheDocument();
  });
});
