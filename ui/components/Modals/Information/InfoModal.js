/* eslint-disable react/display-name */
import React, { useEffect, useRef, useState } from 'react';
import CloseIcon from '@material-ui/icons/Close';
import PatternIcon from '../../../assets/icons/Pattern';
import { Typography, IconButton, Button, Grid, Avatar, Tooltip, Box } from '@material-ui/core';
import useStyles, {
  ActionContainer,
  CreatAtContainer,
  CopyLinkButton,
  VisibilityTag,
  ResourceName,
} from './styles';
import { iconMedium, iconSmall } from '../../../css/icons.styles';
import { getDesignVersion, getSharableCommonHostAndprotocolLink } from '../../../utils/utils';
import moment from 'moment';
import Application from '../../../public/static/img/drawer-icons/application_svg.js';
import { useSnackbar } from 'notistack';
import Filter from '../../../public/static/img/drawer-icons/filter_svg.js';
import { PATTERN_ENDPOINT, FILTER_ENDPOINT } from '../../../constants/endpoints';
import { useNotification } from '../../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../../lib/event-types';
import axios from 'axios';
import _ from 'lodash';
import RJSFWrapper from '../../MesheryMeshInterface/PatternService/RJSF_wrapper';
import CircularProgress from '@mui/material/CircularProgress';
import { Provider } from 'react-redux';
import { store } from '../../../store';
import { useGetUserByIdQuery } from '../../../rtk-query/user.js';
import { ErrorBoundary } from '../../General/ErrorBoundary';
import { getUnit8ArrayForDesign } from '@/utils/utils';
import ServiceMesheryIcon from '@/assets/icons/ServiceMesheryIcon';
import {
  Modal,
  ModalFooter,
  ModalButtonPrimary,
  ModalButtonSecondary,
  ModalBody,
} from '@layer5/sistent';
import TooltipButton from '@/utils/TooltipButton';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import { filterEmptyFields } from '@/utils/objects';
import { VIEW_VISIBILITY, VisibilityMenu } from '@/components/VisibilityMenu';
import { Lock, Public } from '@material-ui/icons';

const APPLICATION_PLURAL = 'applications';
const FILTER_PLURAL = 'filters';
const PATTERN_PLURAL = 'patterns';

