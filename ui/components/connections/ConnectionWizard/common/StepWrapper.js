import React from 'react';
import { Box, Typography, Button, styled } from '@sistent/sistent';

const StepContainer = styled(Box)({
  display: 'flex',
  flexDirection: 'column',
  height: '100%',
  padding: '24px',
  width: '100%',
});

const StepHeader = styled(Box)({
  marginBottom: '24px',
});

const StepContent = styled(Box)({
  flex: 1,
  overflow: 'auto',
});

const StepFooter = styled(Box)({
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  marginTop: '24px',
  paddingTop: '16px',
  borderTop: '1px solid',
  borderColor: 'divider',
});

export const StepWrapper = ({
  title,
  description,
  children,
  onNext,
  onBack,
  onCancel,
  nextButtonText = 'Next',
  backButtonText = 'Back',
  cancelButtonText = 'Cancel',
  nextDisabled = false,
  backDisabled = false,
  showBackButton = true,
  showNextButton = true,
  showCancelButton = true,
  customFooter = null,
}) => {
  return (
    <StepContainer>
      <StepHeader>
        <Typography variant="h5" component="h2" gutterBottom>
          {title}
        </Typography>
        {description && (
          <Typography variant="body1" color="textSecondary">
            {description}
          </Typography>
        )}
      </StepHeader>

      <StepContent>{children}</StepContent>

      {customFooter || (
        <StepFooter>
          <Box>
            {showCancelButton && (
              <Button variant="outlined" onClick={onCancel} sx={{ mr: 2 }}>
                {cancelButtonText}
              </Button>
            )}
            {showBackButton && (
              <Button variant="outlined" onClick={onBack} disabled={backDisabled}>
                {backButtonText}
              </Button>
            )}
          </Box>

          <Box>
            {showNextButton && (
              <Button variant="contained" onClick={onNext} disabled={nextDisabled}>
                {nextButtonText}
              </Button>
            )}
          </Box>
        </StepFooter>
      )}
    </StepContainer>
  );
};
