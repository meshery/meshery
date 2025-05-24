import * as React from 'react';
import { HexagonWrapper } from '../../style';

const Hexagon = ({ children, className, style, onClick }) => {
  return (
    <HexagonWrapper className={className} style={style} onClick={onClick}>
      {children}
    </HexagonWrapper>
  );
};

export default Hexagon;
