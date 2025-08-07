import React from 'react';
import { useState } from 'react';
const HoverIcon = ({ defaultIcon: DefaultIcon, hoverIcon: HoverIconComp, ...props }) => {
  const [hover, setHover] = useState(false);

  return (
    <div
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      style={{
        display: 'inline-block',
        transition: 'transform 0.2s ease',
        padding: '0.8rem',
        margin: '-0.5rem',
        transform: hover ? 'scale(1.1)' : 'scale(1)',
        borderRadius: '50%',
      }}
    >
      {hover ? <HoverIconComp {...props} /> : <DefaultIcon {...props} />}
    </div>
  );
};
export default HoverIcon;
