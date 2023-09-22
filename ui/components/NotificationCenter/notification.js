import * as React from 'react';
import { Avatar, Box, Button, Collapse, Grid, Hidden, IconButton, Popover, Slide, Tooltip, Typography, alpha, useTheme } from "@material-ui/core"
import { makeStyles } from "@material-ui/core"
import { SEVERITY_STYLE, STATUS } from "./constants"
import { iconLarge, iconMedium } from "../../css/icons.styles"
import { MoreVert } from "@material-ui/icons"
import FacebookIcon from "../../assets/icons/FacebookIcon"
import LinkedInIcon from "../../assets/icons/LinkedInIcon"
import TwitterIcon from "../../assets/icons/TwitterIcon"
import ShareIcon from "../../assets/icons/ShareIcon"
import DeleteIcon from "../../assets/icons/DeleteIcon"
import moment from 'moment';
import { useUpdateStatusMutation, useDeleteEventMutation } from "../../rtk-query/notificationCenter"
import { useDispatch, useSelector } from 'react-redux';
import { changeEventStatus, deleteEvent, selectEventById } from '../../store/slices/events';
import { useGetUserByIdQuery } from '../../rtk-query/user';
import { FacebookShareButton, LinkedinShareButton, TwitterShareButton } from 'react-share';
import ReadIcon from '../../assets/icons/ReadIcon';
import UnreadIcon from '../../assets/icons/UnreadIcon';
import { ErrorBoundary, withErrorBoundary, withSuppressedErrorBoundary } from '../General/ErrorBoundary';

const useStyles = makeStyles(() => ({
  root : (props) => ({
    width : "100%",
    borderRadius : "0.25rem",
    border : `0.1rem solid ${props.notificationColor}`,
    borderLeftWidth : props.status === STATUS.UNREAD ? "0.25rem" : "0.1rem",
    marginBlock : "0.5rem",
  }),

  summary : (props) => ({
    paddingBlock : "0.5rem",
    cursor : "pointer",
    backgroundColor : alpha(props.notificationColor, 0.20),
  }),

  gridItem : {
    display : "flex",
    alignItems : "center",
    justifyContent : "center",
  },

  message : {
    overflow : "hidden",
    textOverflow : "ellipsis",
    whiteSpace : "nowrap",
    overflowWrap : "break-word",
    // max of min of 20rem or 50vw
    maxWidth : "min(25rem, 50vw)",
    width : "100%",
  },
  expanded : {
    paddingBlock : "0.75rem",
    paddingInline : "0.2rem",
  },
  actorAvatar : {
    display : "flex",
    justifyContent : "center",
    alignItems : "start",
  },

  descriptionHeading : {
    fontWeight : "bolder !important",
    textTransform : "uppercase",
    fontSize : "0.9rem",
  },


}))

const useMenuStyles = makeStyles((theme) => {
  return {
    paper : {
      color : theme.palette.secondary.iconMain,
      boxShadow : theme.shadows[4],
      borderRadius : "0.25",
      paddingInline : "0.5rem",
      paddingBlock : "0.25rem",
      width : "12.5rem",
    },

    list : {
      display : "flex",
      flexDirection : "column",
      gridGap : "0.5rem",
      marginBlock : "0.5rem",
      borderRadius : "0.25rem",
      backgroundColor : theme.palette.secondary.honeyComb,
      "&:hover" : {
        backgroundColor : alpha(theme.palette.secondary.link2, 0.25),
      },
    },

    listItem : {
      display : "flex",
      gridGap : "0.5rem",
      alignItems : "center",
      justifyContent : "space-around",
    },
    socialListItem : {
      display : "flex",
      backgroundColor : alpha(theme.palette.secondary.honeyComb, 1.25),
      alignItems : "center",
      justifyContent : "space-around",
      padding : ".65rem",
    },

    button : {
      height : "100%",
      width : "100%",
      display : "flex",
      alignItems : "center",
      justifyContent : "start",
    },

  }
})


