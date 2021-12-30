import React from "react";
import { makeStyles, withStyles } from "@material-ui/core/styles";
import MuiAccordion from "@material-ui/core/Accordion";
import MuiAccordionSummary from "@material-ui/core/AccordionSummary";
import MuiAccordionDetails from "@material-ui/core/AccordionDetails";
import DeleteIcon from "@material-ui/icons/Delete";
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
    justifyContent : "flex-end",
    '&$expanded' : {
      margin : '12px 0',
      justifyContent : "flex-end",
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

  return (
    <div className={classes.root}>
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
        <AccordionDetails >{props.children}</AccordionDetails>

      </Accordion>
    </div>
  );
}
