import { CONNECTION_STATES } from '../../utils/Enum';
import { notificationColors } from '../../themes';
import { Box, Button, Chip, Grid, MenuItem, Select, styled, Tab, Tabs } from '@layer5/sistent';
import { StepConnector, StepLabel, Stepper, TableContainer } from '@mui/material';

export const CreateButton = styled(Grid)({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
});

export const InnerTableContainer = styled(TableContainer)(({ theme }) => ({
  background: theme.palette.background.card,
  margin: '10px 10px 10px 13px',
  borderLeft: `9px solid ${theme.palette.background.default} !important`,
  borderRadius: '10px 0 0 10px',
  width: 'calc(100% - 23px)',
  border: 'none',
  overflowX: 'hidden',
}));

export const ActionListItem = styled(Box)(({ theme }) => ({
  display: 'flex',
  gridGap: '0.5rem',
  alignItems: 'center',
  justifyContent: 'space-around',
  width: '100%',
  backgroundColor: theme.palette.background.card,
  padding: '10px',
}));

export const ActionButton = styled(Button)({
  width: '100%',
  justifyContent: 'flex-start',
});

export const ConnectionTab = styled(Tab)(({ theme }) => ({
  minWidth: 40,
  paddingLeft: 0,
  paddingRight: 0,
  '&.Mui-selected': {
    color: theme.palette.mode === 'dark' ? '#00B39F' : theme.palette.primary,
  },
}));

export const ConnectionTabs = styled(Tabs)(({ theme }) => ({
  height: '55px',
  '& .MuiTabs-indicator': {
    backgroundColor: theme.palette.mode === 'dark' ? '#00B39F' : theme.palette.primary,
  },
}));

export const ConnectionStyledSelect = styled(Select)({
  '& .MuiSelect-select': {
    padding: '0 !important',
  },
  '& .MuiMenuItem-root': {
    padding: '0 !important',
  },
});

export const ConnectionStyledMenuItem = styled(MenuItem)({
  padding: 0,
  '&.MuiMenuItem-root': {
    padding: 0,
  },
  '& .MuiButtonBase-root': {
    padding: 0,
  },
});

export const ConnectionIconText = styled('div')({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
});

export const ChipWrapper = styled(Chip)({
  width: '13rem',
  maxWidth: '13rem',
  minWidth: '9rem',
  textAlign: 'left',
  cursor: 'pointer',
  '& .MuiChip-label': {
    flexGrow: 1,
  },
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  border: '1px solid rgba(255, 255, 255, 0.23)',
});

const baseChipStyles = {
  minWidth: '142px !important',
  maxWidth: 'max-content !important',
  display: 'flex !important',
  justifyContent: 'flex-start !important',
  borderRadius: '3px !important',
  padding: '6px 8px',
  '& .MuiChip-label': {
    paddingTop: '3px',
    fontWeight: '400',
  },
  '& .MuiSvgIcon-root': {
    marginLeft: '0px !important',
  },
  '&:hover': {
    boxShadow: '0px 1px 2px 0px rgba(0, 0, 0, 0.25)',
  },
  textTransform: 'capitalize',
};

export const DiscoveredChip = styled(Chip)(() => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: notificationColors.info,
  },
  background: `${notificationColors.info}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${notificationColors.info} !important`,
  },
}));

export const DeletedChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: theme.palette.background.error.default,
  },
  background: `${theme.palette.background.error.default}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.background.error.default} !important`,
  },
}));

export const RegisteredChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: theme.palette.mode === 'dark' ? '#86B2C6' : '#477E96',
  },
  background: `${theme.palette.mode === 'dark' ? '#86B2C6' : '#477E96'}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.mode === 'dark' ? '#86B2C6' : '#477E96'} !important`,
  },
}));

export const ConnectedChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: theme.palette.background.success.default,
  },
  background: `${theme.palette.background.success.default}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.background.success.default} !important`,
  },
}));

export const IgnoredChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: theme.palette.mode === 'dark' ? '#9FAFB6' : '#51636B',
  },
  background: `${theme.palette.mode === 'dark' ? '#9FAFB6' : '#51636B'}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.mode === 'dark' ? '#9FAFB6' : '#51636B'} !important`,
  },
}));

export const DisconnectedChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: theme.palette.background.warning.default,
  },
  background: `${theme.palette.background.warning.default}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.background.warning.default} !important`,
  },
}));

export const NotFoundChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: theme.palette.text.disabled,
  },
  background: `${theme.palette.background.disabled}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.icon.default} !important`,
  },
}));

export const MaintainanceChip = styled(Chip)(() => ({
  ...baseChipStyles,
  '& .MuiChip-label': {
    color: notificationColors.lightwarning,
  },
  background: `${notificationColors.lightwarning}30 !important`,
  '& .MuiSvgIcon-root': {
    color: `${notificationColors.lightwarning} !important`,
  },
}));

export const ColumnWrapper = styled('div')(({ theme }) => ({
  margin: theme.spacing(2),
  padding: theme.spacing(2),
  background: `${theme.palette.background.secondary}10`,
}));

export const OperationButton = styled(Grid)(({ theme }) => ({
  [theme?.breakpoints?.down(1180)]: {
    marginRight: '25px',
  },
}));

export const ContentContainer = styled(Grid)(({ theme }) => ({
  [theme?.breakpoints?.down(1050)]: {
    flexDirection: 'column',
  },
  flexWrap: 'noWrap',
}));

export const StepperContainer = styled(Stepper)(() => ({
  width: '80%',
  marginTop: '2rem',
  ['@media (max-width:780px)']: {
    width: 'auto',
    marginTop: '1rem',
  },
}));

export const CustomLabelStyle = styled(StepLabel)(() => ({
  fontSize: '0.875rem',
  ['@media (max-width:500px)']: {
    fontSize: '0.7rem',
  },
}));

export const ColorlibConnector = styled(StepConnector)(({ theme }) => ({
  alternativeLabel: { top: 22 },
  active: { '& $line': { background: '#00B39F', transition: 'all 1s ease-in' } },
  completed: { '& $line': { background: '#00B39F', transition: 'all 1s ease-in' } },
  line: {
    height: 3,
    border: 0,
    backgroundColor: theme.palette.mode === 'dark' ? '#303030' : '#eaeaf0',
    borderRadius: 1,
    transition: 'all 0.5s ease-out ',
  },
}));

export const CONNECTION_STATE_COLORS = {
  [CONNECTION_STATES.CONNECTED]: '#00B39F',
  [CONNECTION_STATES.REGISTERED]: '#00B39F',
  [CONNECTION_STATES.DISCOVERED]: '#FFC107',
  [CONNECTION_STATES.IGNORED]: '#FFC107',
  [CONNECTION_STATES.DELETED]: '#FF1744',
  [CONNECTION_STATES.MAINTENANCE]: '#FFC107',
  [CONNECTION_STATES.DISCONNECTED]: '#FFC107',
  [CONNECTION_STATES.NOTFOUND]: '#FFC107',
};
