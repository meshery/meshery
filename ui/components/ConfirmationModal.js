import {
  Button, Checkbox, Chip, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, TextField,
  Tooltip, Typography
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import { withStyles } from "@material-ui/core/styles";
import { Search } from "@material-ui/icons";
import { withSnackbar } from "notistack";
import { connect } from "react-redux";
import { setK8sContexts, updateProgress } from "../lib/store";
import { closeButtonForSnackbarAction, errorHandlerGenerator, hideProgress, showProgress, successHandlerGenerator } from "./ConnectionWizard/helpers/common";
import { pingKubernetes } from "./ConnectionWizard/helpers/kubernetesHelpers";
import { getK8sConfigIdsFromK8sConfig } from "../utils/multi-ctx";
import { bindActionCreators } from "redux";
import { useState } from "react";

const styles = (theme) => ({
  icon : {
    display : 'inline',
    verticalAlign : 'text-top',
    width : theme.spacing(1.75),
    marginLeft : theme.spacing(0.5),
  },
  chip : {
    height : "50px",
    fontSize : "15px",
    position : "relative",
    top : theme.spacing(0.5),
    [theme.breakpoints.down("md")] : { fontSize : "12px", },
  },

  ctxChip : {
    backgroundColor : "white",
    cursor : "pointer",
    marginRight : theme.spacing(1),
    marginLeft : theme.spacing(1),
    marginBottom : theme.spacing(1),
    height : "100%",
    padding : theme.spacing(0.5)
  },
  ctxIcon : {
    display : 'inline',
    verticalAlign : 'text-top',
    width : theme.spacing(2.5),
    marginLeft : theme.spacing(0.5),
  },
  title : {
    textAlign : 'center',
    minWidth : 400,
    padding : theme.spacing(1.5),
    color : '#fff',
    backgroundColor : '#607d8b',
    fontSize : "1.6rem",

  },
  subtitle : {
    minWidth : 400,
    overflowWrap : 'anywhere',
    textAlign : 'center',
    padding : '5px'
  },
  button0 : {
    margin : theme.spacing(0.5),
    padding : theme.spacing(1),
    borderRadius : 5,
    minWidth : 100,
  },
  button1 : {
    margin : theme.spacing(0.5),
    padding : theme.spacing(1),
    borderRadius : 5,
    backgroundColor : "#e0e0e0",
    color : "rgba(0, 0, 0, 0.87)",
    "&:hover" : {
      backgroundColor : "#d5d5d5",
      boxShadow : "0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)"
    },
    minWidth : 100,
  },
  actions : {
    display : 'flex',
    justifyContent : 'space-evenly',
  },
  all : {
    display : "table"
  },
  contexts : {
    display : "flex",
    flexWrap : "wrap"
  }
})

function ConfirmationMsg(props) {
  const { classes, open, handleClose, submit, isDelete,
    selectedK8sContexts, k8scontext, title, setK8sContexts, enqueueSnackbar, closeSnackbar } = props

  const [contexts, setContexts] = useState(k8scontext);

  const handleKubernetesClick = (ctxID) => {
    showProgress()
    pingKubernetes(
      successHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes succesfully pinged", () => hideProgress()),
      errorHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes not pinged successfully", () => hideProgress()),
      ctxID
    )
  }

  const handleSubmit = () => {
    if (selectedK8sContexts.length === 0) {
      enqueueSnackbar("Please select Kubernetes context(s) before proceeding with the operation",
        { variant : "info", preventDuplicate : true,
          action : (key) => (
            <IconButton key="close" aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
              <CloseIcon />
            </IconButton>
          ),
          autoHideDuration : 3000, });
    }
    submit();
    handleClose();
  }
  const searchContexts = (search) => {
    if (search === "") {
      setContexts(k8scontext);
    }
    let matchedCtx = [];
    k8scontext.forEach(ctx => {
      if (ctx.contextName.includes(search)) {
        matchedCtx.push(ctx);
      }
    });
    setContexts(matchedCtx);
  }

  const setContextViewer = (id) => {
    if (id === "all") {
      if (selectedK8sContexts.includes("all")) {
        updateProgress({ showProgress : true })
        setK8sContexts({ selectedK8sContexts : [] })
      } else {
        setK8sContexts({ selectedK8sContexts : ["all"] });
      }
      return;
    }

    if (selectedK8sContexts.includes(id)) {
      const filteredContexts = selectedK8sContexts.filter(cid => cid !== id );
      setK8sContexts({ selectedK8sContexts : filteredContexts })
    } else if (selectedK8sContexts[0] === "all") {
      const allContextIds = getK8sConfigIdsFromK8sConfig(k8scontext);
      setK8sContexts({ selectedK8sContexts : allContextIds.filter(cid => cid !== id) });
    } else {
      if (selectedK8sContexts.length === k8scontext.length - 1) {
        setK8sContexts({ selectedK8sContexts : ["all"] })
        return;
      }
      setK8sContexts({ selectedK8sContexts : [...selectedK8sContexts, id] });
    }
  }
  return (
    <div className={classes.root}>
      <Dialog
        open={open}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        className={classes.dialogBox}
      >
        {k8scontext.length > 0 ?
          <>
            <DialogTitle id="alert-dialog-title" className={classes.title}>
              {title}
            </DialogTitle>
            <DialogContent>
              <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
                <Typography variant="body1">
                  <TextField
                    id="search-ctx"
                    label="Search"
                    size="small"
                    variant="outlined"
                    onChange={(event) => searchContexts(event.target.value)}
                    style={{ width : "100%", backgroundColor : "rgba(102, 102, 102, 0.12)", margin : "1px 1px 8px " }}
                    InputProps={{
                      endAdornment : (
                        <Search />
                      )
                    }}
                  />
                  <div className={classes.all}>
                    <Checkbox
                      checked={selectedK8sContexts?.includes("all")}
                      onChange={() => setContextViewer("all")}
                      color="primary"
                    />
                    <span>Select All</span>
                  </div>
                  <div className={classes.contexts}>
                    {
                      contexts.map((ctx) => (
                        <div id={ctx.contextID} className={classes.chip}>
                          <Tooltip title={`Server: ${ctx.configuredServer}`}>
                            <div style={{ display : "flex", justifyContent : "flex-wrap", alignItems : "center" }}>
                              <Checkbox
                                checked={selectedK8sContexts.includes(ctx.contextID) || (selectedK8sContexts.length > 0 && selectedK8sContexts[0] === "all")}
                                onChange={() => setContextViewer(ctx.contextID)}
                                color="primary"
                              />
                              <Chip
                                label={ctx.contextName}
                                className={classes.ctxChip}
                                onClick={() => handleKubernetesClick(ctx.contextID)}
                                icon={<img src="/static/img/kubernetes.svg" className={classes.ctxIcon} />}
                                variant="outlined"
                                data-cy="chipContextName"
                              />
                            </div>

                          </Tooltip>
                        </div>
                      ))
                    }
                  </div>
                </Typography>
              </DialogContentText>
            </DialogContent>
            <DialogActions className={classes.actions}>
              <Button onClick={handleClose} className={classes.button1}>
                <Typography variant body2> Cancel </Typography>
              </Button>
              <Button  disabled
                className={classes.button0} autoFocus type="submit"
                variant="contained"
                color="primary">
                <Typography variant body2 > {isDelete ? "UNDEPLOY LATER" : "DEPLOY"} </Typography>
              </Button>
              <Button onClick={handleSubmit}
                className={classes.button0} autoFocus type="submit"
                variant="contained"
                color="primary">
                <Typography variant body2 > {isDelete ? "UNDEPLOY" : "DEPLOY"} </Typography>
              </Button>
            </DialogActions>
          </>
          :
          <>
            <DialogTitle id="alert-dialog-title" className={classes.title}>
            No Kubernetes contexts detected
            </DialogTitle>
            <DialogContent>
              <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
                <Typography variant="subtitle1">Please upload kube config file.</Typography>
              </DialogContentText>
            </DialogContent>
          </>
        }
      </Dialog>
    </div>
  )
}

const mapStateToProps = state => {
  const selectedK8sContexts = state.get('selectedK8sContexts');
  const k8scontext = state.get("k8sConfig");
  return { selectedK8sContexts : selectedK8sContexts, k8scontext : k8scontext };
}

const mapDispatchToProps = (dispatch) => ({
  updateProgress : bindActionCreators(updateProgress, dispatch),
  setK8sContexts : bindActionCreators(setK8sContexts, dispatch)
});

export default withStyles(styles)(connect(mapStateToProps, mapDispatchToProps)(withSnackbar(ConfirmationMsg)));