import { MuiPickersUtilsProvider, DateTimePicker } from 'material-ui-pickers';
import MomentUtils from '@date-io/moment';
import PropTypes from 'prop-types';
import { Component } from 'react';

class MesheryDateTimePicker extends Component {

    render() {
        const {selectedDate, onChange, label, className, disabled} = this.props;
        return (
            <MuiPickersUtilsProvider utils={MomentUtils}>
            <div className={className}>
                <DateTimePicker disabled={disabled} value={selectedDate} onChange={onChange} 
                    label={label} variant="outlined" fullWidth format={"MMMM Do, YYYY hh:mm:ss a"} />
            </div>
            </MuiPickersUtilsProvider>
        );
    }
}

MesheryDateTimePicker.propTypes = {
    label: PropTypes.string.isRequired,
    selectedDate: PropTypes.object.isRequired,
    onChange: PropTypes.func.isRequired,
    className: PropTypes.object.isRequired,
};

export default MesheryDateTimePicker;