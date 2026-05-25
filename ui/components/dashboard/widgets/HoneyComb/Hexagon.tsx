import * as React from 'react';
import { HexagonWrapper } from '../../style';

type HexagonProps = {
  children?: React.ReactNode;
  className?: string;
  style?: React.CSSProperties;
  onClick?: React.MouseEventHandler<HTMLDivElement>;
};

const Hexagon = ({ children, className, style, onClick }: HexagonProps) => {
  return (
    <HexagonWrapper className={className} style={style} onClick={onClick}>
      {children}
    </HexagonWrapper>
  );
};

export default Hexagon;
