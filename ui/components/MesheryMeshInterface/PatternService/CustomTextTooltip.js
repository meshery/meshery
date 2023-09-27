import { withStyles } from '@material-ui/styles';
import { Tooltip } from '@material-ui/core';
import { ziCalc } from '../../../utils/zIndex';

export const CustomTextTooltip = ({ backgroundColor, flag, ...props }) => {
  const CustomTooltip = withStyles(() => ({
    tooltip: {
      backgroundColor: backgroundColor,
      color: '#fff',
      opacity: '100%',
      fontSize: '0.75rem',
      fontFamily: flag ? 'Qanelas Soft, sans-serif' : 'inherit',
      borderRadius: '0.9375rem',
      padding: '0.9rem',
      zIndex: ziCalc(11),
    },
    popper: {
      zIndex: `${ziCalc(5)} !important`,
    },
  }))(Tooltip);

  return <CustomTooltip {...props} />;
};
