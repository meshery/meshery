import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  TriangleShape,
  RoundTriangleShape,
  PentagonShape,
  RoundPentagonShape,
  RoundRectangleShape,
  RectangleShape,
  TallRoundRectangleShape,
  BottomRoundRectangleShape,
  CutRectangle,
  CircleShape,
  CylinderShape,
  TagShape,
  RoundTagShape,
  StarShape,
  VeeShape,
  HexagonShape,
  RoundHexagonShape,
  ConcaveHexagonShape,
  HeptagonShape,
  RoundHeptagonShape,
  OctagonShape,
  RoundOctagonShape,
  PolygonShape,
  RhomboidShape,
  RightRhomboidShape,
  DiamondShape,
  RoundDiamondShape,
  BarrelShape,
  XWingShape,
} from './DesignerDrawerNodeShapes';

const shapes: Record<string, React.ComponentType<{ color: string; styles?: any }>> = {
  TriangleShape,
  RoundTriangleShape,
  PentagonShape,
  RoundPentagonShape,
  RoundRectangleShape,
  RectangleShape,
  TallRoundRectangleShape,
  BottomRoundRectangleShape,
  CutRectangle,
  CircleShape,
  CylinderShape,
  TagShape,
  RoundTagShape,
  StarShape,
  VeeShape,
  HexagonShape,
  RoundHexagonShape,
  ConcaveHexagonShape,
  HeptagonShape,
  RoundHeptagonShape,
  OctagonShape,
  RoundOctagonShape,
  PolygonShape,
  RhomboidShape,
  RightRhomboidShape,
  DiamondShape,
  RoundDiamondShape,
  BarrelShape,
};

describe('DesignerDrawerNodeShapes', () => {
  it.each(Object.entries(shapes))(
    '%s renders a SVG element with provided color',
    (_, Component) => {
      const { container } = render(<Component color="#abc123" />);
      const svg = container.querySelector('svg');
      expect(svg).not.toBeNull();
      expect(container.innerHTML).toContain('#abc123');
    },
  );

  it.each(Object.entries(shapes))('%s spreads additional styles onto its body', (_, Component) => {
    const { container } = render(
      <Component color="#000" styles={{ 'data-testid': 'shape-body' } as any} />,
    );
    // styles get spread on inner shape elements via {...styles}
    expect(container.querySelector('[data-testid="shape-body"]')).not.toBeNull();
  });

  it('XWingShape renders without a styles prop', () => {
    const { container } = render(<XWingShape color="#ff0000" />);
    expect(container.querySelector('svg')).not.toBeNull();
    expect(container.innerHTML).toContain('#ff0000');
  });
});
