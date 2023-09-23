import { Avatar, CircularProgress, FormControl, Grid, IconButton, MenuItem, NoSsr, TextField, Toolbar, Tooltip } from "@material-ui/core";
import React, { useEffect, useRef, useState } from "react";
import AppBarComponent from "./styledComponents/AppBar";

import DeleteIcon from "@material-ui/icons/Delete";
import SaveIcon from '@material-ui/icons/Save';
import { AvatarGroup } from "@mui/material";
import { iconMedium } from "../../../css/icons.styles";
import { useMeshModelComponents } from "../../../utils/hooks/useMeshModelComponents";
import { getWebAdress } from "../../../utils/webApis";
import CodeEditor from "../CodeEditor";
import LazyComponentForm from "./LazyComponentForm";
import useDesignLifecycle from "./hooks/useDesignLifecycle";
import { useRouter } from "next/router";
import { ArrowBack  } from "@material-ui/icons";
import TooltipButton from "../../../utils/TooltipButton";
import { SaveAs as SaveAsIcon } from "@mui/icons-material";

export default function DesignConfigurator() {
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [selectedModel, setSelectedModel] = useState(null);
  const { models, meshmodelComponents, getModelFromCategory, getComponentsFromModel, categories } = useMeshModelComponents();
  const { onSettingsChange, designSave, designUpdate, designYaml, designJson, designId, designDelete ,updateDesignName,designName ,loadDesign ,updateDesignData } = useDesignLifecycle();
  const formReference = useRef();

  const router = useRouter();
  const { design_id } = router.query;
  console.log("design_id", design_id)

  useEffect(function loadDesignOnMount() {
    if (design_id) {
      loadDesign(design_id)
    }
  }, [design_id])

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
      <TooltipButton title="Back" placement="left">
        <IconButton onClick={() => router.back()}>
          <ArrowBack />
        </IconButton>
      </TooltipButton>
      <AppBarComponent position="static" elevation={0}>
        <Toolbar>
          <div style={{ flexGrow : 1 }}>
            {/* Category Selector */}
            <FormControl>
              <TextField
                select={true}
                SelectProps={{
                  MenuProps : {
                    anchorOrigin : {
                      vertical : "bottom",
                      horizontal : "left"
                    },
                    getContentAnchorEl : null
                  },
                  renderValue : (selected) => {
                    if (!selected || selected.length === 0) {
                      return <em>Select Category</em>;
                    }

                    return selected
                  },
                  displayEmpty : true
                }}
                InputProps={{ disableUnderline : true }}
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
            {selectedCategory && (
              <FormControl>
                <TextField
                  placeholder="select Model"
                  select={true}
                  SelectProps={{
                    MenuProps : {
                      anchorOrigin : {
                        vertical : "bottom",
                        horizontal : "left"
                      },
                      getContentAnchorEl : null
                    },
                    renderValue : (selected) => {
                      if (!selected || selected.length === 0) {
                        return <em>Select Model</em>;
                      }

                      return removeHyphenAndCapitalise(selected)
                    },
                    displayEmpty : true
                  }}
                  InputProps={{ disableUnderline : true }}
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
            )}
          </div>

          {/* Action Toolbar */}
          <TextField label="Design Name" value={designName} onChange={(e) => updateDesignName(e.target.value) } />


          <Tooltip title="Save Design as New File">
            <IconButton
              aria-label="Save"
              color="primary"
              onClick={designSave}
            >
              <SaveAsIcon style={iconMedium} />
            </IconButton>
          </Tooltip>
          {
            designId && (
              <>
                <Tooltip title="Update Design">
                  <IconButton
                    aria-label="Update"
                    color="primary"
                    onClick={designUpdate}
                  >
                    <SaveIcon style={iconMedium} />
                  </IconButton>
                </Tooltip>
                <Tooltip title="Delete Design">
                  <IconButton
                    aria-label="Delete"
                    color="primary"
                    onClick={designDelete}
                  >
                    <DeleteIcon style={iconMedium} />
                  </IconButton>
                </Tooltip>
              </>
            )
          }


        </Toolbar>
      </AppBarComponent>
      <Grid container spacing={3}>
        {meshmodelComponents?.[selectedModel] && <Grid item xs={12} md={6} >
          {
            meshmodelComponents[selectedModel]?.[0]?.components?.map(function ShowRjsfComponentsLazily(trimmedComponent, idx) {
              const hasInvalidSchema = !!trimmedComponent.metadata?.hasInvalidSchema
              return (
                <LazyComponentForm key={`${trimmedComponent.kind}-${idx}`} component={trimmedComponent} onSettingsChange={onSettingsChange(trimmedComponent, formReference)} reference={formReference} disabled={hasInvalidSchema} />
              )
            })
          }
        </Grid>}
        <Grid item xs={12} md={selectedCategory && selectedModel ? 6 : 12}>
          <CodeEditor yaml={designYaml}
            onChange={(_val,_view,update) => {
              updateDesignData({ yamlData : update })
            }}
            saveCodeEditorChanges={(args) => {
              console.log("onSave",args)
            }} fullWidth={!(selectedCategory && selectedModel)} />
          {
            designJson?.services && Object.keys(designJson.services).length > 0 && (
              <AvatarGroup max={10} style={{
                position : "fixed",
                bottom : 60,
                right : 40,
              }}>
                {
                  Object.values(designJson.services).map(function renderAvatarFromServices(service, idx) {
                    const metadata = service.traits?.["meshmodel-metadata"];
                    if (metadata) {
                      const { primaryColor, svgWhite } = metadata;
                      return <Avatar key={idx} src={`${getWebAdress()}/${svgWhite}`} style={{ background : primaryColor, padding : 6, height : 20, width : 20 }} onClick={() => {
                        console.log("TODO: write function to highlight things on editor")
                      }} />
                    }
                  })
                }
              </AvatarGroup>
            )
          }
        </Grid>
      </Grid>
    </NoSsr >
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

function removeHyphenAndCapitalise(str) {
  if (!str) {
    return ""
  }

  return str.split("-").filter(word => word).map(word => word[0].toUpperCase() + word.substring(1)).join(" ")
}