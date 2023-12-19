import { styled } from '@mui/system';
import {
  Button,
  Typography,
  TextField,
  Grid,
  Box,
  Card,
  CardContent,
  Checkbox,
  Chip,
} from '@mui/material';
import theme, { Colors } from '../../themes/app';

// const styles = (theme) => ({
// })
export const EditButtonContainer = styled('div')(() => ({
  padding: '0 0.7rem',
  display: 'flex',
  justifyContent: 'flex-end',
}));

export const ProfileTitle = styled(Typography)(() => ({
  padding: '2rem 1rem 0 1rem',
  margin: 'auto',
}));

export const CancelButton = styled(Button)(() => ({
  marginRight: '1rem',
  backgroundColor: '#6c757d',
  '&:hover': {
    backgroundColor: '#5c636a',
  },
}));

export const CreateButton = styled(Button)(() => ({
  backgroundColor: Colors.keppelGreen,
  '&:hover': {
    backgroundColor: Colors.caribbeanGreen,
  },
}));

export const DeleteButton = styled(CreateButton)(() => ({
  backgroundColor: theme.palette.secondary.error,
  '&:hover': {
    backgroundColor: theme.palette.secondary.lightError,
  },
  '@media (max-width: 768px)': {
    minWidth: '50px',
  },
}));

export const SaveButton = styled(Button)(() => ({
  marginRight: '1rem',
  backgroundColor: Colors.keppelGreen,
  '&:hover': {
    backgroundColor: Colors.caribbeanGreen,
  },
}));

export const ConnectButton = styled(Button)(() => ({
  margin: '1rem',
  padding: '.5rem, .5rem',
  bgcolor: 'red',
  color: theme.palette.secondary.white,
  '&:hover': {
    backgroundColor: Colors.caribbeanGreen,
  },
  '&.Mui-disabled': {
    backgroundColor: theme.palette.secondary.btnDisabled,
    color: theme.palette.secondary.white,
  },
}));

export const ButtonContainer = styled('div')(() => ({
  padding: '0.7rem',
  display: 'flex',
  justifyContent: 'flex-end',
}));

export const FormField = styled(TextField)(() => ({
  width: '100%',
  '& .MuiFormLabel-root': {
    fontFamily: "'Qanelas Soft Regular', sans-serif",
  },
  '& .MuiInputBase-input': {
    padding: '0.8rem 0.6rem',
    ['@media (max-width : 899px)']: {
      width: '20rem',
    },
  },
}));

export const GridItemContainer = styled(Grid)(() => ({
  marginBottom: '1.5rem',
  '& .MuiOutlinedInput-input': {
    ['@media (max-width : 899px)']: {
      width: '20rem',
    },
  },
}));

export const FormPaper = styled('div')(() => ({
  margin: 'auto',
  maxWidth: '45rem',
  padding: '30px 20px',
  borderRadius: '15px',
  backgroundColor: theme.palette.secondary.white,
}));

export const OverviewDiv = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  flexWrap: 'wrap',
}));

export const CardContentBox = styled(Box)({
  display: 'flex',
  alignItems: 'center',
});

export const CardStatsBox = styled(Button)({
  display: 'flex',
  alignItems: 'center',
  flexDirection: 'column',
  padding: '5px',
  minWidth: '175px',
  borderRadius: '5px',
  color: `${theme.palette.secondary.white}`,
  background: Colors.keppelGreen,
  textTransform: 'capitalize',
  marginBottom: '10px',
  '&:hover': {
    cursor: 'pointer',
    background: theme.palette.secondary.success,
  },
});

/** Card Styles */
export const CardWrapper = styled(Card)({
  width: '100%',
  display: 'flex',
  flexDirection: 'column',
  padding: '20px',
  '&:hover': {
    cursor: 'pointer',
  },
});

export const IconButton = styled(Button)({
  minWidth: 'fit-content',
  '&.MuiButtonBase-root:hover': {
    bgcolor: 'transparent',
  },
});

