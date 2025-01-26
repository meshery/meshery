import { Checkbox, Typography, Button, useTheme, Card, styled, Box } from '@layer5/sistent';

/** Bulk action bar styles */
export const BulkActionWrapper = styled('div')({
  width: '100%',
  padding: '0.8rem',
  justifyContent: 'space-between',
  marginTop: '0.18rem',
  marginBottom: '1rem',
  borderRadius: '.25rem',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
});

/** Card Styles */

export const CardWrapper = styled(Card)(() => {
  const theme = useTheme();
  return {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '20px',
    backgroundColor: theme.palette.background?.default,
    '&:hover': {
      cursor: 'pointer',
    },
  };
});

export const Statistic = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  padding: '0 5px',
  fontWeight: 600,
  fontSize: '24px',
  textAlign: 'center',
});

export const StatisticName = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  padding: '0 5px',
  fontWeight: 400,
  fontSize: '16px',
  textAlign: 'center',
});

export const TabCardContent = styled('div')({
  padding: '16px',
  display: 'flex',
  justifyContent: 'center',
  flexDirection: 'row',
});

export const TabIconBox = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  border: `1px solid ${theme.palette.background?.brand?.default}`,
  borderRadius: '5px 0 0 5px',
  padding: '8px 16px',
  borderRight: 0,
  '&:hover': {
    cursor: 'default',
  },
}));

export const TabNameBox = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  flexDirection: 'column',
  background: theme.palette.background?.brand?.default,
  color: theme.palette.mode === 'dark' ? theme.palette.text?.inverse : theme.palette.text?.default,
  padding: '8px 16px',
  border: `1px solid ${theme.palette.background?.brand?.default}`,
  borderRadius: '0 5px 5px 0',
}));

export const TabTitle = styled('p')({
  margin: 0,
  fontSize: '12px',
  fontWeight: 400,
  display: 'flex',
});

export const TabCount = styled('p')({
  margin: 0,
  fontSize: '60px',
  fontWeight: 500,
  lineHeight: 1,
  marginBottom: '5px',
});

export const AllocationButton = styled(Box)(({ theme }) => ({
  padding: '10px 10px 1px 10px',
  borderRadius: '4px',
  height: '100%',
  display: 'flex',
  width: '100%',
  background: theme.palette.background?.brand?.default,
}));

export const AllocationWorkspace = styled('div')({
  display: 'flex',
  width: '100%',
  gap: '10px',
  '@media (min-width: 600px)': {
    flexDirection: 'column',
    gap: 0,
  },
});

export const PopupButton = styled(Button)(({ theme }) => ({
  width: '100%',
  borderRadius: '4px',
  boxShadow: '0px 4px 4px rgba(0, 0, 0, 0.25)',
  marginBottom: '10px',
  padding: '20px 10px',
  color: theme.palette.text?.default,
  backgroundColor: theme.palette.background.constant.table,
  margin: '0px 0px 10px',
  '&:hover': {
    boxShadow: 'none',
    backgroundColor: theme.palette.background.constant.table,
  },
}));

export const Record = styled('div')(() => ({
  // borderBottom: `1px solid ${theme.palette.secondary.modalTabs}60`,
  display: 'flex',
  flexDirection: 'row',
  padding: '5px 0',
}));

export const BulkSelectCheckbox = styled(Checkbox)({
  padding: 0,
  marginRight: '0.5rem',
  color: 'white',
  '&:hover': {
    color: 'white',
    cursor: 'pointer',
  },
  '& .Mui-checked': {
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

export const OrganizationName = styled('p')({
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

export const EmptyDescription = styled('p')({
  fontSize: '0.9rem',
  textAlign: 'left',
  fontStyle: 'italic',
});

export const DescriptionLabel = styled('p')({
  height: 'fit-content',
  fontStyle: 'italic',
  '&:hover': {
    cursor: 'default',
  },
});

export const Name = styled(Typography)({
  height: 'fit-content',
  textAlign: 'left',
  fontWeight: 'bold',
});

export const Status = styled('div')({
  padding: '5px 20px',
  width: 'fit-content',
  border: '1px solid transparent',
  fontSize: '12px',
  '&:hover': {
    cursor: 'default',
  },
});

export const StyledChip = styled('div')(({ theme }) => ({
  padding: '5px 6px',
  color: theme.palette.secondary.text,
  fontSize: '14px',
  textTransform: 'uppercase',
  fontWeight: 400,
  height: 'unset',
  borderRadius: '100px',
  border: `0.5px solid ${theme.palette.secondary.default}`,
  background: theme.palette.secondary.white,
  maxWidth: '230px',
  '& .MuiChip-avatar': {
    margin: 0,
  },
  '&:hover': {
    cursor: 'default',
  },
}));

export const CreateButtonWrapper = styled('div')({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
});

export const EditButton = styled('button')(({ theme }) => ({
  backgroundColor: theme.palette.background?.brand?.default,
  '&:hover': {
    backgroundColor: theme.palette.background?.brand?.default,
  },
  '@media (max-width: 768px)': {
    minWidth: '50px',
  },
}));

export const TextButton = styled('button')({
  marginLeft: '0.5rem',
  display: 'block',
  '@media (max-width: 853px)': {
    display: 'none',
  },
});

export const IconButton = styled('button')({
  minWidth: 'fit-content',
  '&:hover': {
    background: 'transparent',
  },
  padding: 0,
  background: 'none',
  border: 'none',
});
