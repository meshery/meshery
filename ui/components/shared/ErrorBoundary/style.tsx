import { Box, styled } from '@sistent/sistent';

export const FallbackWrapper = styled(Box)(({ theme }) => ({
  margin: '2rem',
  [theme.breakpoints.down('sm')]: {
    margin: '1rem',
    paddingBottom: theme.spacing(3),
  },
}));

/** Always-visible button label for the error page (no mobile display:none). */
export const ErrorPageButtonLabel = styled('span')({
  marginLeft: '0.5rem',
  display: 'block',
});

export const ErrorPageActions = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  flexWrap: 'wrap',
  alignItems: 'center',
  gap: theme.spacing(1.5),
  marginTop: theme.spacing(1.5),
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column',
    alignItems: 'stretch',
    width: '100%',
    '& > *': {
      width: '100%',
      margin: 0,
    },
  },
}));