const InfoModal_ = React.memo((props) => {
  const {
    infoModalOpen,
    handleInfoModalClose,
    resourceOwnerID,
    dataName,
    selectedResource,
    currentUser,
    patternFetcher,
    formSchema,
    handlePublish,
    meshModels = [],
    isReadOnly = false,
  } = props;

  const formRef = React.createRef();
  const formStateRef = useRef();
  const [isCatalogDataEqual, setIsCatalogDataEqual] = useState(false);
  const [dataIsUpdated, setDataIsUpdated] = useState(false);
  const [visibility, setVisibility] = useState(selectedResource?.visibility);
  const [saveFormLoading, setSaveFormLoading] = useState(false);
  const [uiSchema, setUiSchema] = useState({});
  const { notify } = useNotification();
  const classes = useStyles();
  const formatDate = (date) => {
    return moment(date).utc().format('MMMM Do YYYY');
  };
  const currentUserID = currentUser?.id;
  const isAdmin = currentUser?.role_names?.includes('admin') || false;
  const { data: resourceUserProfile } = useGetUserByIdQuery(resourceOwnerID);
  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const isOwner = currentUserID === resourceOwnerID;
  const handleCopy = () => {
    navigator.clipboard.writeText(getSharableCommonHostAndprotocolLink(selectedResource));
    enqueueSnackbar(`Link to "${selectedResource.name}" is copied to clipboard`, {
      variant: 'info',
      action: (key) => (
        <IconButton
          // key={`closeevent-${id}`}
          aria-label="Close"
          color="inherit"
          onClick={() => closeSnackbar(key)}
        >
          <CloseIcon style={iconMedium} />
        </IconButton>
      ),
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
    if (dataName === PATTERN_PLURAL) {
      body = JSON.stringify({
        pattern_data: {
          catalog_data: modifiedData,
          pattern_file: getUnit8ArrayForDesign(selectedResource.pattern_file),
          id: selectedResource.id,
          visibility: visibility,
        },
        save: true,
      });
    } else if (dataName === FILTER_PLURAL) {
      setSaveFormLoading(true);
      let config = '';
      if (selectedResource?.filter_resource && selectedResource?.filter_resource !== '') {
        config = JSON.parse(selectedResource.filter_resource).settings.config; // send the config in order to prevent over-write of the config
      }
      body = JSON.stringify({
        filter_data: {
          catalog_data: modifiedData,
          id: selectedResource.id,
          name: selectedResource.name,
          filter_file: selectedResource.filter_file,
        },
        config: config,
        save: true,
      });
    }

    axios
      .post(dataName === PATTERN_PLURAL ? PATTERN_ENDPOINT : FILTER_ENDPOINT, body)
      .then(() => {
        setSaveFormLoading(false);
        notify({
          message: `${selectedResource.name} data saved successfully`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        patternFetcher()();
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
        formStateRef.current.compatibility = filteredCompatibilityArray;
      }
    } else {
      formStateRef.current = filterEmptyFields(selectedResource?.catalog_data);
    }
  }, [selectedResource?.catalog_data, meshModels]);

  useEffect(() => {
    if (formSchema) {
      const newUiSchema = { ...formSchema.uiSchema };

      // if (!isAdmin) {
      //   selectedResource?.visibility === 'private' ? null : (
      //     currentUserID !== resourceOwnerID &&
      //     (newUiSchema['ui:readonly'] = currentUserID !== resourceOwnerID));
      // }

      if (isReadOnly) {
        newUiSchema['ui:readonly'] = true;
      }

      setUiSchema(newUiSchema);
    }
  }, [resourceOwnerID, formSchema, currentUserID]);

  const renderIcon = () => {
    if (dataName === PATTERN_PLURAL) {
      return <PatternIcon style={{ ...iconSmall }} fill="#FFF" />;
    }
    if (dataName === APPLICATION_PLURAL) {
      return <Application style={{ ...iconSmall }} fill="#FFF" />;
    }
    if (dataName === FILTER_PLURAL) {
      return <Filter style={{ ...iconSmall }} fill="#FFF" />;
    }
  };

  const shouldRenderSaveButton = () => {
    if (!isAdmin) {
      const isPrivate = selectedResource?.visibility === 'private';

      const renderByPermission = isPrivate ? true : isPublished ? false : isOwner;
      return renderByPermission;
    }
    return true;
  };

  const handlePublishController = () => {
    formRef.current.state.schema['required'] = formSchema.rjsfSchema.required;
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

  const handleError = () => {
    setImageError(true);
  };
  return (
    <div style={{ marginBottom: '1rem' }}>
      <Modal
        open={infoModalOpen}
        closeModal={handleInfoModalClose}
        title={selectedResource?.name}
        headerIcon={renderIcon()}
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
              <ResourceName className={classes.resourceName} variant="subtitle1">
                {selectedResource?.name}
              </ResourceName>
              <Grid item xs={12} style={{ marginTop: '1rem' }}>
                <Typography
                  style={{ whiteSpace: 'nowrap' }}
                  gutterBottom
                  variant="subtitle1"
                  className={classes.text}
                >
                  <CreatAtContainer isBold={true}>Created</CreatAtContainer>
                  <CreatAtContainer isBold={false}>
                    {formatDate(selectedResource?.created_at)}
                  </CreatAtContainer>
                </Typography>
              </Grid>
              <Grid item xs={12}>
                <Typography
                  style={{ whiteSpace: 'nowrap' }}
                  gutterBottom
                  variant="subtitle1"
                  className={classes.text}
                >
                  <CreatAtContainer isBold={true}>Updated</CreatAtContainer>
                  <CreatAtContainer idBold={false}>
                    {formatDate(selectedResource?.updated_at)}
                  </CreatAtContainer>
                </Typography>
              </Grid>
              {version === '' ? null : (
                <Grid item xs={12}>
                  <Typography
                    style={{ whiteSpace: 'nowrap' }}
                    gutterBottom
                    variant="subtitle1"
                    className={classes.text}
                  >
                    <CreatAtContainer isBold={true}>Version</CreatAtContainer>
                    <CreatAtContainer isBold={false}>{version}</CreatAtContainer>
                  </Typography>
                </Grid>
              )}
            </Grid>
            <Grid item xs={8} lg>
              <Grid container spacing={2}>
                <Grid item xs={dataName === APPLICATION_PLURAL ? 12 : 6}>
                  <Typography gutterBottom variant="subtitle1">
                    <Tooltip
                      title={`Owner: ${
                        resourceUserProfile?.first_name + ' ' + resourceUserProfile?.last_name
                      }`}
                    >
                      <OwnerChip userProfile={resourceUserProfile} />
                    </Tooltip>
                  </Typography>
                </Grid>
                <Grid
                  item
                  xs={dataName === APPLICATION_PLURAL ? 12 : 6}
                  className={classes.visibilityGridItem}
                >
                  <Typography
                    gutterBottom
                    variant="subtitle1"
                    className={classes.text}
                    style={{ display: 'flex', marginRight: '2rem' }}
                  >
                    {!isPublished && isOwner ? (
                      <VisibilityMenu
                        value={visibility}
                        onChange={(value) => {
                          setVisibility(value);
                          setDataIsUpdated(value != selectedResource?.visibility);
                        }}
                        enabled={true}
                        options={[
                          [VIEW_VISIBILITY.PUBLIC, Public],
                          [VIEW_VISIBILITY.PRIVATE, Lock],
                        ]}
                      />
                    ) : (
                      <VisibilityTag>{selectedResource?.visibility}</VisibilityTag>
                    )}
                  </Typography>
                </Grid>
                {dataName === APPLICATION_PLURAL && formSchema ? null : (
                  <Grid item className={classes.rjsfInfoModalForm}>
                    <RJSFWrapper
                      formData={formStateRef.current}
                      jsonSchema={{
                        ...formSchema.rjsfSchema,
                        required: [],
                      }}
                      uiSchema={uiSchema}
                      onChange={handleFormChange}
                      liveValidate={false}
                      formRef={formRef}
                      hideTitle={true}
                    />
                  </Grid>
                )}
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
              className={classes.copyButton}
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
              className={classes.submitButton}
              onClick={handleSubmit}
              disabled={
                (isCatalogDataEqual && !dataIsUpdated) ||
                !shouldRenderSaveButton() ||
                saveFormLoading
              }
            >
              {saveFormLoading ? (
                <Box sx={{ display: 'flex' }}>
                  <CircularProgress color="inherit" size="1.4rem" />
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

const OwnerChip = ({ userProfile }) => {
  const classes = useStyles();
  return (
    <Box style={{ display: 'flex', gap: 10, alignItems: 'center' }}>
      <Avatar src={userProfile?.avatar_url} className={classes.chipIcon} />
      <Typography>
        {userProfile ? (
          `${userProfile?.first_name} ${userProfile?.last_name}`
        ) : (
          <Box sx={{ display: 'flex' }}>
            <CircularProgress color="inherit" size="1rem" />
          </Box>
        )}
      </Typography>
    </Box>
  );
};

const InfoModal = (props) => {
  return (
    <ErrorBoundary
      FallbackComponent={() => null}
      onError={(e) => console.error('Error in Info modal', e)}
    >
      <Provider store={store}>
        <InfoModal_ {...props} />
      </Provider>
    </ErrorBoundary>
  );
};

export default InfoModal;
