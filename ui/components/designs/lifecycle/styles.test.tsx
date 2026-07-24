import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) =>
      typeof Component === 'string' ? <span {...props} /> : <Component {...props} />;
    Wrapped.displayName = 'StyledLifecycleMock';
    return Wrapped;
  };
  return {
    styled,
    Box: ({ children }: any) => <div>{children}</div>,
    List: ({ children }: any) => <ul>{children}</ul>,
    ListItem: ({ children }: any) => <li>{children}</li>,
    ListItemText: ({ children }: any) => <span>{children}</span>,
    ListSubheader: ({ children }: any) => <div>{children}</div>,
    alpha: () => 'rgba(0,0,0,0)',
  };
});

import {
  ComponentValidationListItem,
  DryRunComponentLabel,
  DryRunComponentStyled,
  DryRunErrorContainer,
  DryRunRootListStyled,
  DryRunSignleError,
  ValidatedComponent,
  ValidationErrorListItem,
  ValidationResultsListWrapper,
  ValidationSubHeader,
} from './styles';

describe('designs/lifecycle/styles', () => {
  it('exports the expected styled wrappers as functions', () => {
    [
      ValidationErrorListItem,
      DryRunErrorContainer,
      ComponentValidationListItem,
      ValidatedComponent,
      DryRunComponentStyled,
      DryRunComponentLabel,
      ValidationResultsListWrapper,
      ValidationSubHeader,
      DryRunRootListStyled,
      DryRunSignleError,
    ].forEach((wrapper) => expect(typeof wrapper).toBe('function'));
  });
});
