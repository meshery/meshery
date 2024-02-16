import { alpha } from '@material-ui/core';
import { Colors } from '../../themes/app';
import { CONNECTION_STATES } from '../../utils/Enum';
import { notificationColors } from '../../themes';

const styles = (theme) => ({
  grid: { padding: theme.spacing(2) },
  tableHeader: {
    fontWeight: 'bolder',
    fontSize: 18,
  },
  muiRow: {
    '& .MuiTableCell-root': {
      // textTransform: 'capitalize',
    },
  },
  statusSelect: {
    '& .MuiSelect-select.MuiSelect-select': {
      padding: '0',
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
  chipFormControl: {
    minWidth: '100%',
    '& .MuiSelect-icon': {
      marginRight: '10px !important',
    },
  },
  statusChip: {
    minWidth: '145px !important',
    width: '100% !important',
    display: 'flex !important',
    justifyContent: 'flex-start !important',
    textTransform: 'capitalize',
    borderRadius: '0 !important',
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
