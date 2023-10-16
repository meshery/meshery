import React, { useEffect, useState } from 'react';
import CloseIcon from '@material-ui/icons/Close';
import PatternIcon from '../../../assets/icons/Pattern';
import {
  Chip,
  Typography,
  IconButton,
  Button,
  Grid,
  Avatar,
  Tooltip,
  Box,
  DialogTitle,
  DialogContent,
  DialogActions,
  Dialog,
} from '@layer5/sistent-components';
import useStyles from './styles';
import { iconSmall } from '../../../css/icons.styles';
import { getSharableCommonHostAndprotocolLink } from '../../../utils/utils';
import OriginalApplicationFileIcon from '../../../assets/icons/OriginalApplicationIcon';
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
import { promisifiedDataFetch } from '../../../lib/data-fetch';
import CircularProgress from '@material-ui/core/CircularProgress';

const APPLICATION_PLURAL = 'applications';
const FILTER_PLURAL = 'filters';
const PATTERN_PLURAL = 'patterns';

const InfoModal = (props) => {
  const {
    infoModalOpen,
    handleInfoModalClose,
    resourceOwnerID,
    dataName,
    selectedResource,
    currentUserID,
    formSchema,
  } = props;

  const formRef = React.createRef();
  const [formState, setFormState] = useState(selectedResource?.catalog_data || {});
  const [isCatalogDataEqual, setIsCatalogDataEqual] = useState(false);
  const [resourceUserProfile, setResourceUserProfile] = useState(null);
  const [saveFormLoading, setSaveFormLoading] = useState(false);
  const [uiSchema, setUiSchema] = useState({});
  const { notify } = useNotification();

  const classes = useStyles();
  const formatDate = (date) => {
    return moment(date).utc().format('MMMM Do YYYY, h:mm:ss A');
  };
  const { enqueueSnackbar } = useSnackbar();

  const handleCopy = () => {
    navigator.clipboard.writeText(getSharableCommonHostAndprotocolLink(selectedResource));
    enqueueSnackbar(`Link to "${selectedResource.name}" is copied to clipboard`, {
      variant: 'info',
      autoHideDuration: 2000,
    });
  };

  const handleSubmit = () => {
    if (formRef.current && formRef.current.validateForm()) {
      setSaveFormLoading(true);
      let body = null;
      if (dataName === PATTERN_PLURAL) {
        body = JSON.stringify({
          pattern_data: {
            catalog_data: formState,
            pattern_file: selectedResource.pattern_file,
            id: selectedResource.id,
          },
          save: true,
        });
      } else if (dataName === FILTER_PLURAL) {
        setSaveFormLoading(true);
        let config = '';
        if (selectedResource.filter_resource !== null || selectedResource.filter_resource !== '') {
          config = JSON.parse(selectedResource.filter_resource).settings.config; // send the config inorder to prevent over-write of the config
        }
        body = JSON.stringify({
          filter_data: {
            catalog_data: formState,
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
    }
  };

  const handleFormChange = (data) => {
    setFormState(data);
    const objectsEqual = _.isEqual(selectedResource?.catalog_data, data);
    setIsCatalogDataEqual(objectsEqual);
  };

  useEffect(() => {
    if (selectedResource?.catalog_data) {
      setFormState(selectedResource.catalog_data);
    }
  }, [selectedResource?.catalog_data]);

  useEffect(() => {
    if (formSchema) {
      const newUiSchema = { ...formSchema.uiSchema };

      // Only make form readonly if resource is private and user is not owner
      if (selectedResource?.visibility === 'private' && currentUserID !== resourceOwnerID) {
        newUiSchema['ui:readonly'] = currentUserID !== resourceOwnerID;
      }

      setUiSchema(newUiSchema);
    }
  }, [resourceOwnerID, formSchema, currentUserID]);

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      try {
        const response = await promisifiedDataFetch(`/api/user/profile/${resourceOwnerID}`, {
          method: 'GET',
          credentials: 'include',
        });

        // Check if the component is still mounted before setting the state
        if (isMounted) {
          setResourceUserProfile(response);
        }
      } catch (error) {
        console.error('Error fetching user profile with id in info modal:', error);
      }
    };

    fetchData();

    // Cleanup function
    return () => {
      isMounted = false;
    };
  }, [resourceOwnerID]);

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
    const isPrivate = selectedResource?.visibility === 'private';
    const isOwner = currentUserID === resourceOwnerID;

    return (isPrivate && !isOwner) || !isPrivate;
  };

  return (
    <div style={{ marginBottom: '1rem' }}>
      <Dialog
        open={infoModalOpen}
        onClose={handleInfoModalClose}
        aria-labelledby="form-dialog-title"
        maxWidth={dataName === APPLICATION_PLURAL ? 'xs' : 'md'}
        style={{ zIndex: 9999 }}
      >
        <DialogTitle textAlign="center" id="form-dialog-title" className={classes.dialogTitle}>
          {renderIcon()}
          <Typography className={classes.textHeader} variant="h6">
            {selectedResource?.name}
          </Typography>
          <IconButton
            aria-label="close"
            onClick={handleInfoModalClose}
            component="button"
            style={{
              color: '#FFFFFF',
            }}
          >
            <CloseIcon className={classes.closing} />
          </IconButton>
        </DialogTitle>
        <DialogContent style={{ padding: '1.5rem' }}>
          <Grid container spacing={2}>
            <Grid item>
              <Button
                variant="outlined"
                disabled
                startIcon={
                  <OriginalApplicationFileIcon
                    style={{
                      boxShadow: '0px 0px 6px 2px rgba(0, 0, 0, 0.25)',
                      borderRadius: '20px',
                    }}
                    width={150}
                    height={150}
                  />
                }
              ></Button>
              <Typography className={classes.resourceName} variant="subtitle1">
                {selectedResource?.name}
              </Typography>
            </Grid>
            <Grid item xs={8} lg>
              <Grid container spacing={2}>
                <Grid item xs={dataName === APPLICATION_PLURAL ? 12 : 6}>
                  <Typography gutterBottom variant="subtitle1" className={classes.text}>
                    <span style={{ fontWeight: 'bold', paddingRight: '0.2rem' }}>Owner:</span>
                    <Tooltip
                      title={`Owner: ${
                        resourceUserProfile?.first_name + ' ' + resourceUserProfile?.last_name
                      }`}
                    >
                      <Chip
                        avatar={
                          <Avatar
                            src={resourceUserProfile?.avatar_url}
                            className={classes.chipIcon}
                          />
                        }
                        label={
                          resourceUserProfile ? (
                            resourceUserProfile?.first_name + ' ' + resourceUserProfile?.last_name
                          ) : (
                            <Box sx={{ display: 'flex' }}>
                              <CircularProgress color="inherit" size="1rem" />
                            </Box>
                          )
                        }
                        variant="outlined"
                        data-cy="chipDesignDetails"
                        className={classes.chip}
                      />
                    </Tooltip>
                  </Typography>
                </Grid>
                <Grid item xs={dataName === APPLICATION_PLURAL ? 12 : 6}>
                  <Typography gutterBottom variant="subtitle1" className={classes.text}>
                    <span
                      style={{
                        fontWeight: 'bold',
                        fontFamily: 'Qanelas Soft, sans-serif',
                      }}
                    >
                      Visibility:
                    </span>{' '}
                    <Chip
                      label={selectedResource?.visibility}
                      variant="outlined"
                      className={classes.chip}
                    />
                  </Typography>
                </Grid>
                {dataName === APPLICATION_PLURAL && formSchema ? null : (
                  <Grid item className={classes.rjsfInfoModalForm}>
                    <RJSFWrapper
                      formData={formState}
                      jsonSchema={formSchema.rjsfSchema}
                      uiSchema={uiSchema}
                      onChange={handleFormChange}
                      liveValidate={false}
                      formRef={formRef}
                      hideTitle={true}
                    />
                  </Grid>
                )}
                <Grid item xs={12}>
                  <Typography
                    style={{ whiteSpace: 'nowrap' }}
                    gutterBottom
                    variant="subtitle1"
                    className={classes.text}
                  >
                    <span style={{ fontWeight: 'bold' }}>Created At:</span>{' '}
                    <span style={{ whiteSpace: 'wrap' }}>
                      {formatDate(selectedResource?.created_at)}
                    </span>
                  </Typography>
                </Grid>
                <Grid item xs={12}>
                  <Typography
                    style={{ whiteSpace: 'nowrap' }}
                    gutterBottom
                    variant="subtitle1"
                    className={classes.text}
                  >
                    <span style={{ fontWeight: 'bold' }}>Updated At:</span>{' '}
                    <span style={{ whiteSpace: 'nowrap' }}>
                      {formatDate(selectedResource?.updated_at)}
                    </span>
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions
          style={{
            justifyContent: 'space-evenly',
            marginBottom: '0.5rem',
          }}
        >
          <Button variant="outlined" onClick={handleCopy} className={classes.copyButton}>
            Copy Link
          </Button>
          {shouldRenderSaveButton() ? (
            <Button
              variant="contained"
              color="primary"
              className={classes.submitButton}
              onClick={handleSubmit}
              disabled={isCatalogDataEqual || saveFormLoading}
            >
              {saveFormLoading ? (
                <Box sx={{ display: 'flex' }}>
                  <CircularProgress color="inherit" size="1.4rem" />
                </Box>
              ) : (
                'Save'
              )}
            </Button>
          ) : null}
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default InfoModal;
