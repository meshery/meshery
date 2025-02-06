import { Badge, Box, Button, Drawer, Grid, IconButton, Typography, styled } from '@layer5/sistent';
import { alpha } from '@mui/system';
import { STATUS } from './constants';

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

export const SideList = styled('div')(() => ({
  width: '45rem',
  maxWidth: '95vw',
}));
export const StyledSubtitle = styled(Typography)(() => ({
  textTransform: 'capitalize',
}));

export const NotificationButton = styled(IconButton)(() => ({
  height: '100%',
  '&:hover': {
    color: '#fff',
  },
}));

export const NotificationDrawer = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
  boxShadow: theme.shadows[2],
  zIndex: 1205,
}));

export const StyledNotificationDrawer = styled(Drawer)(({ theme, isNotificationCenterOpen }) => ({
  '& .MuiDrawer-paper': {
    backgroundColor: theme.palette.background.brand?.disabled,
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'space-between',
    boxShadow: theme.shadows[2],
    zIndex: 1205,
    position: 'absolute',
    right: isNotificationCenterOpen ? 0 : '-42.8rem', // Dynamically handle full and peek views
    transition: 'right 0.3s ease-in-out',
  },
}));

export const FullView = styled('div')(() => ({
  right: 0,
  transition: '0.3s ease-in-out !important',
}));

export const PeekView = styled('div')(() => ({
  right: '-42.8rem',
  transition: '0.3s ease-in-out !important',
}));

export const Container = styled('div')(() => ({
  padding: '1.25rem',
}));

export const Header = styled('header')(({ theme }) => ({
  display: 'flex',
  gap: '0.5rem',
  justifyContent: 'space-between',
  alignItems: 'center',
  height: '5.65rem',
  flexWrap: 'wrap',
  background: theme.palette.background.card,
}));

export const NotificationContainer = styled('div')(({ theme }) => ({
  padding: '1.25rem',
  display: 'flex',
  gap: '0.5rem',
  justifyContent: 'space-between',
  alignItems: 'center',
  height: '5.65rem',
  flexWrap: 'wrap',
  background: theme.palette.background.card,
}));

export const Title = styled('div')(() => ({
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem',
}));

export const TitleBellIcon = styled('div')(() => ({
  width: '2.25rem',
  height: '2.25rem',
  borderRadius: '100%',
  backgroundColor: 'black',
  display: 'flex',
  padding: '0.2rem',
  justifyContent: 'center',
  alignItems: 'center',
  cursor: 'pointer',
}));

export const SeverityChip = styled('div')(() => ({
  borderRadius: '0.25rem',
  display: 'flex',
  gap: '0.45rem',
  justifyContent: 'start',
  alignItems: 'center',
  fontSize: '.95rem',
  fontWeight: 600,
  cursor: 'pointer',
}));

export const SeverityChips = styled('div')(() => ({
  display: 'flex',
  flexWrap: 'wrap',
  gap: '0.75rem',
  alignItems: 'center',
}));

export const Notification = styled('div')(({ theme }) => ({
  margin: theme.spacing(0.5, 1),
}));

export const DetailsContainer = styled(Box)(({ theme }) => ({
  color: theme.palette.text.default,
  boxShadow: theme.shadows[4],
  borderRadius: '0.25',
  paddingInline: '0.5rem',
  paddingBlock: '0.25rem',
  width: '12.5rem',
}));

export const ListDetails = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gridGap: '0.5rem',
  marginBlock: '0.5rem',
  borderRadius: '0.25rem',
  backgroundColor: theme.palette.background.card,
  '&:hover': {
    backgroundColor: alpha(theme.palette.text.brand, 0.25),
  },
}));

export const ListItem = styled(Box)(() => ({
  display: 'flex',
  gridGap: '0.5rem',
  alignItems: 'center',
  justifyContent: 'space-around',
}));

export const SoicialListItem = styled(Box)(({ theme }) => ({
  display: 'flex',
  backgroundColor: theme.palette.background.card,
  alignItems: 'center',
  justifyContent: 'space-around',
  padding: '.65rem',
}));

export const ListButton = styled(Button)(() => ({
  height: '100%',
  width: '100%',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'start',
}));

export const Expanded = styled(Grid)(({ theme }) => ({
  [theme.breakpoints.down('md')]: {
    padding: '0.5rem',
  },
}));
export const StyledAvatarStack = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: '0',
  '& .MuiAvatar-root': {
    width: '2rem',
    height: '2rem',
    border: `0.05rem solid ${theme.palette.text.default}`,
  },
}));
export const ActorAvatar = styled(Grid)(() => ({
  display: 'flex',

  justifyContent: 'center',
  alignItems: 'start',
}));

export const Message = styled(Typography)(() => ({
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  overflowWrap: 'break-word',
  // max of min of 20rem or 50vw
  maxWidth: 'min(25rem, 50vw)',
  width: '100%',
}));

export const GridItem = styled(Grid)(() => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'start',
  gap: '0.25rem',
}));
export const MenuPaper = styled(Box)(({ theme }) => ({
  color: theme.palette.icon.secondary,
  boxShadow: theme.shadows[4],
  borderRadius: '0.25rem',
  paddingInline: '0.5rem',
  paddingBlock: '0.25rem',
  width: '12.5rem',
}));
export const OptionList = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  gridGap: '0.5rem',
  marginBlock: '0.5rem',
  borderRadius: '0.25rem',
  backgroundColor: theme.palette.background.card,
  '&:hover': {
    backgroundColor: alpha(theme.palette.text.brand, 0.25),
  },
}));
export const OptionListItem = styled(Box)(() => ({
  display: 'flex',
  gridGap: '0.5rem',
  alignItems: 'center',
  justifyContent: 'space-around',
}));
export const SocialListItem = styled(Box)(({ theme }) => ({
  display: 'flex',
  backgroundColor: alpha(theme.palette.background?.surfaces, 0.25),
  alignItems: 'center',
  justifyContent: 'space-around',
  padding: '.65rem',
}));
export const StyledBadge = styled(Badge)(({ badgeColor }) => ({
  '& .MuiBadge-badge': {
    backgroundColor: badgeColor,
  },
}));

export const Root = styled('div')(({ notificationcolor, status }) => ({
  width: '100%',
  borderRadius: '0.25rem',
  border: `0.1rem solid ${notificationcolor}`,
  borderLeftWidth: status === STATUS.UNREAD ? '0.5rem' : '0.1rem',
  marginBlock: '0.5rem',
}));
export const Summary = styled(Grid)(({ notificationcolor }) => ({
  paddingBlock: '0.5rem',
  paddingInline: '0.25rem',
  cursor: 'pointer',
  backgroundColor: alpha(notificationcolor, 0.2),
}));
