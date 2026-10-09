import React from 'react';
import { CustomTooltip, IconButton } from '@sistent/sistent';
import type { TooltipIconProps } from './Filters.types';

function TooltipIcon({
  children,
  onClick,
  title,
  disabled,
  permissionKey,
  permissionAction,
}: TooltipIconProps) {
  return (
    <>
      <CustomTooltip title={title} placement="top" interactive>
        <div>
          <IconButton
            onClick={onClick}
            disabled={disabled}
            permissionKey={permissionKey}
            permissionAction={permissionAction}
          >
            {children}
          </IconButton>
        </div>
      </CustomTooltip>
    </>
  );
}

export default TooltipIcon;
