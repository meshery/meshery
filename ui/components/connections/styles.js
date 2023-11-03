import { Colors } from '../../themes/app';

const styles = (theme) => ({
  grid: { padding: theme.spacing(2) },
  tableHeader: {
    fontWeight: 'bolder',
    fontSize: 18,
  },
  muiRow: {
    '& .MuiTableCell-root': {
      textTransform: 'capitalize',
    },
  },
  statusSelect: {
    '& .MuiSelect-select.MuiSelect-select': {
      padding: '0 2px',
    },
    '& .MuiSelect-icon': {
      display: 'none',
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
  },
  statusChip: {
    minWidth: '130px !important',
    maxWidth: '100% !important',
    display: 'flex !important',
    justifyContent: 'flex-start !important',
    textTransform: 'capitalize',
    borderRadius: '3px !important',
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
    textTransform: 'capitalize',
  },
  lowecase: {
    textTransform: 'lowecase',
  },
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
  discovered: {
    '& .MuiChip-label': {
      color: theme.palette.secondary.warning,
    },
    background: `${theme.palette.secondary.warning}30 !important`,
    '& .MuiSvgIcon-root': {
      color: `${theme.palette.secondary.warning} !important`,
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
    '&:hover': {
      cursor: 'initial',
    },
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
  },
  listItem: {
    paddingTop: '0',
    paddingBottom: '0',
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
});

export default styles;
