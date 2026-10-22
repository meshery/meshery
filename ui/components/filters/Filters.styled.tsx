import { Box, Button, DialogTitle, Typography, styled } from '@sistent/sistent';

export const CreateButton = styled(Button)(() => ({
  width: 'fit-content',
  alignSelf: 'flex-start',
  placeSelf: 'center',
}));

export const ViewSwitchButton = styled('div')(() => ({
  justifySelf: 'flex-end',
  paddingLeft: '1rem',
}));

export const YmlDialogTitle = styled(DialogTitle)(() => ({
  display: 'flex',
  alignItems: 'center',
}));

export const YmlDialogTitleText = styled(Typography)(() => ({
  flexGrow: 1,
}));

export const BtnText = styled('span')(({ theme }) => ({
  display: 'block',
  [theme.breakpoints.down('700')]: {
    display: 'none',
  },
}));

export const ActionsBox = styled(Box)(() => ({
  display: 'flex',
}));
