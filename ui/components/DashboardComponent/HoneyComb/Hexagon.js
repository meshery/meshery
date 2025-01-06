import * as React from 'react';
import { HoneycombContext } from './helpers';

const Hexagon = ({ children, className, style, onclick }) => {
  const { gap } = React.useContext(HoneycombContext);
  return (
    <div
      className={className}
      style={{
        ...style,
        position: 'absolute',
        inset: `${gap}px`,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)',
        pointerEvents: 'auto',
        boxSizing: 'border-box',
      }}
      onClick={onclick}
    >
      {children}
    </div>
  );
};

export default Hexagon;
