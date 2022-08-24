import React, { useState } from "react";
import PropTypes from "prop-types";
import classNames from "classnames";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import ErrorIcon from "@material-ui/icons/Error";
import InfoIcon from "@material-ui/icons/Info";
import DoneIcon from '@material-ui/icons/Done';
import IconButton from "@material-ui/core/IconButton";
import Grid from '@material-ui/core/Grid';
import { SnackbarContent } from 'notistack';
import WarningIcon from "@material-ui/icons/Warning";
import { withStyles } from "@material-ui/core/styles";
import Collapse from '@material-ui/core/Collapse';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import { eventTypes } from "../lib/event-types";

const variantIcon = {
  success : CheckCircleIcon,
  warning : WarningIcon,
  error : ErrorIcon,
  info : InfoIcon,
};

const variantHoverColor = {
  success : "iconSuccess",
  warning : "iconWarning",
  error : "iconError",
  info : "iconInfo",
}

const styles = (theme) => ({
  success : { color : "#6fbf73", },
  error : { color : "#ff1744", },
  info : { color : "#2196f3", },
  warning : { color : "#ffc400", },
  iconColor : { color : "rgba(102, 102, 102, 1)" },
  iconSuccess : { "&:hover" : { color : "#6fbf73" } },
  iconError : { "&:hover" : { color : "#ff1744" } },
  iconInfo : { "&:hover" : { color : "#2196f3" } },
  iconWarning : { "&:hover" : { color : "#ffc400" } },
  icon : { fontSize : 20, },
  iconVariant : { opacity : 0.9,
    marginRight : theme.spacing(1), },
  message : { display : "flex",
    alignItems : "center", },
  snackbarContent : { [theme.breakpoints.up("sm")] : { minWidth : "344px !important", }, },
  card : { backgroundColor : "rgba(50, 50, 50)",
    width : "100%", },
  actionRoot : { padding : "8px 8px 8px 16px", },
  icons : { marginLeft : "auto", },
  expand : { padding : "8px 8px",
    transform : "rotate(0deg)",
    transition : theme.transitions.create("transform", { duration : theme.transitions.duration.shortest, }), },
  expandOpen : { transform : "rotate(180deg)", },
  collapse : { padding : 16, },
  checkIcon : { fontSize : 20,
    color : "#b3b3b3",
    paddingRight : 4, },
  button : { padding : 0,
    textTransform : "none", },
});

function MesherySnackbarWrapper(props) {
  const {
    classes, className, message, onClose, variant, details, cause, remedy, errorCode, componentType, componentName
  } = props;
  const Icon = variantIcon[variant];
  const ERROR_DOC_LINK = "https://docs.meshery.io/reference/error-codes"
  const [expanded, setExpanded] = useState(false);
  const [cardHover, setCardHover] = useState(false)

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  return (
    <SnackbarContent className={classes.snackbarContent}>
      <Card className={classNames(classes.card, classes[variant], className)}
        aria-label="Show more"
        onMouseEnter={() => setCardHover(true)}
        onMouseLeave={() => setCardHover(false)}
      >
        <CardActions classes={{ root : classes.actionRoot }} onClick={handleExpandClick}>
          <Grid container direction="row" justify="space-between" alignItems="center" wrap="nowrap">
            <Typography variant="subtitle2">
              <div style={{ display : "flex", alignItems : "center" }}>
                <Icon className={classNames(classes.icon, classes.iconVariant)} />
                <div>{message}</div>
              </div>
            </Typography>
            <Grid container item xs={4} className={classes.icons} justify="flex-end">
              <IconButton
                aria-label="Show more"
                className={classNames(classes.expand, { [classes.expandOpen] : expanded })}
                onClick={handleExpandClick}
              >
                <ExpandMoreIcon className={classNames({ [classes.iconColor] : !cardHover, [classes[variant]] : cardHover })} />
              </IconButton>

              <IconButton className={classes.expand} onClick={onClose}>
                <DoneIcon
                  className={classNames(classes.iconColor, classes[variantHoverColor[variant]]  )}
                  onMouseEnter={() => setCardHover(false)}
                  onMouseLeave={() => setCardHover(true)}
                />
              </IconButton>
            </Grid>
          </Grid>
        </CardActions>
        <Collapse in={expanded} timeout="auto" unmountOnExit >
          <Paper className={classes.collapse} square variant="outlined" elevation={0}>
            <Typography variant="subtitle2" gutterBottom>DETAILS</Typography>
            {details}
          </Paper>
          {variant === eventTypes[2].type &&
           <>
             <Paper className={classes.collapse} square variant="outlined" elevation={0}>
               <Typography variant="subtitle2" gutterBottom>PROBABLE CAUSE</Typography>
               {cause}
             </Paper>
             <Paper className={classes.collapse} square variant="outlined" elevation={0}>
               <Typography variant="subtitle2" gutterBottom>SUGGESTED REMEDIATION</Typography>
               {remedy}
             </Paper>
             <Paper className={classes.collapse} square variant="outlined" elevation={0}>
               <Typography variant="subtitle2" gutterBottom>ERROR CODE</Typography>
               <a href={`${ERROR_DOC_LINK}#meshery-${componentType}-for-meshery-${componentName.toLowerCase()}`} target="_blank" rel="referrer noreferrer"> {errorCode} </a>
             </Paper>
           </>
          }
        </Collapse>
      </Card>
    </SnackbarContent>
  );
}

MesherySnackbarWrapper.propTypes = {
  classes : PropTypes.object.isRequired,
  className : PropTypes.string,
  message : PropTypes.node,
  onClose : PropTypes.func,
  variant : PropTypes.oneOf(["success", "warning", "error", "info"]).isRequired,
  details : PropTypes.string
};

export default withStyles(styles)(MesherySnackbarWrapper);
