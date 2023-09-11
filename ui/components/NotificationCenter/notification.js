import * as React from 'react';
import { Box, Button, Collapse, Divider, Grid, IconButton,  Popover, Typography, alpha, useTheme } from "@material-ui/core"
import { makeStyles } from "@material-ui/core"
import {  SEVERITY_STYLE } from "./constants"
import { iconLarge, iconMedium } from "../../css/icons.styles"
import { MoreVert } from "@material-ui/icons"
import { Avatar } from "@mui/material"
import FacebookIcon from "../../assets/icons/FacebookIcon"
import LinkedInIcon from "../../assets/icons/LinkedInIcon"
import TwitterIcon from "../../assets/icons/TwitterIcon"
import ShareIcon from "../../assets/icons/ShareIcon"
import DeleteIcon from "../../assets/icons/DeleteIcon"
import moment from 'moment';


const useStyles = makeStyles(() => ({
  root: (props) => ({
    width: "100%",
    borderRadius: "3px",
    border: `1px solid ${props.notificationColor}`,
    marginBlock: "8px",
  }),

  summary: (props) => ({
    paddingBlock: "8px",
    cursor: "pointer",
    backgroundColor: alpha(props.notificationColor, 0.20),
  }),

  gridItem: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },

  message: {
    overflow: "hidden",
    textOverflow: "ellipsis",
    whiteSpace: "nowrap",
    width: "100%",
  },
  expanded: {
    paddingBlock: "12px",
  },
  actorAvatar: {
    display: "flex",
    justifyContent: "center",
    alignItems: "start",
  },

  descriptionHeading: {
    fontSize: "16px",
    fontWeight: "bolder !important",
    textTransform: "uppercase",
  },


}))

const useMenuStyles = makeStyles((theme) => {
  return {
    paper: {
      color: theme.palette.secondary.iconMain,
      boxShadow: theme.shadows[4],
      borderRadius: "3px",
      paddingInline: "0.5rem",
      paddingBlock: "0.25rem",
      width: "200px",
    },

    list: {
      padding: "0.5rem",
      display: "flex",
      flexDirection: "column",
      gridGap: "0.5rem",
      marginBlock: "0.5rem",
      borderRadius: "4px",
      backgroundColor: theme.palette.secondary.honeyComb,
    },

    listItem: {
      display: "flex",
      gridGap: "0.5rem",
      alignItems: "center",
      // justifyContent: "center",
    },

  }
})


const formatTimestamp = (utcTimestamp) => {
  // const curretUtcTimestamp = moment.utc().valueOf()

  // const timediff = currentUtcTimestamp - utcTimestamp
  // if (timediff >= 24 * 60 * 60 * 1000) {
  return moment(utcTimestamp).local().format('MMM DD, YYYY')
  // }
  // return moment(utcTimestamp).fromNow()
}

function BasicMenu() {

  const classes = useMenuStyles()
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };
  const handleClose = () => {
    setAnchorEl(null);
  };
  const theme = useTheme()
  return (
    <div className="mui-fixed">
      <IconButton
        id="basic-button"
        aria-controls={open ? 'basic-menu' : undefined}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
        onClick={handleClick}
      >
        <MoreVert />
      </IconButton>
      <Popover
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'left',
        }}
      >
        <Box className={classes.paper}>
          <div className={classes.list}>
            <Box className={classes.listItem} >
              <ShareIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
              <Typography variant="body1" > Share </Typography>
            </Box>
            <Divider />
            <Box className={classes.listItem} >
              <IconButton style={{ paddingLeft: "0rem" }}>
                <FacebookIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />  </IconButton>
              <IconButton>
                <LinkedInIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
              </IconButton>
              <IconButton>
                <TwitterIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
              </IconButton>
            </Box>
          </div>

          <div className={classes.list}>
            <Button style={{ padding: "0rem" }}>
              <div style={{ display: "flex", gap: "0.5rem" }}>
                <DeleteIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                <Typography variant="body1" > Delete </Typography>
              </div>
            </Button>
          </div>
        </Box>
      </Popover>
    </div>
  );
}

export const Notification = ({ event }) => {
  const severityStyles = SEVERITY_STYLE[event.severity]
  const classes = useStyles({
    notificationColor: severityStyles.color
  })

  const [expanded, setExpanded] = React.useState(false)
  const handleExpandClick = (e) => {
    e.stopPropagation()
    setExpanded(!expanded);
  };


  return (
    <div className={classes.root}>
      <Grid container className={classes.summary} onClick={handleExpandClick} >
        <Grid item sm={1} className={classes.gridItem} >
          <severityStyles.icon {...iconLarge} fill={severityStyles.color} />
        </Grid>
        <Grid item sm={7} className={classes.gridItem} >
          <Typography variant="body1" className={classes.message}> {event.description} </Typography>
        </Grid>
        <Grid item sm={3} className={classes.gridItem} >
          <Typography variant="body1"> {formatTimestamp(event.update_at)} </Typography>
        </Grid>
        <Grid item sm={1} >
          <BasicMenu />
        </Grid>
      </Grid>
      <Collapse in={expanded}>
        <Grid container className={classes.expanded}>
          <Grid item sm={1} className={classes.actorAvatar} >
            <Avatar > S </Avatar>
          </Grid>
          <Grid item sm={10}>
            <Grid container  >
              <div>
                <Typography variant="h6" className={classes.descriptionHeading}  >
                  Details
                </Typography>
                <Typography variant="body1" >
                  {event.description}
                </Typography>
              </div>
              <Grid container spacing={4} style={{ marginTop: "0.5rem" }}>
                <Grid item sm={6}>
                  <Typography variant="h6" className={classes.descriptionHeading}  >
                    Probable Cause
                  </Typography>
                  <ol style={{ paddingInline: "12px" }}>
                    {[1, 2, 3].map((i) => <li key={i} >  <Typography variant="body1" >
                      Error Removing Cpx from Kubernetes Context id 23959 .
                    </Typography> </li>)}
                  </ol>
                </Grid>
                <Grid item sm={6} >
                  <Typography variant="h6" className={classes.descriptionHeading}  >
                    Suggested Remediation
                  </Typography>
                  <ol style={{ paddingInline: "12px" }}>
                    {[1, 2, 3].map((i) => <li key={i} >  <Typography variant="body1" >
                      Error Removing Cpx from Kubernetes Context id 23959 .
                    </Typography> </li>)}
                  </ol>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Collapse>
    </div >
  )

}

export default Notification