export const Statistic = styled(Typography)({
  display: 'flex',
  justifyContent: 'center',
  paddingX: '5px',
  fontWeight: '600',
  fontSize: '24px',
  textAlign: 'center',
});

export const StatisticName = styled(Typography)({
  display: 'flex',
  justifyContent: 'center',
  paddingX: '5px',
  fontWeight: '400',
  fontSize: '16px',
  textAlign: 'center',
});

export const TabCardContent = styled(CardContent)({
  padding: '16px',
  display: 'flex',
  justifyContent: 'center',
  flexDirection: 'row',
});

export const TabIconBox = styled(Box)({
  display: 'flex',
  alignItems: 'center',
  border: `1px solid ${Colors.keppelGreen}`,
  borderRadius: '5px 0 0 5px',
  padding: '8px 16px',
  borderRight: '0',
  '&:hover': {
    cursor: 'default',
  },
});

export const TabNameBox = styled(Box)({
  display: 'flex',
  alignItems: 'center',
  flexDirection: 'column',
  background: Colors.keppelGreen,
  color: theme.palette.secondary.white,
  p: '8px',
  border: `1px solid ${Colors.keppelGreen}`,
  borderRadius: '0 5px 5px 0',
  padding: '8px 16px',
});

export const TabTitle = styled('p')({
  margin: '0',
  fontSize: '12px',
  fontWeight: '400',
  display: 'flex',
});

export const TabCount = styled('p')({
  margin: '0',
  fontSize: '30px',
  fontWeight: '500',
  lineHeight: 1,
  marginBottom: '5px',
});

export const AllocationButton = styled(Box)({
  background: Colors.keppelGreen,
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
  color: theme.palette.secondary.penColorPrimary,
  '&:hover': {
    background: theme.palette.secondary.white,
  },
  padding: '15px 10px',
}));

export const ViewButton = styled(Button)(({ disabled }) => ({
  width: '100%',
  borderRadius: '4px',
  background: disabled ? theme.palette.cultured : theme.palette.secondary.white,
  boxShadow: '0px 4px 4px 0px rgba(0, 0, 0, 0.25)',
  display: 'flex',
  flexDirection: 'column',
  marginBottom: '10px',
  color: theme.palette.secondary.penColorPrimary,
  '&:hover': {
    background: theme.palette.secondary.white,
  },
  padding: '15px 10px',
}));

export const Record = styled(Grid)({
  borderBottom: `1px solid ${theme.palette.secondary.modalTabs}60`,
  display: 'flex',
  flexDirection: 'row',
  padding: '5px 0',
});

export const BulkSelectCheckbox = styled(Checkbox)({
  padding: 0,
  marginRight: '0.5rem',
  color: theme.palette.secondary.white,
  '&:hover': {
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

export const Status = styled('div')({
  padding: '5px 20px',
  width: 'fit-content',
  border: '1px solid transparent',
  fontSize: '12px',
  '&:hover': {
    cursor: 'default',
  },
});

export const StyledChip = styled(Chip)({
  padding: '5px 6px !important',
  color: theme.palette.secondary.text,
  fontSize: '14px',
  textTransform: 'uppercase',
  fontWeight: 400,
  height: 'unset',
  borderRadius: '100px',
  border: `0.5px solid ${theme.palette.secondary.default}`,
  background: theme.palette.secondary.white,
  maxWidth: '230px',
  '.MuiChip-avatar': {
    margin: '0',
  },
  '&:hover': {
    cursor: 'default',
  },
});

export const CreateButtonWrapper = styled('div')({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
});

export const EditButton = styled(Button)(() => ({
  backgroundColor: Colors.keppelGreen,
  '&:hover': {
    backgroundColor: Colors.caribbeanGreen,
  },
  '@media (max-width: 768px)': {
    minWidth: '50px',
  },
}));

export const TextButton = styled('span')(() => ({
  marginLeft: '0.5rem',
  display: 'block',
  '@media (max-width: 853px)': {
    display: 'none',
  },
}));
