import { CONNECTION_STATES } from '../../utils/Enum';
import { alpha, styled } from '@/theme';
import {
  Box,
  Button,
  Chip,
  Grid,
  MenuItem,
  Select,
  Tab,
  Tabs,
  StepConnector,
  StepLabel,
  Stepper,
  TableContainer,
  Grid2,
} from '@sistent/sistent';

export const CreateButton = styled(Grid)({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  whiteSpace: 'nowrap',
});

export const InnerTableContainer = styled(TableContainer)(({ theme }) => ({
  background: theme.palette.background.card,
  borderLeft: `9px solid ${theme.palette.background.default} !important`,
  borderRadius: '10px 0 0 10px',
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
    color: theme.palette.mode === 'dark' ? theme.palette.primary.main : theme.palette.primary.main,
  },
}));

export const ConnectionTabs = styled(Tabs)(({ theme }) => ({
  height: '55px',
  '& .MuiTabs-indicator': {
    backgroundColor: theme.palette.primary.main,
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

export const ChipWrapper = styled(Chip)(({ theme }) => ({
  width: '13rem',
  textAlign: 'left',
  cursor: 'pointer',
  '& .MuiChip-label': {
    flexGrow: 1,
  },
  overflow: 'hidden',
  whiteSpace: 'nowrap',
  textOverflow: 'ellipsis',
  border: `1px solid ${alpha(theme.palette.common.white, 0.23)}`,
  textTransform: 'lowercase',
  color: theme.palette.text.primary,
}));

const baseChipStyles = (theme) => ({
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
});

export const DiscoveredChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.info.main,
  },
  background: `${alpha(theme.palette.info.main, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.info.main} !important`,
  },
}));

export const DeletedChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.background.error.default,
  },
  background: `${alpha(theme.palette.background.error.default, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.background.error.default} !important`,
  },
}));

export const RegisteredChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color:
      theme.palette.mode === 'dark' ? theme.palette.secondary.light : theme.palette.secondary.dark,
  },
  background: `${alpha(theme.palette.mode === 'dark' ? theme.palette.secondary.light : theme.palette.secondary.dark, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.mode === 'dark' ? theme.palette.secondary.light : theme.palette.secondary.dark} !important`,
  },
}));

export const ConnectedChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.background.success.default,
  },
  background: `${alpha(theme.palette.background.success.default, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.background.success.default} !important`,
  },
}));

export const IgnoredChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.text.secondary,
  },
  background: `${alpha(theme.palette.text.secondary, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.text.secondary} !important`,
  },
}));

export const DisconnectedChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.background.warning.default,
  },
  background: `${alpha(theme.palette.background.warning.default, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.background.warning.default} !important`,
  },
}));

export const NotFoundChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.text.disabled,
  },
  background: `${alpha(theme.palette.background.disabled, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.icon.default} !important`,
  },
}));

export const MaintainanceChip = styled(Chip)(({ theme }) => ({
  ...baseChipStyles(theme),
  '& .MuiChip-label': {
    color: theme.palette.warning.main,
  },
  background: `${alpha(theme.palette.warning.main, 0.19)} !important`,
  '& .MuiSvgIcon-root': {
    color: `${theme.palette.warning.main} !important`,
  },
}));

export const ColumnWrapper = styled('div')(({ theme }) => ({
  margin: theme.spacing(2),
  padding: theme.spacing(2),
  background: `${theme.palette.background.secondary}10`,
}));

export const OperationButton = styled(Grid2)(({ theme }) => ({
  [theme?.breakpoints?.down(1180)]: {
    marginRight: '25px',
  },
}));

export const FormatterWrapper = styled(Box)({
  marginBlock: '0.4rem',
});
export const ContentContainer = styled(Grid2)(({ theme }) => ({
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
  active: { '& $line': { background: theme.palette.primary.main, transition: 'all 1s ease-in' } },
  completed: {
    '& $line': { background: theme.palette.primary.main, transition: 'all 1s ease-in' },
  },
  line: {
    height: 3,
    border: 0,
    backgroundColor:
      theme.palette.mode === 'dark'
        ? theme.palette.background.default
        : alpha(theme.palette.common.white, 0.95),
    borderRadius: 1,
    transition: 'all 0.5s ease-out ',
  },
}));

export const CONNECTION_STATE_COLORS = (theme) => ({
  [CONNECTION_STATES.CONNECTED]: theme.palette.background.success.default,
  [CONNECTION_STATES.REGISTERED]:
    theme.palette.mode === 'dark' ? theme.palette.secondary.light : theme.palette.secondary.dark,
  [CONNECTION_STATES.DISCOVERED]: theme.palette.info.main,
  [CONNECTION_STATES.IGNORED]: theme.palette.text.secondary,
  [CONNECTION_STATES.DELETED]: theme.palette.background.error.default,
  [CONNECTION_STATES.MAINTENANCE]: theme.palette.warning.main,
  [CONNECTION_STATES.DISCONNECTED]: theme.palette.background.warning.default,
  [CONNECTION_STATES.NOTFOUND]: theme.palette.text.disabled,
});
