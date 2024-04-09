import { DialogTitle, Typography } from '@layer5/sistent';

import { styled } from '@mui/material';
import theme from '../../themes/app';

export const StyledDialog = styled(DialogTitle)(() => ({
  backgroundColor: theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : '#396679',
  // color: theme.palette.primary.contrastText,
  textAlign: 'center',
  minWidth: 400,
  padding: '10px',
  color: '#fff',
  display: 'flex',
}));

export const YamlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));
