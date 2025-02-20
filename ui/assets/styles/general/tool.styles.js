import { styled, Box, tabMenu } from '@layer5/sistent';

export const ToolWrapper = styled(Box)(({ theme }) => ({
  marginBottom: '2rem',
  display: 'flex',
  justifyContent: 'space-between',
  backgroundColor: theme.palette.background.card,
  boxShadow: '0px 2px 4px -1px rgba(0,0,0,0.2)',
  height: '4rem',
  padding: '0.68rem',
  borderRadius: '0.5rem',
  position: 'relative',
  zIndex: '101',
}));

export const MeshModelToolbar = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isAnimated',
})(({ theme, isAnimated }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  backgroundColor: theme.palette.background.card,

  boxShadow: '0px 2px 4px -1px rgba(0,0,0,0.2)',
  height: isAnimated ? '4rem' : '0rem',
  padding: isAnimated ? '0.68rem' : '0rem',
  borderRadius: '0.5rem',
  position: 'relative',
  zIndex: isAnimated ? '125' : '0',
  marginBottom: '0.5rem',
  marginTop: '1rem',
}));

export const MainContainer = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isAnimated',
})(({ theme, isAnimated }) => ({
  backgroundColor: theme.palette.background.card,

  height: isAnimated ? '73vh' : '25rem',
  display: 'flex',
  position: 'relative',
  marginTop: '1rem',
  [theme.breakpoints.down('sm')]: {
    height: isAnimated ? '73rem' : '47rem',
  },
}));

export const InnerContainer = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isAnimated',
})(({ theme, isAnimated }) => ({
  display: 'flex',
  flexDirection: 'row',
  position: 'absolute',
  ...(isAnimated
    ? {
        width: '100%',
        top: '0%',
        paddingX: '2rem',
        transform: 'translate(0%, 0%)',
        justifyContent: 'center',
        left: '0%',
        backgroundColor: tabMenu.main,
      }
    : {
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
      }),
  [theme.breakpoints.down('sm')]: {
    flexDirection: isAnimated ? 'row' : 'column',
    ...(isAnimated && {
      paddingLeft: '1rem',
      overflowX: 'auto',
      padding: '0.4rem',
    }),
  },
}));

export const CardStyle = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isAnimated' && prop !== 'isSelected',
})(({ theme, isAnimated, isSelected }) => ({
  background: isAnimated ? tabMenu.hover : tabMenu.main,
  color: isSelected ? theme.palette.text.default : theme.palette.background.constant.white,
  height: isAnimated ? '3rem' : '10rem',
  width: isAnimated ? '15rem' : '13rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  margin: '0rem 0.7rem',
  paddingTop: '0.2rem',
  flexDirection: isAnimated ? 'row-reverse' : 'column',
  cursor: 'pointer',
  ...(isAnimated && {
    borderRadius: '8px 8px 0px 0px',
    paddingTop: '0.2rem',
    marginTop: '0.7rem',
  }),
  ...(isSelected && {
    backgroundColor: theme.palette.background.card,
  }),
  [theme.breakpoints.down('md')]: {
    ...(!isAnimated && {
      height: '10rem',
      width: '8.5rem',
    }),
  },
  [theme.breakpoints.down('sm')]: {
    ...(isAnimated
      ? {
          padding: '0.1rem',
          flexDirection: 'column',
          margin: '0rem 0.2rem',
          width: '10rem',
        }
      : {
          width: '13rem',
          marginTop: '0.7rem',
          marginRight: '0.5rem',
        }),
  },
}));

export const DetailsContainer = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isEmpty',
})(({ theme, isEmpty }) => ({
  height: '60vh',
  width: '50%',
  backgroundColor:
    theme.palette.mode === 'dark'
      ? theme.palette.background.constant.table
      : theme.palette.border.default,
  borderRadius: '6px',
  padding: isEmpty ? '2.5rem' : '1rem 2rem',
  overflowY: 'auto',
  ...(isEmpty && {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    justifyContent: 'center',
  }),
  [theme.breakpoints.down('sm')]: {
    width: '100%',
    padding: isEmpty ? '0.5rem' : '1rem',
    height: 'fit-content',
    maxHeight: '30rem',
  },
}));

export const TreeContainer = styled(Box)(({ theme }) => ({
  height: '30rem',
  width: '50%',
  margin: '1rem',
  display: 'flex',
  justifyContent: 'center',
  [theme.breakpoints.down('sm')]: {
    width: '100%',
  },
}));

export const TreeWrapper = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isAnimated',
})(({ theme, isAnimated }) => ({
  backgroundColor: theme.palette.background.card,

  display: 'flex',
  gap: '1rem',
  padding: '1rem',
  flexDirection: 'row',
  width: '100%',
  position: 'absolute',
  top: '3.7rem',
  opacity: isAnimated ? '1' : '0',
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column-reverse',
  },
}));

export const Segment = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'column',
  },
}));

export const FullWidth = styled(Box)(({ theme }) => ({
  width: '50%',
  [theme.breakpoints.down('sm')]: {
    width: '70%',
  },
}));
