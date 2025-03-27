import React from 'react';
import {
  ModalFooter,
  useStepper,
  CustomizedStepper,
  ModalBody,
  Box,
  ModalButtonSecondary,
  ModalButtonPrimary,
  ComponentIcon,
  Typography,
  FormControl,
  OutlinedInput,
  DownloadIcon,
  useTheme,
  styled,
  Link,
  Chip,
} from '@layer5/sistent';
// import AppRegistrationIcon from '@mui/icons-material/AppRegistration';
import ModelIcon from '@/assets/icons/ModelIcon';
import LanOutlinedIcon from '@mui/icons-material/LanOutlined';
import { TooltipIconButton } from '@/utils/TooltipButton';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import { CSV_TEMPLATE_BASE_URL } from './data';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';

const StyledHeadingBox = styled(Box)({
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'space-between',
});

export const StyledDocsRedirectLink = styled(Link)(({ theme }) => ({
  color: theme.palette.text.brand,
  textDecoration: 'underline',
}));

const StyledUploadSuccess = styled(Box)(({ theme }) => ({
  border: `1px solid ${theme.palette.border.brand}`,
  borderRadius: theme.shape.borderRadius,
  padding: '0.75rem',
  marginTop: '1rem',
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'space-between',
  alignItems: 'center',
  color: theme.palette.text.brand,
}));

const StyledFileChip = styled(Chip)(({ theme }) => ({
  borderRadius: theme.shape.borderRadius,
}));

