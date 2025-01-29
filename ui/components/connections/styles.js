import { alpha } from '@mui/material/styles';
import { Colors } from '../../themes/app';
import { CONNECTION_STATES } from '../../utils/Enum';
import { notificationColors } from '../../themes';
import { Chip, Grid, styled } from '@layer5/sistent';
import { StepConnector, StepLabel, Stepper } from '@mui/material';

const styles = (theme) => ({
  grid: { padding: theme.spacing(2) },
  tableHeader: {
    fontWeight: 'bolder',
    fontSize: 18,
  },
  muiRow: {
    '& .MuiTableCell-root': {},
  },
  statusSelect: {
    '& .MuiSelect-select.MuiSelect-select': {
      padding: '0 !important',
    },
  },
  createButton: {
    display: 'flex',
    justifyContent: 'flex-start',
    alignItems: 'center',
    whiteSpace: 'nowrap',
  },
  viewSwitchButton: {
    justifySelf: 'flex-end',
    marginLeft: 'auto',
    paddingLeft: '1rem',
  },
  statusChip: {
    minWidth: '145px',
    width: '100%',
    display: 'flex',
    justifyContent: 'flex-start',
    textTransform: 'capitalize',
    borderRadius: '2px',
    padding: '6px 8px',
    '& .MuiChip-label': {
      paddingTop: '3px',
      fontWeight: '400',
    },
    '&:hover': {
      boxShadow: '0px 1px 2px 0px rgba(0, 0, 0, 0.25)',
      cursor: 'pointer',
    },
  },
  appBar: {
    marginBottom: '3rem',
  },
  capitalize: {
    // textTransform: 'capitalize',
  },
  lowecase: {
    textTransform: 'lowecase',
  },
  expandedRows: {
    background: `${theme.palette.secondary.default}10`,
  },
  contentContainer: {
    [theme.breakpoints.down(1050)]: {
      flexDirection: 'column',
    },
    flexWrap: 'noWrap',
  },
  innerTableWrapper: {
    background: `linear-gradient(90deg, ${theme.palette.secondary.innertableBg1} 0.04%, ${theme.palette.secondary.innertableBg2} 100.04%)`,
    borderRadius: 0,
    padding: '0',
  },
  innerTableContainer: {
    background: theme.palette.secondary.innertableBg1,
    margin: '10px 10px 10px 13px',
    borderLeft: `9px solid ${theme.palette.secondary.pinball}`,
    borderRadius: '10px 0 0 10px',
    width: 'calc(100% - 23px)',
    border: 'none',
    overflowX: 'hidden',
  },
  noGutter: {
    padding: '0',
  },
  showMore: {
    color: Colors.keppelGreen,
    cursor: 'pointer',
  },
  bulkAction: {
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
  },
  centerContent: {
    display: 'flex',
    justifyContent: 'center',
  },
  tab: {
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  tabs: {
    height: '55px',
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  iconText: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
  },
  addClusterButtonClass: {
    borderRadius: 5,
    marginRight: '2rem',
    padding: '8px',
  },
  list: {
    display: 'flex',
    flexDirection: 'column',
    gridGap: '0.5rem',
    marginBlock: '0.5rem',
    borderRadius: '0.25rem',
    backgroundColor: theme.palette.secondary.honeyComb,
  },
  listButton: {
    '&:hover': {
      backgroundColor: alpha(theme.palette.secondary.link2, 0.25),
    },
  },
  listItem: {
    display: 'flex',
    gridGap: '0.5rem',
    alignItems: 'center',
    justifyContent: 'space-around',
  },
  button: {
    width: '100%',
    justifyContent: 'flex-start',
    '&:hover': {
      backgroundColor: 'transparent',
    },
  },
  listContainer: {
    width: '100%',
    justifyContent: 'flex-start',
    display: 'flex',
    alignItems: 'center',
  },

  /** Connection status select colors according to the status */
  ignored: {
    '& .MuiChip-label': {
      color: `${theme.palette.secondary.default}`,
    },
    background: `${theme.palette.secondary.default}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.default} !important`,
    },
  },
  connected: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.success,
    },
    background: `${theme.palette.secondary.success}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.success} !important`,
    },
  },
  registered: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.primary,
    },
    background: `${theme.palette.secondary.primary}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.primary} !important`,
    },
  },
  register: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.primary,
    },
    background: `${theme.palette.secondary.primary}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.primary} !important`,
    },
  },
  discovered: {
    '& .MuiChip-label': {
      color: notificationColors.info,
    },
    background: `${notificationColors.info}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${notificationColors.info} !important`,
    },
  },
  deleted: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.error,
    },
    background: `${theme.palette.secondary.lightError}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.error} !important`,
    },
  },
  maintenance: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.warning,
    },
    background: `${theme.palette.secondary.warning}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.warning} !important`,
    },
  },
  disconnected: {
    '& .MuiChip-label': {
      color: notificationColors.lightwarning,
    },
    background: `${notificationColors.lightwarning}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${notificationColors.lightwarning} !important`,
    },
  },
  notfound: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.text,
    },
    background: `${theme.palette.secondary.disableButtonBg}60 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.iconMain} !important`,
    },
  },
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

export default styles;

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
