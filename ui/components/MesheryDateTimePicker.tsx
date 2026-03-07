import React, { useState } from 'react';
import { TextField } from '@sistent/sistent';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import moment, { Moment } from 'moment';

interface MesheryDateTimePickerProps {
  label: string;
  selectedDate?: Moment | string | null;
  // eslint-disable-next-line no-unused-vars
  onChange: (date: Moment) => void;
  className?: string;
  disabled?: boolean;
}

const MesheryDateTimePicker: React.FC<MesheryDateTimePickerProps> = ({
  selectedDate,
  onChange,
  label,
  className,
  disabled = false,
}) => {
  const [hasError, setHasError] = useState<boolean>(false);
  const dateFormat = 'YYYY-MM-DD, hh:mm:ss a';

  const handleError = (error: unknown): void => {
    console.error('Error:', error);
    setHasError(Boolean(error));
  };

  const handleDateChange = (date: Moment | null): void => {
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
    const handleNativeChange = (event: React.ChangeEvent<HTMLInputElement>): void => {
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

export default MesheryDateTimePicker;
