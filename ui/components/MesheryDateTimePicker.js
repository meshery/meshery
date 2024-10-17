import { DateTimePicker } from '@material-ui/pickers';
import PropTypes from 'prop-types';
import { useState, useCallback } from 'react';
import { TextField } from '@material-ui/core';
import moment from 'moment';

const MesheryDateTimePicker = ({ selectedDate, onChange, label, className, disabled }) => {
  const [hasError, setHasError] = useState(false); 
  const dateFormat = 'YYYY-MM-DD, hh:mm:ss a';

  const dateChange = useCallback((event) => {
    try {
      onChange(moment(event.target.value));
    } catch (error) {
      console.error(error);
      setHasError(true); 
    }
  }, [onChange]);

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
          onChange={dateChange}
          variant="outlined"
          fullWidth
          InputLabelProps={{ shrink: true }}
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
        format={dateFormat}
      />
    </div>
  );
};

MesheryDateTimePicker.propTypes = {
  label: PropTypes.string.isRequired,
  selectedDate: PropTypes.object.isRequired,
  onChange: PropTypes.func.isRequired,
  className: PropTypes.object.isRequired,
};

export default MesheryDateTimePicker;
