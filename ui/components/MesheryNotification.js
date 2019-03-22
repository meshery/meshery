import IconButton from '@material-ui/core/IconButton';
import {connect} from "react-redux";
import { bindActionCreators } from 'redux';
import NoSsr from '@material-ui/core/NoSsr';
import { Badge, Drawer, Tooltip } from '@material-ui/core';
import MesheryEventViewer from './MesheryEventViewer';
import NotificationsIcon from '@material-ui/icons/Notifications';
import { withStyles } from '@material-ui/core/styles';

const styles = {
    sidelist: {
        width: 350,
    },
};

class MesheryNotification extends React.Component {

  state = {
    events: [],
    badgeContent: 0,
    open: false,
    k8sConfig: {
        inClusterConfig: false,
        k8sfile: '', 
        contextName: '', 
        meshLocationURL: '', 
        reconfigureCluster: true,
    },
    createStream: false,
  }

  handleToggle = () => {
    this.setState(state => ({ open: !state.open }));
  };

  handleClose = event => {
    if (this.anchorEl.contains(event.target)) {
      return;
    }
    this.setState({ open: false });
  };

  static getDerivedStateFromProps(props, state){
    if (JSON.stringify(props.k8sConfig) !== JSON.stringify(state.k8sConfig)) {
        return {
            createStream: true,
            k8sConfig: props.k8sConfig,
        };
    }
    return null;
  }

  componentDidMount() {
    // if (this.state.createStream) {
    //     this.startEventStream();
    // }
  }

  componentDidUpdate(){
      const {createStream, k8sConfig} = this.state;
    if (k8sConfig.k8sfile === '' && k8sConfig.meshLocationURL === '') {
        this.closeEventStream();
    }
    if (createStream && k8sConfig.k8sfile !== '' && k8sConfig.meshLocationURL !== '') {
        this.startEventStream();
    }
  }

  async startEventStream() {
    this.closeEventStream();
    this.eventStream = new EventSource("/api/events");
    this.eventStream.onmessage = this.handleEvents;
    this.setState({createStream: false});
  }

  handleEvents = e => {
    const {events} = this.state;
    const data = JSON.parse(e.data);
    events.push(data);
    this.setState({events, badgeContent: events.length});
  }

  closeEventStream() {
      if(this.eventStream && this.eventStream.close){
        this.eventStream.close();
      }
  }

  deleteEvent = (ind) => {
      const {events} = this.state;
      if (events[ind]){
        events.splice(ind, 1);
      }
      this.setState({events, badgeContent: events.length});
  }

  render() {
    const {classes} = this.props;
    const { open, badgeContent, events } = this.state;
    const self = this;

    return (
      <div>
        <NoSsr>
        <Tooltip title={`There are ${badgeContent} events`}>
        <IconButton 
            buttonRef={node => {
                this.anchorEl = node;
              }}
            color="inherit" onClick={this.handleToggle}>
            <Badge badgeContent={badgeContent} color="secondary">
                <NotificationsIcon />
            </Badge>
        </IconButton>
        </Tooltip>

        <Drawer anchor="right" open={open} onClose={this.handleClose}>
          <div
            tabIndex={0}
            role="button"
            onClick={this.handleClose}
            onKeyDown={this.handleClose}
          >
            <div className={classes.sidelist}>
            {events && events.map((event, ind) => (
                <MesheryEventViewer eventVariant={event.event_type} eventSummary={event.summary} 
                    eventDetail={event.details} eventIndex={ind} deleteEvent={self.deleteEvent} />
                // <span>
                //     eventVariant={event.event_type} eventSummary={event.summary} 
                //     eventDetail={event.details} eventIndex={ind} deleteEvent={self.deleteEvent}
                // </span>
            ))}
            </div>
          </div>
        </Drawer>

        </NoSsr>
    </div>
    )
  }
}

const mapDispatchToProps = dispatch => {
  return {
    // updateUser: bindActionCreators(updateUser, dispatch)
  }
}
const mapStateToProps = state => {
    const k8sConfig = state.get("k8sConfig").toObject();
    return {k8sConfig};
  }

export default withStyles(styles)(connect(
    mapStateToProps,
  mapDispatchToProps
)(MesheryNotification));