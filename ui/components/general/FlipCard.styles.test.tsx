import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const Wrapped = (props: any) =>
      typeof Component === 'string' ? <div {...props} /> : <Component {...props} />;
    Wrapped.displayName = 'StyledMock';
    return Wrapped;
  };
  return {
    styled,
  };
});

import { FlipCardWrapper, InnerCard, CardFaceFront, CardFaceBack } from './FlipCard.styles';

describe('general/FlipCard.styles', () => {
  it('exports the expected styled components as functions', () => {
    [FlipCardWrapper, InnerCard, CardFaceFront, CardFaceBack].forEach((component) => {
      expect(typeof component).toBe('function');
    });
  });
});
