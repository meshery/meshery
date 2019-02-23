import React from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import { NoSsr } from '@material-ui/core';
import dynamic from 'next/dynamic'
const ReactCountdownClock = dynamic(() => import('react-countdown-clock'), {
  ssr: false
})

class LoadTestTimerDialog extends React.Component {

  render() {
      const {countDownComplete, ...other} = this.props;
    return (
      <NoSsr>
        <Dialog onClose={this.handleTimerDialogClose} 
            // aria-labelledby="timer-dialog-title" 
            aria-describedby="alert-dialog-description"
            {...other}>
            {/* <DialogTitle id="timer-dialog-title">Set backup account</DialogTitle> */}
            <DialogContent>
                <DialogContentText id="alert-dialog-description">
                    <ReactCountdownClock seconds={120}
                        color="#029BE5"
                        alpha={0.9}
                        size={300}
                        onComplete={countDownComplete} />
                </DialogContentText>
            </DialogContent>
        </Dialog>
      </NoSsr>
    );
  }
}

// LoadTestTimerDialog.propTypes = {
//   classes: PropTypes.object.isRequired,
// };

export default LoadTestTimerDialog;
