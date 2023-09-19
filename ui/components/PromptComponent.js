import React from 'react';
import {
  withStyles,
  Typography,
  Button,
  Dialog,
  DialogActions,
  DialogContentText,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  Checkbox,
} from '@material-ui/core';
import classNames from 'classnames';

const styles = (theme) => ({
  title: {
    textAlign: 'center',
    minWidth: 400,
    padding: '10px',
    color: '#fff',
    backgroundColor:
      theme.palette.type === 'dark' ? theme.palette.secondary.headerColor : '#396679',
  },
  subtitle: {
    minWidth: 400,
    overflowWrap: 'anywhere',
    textAlign: 'center',
    padding: '5px',
    color: theme.palette.secondary.text,
  },
  actions: {
    display: 'flex',
    justifyContent: 'space-evenly',
  },
  button0: {
    margin: theme.spacing(0.5),
    padding: theme.spacing(1),
    borderRadius: 5,
    backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
    '&:hover': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
      boxShadow:
        '0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)',
    },
    minWidth: 100,
  },
  button1: {
    margin: theme.spacing(0.5),
    padding: theme.spacing(1),
    borderRadius: 5,
    backgroundColor: '#e0e0e0',
    color: 'rgba(0, 0, 0, 0.87)',
    '&:hover': {
      backgroundColor: '#d5d5d5',
      boxShadow:
        '0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)',
    },
    minWidth: 100,
  },
  resetButton: {
    backgroundColor: '#8F1F00',
    '&:hover': {
      backgroundColor: '#B32700',
    },
  },
  checkboxLabelStyle: {
    fontSize: '1rem',
  },
  checkbox: {
    color: theme.palette.secondary.focused,
    '&$checked': {
      color: theme.palette.secondary.focused,
    },
  },
});

class PromptComponent extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      show: false,
      title: '',
      subtitle: '',
      options: [],
      isChecked: false,
      showCheckbox: false,
    };
    this.promiseInfo = {};
  }

  show = async (passed) => {
    return new Promise((resolve, reject) => {
      this.promiseInfo = { resolve, reject };
      this.setState({
        title: passed.title,
        subtitle: passed.subtitle,
        options: passed.options,
        showCheckbox: !!passed.showCheckbox,
        show: true,
      });
    });
  };

  hide = () => {
    this.setState({ show: false });
  };

  handleCheckboxChange = () => {
    this.setState((prevState) => ({
      isChecked: !prevState.isChecked,
    }));
  };

  getCheckboxState = () => {
    return this.state.isChecked;
  };

  render() {
    const { show, options, title, subtitle, isChecked, showCheckbox } = this.state;
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
          {title !== '' && (
            <DialogTitle id="alert-dialog-title" className={classes.title}>
              <b>{title}</b>
            </DialogTitle>
          )}
          {subtitle !== '' && (
            <DialogContent>
              <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
                <Typography variant="body1">{subtitle}</Typography>
              </DialogContentText>
              {showCheckbox && (
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={isChecked}
                      onChange={this.handleCheckboxChange}
                      className={classes.checkbox}
                      color="primary"
                    />
                  }
                  label={<span className={classes.checkboxLabelStyle}>Do not show again</span>}
                />
              )}
            </DialogContent>
          )}
          <DialogActions className={classes.actions}>
            <Button
              onClick={() => {
                this.hide();
                resolve(options[1]);
              }}
              key={options[1]}
              className={classes.button1}
            >
              <Typography variant body2>
                {' '}
                {options[1]}{' '}
              </Typography>
            </Button>
            <Button
              onClick={() => {
                this.hide();
                resolve(options[0]);
              }}
              key={options[0]}
              className={
                options[0]?.toLowerCase() === 'reset'
                  ? classNames(classes.button0, classes.resetButton)
                  : classes.button0
              }
              type="submit"
              variant="contained"
              color="primary"
            >
              <Typography variant body2>
                {' '}
                {options[0]}{' '}
              </Typography>
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default withStyles(styles)(PromptComponent);
