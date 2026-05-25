import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) =>
      typeof Component === 'string' ? <span {...props} /> : <Component {...props} />;
    Wrapped.displayName = 'StyledStyled';
    return Wrapped;
  };
  return {
    styled,
    Box: ({ children, ...rest }: any) => <div {...rest}>{children}</div>,
    DialogTitle: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
  };
});

vi.mock('@/assets/icons', () => ({
  AddCircleOutlined: () => null,
}));

import {
  AddIconStyled,
  BtnText,
  CreateButton,
  SearchWrapper,
  ViewSwitchButton,
  YamlDialogTitle,
  YamlDialogTitleText,
} from './MesheryPatterns.styled';

describe('designs/patterns/MesheryPatterns.styled', () => {
  it('exports the expected styled wrappers as functions', () => {
    expect(typeof ViewSwitchButton).toBe('function');
    expect(typeof CreateButton).toBe('function');
    expect(typeof AddIconStyled).toBe('function');
    expect(typeof SearchWrapper).toBe('function');
    expect(typeof BtnText).toBe('function');
    expect(typeof YamlDialogTitle).toBe('function');
    expect(typeof YamlDialogTitleText).toBe('function');
  });
});
