import { styled } from '@material-ui/core';

export const ErrorMain = styled('main')(({ theme }) => ({
  background: theme.palette.aliceBlue,
  padding: '4rem 8rem',
  minHeight: '100vh',
  ['@media (max-width:680px)']: {
    padding: '4rem 2rem',
  },
}));

export const ErrorContainer = styled('div')(({ theme }) => ({
  background: theme.palette.white,
  boxShadow:
    '0px 2px 1px -1px rgb(0 0 0 / 20%), 0px 1px 1px 0px rgb(0 0 0 / 14%), 0px 1px 3px 0px rgb(0 0 0 / 12%)',
  borderRadius: '8px',
  padding: '3rem',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
}));

export const ErrorComponent = styled('div')(() => ({
  paddingTop: '2rem',
  width: '100%',
}));

export const ErrorContentContainer = styled('div')(() => ({
  textAlign: 'center',
  backgroundColor: '#fafafa',
  margin: '2rem',
  padding: '20px',
  borderRadius: '10px',
  boxShadow:
    ' 0 2px 2px 0 rgb(0 0 0 / 14%), 0 1px 5px 0 rgb(0 0 0 / 12%), 0 3px 1px -2px rgb(0 0 0 / 20%)',
}));

export const ErrorLink = styled('a')(({ theme }) => ({
  color: theme.palette.secondary.focused,
  textDecoration: 'none',
}));

export const ErrorMsg = styled('errormsg')(() => ({
  fontWeight: '600',
}));
