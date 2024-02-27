import { withStyles } from '@material-ui/styles';
import { Tooltip } from '@material-ui/core';
import { ziCalc } from '../../../utils/zIndex';
import { makeStyles } from '@material-ui/core/styles';
import { Colors } from '@/themes/app';

const useStyles = makeStyles((theme) => ({
  toolTip: {
    textDecoration: 'underline',
    color: theme.palette.secondary.link2,
    cursor: 'pointer',
  },
}));

export const RenderTooltipContent = ({ showPriortext, showAftertext, link }) => {
  const classes = useStyles();

  const handleClick = (e) => {
    window.open(link, '_blank');
    e.stopPropagation();
  };

  return (
    <div>
      <span style={{ marginRight: '2px' }}>{showPriortext}</span>
      <a
        onClick={handleClick}
        className={classes.toolTip}
        target="_blank"
        rel="noopener noreferrer"
      >
        Read docs
      </a>
      <span style={{ marginLeft: '2px' }}>{showAftertext}</span>
    </div>
  );
};

export const CustomTextTooltip = ({ backgroundColor = Colors.charcoal, flag, ...props }) => {
  const CustomTooltip = withStyles(() => ({
    tooltip: {
      backgroundColor: backgroundColor,
      color: '#fff',
      opacity: '100%',
      fontSize: '0.75rem',
      fontFamily: flag ? 'Qanelas Soft, sans-serif' : 'inherit',
      borderRadius: '0.9375rem',
      padding: '0.9rem',
      zIndex: ziCalc(11),
    },
    popper: {
      zIndex: `${ziCalc(5)} !important`,
    },
  }))(Tooltip);

  return <CustomTooltip {...props} />;
};
