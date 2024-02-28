import React, { forwardRef, useRef, useImperativeHandle, useState } from 'react';
import {
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
  IconButton,
  withStyles,
} from '@material-ui/core';
import theme from '../themes/app';
import { CustomTextTooltip } from './MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import { getHyperLinkDiv } from './MesheryMeshInterface/PatternService/helper';

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

const IconButtonWrapper = styled(IconButton)(({ theme }) => ({
  position: 'absolute',
  right: 10,
  color: theme.palette.secondary.focused,
}));

const PromptComponent = forwardRef(({ classes }, ref) => {
  const [show, setShow] = useState(false);
  const [title, setTitle] = useState('');
  const [subtitle, setSubtitle] = useState('');
  const [options, setOptions] = useState([]);
  const [isChecked, setIsChecked] = useState(false);
  const [showCheckbox, setShowCheckbox] = useState(false);
  const [showInfoIcon, setShowInfoIcon] = useState(null);
  const [variant, setVariant] = useState(null);

  const handleCheckboxChange = () => {
    setIsChecked((prevState) => !prevState);
  };

  const hide = () => {
    setShow(false);
  };

  const resolveFn = useRef(null);
  useImperativeHandle(ref, () => ({
    show: async (passed) => {
      return new Promise((resolve) => {
        setVariant(passed.variant);
        setTitle(passed.title);
        setSubtitle(passed.subtitle);
        setOptions(passed.options);
        setShowCheckbox(!!passed.showCheckbox);
        setShow(true);
        setShowInfoIcon(passed.showInfoIcon || null);
        resolveFn.current = resolve;
      });
    },
  }));

  return (
    <div className={classes.root}>
      <Dialog
        open={show}
        onClose={hide}
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
                    onChange={handleCheckboxChange}
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
                hide();
                resolveFn.current && resolveFn.current(options[1]);
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
              hide();
              resolveFn.current && resolveFn.current(options[0]);
            }}
            key={options[0]}
            promptVariant={variant}
            style={variant && { backgroundColor: theme.palette.secondary[variant] }}
            type="submit"
            variant="contained"
          >
            <Typography variant body2>
              {options[0]}{' '}
            </Typography>
          </PromptActionButton>
          {showInfoIcon && (
            <CustomTextTooltip
              backgroundColor="#3C494F"
              placement="top"
              interactive={true}
              style={{ whiteSpace: 'pre-line' }}
              title={getHyperLinkDiv(showInfoIcon)}
            >
              <IconButtonWrapper color="primary">
                <InfoOutlinedIcon />
              </IconButtonWrapper>
            </CustomTextTooltip>
          )}
        </DialogActions>
      </Dialog>
    </div>
  );
});
PromptComponent.displayName = 'PromptComponent';
export default withStyles(styles)(PromptComponent);
