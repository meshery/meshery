import React from 'react';
import { Button, CustomTooltip, IconButton } from '@sistent/sistent';

export default function TooltipButton({ children, onClick, title, variant, ...props }) {
  return (
    <CustomTooltip title={title} placement="top" interactive>
      <div>
        <Button sx={{ fontSize: '0.875rem' }} variant={variant} onClick={onClick} {...props}>
          {children}
        </Button>
      </div>
    </CustomTooltip>
  );
}

export const TooltipIconButton = ({ children, onClick, title, ...props }) => {
  return (
    <CustomTooltip title={title} placement="top" interactive>
      <div>
        <IconButton onClick={onClick} {...props}>
          {children}
        </IconButton>
      </div>
    </CustomTooltip>
  );
};
