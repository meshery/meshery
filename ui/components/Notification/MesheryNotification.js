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
  } from "@mui/material";
  import NotificationsIcon from '@mui/icons-material/Notifications';
  import DoneAllIcon from '@mui/icons-material/DoneAll';
  import MesheryEventViewer from "./MesheryEventViewer";

   
  function MesheryNotification () {

  const DrawerButton = styled(IconButton)((theme) => ({
    padding : '0.45rem',
    margin : '0.2rem',
    backgroundColor : '#000000',
    color : '#FFFFFF',
    "&:hover" : { backgroundColor : '#FFFFFF',
      color : '#000000' }
  }))

    const listTop = styled('div')(() => ({
      display : 'grid',
    alignItems : 'center',
    gridTemplateColumns : "2fr 6fr 2fr",
    paddingTop : theme.spacing(2),
    paddingLeft : theme.spacing(1),
    paddingRight : theme.spacing(1),
    paddingBottom : theme.spacing(2),

    }))

    const NotifSelector = styled('div')(() => ({
      display : 'flex', 

    }))
    const ClearAllButton = styled('div')(() => ({
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

    const handleTabChange = ( newTabValue) => {
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
        <div>
        <IconButton color="inherit" size="large"  onClick={handleToggle}>
        <Badge badgeContent={getNotificationCount(events)} >
                <NotificationsIcon   />
                </Badge>  
            </IconButton>
                   
             <Drawer
             sx={{
                // width: 450,
                flexShrink: 0,
                '& .MuiDrawer-paper': {
                  width: 450,
                  background: "#fff",
                  boxSizing: 'border-box',
                },
            }}
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
            <div>
              <div >
                <div style={{      display : 'grid',
    alignItems : 'center',
    gridTemplateColumns : "2fr 6fr 2fr",
    padding: "1rem 0.5rem "
   }} >
                  <NotifSelector >
                    <Tooltip title="Show all notifications">
                      <DrawerButton
                        color="inherit"
                        // onClick={handleBellButtonClick}
                        onClick={handleClose}
                      >
                        <NotificationsIcon />
                      </DrawerButton>
                    </Tooltip>
                  </NotifSelector>
                  <div >
                    <Typography variant="subtitle1" align="center">
                      Notifications
                    </Typography>
                  </div>
                  <ClearAllButton>
                    <Tooltip title="Clear all notifications">
                      <DrawerButton
                        color="inherit"
                        // onClick={handleClearAllNotifications}
                      >
                        <DoneAllIcon />
                      </DrawerButton>
                    </Tooltip>
                  </ClearAllButton>
                </div>
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
                    // deleteEvent={self.deleteEvent(ind)}
                    eventDetails={event.details || "Details Unavailable"}
                  />
                ))}
              </div>
            </div>
            </Drawer>  
            
  
        </div>
    )
  }

  export default MesheryNotification;