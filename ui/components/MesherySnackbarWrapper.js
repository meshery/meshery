import React, { useState } from "react";
import PropTypes from "prop-types";
import classNames from "classnames";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";
import ErrorIcon from "@material-ui/icons/Error";
import InfoIcon from "@material-ui/icons/Info";
import CloseIcon from "@material-ui/icons/Close";
import IconButton from "@material-ui/core/IconButton";
import { SnackbarContent } from 'notistack';
import WarningIcon from "@material-ui/icons/Warning";
import { withStyles } from "@material-ui/core/styles";
import classnames from 'classnames';
import Collapse from '@material-ui/core/Collapse';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardActions from '@material-ui/core/CardActions';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

const variantIcon = {
  success: CheckCircleIcon,
  warning: WarningIcon,
  error: ErrorIcon,
  info: InfoIcon,
};

const styles = (theme) => ({
  success: {
    color: "#6fbf73",
  },
  error: {
    color: "#ff1744",
  },
  info: {
    color: "#2196f3",
  },
  warning: {
    color: "#ffc400",
  },
  icon: {
    fontSize: 20,
  },
  iconVariant: {
    opacity: 0.9,
    marginRight: theme.spacing(1),
  },
  message: {
    display: "flex",
    alignItems: "center",
  },
  root: {
    [theme.breakpoints.up("sm")]: {
      minWidth: "344px !important",
    },
  },
  card: {
    backgroundColor: "rgba(50, 50, 50)",
    width: "100%",
  },
  actionRoot: {
    padding: "8px 8px 8px 16px",
  },
  icons: {
    marginLeft: "auto",
  },
  expand: {
    padding: "8px 8px",
    transform: "rotate(0deg)",
    transition: theme.transitions.create("transform", {
      duration: theme.transitions.duration.shortest,
    }),
  },
  expandOpen: {
    transform: "rotate(180deg)",
  },
  collapse: {
    padding: 16,
  },
  checkIcon: {
    fontSize: 20,
    color: "#b3b3b3",
    paddingRight: 4,
  },
  button: {
    padding: 0,
    textTransform: "none",
  },
});

function MesherySnackbarWrapper(props) {
  const { classes, className, message, onClose, variant, details } = props;
  const Icon = variantIcon[variant];

  const [expanded, setExpanded] = useState(false);

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  return (
    <SnackbarContent className={classes.root}>
      <Card className={classNames(classes.card, classes[variant], className)}>
        <CardActions classes={{ root: classes.actionRoot }}>
          <Typography variant="subtitle2">
            <div style={{ display: "flex", alignItems: "center" }}>
              <Icon className={classNames(classes.icon, classes.iconVariant)} />
              <div>{message}</div>
            </div>
          </Typography>
          <div className={classes.icons}>
            <IconButton
              aria-label="Show more"
              className={classnames(classes.expand, { [classes.expandOpen]: expanded })}
              onClick={handleExpandClick}
            >
              <ExpandMoreIcon style={{ color: "#fff" }} />
            </IconButton>
            <IconButton className={classes.expand} onClick={onClose}>
              <CloseIcon style={{ color: "#fff" }} />
            </IconButton>
          </div>
        </CardActions>
        <Collapse in={expanded} timeout="auto" unmountOnExit>
          <Paper className={classes.collapse}>
            <Typography gutterBottom>Details</Typography>
            {details}
          </Paper>
        </Collapse>
      </Card>
    </SnackbarContent>
  );
}

MesherySnackbarWrapper.propTypes = {
  classes: PropTypes.object.isRequired,
  className: PropTypes.string,
  message: PropTypes.node,
  onClose: PropTypes.func,
  variant: PropTypes.oneOf(["success", "warning", "error", "info"]).isRequired,
  details: PropTypes.string
};

export default withStyles(styles)(MesherySnackbarWrapper);
