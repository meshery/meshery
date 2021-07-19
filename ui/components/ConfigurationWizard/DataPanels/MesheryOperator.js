/* eslint-disable react/display-name */
import CloseIcon from "@material-ui/icons/Close";
import {
  withStyles,
  Grid,
  Chip,
  IconButton,
  List,
  Paper,
} from "@material-ui/core/";
import { withSnackbar } from "notistack";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import { updateProgress } from "../../../lib/store";
import { pingMesheryOperatorWithNotification } from "../helpers/mesheryOperator";
import fetchMesheryOperatorStatus from "../../graphql/queries/OperatorStatusQuery";


const chipStyles = (theme) => ({
  chipIcon: {
    width: theme.spacing(2.5)
  },
  chip: {
    marginRight: theme.spacing(1),
    marginBottom: theme.spacing(1),
  },
})

const mesheryOperatorPingSnackbarAction = closeSnackbar => key => (
  <IconButton
    key="close"
    aria-label="Close"
    color="inherit"
    onClick={() => closeSnackbar(key)}
  >
    <CloseIcon />
  </IconButton>
)


const MesheryOperatorChip = withStyles(chipStyles)(({classes, handleMesheryOperatorClick, label}) => (
  <Chip
    // label={inClusterConfig?'Using In Cluster Config': contextName + (configuredServer?' - ' + configuredServer:'')}
    label={label}
    // onDelete={self.handleReconfigure}
    onClick={handleMesheryOperatorClick}
    icon={<img src="/static/img/meshery-operator.svg" className={classes.chipIcon} />}
    variant="outlined"
    data-cy="chipOperator"
  />
))



const MesheryOperatorDataPanel = ({operatorInformation, updateProgress, enqueueSnackbar, closeSnackbar}) => (

  <Paper style={{padding: "2rem"}}>
    <MesheryOperatorChip
      handleMesheryOperatorClick={() => pingMesheryOperatorWithNotification(updateProgress, enqueueSnackbar, mesheryOperatorPingSnackbarAction(closeSnackbar), fetchMesheryOperatorStatus )} 
      label="Meshery Operator"
    />
    <Grid container spacing={1}>
      <Grid item xs={12} md={4}>
        <List>
          <ListItem>
            <ListItemText primary="Operator State" secondary={operatorInformation.operatorInstalled ? "Active" : "Disabled"} />
          </ListItem>
          <ListItem>
            <ListItemText primary="Operator Version" secondary={operatorInformation.operatorVersion} />
          </ListItem>
        </List>
      </Grid>
      <Grid item xs={12} md={4}>
        <List>
          <ListItem>
            <ListItemText primary="MeshSync State" secondary={operatorInformation.meshSyncInstalled ? "Active" : "Disabled"} />
          </ListItem>
          <ListItem>
            <ListItemText primary="MeshSync Version" secondary={operatorInformation.meshSyncVersion} />
          </ListItem>
        </List>
      </Grid>
      <Grid item xs={12} md={4}>
        <List>
          <ListItem>
            <ListItemText primary="NATS State" secondary={operatorInformation.NATSInstalled ? "Active" : "Disabled"} />
          </ListItem>
          <ListItem>
            <ListItemText primary="NATS Version" secondary={operatorInformation.NATSVersion} />
          </ListItem>
        </List>
      </Grid>
    </Grid>
  </Paper>
)


const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

export default connect(null, mapDispatchToProps)(withSnackbar(MesheryOperatorDataPanel))
