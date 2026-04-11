import React from 'react';
import {
  Box,
  TextField,
  ModalButtonSecondary,
  ModalButtonPrimary,
  Typography,
  FormControl,
  Grid2,
  Tabs,
  Tab,
  MenuItem,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@sistent/sistent';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import RJSFWrapper from '../MesheryMeshInterface/PatternService/RJSF_wrapper';
import cloneDeep from 'lodash/cloneDeep';
import { useMeshModelComponents } from '@/utils/hooks/useMeshModelComponents';
import omit from 'lodash/omit';

const SelectorsForm = ({ selectorsSchema, formData, onChange }) => {
  const [tabValue, setTabValue] = React.useState(0);
  const [selectorsData, setSelectorsData] = React.useState(
    formData?.selectors || {
      allow: { from: [], to: [] },
      deny: { from: [], to: [] },
    },
  );
  const matchSchema = {
    type: 'object',
    properties: omit(selectorsSchema.items.properties.allow.properties.from.items.properties, [
      'kind',
      'model',
    ]),
  };

  const [expandedPanels, setExpandedPanels] = React.useState({});
  const [selectedModelInfo, setSelectedModelInfo] = React.useState({});
  const { models, categories, getModelFromCategory } = useMeshModelComponents();
  const [categorySelections, setCategorySelections] = React.useState({});

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handlePanelChange = (panelId) => (event, isExpanded) => {
    setExpandedPanels({
      ...expandedPanels,
      [panelId]: isExpanded,
    });
  };

  const handleCategoryChange = (type, direction, index, event) => {
    const category = event.target.value;
    setCategorySelections((prev) => ({
      ...prev,
      [`${type}.${direction}.${index}`]: {
        ...prev[`${type}.${direction}.${index}`],
        category,
        model: null,
      },
    }));

    getModelFromCategory(category);
  };

  const handleModelChange = (type, direction, index, event) => {
    const modelName = event.target.value;
    const panelKey = `${type}.${direction}.${index}`;
    const category = categorySelections[panelKey]?.category;

    if (category && modelName) {
      const modelData = models[category]?.find((model) => model.name === modelName);

      if (modelData) {
        const updatedSelectors = cloneDeep(selectorsData);
        updatedSelectors[type][direction][index].model = {
          name: modelData.name,
          id: modelData.id,
          registrant: { kind: modelData.registrant.kind },
        };

        setSelectorsData(updatedSelectors);
        onChange({
          ...formData,
          selectors: updatedSelectors,
        });

        setCategorySelections((prev) => ({
          ...prev,
          [panelKey]: {
            ...prev[panelKey],
            model: modelName,
          },
        }));

        setSelectedModelInfo({
          ...selectedModelInfo,
          [panelKey]: modelData,
        });
      }
    }
  };

  const addSelector = (type, direction) => {
    const newSelector = {
      kind: '',
      model: {
        name: '',
        registrant: { kind: '' },
      },
    };

    const updatedSelectors = cloneDeep(selectorsData);
    updatedSelectors[type][direction].push(newSelector);

    setSelectorsData(updatedSelectors);
    onChange({
      ...formData,
      selectors: updatedSelectors,
    });
  };

  const updateSelector = (type, direction, index, selectorData) => {
    const updatedSelectors = cloneDeep(selectorsData);
    updatedSelectors[type][direction][index] = {
      ...updatedSelectors[type][direction][index],
      ...selectorData,
    };

    setSelectorsData(updatedSelectors);
    onChange({
      ...formData,
      selectors: updatedSelectors,
    });
  };

  const removeSelector = (type, direction, index) => {
    const updatedSelectors = cloneDeep(selectorsData);
    updatedSelectors[type][direction].splice(index, 1);

    setSelectorsData(updatedSelectors);
    onChange({
      ...formData,
      selectors: updatedSelectors,
    });
  };

  const renderSelectorForm = (type, direction) => {
    const selectors = selectorsData[type][direction] || [];

    return (
      <Box mt={2}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
          <Typography variant="subtitle1">{direction} selectors</Typography>
          <ModalButtonPrimary
            size="small"
            sx={{ minWidth: '6rem' }}
            onClick={() => addSelector(type, direction)}
          >
            Add {direction}
          </ModalButtonPrimary>
        </Box>

        {selectors.length === 0 ? (
          <Typography color="textSecondary" variant="body2">
            No {direction} selectors added yet. Click the button above to add one.
          </Typography>
        ) : (
          selectors.map((selector, index) => {
            const panelKey = `${type}.${direction}.${index}`;
            const panelExpanded = expandedPanels[panelKey] || false;
            const currentCategory = categorySelections[panelKey]?.category;
            const currentModel = categorySelections[panelKey]?.model;

            return (
              <Accordion
                key={`${type}-${direction}-${index}`}
                expanded={panelExpanded}
                onChange={handlePanelChange(panelKey)}
              >
                <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                  <Box
                    width="100%"
                    display="flex"
                    justifyContent="space-between"
                    alignItems="center"
                  >
                    <Typography variant="subtitle2">
                      Selector {index + 1}: {selector.kind || 'New Selector'}
                    </Typography>
                    <ModalButtonSecondary
                      size="small"
                      color="error"
                      onClick={(e) => {
                        e.stopPropagation();
                        removeSelector(type, direction, index);
                      }}
                    >
                      Remove
                    </ModalButtonSecondary>
                  </Box>
                </AccordionSummary>
                <AccordionDetails>
                  <Grid2 container spacing={2}>
                    <Grid2 size={{ xs: 12 }}>
                      <FormControl fullWidth>
                        <TextField
                          label="Kind"
                          value={selector.kind || ''}
                          onChange={(e) => {
                            const updatedSelector = { ...selector, kind: e.target.value };
                            updateSelector(type, direction, index, updatedSelector);
                          }}
                          fullWidth
                          variant="outlined"
                          helperText="kind of resource (e.g., ClusterRole, Role)"
                        />
                      </FormControl>
                    </Grid2>

                    <Grid2 size={{ xs: 12, sm: 6 }}>
                      <FormControl fullWidth>
                        <TextField
                          select={true}
                          label="Model Category"
                          value={currentCategory || ''}
                          onChange={(e) => handleCategoryChange(type, direction, index, e)}
                          fullWidth
                          variant="outlined"
                          helperText="Select the model category"
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
                          disabled={!currentCategory}
                          label="Model Name"
                          value={currentModel || ''}
                          onChange={(e) => handleModelChange(type, direction, index, e)}
                          fullWidth
                          variant="outlined"
                          helperText="Select the model name"
                        >
                          <MenuItem value="" disabled>
                            <em>Select Model</em>
                          </MenuItem>
                          {models?.[currentCategory]?.map((model, idx) => (
                            <MenuItem key={`${model.name}-${idx}`} value={model.name}>
                              {model.displayName || model.name}
                            </MenuItem>
                          ))}
                        </TextField>
                      </FormControl>
                    </Grid2>

                    {/* -> here we are implmenting the rjsf for match prope */}
                    <Grid2 size={{ xs: 12 }}>
                      <Accordion>
                        <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                          <Typography>Match Settings</Typography>
                        </AccordionSummary>
                        <AccordionDetails>
                          <RJSFWrapper
                            jsonSchema={matchSchema}
                            formData={selector}
                            onChange={(updatedData) => {
                              updateSelector(type, direction, index, updatedData);
                            }}
                            className="relationship-selectors-form"
                            liveValidate={false}
                            focusOnFirstError={false}
                          />
                        </AccordionDetails>
                      </Accordion>
                    </Grid2>
                  </Grid2>
                </AccordionDetails>
              </Accordion>
            );
          })
        )}
      </Box>
    );
  };

  return (
    <div>
      <Box display="flex" alignItems="center" mb={2}>
        <Typography>Define allowed and denied relationships using selectors.</Typography>
      </Box>

      <Tabs
        value={tabValue}
        onChange={handleTabChange}
        indicatorColor="primary"
        textColor="primary"
        variant="fullWidth"
      >
        <Tab label="Allow" />
        <Tab label="Deny" />
      </Tabs>

      {tabValue === 0 && (
        <Box mt={2}>
          {renderSelectorForm('allow', 'from')}
          {renderSelectorForm('allow', 'to')}
        </Box>
      )}

      {tabValue === 1 && (
        <Box mt={2}>
          {renderSelectorForm('deny', 'from')}
          {renderSelectorForm('deny', 'to')}
        </Box>
      )}
    </div>
  );
};

export default SelectorsForm;
