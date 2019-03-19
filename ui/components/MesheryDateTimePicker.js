import { MuiPickersUtilsProvider, DateTimePicker } from 'material-ui-pickers';
import MomentUtils from '@date-io/moment';
import PropTypes from 'prop-types';
import { Component } from 'react';

class MesheryDateTimePicker extends Component {

    render() {
        const {selectedDate, onChange, label} = this.props;
        return (
            <MuiPickersUtilsProvider utils={MomentUtils}>
            <div>
                <DateTimePicker value={selectedDate} onChange={onChange} label={label} />
            </div>
            </MuiPickersUtilsProvider>
        );
    }
}

MesheryDateTimePicker.propTypes = {
    label: PropTypes.string.isRequired,
    selectedDate: PropTypes.object.isRequired,
    onChange: PropTypes.func.isRequired,
};

export default MesheryDateTimePicker;