import { DateTimePicker,  } from '@material-ui/pickers';

import PropTypes from 'prop-types';
import { Component } from 'react';

class MesheryDateTimePicker extends Component {

    render() {
        const {selectedDate, onChange, label, className, disabled} = this.props;
        return (
                <DateTimePicker disabled={disabled} value={selectedDate} onChange={onChange} 
                    label={label} variant="dialog" fullWidth format={"MMMM Do, YYYY hh:mm:ss a"} />
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