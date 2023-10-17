import React, { useEffect } from 'react';
import MuiAccordion from '@material-ui/core/Accordion';
import MuiAccordionSummary from '@material-ui/core/AccordionSummary';
import MuiAccordionDetails from '@material-ui/core/AccordionDetails';
import DeleteIcon from '../../../../assets/icons/DeleteIcon';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import { IconButton, makeStyles } from '@material-ui/core';
import { iconSmall } from '../../../../css/icons.styles';

const useStyles = makeStyles((theme) => ({
  accordionRoot: {
    width: '100%',
    marginBottom: '0rem',
  },
  heading: {
    fontSize: theme.typography.pxToRem(15),
    fontWeight: theme.typography.fontWeightRegular,
  },
  typography: {
    fontSize: '0.8rem',
  },
}));

export default function SimpleAccordion(props) {
  useEffect(() => {
    // for managing focus
    if (accordionDetailsRef.current) {
      accordionDetailsRef.current.scrollIntoView({
        behavior: 'smooth',
        block: 'nearest',
      });
    }
  }, []);
  // React Hook useEffect has an unnecessary dependency: 'accordionDetailsRef.current'. Either exclude it or remove the dependency array. Mutable values like 'accordionDetailsRef.current' aren't valid dependencies because mutating them doesn't re-render the component.

  const classes = useStyles();
  const accordionDetailsRef = React.useRef(null);

  return (
    <div className={classes.accordionRoot}>
      <MuiAccordion defaultExpanded elevation={0}>
        <MuiAccordionSummary
          expandIcon={<ExpandMoreIcon style={iconSmall} fill="gray" />}
          aria-controls="panel1a-content"
          id="panel1a-header"
          style={{ backgroundColor: 'F7F7F7' }}
        >
          <Typography className={classes.heading}>
            {props.heading?.charAt(0).toUpperCase() + props.heading?.slice(1)}{' '}
          </Typography>

          {props.childProps.hasRemove && (
            <IconButton
              style={{ padding: '0', iconSmall }}
              // style={btnStyle, iconSmall}
              disabled={props.childProps.disabled || props.childProps.readonly}
              onClick={props.childProps.onDropIndexClick(props.childProps.index)}
            >
              <DeleteIcon fill="gray" style={iconSmall} />
            </IconButton>
          )}
        </MuiAccordionSummary>
        <MuiAccordionDetails ref={accordionDetailsRef}>{props.children}</MuiAccordionDetails>
      </MuiAccordion>
    </div>
  );
}
