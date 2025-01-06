import * as React from 'react';

const Hexagon = ({ children, className, style, onClick }) => {
  return (
    <div
      className={className}
      style={{
        ...style,
        position: 'absolute',
        inset: `3px`,
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)',
        pointerEvents: 'auto',
        boxSizing: 'border-box',
      }}
      onClick={onClick}
    >
      {children}
    </div>
  );
};

export default Hexagon;
