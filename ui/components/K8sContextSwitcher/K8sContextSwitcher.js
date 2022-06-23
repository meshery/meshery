import React, { useEffect } from "react";
import { useSnackbar } from "notistack";
import { AddIcon, Button, ClickAwayListener, IconButton, Link, Paper, Search, Slide, TextField } from "@mui/material";
import { useStyles } from "./K8sContextSwitcher.styles";

export default function K8sContextMenu({ contexts }) {
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = React.useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = React.useState(false);
  const [transformProperty, setTransformProperty] = React.useState(100);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();

  const styleSlider = {
    backgroundColor: "#EEEEEE",
    position: "absolute",
    left: "-5rem",
    zIndex: "-1",
    bottom: "-55%",
    transform: showFullContextMenu ? `translateY(${transformProperty}%)` : "translateY(0)",
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

  const handleKubernetesClick = (id) => {
    showProgress();

    pingKubernetes(
      successHandlerGenerator(
        enqueueSnackbar,
        closeButtonForSnackbarAction(closeSnackbar),
        "Kubernetes succesfully pinged",
        successCallback
      ),
      errorHandlerGenerator(
        enqueueSnackbar,
        closeButtonForSnackbarAction(closeSnackbar),
        "Kubernetes not pinged successfully",
        () => hideProgress()
      ),
      id
    );
  };

  const handleKubernetesDelete = (name, ctxId) => () => {
    if (confirm(`Are you sure you want to delete "${name}" cluster from Meshery?`)) {
      const successCallback = async () => {
        showProgress();
        const updatedConfig = await loadActiveK8sContexts();
        if (Array.isArray(updatedConfig)) {
          updateK8SConfig({ k8sConfig: updatedConfig });
        }
      };
      deleteKubernetesConfig(
        successHandlerGenerator(
          enqueueSnackbar,
          closeButtonForSnackbarAction(closeSnackbar),
          "Kubernetes config successfully removed",
          successCallback
        ),
        errorHandlerGenerator(
          enqueueSnackbar,
          closeButtonForSnackbarAction(closeSnackbar),
          "Not able to remove config"
        ),
        ctxId
      );
    }
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
        <div className={classes.cbadgeContainer}>
          <img
            className="k8s-image"
            src="/static/img/kubernetes.svg"
            width="24px"
            height="24px"
            style={{ zIndex: "2" }}
          />
          <div className={classes.cbadge}>{contexts?.total_count || 0}</div>
        </div>
      </IconButton>

      <Slide direction="down" style={classes.styleSlider} timeout={400} in={open} mountOnEnter unmountOnExit>
        <div>
          <ClickAwayListener
            onClickAway={(e) => {
              if (
                !e.target.className?.includes("cbadge") &&
                e.target?.className != "k8s-image" &&
                !e.target.className.includes("k8s-icon-button")
              ) {
                setAnchorEl(false);
                setShowFullContextMenu(false);
              }
            }}
          >
            <Paper className={classes.cMenuContainer}>
              <div>
                <TextField
                  id="search-ctx"
                  label="Search"
                  size="small"
                  variant="outlined"
                  onChange={(ev) => searchContexts(ev.target.value)}
                  style={{ width: "100%", backgroundColor: "rgba(102, 102, 102, 0.12)", margin: "1px 0px" }}
                  InputProps={{
                    endAdornment: <Search className={classes.searchIcon} />,
                  }}
                />
              </div>
              <div>
                {contexts?.total_count ? (
                  <>
                    <Checkbox
                      checked={activeContexts.includes("all")}
                      onChange={() => setActiveContexts("all")}
                      color="primary"
                    />
                    <span>Select All</span>
                  </>
                ) : (
                  <Link href="/settings">
                    <Button
                      type="submit"
                      variant="contained"
                      color="primary"
                      size="large"
                      style={{ margin: "0.5rem 0.5rem", whiteSpace: "nowrap" }}
                    >
                      <AddIcon className={classes.AddIcon} />
                      Connect Clusters
                    </Button>
                  </Link>
                )}
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
                    <div id={ctx.id} className={classes.chip}>
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
                            onDelete={handleKubernetesDelete(ctx.name, ctx.id)}
                            onClick={() => handleKubernetesClick(ctx.id)}
                            avatar={
                              meshStatus ? (
                                <BadgeAvatars>
                                  <Avatar
                                    src="/static/img/kubernetes.svg"
                                    className={classes.icon}
                                    style={operStatus ? {} : { opacity: 0.2 }}
                                  />
                                </BadgeAvatars>
                              ) : (
                                <Avatar
                                  src="/static/img/kubernetes.svg"
                                  className={classes.icon}
                                  style={operStatus ? { margin: 8 } : { opacity: 0.2, margin: 8 }}
                                />
                              )
                            }
                            variant="filled"
                            className={classes.Chip}
                            data-cy="chipContextName"
                          />
                        </div>
                      </Tooltip>
                    </div>
                  );
                })}
              </div>
            </Paper>
          </ClickAwayListener>
        </div>
      </Slide>
    </>
  );
}
