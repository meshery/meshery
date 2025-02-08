import { styled, AppBar, Toolbar, Paper, MenuIcon, IconButton, darkTeal } from '@layer5/sistent';

export const HeaderAppBar = styled(AppBar, {
  shouldForwardProp: (prop) => prop !== 'isDrawerCollapsed',
})(({ theme, isDrawerCollapsed }) => ({
  backgroundColor: theme.palette.background.tabs,
  shadowColor: '#808080',
  ...(isDrawerCollapsed
    ? {
        [theme.breakpoints.between(635, 732)]: { padding: theme.spacing(0.75, 1.4) },
        [theme.breakpoints.between(600, 635)]: { padding: theme.spacing(0.4, 1.4) },
      }
    : {}),
}));

export const StyledToolbar = styled(Toolbar, {
  shouldForwardProp: (prop) => prop !== 'isDrawerCollapsed',
})(({ theme, isDrawerCollapsed }) => ({
  minHeight: 59,
  padding: 16,
  paddingLeft: isDrawerCollapsed ? 0 : 34,
  paddingRight: 34,
  backgroundColor: theme.palette.mode === 'dark' ? theme.palette.background.card : darkTeal.main,
  boxShadow: `3px 0px 4px ${theme.palette.background.brand.default}`,
  ...(isDrawerCollapsed
    ? {
        [theme.breakpoints.down('xs')]: {
          padding: 0,
        },
      }
    : {
        [theme.breakpoints.between(620, 732)]: {
          minHeight: 68,
          paddingLeft: 20,
          paddingRight: 20,
        },
      }),
}));

export const UserContainer = styled('div')(({ theme }) => ({
  paddingLeft: 1,
  display: 'flex',
  alignItems: 'center',
  [theme.breakpoints.down('xs')]: {
    width: '100%',
    justifyContent: 'flex-end',
    marginRight: '1rem',
    marginBlock: '0.5rem',
  },
}));

export const PageTitleWrapper = styled('div')({
  flexGrow: 1,
  marginRight: 'auto',
  flexWrap: 'nowrap',
  marginBlock: '0.5rem',
});

export const MenuIconButton = styled(IconButton)(({ theme }) => ({
  marginLeft: -theme.spacing(1),
}));

export const UserSpan = styled('span')(({ theme }) => ({
  marginLeft: theme.spacing(1),
}));

export const HeaderIcons = styled(MenuIcon)(({ theme }) => ({
  fontSize: '1.5rem',
  height: '1.5rem',
  width: '1.5rem',
  fill: theme.palette.background.constant.white,
  '&:hover': {
    fill: theme.palette.background.brand.default,
  },
}));

export const CBadge = styled('span')({
  fontSize: '0.65rem',
  backgroundColor: 'white',
  borderRadius: '50%',
  color: 'black',
  height: '1.30rem',
  width: '1.30rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  position: 'absolute',
  zIndex: 1,
  right: '-0.75rem',
  top: '-0.29rem',
});

export const CBadgeContainer = styled('div')({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  position: 'relative',
});

export const CMenuContainer = styled(Paper)(({ theme }) => ({
  backgroundColor: theme.palette.background.card,
  marginTop: '-1.2rem',
  borderRadius: '3px',
  padding: '1rem',
  boxShadow: '20px #979797',
  transition: 'linear .2s',
  transitionProperty: 'height',
}));

export const IconButtonAvatar = styled(IconButton)({
  padding: 4,
});
