import React from 'react';
import { CustomTooltip } from '@sistent/sistent';

export const CustomTextTooltip = ({ ...props }: React.ComponentProps<typeof CustomTooltip>) => {
  return <CustomTooltip {...props} />;
};
