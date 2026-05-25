import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) =>
      typeof Component === 'string' ? <span {...props} /> : <Component {...props} />;
    Wrapped.displayName = 'StyledPerformanceMock';
    return Wrapped;
  };
  return {
    styled,
    Accordion: ({ children }: any) => <div>{children}</div>,
    FormControl: ({ children }: any) => <div>{children}</div>,
    IconButton: ({ children }: any) => <button>{children}</button>,
    Radio: () => <input type="radio" />,
  };
});

vi.mock('@/assets/icons', () => ({
  HelpOutlineOutlined: () => null,
}));

vi.mock('react-big-calendar', () => ({
  Calendar: () => null,
}));

import {
  BottomPart,
  ButtonTextWrapper,
  CalendarComponent,
  CardButton,
  CenterTimer,
  ExpansionPanelComponent,
  FormContainer,
  HelpIcon,
  IconButtonComp,
  PaginationWrapper,
  ProfileContainer,
  RadioButton,
  ResultContainer,
  ResultContainerWrap,
  ViewSwitchBUtton,
} from './style';

describe('performance/style', () => {
  it('exports the expected styled wrappers as functions', () => {
    [
      CardButton,
      BottomPart,
      ResultContainer,
      PaginationWrapper,
      ViewSwitchBUtton,
      ProfileContainer,
      ButtonTextWrapper,
      ResultContainerWrap,
      CalendarComponent,
      IconButtonComp,
      CenterTimer,
      HelpIcon,
      RadioButton,
      ExpansionPanelComponent,
      FormContainer,
    ].forEach((wrapper) => expect(typeof wrapper).toBe('function'));
  });
});
