import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) => <Component {...props} />;
    Wrapped.displayName = 'StyledAppBar';
    return Wrapped;
  };
  return {
    AppBar: ({ children, ...props }: any) => (
      <div data-testid="appbar" {...props}>
        {children}
      </div>
    ),
    styled,
  };
});

import AppBarComponent from './AppBar';

describe('AppBar styled wrapper', () => {
  it('is a function component (styled wrapper) that can be referenced', () => {
    expect(typeof AppBarComponent).toBe('function');
  });
});
