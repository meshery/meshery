/* eslint-disable react/forbid-dom-props */
import React from 'react';
import { CardStyle } from '@/assets/styles/general/tool.styles';

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
      <span
        style={{
          fontSize: '1rem',
          marginLeft: '4px',
        }}
      >
        {`(${count?.toLocaleString() || 0})`}
      </span>
      {label}
    </CardStyle>
  );
};

export default TabCard;
