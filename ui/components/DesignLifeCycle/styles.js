import { Box, List, ListItem, ListItemText, ListSubheader, styled, alpha } from '@layer5/sistent';
import { NOTIFICATIONCOLORS } from '@/themes/index';

export const ValidationErrorListItem = styled(ListItem)(({ theme }) => ({
  gap: '0.5rem',
  backgroundColor: theme.palette.background.card,
  cursor: 'pointer',
  '&:hover': {
    backgroundColor: alpha(NOTIFICATIONCOLORS.WARNING, 0.25),
  },
}));

export const DryRunErrorContainer = styled(ListItem)(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
  cursor: 'pointer',
  '&:hover': {
    backgroundColor: alpha(NOTIFICATIONCOLORS.ERROR_DARK, 0.25),
  },
}));

export const ComponentValidationListItem = styled(ListItem)({
  gap: '0.5rem',
  backgroundColor: NOTIFICATIONCOLORS.WARNING,
  '&:hover': {
    backgroundColor: NOTIFICATIONCOLORS.WARNING,
  },
});

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

export const DryRunComponentLabel = styled(ListItem)({
  backgroundColor: NOTIFICATIONCOLORS.ERROR_DARK,
  gap: '0.5rem',
  color: 'white',
  '&:hover': {
    backgroundColor: NOTIFICATIONCOLORS.ERROR_DARK,
  },
});

export const ValidationResultsListWrapper = styled(List)({
  width: '100%',
  maxHeight: '18rem',
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
