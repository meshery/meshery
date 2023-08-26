import React, { useEffect, useState } from "react";
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
import { EVENT_TYPES } from "../lib/event-types";
import ReplyIcon from '@material-ui/icons/Reply';
import {
  TwitterShareButton,
  LinkedinShareButton,
  FacebookShareButton,
  TwitterIcon,
  LinkedinIcon,
  FacebookIcon
} from "react-share"
import { ClickAwayListener, Fade, Popper } from "@material-ui/core";
import moment from "moment";

const variantIcon = {
  [EVENT_TYPES.SUCCESS.type] : CheckCircleIcon,
  [EVENT_TYPES.WARNING.type] : WarningIcon,
  [EVENT_TYPES.ERROR.type] : ErrorIcon,
  [EVENT_TYPES.DEFAULT.type] : InfoIcon,
  [EVENT_TYPES.INFO.type] : InfoIcon,
};

const variantHoverColor = {
  [EVENT_TYPES.SUCCESS.type] : "iconSuccess",
  [EVENT_TYPES.WARNING.type] : "iconWarning",
  [EVENT_TYPES.ERROR.type] : "iconError",
  [EVENT_TYPES.INFO.type] : "iconInfo",
  [EVENT_TYPES.DEFAULT.type] : "iconInfo"
}

const styles = (theme) => ({
  [EVENT_TYPES.SUCCESS.type] : { color : "#6fbf73", },
  [EVENT_TYPES.ERROR.type] : { color : "#ff1744", },
  [EVENT_TYPES.INFO.type] : { color : "#2196f3", },
  [EVENT_TYPES.WARNING.type] : { color : "#ffc400", },
  [EVENT_TYPES.DEFAULT.type] : { color : "#edeff1", },
  iconColor : { color : "rgba(102, 102, 102, 1)" },
  iconSuccess : { "&:hover" : { color : "#6fbf73" } },
  iconError : { "&:hover" : { color : "#ff1744" } },
  iconInfo : { "&:hover" : { color : "#2196f3" } },
  iconWarning : { "&:hover" : { color : "#ffc400" } },
  icon : { fontSize : 20, },
  iconVariant : {
    opacity : 0.9,
    marginRight : theme.spacing(1),
  },
  message : {
    display : "flex",
    alignItems : "center",
  },
  timestamp : {
    color : "#ebeff1",
    fontSize : "0.8rem",
    fontStyle : "italic",
  },
  snackbarContent : { [theme.breakpoints.up("sm")] : { minWidth : "344px !important", }, },
  snackbarContentBorder : {
    border : "1px solid rgba(102, 102, 102, 1)"
  },
  card : {
    backgroundColor : "rgba(50, 50, 50)",
    width : "100%",
  },
  actionRoot : { padding : "8px 8px 8px 16px", },
  icons : { marginLeft : "auto", },
  expand : {
    padding : "8px 8px",
    transform : "rotate(0deg)",
    transition : theme.transitions.create("transform", { duration : theme.transitions.duration.shortest, }),
  },
  expandOpen : { transform : "rotate(180deg)", },
  collapse : { padding : 16, },
  checkIcon : {
    fontSize : 20,
    color : "#b3b3b3",
    paddingRight : 4,
  },
  button : {
    padding : 0,
    textTransform : "none",
  },
  share : {
    transform : "scaleX(-1)"
  },
  popper : {
    width : 500,
  },
  paper : {
    padding : theme.spacing(1)
  },
  shareIcon : {
    margin : theme.spacing(0.4)
  }
});

const generateMsgForMesh = (name) => {
  return `I deployed ${name} service mesh with one-click using @mesheryio!\nManage your infrastructure with Meshery`
}

const generateMsgForSampleApp = (name) => {
  return `I deployed ${name} with one-click using @mesheryio!\nManage your infrastructure with Meshery`
}

const generateMsgForAppsPatt = (name) => {
  return `I deployed ${name} [design | application] in a single-click using @mesheryio!\nFind design patterns like mine in the Meshery Catalog - https://meshery.io/catalog`
}

const getDefaultMessage = (message) => {
  const msg = `" ${message} "
    Manage your infrastructure with Meshery
  `
  return msg
}

const formatTimestamp = (utcTimestamp ) => {
  const currentUtcTimestamp = moment.utc().valueOf()

  const timediff = currentUtcTimestamp - utcTimestamp
  if (timediff >= 24 * 60 *60 *1000) {
    return moment(utcTimestamp).local().format('YYYY-MM-DD HH:mm')
  }
  return moment(utcTimestamp).fromNow()
}

