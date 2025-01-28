import {
  Box,
  Button,
  ButtonGroup,
  Divider,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  styled,
  Drawer,
} from '@layer5/sistent';
import { disabledStyleWithOutOpacity } from '../../css/disableComponent.styles';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

export const FallbackWrapper = styled(Box)(() => ({
  margin: '2rem',
}));

export const TryAgainButton = styled(Button)(({ theme }) => ({
  background: 'transparent',
  border: `1px solid ${theme.palette.border.brand}`,
  '&:hover': {
    border: `1px solid ${theme.palette.border.brand}`,
  },
}));

export const EditButton = styled(Button)(({ theme }) => ({
  backgroundImage: theme.palette.background.brand.prominent,
  backgroundColor: `${
    theme.palette.background.brand.prominent || theme.palette.background.brand.default
  }`,
  '@media (max-width: 768px)': {
    minWidth: '50px',
  },
}));

export const TextButton = styled('span')(({ style }) => ({
  marginLeft: '0.5rem',
  display: 'block',
  '@media (max-width: 853px)': {
    display: 'none',
  },
  ...style,
}));

export const ToolBarButtonContainer = styled('span')(() => ({
  margin: '0 1rem 0 0',
  '@media (max-width: 400px)': {
    margin: '0 0.25rem 0 0',
  },
}));

export const StyledListItem = styled(ListItem, {
  shouldForwardProp: (prop) => prop !== 'disableLogo',
})(({ disableLogo }) => ({
  cursor: 'pointer',
  backgroundColor: '#263238',
  boxShadow: '0 -1px 0 #404854 inset',
  paddingTop: '1.325rem',
  paddingBottom: '1.325rem',
  position: 'sticky',
  top: 0,
  zIndex: 5,
  color: 'rgba(255, 255, 255, 0.7)',
  fill: '#fff',
  ...(disableLogo && {
    ...disabledStyleWithOutOpacity,
  }),
  '&:hover': {
    '& .expandMoreIcon': {
      opacity: 1,
      transition: 'opacity 200ms ease-in',
    },
  },
}));

export const MainLogo = styled('img')(({ theme }) => ({
  marginRight: theme.spacing(1),
  marginTop: theme.spacing(1),
  marginLeft: theme.spacing(-1),
  width: 40,
  height: 40,
  borderRadius: 'unset',
}));

export const MainLogoCollapsed = styled('img')(({ theme }) => ({
  marginRight: theme.spacing(1),
  marginTop: theme.spacing(1),
  marginLeft: theme.spacing(-0.5),
  width: 40,
  height: 40,
  borderRadius: 'unset',
}));

export const MainLogoTextCollapsed = styled('img')(({ theme }) => ({
  marginLeft: theme.spacing(1),
  marginTop: theme.spacing(1),
  width: 170,
  borderRadius: 'unset',
}));

export const MainLogoText = styled('img')(({ theme }) => ({
  marginLeft: theme.spacing(0.5),
  marginTop: theme.spacing(1),
  width: 170,
  borderRadius: 'unset',
}));

export const NavigatorList = styled(List)({
  padding: 0,
});

export const NavigatorListItem = styled(ListItem, {
  shouldForwardProp: (prop) =>
    prop !== 'depth' && prop !== 'isDrawerCollapsed' && prop !== 'isActive',
})(({ theme, depth, isDrawerCollapsed, isActive }) => ({
  paddingLeft: isDrawerCollapsed ? theme.spacing(2) : '',
  paddingRight: isDrawerCollapsed ? '16px' : '',
  color: isActive ? '#4fc3f7' : 'rgba(255, 255, 255, 0.7)',
  fill: isActive ? '#4fc3f7' : '#fff',
  '&:hover': {
    backgroundColor: 'rgba(0, 187, 166, 0.5)',
    '& $expandMoreIcon': {
      opacity: 1,
      transition: 'opacity 200ms ease-in',
    },
  },
  paddingTop: 4,
  paddingBottom: 4,
  // color: 'rgba(255, 255, 255, 0.7)',
  // fill: '#fff',
}));

export const NavigatorListItemII = styled(ListItem, {
  shouldForwardProp: (prop) =>
    prop !== 'depth' && prop !== 'isDrawerCollapsed' && prop !== 'isActive',
})(({ theme, depth, isDrawerCollapsed, isActive }) => ({
  paddingLeft: isDrawerCollapsed
    ? theme.spacing(2)
    : depth === 1
      ? theme.spacing(3)
      : theme.spacing(5),
  paddingRight: isDrawerCollapsed ? '16px' : '',
  color: isActive ? '#4fc3f7' : 'rgba(255, 255, 255, 0.7)',
  fill: isActive ? '#4fc3f7' : '#fff',
  '&:hover': {
    backgroundColor: 'rgba(0, 187, 166, 0.5)',
    '& $expandMoreIcon': {
      opacity: 1,
      transition: 'opacity 200ms ease-in',
    },
  },
  paddingTop: 4,
  paddingBottom: 4,
  // color: 'rgba(255, 255, 255, 0.7)',
  // fill: '#fff',
}));

