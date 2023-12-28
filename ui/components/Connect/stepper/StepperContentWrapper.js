import React from 'react';
import { Box, Typography, Divider, Button } from '@material-ui/core';
import { useStyles } from '../styles';
import { useTheme } from '@material-ui/core/styles';
import { Colors } from '../../../themes/app';

const ContentHeader = ({ title }) => (
  <>
    <Typography style={{ fontSize: '1.4rem', fontWeight: 'bold' }}>{title}</Typography>
    <Divider style={{ margin: '1rem 0 1.5rem 0' }} />
  </>
);

const StepperContent = ({
  title,
  children,
  subtitle,
  tips,
  btnText,
  cancel,
  handleCallback,
  cancelCallback,
  disabled,
}) => {
  const theme = useTheme();
  const classes = useStyles();
  return (
    <Box className={classes.contentContainer}>
      <div>
        <ContentHeader title={title} />
        {subtitle && (
          <Box className={classes.contentBody}>
            <Typography variant="subtitle1" style={{ fontSize: '1rem' }}>
              {subtitle}
            </Typography>
          </Box>
        )}
        {children}
        {tips && (
          <Box
            style={{
              background: 'rgba(0, 211, 169, 0.05)',
              padding: '0.6rem',
              margin: '2rem 0',
            }}
          >
            <Typography variant="body1">{tips}</Typography>
          </Box>
        )}
      </div>
      <Box
        style={{
          display: 'flex',
          justifyContent: cancel ? 'space-between' : 'flex-end',
        }}
      >
        {cancel && (
          <Button
            className={classes.cancelButton}
            style={{
              background: Colors.charcoal,
              color: theme.palette.secondary.primaryModalText,
              margin: 0,
            }}
            onClick={cancelCallback}
            variant="contained"
          >
            Cancel
          </Button>
        )}

        {btnText && (
          <Button
            onClick={handleCallback}
            disabled={disabled}
            variant="contained"
            style={{ margin: 0 }}
            className={classes.stepperButton}
          >
            {btnText}
          </Button>
        )}
      </Box>
    </Box>
  );
};

export default StepperContent;
