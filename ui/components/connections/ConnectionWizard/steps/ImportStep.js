import React, { useState } from 'react';
import { Box, Typography, FormGroup, TextField, InputAdornment, Alert } from '@sistent/sistent';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { StepWrapper } from '../common/StepWrapper';

const CONNECTION_TYPE_CONFIG = {
  kubernetes: {
    title: 'Import Kubeconfig File',
    description: 'Upload your kubeconfig file to discover available clusters',
    fileTypes: '',
    placeholder: 'Upload kubeconfig',
    helpText: 'commonly found at ~/.kube/config',
    validationMessage: 'Please select a valid kubeconfig file',
  },
};

export const ImportStep = ({
  connectionType = 'kubernetes',
  wizardData,
  setWizardData,
  onNext,
  onCancel,
  isFirstStep,
}) => {
  console.log('wizardData', wizardData);
  const [selectedFile, setSelectedFile] = useState(wizardData?.uploadedFile || null);
  const [validationError, setValidationError] = useState(null);
  const { notify } = useNotification();

  const config = CONNECTION_TYPE_CONFIG[connectionType] || CONNECTION_TYPE_CONFIG.kubernetes;

  const handleFileChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setSelectedFile(file);
      setValidationError(null);

    
      const formData = new FormData();
      formData.append('k8sfile', file);

      setWizardData((prev) => ({
        ...prev,
        uploadedFile: file,
        uploadedFormData: formData, // Store FormData for API calls
      }));
    }
  };

  const validateFile = (file) => {
    if (!file) {
      throw new Error('No file selected');
    }

    const fileName = file.name;
    const invalidExtensions = /^.*\.(jpg|gif|jpeg|pdf|png|svg)$/i;

    if (invalidExtensions.test(fileName)) {
      throw new Error('Invalid file selected. Please select a valid kubeconfig file');
    }

    if (file.size > 10 * 1024 * 1024) {
      throw new Error('File too large. Maximum size is 10MB');
    }

    return true;
  };

  const handleNext = () => {
    if (!selectedFile) {
      setValidationError(config.validationMessage);
      return;
    }

    try {

      validateFile(selectedFile);

      notify({
        message: 'File uploaded successfully!',
        event_type: EVENT_TYPES.SUCCESS,
      });

      onNext();
    } catch (error) {
      setValidationError(error.message);
      notify({
        message: `File validation failed: ${error.message}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  };

  const handleFileSelect = () => {
    document.getElementById('connection-file-input')?.click();
  };

  return (
    <StepWrapper
      title={config.title}
      description={config.description}
      onNext={handleNext}
      onCancel={onCancel}
      nextDisabled={!selectedFile}
      nextButtonText="Continue"
      showBackButton={!isFirstStep}
    >
      <Box sx={{ mt: 3 }}>
        <FormGroup>
          <input
            id="connection-file-input"
            type="file"
            accept={config.fileTypes}
            onChange={handleFileChange}
            style={{ display: 'none' }}
          />

          <TextField
            id="k8sfileLabelText"
            value={selectedFile ? selectedFile.name : ''}
            placeholder={config.placeholder}
            variant="outlined"
            fullWidth
            onClick={handleFileSelect}
            style={{ cursor: 'pointer', marginBottom: '16px' }}
            InputProps={{
              readOnly: true,
              endAdornment: (
                <InputAdornment position="end">
                  <CloudUploadIcon />
                </InputAdornment>
              ),
            }}
          />

          <Typography variant="body2" color="textSecondary">
            {config.helpText}
          </Typography>
        </FormGroup>

        {validationError && (
          <Alert severity="error" sx={{ mt: 2 }}>
            {validationError}
          </Alert>
        )}

        {selectedFile && !validationError && (
          <Alert severity="info" sx={{ mt: 2 }}>
            <Typography variant="body2">
              <strong>Selected file:</strong> {selectedFile.name}
            </Typography>
            <Typography variant="body2">
              <strong>Size:</strong> {(selectedFile.size / 1024).toFixed(2)} KB
            </Typography>
          </Alert>
        )}
      </Box>
    </StepWrapper>
  );
};
