import React, { useEffect } from "react";
import { makeStyles } from "@material-ui/core/styles";
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
    marginBottom : "0rem"
  },
  heading : {
    fontSize : theme.typography.pxToRem(15),
    fontWeight : theme.typography.fontWeightRegular
  }
}));

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
      <MuiAccordion defaultExpanded elevation={0}>
        <MuiAccordionSummary
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
        </MuiAccordionSummary>
        <MuiAccordionDetails ref={accordionDetailsRef} >{props.children}</MuiAccordionDetails>

      </MuiAccordion>
    </div>
  );
}
