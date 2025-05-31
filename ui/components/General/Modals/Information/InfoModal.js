import ServiceMesheryIcon from '@/assets/icons/ServiceMesheryIcon';
import { usePublishPatternMutation, useUpdatePatternFileMutation } from '@/rtk-query/design';
import TooltipButton from '@/utils/TooltipButton';
import CAN from '@/utils/can';
import { filterEmptyFields } from '@/utils/objects';
import { keys } from '@/utils/permission_constants';
import {
  Avatar,
  Box,
  Button,
  CircularProgress,
  CustomTooltip,
  getFormatDate,
  Grid,
  IconButton,
  Link,
  Modal,
  ModalBody,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalFooter,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
  Skeleton,
  Typography,
  VisibilityChipMenu,
} from '@layer5/sistent';
import CloseIcon from '@mui/icons-material/Close';
import yaml from 'js-yaml';
import _ from 'lodash';
import { useSnackbar } from 'notistack';
import React, { useEffect, useRef, useState } from 'react';
import PatternIcon from '../../../../assets/icons/Pattern';
import { MESHERY_CLOUD_PROD } from '../../../../constants/endpoints';
import { iconMedium, iconSmall } from '../../../../css/icons.styles';
import { EVENT_TYPES } from '../../../../lib/event-types';
import { useGetUserByIdQuery } from '../../../../rtk-query/user.js';
import { useNotification } from '../../../../utils/hooks/useNotification';
import {
  getDesignVersion,
  getSharableCommonHostAndprotocolLink,
  modifyRJSFSchema,
} from '../../../../utils/utils';
import { ActionContainer, CopyLinkButton, CreatAtContainer, ResourceName } from './styles';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { updateProgress } from '@/store/slices/mesheryUi';
import { getMeshModels } from '@/api/meshmodel';
import { useSelector } from 'react-redux';
import { Lock, Public } from '@mui/icons-material';
import RJSFWrapper from '@/components/MesheryMeshInterface/PatternService/RJSF_wrapper';

export const VIEW_VISIBILITY = {
  PUBLIC: 'public',
  PRIVATE: 'private',
};

