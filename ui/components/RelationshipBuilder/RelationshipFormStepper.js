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
  MenuItem,
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
import omit from 'lodash/omit';
import { useMeshModelComponents } from '@/utils/hooks/useMeshModelComponents';

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

  const filteredSchema = omit(RelationshipDefinitionV1Alpha3Schema, ['properties.capabilities']);
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

  const [selectedCategory, setSelectedCategory] = React.useState(null);
  const [selectedModel, setSelectedModel] = React.useState(null);
  const { models, meshmodelComponents, getModelFromCategory, getComponentsFromModel, categories } =
    useMeshModelComponents();

  const handleFormChange = (formData) => {
    if (selectedModel && selectedCategory) {
      const modelData = models[selectedCategory]?.find((model) => model.name === selectedModel);

      if (modelData) {
        formData = {
          ...formData,
          model: {
            ...formData.model,
            id: modelData.id,
            name: modelData.name,
            schemaVersion: modelData.schemaVersion,
            subCategory: modelData.subCategory,
            status: modelData.status,
            version: modelData.version,
            displayName: modelData.displayName,
            registrant: modelData.registrant,
            metadata: modelData.metadata,
          },
        };
      }
    }

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

  const handleCategoryChange = (event) => {
    setSelectedCategory(event.target.value);
    getModelFromCategory(event.target.value);
    setSelectedModel(null);
  };

  const handleModelChange = (event) => {
    if (event.target.value) {
      getComponentsFromModel(event.target.value);
      setSelectedModel(event.target.value);

      const modelData = models[selectedCategory]?.find(
        (model) => model.name === event.target.value,
      );
      if (modelData) {
        setFormData((prevData) => ({
          ...prevData,
          model: {
            ...formData.model,
            id: modelData.id,
            name: modelData.name,
            schemaVersion: modelData.schemaVersion,
            subCategory: modelData.subCategory,
            status: modelData.status,
            version: modelData.version,
            displayName: modelData.displayName,
            registrant: modelData.registrant,
            metadata: modelData.metadata,
          },
        }));
      }
    }
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

              <Grid2 size={{ xs: 12 }}>
                <Typography variant="subtitle1" mt={2} mb={1}>
                  Select Model Category and Type
                </Typography>
              </Grid2>

              <Grid2 size={{ xs: 12, sm: 6 }}>
                <FormControl fullWidth>
                  <TextField
                    select={true}
                    SelectProps={{
                      MenuProps: {
                        anchorOrigin: {
                          vertical: 'bottom',
                          horizontal: 'left',
                        },
                        getContentAnchorEl: null,
                      },
                      renderValue: (selected) => {
                        if (!selected || selected.length === 0) {
                          return <em>Select Category</em>;
                        }
                        return selected;
                      },
                      displayEmpty: true,
                    }}
                    id="category-selector"
                    value={selectedCategory || ''}
                    onChange={handleCategoryChange}
                    fullWidth
                    variant="outlined"
                    helperText="Select a model category"
                  >
                    <MenuItem value="" disabled>
                      <em>Select Category</em>
                    </MenuItem>
                    {categories.map((cat) => (
                      <MenuItem key={cat.name} value={cat.name}>
                        {cat.name}
                      </MenuItem>
                    ))}
                  </TextField>
                </FormControl>
              </Grid2>

              <Grid2 size={{ xs: 12, sm: 6 }}>
                <FormControl fullWidth>
                  <TextField
                    select={true}
                    disabled={!selectedCategory}
                    SelectProps={{
                      MenuProps: {
                        anchorOrigin: {
                          vertical: 'bottom',
                          horizontal: 'left',
                        },
                        getContentAnchorEl: null,
                      },
                      renderValue: (selected) => {
                        if (!selected || selected.length === 0) {
                          return <em>Select Model</em>;
                        }
                        return selected;
                      },
                      displayEmpty: true,
                    }}
                    id="model-selector"
                    value={selectedModel || ''}
                    onChange={handleModelChange}
                    fullWidth
                    variant="outlined"
                    helperText="Select a model type"
                  >
                    <MenuItem value="" disabled>
                      <em>Select Model</em>
                    </MenuItem>
                    {models?.[selectedCategory]?.map((model, idx) => (
                      <MenuItem key={`${model.name}-${idx}`} value={model.name}>
                        {model.displayName}
                      </MenuItem>
                    ))}
                  </TextField>
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

            {selectedModel && meshmodelComponents[selectedModel] && (
              <Box mb={3} mt={2}>
                <Typography variant="subtitle1" mb={1}>
                  Selected Model: <strong>{selectedModel}</strong>
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  The form below will be pre-populated with some values from your selected model.
                </Typography>
              </Box>
            )}

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
                <strong>Download:</strong> Click &quot;Download JSON&quot; to save the relationship
                file on your computer
              </li>
            </ul>
          </>
        ),
      },
    ],
  });

  const transitionConfig = {
    0: {
      canGoNext: () =>
        relationshipName &&
        /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/.test(relationshipName) &&
        selectedCategory &&
        selectedModel,
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
