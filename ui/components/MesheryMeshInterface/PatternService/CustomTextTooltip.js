import { CustomTooltip } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

export const CustomTextTooltip = ({ ...props }) => {
  return (
    <UsesSistent>
      <CustomTooltip {...props} />
    </UsesSistent>
  );
};
