import { Button, Chip, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, TextField, Tooltip, Typography } from "@material-ui/core";
import { withStyles } from "@material-ui/core/styles";
import { connect } from "react-redux";
import { closeButtonForSnackbarAction, errorHandlerGenerator, hideProgress, showProgress, successHandlerGenerator } from "./ConnectionWizard/helpers/common";
import { pingKubernetes } from "./ConnectionWizard/helpers/kubernetesHelpers";
import { useSnackbar } from "notistack";
import { Search } from "@material-ui/icons";
import { useEffect, useState } from "react";


const styles = (theme) => ({
  icon : {
    display : 'inline',
    verticalAlign : 'text-top',
    width : theme.spacing(1.75),
    marginLeft : theme.spacing(0.5),
  },
  Chip : {
    backgroundColor : "white",
    cursor : "pointer",
    marginRight : theme.spacing(1),
    marginLeft : theme.spacing(1),
    marginBottom : theme.spacing(1),
    height : "100%",
  },
})

function ConfirmationMsg(props) {
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const { classes, open, handleClose, submit, category, operation, isDelete, activeK8sContext } = props
  const [context, setContext] = useState(activeK8sContext)

  useEffect(() => {
    setContext(activeK8sContext)
  },[activeK8sContext])

  const handleKubernetesClick = () => {
    showProgress()
    pingKubernetes(
      successHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes succesfully pinged", () => hideProgress()),
      errorHandlerGenerator(enqueueSnackbar, closeButtonForSnackbarAction(closeSnackbar), "Kubernetes not pinged successfully", () => hideProgress())
    )
  }

  const searchContexts = (search) => {
    const matchedCtx = activeK8sContext.filter((ctx) => ctx.name.includes(search))
    setContext(matchedCtx)
  }
  console.log(context, "state")
  return (
    <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{"The selected operation will be applied to following contexts."}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          <div>
            <TextField
              id="search-ctx"
              label="Search"
              size="small"
              variant="outlined"
              onChange={ev => searchContexts(ev.target.value)}
              style={{ width : "100%", backgroundColor : "rgba(102, 102, 102, 0.12)", margin : "1px 1px 8px " }}
              InputProps={{ endAdornment : (
                <Search className={classes.searchIcon} />
              ) }}
            />
          </div>
          <Typography>
            { context.map((ctx) => (
              <Tooltip title={`Server: ${ctx.server}`}>
                <Chip
                  label={ctx.name}
                  className={classes.chip}
                  onClick={handleKubernetesClick}
                  icon={<img src = "/static/img/kubernetes.svg" className={classes.icon}  />}
                  variant="outlined"
                  data-cy="chipContextName"
                />
              </Tooltip>
            ))
            }
          </Typography>
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary">
              Cancel
        </Button>
        <Button onClick={() => submit(category, operation, isDelete)} color="primary" autoFocus>
              Continue
        </Button>
      </DialogActions>
    </Dialog>


  )
}

const mapStateToProps = state => ({
  activeK8sContext : state.get("activeK8sContext").toJS()
})

export default withStyles(styles)(connect(mapStateToProps)(ConfirmationMsg));


