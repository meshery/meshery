import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  // FlipCard calls `styled('div')` immediately at module load, which returns a
  // component factory; the factory is then invoked with style args. Our mock
  // simply ignores the styled config and returns a plain wrapper component.
  const styled = (_tag: any) => () => {
    const Styled = ({ children, onClick, sx, ...props }: any) => (
      <div data-sx={JSON.stringify(sx || {})} onClick={onClick} {...props}>
        {children}
      </div>
    );
    Styled.displayName = 'StyledSistent';
    return Styled;
  };
  return { styled };
});

import FlipCard from './index';

describe('FlipCard', () => {
  it('renders both front and back components inside the card', () => {
    render(
      <FlipCard
        frontComponents={<div data-testid="front">FRONT</div>}
        backComponents={<div data-testid="back">BACK</div>}
      />,
    );

    expect(screen.getByTestId('front')).toBeInTheDocument();
    expect(screen.getByTestId('back')).toBeInTheDocument();
  });

  const findInner = (container: HTMLElement) =>
    Array.from(container.querySelectorAll<HTMLElement>('div[data-sx]')).find((el) =>
      (el.getAttribute('data-sx') || '').includes('rotateY'),
    );

  it('toggles flipped state on click when flipping is enabled', () => {
    const { container } = render(
      <FlipCard frontComponents={<div>FRONT</div>} backComponents={<div>BACK</div>} />,
    );

    const inner = findInner(container) as HTMLElement;
    expect(inner).toBeTruthy();
    expect(inner.getAttribute('data-sx')).toContain('rotateY(0deg)');

    fireEvent.click(inner);
    const inner2 = findInner(container) as HTMLElement;
    expect(inner2.getAttribute('data-sx')).toContain('rotateY(180deg)');

    fireEvent.click(inner2);
    const inner3 = findInner(container) as HTMLElement;
    expect(inner3.getAttribute('data-sx')).toContain('rotateY(0deg)');
  });

  it('does not toggle the flipped state when disableFlip is true', () => {
    const { container } = render(
      <FlipCard frontComponents={<div>FRONT</div>} backComponents={<div>BACK</div>} disableFlip />,
    );

    const inner = findInner(container) as HTMLElement;
    expect(inner.getAttribute('data-sx')).toContain('rotateY(0deg)');
    fireEvent.click(inner);
    const inner2 = findInner(container) as HTMLElement;
    expect(inner2.getAttribute('data-sx')).toContain('rotateY(0deg)');
  });
});
