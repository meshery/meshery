import React from 'react';
import { CustomTooltip } from '@sistent/sistent';

type CustomTextTooltipProps = React.ComponentProps<typeof CustomTooltip> & {
  flag?: unknown;
  bgColor?: string;
};

export const CustomTextTooltip = ({ ...props }: CustomTextTooltipProps) => {
  return <CustomTooltip {...props} />;
};
