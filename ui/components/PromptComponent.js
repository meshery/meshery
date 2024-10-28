import React, { useState, useRef, forwardRef, useImperativeHandle } from 'react';
import {
  withStyles,
  Typography,
  DialogContentText,
  FormControlLabel,
  Checkbox,
} from '@material-ui/core';
import theme from '../themes/app';
import {
  Box,
  Modal,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
} from '@layer5/sistent';
import { UsesSistent } from './SistentWrapper';

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

export const PROMPT_VARIANTS = {
  WARNING: 'warning',
  DANGER: 'danger',
  SUCCESS: 'success',
  CONFIRMATION: 'confirmation',
};

const PromptComponent = forwardRef(({ classes, variant: initialVariant }, ref) => {
  const [isOpen, setIsOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [subtitle, setSubtitle] = useState('');
  const [options, setOptions] = useState([]);
  const [isChecked, setIsChecked] = useState(false);
  const [showCheckbox, setShowCheckbox] = useState(false);
  const [showInfoIcon, setShowInfoIcon] = useState(null);
  const [variant, setVariant] = useState(initialVariant);
  const promiseInfo = useRef({});

  PromptComponent.displayName = 'PromptComponent';

  const show = async (passed) => {
    return new Promise((resolve, reject) => {
      promiseInfo.current = { resolve, reject };
      setVariant(passed.variant);
      setTitle(passed.title);
      setSubtitle(passed.subtitle);
      setOptions(passed.options);
      setShowCheckbox(!!passed.showCheckbox);
      setShowInfoIcon(passed.showInfoIcon || null);
      setIsOpen(true);
    });
  };

  useImperativeHandle(ref, () => ({
    show,
    getCheckboxState: () => isChecked,
  }));

  const hide = () => {
    setIsOpen(false);
  };

  const handleCheckboxChange = () => {
    setIsChecked((prev) => !prev);
  };

  const { resolve } = promiseInfo.current;

  return (
    <div className={classes.root}>
      <UsesSistent>
        <Modal open={isOpen} closeModal={hide} title={title}>
          {subtitle !== '' && (
            <ModalBody>
              <DialogContentText id="alert-dialog-description" className={classes.subtitle}>
                <Typography variant="body1">{subtitle}</Typography>
              </DialogContentText>
              {showCheckbox && (
                <FormControlLabel
                  control={
                    <Checkbox
                      checked={isChecked}
                      onChange={handleCheckboxChange}
                      className={classes.checkbox}
                      color="primary"
                    />
                  }
                  label={<span className={classes.checkboxLabelStyle}>Do not show again</span>}
                />
              )}
            </ModalBody>
          )}
          <ModalFooter variant="filled" helpText={showInfoIcon}>
            <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'end' }}>
              {options.length > 1 && (
                <ModalButtonSecondary
                  onClick={() => {
                    hide();
                    resolve(options[1]);
                  }}
                  key={options[1]}
                >
                  <Typography variant="body2">{options[1]}</Typography>
                </ModalButtonSecondary>
              )}
              <ModalButtonPrimary
                color="primary"
                onClick={() => {
                  hide();
                  resolve(options[0]);
                }}
                key={options[0]}
                promptVariant={variant}
                style={variant && { backgroundColor: theme.palette.secondary[variant] }}
                type="submit"
                variant="contained"
              >
                <Typography variant="body2">{options[0]}</Typography>
              </ModalButtonPrimary>
            </Box>
          </ModalFooter>
        </Modal>
      </UsesSistent>
    </div>
  );
});

export default withStyles(styles)(PromptComponent);
