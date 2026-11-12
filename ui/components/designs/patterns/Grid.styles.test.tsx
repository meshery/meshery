import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) =>
      typeof Component === 'string' ? <span {...props} /> : <Component {...props} />;
    Wrapped.displayName = 'StyledGridMock';
    return Wrapped;
  };
  return {
    styled,
    Paper: ({ children }: any) => <div>{children}</div>,
    Typography: ({ children }: any) => <span>{children}</span>,
  };
});

import {
  GridAddIconStyles,
  GridNoContainerStyles,
  GridNoPapperStyles,
  GridNoTextStyles,
  GridPaginationStyles,
} from './Grid.styles';

describe('designs/patterns/Grid.styles', () => {
  it('exports the expected styled wrappers as functions', () => {
    [
      GridNoPapperStyles,
      GridNoContainerStyles,
      GridNoTextStyles,
      GridPaginationStyles,
      GridAddIconStyles,
    ].forEach((wrapper) => expect(typeof wrapper).toBe('function'));
  });
});