const CsvStepper = React.memo(({ handleClose, handleGenerateModal }) => {
  const [modelCsvFile, setModelCsvFile] = React.useState(null);
  const [componentCsvFile, setComponentCsvFile] = React.useState(null);
  const [relationshipCsvFile, setRelationshipCsvFile] = React.useState(null);
  const [registerModel] = React.useState(true);
  const fileToBase64 = (file) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = (error) => reject(error);
    });
  };
  const handleModelCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      const base64 = await fileToBase64(file);
      setModelCsvFile({ base64, name: file.name });
    }
  };
  const handleComponentCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      const base64 = await fileToBase64(file);
      setComponentCsvFile({ base64, name: file.name });
    }
  };

  const handleRelationshipCsvFileChange = async (event) => {
    const file = event.target.files[0];
    if (file) {
      const base64 = await fileToBase64(file);
      setRelationshipCsvFile({ base64, name: file.name });
    }
  };

  const handleFinish = () => {
    handleClose();
    handleGenerateModal({
      model_csv: `data:text/csv;base64,${modelCsvFile?.base64?.split(',')[1]}`,
      component_csv: `data:text/csv;base64,${componentCsvFile?.base64?.split(',')[1]}`,
      relationship_csv: relationshipCsvFile
        ? `data:text/csv;base64,${relationshipCsvFile?.base64?.split(',')[1]}`
        : null,
      register: registerModel,
    });
  };
  const theme = useTheme();
  const handleDownload = async (fileName) => {
    try {
      // Convert GitHub raw URL
      const rawUrl = `${CSV_TEMPLATE_BASE_URL}${fileName}`;
      const response = await fetch(rawUrl);
      const blob = await response.blob();

      // Create a download link
      const downloadUrl = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.download = fileName;

      // Trigger download
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);

      // Cleanup
      window.URL.revokeObjectURL(downloadUrl);
    } catch (error) {
      console.error('Error downloading the file:', error);
    }
  };

  const csvStepper = useStepper({
    steps: [
      {
        component: (
          <>
            <StyledHeadingBox>
              <Typography fontWeight={'bold'}>Please upload the model CSV file</Typography>
              <TooltipIconButton
                title="Download Model CSV template"
                onClick={() => handleDownload('Models.csv')}
              >
                <DownloadIcon fill={theme.palette.icon.default} />
              </TooltipIconButton>
            </StyledHeadingBox>
            <FormControl fullWidth>
              {modelCsvFile ? (
                <StyledUploadSuccess>
                  <Typography variant="subtitle1">File uploaded successfully</Typography>
                  <StyledFileChip
                    icon={<InsertDriveFileIcon />}
                    label={modelCsvFile.name || 'model.csv'}
                    onDelete={() => setModelCsvFile(null)}
                  />
                </StyledUploadSuccess>
              ) : (
                <>
                  <OutlinedInput
                    required
                    id="model-csv-file"
                    type="file"
                    inputProps={{
                      accept: '.csv',
                    }}
                    onChange={handleModelCsvFileChange}
                    sx={{
                      marginTop: '1rem',
                    }}
                  />
                </>
              )}
            </FormControl>
          </>
        ),
        icon: ModelIcon,
        label: 'Model CSV',
        helpText: (
          <>
            Learn more about Meshery Models in our{' '}
            <StyledDocsRedirectLink
              href={`${MESHERY_DOCS_URL}/concepts/logical/models`}
              target="_blank"
              rel="noopener noreferrer"
            >
              documentation
            </StyledDocsRedirectLink>
            . Models are versioned packages containing components, relationships and policies for
            defining infrastructure in Meshery.
          </>
        ),
      },
      {
        component: (
          <>
            <StyledHeadingBox>
              <Typography fontWeight={'bold'}>Please upload the component CSV file</Typography>
              <TooltipIconButton
                title="Download Component CSV template"
                onClick={() => handleDownload('Components.csv')}
              >
                <DownloadIcon fill={theme.palette.icon.default} />
              </TooltipIconButton>
            </StyledHeadingBox>
            <FormControl fullWidth>
              {componentCsvFile ? (
                <StyledUploadSuccess>
                  <Typography variant="subtitle1">File uploaded successfully</Typography>
                  <StyledFileChip
                    icon={<InsertDriveFileIcon />}
                    label={componentCsvFile.name || 'component.csv'}
                    onDelete={() => setComponentCsvFile(null)}
                  />
                </StyledUploadSuccess>
              ) : (
                <>
                  <OutlinedInput
                    required
                    id="component-csv-file"
                    type="file"
                    inputProps={{
                      accept: '.csv',
                    }}
                    onChange={handleComponentCsvFileChange}
                    sx={{
                      marginTop: '1rem',
                    }}
                  />
                </>
              )}
            </FormControl>
          </>
        ),
        icon: ComponentIcon,
        label: 'Component CSV',
        helpText: (
          <>
            Learn more about Meshery Components in our{' '}
            <StyledDocsRedirectLink
              href={`${MESHERY_DOCS_URL}/concepts/logical/components`}
              target="_blank"
              rel="noopener noreferrer"
            >
              documentation
            </StyledDocsRedirectLink>
            . Components are fundamental building blocks that represent distinct capabilities and
            features of your infrastructure in Meshery.
          </>
        ),
      },
      {
        component: (
          <div>
            <StyledHeadingBox>
              <Typography fontWeight={'bold'}>Please upload the relationship CSV file</Typography>
              <TooltipIconButton
                title="Download Relationship CSV template"
                onClick={() => handleDownload('Relationships.csv')}
              >
                <DownloadIcon fill={theme.palette.icon.default} />
              </TooltipIconButton>
            </StyledHeadingBox>
            <FormControl fullWidth>
              {relationshipCsvFile ? (
                <StyledUploadSuccess>
                  <Typography variant="subtitle2">File uploaded successfully</Typography>

                  <StyledFileChip
                    icon={<InsertDriveFileIcon />}
                    label={relationshipCsvFile.name || 'relationship.csv'}
                    onDelete={() => setRelationshipCsvFile(null)}
                  />
                </StyledUploadSuccess>
              ) : (
                <>
                  <OutlinedInput
                    required
                    id="relationship-csv-file"
                    type="file"
                    inputProps={{
                      accept: '.csv',
                    }}
                    onChange={handleRelationshipCsvFileChange}
                    sx={{
                      marginTop: '1rem',
                    }}
                  />
                </>
              )}
            </FormControl>
          </div>
        ),
        icon: LanOutlinedIcon,
        label: 'Relationship CSV',
        helpText: (
          <>
            Learn more about Meshery Relationship in our{' '}
            <StyledDocsRedirectLink
              href={`${MESHERY_DOCS_URL}/concepts/logical/relationships`}
              target="_blank"
              rel="noopener noreferrer"
            >
              documentation
            </StyledDocsRedirectLink>
            . Relationships define how components interact and connect with each other within your
            infrastructure model in Meshery.
          </>
        ),
      },
      // {
      //   component: (
      //     <div>
      //       <FormControl component="fieldset" marginTop={'1rem'}>
      //         <FormControlLabel
      //           labelPlacement="start"
      //           style={{ marginLeft: '0' }}
      //           control={
      //             <Checkbox
      //               checked={registerModel}
      //               onChange={(e) => setRegisterModel(e.target.checked)}
      //               name="registerModel"
      //               color="primary"
      //             />
      //           }
      //           label="Would you like to register the model now so you can use it immediately after it's generated?"
      //         />
      //       </FormControl>
      //     </div>
      //   ),
      //   icon: AppRegistrationIcon,
      //   label: 'Register Model',
      //   helpText: 'Choose whether to register the model.',
      // },
    ],
  });

  const transitionConfig = {
    0: {
      canGoNext: () => modelCsvFile !== null,
      nextButtonText: 'Next',
      nextAction: () => csvStepper.handleNext(),
    },
    1: {
      canGoNext: () => componentCsvFile !== null,
      nextButtonText: 'Next',
      nextAction: () => csvStepper.handleNext(),
    },
    2: {
      canGoNext: () => true,
      nextButtonText: 'Finish',
      nextAction: handleFinish,
    },
  };

  const canGoNext = transitionConfig[csvStepper.activeStep].canGoNext();
  const nextButtonText = transitionConfig[csvStepper.activeStep].nextButtonText;

  return (
    <>
      <ModalBody>
        <CustomizedStepper {...csvStepper}>{csvStepper.activeStepComponent}</CustomizedStepper>
      </ModalBody>
      <ModalFooter
        variant="filled"
        helpText={csvStepper.steps[csvStepper.activeStep]?.helpText || ''}
      >
        <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'flex-end' }}>
          <ModalButtonSecondary onClick={csvStepper.goBack} disabled={!csvStepper.canGoBack}>
            Back
          </ModalButtonSecondary>
          <ModalButtonPrimary
            disabled={!canGoNext}
            onClick={transitionConfig[csvStepper.activeStep].nextAction}
          >
            {nextButtonText}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </>
  );
});

CsvStepper.displayName = 'CsvStepper';

export default CsvStepper;
