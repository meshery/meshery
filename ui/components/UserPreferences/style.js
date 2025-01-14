import { Card, FormControl, Paper, styled } from '@layer5/sistent';

const styles = (theme) => ({
  orgSelect: {
    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
      borderColor: 'transparent !important',
    },
    width: '85%',
    backgroundColor: theme.palette.primary,
    color: theme.palette.secondary.text,
    '@media (max-width: 368px)': {
      width: '180px',
    },
    position: 'relative',
    overflow: 'hidden',
    '& svg': {
      backgroundColor: theme.palette.white,
      fill: theme.palette.charcoal,
      float: 'right',
      top: '0',
      right: '0',
      height: '100%',
    },
  },
  selectItem: {
    width: '100%',
    '& svg': {
      backgroundColor: 'inherit',
    },
  },
  org: {
    marginLeft: '10px',
    overflow: 'hidden',
    whiteSpace: 'nowrap',
    textOverflow: 'ellipsis',
  },
  orgIconWrapper: {
    width: '32px',
    height: '100%',
    position: 'relative',
  },
});

export const FormGroupWrapper = styled(FormControl)(() => ({
  padding: 20,
  border: '1.5px solid #969696',
  display: 'flex',
  width: '70%',
}));

export const FormContainerWrapper = styled('div')(() => ({
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'space-evenly',
  padding: 50,
}));

export const StatsWrapper = styled(Paper)({
  maxWidth: '100%',
  height: 'auto',
  borderTopLeftRadius: 0,
  borderTopRightRadius: 0,
  borderBottomLeftRadius: 3,
  borderBottomRightRadius: 3,
});

export const PaperRoot = styled(Paper)({
  flexGrow: 1,
  maxWidth: '100%',
  marginLeft: 0,
  borderTopLeftRadius: 3,
  borderTopRightRadius: 3,
});

export const RootContainer = styled('div')(({ theme }) => ({
  width: '100%',
  paddingLeft: theme.spacing(15),
  paddingRight: theme.spacing(15),
  paddingBottom: theme.spacing(10),
  paddingTop: theme.spacing(5),
}));

export const ProviderCard = styled(Card)(({ theme }) => ({
  border: '1px solid rgba(0,179,159,0.3)',
  margin: '20px 0px',
  backgroundColor: theme.palette.type === 'dark' ? '#293B43' : '#C9DBE3',
}));

export default styles;
