import React from 'react';
import {
  ModalFooter,
  useStepper,
  CustomizedStepper,
  ModalBody,
  Box,
  TextField,
  ModalButtonSecondary,
  ModalButtonPrimary,
  Typography,
  FormControl,
  Grid2,
} from '@sistent/sistent';
import { GlobalStyles } from '@mui/material';
import { styled } from '@sistent/sistent';
import { RelationshipDefinitionV1Alpha3OpenApiSchema } from '@meshery/schemas';
import DescriptionIcon from '@mui/icons-material/Description';
import CodeIcon from '@mui/icons-material/Code';
import LinkIcon from '@mui/icons-material/Link';
import SaveAltIcon from '@mui/icons-material/SaveAlt';
import { downloadFileFromContent } from '@/utils/fileDownloader';
import RJSFWrapper from '../MesheryMeshInterface/PatternService/RJSF_wrapper';

const StyledSummaryBox = styled(Box)(({ theme }) => ({
  backgroundColor: theme.palette.background.blur?.heavy,
  padding: '1rem',
  borderRadius: '8px',
  marginBottom: '1rem',
}));
const StyledDocsRedirectLink = styled('a')(({ theme }) => ({
  color: theme.palette.background.brand.default,
  textDecoration: 'underline',
}));

const RelationshipFormStepper = React.memo(({ handleClose }) => {
  const RelationshipDefinitionV1Alpha3Schema =
    RelationshipDefinitionV1Alpha3OpenApiSchema.components.schemas.RelationshipDefinition;

  const filteredSchema = _.omit(RelationshipDefinitionV1Alpha3Schema, ['properties.capabilities']);
  // -> added zIndex so dropdown appears above the modal (mui default zIndex:1300)
  const globalStyles = (
    <GlobalStyles
      styles={{
        '.react-select__menu': {
          zIndex: 1500,
        },
        '.MuiAutocomplete-popper': {
          zIndex: 1500,
        },
        '.MuiPopover-root': {
          zIndex: 1500,
        },
        '.MuiMenu-paper': {
          zIndex: 1500,
        },
      }}
    />
  );

  const [relationshipName, setRelationshipName] = React.useState('');
  const [formData, setFormData] = React.useState({});
  const [jsonOutput, setJsonOutput] = React.useState('');
  const formRef = React.useRef();

  const handleFormChange = (formData) => {
    setFormData(formData);
  };

  const handleDownload = () => {
    const jsonContent = JSON.stringify(formData, null, 2);
    downloadFileFromContent(
      jsonContent,
      `${relationshipName || 'relationship'}.json`,
      'application/json',
    );
  };

  const relationshipStepper = useStepper({
    steps: [
      {
        component: (
          <div>
            <Box display="flex" alignItems="center" mb={2}>
              <Typography>
                Please enter a unique <strong>Relationship Name</strong> for your relationship
                definition.
              </Typography>
            </Box>
            <Grid2 container spacing={2} size="grow">
              <Grid2 size={{ xs: 12 }}>
                <FormControl fullWidth>
                  <TextField
                    required
                    id="relationship-name"
                    label="Relationship File Name"
                    placeholder="my-relationship"
                    helperText="A unique name for your relationship's file using lowercase letters and hyphens"
                    error={
                      relationshipName.length > 0 &&
                      !/^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/.test(relationshipName)
                    }
                    value={relationshipName}
                    onChange={(e) => setRelationshipName(e.target.value)}
                    variant="outlined"
                  />
                </FormControl>
              </Grid2>
            </Grid2>
          </div>
        ),
        icon: DescriptionIcon,
        label: 'Basic Details',
        helpText: (
          <>
            <p>
              Learn more about{' '}
              <StyledDocsRedirectLink
                href="https://docs.meshery.io/concepts/logical/relationships"
                target="_blank"
                rel="noopener noreferrer"
              >
                Meshery Relationships
              </StyledDocsRedirectLink>
              .
            </p>
          </>
        ),
      },
      {
        component: (
          <div>
            <Box display="flex" alignItems="center">
              <Typography>Define your relationship properties using the form below.</Typography>
            </Box>
            <Box>
              <RJSFWrapper
                jsonSchema={filteredSchema}
                formData={formData}
                onChange={handleFormChange}
                formRef={formRef}
                className="relationship-form-wrapper"
                liveValidate={false}
                focusOnFirstError={false}
              />
            </Box>
          </div>
        ),
        icon: LinkIcon,
        label: 'Relationship Properties',
        helpText: (
          <>
            <p>Fill out all required fields to properly define your relationship.</p>
          </>
        ),
      },
      {
        component: (
          <div>
            <Box display="flex" alignItems="center" mb={2}>
              <Typography>Review the JSON output for your relationship definition.</Typography>
            </Box>
            <StyledSummaryBox>
              <pre
                style={{
                  whiteSpace: 'pre-wrap',
                  maxHeight: '400px',
                  overflowX: 'auto',
                  padding: '4px',
                  borderRadius: '4px',
                  scrollbarWidth: 'none',
                  '&::-webkit-scrollbar': {
                    display: 'none',
                  },
                }}
              >
                {jsonOutput || 'No JSON preview available. Please go back and try again.'}
              </pre>
            </StyledSummaryBox>
            <Box mt={2} display="flex" justifyContent="center">
              <ModalButtonPrimary onClick={handleDownload} startIcon={<SaveAltIcon />}>
                Download JSON
              </ModalButtonPrimary>
            </Box>
          </div>
        ),
        icon: CodeIcon,
        label: 'JSON Preview',
        helpText: (
          <>
            <ul>
              <li>
                <strong>Preview JSON:</strong> Review the generated JSON for your relationship
                definition.
              </li>
              <li>
                <strong>Download:</strong> Click "Download JSON" to save the relationship file on
                your computer
              </li>
            </ul>
          </>
        ),
      },
    ],
  });

  const transitionConfig = {
    0: {
      canGoNext: () => relationshipName && /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/.test(relationshipName),
      nextButtonText: 'Next',
      nextAction: () => relationshipStepper.handleNext(),
    },
    1: {
      canGoNext: () => true,
      nextButtonText: 'Next',
      nextAction: () => {
        const jsonContent = JSON.stringify(formData, null, 2);
        setJsonOutput(jsonContent);
        relationshipStepper.handleNext();
      },
    },
    2: {
      canGoNext: () => true,
      nextButtonText: 'Finish',
      nextAction: handleClose,
    },
  };

  const canGoNext = transitionConfig[relationshipStepper.activeStep].canGoNext();
  const nextButtonText = transitionConfig[relationshipStepper.activeStep].nextButtonText;

  return (
    <>
      {globalStyles}
      <ModalBody>
        <CustomizedStepper {...relationshipStepper}>
          {relationshipStepper.activeStepComponent}
        </CustomizedStepper>
      </ModalBody>
      <ModalFooter
        variant="filled"
        helpText={relationshipStepper.steps[relationshipStepper.activeStep]?.helpText || ``}
      >
        <Box style={{ width: '100%', display: 'flex', gap: '1rem', justifyContent: 'end' }}>
          <ModalButtonSecondary
            onClick={relationshipStepper.goBack}
            disabled={!relationshipStepper.canGoBack}
          >
            Back
          </ModalButtonSecondary>
          <ModalButtonPrimary
            disabled={!canGoNext}
            onClick={transitionConfig[relationshipStepper.activeStep].nextAction}
            data-testid={`RelationshipStepper-Button-${nextButtonText}`}
          >
            {nextButtonText}
          </ModalButtonPrimary>
        </Box>
      </ModalFooter>
    </>
  );
});

RelationshipFormStepper.displayName = 'RelationshipFormStepper';

export default RelationshipFormStepper;
