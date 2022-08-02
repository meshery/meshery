import react, {useState} from "react";
import { styled } from "@mui/material/styles";
import {
    IconButton,
    Typography,
    Divider,
    Tooltip,
    Badge,
    Drawer,
    Tabs,
    Tab,
    Grid,
  } from "@mui/material";
  import NotificationsIcon from '@mui/icons-material/Notifications';
  import DoneAllIcon from '@mui/icons-material/DoneAll';
  import MesheryEventViewer from "./MesheryEventViewer";

   
  function MesheryNotification () {

  const DrawerButton = styled(IconButton)(({theme}) => ({
    padding : theme.spacing(1),
    margin : theme.spacing(0.5),
    backgroundColor : theme.palette.secondary.dark,
    color : theme.palette.white,
    "&:hover" : { backgroundColor : theme.palette.white,
      color : theme.palette.secondary.dark }
  }))
  const DrawerContainer = styled(Drawer)(({theme}) => ({
    flexShrink: 0,
    '& .MuiDrawer-paper': {
      width: 450,
      background: theme.palette.white,
      boxSizing: 'border-box',
    },
  }))

    const ListTopGrid = styled(Grid)(({theme}) => ({
    paddingTop : theme.spacing(2),
    paddingLeft : theme.spacing(1),
    paddingRight : theme.spacing(1),
    paddingBottom : theme.spacing(2),

    }))

    const ClearAllButton = styled(Grid)(() => ({
      display : 'flex', 
      justifyContent : 'flex-end',
    }))
    
    const [events, setEvents] = useState([]);
    const [open, setOpen] = useState(false);
    const [tabValue, setTabValue] = useState(0);
    const [displayEventType, setDisplayEventType] = useState('*');

    const handleToggle = () => {
        setOpen(true);
      };
    const handleClose = () => {
        setOpen(false);
      };
    
    const handleNotifFiltering = (type) => {
      // setDisplayEventType(type);
    }  

    const handleTabChange = (event ,newTabValue) => {
      setTabValue(newTabValue) ;
    }
  
    function getNotifications(events, type) {
      if (!Array.isArray(events)) return [];
    
      if (type === "error") return events.filter(ev => ev.event_type === 2);
      if (type === "warning") return events.filter(ev => ev.event_type === 1)
      if (type === "success") return events.filter(ev => ev.event_type === 0)
    
      return events;
    }

    function getNotificationCount(events) {
        if (!Array.isArray(events)) return 0;
      
        const errorEventCount = events.filter(ev => ev.event_type === 2).length;
        const totalEventsCount = events.length;
      
        return errorEventCount || totalEventsCount;
      }

    return(
        <>
        <IconButton color="inherit" size="large"  onClick={handleToggle}>
        <Badge badgeContent={getNotificationCount(events)} >
                <NotificationsIcon   />
                </Badge>  
            </IconButton>
                   
             <DrawerContainer
            PaperProps={{
              sx: {
                "& .MuiPaper-root": {
                width: "100%",
                }
              }
            }}
            variant="persistent"
            anchor="right"
            open={open}
            onClose={handleClose}
          >
                {/* <ListTop> */}
                <ListTopGrid container  spacing={{ xs: 2, md: 3 }} columns={{ xs: 4, sm: 8, md: 12 }}>
                <Grid item xs={2} sm={4} md={4} >
                    <Tooltip title="Show all notifications">
                      <DrawerButton
                        color="inherit"
                        // onClick={handleBellButtonClick}
                        onClick={handleClose}
                      >
                        <NotificationsIcon />
                      </DrawerButton>
                    </Tooltip>
                </Grid>

                <Grid item xs={2} sm={4} md={4} >
                    <Typography variant="subtitle1" align="center">
                      Notifications
                    </Typography>
                  </Grid>
                  <ClearAllButton item xs={2} sm={4} md={4} > 
                    <Tooltip title="Clear all notifications">
                      <DrawerButton
                        color="inherit"
                        // onClick={handleClearAllNotifications}
                      >
                        <DoneAllIcon />
                      </DrawerButton>
                    </Tooltip>
                  </ClearAllButton>
                </ListTopGrid>
                <Divider light />
                <Tabs
                  value={tabValue}
                  onChange={handleTabChange}
                  indicatorColor="primary"
                  textColor="primary"
                  variant="fullWidth"
                >
                  <Tab label="All"  onClick={handleNotifFiltering('*')} style={{ minWidth : "15%" }}/>
                  <Tab label="Error"  onClick={handleNotifFiltering('error')} style={{ minWidth : "15%" }}/>
                  <Tab label="Warning" onClick={handleNotifFiltering('warning')} style={{ minWidth : "15%" }}/>
                  <Tab label="Success" onClick={handleNotifFiltering('success')} style={{ minWidth : "15%" }}/>
                </Tabs> 
                 {getNotifications(events, displayEventType).map((event, ind) => (
                  <MesheryEventViewer
                    eventVariant={event.event_type}
                    eventSummary={event.summary}
                    eventDetails={event.details || "Details Unavailable"}
                  />
                ))}
            </DrawerContainer> 
        </>
    )
  }

  export default MesheryNotification;