import React from 'react';
import {
  withStyles,
  Typography,
  Button,
  Dialog,
  DialogActions,
  DialogContentText,
  DialogContent,
  DialogTitle
} from "@material-ui/core";

const styles = () => ({
  title: {
    textAlign: 'center',
    minWidth: 400,
    padding: '10px',
    color: '#fff',
    backgroundColor: '#607d8b'
  },
  subtitle: {
    minWidth: 400,
    overflowWrap: 'anywhere',
    textAlign: 'center',
    padding: '5px'
  },
  actions: {
    display: 'flex',
    justifyContent: 'center',
  },
  button: {
    margin: '8px 0px',
  }
});

class PromptComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      show: false,
      title: "",
      subtitle: "",
      options: []
    }
    this.promiseInfo = {};
  }

  show = async (passed) => {
    return new Promise((resolve, reject) => {
      this.promiseInfo = {
        resolve,
        reject
      };
      this.setState({
        title: passed.title,
        subtitle: passed.subtitle,
        options: passed.options,
        show: true
      });
    });
  };

  hide = () => {
    this.setState({
      show: false
    });
  };

  render() {
    const { show, options, title, subtitle } = this.state;
    const { classes } = this.props;
    const { resolve } = this.promiseInfo;
    return (
      <div className={classes.root}>
        <Dialog
          open={show}
          onClose={this.hide}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
          className={classes.dialogBox}
        >
          {title !== "" &&
            <DialogTitle id="alert-dialog-title" className={classes.title}>
              <b>{title}</b>
            </DialogTitle>
          }
          {subtitle !== "" &&
            <DialogContent>
              <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
                <Typography variant="body1">
                  {subtitle}
                </Typography>
              </DialogContentText>
            </DialogContent>
          }
          <DialogActions className={classes.actions}>
            {
              options?.map((option, index) => {
                return (
                  <Button onClick={() => {
                    this.hide();
                    resolve(option);
                  }} key={index} className={classes.button}
                  type="submit"
                  variant="contained"
                  color="primary">
                    {option}
                  </Button>
                )
              })
            }
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default withStyles(styles)(PromptComponent);