export const NavigatorListItemIII = styled(ListItem, {
  shouldForwardProp: (prop) =>
    prop !== 'depth' && prop !== 'isDrawerCollapsed' && prop !== 'isActive' && prop !== 'isShow',
})(({ theme, depth, isDrawerCollapsed, isActive, isShow }) => ({
  paddingLeft: isDrawerCollapsed
    ? theme.spacing(2)
    : depth === 1
      ? theme.spacing(3)
      : theme.spacing(5),
  paddingRight: isDrawerCollapsed ? '16px' : '',
  color: isActive ? '#4fc3f7' : 'rgba(255, 255, 255, 0.7)',
  fill: isActive ? '#4fc3f7' : '#fff',
  '&:hover': {
    backgroundColor: 'rgba(0, 187, 166, 0.5)',
    '& $expandMoreIcon': {
      opacity: 1,
      transition: 'opacity 200ms ease-in',
    },
  },
  paddingTop: 4,
  paddingBottom: 4,
  pointerEvents: isShow ? 'none' : 'auto',
  opacity: isShow ? 0.5 : '',
}));

export const SideBarListItem = styled(ListItem, {
  shouldForwardProp: (prop) => prop !== 'link' && prop !== 'isActive' && prop !== 'isShow',
})(({ link, isActive, isShow }) => ({
  color: isActive ? '#4fc3f7' : 'rgba(255, 255, 255, 0.7)',
  fill: isActive ? '#4fc3f7' : '#fff',
  '&:hover': {
    ...(link && {
      backgroundColor: 'rgba(0, 187, 166, 0.5)', // Only applies when link is true
    }),
    '& $expandMoreIcon': {
      opacity: 1,
      transition: 'opacity 200ms ease-in',
    },
  },
  paddingTop: 4,
  paddingBottom: 4,
  pointerEvents: isShow ? 'none' : 'auto',
  opacity: isShow ? 0.5 : '',
  fontSize: '14px',
}));

export const SideBarText = styled(ListItemText)(({ drawerCollapsed }) => ({
  opacity: drawerCollapsed ? 0 : 1,
  transition: drawerCollapsed ? 'opacity 200ms ease-in-out' : 'opacity 200ms ease-in-out',
  // color: '#fff',
  fontSize: '14px',
  '& .MuiListItemText-primary': {
    fontSize: '14px',
  },
  '& .MuiTypography-root': {
    fontSize: '14px',
  },
}));

export const PrimaryElement = styled(SideBarText)(({ theme }) => ({
  color: 'inherit',
  fontSize: theme.typography.fontSize,
  '&$textDense': { fontSize: theme.typography.fontSize },
}));
export const RootDiv = styled('div', {
  shouldForwardProp: (prop) => prop !== 'show',
})(({ show }) => ({
  cursor: show ? '' : 'not-allowed',
  '& svg': {
    width: '1.21rem',
    height: '1.21rem',
  },
}));

export const HideScrollbar = styled(List)(() => ({
  overflow: 'hidden auto',
  'scrollbar-width': 'none',
  '-ms-overflow-style': 'none',
  '&::-webkit-scrollbar': {
    display: 'none',
  },
}));

export const SecondaryDivider = styled(Divider)(({ theme }) => ({
  marginTop: theme.spacing(1),
  marginBottom: theme.spacing(1),
  borderColor: '#404854',
}));

export const ExpandMoreIcon = styled(FontAwesomeIcon, {
  shouldForwardProp: (prop) => prop !== 'isCollapsed' && prop !== 'hasChildren',
})(({ isCollapsed, hasChildren }) => ({
  opacity: 0,
  cursor: 'pointer',
  display: hasChildren ? 'inline-block' : 'none',
  transform: isCollapsed ? 'rotate(180deg) translateX(-0.8px)' : 'translateX(3px)', // Rotate if collapsed
  transition: 'transform 200ms ease-in-out', // Smooth rotation animation
  '&:hover': {
    opacity: 1,
    color: '#4fc3f7',
  },
}));

export const ChevronIcon = styled(FontAwesomeIcon)(({ theme }) => ({
  color: theme.palette.icon.default,
  cursor: 'pointer',
  margin: '0.75rem 0.2rem',
  width: '0.8rem',
  verticalAlign: 'middle',
}));

export const MainListIcon = styled(ListItemIcon)(({ theme }) => ({
  minWidth: theme.spacing(3.5),
  paddingTop: theme.spacing(0.5),
  textAlign: 'center',
  display: 'inline-table',
  paddingRight: theme.spacing(0.5),
  marginLeft: theme.spacing(0.8),
}));

