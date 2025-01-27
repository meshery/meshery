import React from 'react';
import { Component } from 'react';
import { NoSsr } from '@mui/material';
import {
  Button,
  TextField,
  MenuItem,
  Grid,
  Dialog,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Divider,
  FormControlLabel,
  Switch,
  IconButton,
  styled,
} from '@layer5/sistent';
import Moment from 'react-moment';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import PropTypes from 'prop-types';
import MesheryDateTimePicker from '../../MesheryDateTimePicker';
import { Close } from '@mui/icons-material';
import { UsesSistent } from '@/components/SistentWrapper';

const StyledDialog = styled(Dialog)(({ theme }) => ({
  '& .MuiDialogContent-root': {
    padding: theme.spacing(2),
  },
}));

const DialogTitleBar = styled('div')(() => ({
  display: 'flex',
  justifyContent: 'space-between',
  backgroundColor: 'rgb(57, 102, 121)',
  color: 'white',
}));

const CloseIconButton = styled(IconButton)(() => ({
  height: '2.3rem',
  marginTop: '0.8rem',
  marginRight: '10px',
  color: 'white',
}));

const CloseIcon = styled(Close)({
  transform: 'rotate(-90deg)',
  '&:hover': {
    transform: 'rotate(90deg)',
    transition: 'all .3s ease-in',
  },
});

const RangeButton = styled(Button)(() => ({
  border: '1px solid rgba(0, 0, 0, 0.23)',
}));

const RangeDialogRow = styled('div')(() => ({
  display: 'flex',
  gap: '4rem',
  '& > *': { flex: '1' },
}));

const DateTimePicker = styled(MesheryDateTimePicker)(({ theme }) => ({
  marginTop: theme.spacing(2),
  marginBottom: theme.spacing(1),
}));

const InnerGrid = styled(Grid)(({ theme }) => ({
  borderTop: `1px solid ${theme.palette.divider}`,
  paddingTop: theme.spacing(2),
  marginTop: theme.spacing(2),
}));

const TimeList = styled(Grid)(() => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'flex-start',
  '& *': { textAlign: 'left' },
}));

const Space = styled('span')(({ theme }) => ({
  margin: theme.spacing(1),
}));

const refreshIntervals = ['off', '5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'];

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
  ],
  [
    'Yesterday',
    'Day before yesterday',
    'This day last week',
    'Previous week',
    'Previous month',
    'Previous year',
  ],
  [
    'Today',
    'Today so far',
    'This week',
    'This week so far',
    'This month',
    'This month so far',
    'This year',
    'This year so far',
  ],
  [
    'Last 5 minutes',
    'Last 15 minutes',
    'Last 30 minutes',
    'Last 1 hour',
    'Last 3 hours',
    'Last 6 hours',
    'Last 12 hours',
    'Last 24 hours',
  ],
];

class GrafanaDateRangePicker extends Component {
  constructor(props) {
    super(props);

    this.state = {
      // startDate: props.startDate,
      // endDate: props.endDate,
      // startGDate: props.from,
      // endGDate: props.to,
      // liveTail: true,
      // refreshInterval: props.refresh,
      open: false,
    };
  }

  handleClick = () => {
    this.setState({ open: true });
  };

