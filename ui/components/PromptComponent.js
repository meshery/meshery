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
  styled,
} from '@material-ui/core';
import theme from '../themes/app';

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

const PromptActionButton = styled(Button)(({ theme }) => ({
  margin: theme.spacing(0.5),
  padding: theme.spacing(1),
  borderRadius: 5,
  color: '#fff',
  '&:hover': {
    boxShadow:
      '0px 2px 4px -1px rgb(0 0 0 / 20%), 0px 4px 5px 0px rgb(0 0 0 / 14%), 0px 1px 10px 0px rgb(0 0 0 / 12%)',
  },
  minWidth: 100,
}));

export const PROMPT_VARIANTS = {
  WARNING: 'warning',
  DANGER: 'danger',
  SUCCESS: 'success',
  CONFIRMATION: 'confirmation',
};

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
    this.variant = this.props.variant;
  }

  show = async (passed) => {
    return new Promise((resolve, reject) => {
      this.promiseInfo = { resolve, reject };
      console.log('show passed', passed);
      this.variant = passed.variant;
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
            {options.length > 1 && (
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
            )}
            <PromptActionButton
              color="primary"
              onClick={() => {
                this.hide();
                resolve(options[0]);
              }}
              key={options[0]}
              promptVariant={this.variant}
              style={this.variant && { backgroundColor: theme.palette.secondary[this.variant] }}
              type="submit"
              variant="contained"
            >
              <Typography variant body2>
                {options[0]}{' '}
              </Typography>
            </PromptActionButton>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

export default withStyles(styles)(PromptComponent);
