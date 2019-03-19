import { Component } from "react";
import MesheryDateTimePicker from "./MesheryDateTimePicker";
import { NoSsr, Grow, Paper, Button, Popper } from "@material-ui/core";
import Moment from "react-moment";
import AccessTimeIcon from '@material-ui/icons/AccessTime';
import { withStyles } from '@material-ui/core/styles';

const styles = theme => ({

});

 class GrafanaDateRangePicker extends Component {

    state = {
        startDate: new Date(),
        endDate: new Date(),
        open: false,
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

    handleChange = name => event => {
        if (name === 'startDate' || name === 'endDate'){
            const {startDate, endDate} = this.state;
            const dt = event.toDate();
            if (name === 'startDate') {
                if (dt > endDate) {
                    this.setState({ [name]: dt, endDate: dt, open: false });        
                    return;
                }
            } else if(name === 'endDate') {
                if ( dt < startDate ){
                    this.setState({ [name]: dt, startDate: dt, open: false });        
                    return;
                }
            }
            this.setState({ [name]: dt, open: false });
            return;
        }
        this.setState({ [name]: event.target.value });
    };

    render() {
        const {open, startDate, endDate} = this.state;
        return (
            <NoSsr>
            <React.Fragment>
            <Button
                variant="outlined"
                buttonRef={node => {
                    this.anchorEl = node;
                }}
                aria-owns={open ? 'dateRange-list-grow' : undefined}
                aria-haspopup="true"
                onClick={this.handleToggle}>
                <AccessTimeIcon /> <Moment format="LLLL">{startDate}</Moment>{' - '}<Moment format="LLLL">{endDate}</Moment>
            </Button>
            <Popper open={open} anchorEl={this.anchorEl} transition placement='bottom-end'>
            {({ TransitionProps, placement }) => (
              <Grow
                {...TransitionProps}
                id="dateRange-list-grow"
                style={{ transformOrigin: placement === 'bottom' ? 'left top' : 'left bottom' }}
              >
                <Paper>
                    <MesheryDateTimePicker selectedDate={startDate} onChange={this.handleChange('startDate')} label={"Start"} />
                    <MesheryDateTimePicker selectedDate={endDate} onChange={this.handleChange('endDate')} label={"End"} />    
                </Paper>
              </Grow>
            )}
          </Popper>
            </React.Fragment>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(GrafanaDateRangePicker);