import { Box, styled, Button, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';

export const CreatAtContainer = styled(Typography)(({ isBold }) => ({
  fontWeight: isBold ? 'bold' : '',
  whiteSpace: 'wrap',
  fontSize: '0.8rem',
}));

export const ActionContainer = styled(Box)({
  width: '100%',
  display: 'flex',
  gap: '1rem',
  justifyContent: 'end',
});

export const CopyLinkButton = styled(Button)(({ theme }) => ({
  color: theme.palette.secondary.primaryModalText,
}));

export const VisibilityTag = styled(`div`)(({ theme }) => ({
  border: `0.5px solid ${theme.palette.secondary.textMain}`,
  color: theme.palette.secondary.textMain,
  borderRadius: '0.2rem',
  paddingInline: '0.2rem',
  paddingBlock: '0.2rem',
  fontSize: '0.7rem',
  textTransform: 'uppercase',
  textAlign: 'center',
  marginLeft: '0.5rem',
}));

export const ResourceName = styled(Typography)(() => ({
  fontFamily: 'Qanelas Soft, sans-serif',
  textAlign: 'left',
  marginTop: '0.5rem',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  maxWidth: '10rem',
}));

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
    width: '2rem',
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
