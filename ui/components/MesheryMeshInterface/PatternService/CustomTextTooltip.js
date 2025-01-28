import { CustomTooltip, styled } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

export const CustomTextTooltip = ({ flag, ...props }) => {
  const StyledTooltip = styled(CustomTooltip)(() => ({
    fontFamily: flag ? 'Qanelas Soft, sans-serif' : 'inherit',
  }));

  return (
    <UsesSistent>
      <StyledTooltip {...props} />
    </UsesSistent>
  );
};