  setRange = (range) => () => {
    let startDate;
    let endDate;
    let liveTail;
    let startGDate;
    let endGDate;
    switch (range) {
      case 'Last 2 days':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - 2);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-2d';
        endGDate = 'now';
        break;
      case 'Last 7 days':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - 7);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-7d';
        endGDate = 'now';
        break;
      case 'Last 30 days':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - 30);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-30d';
        endGDate = 'now';
        break;
      case 'Last 90 days':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - 90);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-90d';
        endGDate = 'now';
        break;
      case 'Last 6 months':
        startDate = new Date();
        startDate.setMonth(startDate.getMonth() - 6);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-6M';
        endGDate = 'now';
        break;
      case 'Last 1 year':
        startDate = new Date();
        startDate.setFullYear(startDate.getFullYear() - 1);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-1y';
        endGDate = 'now';
        break;
      case 'Last 2 years':
        startDate = new Date();
        startDate.setFullYear(startDate.getFullYear() - 2);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-2y';
        endGDate = 'now';
        break;
      case 'Last 5 years':
        startDate = new Date();
        startDate.setFullYear(startDate.getFullYear() - 5);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-5y';
        endGDate = 'now';
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
        liveTail = false;
        startGDate = 'now-1d/d';
        endGDate = 'now-1d/d';
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
        liveTail = false;
        startGDate = 'now-2d/d';
        endGDate = 'now-2d/d';
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
        liveTail = false;
        startGDate = 'now-7d/d';
        endGDate = 'now-7d/d';
        break;
      case 'Previous week':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - 6 - ((startDate.getDay() + 8) % 7));
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
        liveTail = false;
        startGDate = 'now-1w/w';
        endGDate = 'now-1w/w';
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
        liveTail = false;
        startGDate = 'now-1M/M';
        endGDate = 'now-1M/M';
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
        liveTail = false;
        startGDate = 'now-1y/y';
        endGDate = 'now-1y/y';
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
        liveTail = false;
        startGDate = 'now/d';
        endGDate = 'now/d';
        break;
      case 'Today so far':
        startDate = new Date();
        startDate.setHours(0);
        startDate.setMinutes(0);
        startDate.setSeconds(0);
        startDate.setMilliseconds(0);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now/d';
        endGDate = 'now';
        break;
      case 'This week':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - ((startDate.getDay() + 7) % 7));
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
        liveTail = false;
        startGDate = 'now/w';
        endGDate = 'now/w';
        break;
      case 'This week so far':
        startDate = new Date();
        startDate.setDate(startDate.getDate() - ((startDate.getDay() + 7) % 7));
        startDate.setHours(0);
        startDate.setMinutes(0);
        startDate.setSeconds(0);
        startDate.setMilliseconds(0);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now/w';
        endGDate = 'now';
        break;
      case 'This month':
        startDate = new Date();
        startDate.setDate(1);
        startDate.setHours(0);
        startDate.setMinutes(0);
        startDate.setSeconds(0);
        startDate.setMilliseconds(0);
        endDate = new Date();
        endDate.setMonth(endDate.getMonth() + 1);
        endDate.setDate(0);
        endDate.setHours(23);
        endDate.setMinutes(59);
        endDate.setSeconds(59);
        endDate.setMilliseconds(999);
        liveTail = false;
        startGDate = 'now/M';
        endGDate = 'now/M';
        break;
      case 'This month so far':
        startDate = new Date();
        startDate.setDate(1);
        startDate.setHours(0);
        startDate.setMinutes(0);
        startDate.setSeconds(0);
        startDate.setMilliseconds(0);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now/M';
        endGDate = 'now';
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
        liveTail = false;
        startGDate = 'now/y';
        endGDate = 'now/y';
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
        liveTail = true;
        startGDate = 'now/y';
        endGDate = 'now';
        break;

      case 'Last 5 minutes':
        startDate = new Date();
        startDate.setMinutes(startDate.getMinutes() - 5);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-5m';
        endGDate = 'now';
        break;
      case 'Last 15 minutes':
        startDate = new Date();
        startDate.setMinutes(startDate.getMinutes() - 15);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-15m';
        endGDate = 'now';
        break;
      case 'Last 30 minutes':
        startDate = new Date();
        startDate.setMinutes(startDate.getMinutes() - 30);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-30m';
        endGDate = 'now';
        break;
      case 'Last 1 hour':
        startDate = new Date();
        startDate.setHours(startDate.getHours() - 1);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-1h';
        endGDate = 'now';
        break;
      case 'Last 3 hours':
        startDate = new Date();
        startDate.setHours(startDate.getHours() - 3);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-3h';
        endGDate = 'now';
        break;
      case 'Last 6 hours':
        startDate = new Date();
        startDate.setHours(startDate.getHours() - 6);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-6h';
        endGDate = 'now';
        break;
      case 'Last 12 hours':
        startDate = new Date();
        startDate.setHours(startDate.getHours() - 12);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-12h';
        endGDate = 'now';
        break;
      case 'Last 24 hours':
        startDate = new Date();
        startDate.setHours(startDate.getHours() - 24);
        endDate = new Date();
        liveTail = true;
        startGDate = 'now-24h';
        endGDate = 'now';
        break;
    }
    // this.setState({startDate, startendDate, liveTail})
    this.props.updateDateRange(
      startGDate,
      startDate,
      endGDate,
      endDate,
      liveTail,
      this.props.refresh,
    );
  };

  // handleToggle = () => {
  //     this.setState(state => ({ open: !state.open }));
  // };

  handleClose = () => {
    // if (this.anchorEl.contains(event.target)) {
    //     return;
    // }
    this.setState({ open: false });
    // this.setState({
    //     anchorEl: null,
    //   });
  };

  handleChange = (name) => (event) => {
    if (name === 'startDate' || name === 'endDate') {
      const { startDate, endDate } = this.state;
      const dt = event.toDate();
      if (name === 'startDate') {
        if (dt > endDate) {
          // this.setState({ [name]: dt, endDate: dt, startGDate: dt.getTime().toString(), endGDate: dt.getTime().toString() });
          this.props.updateDateRange(
            dt.getTime().toString(),
            dt,
            dt.getTime().toString(),
            dt,
            this.props.liveTail,
            this.props.refresh,
          );
          return;
        }
        // this.setState({ [name]: dt,  startGDate: dt.getTime().toString()});
        this.props.updateDateRange(
          dt.getTime().toString(),
          dt,
          this.props.to,
          this.props.endDate,
          this.props.liveTail,
          this.props.refresh,
        );
      } else if (name === 'endDate') {
        if (dt < startDate) {
          // this.setState({ [name]: dt, startDate: dt, startGDate: dt.getTime().toString(), endGDate: dt.getTime().toString() });
          this.props.updateDateRange(
            dt.getTime().toString(),
            dt,
            dt.getTime().toString(),
            dt,
            this.props.liveTail,
            this.props.refresh,
          );
          return;
        }
        // this.setState({ [name]: dt, endGDate: dt.getTime().toString() });
        // this.props.updateDateRange(dt.getTime().toString(), dt, dt.getTime().toString(), dt, this.props.liveTail, this.props.refresh);
        this.props.updateDateRange(
          this.props.from,
          this.props.startDate,
          dt.getTime().toString(),
          dt,
          this.props.liveTail,
          this.props.refresh,
        );
      }
      return;
    }
    if (name === 'liveTail') {
      // this.setState({liveTail: event.target.checked});
      if (event.target.checked) {
        this.props.updateDateRange(
          this.props.from,
          this.props.startDate,
          'now',
          this.props.endDate,
          event.target.checked,
          this.props.refresh,
        );
        return;
      }
      this.props.updateDateRange(
        this.props.startDate.getTime().toString(),
        this.props.startDate,
        this.props.endDate.getTime().toString(),
        this.props.endDate,
        event.target.checked,
        this.props.refresh,
      );
      return;
    }
    if (name === 'refresh') {
      this.props.updateDateRange(
        this.props.from,
        this.props.startDate,
        this.props.to,
        this.props.endDate,
        this.props.liveTail,
        event.target.value,
      );
      return;
    }
    this.setState({ [name]: event.target.value });
  };

  render() {
    const { open } = this.state;
    const { startDate, endDate, liveTail, refresh } = this.props;
    return (
      <UsesSistent>
        <NoSsr>
          <React.Fragment>
            <RangeButton variant="filled" onClick={this.handleClick}>
              <AccessTimeIcon sx={{ marginRight: '0.25rem', fontSize: '1.15rem' }} />
              <Moment format="LLLL">{startDate}</Moment>
              <Space>-</Space>
              {liveTail ? 'now' : <Moment format="LLLL">{endDate}</Moment>}
              <Space>,{refresh}</Space>
            </RangeButton>
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

            <StyledDialog
              open={open}
              onClose={this.handleClose}
              scroll="paper"
              aria-labelledby="daterange-dialog-title"
              maxWidth="md"
            >
              <DialogTitleBar>
                <DialogTitle id="daterange-dialog-title">Select a Date Range</DialogTitle>
                <CloseIconButton aria-label="close" onClick={this.handleClose}>
                  <CloseIcon />
                </CloseIconButton>
              </DialogTitleBar>
              <DialogContent>
                <DialogContentText>
                  <Grid container>
                    <Grid item xs={12}>
                      Custom Range
                      <RangeDialogRow>
                        <DateTimePicker
                          selectedDate={startDate}
                          onChange={this.handleChange('startDate')}
                          label="Start"
                        />
                        <DateTimePicker
                          disabled={liveTail}
                          selectedDate={endDate}
                          onChange={this.handleChange('endDate')}
                          label="End"
                        />
                      </RangeDialogRow>
                      <RangeDialogRow>
                        <FormControlLabel
                          control={
                            <Switch
                              checked={liveTail}
                              color="primary"
                              onChange={this.handleChange('liveTail')}
                            />
                          }
                          label="Live tail"
                        />
                        <TextField
                          select
                          id="refresh"
                          name="refresh"
                          label="Refresh Interval"
                          fullWidth
                          value={refresh}
                          margin="normal"
                          variant="outlined"
                          onChange={this.handleChange('refresh')}
                          // disabled={liveTail}
                        >
                          {refreshIntervals.map((ri) => (
                            <MenuItem key={`ri_-_-_${ri}`} value={ri}>
                              {ri}
                            </MenuItem>
                          ))}
                        </TextField>
                      </RangeDialogRow>
                    </Grid>
                    <InnerGrid item xs={12}>
                      Quick Ranges
                      <Grid container spacing={0}>
                        {quickRanges.map((qr) => (
                          <TimeList item key={qr.uniqueID} xs={12} sm={3}>
                            {qr.map((q) => (
                              <Button key={q.uniqueID} variant="text" onClick={this.setRange(q)}>
                                {q}
                              </Button>
                            ))}
                          </TimeList>
                        ))}
                      </Grid>
                    </InnerGrid>
                  </Grid>
                </DialogContentText>
              </DialogContent>
              <Divider light variant="fullWidth" />
            </StyledDialog>
            {/* </Paper> */}
            {/* </Popover> */}
            {/* </Grow> */}
            {/* )} */}
            {/* </Popper> */}
          </React.Fragment>
        </NoSsr>
      </UsesSistent>
    );
  }
}

GrafanaDateRangePicker.propTypes = {
  updateDateRange: PropTypes.func.isRequired,
  from: PropTypes.string.isRequired,
  startDate: PropTypes.object.isRequired,
  to: PropTypes.string.isRequired,
  endDate: PropTypes.object.isRequired,
  // liveTail: PropTypes.bool.isRequired,
  refresh: PropTypes.string.isRequired,
};

export default GrafanaDateRangePicker;
