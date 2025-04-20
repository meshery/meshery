import { styled, Box, tabMenu, accentGrey, slateGray } from '@layer5/sistent';

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

export const MeshModelToolbar = styled(Box)(({ theme }) => ({
  display: 'flex',
  justifyContent: 'space-between',
  backgroundColor: theme.palette.background.card,
  boxShadow: '0px 2px 4px -1px rgba(0,0,0,0.2)',
  height: '4rem',
  padding: '0.68rem',
  borderRadius: '0.5rem',
  position: 'relative',
  zIndex: '125',
  marginBottom: '0.5rem',
  marginTop: '1rem',
}));

export const MainContainer = styled(Box)(({ theme }) => ({
  backgroundColor: 'inherit',
  borderRadius: '0.25rem',
  height: '68vh',
  display: 'flex',
  position: 'relative',
  [theme.breakpoints.down('sm')]: {
    height: '73rem',
  },
}));

export const InnerContainer = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  gap: '0rem',
  alignItems: 'flex-start',
  justifyContent: 'center',
  width: '100%',
  top: '0%',
  paddingX: '2rem',
  transform: 'translate(0%, 0%)',
  left: '0%',
  borderRadius: '0.25rem',
  overflow: 'scroll',
  backgroundColor: 'inherit',
  [theme.breakpoints.down('sm')]: {
    flexDirection: 'row',
    paddingLeft: '1rem',
    overflowX: 'auto',
    padding: '0.4rem',
  },
}));

export const CardStyle = styled(Box, {
  shouldForwardProp: (prop) => prop !== 'isSelected',
})(({ theme, isSelected }) => ({
  background: isSelected
    ? accentGrey[30] // color when tab is selected
    : theme.palette.mode === 'dark' 
      ? accentGrey[10] // color for inactive tabs
      : accentGrey[20],
  color: isSelected ? theme.palette.text.default : theme.palette.background.constant.white,
  height: '3rem',
  width: '15rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  margin: '0rem 0.7rem',
  paddingTop: '0.2rem',
  flexDirection: 'row-reverse',
  cursor: 'pointer',
  borderRadius: '8px 8px 0px 0px',
  marginTop: '0.7rem',
  transition: 'background-color 0.2s ease-in-out',
  '&:hover': {
    background: isSelected
      ? accentGrey[30]
      : theme.palette.mode === 'dark'
        ? theme.palette.background.hover
        : theme.palette.background.inverseHover,
  },
  [theme.breakpoints.down('md')]: {
    width: '8.5rem',
  },
  [theme.breakpoints.down('sm')]: {
    padding: '0.1rem',
    flexDirection: 'column',
    margin: '0rem 0.2rem',
    width: '10rem',
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

export const TreeWrapper = styled(Box)(({ theme }) => ({
  backgroundColor: accentGrey[30],
  display: 'flex',
  gap: '1rem',
  padding: '1rem',
  flexDirection: 'row',
  width: '100%',
  position: 'absolute',
  top: '3.7rem',
  borderBottomLeftRadius: '0.5rem',
  borderBottomRightRadius: '0.5rem',
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
