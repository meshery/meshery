import React from 'react';
import { styled } from '@sistent/sistent';
import { CardStyle } from '@/assets/styles/general/tool.styles';

const CountSpan = styled('span')(() => ({
  fontSize: '1rem',
  marginLeft: '4px',
}));

const TabCard = ({
  label,
  count,
  active,
  onClick,
}: {
  label: string;
  count: number;
  active: boolean;
  onClick: () => void;
}) => {
  return (
    <CardStyle isSelected={active} elevation={3} onClick={onClick}>
      <CountSpan>{`(${count?.toLocaleString() || 0})`}</CountSpan>
      {label}
    </CardStyle>
  );
};

export default TabCard;