const formatTimestamp = (utcTimestamp) => {
  const currentUtcTimestamp = moment.utc().valueOf()


  const timediff = currentUtcTimestamp - moment(utcTimestamp).valueOf()
  if (timediff >= 24 * 60 * 60 * 1000) {
    return moment(utcTimestamp).local().format('MMM DD, YYYY')
  }
  return moment(utcTimestamp).fromNow()
}

const  BasicMenu = withSuppressedErrorBoundary(({ event }) =>  {

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

  const [isSocialShareOpen, setIsSocialShareOpen] = React.useState(false)
  const toggleSocialShare = (e) => {
    e.stopPropagation()
    setIsSocialShareOpen(prev => !prev)
  }

  const theme = useTheme()
  return (
    <div className="mui-fixed" onClick={(e) => e.stopPropagation()} >
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
          vertical : 'bottom',
          horizontal : 'left',
        }}
      >
        <Box className={classes.paper}>
          <div className={classes.list}>
            <Box className={classes.listItem} sx={{ width : '100%' }} >
              <Button onClick={toggleSocialShare} className={classes.button} >
                <ShareIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                <Typography variant="body1" style={{ marginLeft : '0.5rem' }} >Share</Typography>
              </Button>
            </Box>
            <Collapse in={isSocialShareOpen}>
              <Box className={classes.socialListItem} >
                <FacebookShareButton  url={"https://meshery.io"} quote={event.description || ""} >
                  <FacebookIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                </FacebookShareButton>
                <LinkedinShareButton url={"https://meshery.io"} summary={event.description || ""} >
                  <LinkedInIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                </LinkedinShareButton>
                <TwitterShareButton url={"https://meshery.io"} title={event.description || ""} >
                  <TwitterIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
                </TwitterShareButton>
              </Box>
            </Collapse>
          </div>

          <DeleteEvent event={event} />
          <ChangeStatus event={event} />
        </Box >
      </Popover >
    </div >
  );
})

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
        <Typography variant="body1" style={{ marginLeft : '0.5rem' }} > Delete </Typography>
      </Button>
    </div>
  )

}


export const FormattedErrorMetadata = ({ metadata }) => {
  const classes = useStyles()
  const longDescription = metadata?.LongDescription || []
  const probableCause = metadata?.ProbableCause || []
  const suggestedRemediation = metadata?.SuggestedRemediation || []

  return (
    <Grid container  >
      <div>
        <div style={{ marginTop : "1rem" }}>
          <NestedData classes={classes} heading="Details" data={longDescription} />
        </div>
      </div>
      <Grid container spacing={1} style={{ marginTop : "0.5rem" }}>
        <Grid item sm={suggestedRemediation?.length > 0 ? 6 : 12}>
          <NestedData classes={classes} heading="Probable Cause" data={probableCause} />
        </Grid>
        <Grid item sm={probableCause?.length > 0 ? 6 : 12} >
          <NestedData classes={classes} heading="Suggested Remediation" data={suggestedRemediation} />
        </Grid>
      </Grid>
    </Grid>
  )
}

const METADATA_FORMATTER = {
  "error" : FormattedErrorMetadata,
}

// Maps the metadata to the appropriate formatter component
export const FormattedMetadata = ({ event }) => {
  if (!event || !event.metadata ) return null

  const metdataKeys = Object.keys(event.metadata)
  return metdataKeys.map((key) => {
    const formatter = METADATA_FORMATTER[key]
    if (!formatter) return null
    return <formatter key={key}  metadata={event.metadata[key]} event={event} />
  })
}



