/* eslint-disable react/display-name */
/* eslint-disable  no-unused-vars*/
import {
  Grid,
  List,
  Paper,
} from "@material-ui/core/";
import { connect } from "react-redux";
import { bindActionCreators } from "redux";
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import { updateProgress } from "../../../lib/store";
import { pingMesheryOperator } from "../helpers/mesheryOperator";
import fetchMesheryOperatorStatus from "../../graphql/queries/OperatorStatusQuery";
import AdapterChip from "./AdapterChip"
import { useNotification } from "../../../utils/hooks/useNotification";
import { EVENT_TYPES } from "../../../lib/event-types";


const chipStyles = (theme) => ({ chipIcon : { width : theme.spacing(2.5) },
  chip : { marginRight : theme.spacing(1),
    marginBottom : theme.spacing(1), }, })

// Connection Wizard
// TODO: bind to contextID prop, leaving due to no use in current UI
const MesheryOperatorDataPanel = ({
  operatorInformation
}) => {
  const { notify } = useNotification()
  const handleMesheryOperatorClick = () => {

    const successCb = (res) => {
      if (res?.operator?.status == "ENABLED") {
        notify({ message : 'Operator was pinged!', type : EVENT_TYPES.SUCCESS })
      } else {
        notify({ message : 'Operator was not pinged!', type : EVENT_TYPES.ERROR })
      }
    }

    const errorCb = (err) => {
      notify({ message : 'Unable to ping meshery operator!', type : EVENT_TYPES.ERROR, details : err.toString() })
    }


    pingMesheryOperator(
      fetchMesheryOperatorStatus,
      successCb,
      errorCb
    )
  }


  return (
    <Paper style={{ padding : "2rem" }}>
      <AdapterChip
        handleClick={handleMesheryOperatorClick}
        isActive={true}
        image="/static/img/meshery-operator.svg"
        label="Meshery Operator"
      />
      <Grid container spacing={1}>
        <Grid item xs={12} md={4}>
          <List>
            <ListItem>
              <ListItemText primary="Operator State" secondary={operatorInformation.operatorInstalled
                ? "Active"
                : "Disabled"} />
            </ListItem>
            <ListItem>
              <ListItemText primary="Operator Version" secondary={operatorInformation.operatorVersion} />
            </ListItem>
          </List>
        </Grid>
        <Grid item xs={12} md={4}>
          <List>
            <ListItem>
              <ListItemText primary="MeshSync State" secondary={operatorInformation.meshSyncInstalled
                ? "Active"
                : "Disabled"} />
            </ListItem>
            <ListItem>
              <ListItemText primary="MeshSync Version" secondary={operatorInformation.meshSyncVersion} />
            </ListItem>
          </List>
        </Grid>
        <Grid item xs={12} md={4}>
          <List>
            <ListItem>
              <ListItemText primary="NATS State" secondary={operatorInformation.NATSInstalled
                ? "Active"
                : "Disabled"} />
            </ListItem>
            <ListItem>
              <ListItemText primary="NATS Version" secondary={operatorInformation.NATSVersion} />
            </ListItem>
          </List>
        </Grid>
      </Grid>
    </Paper>
  )
}


const mapDispatchToProps = (dispatch) => ({ updateProgress : bindActionCreators(updateProgress, dispatch), });

export default connect(null, mapDispatchToProps)(MesheryOperatorDataPanel)
