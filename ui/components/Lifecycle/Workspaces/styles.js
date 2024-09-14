import { Typography, Card, styled, Box, Button } from '@material-ui/core';
import theme, { Colors } from '../../../themes/app';
import { Checkbox } from '@layer5/sistent';

/** Workspace card wrapper */
export const CardWrapper = styled(Card)(({ theme }) => ({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  backgroundColor: theme.palette.secondary.elevatedComponents,
  padding: '20px',
  '&:hover': {
    cursor: 'pointer',
  },
}));

export const BulkSelectCheckbox = styled(Checkbox)({
  padding: 0,
  marginRight: '0.5rem',
  height: '28px',
  '& .MuiSvgIcon-root': {
    borderColor: 'white',
  },
  color: 'white',
  '&:hover': {
    color: 'white',
    cursor: 'pointer',
  },
  '&.Mui-checked': {
    color: 'white',
  },
});

export const CardTitle = styled(Typography)({
  fontSize: '1.25rem',
  fontWeight: 800,
  '&:hover': {
    cursor: 'default',
  },
});

export const OrganizationName = styled(Typography)({
  fontSize: '0.9rem',
  display: 'flex',
  alignItems: 'end',
  padding: '0 5px',
  '&:hover': {
    cursor: 'default',
  },
});

export const StyledIconButton = styled('button')({
  background: 'transparent',
  border: 'none',
  '&:hover': {
    cursor: 'default',
  },
});

export const DateLabel = styled(Typography)({
  fontStyle: 'italic',
  fontSize: '12px',
  '&:hover': {
    cursor: 'default',
  },
});

export const EmptyDescription = styled(Typography)({
  fontSize: '0.9rem',
  textAlign: 'left',
  fontStyle: 'italic',
});

export const DescriptionLabel = styled(EmptyDescription)({
  height: 'fit-content',
  fontStyle: 'normal',
  '&:hover': {
    cursor: 'default',
  },
});

export const AllocationButton = styled(Box)({
  background: theme.palette.secondary.focused,
  padding: '10px 10px 1px 10px',
  borderRadius: '4px',
  height: '100%',
  display: 'flex',
  width: '100%',
});

export const AllocationWorkspace = styled(AllocationButton)({
  display: 'flex',
  width: '100%',
  gap: '10px',
  ['@media (min-width : 600px)']: {
    flexDirection: 'column',
    gap: '0',
  },
});

export const PopupButton = styled(Button)(() => ({
  width: '100%',
  borderRadius: '4px',
  background: theme.palette.secondary.white,
  boxShadow: '0px 4px 4px 0px rgba(0, 0, 0, 0.25)',
  display: 'flex',
  flexDirection: 'column',
  marginBottom: '10px',
  color: Colors.charcoal,
  '&:hover': {
    background: theme.palette.secondary.white,
  },
  padding: '15px 10px',
}));

export const TabTitle = styled('p')({
  margin: '0',
  fontSize: '12px',
  fontWeight: '400',
  display: 'flex',
});

export const TabCount = styled('p')({
  margin: '0',
  fontSize: '60px',
  fontWeight: '500',
  lineHeight: 1,
  marginBottom: '5px',
});

export const ViewButton = styled(Button)(() => ({
  width: '100%',
  borderRadius: '4px',
  background: theme.palette.secondary.white,
  boxShadow: '0px 4px 4px 0px rgba(0, 0, 0, 0.25)',
  display: 'flex',
  flexDirection: 'column',
  marginBottom: '10px',
  color: `${Colors.charcoal}30 !important`,
  '&:hover': {
    background: theme.palette.secondary.white,
  },
  padding: '15px 10px',
}));
