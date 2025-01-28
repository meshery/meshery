import { TextField } from '@layer5/sistent';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import PropTypes from 'prop-types';
import { useState } from 'react';
import moment from 'moment';

const MesheryDateTimePicker = ({ selectedDate, onChange, label, className, disabled }) => {
  const [hasError, setHasError] = useState(false);
  const dateFormat = 'YYYY-MM-DD, hh:mm:ss a';

  const handleError = (error) => {
    console.error('Error:', error);
    setHasError(Boolean(error));
  };

  const handleDateChange = (date) => {
    if (moment.isMoment(date) && date.isValid()) {
      setHasError(false);
      onChange(date);
    } else {
      setHasError(true);
    }
  };

  // Fallback to a native datetime picker when there's an error
  if (hasError) {
    const dateVal = selectedDate
      ? moment(selectedDate).format(moment.HTML5_FMT.DATETIME_LOCAL_SECONDS)
      : '';
    const handleNativeChange = (event) => {
      const nativeDate = moment(event.target.value, moment.HTML5_FMT.DATETIME_LOCAL_SECONDS);
      handleDateChange(nativeDate);
    };

    return (
      <div className={className}>
        <TextField
          disabled={disabled}
          label={label}
          type="datetime-local"
          value={dateVal}
          onChange={handleNativeChange}
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
        value={selectedDate ? moment(selectedDate) : null}
        onChange={handleDateChange}
        label={label}
        format={dateFormat}
        slotProps={{ textField: { fullWidth: true } }}
        onError={handleError}
      />
    </div>
  );
};

MesheryDateTimePicker.propTypes = {
  label: PropTypes.string.isRequired,
  selectedDate: PropTypes.oneOfType([PropTypes.object, PropTypes.string]),
  onChange: PropTypes.func.isRequired,
  className: PropTypes.string,
  disabled: PropTypes.bool,
};

export default MesheryDateTimePicker;
