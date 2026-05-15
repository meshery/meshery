import { Box, List, ListItem, ListItemText, ListSubheader } from '@sistent/sistent';
import { alpha, styled } from '@/theme';

export const ValidationErrorListItem = styled(ListItem, {
  shouldForwardProp: (prop) => prop !== 'clickable',
})(({ theme, clickable = true }) => ({
  gap: '0.5rem',
  backgroundColor: theme.palette.background.card,
  cursor: clickable ? 'pointer' : 'default',
  '&:hover': {
    backgroundColor: clickable
      ? alpha(theme.palette.warning.main, 0.25)
      : theme.palette.background.card,
  },
}));

export const DryRunErrorContainer = styled(ListItem)(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
  cursor: 'pointer',
  '&:hover': {
    backgroundColor: alpha(theme.palette.error.dark, 0.25),
  },
}));

export const ComponentValidationListItem = styled(ListItem)(({ theme }) => ({
  gap: '0.5rem',
  backgroundColor: theme.palette.warning.main,
  '&:hover': {
    backgroundColor: theme.palette.warning.main,
  },
}));

export const ValidatedComponent = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
  color: theme.palette.text.default,
  margin: '0.6rem 0rem',
}));

export const DryRunComponentStyled = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
  color: theme.palette.text.default,
  marginBlock: '0.5rem',
}));

export const DryRunComponentLabel = styled(ListItem)(({ theme }) => ({
  backgroundColor: theme.palette.error.dark,
  gap: '0.5rem',
  color: 'white',
  '&:hover': {
    backgroundColor: theme.palette.error.dark,
  },
}));

export const ValidationResultsListWrapper = styled(List)({
  width: '100%',
  maxHeight: '18rem',
  overflowY: 'auto',
  marginBottom: '0.5rem',
  display: 'flex',
  flexDirection: 'column',
  gap: '1rem',
});

export const ValidationSubHeader = styled(ListSubheader)({
  marginTop: '1rem',
  display: 'flex',
  padding: 0,
  justifyContent: 'space-between',
  width: '100%',
});

export const DryRunRootListStyled = styled(List)({
  width: '100%',
  position: 'relative',
  marginBottom: '0.5rem',
});

export const DryRunSignleError = styled(ListItemText)(({ theme }) => ({
  paddingInline: theme.spacing(1),
  paddingBlock: theme.spacing(1),
  marginInline: theme.spacing(0.5),
}));
