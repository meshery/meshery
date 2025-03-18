import React from 'react';
import { modelShapes } from './data';
import {
  TextField,
  Grid,
  Select,
  InputLabel,
  DeleteIcon,
  Typography,
  FormControl,
  MenuItem,
  Button,
  Box,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  useTheme,
} from '@layer5/sistent';
import { ExpandMore, Add as AddIcon } from '@mui/icons-material';
import { TooltipIconButton } from '@/utils/TooltipButton';
import ClickAwayListener from '@mui/material/ClickAwayListener';
import { StyledSketchContainer, StyledSketchPicker, StyledSketchWrapper } from './styles';

const ComponentStep = ({
  components,
  setComponents,
  modelCategory,
  modelShape,
  modelSubCategory,
  modelPrimaryColor,
  modelSecondryColor,
}) => {
  const [openColorPicker, setOpenColorPicker] = React.useState({});

  const handleColorPickerToggle = (componentId, type) => {
    setOpenColorPicker((prev) => {
      // Close all other pickers first
      const newState = {};
      // Only set the new picker state
      newState[componentId] = {
        [type]: !prev[componentId]?.[type],
      };
      return newState;
    });
  };

  const handleClickAway = () => {
    setOpenColorPicker({});
  };

  const handleComponentChange = (id, field, value) => {
    setComponents((prevComponents) =>
      prevComponents.map((comp) => (comp.id === id ? { ...comp, [field]: value } : comp)),
    );
  };

  const handleFileChange = (id, field, file) => {
    if (file && file.type === 'image/svg+xml') {
      const reader = new FileReader();
      reader.onload = (e) => {
        const svgData = e.target.result;
        handleComponentChange(id, field, svgData);
      };
      reader.readAsText(file);
    } else {
      console.error('Please upload a valid SVG file.');
    }
  };

  const addNewComponent = () => {
    setComponents((prev) => [
      ...prev,
      {
        id: Date.now(),
        component: '',
        category: modelCategory,
        subCategory: modelSubCategory,
        shape: modelShape,
        primaryColor: modelPrimaryColor,
        secondaryColor: modelSecondryColor,
        svgColor: null,
        svgWhite: null,
        isAnnotation: true,
      },
    ]);
  };

  const removeComponent = (id) => {
    setComponents((prev) => prev.filter((comp) => comp.id !== id));
  };

  const theme = useTheme();

  return components.length === 0 ? (
    <Box
      sx={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        flexDirection: 'column',
        height: '100%',
      }}
    >
      <Button variant="contained" onClick={addNewComponent} sx={{ mt: 2 }}>
        Add Component
      </Button>
    </Box>
  ) : (
    <Box>
      {components.map((component, index) => (
        <Accordion
          style={{ backgroundColor: 'transparent' }}
          key={component.id}
          defaultExpanded={true}
          sx={{ mb: 1 }}
        >
          <AccordionSummary expandIcon={<ExpandMore />}>
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                width: '100%',
              }}
            >
              <Typography
                variant="body1"
                fontWeight={'bold'}
              >{`Component ${index + 1}`}</Typography>
              <Box display={'flex'} gap={0.5} mr={1}>
                <TooltipIconButton
                  title={'Add Another Component'}
                  onClick={(e) => {
                    e.stopPropagation();
                    addNewComponent();
                  }}
                >
                  <AddIcon />
                </TooltipIconButton>
                <TooltipIconButton
                  onClick={(e) => {
                    e.stopPropagation();
                    removeComponent(component.id);
                  }}
                  title={'Remove Component'}
                >
                  <DeleteIcon fill={theme.palette.icon.default} />
                </TooltipIconButton>
              </Box>
            </Box>
          </AccordionSummary>
          <AccordionDetails style={{ paddingBottom: '2rem' }}>
            <Grid container spacing={2}>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <TextField
                    required
                    label="Component Name"
                    placeholder="Enter component name"
                    value={component.component}
                    onChange={(e) =>
                      handleComponentChange(component.id, 'component', e.target.value)
                    }
                    error={
                      component.component.length > 0 &&
                      !/^[a-zA-Z0-9\s]+$/.test(component.component)
                    }
                    variant="outlined"
                  />
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth variant="outlined">
                  <InputLabel>Shape</InputLabel>
                  <Select
                    value={component.shape}
                    label="Shape"
                    onChange={(e) => handleComponentChange(component.id, 'shape', e.target.value)}
                  >
                    {modelShapes.map((shape, idx) => (
                      <MenuItem key={idx} value={shape}>
                        {shape}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <Typography>Primary Color</Typography>
                  <Box position={'relative'}>
                    <StyledSketchContainer
                      primaryColor={component.primaryColor}
                      onClick={() => handleColorPickerToggle(component.id, 'primary')}
                    />
                    {openColorPicker[component.id]?.primary && (
                      <ClickAwayListener onClickAway={handleClickAway}>
                        <StyledSketchWrapper>
                          <StyledSketchPicker
                            style={{
                              backgroundColor: theme.palette.background.surfaces,
                            }}
                            color={component.primaryColor}
                            onChange={(color) =>
                              handleComponentChange(component.id, 'primaryColor', color.hex)
                            }
                          />
                        </StyledSketchWrapper>
                      </ClickAwayListener>
                    )}
                  </Box>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <Typography>Secondary Color</Typography>
                  <Box position={'relative'}>
                    <StyledSketchContainer
                      primaryColor={component.secondaryColor}
                      onClick={() => handleColorPickerToggle(component.id, 'secondary')}
                    />
                    {openColorPicker[component.id]?.secondary && (
                      <ClickAwayListener onClickAway={handleClickAway}>
                        <StyledSketchWrapper>
                          <StyledSketchPicker
                            color={component.secondaryColor}
                            onChange={(color) =>
                              handleComponentChange(component.id, 'secondaryColor', color.hex)
                            }
                          />
                        </StyledSketchWrapper>
                      </ClickAwayListener>
                    )}
                  </Box>
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <Typography>Component Color SVG</Typography>
                  <input
                    type="file"
                    accept=".svg"
                    onChange={(e) => handleFileChange(component.id, 'svgColor', e.target.files[0])}
                    style={{ marginTop: '0.5rem' }}
                  />
                </FormControl>
              </Grid>

              <Grid item xs={12} md={6}>
                <FormControl fullWidth>
                  <Typography>Component White SVG</Typography>
                  <input
                    type="file"
                    accept=".svg"
                    onChange={(e) => handleFileChange(component.id, 'svgWhite', e.target.files[0])}
                    style={{ marginTop: '0.5rem' }}
                  />
                </FormControl>
              </Grid>
            </Grid>
          </AccordionDetails>
        </Accordion>
      ))}
    </Box>
  );
};

export default ComponentStep;
