import { Component } from "react";
import MesheryDateTimePicker from "./MesheryDateTimePicker";
import { NoSsr, Grow, Paper, Button, Popper, TextField, MenuItem, Grid } from "@material-ui/core";
import Moment from "react-moment";
import AccessTimeIcon from '@material-ui/icons/AccessTime';
import { withStyles } from '@material-ui/core/styles';

const styles = theme => ({

});

const refreshIntervals = [
    'off',
    '5s',
    '10s',
    '30s',
    '1m',
    '5m',
    '15m',
    '30m',
    '1h',
    '2h',
    '1d',
]

 class GrafanaDateRangePicker extends Component {

    state = {
        startDate: new Date(),
        endDate: new Date(),
        refreshInterval: '',
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
                    this.setState({ [name]: dt, endDate: dt });        
                    return;
                }
            } else if(name === 'endDate') {
                if ( dt < startDate ){
                    this.setState({ [name]: dt, startDate: dt });        
                    return;
                }
            }
            this.setState({ [name]: dt });
            return;
        }
        this.setState({ [name]: event.target.value });
    };

    render() {
        const {open, startDate, endDate, refreshInterval} = this.state;
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
            <Popper open={open} anchorEl={this.anchorEl} transition placement='bottom-start'>
            {({ TransitionProps, placement }) => (
              <Grow
                {...TransitionProps}
                id="dateRange-list-grow"
                style={{ transformOrigin: placement === 'bottom' ? 'left top' : 'left bottom' }}
              >
                <Paper>
                <Grid container spacing={5}>
                    <Grid item xs={12} sm={4}>
                    <MesheryDateTimePicker selectedDate={startDate} onChange={this.handleChange('startDate')} label={"Start"} />
                    <MesheryDateTimePicker selectedDate={endDate} onChange={this.handleChange('endDate')} label={"End"} />
                    <TextField
                        select
                        id="refreshInterval"
                        name="refreshInterval"
                        label="Refresh Interval"
                        fullWidth
                        value={refreshInterval}
                        margin="normal"
                        variant="outlined"
                        onChange={this.handleChange('refreshInterval')}
                    >
                        {refreshIntervals.map((ri) => (
                            <MenuItem key={'ri_-_-_'+ri} value={ri}>{ri}</MenuItem>
                        ))}
                    </TextField>
                    </Grid>
                    <Grid item xs={12} sm={8}>
                            
                    </Grid>
                </Grid>
                    
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