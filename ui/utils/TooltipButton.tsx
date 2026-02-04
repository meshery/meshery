import React from 'react';
import { Button } from '@sistent/sistent';
import { CustomTooltip, IconButton } from '@sistent/sistent';

export default function TooltipButton({ children, onClick, title, variant, ...props }) {
  return (
    <CustomTooltip title={title} placement="top" interactive>
      <span>
        <Button sx={{ fontSize: '0.875rem' }} variant={variant} onClick={onClick} {...props}>
          {children}
        </Button>
      </span>
    </CustomTooltip>
  );
}

export const TooltipIconButton = ({ children, onClick, title, ...props }) => {
  return (
    <CustomTooltip title={title} placement="top" interactive>
      <span>
        <IconButton onClick={onClick} {...props}>
          {children}
        </IconButton>
      </span>
    </CustomTooltip>
  );
};
