import WarningIcon from '@material-ui/icons/Warning';
import ErrorIcon from '@material-ui/icons/Error';
import InfoIcon from '@material-ui/icons/Info';

export const eventTypes = {
  0: {
    icon: InfoIcon,
    type: 'success',
  },
  1: {
    icon: WarningIcon,
    type: 'warning',
  },
  2: {
    icon: ErrorIcon,
    type: 'error',
  },
};