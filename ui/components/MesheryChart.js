import React from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import { NoSsr } from '@material-ui/core';
// import dynamic from 'next/dynamic'
// const ReactCountdownClock = dynamic(() => import('react-countdown-clock'), {
//   ssr: false
// })
import { Line } from 'react-chartjs-2';

class MesheryChart extends React.Component {

  render() {
      const {countDownComplete, ...other} = this.props;
    return (
      <NoSsr>
        <Line />
      </NoSsr>
    );
  }
}

// LoadTestTimerDialog.propTypes = {
//   classes: PropTypes.object.isRequired,
// };

export default MesheryChart;
