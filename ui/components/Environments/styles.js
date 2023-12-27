import { Colors } from '../../themes/app';

const styles = (theme) => ({
  /** Bulk action bar styles */
  bulkActionWrapper: {
    width: '100%',
    padding: '0.8rem',
    justifyContent: 'space-between',
    marginTop: '0.18rem',
    marginBottom: '1rem',
    borderRadius: '.25rem',
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
  },

  /** Card Styles */
  cardWrapper: {
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    padding: '20px',
    backgroundColor: theme.palette.secondary.white,
    '&:hover': {
      cursor: 'pointer',
    },
  },

  statistic: {
    display: 'flex',
    justifyContent: 'center',
    paddingX: '5px',
    fontWeight: '600',
    fontSize: '24px',
    textAlign: 'center',
  },

  statisticName: {
    display: 'flex',
    justifyContent: 'center',
    paddingX: '5px',
    fontWeight: '400',
    fontSize: '16px',
    textAlign: 'center',
  },

  tabCardContent: {
    padding: '16px',
    display: 'flex',
    justifyContent: 'center',
    flexDirection: 'row',
  },

  tabIconBox: {
    display: 'flex',
    alignItems: 'center',
    border: `1px solid ${Colors.keppelGreen}`,
    borderRadius: '5px 0 0 5px',
    padding: '8px 16px',
    borderRight: '0',
    '&:hover': {
      cursor: 'default',
    },
  },

  tabNameBox: {
    display: 'flex',
    alignItems: 'center',
    flexDirection: 'column',
    background: Colors.keppelGreen,
    color: theme.palette.secondary.white,
    p: '8px',
    border: `1px solid ${Colors.keppelGreen}`,
    borderRadius: '0 5px 5px 0',
    padding: '8px 16px',
  },

  tabTitle: {
    margin: '0',
    fontSize: '12px',
    fontWeight: '400',
    display: 'flex',
  },

  tabCount: {
    margin: '0',
    fontSize: '30px',
    fontWeight: '500',
    lineHeight: 1,
    marginBottom: '5px',
  },

  allocationButton: {
    padding: '10px 10px 1px 10px',
    borderRadius: '4px',
    height: '100%',
    display: 'flex',
    width: '100%',
    background: theme.palette.secondary.focused,
  },

  allocationWorkspace: {
    display: 'flex',
    width: '100%',
    gap: '10px',
    ['@media (min-width : 600px)']: {
      flexDirection: 'column',
      gap: '0',
    },
  },

  popupButton: {
    width: '100%',
    borderRadius: '4px',
    boxShadow: '0px 4px 4px 0px rgba(0, 0, 0, 0.25)',
    marginBottom: '10px',
    padding: '20px 10px',
    backgroundColor: theme.palette.secondary.mainBackground2,
    color: theme.palette.secondary.penColorPrimary,
    '&:hover': {
      background: theme.palette.secondary.mainBackground2,
    },
  },

  viewButton: {
    width: '100%',
    borderRadius: '4px',
    boxShadow: '0px 4px 4px 0px rgba(0, 0, 0, 0.25)',
    display: 'flex',
    flexDirection: 'column',
    marginBottom: '10px',
    color: theme.palette.secondary.penColorPrimary,
    '&:hover': {
      background: theme.palette.secondary.white,
    },
    padding: '15px 10px',
  },

  record: {
    borderBottom: `1px solid ${theme.palette.secondary.modalTabs}60`,
    display: 'flex',
    flexDirection: 'row',
    padding: '5px 0',
  },

  bulkSelectCheckbox: {
    padding: 0,
    marginRight: '0.5rem',
    color: 'white',
    '&:hover': {
      color: 'white',
      cursor: 'pointer',
    },
    '&.Mui-checked': {
      color: 'white',
    },
  },

  cardTitle: {
    fontSize: '1.25rem',
    fontWeight: 800,
    '&:hover': {
      cursor: 'default',
    },
  },

  organizationName: {
    fontSize: '0.9rem',
    display: 'flex',
    alignItems: 'end',
    padding: '0 5px',
    '&:hover': {
      cursor: 'default',
    },
  },

  styledIconButton: {
    background: 'transparent',
    border: 'none',
    '&:hover': {
      cursor: 'default',
    },
  },

  dateLabel: {
    fontStyle: 'italic',
    fontSize: '12px',
    '&:hover': {
      cursor: 'default',
    },
  },

  emptyDescription: {
    fontSize: '0.9rem',
    textAlign: 'left',
    fontStyle: 'italic',
  },

  descriptionLabel: {
    height: 'fit-content',
    fontStyle: 'normal',
    '&:hover': {
      cursor: 'default',
    },
  },

  status: {
    padding: '5px 20px',
    width: 'fit-content',
    border: '1px solid transparent',
    fontSize: '12px',
    '&:hover': {
      cursor: 'default',
    },
  },

  styledChip: {
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
  },

  createButtonWrapper: {
    display: 'flex',
    justifyContent: 'flex-start',
    alignItems: 'center',
    whiteSpace: 'nowrap',
  },

  editButton: {
    backgroundColor: Colors.keppelGreen,
    '&:hover': {
      backgroundColor: Colors.caribbeanGreen,
    },
    '@media (max-width: 768px)': {
      minWidth: '50px',
    },
  },

  textButton: {
    marginLeft: '0.5rem',
    display: 'block',
    '@media (max-width: 853px)': {
      display: 'none',
    },
  },

  iconButton: {
    minWidth: 'fit-content',
    '&:hover': {
      background: 'transparent',
    },
    padding: 0,
  },
});

export default styles;