function MesherySnackbarWrapper(props) {
  const {
    classes, className, message, onClose, eventType, details, cause, remedy, errorCode, componentType, componentName, expand,timestamp
  } = props;
  const variant = eventType.type
  const Icon = variantIcon[variant];
  const ERROR_DOC_LINK = "https://docs.meshery.io/reference/error-codes"
  const [expanded, setExpanded] = useState(false);
  const [cardHover, setCardHover] = useState(false)
  const [socialExpand, setSocialExpand] = useState(false);
  const [anchorEl, setAnchorEl] = useState(null);
  const [socialMessage, setSocialMessage] = useState("");
  const [highlight, setHighlight] = useState(false)

  const handleExpandClick = () => {
    setExpanded(!expanded);
  };

  const handleSocialExpandClick = (e) => {
    setAnchorEl(e.currentTarget);
    e.stopPropagation();
    setSocialExpand(socialExpand => !socialExpand);
  }
  useEffect(() => {
    if (expand && !expanded) {
      handleExpandClick();
      setHighlight(true)
      setTimeout(() => {
        setHighlight(false)
      }, 3000)
    }
    if (componentType === "adapter") {
      if (message?.includes("mesh installed")) {
        setSocialMessage(generateMsgForMesh(componentName[0].toUpperCase() + componentName.substring(1).toLowerCase()))
        return
      }
      if (message.includes("application installed")) {
        const name = message?.split(" ")[0];
        setSocialMessage(generateMsgForSampleApp(name[0].toUpperCase() + name.substring(1).toLowerCase()))
        return
      }
      setSocialMessage(getDefaultMessage(message))
      return
    }
    if (componentType === "core" && message?.includes("deployed")) {
      const designName = message?.split(":")[1]
      setSocialMessage(generateMsgForAppsPatt(designName))
    }

  },[expand])
  return (
    <SnackbarContent className={classes.snackbarContent}>
      <Card className={highlight ? classNames(classes.card, classes[variant], className, classes.snackbarContentBorder) : classNames(classes.card, classes[variant], className)}
        aria-label="Show more"
      >
        <CardActions classes={{ root : classes.actionRoot }} onClick={handleExpandClick}>
          <Grid container direction="row" justify="space-between" alignItems="center" wrap="nowrap">
            <Typography variant="subtitle2">
              <div style={{ display : "flex", alignItems : "center" }}>
                <Icon className={classNames(classes.icon, classes.iconVariant)} />
                <div>
                  <p>{message}</p>
                  <p className={classes.timestamp}> {formatTimestamp(timestamp)} </p>
                </div>
              </div>
            </Typography>

            <Grid container item xs={4} className={classes.icons} justify="flex-end">
              <IconButton
                aria-label="Show more"
                className={classNames(classes.expand, { [classes.expandOpen] : expanded })}
                onClick={handleExpandClick}
                onMouseEnter={() => setCardHover(true)}
                onMouseLeave={() => setCardHover(false)}
              >
                <ExpandMoreIcon className={classNames({ [classes.iconColor] : !cardHover, [classes[variant]] : cardHover })} />
              </IconButton>

              {variant === EVENT_TYPES.SUCCESS.type &&
                <IconButton
                  aria-label="Share"
                  className={classes.expand}
                  onClick={(e) => handleSocialExpandClick(e)}
                >
                  <ReplyIcon className={classNames(classes.share, classes.iconColor, classes[variantHoverColor[variant]])}
                    onMouseEnter={() => setCardHover(false)}
                    onMouseLeave={() => setCardHover(true)}
                  />
                </IconButton>
              }
              <IconButton className={classes.expand} onClick={onClose}>
                <DoneIcon
                  className={classNames(classes.iconColor, classes[variantHoverColor[variant]])}
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
          {variant === EVENT_TYPES.ERROR.type  &&
            <>

              { cause &&
              <Paper className={classes.collapse} square variant="outlined" elevation={0}>
                <Typography variant="subtitle2" gutterBottom>PROBABLE CAUSE</Typography>
                {cause}
              </Paper>
              }
              { remedy &&
              <Paper className={classes.collapse} square variant="outlined" elevation={0}>
                <Typography variant="subtitle2" gutterBottom>SUGGESTED REMEDIATION</Typography>
                {remedy}
              </Paper>
              }
              {componentName &&
              <Paper className={classes.collapse} square variant="outlined" elevation={0}>
                <Typography variant="subtitle2" gutterBottom>ERROR CODE</Typography>
                <a href={`${ERROR_DOC_LINK}#meshery-${componentType}-for-meshery-${componentName.toLowerCase()}`} target="_blank" rel="referrer noreferrer"> {errorCode} </a>
              </Paper>}
            </>
          }
        </Collapse>
        <div className={classes.popper}>
          <Popper open={socialExpand} anchorEl={anchorEl} transition style={{ zIndex : "1301" }}>
            {({ TransitionProps }) => (
              <ClickAwayListener onClickAway={() => setSocialExpand(false)}>
                <Fade {...TransitionProps} timeout={350}>
                  <Paper className={classes.paper}>
                    <TwitterShareButton className={classes.shareIcon} url={"https://meshery.io"} title={socialMessage}
                      hashtags={["opensource"]}
                    >
                      <TwitterIcon  size={32} />
                    </TwitterShareButton>
                    <LinkedinShareButton className={classes.shareIcon} url={"https://meshery.io"} summary={socialMessage}>
                      <LinkedinIcon  size={32}  />
                    </LinkedinShareButton>
                    <FacebookShareButton className={classes.shareIcon} url={"https://meshery.io"} quote={socialMessage}
                      hashtag={"#opensource"}
                    >
                      <FacebookIcon  size={32}  />
                    </FacebookShareButton>
                  </Paper>
                </Fade>
              </ClickAwayListener>
            )}
          </Popper>
        </div>
      </Card>
    </SnackbarContent>
  );
}

MesherySnackbarWrapper.propTypes = {
  classes : PropTypes.object.isRequired,
  className : PropTypes.string,
  message : PropTypes.node,
  onClose : PropTypes.func,
  eventType : PropTypes.object.isRequired,
  details : PropTypes.string,
  expand : PropTypes.string.isRequired
};

export default withStyles(styles)(MesherySnackbarWrapper);