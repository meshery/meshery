import { Box, styled } from '@sistent/sistent';

export const FallbackWrapper = styled(Box)(({ theme }) => ({
  margin: '2rem',
  [theme.breakpoints.down('sm')]: {
    margin: '1rem',
    paddingBottom: theme.spacing(3),
  },
}));

export const ErrorPageActions = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexWrap: 'wrap',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  marginTop: theme.spacing(1.5),
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column',
    alignItems: 'stretch',
    '& > *': {
      width: '100%',
      margin: 0,
    },
  },
}));
