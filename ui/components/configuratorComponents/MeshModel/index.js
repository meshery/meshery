import { CircularProgress, FormControl, Grid, IconButton, MenuItem, NoSsr, TextField, Toolbar, Tooltip } from "@material-ui/core";
import React, { useRef, useState } from "react";
import AppBarComponent from "./styledComponents/AppBar";

import DeleteIcon from "@material-ui/icons/Delete";
import FileCopyIcon from '@material-ui/icons/FileCopy';
import SaveIcon from '@material-ui/icons/Save';
import { iconMedium } from "../../../css/icons.styles";
import { useMeshModelComponents } from "../../../utils/hooks/useMeshModelComponents";
import { randomPatternNameGenerator as getRandomName } from "../../../utils/utils";
import LazyComponentForm from "./LazyComponentForm";
import useDesignLifecycle from "./hooks/useDesignLifecycle";
import CodeEditor from "../CodeEditor";


export default function DesignConfigurator() {
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [selectedModel, setSelectedModel] = useState(null);
  const { models, meshmodelComponents, getModelFromCategory, getComponentsFromModel, categories } = useMeshModelComponents();
  const { onSettingsChange, onDelete, onSubmit, designYaml } = useDesignLifecycle();
  const formReference = useRef();

  function handleCategoryChange(event) {
    setSelectedCategory(event.target.value);
    getModelFromCategory(event.target.value);
  }

  function handleModelChange(event) {
    if (event.target.value) {
      getComponentsFromModel(event.target.value);
      setSelectedModel(event.target.value);
    }
  }

  return (
    <NoSsr>
      <AppBarComponent position="static" elevation={0}>
        <Toolbar>
          {/* Category Selector */}
          <FormControl>
            <TextField
              select={true}
              SelectProps={{
                MenuProps: {
                  anchorOrigin: {
                    vertical: "bottom",
                    horizontal: "left"
                  },
                  getContentAnchorEl: null
                }
              }}
              InputProps={{ disableUnderline: true }}
              labelId="category-selector"
              id="category-selector"
              value={selectedCategory}
              onChange={handleCategoryChange}
              fullWidth
            >
              {categories.map(cat => (<MenuItem key={cat.name} value={cat.name}>
                {cat.name}
              </MenuItem>))}
            </TextField>
          </FormControl>

          {/* Model Selector */}
          <FormControl>
            <TextField
              select={true}
              SelectProps={{
                MenuProps: {
                  anchorOrigin: {
                    vertical: "bottom",
                    horizontal: "left"
                  },
                  getContentAnchorEl: null
                }
              }}
              InputProps={{ disableUnderline: true }}
              labelId="model-selector"
              id="model-selector"
              value={selectedModel}
              onChange={handleModelChange}
              fullWidth
            >
              {models?.[selectedCategory]
                ? models[selectedCategory].map(function renderModels(model, idx) {
                  return (<MenuItem key={`${model.name}-${idx}`} value={model.name} >{model.displayName}</MenuItem>)
                })
                : <RenderModelNull selectedCategory={selectedCategory} models={models} />
              }
            </TextField>
          </FormControl>

          <Tooltip title="Save Pattern as New File">
            <IconButton
              aria-label="Save"
              color="primary"
              onClick={() => handleSubmitFinalPattern(yaml, "", getRandomName(), "upload")}
            >
              <FileCopyIcon style={iconMedium} />
            </IconButton>
          </Tooltip>
          <Tooltip title="Update Pattern">
            <IconButton
              aria-label="Update"
              color="primary"
              onClick={() => handleSubmitFinalPattern(yaml, pattern.id, pattern.name, "update")}
            >
              <SaveIcon style={iconMedium} />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Pattern">
            <IconButton
              aria-label="Delete"
              color="primary"
              // onClick={() => handleSubmitFinalPattern(yaml, pattern.id, pattern.name, "delete")}
              onClick={() => {
                deleteConfiguration()
              }}
            >
              <DeleteIcon style={iconMedium} />
            </IconButton>
          </Tooltip>
        </Toolbar>
      </AppBarComponent>
      <Grid container spacing={3}>
        {meshmodelComponents?.[selectedModel] && <Grid item xs={12} md={6} >
          {
            meshmodelComponents[selectedModel]?.[0]?.components?.map(function ShowRjsfComponentsLazily(trimmedComponent, idx) {
              return (
                <LazyComponentForm key={`${trimmedComponent.kind}-${idx}`} component={trimmedComponent} onSettingsChange={onSettingsChange(trimmedComponent, formReference)} reference={formReference} />
              )
            })
          }
        </Grid>}
        <Grid item xs={12} md={selectedCategory && selectedModel ? 6 : 12}>
          <CodeEditor yaml={designYaml} saveCodeEditorChanges={() => { }} />
        </Grid>
      </Grid>
    </NoSsr>
  )
}


function RenderModelNull({ selectedCategory, models }) {
  if (!selectedCategory) {
    return <MenuItem value={undefined}>Select a Category First</MenuItem>
  }

  if (!models?.[selectedCategory]) {
    return <CircularProgress />
  }
}