export const ChangeStatus = ({ event }) => {

  const classes = useMenuStyles()
  const newStatus = event.status === STATUS.READ ? STATUS.UNREAD : STATUS.READ
  const [updateStatusMutation] = useUpdateStatusMutation()
  const theme = useTheme()
  const dispatch = useDispatch()

  const updateStatus = (e) => {
    e.stopPropagation()
    dispatch(changeEventStatus(updateStatusMutation, event.id, newStatus))
  }
  return (
    <div className={classes.list}>
      <Button className={classes.button} onClick={updateStatus}>
        {newStatus === STATUS.READ ?
          <ReadIcon {...iconMedium} fill={theme.palette.secondary.iconMain} /> :
          <UnreadIcon {...iconMedium} fill={theme.palette.secondary.iconMain} />
        }
        <Typography variant="body1" style={{ marginLeft : '0.5rem' }}> Mark as {newStatus}  </Typography>
      </Button>
    </div>
  )

}

const BulletList = ({ items }) => {
  return <ol style={{ paddingInline : "0.75rem", paddingBlock : "0.3rem", margin : "0rem" }}>
    {items.map((i) => <li key={i} >
      <Typography variant="body1" > {i} </Typography>
    </li>)}
  </ol>
}



export const Notification = withErrorBoundary(({ event_id }) => {
  const event = useSelector(state => selectEventById(state, event_id))
  const isVisible = event.is_visible === undefined ? true : event.is_visible
  const severityStyles = SEVERITY_STYLE[event.severity]
  const classes = useStyles({
    notificationColor : severityStyles?.color,
    status : event?.status
  })
  const [expanded, setExpanded] = React.useState(false)
  const handleExpandClick = (e) => {
    e.stopPropagation()
    setExpanded(!expanded);
  };

  const { data : user } = useGetUserByIdQuery(event.user_id || "")

  const userName = `${user?.first_name || ""} ${user?.last_name || ""}`
  const userAvatarUrl = user?.avatar_url || ""

  const longDescription = event?.metadata?.error?.LongDescription || []
  const probableCause = event?.metadata?.error?.ProbableCause || []
  const suggestedRemediation = event?.metadata?.error?.SuggestedRemediation || []

  return (
    <Slide in={isVisible} timeout={250} direction="left" appear={false} enter={false} mountOnEnter unmountOnExit  >
      <div className={classes.root}>
        <Grid container className={classes.summary} onClick={handleExpandClick} >
          <Grid item sm={1} className={classes.gridItem} >
            <severityStyles.icon {...iconLarge} fill={severityStyles?.color} />
          </Grid>
          <Grid item xs={9} md={7} className={classes.gridItem} >
            <Typography variant="body1" className={classes.message}> {event.description}   </Typography>
          </Grid>
          <Hidden smDown>
            <Grid item sm={3} className={classes.gridItem} >
              <Typography variant="body1"> {formatTimestamp(event.created_at)} </Typography>
            </Grid>
          </Hidden>
          <Grid item sm={1} >
            <Box>
              <BasicMenu event={event} />
            </Box>
          </Grid>
        </Grid>
        <Collapse in={expanded}>
          <ErrorBoundary>

            <Grid container className={classes.expanded}>
              <Grid item sm={1} className={classes.actorAvatar} >
                <Box sx={{ display : "flex", gridGap : "0.5rem", flexDirection : { xs : "row", md : "column" } }} >

                  {event.user_id  && user &&
                    <Tooltip title={userName} placement="top" >
                      <Avatar alt={userName} src={userAvatarUrl} />
                    </Tooltip>
                  }
                  {event.system_id &&
                    <Tooltip title={`System ID: ${event.system_id}`} placement="top" >
                      <Avatar src="/static/img/meshery-logo.png" />
                    </Tooltip>
                  }

                </Box>
              </Grid>
              <Grid item sm={10}>
                <Grid container  >
                  <div>
                    <NestedData classes={classes} heading="Description" data={event.description} />
                    <div style={{ marginTop : "1rem" }}>
                      <NestedData classes={classes} heading="Details" data={longDescription} />
                    </div>
                  </div>
                  <Grid container spacing={1} style={{ marginTop : "0.5rem" }}>
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
          </ErrorBoundary>
        </Collapse>
      </div >
    </Slide>
  )

})
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