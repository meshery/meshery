import { DateTimePicker } from '@material-ui/pickers';
import PropTypes from 'prop-types';
import { Component } from 'react';
import { TextField } from '@material-ui/core';
import moment from 'moment';

class MesheryDateTimePicker extends Component {
  constructor(props) {
    super(props);
    this.dateFormat = 'YYYY-MM-DD, hh:mm:ss a';
    this.state = { hasError : false, };
  }

  static getDerivedStateFromError() {
    return { hasError : true };
  }

  componentDidCatch(error, info) {
    console.log(`error: ${error}, info: ${info}`);
  }

  dateChange() {
    const self = this;
    return (event) => {
      self.props.onChange(moment(event.target.value));
    };
  }

  render() {
    const {
      selectedDate, onChange, label, className, disabled,
    } = this.props;
    const { hasError } = this.state;

    if (hasError) {
      const dateVal = moment(selectedDate).format(moment.HTML5_FMT.DATETIME_LOCAL_SECONDS);
      return (
        <div className={className}>
          <TextField
            disabled={disabled}
            label={label}
            type="datetime-local"
            // defaultValue="2017-05-24T10:30"
            value={dateVal}
            onChange={this.dateChange()}
            variant="outlined"
            fullWidth
            // className={classes.textField}
            InputLabelProps={{ shrink : true, }}
          />
        </div>
      );
    }
    return (
      <div className={className}>
        <DateTimePicker
          disabled={disabled}
          value={selectedDate}
          onChange={onChange}
          label={label}
          variant="outlined"
          fullWidth
          format={this.dateFormat}
        />
      </div>
    );
  }
}

MesheryDateTimePicker.propTypes = {
  label : PropTypes.string.isRequired,
  selectedDate : PropTypes.object.isRequired,
  onChange : PropTypes.func.isRequired,
  className : PropTypes.object.isRequired,
};

export default MesheryDateTimePicker;
