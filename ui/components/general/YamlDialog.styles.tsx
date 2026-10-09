import { DialogTitle, Typography, styled } from '@sistent/sistent';

export const StyledDialog = styled(DialogTitle)(({ theme }) => ({
  backgroundColor: theme.palette.mode === 'dark' ? theme.palette.background.card : '#396679',
  textAlign: 'center',
  minWidth: 400,
  padding: '10px',
  color: '#fff',
  display: 'flex',
}));

export const YamlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));
