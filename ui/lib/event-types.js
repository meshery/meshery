import WarningIcon from '@mui/icons-material/Warning';
import ErrorIcon from '@mui/icons-material/Error';
import InfoIcon from '@mui/icons-material/Info';

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