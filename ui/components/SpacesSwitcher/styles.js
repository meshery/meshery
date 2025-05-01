import { styled, Drawer as MuiDrawer, List } from '@layer5/sistent';

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
    ...(open && {
      width: DRAWER_WIDTH,
      transition: theme.transitions.create('width', {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen,
      }),
      overflowX: 'hidden',
      '& .MuiDrawer-paper': {
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
  padding: '1rem',
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
