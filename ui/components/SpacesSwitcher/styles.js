import { height, margin, width } from '@mui/system';
import {
  styled,
  Drawer as MuiDrawer,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  DARK_BLUE_GRAY,
  Box,
  Avatar,
  Button,
  Checkbox,
  Typography,
  Modal,
} from '@sistent/sistent';

const DRAWER_WIDTH = 300;

export const DrawerHeader = styled('div')(({ theme, open }) => ({
  display: 'flex',
  alignItems: 'flex-end',
  justifyContent: open ? 'flex-end' : 'center',
  marginBottom: '1rem',
  height: '100%',
  ...theme.mixins.toolbar,
}));

export const StyledDrawer = styled(MuiDrawer, { shouldForwardProp: (prop) => prop !== 'open' })(
  ({ theme, open }) => ({
    width: DRAWER_WIDTH,
    flexShrink: 0,
    whiteSpace: 'nowrap',
    boxSizing: 'border-box',
    [theme.breakpoints.down('xl')]: {
      position: 'absolute',
      height: '100%',
      zIndex: theme.zIndex.drawer,
    },
    ...(open && {
      width: DRAWER_WIDTH,
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
      overflowX: 'hidden',
      '& .MuiDrawer-paper': {
        backgroundColor:
          theme.palette.mode == 'light' ? theme.palette.background.paper : DARK_BLUE_GRAY,
        width: DRAWER_WIDTH,
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.enteringScreen,
        }),
        overflowX: 'hidden',
      },
    }),
    ...(!open && {
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen,
      }),
      overflowX: 'hidden',
      width: `calc(${theme.spacing(7)} + 1px)`,
      [theme.breakpoints.up('sm')]: {
        width: `calc(${theme.spacing(8)} + 1px)`,
      },
      '& .MuiDrawer-paper': {
        backgroundColor:
          theme.palette.mode == 'light' ? theme.palette.background.paper : DARK_BLUE_GRAY,
        transition: theme.transitions.create('width', {
          easing: theme.transitions.easing.sharp,
          duration: theme.transitions.duration.leavingScreen,
        }),
        overflowX: 'hidden',
        width: `calc(${theme.spacing(7)} + 1px)`,
        [theme.breakpoints.up('sm')]: {
          width: `calc(${theme.spacing(8)} + 1px)`,
        },
      },
    }),
  }),
);

export const StyledMainContent = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  width: '100%',
  height: '100%',
  padding: '1rem 2rem',
  overflowY: 'auto',
  [theme.breakpoints.down('xl')]: {
    paddingLeft: '5rem',
  },
}));

export const DesignList = styled(List)({
  width: '100%',
  padding: '0px',
  marginTop: '0rem',
  paddingBottom: '1rem',
  overflowY: 'auto',
  position: 'relative',
});

export const LoadingContainer = styled('div')({
  textAlign: 'center',
});

export const GhostContainer = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  position: 'absolute',
  top: '-85px',
  zIndex: '0',
});

export const GhostImage = styled('img')({
  marginRight: '0.5rem',
});

export const GhostText = styled('div')({
  fontSize: 15,
});

export const StyledListItem = styled(ListItem)({
  display: 'flex',
  alignItems: 'center',
  gap: '1.5rem',
  marginBlock: '0',
  paddingBlock: '6px',
  position: 'relative',
});

export const StyledUserDetailsContainer = styled('div')({
  display: 'flex',
  alignItems: 'start',
  flexDirection: 'column',
  marginLeft: '1rem',
  gap: '0.1rem',
});

export const StyledAvatarContainer = styled('div')({
  display: 'flex',
  alignItems: 'center',
  width: '100%',
});

export const StyledMainMenuComponent = styled('div')({
  opacity: 0,
  visibility: 'hidden',
  transition: 'all 0.2s ease-in-out',
  cursor: 'pointer',
});

export const StyledListItemText = styled(ListItemText)(({ theme }) => ({
  cursor: 'pointer',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  margin: '0',
  maxWidth: theme.spacing(25),

  '& .MuiTypography-root': {
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
}));

export const StyledListIcon = styled(ListItemIcon)({
  minWidth: '0px',
  paddingRight: '1rem',
});

export const StyledUpdatedText = styled('p')({
  margin: '0',
  fontSize: '0.8rem',
  fontStyle: 'italic',
  color: '#647881',
  cursor: 'pointer',
});

export const StyledSmallAvatarContainer = styled('div')(({ transform, clipPath }) => ({
  position: 'absolute',
  bottom: '-1.65rem',
  left: '2rem',
  padding: '0rem',
  transition: 'all 0.4s ease',
  transform: transform,
  clipPath: clipPath,
  zIndex: 1,
}));

export const StyledSmallAvatar = styled(Avatar)(({ borderColor }) => ({
  height: '23px',
  width: '23px',
  border: borderColor ? `2px solid ${borderColor}` : undefined,
}));

export const StyledResponsiveButton = styled(Button)(({ theme }) => ({
  '& .MuiButton-startIcon': {
    [theme.breakpoints.down('sm')]: {
      margin: 0,
    },
  },
  [theme.breakpoints.down('sm')]: {
    minWidth: '40px',
    padding: '8px',
  },
}));

export const StyledMuiDoubleCheckbox = styled(Checkbox)(({ theme }) => ({
  color: theme.palette.text.secondary,
  '&.Mui-checked': {
    color: theme.palette.text.secondary,
  },
  '& .MuiSvgIcon-root': {
    width: '1.25rem',
    height: '1.25rem',
  },
}));

export const StyledTypography = styled(Typography)(({ theme }) => ({
  [theme.breakpoints.up('lg')]: {
    width: theme.spacing(12),
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
}));

export const StyledModal = styled(Modal)(({ theme }) => ({

  '& .MuiDialog-paperFullScreen': {
    margin: '0'
  },

  '& .MuiDialog-paperFullWidth' : {
    width: "90%",
    height: "80%"
  },

  '& .MuiDialog-paper': {
    maxWidth: '100%',

    [theme.breakpoints.down('md')]: {
      margin: '0',
      width: '100%',
      maxWidth: '100%',
      height: '100%',
      maxHeight: '100%',
    },
  },
}))