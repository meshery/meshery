import React, { useState } from "react";
import IconButton from "@material-ui/core/IconButton";
import { connect } from "react-redux";
import NoSsr from "@material-ui/core/NoSsr";
import { Drawer, Tooltip, Divider, ClickAwayListener, Box, Typography, makeStyles, alpha } from "@material-ui/core";
import Filter from "./filter";
import BellIcon from "../../assets/icons/BellIcon.js"
import ErrorIcon from "../../assets/icons/ErrorIcon.js"
import { loadEventsFromPersistence, toggleNotificationCenter, updateEvents } from "../../lib/store";
import { iconMedium } from "../../css/icons.styles";
import { bindActionCreators } from "redux";
import { NOTIFICATIONCOLORS } from "../../themes";
import AlertIcon from "../../assets/icons/AlertIcon";
import ArchiveIcon from "../../assets/icons/ArchiveIcon";

const useStyles = makeStyles((theme) => ({
  sidelist : { width : "45rem" },
  notificationButton : { height : "100%" },
  notificationDrawer : {
    backgroundColor : theme.palette.secondary.sideBar,
    display : "flex",
    flexDirection : "column",
    justifyContent : "space-between",
  },
  drawerButton : {
    padding : "0.45rem",
    margin : "0.2rem",
    backgroundColor : theme.palette.secondary.dark,
    color : "#FFFFFF",
    "&:hover" : {
      backgroundColor : "#FFFFFF",
      color : theme.palette.secondary.dark,
    },
  },
  fullView : {
    right : 0,
    transition : "0.3s ease-in-out !important",
  },
  peekView : {
    right : "-42.1rem",
    transition : "0.3s ease-in-out !important",
  },

  container : {
    padding : "20px"
  },
  header : {
    marginBottom : "20px",
    display : "flex",
    gap : "0.5rem",
    justifyContent : "space-between",
    alignItems : "center",
  },
  title : {
    display : "flex",
    alignItems : "center",
    gap : "0.5rem",
  },
  titleBellIcon : {
    width : "36px",
    height : "36px",
    borderRadius : "100%",
    backgroundColor : "black",
    display : "flex",
    padding : "0.2rem",
    justifyContent : "center",
    alignItems : "center"
  },
  severityChip : {
    borderRadius : "4px",
    display : "flex",
    gap : "4px",
    // alignItems: "center",
    padding : "4px 12px",
    fontSize : "16px",
  },

  severityChips : {
    display : "flex",
    gap : "12px",
    alignItems : "center",

  },
  notification : {
    margin : theme.spacing(0.5, 1),
  },
}));

const SEVERITY = {
  INFO : "info",
  ERROR : "error",
  WARNING : "warning",
  // SUCCESS: "success"
}

const STATUS = {
  ACKNOWLEDGED : "acknowledged"
}
const STATUS_STYLE = {
  [STATUS.ACKNOWLEDGED] : {
    icon : ArchiveIcon,
    color : "#3c494f" // Charcoal
  }
}

const SEVERITY_STYLE = {
  [SEVERITY.INFO] : {
    icon : ErrorIcon,
    color : NOTIFICATIONCOLORS.INFO
  },
  [SEVERITY.ERROR] : {
    icon : ErrorIcon,
    color : NOTIFICATIONCOLORS.ERROR
  },
  [SEVERITY.WARNING] : {
    icon : AlertIcon,
    color : NOTIFICATIONCOLORS.WARNING
  },

}


const NotificationCountChip = ({ classes, notificationStyle, count }) => {
  const chipStyles = {
    fill : notificationStyle.color,
    height : "20px",
    width : "20px",
  }
  count = Number(count).toLocaleString('en', { useGrouping : true })
  return (
    <div className={classes.severityChip} style={{ backgroundColor : alpha(chipStyles.fill, 0.20) }} >
      {<notificationStyle.icon {...chipStyles} />}
      {count}
    </div>
  )
}

const Header = () => {
  const classes = useStyles()
  return (
    <div className={classes.container}>
      <Box className={classes.header}>
        <div className={classes.title}>
          <div className={classes.titleBellIcon}>
            <BellIcon height="30" width="30" fill="#fff" />
          </div>
          <Typography variant="h6"> Notifications</Typography>
        </div>
        <div className={classes.severityChips}>
          {Object.values(SEVERITY).map(severity =>
            <NotificationCountChip key={severity} classes={classes} notificationStyle={SEVERITY_STYLE[severity]} count="1000" />
          )}
          {Object.values(STATUS).map(status =>
            <NotificationCountChip key={status} classes={classes} notificationStyle={STATUS_STYLE[status]} count="1000" />
          )}

        </div>
      </Box>
      <Filter></Filter>
    </div>
  )
}


const MesheryNotification = (props) => {
  const [anchorEl, setAnchorEl] = useState(null);

  const handleToggle = () => {
    props.toggleOpen();
  };

  const handleClose = () => {
    if (!props.showFullNotificationCenter) {
      return;
    }
    props.toggleOpen();
    setAnchorEl(null);
  };
  const classes = useStyles()
  const { showFullNotificationCenter } = props;
  const open = Boolean(anchorEl) || showFullNotificationCenter;

  return (
    <NoSsr>
      <div>
        <Tooltip title={"errors"}>
          <IconButton
            id="notification-button"
            className={classes.notificationButton}
            color="inherit"
            onClick={handleToggle}
            onMouseOver={(e) => {
              e.preventDefault();
              setAnchorEl(e.currentTarget);
            }}
            onMouseLeave={(e) => {
              e.preventDefault();
              setAnchorEl(null);
            }}
          >
            <BellIcon className={iconMedium} fill="#fff" />
            {/* <Badge id="notification-badge" badgeContent={getNotificationCount(events)} color={newNotificationsType}>
              </Badge> */}
          </IconButton>
        </Tooltip>
      </div>

      <ClickAwayListener
        onClickAway={(e) => {
          if (
            e.target.className.baseVal !== "" &&
            e.target.className.baseVal !== "MuiSvgIcon-root" &&
            (typeof e.target.className === "string" ? !e.target.className?.includes("MesheryNotification") : null)
          ) {
            handleClose();
          }
        }}
      >
        <Drawer
          anchor="right"
          variant="persistent"
          open={open}
          classes={{
            paper : classes.notificationDrawer,
            paperAnchorRight : showFullNotificationCenter ? classes.fullView : classes.peekView,
          }}
        >
          <div>
            <div>
              <div className={classes.sidelist}>
                <Header></Header>

                <Divider light />
              </div>
            </div>
          </div>
        </Drawer>
      </ClickAwayListener>
    </NoSsr>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateEvents : bindActionCreators(updateEvents, dispatch),
  toggleOpen : bindActionCreators(toggleNotificationCenter, dispatch),
  loadEventsFromPersistence : bindActionCreators(loadEventsFromPersistence, dispatch),
});

const mapStateToProps = (state) => {
  const events = state.get("events");
  return {
    user : state.get("user"),
    events : events.toJS(),
    openEventId : state.get("notificationCenter").get("openEventId"),
    showFullNotificationCenter : state.get("notificationCenter").get("showFullNotificationCenter"),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MesheryNotification);
