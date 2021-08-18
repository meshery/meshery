import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Accordion from "@material-ui/core/Accordion";
import AccordionSummary from "@material-ui/core/AccordionSummary";
import AccordionDetails from "@material-ui/core/AccordionDetails";
import DeleteIcon from "@material-ui/icons/Delete"
import Typography from "@material-ui/core/Typography";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import { IconButton } from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  root : {
    width : "100%",
    marginBottom : "1rem"
  },
  heading : {
    fontSize : theme.typography.pxToRem(15),
    fontWeight : theme.typography.fontWeightRegular
  }
}));

export default function SimpleAccordion(props) {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Accordion>
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          aria-controls="panel1a-content"
          id="panel1a-header"
        >
          <Typography className={classes.heading}>{props.heading}</Typography>

          {props.childProps.hasRemove && (

            <IconButton
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
        <AccordionDetails>{props.children}</AccordionDetails>
      </Accordion>
    </div>
  );
}
