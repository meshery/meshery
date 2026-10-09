import { styled } from '@/theme';

const drawerWidth = 256;

export const StyledFooterText = styled('span')({
  cursor: 'pointer',
  display: 'inline',
  verticalAlign: 'middle',
});

export const StyledRoot = styled('div')({
  display: 'flex',
  minHeight: '100vh',
  height: '100vh',
});

export const StyledFooterBody = styled('footer')(({ theme }) => ({
  padding: theme.spacing(2),
  backgroundColor:
    theme.palette.mode === 'dark' ? theme.palette.background.card : theme.palette.common.white,
}));
export const StyledMainContent = styled('main')(({ theme }) => ({
  backgroundColor:
    theme.palette.mode === 'dark'
      ? theme.palette.background.elevatedComponents
      : theme.palette.background.hover,
  flex: 1,
  padding: '48px 36px 24px',
  [theme.breakpoints.down('sm')]: {
    padding: '24px 10px 16px 10px',
  },
}));

export const StyledAppContent = styled('div', {
  shouldForwardProp: (prop) => prop !== 'isDrawerCollapsed' && prop !== 'canShowNav',
})(({ theme, canShowNav, isDrawerCollapsed }) => ({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  position: 'relative',
  overflowX: 'hidden',
  overflowY: 'hidden',

  [theme.breakpoints.down('sm')]: {
    marginLeft: canShowNav ? '4.25rem' : '0',
  },

  [theme.breakpoints.up('sm')]: {
    marginLeft: isDrawerCollapsed && canShowNav ? 68 : 0,
  },
}));

export const StyledContentWrapper = styled('div')({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  overflow: 'auto',
  minHeight: 0,
});

export const StyledDrawer = styled('nav', {
  shouldForwardProp: (prop) => prop !== 'isDrawerCollapsed',
})(({ theme, isDrawerCollapsed }) => ({
  [theme.breakpoints.up('sm')]: {
    width: isDrawerCollapsed ? theme.spacing(8.4) + 1 : drawerWidth,
    flexShrink: 0,
  },
  transition: theme.transitions.create('width', {
    easing: theme.transitions.easing.sharp,
    duration: isDrawerCollapsed
      ? theme.transitions.duration.leavingScreen
      : theme.transitions.duration.enteringScreen,
  }),
  '& > div:first-child': {
    height: 'inherit',
    width: 'inherit',
  },
  height: '100%',
  overflow: 'visible',
  paddingRight: '4rem',
  [theme.breakpoints.up('xs')]: {
    paddingRight: '0',
  },
  [theme.breakpoints.down('sm')]: {
    position: 'absolute',
  },
}));
