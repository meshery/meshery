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

export const FlipCardWrapper = styled('div')(({ theme }) => ({
  height: '100%',
  backgroundColor: 'transparent',
  perspective: theme.spacing(125),
}));

export const InnerCard = styled('div')(({ theme }) => ({
  padding: theme.spacing(2),
  borderRadius: theme.spacing(1),
  transformStyle: 'preserve-3d',
  boxShadow: '0 4px 8px 0 rgba(0,0,0,0.2)',
  backgroundColor: theme.palette.mode === 'dark' ? '#202020' : '#fff',
  cursor: 'pointer',
}));
