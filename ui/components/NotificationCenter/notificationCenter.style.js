import { makeStyles, styled } from '@material-ui/core';

export const DarkBackdrop = styled('div')(({ open }) => ({
  position: 'fixed',
  top: 0,
  left: 0,
  right: 0,
  bottom: 0,
  backgroundColor: 'rgba(0, 0, 0, 0.5)', // Adjust the opacity as needed
  display: open ? 'block' : 'none', // Show only when anchored to mobile
  zIndex: '1202', // Ensure it's behind the container
}));
export const useStyles = makeStyles((theme) => ({
  sidelist: {
    width: '45rem',
    maxWidth: '95vw',
  },
  notificationButton: {
    height: '100%',
    '&:hover': {
      color: theme.palette.secondary.whiteIcon,
    },
  },
  notificationDrawer: {
    backgroundColor: theme.palette.secondary.drawer,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    boxShadow: theme.shadows[2],
    zIndex: 1205,
  },
  drawerButton: {
    padding: '0.45rem',
    margin: '0.2rem',
    backgroundColor: theme.palette.secondary.dark,
    color: '#FFFFFF',
    '&:hover': {
      backgroundColor: '#FFFFFF',
      color: theme.palette.secondary.dark,
    },
  },
  fullView: {
    right: 0,
    transition: '0.3s ease-in-out !important',
  },
  peekView: {
    right: '-42.8rem',
    transition: '0.3s ease-in-out !important',
  },

  container: {
    padding: '1.25rem',
  },
  header: {
    display: 'flex',
    gap: '0.5rem',
    justifyContent: 'space-between',
    alignItems: 'center',
    height: '5.65rem',
    flexWrap: 'wrap',
    background: theme.palette.secondary.headerColor,
  },
  title: {
    display: 'flex',
    alignItems: 'center',
    gap: '0.5rem',
  },
  titleBellIcon: {
    width: '2.25rem',
    height: '2.25rem',
    borderRadius: '100%',
    backgroundColor: 'black',
    display: 'flex',
    padding: '0.2rem',
    justifyContent: 'center',
    alignItems: 'center',
    cursor: 'pointer',
  },
  // Aggregrate Level Chips
  severityChip: {
    borderRadius: '0.25rem',
    display: 'flex',
    gap: '0.45rem',
    justifyContent: 'start',
    alignItems: 'center',
    fontSize: '.95rem',
    fontWeight: 600,
    cursor: 'pointer',
  },
  // Aggregrate Level Chips
  severityChips: {
    display: 'flex',
    flexWrap: 'wrap',
    gap: '0.75rem',
    alignItems: 'center',
  },
  notification: {
    margin: theme.spacing(0.5, 1),
  },
}));

export const useNavNotificationIconStyles = makeStyles(() => ({
  root: (props) => ({
    '& .MuiBadge-badge': {
      backgroundColor: props.badgeColor,
    },
  }),
}));
