import { getWebAdress } from '../../../utils/webApis';
import clsx from 'clsx';
import { getComponentMetadata, getStyleOverrides } from './utils';
import {
  CircleShape,
  RectangleShape,
  RoundPentagonShape,
  RoundRectangleShape,
  RoundTriangleShape,
  TallRoundRectangleShape,
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
  PentagonShape,
  TriangleShape,
  CutRectangle,
  PolygonShape,
  RhomboidShape,
  RightRhomboidShape,
  RoundDiamondShape,
  DiamondShape,
  BarrelShape,
  BottomRoundRectangleShape,
} from './DesignerDrawerNodeShapes';

const NODE_ICON_WRAPPER_CLASS = 'node-icon-wrapper-svg';

export function ShapeOfSVG({ color, shape, styles }) {
  switch (shape) {
    case 'circle':
      return <CircleShape color={color} styles={styles} />;
    case 'triangle':
      return <TriangleShape color={color} styles={styles} />;
    case 'round-triangle':
      return <RoundTriangleShape color={color} styles={styles} />;
    case 'pentagon':
      return <PentagonShape color={color} styles={styles} />;
    case 'round-pentagon':
      return <RoundPentagonShape color={color} styles={styles} />;
    case 'round-rectangle':
      return <RoundRectangleShape color={color} styles={styles} />;
    case 'rectangle':
      return <RectangleShape color={color} styles={styles} />;
    case 'tall-round-rectangle':
      return <TallRoundRectangleShape color={color} styles={styles} />;
    case 'bottom-round-rectangle':
      return <BottomRoundRectangleShape color={color} styles={styles} />;
    case 'cut-rectangle':
      return <CutRectangle color={color} styles={styles} />;
    // todos, rn, defaulted by diamond, why not adding default?
    // because the shapes may go unnoticed later in time
    case 'diamond':
      return <DiamondShape color={color} styles={styles} />;
    case 'round-diamond':
      return <RoundDiamondShape color={color} styles={styles} />;
    case 'cylinder':
      return <CylinderShape color={color} styles={styles} />;
    case 'tag':
      return <TagShape color={color} styles={styles} />;
    case 'round-tag':
      return <RoundTagShape color={color} styles={styles} />;
    case 'star':
      return <StarShape color={color} styles={styles} />;
    case 'vee':
      return <VeeShape color={color} styles={styles} />;
    case 'hexagon':
      return <HexagonShape color={color} styles={styles} />;
    case 'round-hexagon':
      return <RoundHexagonShape color={color} styles={styles} />;
    case 'concave-hexagon':
      return <ConcaveHexagonShape color={color} styles={styles} />;
    case 'heptagon':
      return <HeptagonShape color={color} styles={styles} />;
    case 'round-heptagon':
      return <RoundHeptagonShape color={color} styles={styles} />;
    case 'octagon':
      return <OctagonShape color={color} styles={styles} />;
    case 'round-octagon':
      return <RoundOctagonShape color={color} styles={styles} />;
    case 'polygon':
      return <PolygonShape color={color} styles={styles} />;
    case 'rhomboid':
      return <RhomboidShape color={color} styles={styles} />;
    case 'right-rhomboid':
      return <RightRhomboidShape color={color} styles={styles} />;
    case 'barrel':
      return <BarrelShape color={color} styles={styles} />;
    default:
      return null;
  }
}

function getSvgUrl(url) {
  if (!url) {
    return '/static/img/model-categories/dark/other-dark.svg';
  }

  if (url.startsWith('http')) {
    return url;
  }

  return getWebAdress() + '/' + url;
}

function styleProps(metadata) {
  if (metadata.shape === 'no-shape') {
    return {
      height: 40,
      width: 40,
    };
  }

  return {
    height: 25,
    width: 25,
    x: 7.5,
    y: 7.5,
    ...getStyleOverrides(metadata?.styleOverrides),
  };
}

export default function GetNodeIcon({ metadata, otherOverrides = {}, className }) {
  console.log('inside getnodfeicone', metadata);
  if (metadata?.svgComplete) {
    return (
      // wrapping in svg important for onDragStart function
      <svg
        className={clsx(className, NODE_ICON_WRAPPER_CLASS)}
        width="40"
        height="40"
        xmlns="http://www.w3.org/2000/svg"
        {...otherOverrides}
      >
        <image href={getSvgUrl(metadata.svgComplete)} height={40} width={40} {...otherOverrides} />
      </svg>
    );
  }
  /**
   * Converts a set of input styles to SVG-compatible styles.
   * @param {object} inputStyle - The input style object to be converted.
   * @param {string} color - The color to be used for the SVG stroke.
   * @returns {object} An object containing SVG-compatible style properties.
   */
  function convertToSvgStyle(inputStyle, color) {
    const svgStyle = {
      stroke: color, // Use the specified color for the SVG stroke
      fillOpacity: inputStyle['background-opacity'], // Map 'background-opacity' to 'fillOpacity'
      strokeWidth: inputStyle['border-width'], // Map 'border-width' to 'strokeWidth'
      //Add more styles as per requirement...
    };
    return svgStyle;
  }

  const componentMetadata = getComponentMetadata(metadata);
  const styles = styleProps(componentMetadata);
  const styleOverrides = getStyleOverrides(metadata?.styleOverrides);
  const customStyles = convertToSvgStyle(styleOverrides, componentMetadata.primaryColor);
  return (
    <>
      <svg
        width="40"
        height="40"
        className={clsx(className, NODE_ICON_WRAPPER_CLASS)}
        xmlns="http://www.w3.org/2000/svg"
      >
        <ShapeOfSVG
          color={componentMetadata.primaryColor}
          shape={componentMetadata.shape}
          styles={customStyles}
        />
        {styles['background-image'] !== 'none' && componentMetadata.svgColor && (
          <image href={getSvgUrl(componentMetadata.svgWhite)} {...styles} />
        )}
      </svg>
    </>
  );
}
