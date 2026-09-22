import {
  Avatar,
  AvatarGroup,
  FormControl,
  Grid2,
  IconButton,
  InputAdornment,
  Menu,
  MenuItem,
  TextField,
  Toolbar,
  CustomTooltip,
  styled,
  useTheme,
} from '@sistent/sistent';
import React, { useEffect, useRef, useState } from 'react';
import AppBarComponent from './styledComponents/AppBar';
import {
  Delete as DeleteIcon,
  Save as SaveIcon,
  ArrowBack,
  ExpandMoreIcon,
  SaveAs as SaveAsIcon,
} from '@/assets/icons';
import { NoSsr } from '@sistent/sistent';
import { iconMedium } from '../../../../css/icons.styles';
import { useMeshModelComponents } from '../../../../utils/hooks/useMeshModelComponents';
import { getWebAdress } from '../../../../utils/webApis';
import CodeEditor from '../CodeEditor';
import LazyComponentForm from './LazyComponentForm';
import useDesignLifecycle from './hooks/useDesignLifecycle';
import { useRouter } from 'next/router';

import { Keys } from '@meshery/schemas/permissions';

const ScrollContainer = styled('div')({
  overflowY: 'auto',
  width: '100%',
  height: '58.5vh',
});

export default function DesignConfigurator() {
  const [selectedCategory, setSelectedCategory] = useState(null);
  const [selectedModel, setSelectedModel] = useState(null);
  const [browseCategory, setBrowseCategory] = useState(null);
  const [menuAnchor, setMenuAnchor] = useState(null);
  const { models, meshmodelComponents, getModelFromCategory, getComponentsFromModel, categories } =
    useMeshModelComponents();
  const {
    onSettingsChange,
    designSave,
    designUpdate,
    designYaml,
    designJson,
    designId,
    designDelete,
    updateDesignName,
    loadDesign,
    updateDesignData,
  } = useDesignLifecycle();
  const formReference = useRef();

  const router = useRouter();
  const { design_id } = router.query;

  const theme = useTheme();

  useEffect(
    function loadDesignOnMount() {
      if (design_id) {
        loadDesign(design_id);
      }
    },
    [design_id],
  );

  function handleDropdownOpen(event) {
    setBrowseCategory(null);
    setMenuAnchor(event.currentTarget);
  }

  function handleDropdownClose() {
    setMenuAnchor(null);
    setBrowseCategory(null);
  }

  function handleDeselect() {
    setSelectedCategory(null);
    setSelectedModel(null);
    handleDropdownClose();
  }

  function handleModelPick(category, model) {
    setSelectedCategory(category);
    setSelectedModel(model);
    getComponentsFromModel(model);
    handleDropdownClose();
  }

  const selectedModelData = models?.[selectedCategory]?.find((m) => m.name === selectedModel);
  const displayText =
    menuAnchor && browseCategory
      ? browseCategory
      : selectedModelData?.displayName || selectedModel || '';

  return (
    <NoSsr>
      <CustomTooltip title="Back" placement="right">
        <IconButton onClick={() => router.back()}>
          <ArrowBack />
        </IconButton>
      </CustomTooltip>
      <AppBarComponent position="static" elevation={0} data-testid="design-configurator-app-bar">
        <Toolbar>
          <div style={{ flexGrow: 1 }}>
            <FormControl>
              <TextField
                variant="standard"
                id="category-model-selector"
                data-testid="category-model-selector"
                value={displayText}
                placeholder="Select Category & Model"
                onClick={handleDropdownOpen}
                onKeyDown={(event) => {
                  if (event.key === 'Enter') {
                    handleDropdownOpen(event);
                  }
                }}
                slotProps={{
                  input: {
                    readOnly: true,
                    endAdornment: (
                      <InputAdornment position="end">
                        <ExpandMoreIcon />
                      </InputAdornment>
                    ),
                  },
                }}
                fullWidth
              />
              <Menu
                anchorEl={menuAnchor}
                open={Boolean(menuAnchor)}
                onClose={handleDropdownClose}
                anchorOrigin={{
                  vertical: 'bottom',
                  horizontal: 'left',
                }}
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'left',
                }}
                slotProps={{
                  paper: {
                    style: {
                      maxHeight: 400,
                      overflowY: 'auto',
                    },
                  },
                }}
              >
                {browseCategory === null
                  ? [
                      <MenuItem
                        key="none"
                        onClick={handleDeselect}
                        data-testid="clear-category-model-selector"
                      >
                        <em>None</em>
                      </MenuItem>,
                      ...(categories ?? []).map((cat) => (
                        <MenuItem
                          key={cat.name}
                          onClick={() => {
                            getModelFromCategory(cat.name);
                            setBrowseCategory(cat.name);
                          }}
                          data-testid={cat.name}
                        >
                          {cat.name}
                        </MenuItem>
                      )),
                    ]
                  : [
                      <MenuItem
                        key="back"
                        onClick={() => setBrowseCategory(null)}
                        data-testid="back-to-categories"
                      >
                        ← {browseCategory}
                      </MenuItem>,
                      ...(models?.[browseCategory]?.length
                        ? models[browseCategory].map((model, idx) => (
                            <MenuItem
                              key={`${browseCategory}-${model.name}-${idx}`}
                              onClick={() => handleModelPick(browseCategory, model.name)}
                              data-testid={`${browseCategory}-${model.name}`}
                            >
                              {model.displayName}
                            </MenuItem>
                          ))
                        : models?.[browseCategory]
                          ? [
                              <MenuItem key="empty" disabled>
                                No models found
                              </MenuItem>,
                            ]
                          : [
                              <MenuItem key="loading" disabled>
                                Loading…
                              </MenuItem>,
                            ]),
                    ]}
              </Menu>
            </FormControl>
          </div>

          {/* Action Toolbar */}
          <TextField
            label="Design Name"
            value={designJson.name}
            onChange={(e) => updateDesignName(e.target.value)}
            variant="standard"
          />

          <CustomTooltip title="Save Design as New File">
            <div>
              <IconButton
                aria-label="Save"
                data-testid="design-configurator-save-design-btn"
                onClick={designSave}
                permissionKey={Keys.CatalogManagementCreateNewDesign}
              >
                <SaveAsIcon fill={theme.palette.icon.default} style={iconMedium} />
              </IconButton>
            </div>
          </CustomTooltip>
          {designId && (
            <>
              <CustomTooltip title="Update Design">
                <div>
                  <IconButton
                    aria-label="Update"
                    data-testid="design-configurator-update-design-btn"
                    onClick={designUpdate}
                    permissionKey={Keys.CatalogManagementEditDesign}
                  >
                    <SaveIcon style={iconMedium} />
                  </IconButton>
                </div>
              </CustomTooltip>
              <CustomTooltip title="Delete Design">
                <div>
                  <IconButton
                    aria-label="Delete"
                    data-testid="design-configurator-delete-design-btn"
                    onClick={designDelete}
                    permissionKey={Keys.CatalogManagementDeleteADesign}
                  >
                    <DeleteIcon style={iconMedium} />
                  </IconButton>
                </div>
              </CustomTooltip>
            </>
          )}
        </Toolbar>
      </AppBarComponent>
      <Grid2 container spacing={3} size="grow">
        {meshmodelComponents?.[selectedModel] && (
          <Grid2
            size={{ xs: 12, md: 6 }}
            data-testid="model-component-list"
            sx={{
              height: '100%',
              display: 'flex',
            }}
          >
            <ScrollContainer>
              {meshmodelComponents[selectedModel]?.[0]?.components?.map(
                function ShowRjsfComponentsLazily(trimmedComponent, idx) {
                  const hasInvalidSchema = !!trimmedComponent.metadata?.hasInvalidSchema;
                  return (
                    <LazyComponentForm
                      key={`${trimmedComponent.component.kind}-${idx}`}
                      component={trimmedComponent}
                      onSettingsChange={onSettingsChange(trimmedComponent, formReference)}
                      reference={formReference}
                      disabled={hasInvalidSchema}
                    />
                  );
                },
              )}
            </ScrollContainer>
          </Grid2>
        )}
        <Grid2
          data-testid="design-configurator-code-editor"
          size={{ xs: 12, md: selectedCategory && selectedModel ? 6 : 12 }}
        >
          <CodeEditor
            yaml={designYaml}
            onChange={(_val, _view, update) => {
              updateDesignData({ yamlData: update });
            }}
            saveCodeEditorChanges={(args) => {
              console.log('onSave', args);
            }}
            fullWidth={!(selectedCategory && selectedModel)}
          />
          {designJson?.services && Object.keys(designJson.services).length > 0 && (
            <AvatarGroup
              max={10}
              style={{
                position: 'fixed',
                bottom: 60,
                right: 40,
              }}
            >
              {Object.values(designJson.services).map(
                function renderAvatarFromServices(service, idx) {
                  const metadata = service.traits?.['meshmodel-metadata'];
                  if (metadata) {
                    const { primaryColor, svgWhite } = metadata;
                    return (
                      <Avatar
                        key={idx}
                        src={`${getWebAdress()}/${svgWhite}`}
                        style={{ background: primaryColor, padding: 6, height: 20, width: 20 }}
                        alt={service.name}
                        title={service.name}
                        data-testid={'service-avatar-' + idx}
                      />
                    );
                  }
                },
              )}
            </AvatarGroup>
          )}
        </Grid2>
      </Grid2>
    </NoSsr>
  );
}
