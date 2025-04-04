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
import ModelIcon from '@/assets/icons/ModelIcon';
import LanOutlinedIcon from '@mui/icons-material/LanOutlined';
import { TooltipIconButton } from '@/utils/TooltipButton';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';
import {
  ModelDefinitionV1Beta1Schema,
  ComponentDefinitionV1Beta1Schema,
  RelationshipDefinitionV1Alpha3Schema,
} from '@layer5/schemas';
import FinishFlagIcon from '@/assets/icons/FinishFlagIcon';
import FinishModelGenerateStep from './FinishModelGenerateStep';

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

const CSV_TEMPLATE_BASE_URL =
  'https://raw.githubusercontent.com/meshery/meshery/a514f8689260791077bde8171646933cff15dd08/mesheryctl/templates/template-csvs/';

const CsvStepper = React.memo(({ handleClose }) => {
  const [modelData, setModelData] = React.useState({});
  const [modelCsvFile, setModelCsvFile] = React.useState(null);
  const [componentCsvFile, setComponentCsvFile] = React.useState(null);
  const [relationshipCsvFile, setRelationshipCsvFile] = React.useState(null);
  const registerModel = true;

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
            . {ModelDefinitionV1Beta1Schema.description}
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
            . {ComponentDefinitionV1Beta1Schema.description}
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
            . {RelationshipDefinitionV1Alpha3Schema.description}
          </>
        ),
      },
      {
        component: (
          <FinishModelGenerateStep
            requestBody={{
              importBody: {
                model_csv: modelData.model_csv,
                component_csv: modelData.component_csv,
                relationship_csv: modelData.relationship_csv,
              },
              uploadType: 'csv',
              register: modelData.register,
            }}
            generateType="register"
          />
        ),
        label: 'Finish',
        icon: FinishFlagIcon,
      },
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
      canGoNext: () => relationshipCsvFile !== null,
      nextButtonText: 'Generate',
      nextAction: () => {
        csvStepper.handleNext();
        setModelData({
          model_csv: `data:text/csv;base64,${modelCsvFile?.base64?.split(',')[1]}`,
          component_csv: `data:text/csv;base64,${componentCsvFile?.base64?.split(',')[1]}`,
          relationship_csv: relationshipCsvFile
            ? `data:text/csv;base64,${relationshipCsvFile?.base64?.split(',')[1]}`
            : null,
          register: registerModel,
        });
      },
    },
    3: {
      canGoNext: () => true,
      nextButtonText: 'Finish',
      nextAction: handleClose,
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
