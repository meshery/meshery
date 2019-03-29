import React from 'react';
import { NoSsr, Popper, Paper } from '@material-ui/core';
import dynamic from 'next/dynamic'
const ReactCountdownClock = dynamic(() => import('react-countdown-clock'), {
  ssr: false
})

class LoadTestTimerDialog extends React.Component {

  render() {
      const {countDownComplete, t, container, open} = this.props;
    return (
      <NoSsr>
        {/* <Menu
          // id="simple-menu"
          anchorEl={container}
          open={open}
          // onClose={this.handleClose}
        > */}
        <Popper open={open} anchorEl={container} 
           placement='bottom' 
            modifiers={{
              flip: {
                enabled: false,
              },
            }}>
                {/* <Paper> */}
                  {/* <ClickAwayListener onClickAway={this.handleClose}></ClickAwayListener> */}

          {/* <MenuItem onClick={this.handleClose}>Profile</MenuItem>
          <MenuItem onClick={this.handleClose}>My account</MenuItem>
          <MenuItem onClick={this.handleClose}>Logout</MenuItem> */}
          <ReactCountdownClock seconds={t * 60}
                        color="#667C89"
                        alpha={0.9}
                        size={400}
                        onComplete={countDownComplete} 
                        />
                        {/* </Paper> */}
                        </Popper>
        {/* </Menu> */}
        
        {/* <Dialog onClose={this.handleTimerDialogClose} 
            {...other} hideBackdrop={true} >
            <DialogContent>
                    
            </DialogContent>
        </Dialog> */}
      </NoSsr>
    );
  }
}

export default LoadTestTimerDialog;
