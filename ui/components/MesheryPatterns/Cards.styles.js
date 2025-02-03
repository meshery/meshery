import { Grid, styled } from '@layer5/sistent';
import { makeStyles } from '@mui/styles';

const useStyles = makeStyles({
  fullScreenCodeMirror: {
    height: '100%',
    width: '100%',
    '& .CodeMirror': {
      minHeight: '300px',
      height: '100%',
      width: '100%',
    },
  },
});

export const fullScreenCodeMirrorstyles = {
  height: '100%',
  width: '100%',
  '& .CodeMirror': {
    minHeight: '300px',
    height: '100%',
    width: '100%',
  },
};

export const CardBackGrid = styled(Grid)(() => ({
  marginBottom: '0.25rem',
  minHeight: '6rem',
  position: 'relative',
}));

export const YamlDialogTitleGrid = styled(Grid)(() => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
}));

export const CardHeaderRight = styled('div')(() => ({
  display: 'flex',
  alignItems: 'center',
}));

export const GridBtnText = styled('span')(({ theme }) => ({
  marginLeft: '5px',
  display: 'flex',
  justifyContent: 'center',
  [theme.breakpoints.down(1370)]: { display: 'none' },
  [`${theme.breakpoints.up(1920)} and (max-width: 2200px)`]: {
    display: 'none',
  },
}));

export const GridCloneBtnText = styled('span')(({ theme }) => ({
  display: 'flex',
  justifyContent: 'center',
  marginLeft: '3px',
  [theme.breakpoints.down(1370)]: { display: 'none' },
  [`${theme.breakpoints.up(1920)} and (max-width: 2200px)`]: {
    display: 'none',
  },
}));

export const CardImg = styled('img')(() => ({
  marginRight: '0.5rem',
  // filter: theme.palette.secondary.img,
}));

export const CardNoPaper = styled('div')(() => ({
  padding: '0.5rem',
  fontSize: '3rem',
}));

export const CardNoText = styled('div')(() => ({
  fontSize: '2rem',
  marginBottom: '2rem',
}));

export const CardNoContainer = styled('div')(() => ({
  padding: '2rem',
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  flexDirection: 'column',
}));

export const CardPagination = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  marginTop: '2rem',
}));

export const PatternAddIcon = styled('div')(() => ({
  paddingLeft: '0.5rem',
  marginRight: '0.5rem',
}));

export const UpdateDeleteButtons = styled('div')(() => ({
  width: 'fit-content',
  margin: '10 0 0 auto',
  position: 'absolute',
  right: 0,
  bottom: 0,
}));

export const BottomContainer = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'flex-start',
  alignItems: 'center',
  marginTop: '1rem',
}));

export const CatalogCardButtons = styled('div')(() => ({
  display: 'flex',
  flexWrap: 'wrap',
  marginTop: '50px',
  height: '100%',
  gap: '.5rem',
}));

export default useStyles;
