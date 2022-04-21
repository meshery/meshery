import React, { useEffect } from "react";
import { makeStyles, withStyles } from "@material-ui/core/styles";
import MuiAccordion from "@material-ui/core/Accordion";
import MuiAccordionSummary from "@material-ui/core/AccordionSummary";
import MuiAccordionDetails from "@material-ui/core/AccordionDetails";
import DeleteIcon from "@material-ui/icons/Delete";
import Typography from "@material-ui/core/Typography";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import { IconButton } from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  accordionRoot : {
    width : "100%",
    marginBottom : "1rem"
  },
  heading : {
    fontSize : theme.typography.pxToRem(15),
    fontWeight : theme.typography.fontWeightRegular
  }
}));


const Accordion = withStyles({
  root : {
    border : '1px solid rgba(0, 0, 0, .125)',
    boxShadow : 'none',
    '&:not(:last-child)' : {
      borderBottom : 0,
    },
    '&:before' : {
      display : 'none',
    },
    '&$expanded' : {
      margin : 'auto',
    },
  },
  expanded : {},
})(MuiAccordion);

const AccordionSummary = withStyles({
  root : {
    backgroundColor : 'rgba(0, 0, 0, .03)',
    borderBottom : '1px solid rgba(0, 0, 0, .125)',
    marginBottom : -1,
    maxHeight : "1.5rem",
    '&$expanded' : {
      minHeight : 56,
    },
  },
  content : {
    justifyContent : "space-between",
    '&$expanded' : {
      margin : '12px 0',
      justifyContent : "space-between",
    },
  },
  expanded : {},
})(MuiAccordionSummary);

const AccordionDetails = withStyles((theme) => ({
  root : {
    padding : theme.spacing(2),
  },
}))(MuiAccordionDetails);

export default function SimpleAccordion(props) {
  const classes = useStyles();
  const accordionDetailsRef = React.useRef(null);

  useEffect(() => {
    // for managing focus
    if (accordionDetailsRef.current) {
      accordionDetailsRef.current.scrollIntoView({
        behavior : "smooth",
        block : "nearest",
      });
    }

  }, [accordionDetailsRef.current]);

  return (
    <div className={classes.accordionRoot}>
      <Accordion defaultExpanded elevation={0}>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel1a-content"
          id="panel1a-header"
          style={{ backgroundColor : "F7F7F7" }}
        >

          <Typography className={classes.heading}>{props.heading}</Typography>

          {props.childProps.hasRemove && (

            <IconButton
              style={{ padding : "0" }}
              // style={btnStyle}
              disabled={props.childProps.disabled || props.childProps.readonly}
              onClick={props.childProps.onDropIndexClick(
                props.childProps.index
              )}
            >
              <DeleteIcon />
            </IconButton>
          )}
        </AccordionSummary>
        <AccordionDetails ref={accordionDetailsRef} >{props.children}</AccordionDetails>

      </Accordion>
    </div>
  );
}
