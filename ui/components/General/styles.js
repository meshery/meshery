import { Card, Paper, Typography, styled } from '@mui/material';
// import theme from '../../styles/themes/theme';
import Button from '@mui/material/Button';

export const SupportTitle = styled(Typography)(() => ({
  fontSize: 42,
  fontWeight: 500,
  textAlign: 'center',
}));

export const SupportModal = styled(Paper)(() => ({
  maxWidth: 900,
  paddingTop: 1,
  paddingBottom: 30,
  textAlign: 'center',
  margin: 'auto',
}));

export const SubmitButton = styled(Button)(() => ({
  //   backgroundColor: theme.palette.btnBg,
  '&:hover': {
    // backgroundColor: theme.palette.btnHover,
  },
}));

export const CardTitle = styled(Typography)(() => ({
  fontSize: 20,
  fontWeight: 500,
  textAlign: 'left',
  //   color: theme.palette.limedSpruce,
  marginLeft: '1rem',
}));

export const CardButton = styled(Button)(() => ({
  //   color: theme.palette.gray,
  '&:hover': {
    // backgroundColor: theme.palette.lightCasper,
  },
}));

export const DescriptionModal = styled('div')(() => ({
  width: 700,
  maxWidth: '100%',
  margin: 'auto',
}));

export const CardBox = styled(Card)(() => ({
  minWidth: 5,
  boxShadow:
    ' 0 2px 2px 0 rgb(0 0 0 / 14%), 0 1px 5px 0 rgb(0 0 0 / 12%), 0 3px 1px -2px rgb(0 0 0 / 20%)',
  //   backgroundColor: theme.palette.white,
  //hover
  '&:hover': {
    // backgroundColor: theme.palette.lightCasper,
  },
}));

export const ResultItem = styled('div')(() => ({
  padding: '0.8rem',
  '&:hover': {
    // backgroundColor: theme.palette.lightCasper,
    opacity: 0.8,
    // zIndex: 1,
    // padding: "5px",
  },
}));
