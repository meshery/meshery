import { makeStyles } from '@material-ui/core/styles';
import { styled } from '@mui/system';
import Tooltip, { tooltipClasses } from '@mui/material/Tooltip';

export const CustomTooltip = styled(({ className, ...props }) => (
  <Tooltip {...props} classes={{ popper: className }} />
))({
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: 'rgba(122,132,142,1)',
    color: '#F5F5F5',
    padding: '1rem',
    fontSize: '0.925rem',
    '& .tooltip-dark': {
      fontWeight: 'bold',
      fontSize: '1rem',
    },
  },
});

const useStyles = makeStyles((theme) => ({
  dialogTitle: {
    backgroundColor: theme.palette.secondary.mainBackground,
    justifyContent: 'space-between',
    flexDirection: 'row',
    alignItems: 'center',
    padding: '12px 20px',
    gap: '146px',
    color: theme.palette.secondary.white,
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
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  rjsfInfoModalForm: {
    marginLeft: '-1rem',
    marginTop: '-1rem',
    maxWidth: '44rem',
  },
  copyButton: {
    '&:hover': {
      color: theme.palette.secondary.link2,
    },
  },
  submitButton: {
    backgroundColor: theme.palette.secondary.focused,
    color: theme.palette.secondary.white,
  },
  visibilityGridItem: {
    display: 'flex',
    alignItems: 'center',
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
  modalActions: {
    padding: '1rem 1.5rem',
    display: 'flex',
    justifyContent: 'flex-end',
    width: '100%',
    gap: '10px',
  },
}));

export default useStyles;
