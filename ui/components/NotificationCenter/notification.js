import * as React from 'react';
import { Box, Button, Collapse, Divider, Grid, IconButton, Popover, Typography, alpha, useTheme } from "@material-ui/core"
import { makeStyles } from "@material-ui/core"
import { SEVERITY_STYLE, STATUS } from "./constants"
import { iconLarge, iconMedium } from "../../css/icons.styles"
import { MoreVert } from "@material-ui/icons"
import { Avatar } from "@mui/material"
import FacebookIcon from "../../assets/icons/FacebookIcon"
import LinkedInIcon from "../../assets/icons/LinkedInIcon"
import TwitterIcon from "../../assets/icons/TwitterIcon"
import ShareIcon from "../../assets/icons/ShareIcon"
import DeleteIcon from "../../assets/icons/DeleteIcon"
import moment from 'moment';
import { useUpdateStatusMutation, useDeleteEventMutation } from "../../rtk-query/notificationCenter"
import { useDispatch } from 'react-redux';
import { changeEventStatus, deleteEvent } from '../../store/slices/events';

const useStyles = makeStyles(() => ({
  root: (props) => ({
    width: "100%",
    borderRadius: "3px",
    border: `1px solid ${props.notificationColor}`,
    borderLeftWidth: props.status === STATUS.READ ? "4px" : "1px",
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

    button: {
      padding: "0rem",
      display: "flex",
      gridGap: "0.5rem",
      alignItems: "center",
      justifyContent: "start",
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

function BasicMenu({ event }) {

  const classes = useMenuStyles()
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);


  const handleClick = (event) => {
    event.stopPropagation()
    setAnchorEl(event.currentTarget);
  };
  const handleClose = (e) => {
    e.stopPropagation()
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

          <DeleteEvent event={event} />
          <ChangeStatus event={event} />
        </Box>
      </Popover>
    </div>
  );
}

export const DeleteEvent = ({ event }) => {

  const classes = useMenuStyles()
  const dispatch = useDispatch()
  const [deleteEventMutation] = useDeleteEventMutation()
  const theme = useTheme()
  const handleDelete = (e) => {
    e.stopPropagation()
    dispatch(deleteEvent(deleteEventMutation, event.id))
  }
  return (
    <div className={classes.list}>
      <Button className={classes.button} onClick={handleDelete}>
        <DeleteIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
        <Typography variant="body1" > Delete </Typography>
      </Button>
    </div>
  )

}



export const ChangeStatus = ({ event }) => {

  const classes = useMenuStyles()
  const newStatus = event.status === STATUS.READ ? STATUS.UNREAD : STATUS.READ
  const [updateStatusMutation] = useUpdateStatusMutation()

  const dispatch = useDispatch()

  const updateStatus = (e) => {
    e.stopPropagation()
    dispatch(changeEventStatus(updateStatusMutation, event.id, newStatus))
  }
  return (
    <div className={classes.list}>
      <Button className={classes.button} onClick={updateStatus}>
        <Typography variant="body1" > Mark As Read  </Typography>
      </Button>
    </div>
  )

}

const BulletList = ({ items }) => {
  return <ol style={{ paddingInline: "12px", paddingBlock: "0.3rem", margin: "0rem" }}>
    {[items].map((i) => <li key={i} >
      <Typography variant="body1" > {i} </Typography>
    </li>)}
  </ol>
}

export const Notification = ({ event }) => {
  const severityStyles = SEVERITY_STYLE[event.severity]
  const classes = useStyles({
    notificationColor: severityStyles.color,
    status: event.status
  })
  const [expanded, setExpanded] = React.useState(false)
  const handleExpandClick = (e) => {
    e.stopPropagation()
    setExpanded(!expanded);
  };

  const longDescription = event?.metadata?.error?.LongDescription || []
  const probableCause = event?.metadata?.error?.ProbableCause || []
  const suggestedRemediation = event?.metadata?.error?.SuggestedRemediation || []


  return (
    <div className={classes.root}>
      <Grid container className={classes.summary} onClick={handleExpandClick} >
        <Grid item sm={1} className={classes.gridItem} >
          <severityStyles.icon {...iconLarge} fill={severityStyles.color} />
        </Grid>
        <Grid item sm={7} className={classes.gridItem} >
          <Typography variant="body1" className={classes.message}> {event.description}   </Typography>
        </Grid>
        <Grid item sm={3} className={classes.gridItem} >
          <Typography variant="body1"> {formatTimestamp(event.update_at)} </Typography>
        </Grid>
        <Grid item sm={1} >
          <Box>
            <BasicMenu event={event} />
          </Box>
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
                <NestedData classes={classes} heading="Description" data={event.description} />
                <div style={{ marginTop: "0.3rem" }}>
                  <NestedData classes={classes} heading="Details" data={longDescription} />
                </div>
              </div>
              <Grid container spacing={1} style={{ marginTop: "0.5rem" }}>
                <Grid item sm={suggestedRemediation?.length > 0 ? 6 : 12}>
                  <NestedData classes={classes} heading="Probable Cause" data={probableCause} />
                </Grid>
                <Grid item sm={probableCause?.length > 0 ? 6 : 12} >
                  <NestedData classes={classes} heading="Suggested Remediation" data={suggestedRemediation} />
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Collapse>
    </div >
  )

}

const NestedData = ({ heading, data, classes }) => {

  if (!data || data?.length == 0) return null
  return (
    <>
      <Typography variant="h6" className={classes.descriptionHeading}  >
        {heading}
      </Typography>
      {typeof data === "string" ?
        <Typography variant="body1" >
          {data}
        </Typography> :
        <BulletList items={data} />
      }
    </>
  )
}



export default Notification