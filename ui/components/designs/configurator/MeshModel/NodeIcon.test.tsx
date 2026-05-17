import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const normalizeStaticImagePath = vi.fn((src?: string) => `NORM:${src}`);

vi.mock('../../../../utils/fallback', () => ({
  normalizeStaticImagePath: (src?: string) => normalizeStaticImagePath(src),
}));

vi.mock('./utils', () => ({
  getComponentMetadata: (metadata: any) => ({
    shape: metadata?.shape || 'circle',
    primaryColor: metadata?.primaryColor || '#000',
    svgColor: metadata?.svgColor || 'colored.svg',
    svgWhite: metadata?.svgWhite || 'white.svg',
  }),
  getStyleOverrides: (raw: any) => (raw ? { 'background-opacity': 0.4, 'border-width': 2 } : {}),
}));

vi.mock('./DesignerDrawerNodeShapes', () => ({
  CircleShape: ({ color }: any) => <svg data-testid="shape-circle" data-color={color} />,
  TriangleShape: ({ color }: any) => <svg data-testid="shape-triangle" data-color={color} />,
  RoundTriangleShape: () => <svg data-testid="shape-rt" />,
  PentagonShape: () => <svg data-testid="shape-pent" />,
  RoundPentagonShape: () => <svg data-testid="shape-rp" />,
  RoundRectangleShape: () => <svg data-testid="shape-rr" />,
  RectangleShape: ({ color }: any) => <svg data-testid="shape-rect" data-color={color} />,
  TallRoundRectangleShape: () => <svg data-testid="shape-trr" />,
  CylinderShape: () => <svg data-testid="shape-cyl" />,
  TagShape: () => <svg data-testid="shape-tag" />,
  RoundTagShape: () => <svg data-testid="shape-rtag" />,
  StarShape: () => <svg data-testid="shape-star" />,
  VeeShape: () => <svg data-testid="shape-vee" />,
  HexagonShape: () => <svg data-testid="shape-hex" />,
  RoundHexagonShape: () => <svg data-testid="shape-rhex" />,
  ConcaveHexagonShape: () => <svg data-testid="shape-chex" />,
  HeptagonShape: () => <svg data-testid="shape-hept" />,
  RoundHeptagonShape: () => <svg data-testid="shape-rhept" />,
  OctagonShape: () => <svg data-testid="shape-octa" />,
  RoundOctagonShape: () => <svg data-testid="shape-rocta" />,
  PolygonShape: () => <svg data-testid="shape-poly" />,
  RhomboidShape: () => <svg data-testid="shape-rhomb" />,
  RightRhomboidShape: () => <svg data-testid="shape-rrhomb" />,
  RoundDiamondShape: () => <svg data-testid="shape-rdiam" />,
  DiamondShape: () => <svg data-testid="shape-diam" />,
  BarrelShape: () => <svg data-testid="shape-barrel" />,
  CutRectangle: () => <svg data-testid="shape-cut" />,
  BottomRoundRectangleShape: () => <svg data-testid="shape-brr" />,
}));

import GetNodeIcon, { ShapeOfSVG } from './NodeIcon';

describe('ShapeOfSVG', () => {
  it.each([
    ['circle', 'shape-circle'],
    ['triangle', 'shape-triangle'],
    ['rectangle', 'shape-rect'],
    ['hexagon', 'shape-hex'],
    ['star', 'shape-star'],
    ['barrel', 'shape-barrel'],
  ])('renders %s -> %s', (shape, testId) => {
    const { getByTestId } = render(<ShapeOfSVG shape={shape} color="#fff" styles={{}} />);
    expect(getByTestId(testId)).toBeInTheDocument();
  });

  it('returns null for unknown shapes', () => {
    const { container } = render(<ShapeOfSVG shape="mystery" color="#fff" styles={{}} />);
    expect(container.firstChild).toBeNull();
  });
});

describe('GetNodeIcon', () => {
  it('renders the kubernetes fallback image when svg metadata is missing', () => {
    const { container } = render(<GetNodeIcon metadata={{}} />);
    const img = container.querySelector('img');
    expect(img).not.toBeNull();
    expect(img).toHaveAttribute('src', '/static/img/kubernetes.svg');
  });

  it('renders an inline image when svgComplete is provided', () => {
    const { container } = render(
      <GetNodeIcon
        metadata={{
          svgComplete: 'foo.svg',
          svgColor: 'c.svg',
          svgWhite: 'w.svg',
        }}
      />,
    );
    expect(container.querySelector('image')).not.toBeNull();
    expect(normalizeStaticImagePath).toHaveBeenCalledWith('foo.svg');
  });

  it('passes through absolute https URLs unchanged through getSvgUrl', () => {
    const { container } = render(
      <GetNodeIcon
        metadata={{
          svgComplete: 'https://example.com/x.svg',
          svgColor: 'c.svg',
          svgWhite: 'w.svg',
        }}
      />,
    );
    const image = container.querySelector('image');
    expect(image?.getAttribute('href')).toBe('https://example.com/x.svg');
  });

  it('renders the colored shape when svg metadata is supplied without svgComplete', () => {
    const { getByTestId } = render(
      <GetNodeIcon
        metadata={{
          shape: 'circle',
          primaryColor: '#fa0',
          svgColor: 'colored.svg',
          svgWhite: 'white.svg',
        }}
      />,
    );
    expect(getByTestId('shape-circle')).toHaveAttribute('data-color', '#fa0');
  });
});
