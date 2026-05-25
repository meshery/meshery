import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) => <Component {...props} />;
    Wrapped.displayName = 'StyledStyleMock';
    return Wrapped;
  };
  return {
    styled,
    DialogTitle: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
  };
});

import { FlipCardWrapper, InnerCard, StyledDialog, YamlDialogTitleText } from './style';

describe('designs/patterns/style', () => {
  it('exports the expected styled wrappers as functions', () => {
    expect(typeof FlipCardWrapper).toBe('function');
    expect(typeof InnerCard).toBe('function');
    expect(typeof StyledDialog).toBe('function');
    expect(typeof YamlDialogTitleText).toBe('function');
  });
});
