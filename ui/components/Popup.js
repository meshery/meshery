import React from 'react'
import { Dialog, DialogTitle, DialogContent, makeStyles, Typography } from '@material-ui/core';

const useStyles = makeStyles(theme => ({
  dialogWrapper: {
    padding: theme.spacing(2),
    position: 'absolute',
  },
  DialogTitle: {
    paddingRight: '0px'
  }
}))

export default function Popup(props) {

  const {title, children, openPopup, setOpenPopup} = props;
  const classes = useStyles();

  return (
    <Dialog open={openPopup} maxWidth="md" classes={{ paper: classes.dialogWrapper}}>
      <DialogTitle className={classes.DialogTitle}>
        <div style={{display: 'flex'}}>
          <Typography varient="h6" component = "div" style={{flexGrow:1}}>
            {title}
          </Typography>
        </div>
      </DialogTitle>
      <DialogContent>
        {children}
      </DialogContent>
    </Dialog>
  )
}

