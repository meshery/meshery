import React, { useState, useEffect } from "react";
import Search from "@mui/icons-material/Search";
import AddIcon from '@mui/icons-material/Add';
import {Paper, Slide, Button, IconButton, Link, TextField, Checkbox, Avatar} from '@mui/material';
import ClickAwayListener from '@mui/base/ClickAwayListener';
import BadgeAvatars from "@/components/CustomAvatar"
import { styled } from "@mui/material/styles";
import { useTheme } from "@mui/system";


const CbadgeContainer = styled('div')(({ theme }) => ({
  display : "flex",
  justifyContent : "center",
  alignItems : "center",
  position : "relative"
}));

const CMenuContainer = styled(Paper)(({ theme }) => ({
  backgroundColor : "revert",
  borderRadius : "3px",
  padding : theme.spacing(3),
  zIndex : 1201,
  boxShadow : "20px #979797",
  transition : "linear .2s",
  transitionProperty : "height",
  alignItems: "center"
}));

const Cbadge = styled('div')(({ theme }) => ({
  fontSize : theme.spacing(1.5),
  backgroundColor : "white",
  borderRadius : "50%",
  color : "black",
  height : theme.spacing(2.5),
  width : theme.spacing(2.5),
  display : "flex",
  justifyContent : "center",
  alignItems : "center",
  position : "absolute",
  zIndex : 1,
  right : "-0.75rem",
  top : "-0.29rem"
}));


export default function K8sContextMenu({
  contexts={},  setActiveContexts = () => { },
  searchContexts = () => { }, runningStatus, }) {
  
  const [anchorEl, setAnchorEl] = useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = useState(false);
  const [transformProperty, setTransformProperty] = useState(100);

  const theme = useTheme();

  const styleSlider = {
    backgroundColor : "#EEEEEE",
    position : "absolute",
    top: "-20%",
    right : theme.spacing(14),
    zIndex : "-1",
    bottom : "-55%",
    transform : showFullContextMenu ? `translateY(${transformProperty}%)` : "translateY(0)"
  };

  const getOperatorStatus = (contextId) => {
    const state = runningStatus.operatorStatus;
    if (!state) {
      return false;
    }

    const context = state.find((st) => st.contextID === contextId);
    if (!context) {
      return false;
    }

    return context.operatorStatus.status === "ENABLED";
  };

  const getMeshSyncStatus = (contextId) => {
    const state = runningStatus.meshSyncStatus;
    if (!state) {
      return false;
    }

    const context = state.find((st) => st.contextID === contextId);
    if (!context) {
      return false;
    }

    return context.OperatorControllerStatus.status?.includes("ENABLED");
  };
  
  let open = Boolean(anchorEl);
  if (showFullContextMenu) {
    open = showFullContextMenu;
  }

  useEffect(() => {
    setTransformProperty((prev) => prev + (contexts.total_count ? contexts.total_count * 3.125 : 0));
  }, []);

  return (
    <>
      <IconButton
        aria-label="contexts"
        className="k8s-icon-button"
        onClick={(e) => {
          e.preventDefault();
          setShowFullContextMenu((prev) => !prev);
        }}
        onMouseOver={(e) => {
          e.preventDefault();
          setAnchorEl(true);
        }}
        onMouseLeave={(e) => {
          e.preventDefault();
          setAnchorEl(false);
        }}
        aria-owns={open ? "menu-list-grow" : undefined}
        aria-haspopup="true"
        style={{ marginRight: "0.5rem" }}
      >
        <CbadgeContainer>
          <img
            className="k8s-image"
            src="/static/img/kubernetes.svg"
            width = {theme.spacing(3)}
            style={{ zIndex: "2" }}
          />
          <Cbadge>{contexts?.total_count || 0}</Cbadge>
        </CbadgeContainer>
      </IconButton>

      <Slide direction="down" style={styleSlider} timeout={400} in={open} mountOnEnter unmountOnExit>
        <div>
          <ClickAwayListener
            onClickAway={(e) => {
                setAnchorEl(false);
                setShowFullContextMenu(false);
            }}
          >
            <CMenuContainer >
              <div>
                <TextField
                  id="search-ctx"
                  label="Search"
                  size="small"
                  variant="outlined"
                  onChange={(ev) => searchContexts(ev.target.value)}
                  style={{ width: "100%", backgroundColor: "rgba(102, 102, 102, 0.12)", margin: "1px 0px" }}
                  InputProps={{
                    endAdornment: <Search />,
                  }}
                />
              </div>
              <div style={{alignItems: "center"}}>
                {contexts?.total_count ? 
                  <>
                    <Checkbox
                      checked={activeContexts.includes("all")}
                      onChange={() => setActiveContexts("all")}
                      color="primary"
                    />
                    <span>Select All</span>
                  </>
                 : 
                  <Link href="/settings">
                    <Button
                      type="submit"
                      variant="contained"
                      color="primary"
                      size="large"
                      style={{ margin: "0.5rem 2rem" , whiteSpace: "nowrap" }}
                    >
                      <AddIcon/>
                      Connect Clusters
                    </Button>
                  </Link>
                }
                {contexts?.contexts?.map((ctx) => {
                  const meshStatus = getMeshSyncStatus(ctx.id);
                  const operStatus = getOperatorStatus(ctx.id);

                  function getStatus(status) {
                    if (status) {
                      return "Active";
                    } else {
                      return "InActive";
                    }
                  }

                  return (
                    <div id={ctx.id}>
                      <Tooltip
                        title={`Server: ${ctx.server}, Meshsync: ${getStatus(meshStatus)}, Operator: ${getStatus(
                          operStatus
                        )}`}
                      >
                        <div style={{ display: "flex", justifyContent: "flex-start", alignItems: "center" }}>
                          <Checkbox
                            checked={activeContexts.includes(ctx.id)}
                            onChange={() => setActiveContexts(ctx.id)}
                            color="primary"
                          />
                          <Chip
                            label={ctx?.name}
                            avatar={
                              meshStatus ? (
                                <BadgeAvatars>
                                  <Avatar
                                    src="/static/img/kubernetes.svg"
                                    style={operStatus ? {} : { opacity: 0.2 }}
                                  />
                                </BadgeAvatars>
                              ) : (
                                <Avatar
                                  src="/static/img/kubernetes.svg"
                                  style={operStatus ? { margin: 8 } : { opacity: 0.2, margin: 8 }}
                                />
                              )
                            }
                            variant="filled"
                            data-cy="chipContextName"
                          />
                        </div>
                      </Tooltip>
                    </div>
                  );
                })}
              </div>
            </CMenuContainer>
          </ClickAwayListener>
        </div>
      </Slide>
    </>
  );
}