export const ListIconSide = styled(ListItemIcon)(({ theme }) => ({
  paddingTop: theme.spacing(0.5),
  textAlign: 'center',
  display: 'inline-table',
  paddingRight: theme.spacing(0.5),
  marginLeft: theme.spacing(0.8),
  color: '#fff',
  opacity: '0.7',
  transition: 'opacity 200ms linear',
  '&:hover': {
    opacity: 1,
    background: 'transparent',
  },
  '&:focus': {
    opacity: 1,
    background: 'transparent',
  },
}));

export const HiddenText = styled(ListItemText)(({ drawerCollapsed }) => ({
  opacity: drawerCollapsed ? 0 : 1,
  color: '#fff',
  fontSize: '14px',
  transition: drawerCollapsed ? 'opacity 200ms ease-in-out' : 'opacity 200ms ease-in-out',
}));

export const LinkContainer = styled('div')(() => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: '100%',
  height: '30px',
}));

export const NavigatorDrawer = styled(Drawer, {
  shouldForwardProp: (prop) => prop !== 'isCollapsed',
})(({ theme, isCollapsed }) => ({
  width: isCollapsed ? theme.spacing(8) + 4 : 256,
  flexShrink: 0,
  whiteSpace: 'nowrap',
  boxSizing: 'border-box',
  '& .MuiDrawer-paper': {
    background: '#263238',
    width: isCollapsed ? theme.spacing(8) + 4 : 256,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    overflowX: 'hidden',
  },
}));

export const NavigatorWrapper = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  height: '100%',
});

export const NavigatorHelpIcons = styled(ButtonGroup, {
  shouldForwardProp: (prop) => prop !== 'isCollapsed',
})(({ isCollapsed }) => ({
  ...(isCollapsed && {
    marginRight: 4,
    alignItems: 'center',
  }),
  ...(!isCollapsed && {
    padding: '5px',
    '& > li': {
      padding: '0',
    },
  }),
}));

export const NavigatorFooter = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  marginTop: 'auto',
  marginBottom: '0.5rem',
});

export const ChevronButtonWrapper = styled('div', {
  shouldForwardProp: (prop) => prop !== 'isCollapsed',
})(({ isCollapsed }) => ({
  backgroundColor: isCollapsed ? '#515b60' : 'transparent',
  color: isCollapsed ? '#ffffff' : 'inherit',
  boxShadow: !isCollapsed
    ? '0.5px 0px 0px 0px rgb(0 0 0 / 20%), 1.5px 0px 0px 0px rgb(0 0 0 / 14%), 2.5px 1px 3px 0px rgb(0 0 0 / 12%)'
    : 'none',
  position: 'fixed',
  borderRadius: '0 5px 5px 0',
  cursor: 'pointer',
  bottom: '12%',
  left: isCollapsed ? '49px' : '257px',
  zIndex: '1400',
  width: 'auto',
  transition: 'left 225ms',
  transform: isCollapsed ? 'rotate(180deg)' : 'none',
  display: 'flex',
  justifyContent: 'center',

  '&:hover': {
    opacity: 1,
    background: !isCollapsed ? 'transparent' : undefined,
  },
  '&:focus': {
    opacity: 1,
    background: !isCollapsed ? 'transparent' : undefined,
  },
}));

export const NavigatorLink = styled('span')({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: '100%',
  height: '30px',
});

export const HelpListItem = styled(ListItem)({
  paddingLeft: 0,
  paddingTop: 4,
  paddingBottom: 4,
  color: 'rgba(255, 255, 255, 0.7)',
  fill: '#fff',
  '&:hover': {
    '& $expandMoreIcon': {
      opacity: 1,
      transition: 'opacity 200ms ease-in',
    },
  },
});

export const HelpButton = styled(IconButton, {
  shouldForwardProp: (prop) => prop !== 'isCollapsed',
})(({ isCollapsed }) => ({
  ...(isCollapsed && {
    height: '1.45rem',
    marginTop: '-4px',
    transform: 'translateX(0px)',
  }),
  ...(!isCollapsed && {
    transform: 'translateX(0.5px)',
  }),
}));

export const FixedSidebarFooter = styled('div')({
  display: 'flex',
  flexDirection: 'column',
  marginTop: 'auto',
  marginBottom: '0.5rem',
});

export const SidebarDrawer = styled(Drawer, {
  shouldForwardProp: (prop) => prop !== 'isCollapsed',
})(({ theme, isCollapsed }) => ({
  width: isCollapsed ? 68 : 256,
  flexShrink: 0,
  whiteSpace: 'nowrap',
  height: '100%',
  overflowX: 'hidden',
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: isCollapsed
      ? theme.transitions.duration.leavingScreen
      : theme.transitions.duration.enteringScreen,
  }),
  '& .MuiDrawer-paper': {
    width: isCollapsed ? 68 : 256,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: isCollapsed
        ? theme.transitions.duration.leavingScreen
        : theme.transitions.duration.enteringScreen,
    }),
    overflowX: 'hidden',
    background: '#263238',
  },
}));
