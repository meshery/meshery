import { Component } from "react";
import MesheryDateTimePicker from "./MesheryDateTimePicker";
import { NoSsr, Grow, Paper, Button, Popper, TextField, MenuItem, Grid, Popover, Dialog, DialogContent, DialogContentText, DialogActions, FormControl, FormLabel, FormGroup, DialogTitle } from "@material-ui/core";
import Moment from "react-moment";
import AccessTimeIcon from '@material-ui/icons/AccessTime';
import { withStyles } from '@material-ui/core/styles';

const styles = theme => ({
    rangeDialog: {
        // width: window.innerWidth * 0.7,
    },
    dateTimePicker: {
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(1),
    },
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

const quickRanges = [
    [
        'Last 2 days',
        'Last 7 days',
        'Last 30 days',
        'Last 90 days',
        'Last 6 months',
        'Last 1 year',
        'Last 2 years',
        'Last 5 years',
    ], [
        'Yesterday',
        'Day before yesterday',
        'This day last week',
        'Previous week',
        'Previous month',
        'Previous year',
    ], [
        'Today',
        'Today so far',
        'This week',
        'This week so far',
        'This month',
        'This month so far',
        'This year',
        'This year so far',
    ], [
        'Last 5 minutes',
        'Last 15 minutes',
        'Last 30 minutes',
        'Last 1 hour',
        'Last 3 hours',
        'Last 6 hours',
        'Last 12 hours',
        'Last 24 hours',
    ]
]

 class GrafanaDateRangePicker extends Component {
    constructor(props) {
        super(props);

        const startDate = new Date();
        startDate.setMinutes(startDate.getMinutes() - 5);
        this.state = {
            startDate,
            endDate: new Date(),
            refreshInterval: '10s',
            open: false,
        }
    }

    handleClick = event => {
        this.setState({
            open: true,
        });
    };

    setRange = range => () => {
        let {startDate, endDate} = this.state;
        
        switch (range) {
            case 'Last 2 days':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 2);
                endDate = new Date();
                break;
            case 'Last 7 days':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 7);
                endDate = new Date();
                break;
            case 'Last 30 days':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 30);
                endDate = new Date();
                break;
            case 'Last 90 days':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 90);
                endDate = new Date();
                break;
            case 'Last 6 months':
                startDate = new Date();
                startDate.setMonth(startDate.getMonth() - 6);
                endDate = new Date();
                break;
            case 'Last 1 year':
                startDate = new Date();
                startDate.setFullYear(startDate.getFullYear() - 1);
                endDate = new Date();
                break;
            case 'Last 2 years':
                startDate = new Date();
                startDate.setFullYear(startDate.getFullYear() - 2);
                endDate = new Date();
                break;
            case 'Last 5 years':
                startDate = new Date();
                startDate.setFullYear(startDate.getFullYear() - 5);
                endDate = new Date();
                break;

            case 'Yesterday':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setDate(endDate.getDate() - 1);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'Day before yesterday':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 2);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setDate(endDate.getDate() - 2);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'This day last week':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 7);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setDate(endDate.getDate() - 7);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'Previous week':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - 6 - (startDate.getDay() + 8) % 7);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setDate(startDate.getDate() + 6);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'Previous month':
                startDate = new Date();
                startDate.setMonth(startDate.getMonth() - 1);
                startDate.setDate(1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setMonth(endDate.getMonth());
                endDate.setDate(0);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'Previous year':
                startDate = new Date();
                startDate.setFullYear(startDate.getFullYear() - 1);
                startDate.setMonth(0);
                startDate.setDate(1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setFullYear(endDate.getFullYear() - 1);
                endDate.setMonth(12);
                endDate.setDate(0);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;

            case 'Today':
                startDate = new Date();
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'Today so far':
                startDate = new Date();
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                break;
            case 'This week':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - (startDate.getDay() + 7) % 7);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setDate(startDate.getDate() + 6);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'This week so far':
                startDate = new Date();
                startDate.setDate(startDate.getDate() - (startDate.getDay() + 7) % 7);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                break;
            case 'This month':
                startDate = new Date();
                startDate.setDate(1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setMonth(endDate.getMonth()+1);
                endDate.setDate(0);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'This month so far':
                startDate = new Date();
                startDate.setDate(1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                break;
            case 'This year':
                startDate = new Date();
                startDate.setMonth(0);
                startDate.setDate(1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                endDate.setMonth(12);
                endDate.setDate(0);
                endDate.setHours(23);
                endDate.setMinutes(59);
                endDate.setSeconds(59);
                endDate.setMilliseconds(999);
                break;
            case 'This year so far':
                startDate = new Date();
                startDate.setMonth(0);
                startDate.setDate(1);
                startDate.setHours(0);
                startDate.setMinutes(0);
                startDate.setSeconds(0);
                startDate.setMilliseconds(0);
                endDate = new Date();
                break;

            case 'Last 5 minutes':
                startDate = new Date();
                startDate.setMinutes(startDate.getMinutes() - 5);
                endDate = new Date();
                break;
            case 'Last 15 minutes':
                startDate = new Date();
                startDate.setMinutes(startDate.getMinutes() - 15);
                endDate = new Date();
                break;
            case 'Last 30 minutes':
                startDate = new Date();
                startDate.setMinutes(startDate.getMinutes() - 30);
                endDate = new Date();
                break;
            case 'Last 1 hour':
                startDate = new Date();
                startDate.setHours(startDate.getHours() - 1);
                endDate = new Date();
                break;
            case 'Last 3 hours':
                startDate = new Date();
                startDate.setHours(startDate.getHours() - 3);
                endDate = new Date();
                break;
            case 'Last 6 hours':
                startDate = new Date();
                startDate.setHours(startDate.getHours() - 6);
                endDate = new Date();
                break;
            case 'Last 12 hours':
                startDate = new Date();
                startDate.setHours(startDate.getHours() - 12);
                endDate = new Date();
                break;
            case 'Last 24 hours':
                startDate = new Date();
                startDate.setHours(startDate.getHours() - 24);
                endDate = new Date();
                break;
        }
        this.setState({startDate, endDate})
    }
        
    // handleToggle = () => {
    //     this.setState(state => ({ open: !state.open }));
    // };

    handleClose = event => {
        // if (this.anchorEl.contains(event.target)) {
        //     return;
        // }
        this.setState({ open: false });
        // this.setState({
        //     anchorEl: null,
        //   });
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
        const {startDate, endDate, refreshInterval, open} = this.state;
        // const open = Boolean(anchorEl);
        const {classes} = this.props;
        return (
            <NoSsr>
            <React.Fragment>
            <Button
                variant="outlined"
                // buttonRef={node => {
                //     this.anchorEl = node;
                // }}
                // aria-owns={open ? 'daterange-popper' : undefined}
                // aria-haspopup="true"
                onClick={this.handleClick}>
                <AccessTimeIcon /> <Moment format="LLLL">{startDate}</Moment>{' - '}<Moment format="LLLL">{endDate}</Moment>
            </Button>
            {/* <Popper open={open} anchorEl={this.anchorEl} transition placement='bottom-start'>
            {({ TransitionProps, placement }) => (
              <Grow
                {...TransitionProps}
                id="dateRange-list-grow"
                style={{ transformOrigin: placement === 'bottom' ? 'left top' : 'left bottom' }}
              > */}
              {/* <Popover 
                 id="daterange-popper"
                 open={open}
                 anchorEl={anchorEl}
                 onClose={this.handleClose}
                anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                }}
                transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                }} */}
                {/* > */}
                {/* <Paper> */}
                <Dialog
                    open={open}
                    onClose={this.handleClose}
                    scroll={'paper'}
                    aria-labelledby="daterange-dialog-title"
                    fullWidth
                    maxWidth='md'>
                    <DialogTitle id="daterange-dialog-title">Select a Date Range</DialogTitle>
                    <DialogContent>
                        <DialogContentText className={classes.rangeDialog}>
                            <Grid container spacing={5}>
                                <Grid item xs={12} sm={4}>
                                        <MesheryDateTimePicker selectedDate={startDate} onChange={this.handleChange('startDate')} label={"Start"} className={classes.dateTimePicker} />
                                        <MesheryDateTimePicker selectedDate={endDate} onChange={this.handleChange('endDate')} label={"End"} className={classes.dateTimePicker} />
                                        <TextField
                                            select
                                            id="refreshInterval"
                                            name="refreshInterval"
                                            label="Refresh Interval"
                                            fullWidth
                                            value={refreshInterval}
                                            margin="normal"
                                            variant="outlined"
                                            onChange={this.handleChange('refreshInterval')}>
                                            {refreshIntervals.map((ri) => (
                                                <MenuItem key={'ri_-_-_'+ri} value={ri}>{ri}</MenuItem>
                                            ))}
                                        </TextField>
                                </Grid>
                                {quickRanges.map(qr => (
                                <Grid item xs={12} sm={2}>
                                   {qr.map(q => (
                                       <Button size="small" variant="text" onClick={this.setRange(q)}>{q}</Button>
                                   ))}
                                </Grid>
                                ))}
                            </Grid>
                        </DialogContentText>
                    </DialogContent>
                <DialogActions>
                    <Button onClick={this.handleClose} color="primary">
                    Cancel
                    </Button>
                    <Button onClick={this.handleClose} color="primary">
                    Subscribe
                    </Button>
                </DialogActions>
                </Dialog>
                {/* </Paper> */}
                {/* </Popover> */}
              {/* </Grow> */}
            {/* )} */}
          {/* </Popper> */}
            </React.Fragment>
            </NoSsr>
        );
    }
}

export default withStyles(styles)(GrafanaDateRangePicker);