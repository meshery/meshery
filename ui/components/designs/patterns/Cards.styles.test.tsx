import React from 'react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) =>
      typeof Component === 'string' ? <span {...props} /> : <Component {...props} />;
    Wrapped.displayName = 'StyledCardMock';
    return Wrapped;
  };
  return {
    styled,
    Grid: ({ children }: any) => <div>{children}</div>,
  };
});

import {
  BottomContainer,
  CardBackGrid,
  CardHeaderRight,
  CardImg,
  CardNoContainer,
  CardNoPaper,
  CardNoText,
  CardPagination,
  CatalogCardButtons,
  GridBtnText,
  GridCloneBtnText,
  PatternAddIcon,
  StyledCodeMirrorWrapper,
  UpdateDeleteButtons,
  YamlDialogTitleGrid,
} from './Cards.styles';

describe('designs/patterns/Cards.styles', () => {
  it('exports the expected styled wrappers as functions', () => {
    [
      StyledCodeMirrorWrapper,
      CardBackGrid,
      YamlDialogTitleGrid,
      CardHeaderRight,
      GridBtnText,
      GridCloneBtnText,
      CardImg,
      CardNoPaper,
      CardNoText,
      CardNoContainer,
      CardPagination,
      PatternAddIcon,
      UpdateDeleteButtons,
      BottomContainer,
      CatalogCardButtons,
    ].forEach((wrapper) => {
      expect(typeof wrapper).toBe('function');
    });
  });
});