const InfoModal_ = React.memo((props) => {
  const {
    infoModalOpen,
    handleInfoModalClose,
    resourceOwnerID,
    selectedResource,
    patternFetcher,
    isReadOnly = false,
  } = props;

  const { user: currentUser } = useSelector((state) => state.ui);
  const formRef = React.createRef();
  const formStateRef = useRef();
  const [isCatalogDataEqual, setIsCatalogDataEqual] = useState(false);
  const [dataIsUpdated, setDataIsUpdated] = useState(false);
  const [visibility, setVisibility] = useState(selectedResource?.visibility);
  const [saveFormLoading, setSaveFormLoading] = useState(false);
  const [uiSchema, setUiSchema] = useState({});
  const { notify } = useNotification();

  const [updatePattern] = useUpdatePatternFileMutation();
  const currentUserID = currentUser?.id;
  const isAdmin = currentUser?.role_names?.includes('admin') || false;
  const { data: resourceUserProfile } = useGetUserByIdQuery(resourceOwnerID);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const isOwner = currentUserID === resourceOwnerID;
  const [meshModels, setMeshModels] = useState([]);
  const [publishSchema, setPublishSchema] = useState({});

  useEffect(() => {
    const fetchModels = async () => {
      const { models } = await getMeshModels();
      const modelNames = _.uniqBy(
        models?.map((model) => {
          if (model.displayName && model.displayName !== '') {
            return model.displayName;
          }
        }),
        _.toLower,
      );
      modelNames.sort();

      // Modify the schema using the utility function
      const modifiedSchema = modifyRJSFSchema(
        publishCatalogItemSchema,
        'properties.compatibility.items.enum',
        modelNames,
      );
      setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: publishCatalogItemUiSchema });
      setMeshModels(models);
    };
    fetchModels();
  }, []);

  const handleCopy = () => {
    navigator.clipboard.writeText(getSharableCommonHostAndprotocolLink(selectedResource));
    enqueueSnackbar(`Link to "${selectedResource.name}" is copied to clipboard`, {
      variant: 'info',
      action: (key) => (
        <IconButton aria-label="Close" color="inherit" onClick={() => closeSnackbar(key)}>
          <CloseIcon style={iconMedium} />
        </IconButton>
      ),
    });
  };

  const [publishCatalog] = usePublishPatternMutation();

  const handlePublish = (formData) => {
    const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formData?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: selectedResource?.id,
      catalog_data: {
        ...formData,
        compatibility: compatibilityStore,
        type: _.toLower(formData?.type),
      },
    };
    updateProgress({ showProgress: true });
    publishCatalog({
      publishBody: JSON.stringify(payload),
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        if (currentUser.role_names.includes('admin')) {
          notify({
            message: `${selectedResource.name} Design Published`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        } else {
          notify({
            message:
              'Design queued for publishing into Meshery Catalog. Maintainers notified for review',
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `Error while publishing ${selectedResource.name} Design`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };

  const handleSubmit = () => {
    setSaveFormLoading(true);
    const compatibilityStore = _.uniqBy(meshModels, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formStateRef.current?.compatibility?.some(
          (comp) => _.toLower(comp) === _.toLower(model.displayName),
        ),
      )
      ?.map((model) => model.name);

    let body = null;
    let modifiedData = {
      ...formStateRef.current,
      type: formStateRef.current?.type?.toLowerCase(),
      compatibility: compatibilityStore,
    };

    body = JSON.stringify({
      name: selectedResource.name,
      catalog_data: modifiedData,
      design_file: yaml.load(selectedResource.pattern_file),
      id: selectedResource.id,
      visibility: visibility,
    });

    updatePattern({ updateBody: body })
      .then(() => {
        setSaveFormLoading(false);
        notify({
          message: `${selectedResource.name} data saved`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        patternFetcher && patternFetcher();
        handleInfoModalClose();
      })
      .catch((error) => {
        setSaveFormLoading(false);
        const errorMessage = error.response?.data?.message || error.message;

        notify({
          message: `Error while saving ${selectedResource.name} data: ${errorMessage}`,
          event_type: EVENT_TYPES.ERROR,
        });

        console.error('Error while saving pattern data', error);
      });
  };
  function normalizeCompatibility(arr) {
    return (
      arr?.map((item) => {
        if (typeof item === 'string') {
          return item.toLowerCase().split('-').join(' ');
        }
        return item;
      }) || []
    );
  }
  // Function to compare objects while normalizing case in compatibility array
  function isEqualIgnoringCase(obj1, obj2) {
    // Check each property one by one
    const isEqualPatternCaveats = obj1.pattern_caveats === obj2.pattern_caveats;
    const isEqualPatternInfo = obj1.pattern_info === obj2.pattern_info;
    const isEqualType = obj1.type?.toLowerCase() === obj2.type?.toLowerCase();

    // Normalize and compare compatibility array
    const normalizedCompat1 = normalizeCompatibility(obj1.compatibility);
    const normalizedCompat2 = normalizeCompatibility(obj2.compatibility);
    const isEqualCompatibility = _.isEqual(normalizedCompat1, normalizedCompat2);

    // Return true only if all properties are equal
    return isEqualPatternCaveats && isEqualPatternInfo && isEqualType && isEqualCompatibility;
  }
  const handleFormChange = (data) => {
    formStateRef.current = data;
    setIsCatalogDataEqual(isEqualIgnoringCase(selectedResource?.catalog_data, data));
  };

  useEffect(() => {
    if (selectedResource?.catalog_data && Object.keys(selectedResource?.catalog_data).length > 0) {
      if (meshModels) {
        const compatibilitySet = new Set(
          (selectedResource?.catalog_data?.compatibility || []).map((comp) => comp.toLowerCase()),
        );

        const filteredCompatibilityArray = _.uniq(
          meshModels
            .filter((obj) => {
              const modelName = obj.name.toLowerCase();
              return compatibilitySet.has(modelName);
            })
            .map((obj) => obj.displayName),
        );

        let modifiedData = {
          ...selectedResource.catalog_data,
          type: _.startCase(selectedResource?.catalog_data?.type),
          compatibility: filteredCompatibilityArray,
        };
        formStateRef.current = filterEmptyFields(modifiedData);
      }
    } else {
      formStateRef.current = filterEmptyFields(selectedResource?.catalog_data);
    }
  }, [selectedResource?.catalog_data, meshModels]);

  useEffect(() => {
    if (publishSchema) {
      const newUiSchema = { ...publishSchema.uiSchema };

      if (isReadOnly) {
        newUiSchema['ui:readonly'] = true;
      }

      setUiSchema(newUiSchema);
    }
  }, [resourceOwnerID, publishSchema, currentUserID]);

  const shouldRenderSaveButton = () => {
    if (!isAdmin) {
      const isPrivate = selectedResource?.visibility === 'private';

      const renderByPermission = isPrivate ? true : isPublished ? false : isOwner;
      return renderByPermission;
    }
    return true;
  };

  const handlePublishController = () => {
    formRef.current.state.schema['required'] = publishSchema.rjsfSchema.required;
    if (formRef.current && formRef.current.validateForm()) {
      setSaveFormLoading(true);
      handleInfoModalClose();
      handlePublish(formRef.current.state.formData);
      setSaveFormLoading(false);
    }
  };

  const isPublished = selectedResource?.visibility === 'published';
  const [imageError, setImageError] = useState(false);
  const version = getDesignVersion(selectedResource);
  const canChangeVisibility = !isPublished && isOwner;
  const handleError = () => {
    setImageError(true);
  };

  return (
    <div style={{ marginBottom: '1rem' }}>
      <Modal
        open={infoModalOpen}
        closeModal={handleInfoModalClose}
        title={selectedResource?.name}
        headerIcon={<PatternIcon style={{ ...iconSmall }} fill="#FFF" />}
        maxWidth={false}
        sx={{
          '& .MuiDialog-container': {
            '& .MuiPaper-root': {
              width: '100%',
              maxWidth: '800px',
            },
          },
        }}
      >
        <ModalBody>
          <Grid container spacing={2}>
            <Grid item>
              <Button
                variant="outlined"
                disabled
                style={{
                  border: '0.1px solid #E6E6E6',
                  width: '150px',
                  height: '150px',
                  display: 'flex',
                  justifyContent: 'center',
                  alignItems: 'center',
                }}
              >
                {selectedResource?.catalog_data?.imageURL && !imageError ? (
                  <img
                    src={selectedResource.catalog_data.imageURL[0]}
                    width="140"
                    alt={selectedResource?.name}
                    onError={handleError}
                  />
                ) : (
                  <ServiceMesheryIcon
                    style={{
                      boxShadow: '0px 0px 6px 2px rgba(0, 0, 0, 0.25)',
                      borderRadius: '20px',
                    }}
                    width={100}
                    height={100}
                  />
                )}
              </Button>
              <ResourceName variant="subtitle1">{selectedResource?.name}</ResourceName>
              <Grid item xs={12} style={{ marginTop: '1rem' }}>
                <Typography style={{ whiteSpace: 'nowrap' }} gutterBottom variant="subtitle1">
                  <CreatAtContainer isBold={true}>Created</CreatAtContainer>
                  <CreatAtContainer isBold={false}>
                    {getFormatDate(selectedResource?.created_at)}
                  </CreatAtContainer>
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <Typography style={{ whiteSpace: 'nowrap' }} gutterBottom variant="subtitle1">
                  <CreatAtContainer isBold={true}>Updated</CreatAtContainer>
                  <CreatAtContainer idBold={false}>
                    {getFormatDate(selectedResource?.updated_at)}
                  </CreatAtContainer>
                </Typography>
              </Grid>
              {version === '' ? null : (
                <Grid item xs={12}>
                  <Typography style={{ whiteSpace: 'nowrap' }} gutterBottom variant="subtitle1">
                    <CreatAtContainer isBold={true}>Version</CreatAtContainer>
                    <CreatAtContainer isBold={false}>{version}</CreatAtContainer>
                  </Typography>
                </Grid>
              )}
            </Grid>
            <Grid item xs={8} lg>
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography gutterBottom variant="subtitle1">
                    <CustomTooltip
                      title={`Owner: ${
                        resourceUserProfile?.first_name + ' ' + resourceUserProfile?.last_name
                      }`}
                    >
                      <div>
                        <OwnerChip userProfile={resourceUserProfile} />
                      </div>
                    </CustomTooltip>
                  </Typography>
                </Grid>
                <Grid
                  item
                  xs={6}
                  style={{ display: 'flex', alignItems: 'center', justifyContent: 'flex-end' }}
                >
                  <Typography
                    gutterBottom
                    variant="subtitle1"
                    style={{ display: 'flex', marginRight: '2rem' }}
                  >
                    <VisibilityChipMenu
                      value={visibility}
                      onChange={(value) => {
                        setVisibility(value);
                        setDataIsUpdated(value != selectedResource?.visibility);
                      }}
                      enabled={canChangeVisibility}
                      options={[
                        [VIEW_VISIBILITY.PUBLIC, Public],
                        [VIEW_VISIBILITY.PRIVATE, Lock],
                      ]}
                    />
                  </Typography>
                </Grid>

                <Grid
                  item
                  style={{
                    marginLeft: '-1rem',
                    marginTop: '-1rem',
                    maxWidth: '39rem',
                  }}
                >
                  <RJSFWrapper
                    formData={formStateRef.current}
                    jsonSchema={{
                      ...publishSchema.rjsfSchema,
                      required: [],
                    }}
                    uiSchema={uiSchema}
                    onChange={handleFormChange}
                    liveValidate={false}
                    formRef={formRef}
                    hideTitle={true}
                  />
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </ModalBody>
        <ModalFooter
          helpText={
            'Upon submitting your catalog item, an approval flow will be initiated. [Learn More](https://docs.meshery.io/concepts/catalog)'
          }
          variant="filled"
        >
          <ActionContainer>
            <TooltipButton title={'Copy Design Link'} onClick={handleCopy}>
              <CopyLinkButton>Copy Link</CopyLinkButton>
            </TooltipButton>

            <ModalButtonSecondary
              data-testid="publish-button"
              variant="outlined"
              onClick={handlePublishController}
              disabled={
                !isPublished
                  ? false
                  : !(
                      CAN(keys.PUBLISH_DESIGN.action, keys.PUBLISH_DESIGN.subject) &&
                      currentUser?.id === selectedResource?.user_id
                    ) || isPublished
              }
            >
              {isPublished ? 'Published' : 'Publish to Catalog'}
            </ModalButtonSecondary>
            <ModalButtonPrimary
              variant="contained"
              color="primary"
              onClick={handleSubmit}
              disabled={
                (isCatalogDataEqual && !dataIsUpdated) ||
                !shouldRenderSaveButton() ||
                saveFormLoading
              }
            >
              {saveFormLoading ? (
                <Box sx={{ display: 'flex' }}>
                  <CircularProgress size="1.4rem" />
                </Box>
              ) : (
                'Save'
              )}
            </ModalButtonPrimary>
          </ActionContainer>
        </ModalFooter>
      </Modal>
    </div>
  );
});

InfoModal_.displayName = 'InfoModal_';

const OwnerChip = ({ userProfile }) => {
  return (
    <Box style={{ display: 'flex', gap: 10, alignItems: 'center' }}>
      {userProfile ? (
        <>
          <Link href={`${MESHERY_CLOUD_PROD}/user/${userProfile.id}`} rel="noopener noreferrer">
            <Avatar src={userProfile.avatar_url} />
          </Link>
          <Typography>{`${userProfile.first_name} ${userProfile.last_name}`}</Typography>
        </>
      ) : (
        <Skeleton variant="circular" width={40} height={40} />
      )}
    </Box>
  );
};

const InfoModal = (props) => {
  return (
    <ProviderStoreWrapper>
      <InfoModal_ {...props} />
    </ProviderStoreWrapper>
  );
};

export default InfoModal;
