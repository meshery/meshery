import { Box, styled } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';

export const CreatAtContainer = styled('span')(({ theme, isBold }) => ({
  fontWeight: isBold ? 'bold' : '',
  whiteSpace: 'wrap',
  color: theme.palette.type === 'dark' ? '#ccc' : theme.palette.secondary.mainBackground,
  fontSize: '0.8rem',
}));

export const ActionContainer = styled(Box)({
  width: '100%',
  display: 'flex',
  gap: '1rem',
  justifyContent: 'end',
});

const useStyles = makeStyles((theme) => ({
  dialogBox: {
    '& .MuiDialog-paper': {
      maxWidth: '820px',
    },
  },
  dialogTitle: {
    backgroundColor: theme.palette.secondary.mainBackground,
    justifyContent: 'space-between',
    flexDirection: 'row',
    alignItems: 'center',
    padding: '12px 20px',
    gap: '146px',
    color: '#FFFFFF',
    textAlign: 'center',
    textOverflow: 'ellipsis',
    '& h2': {
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
    },
  },

  text: {
    fontFamily: 'Qanelas Soft, sans-serif',
    display: 'flex',
    flexDirection: 'column',
    '&.MuiTypography-root': {
      fontFamily: 'Qanelas Soft, sans-serif',
    },
  },
  textHeader: {
    fontFamily: 'Qanelas Soft, sans-serif',
    textAlign: 'center',
  },
  resourceName: {
    fontFamily: 'Qanelas Soft, sans-serif',
    textAlign: 'center',
    maxWidth: '200px',
    marginTop: '0.5rem',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  rjsfInfoModalForm: {
    marginLeft: '-1rem',
    marginTop: '-1rem',
    maxWidth: '39rem',
  },
  copyButton: {
    '&:hover': {
      color: theme.palette.secondary.link2,
    },
  },
  submitButton: {
    backgroundColor: theme.palette.secondary.focused,
    color: '#fff',
  },
  visibilityGridItem: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
  },
  img: {
    marginLeft: '0.5rem',
    marginTop: '-0.35rem',
    filter: theme.palette.secondary.img,
    height: '2rem',
    width: '4.44rem',
  },
  chip: {
    padding: '1rem .5rem',
    // backgroundColor: "red",
  },
  chipIcon: {
    width: '1rem',
  },
  // TODO this style needs to be refactored and applied to all modal globally
  closing: {
    transform: 'rotate(-90deg)',
    '&:hover': {
      transform: 'rotate(90deg)',
      transition: 'all .3s ease-in',
      cursor: 'pointer',
    },
  },
}));

export default useStyles;
