import { ziCalc } from '../../../utils/zIndex';
import { makeStyles } from '@material-ui/core/styles';
import { styled } from '@material-ui/core';
import { CustomTooltip } from '@layer5/sistent';
import { UsesSistent } from '@/components/SistentWrapper';

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
        Learn more
      </a>
      <span style={{ marginLeft: '2px' }}>{showAftertext}</span>
    </div>
  );
};

export const CustomTextTooltip = ({ flag, ...props }) => {
  const StyledTooltip = styled(CustomTooltip)(() => ({
    tooltip: {
      fontFamily: flag ? 'Qanelas Soft, sans-serif' : 'inherit',
      zIndex: ziCalc(11),
    },
    popper: {
      zIndex: `${ziCalc(5)} !important`,
    },
  }));

  return (
    <UsesSistent>
      <StyledTooltip {...props} />
    </UsesSistent>
  );
};
