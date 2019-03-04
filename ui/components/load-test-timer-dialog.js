import React from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import { NoSsr } from '@material-ui/core';
import dynamic from 'next/dynamic'
const ReactCountdownClock = dynamic(() => import('react-countdown-clock'), {
  ssr: false
})

class LoadTestTimerDialog extends React.Component {

  render() {
      const {countDownComplete, t, ...other} = this.props;
    return (
      <NoSsr>
        <Dialog onClose={this.handleTimerDialogClose} 
            {...other}>
            <DialogContent>
                    <ReactCountdownClock seconds={t * 60}
                        color="#667C89"
                        alpha={0.9}
                        size={300}
                        onComplete={countDownComplete} />
            </DialogContent>
        </Dialog>
      </NoSsr>
    );
  }
}

export default LoadTestTimerDialog;
