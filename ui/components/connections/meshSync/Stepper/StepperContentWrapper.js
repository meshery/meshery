import React from 'react';
import { Box, Typography, Divider, Button, styled } from '@layer5/sistent';

const ContentHeader = ({ title }) => (
  <>
    <Typography variant="subtitle1" style={{ fontWeight: 'bold' }}>
      {title}
    </Typography>
    <Divider style={{ margin: '1rem 0 1rem 0' }} />
  </>
);

const ContentContainer = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  width: '60%',
  justifyContent: 'space-between',
  margin: '1rem 0 0 2rem',
  '@media (max-width:780px)': {
    width: 'auto',
    margin: '1rem 0',
    justifyContent: 'flex-start',
  },
});

const ContentBody = styled(Box)({
  marginBottom: '2rem',
});

const CancelButton = styled(Button)(({ theme }) => ({
  backgroundColor: theme.palette.background.neutral.default,
  color: 'white',
  margin: 0,
}));

const StepperButton = styled(Button)(({ theme }) => ({
  marginTop: '1rem',
  color: 'white',
  backgroundColor: theme.palette.background.brand.default,
  '&:hover': {
    backgroundColor: theme.palette.background.brand.default,
  },
  '&.Mui-disabled': {
    cursor: 'not-allowed',
  },
}));

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
  return (
    <ContentContainer>
      <div>
        <ContentHeader title={title} />
        {subtitle && (
          <ContentBody>
            <Typography variant="body2">{subtitle}</Typography>
          </ContentBody>
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
            <Typography style={{ fontSize: '1rem' }}>{tips}</Typography>
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
          <CancelButton
            style={{
              margin: 0,
            }}
            onClick={cancelCallback}
            variant="contained"
          >
            Cancel
          </CancelButton>
        )}

        {btnText && (
          <StepperButton
            onClick={handleCallback}
            disabled={disabled}
            variant="contained"
            style={{ margin: 0 }}
          >
            {btnText}
          </StepperButton>
        )}
      </Box>
    </ContentContainer>
  );
};

export default StepperContent;
