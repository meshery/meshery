import React, { useEffect, useState } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import InfoIcon from '@material-ui/icons/Info';
import CloseIcon from '@material-ui/icons/Close';
import PatternIcon from '../../styles/assets/Pattern';
import { Chip, Typography, IconButton, Button, Grid, Avatar, Tooltip } from '@material-ui/core';
import useStyles from '../../Share/styles';
import {
  getSharableCommonHostAndprotocolLink,
  iconSmall,
} from '../../../../sections/MesheryDesignerComponent/PatternServiceFormWrapper/utils';
import OriginalApplicationFileIcon from '../../styles/assets/OriginalApplicationIcon';
import moment from 'moment';
import { APPLICATION_PLURAL, FILTER_PLURAL, PATTERN_PLURAL } from '@/utils/constants';
import ApplicationIcon from '../../styles/assets/ApplicationIcon';
import { uiTheme } from '../../../../globals/theme';
import { useSnackbar } from 'notistack';
import FilterIcon from '../../styles/assets/FilterIcon';
import { RJSF_SCHEMAS, PATTERN_ENDPOINT, FILTER_ENDPOINT } from '@/utils/constants';
import { useNotification, EVENT_TYPES } from '../../../../globals/notifications';
import { getUserProfile } from '@/components/api';
import { getMeshModels } from '@/components/api';
import axios from 'axios';
import _ from 'lodash';
import RJSFWrapper from '../../MesheryMeshInterface/PatternService/RJSF_wrapper';

const InfoModal = (props) => {
  const {
    infoModalOpen,
    handleInfoModalClose,
    designOwnerName,
    ownerAvatar,
    dataName,
    selectedDesign: selectedResource,
  } = props;

  const theme = uiTheme.theme;
  const [formSchema, setFormSchema] = useState({});
  const formRef = React.createRef();
  const [formState, setFormState] = useState(selectedResource?.catalog_data || {});
  const [isCatalogDataEqual, setIsCatalogDataEqual] = useState(false);
  const [userProfile, setUserProfile] = useState({});
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
      let body = null;
      if (dataName === PATTERN_PLURAL) {
        body = {
          body: JSON.stringify({
            pattern_data: {
              catalog_data: formState,
              id: selectedResource.id,
            },
          }),
        };
      } else if (dataName === FILTER_PLURAL) {
        body = {
          body: JSON.stringify({
            filter_data: {
              catalog_data: formState,
              id: selectedResource.id,
            },
          }),
        };
      }

      axios
        .post(dataName === PATTERN_PLURAL ? PATTERN_ENDPOINT : FILTER_ENDPOINT, body, {
          headers: {
            'Content-Type': 'application/json',
          },
        })
        .then(() => {
          notify({
            message: `${selectedResource.name} data saved successfully`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        })
        .catch((error) => {
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
    const newUiSchema = { ...formSchema.uiSchema };

    newUiSchema['ui:readonly'] = userProfile?.id !== selectedResource?.user_id;

    setUiSchema(newUiSchema);
  }, [userProfile, selectedResource]);

  const modifyRJSFSchema = (schema, propertyPath, newValue) => {
    const clonedSchema = _.cloneDeep(schema);
    _.set(clonedSchema, propertyPath, newValue);
    return clonedSchema;
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(`${RJSF_SCHEMAS}/publish`);
        const reponseData = await getMeshModels();
        const userProfile = await getUserProfile();
        const modelNames = _.uniq(reponseData.data.models?.map((model) => model.displayName));

        const modifiedSchema = modifyRJSFSchema(
          response.data.rjsfSchema,
          'properties.compatibility.items.enum',
          modelNames,
        );

        setFormSchema({ rjsfSchema: modifiedSchema, uiSchema: response.data.uiSchema });
        setUserProfile(userProfile.data);
      } catch (error) {
        console.error('Error fetching schema:', error);
      }
    };

    fetchData();
  }, []);

  const renderIcon = () => {
    if (dataName === PATTERN_PLURAL) {
      return <PatternIcon style={{ ...iconSmall }} fill="#FFF" />;
    }
    if (dataName === APPLICATION_PLURAL) {
      return <ApplicationIcon style={{ ...iconSmall }} fill="#FFF" />;
    }
    if (dataName === FILTER_PLURAL) {
      return <FilterIcon style={{ ...iconSmall }} fill="#FFF" />;
    }
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
                    <Tooltip title={`Owner: ${designOwnerName}`}>
                      <Chip
                        avatar={<Avatar src={ownerAvatar} className={classes.chipIcon} />}
                        label={designOwnerName}
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
                {dataName === APPLICATION_PLURAL ? null : (
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
          {userProfile?.id === selectedResource?.user_id ? (
            <Button
              variant="contained"
              color="primary"
              className={classes.submitButton}
              onClick={handleSubmit}
              disabled={isCatalogDataEqual}
            >
              Save
            </Button>
          ) : null}
          {dataName === PATTERN_PLURAL && (
            <Chip
              style={{
                backgroundColor: `${theme.palette.secondary.mainBackground2}`,
                border: '0',
                display: 'flex',
                justifyContent: 'center',
              }}
              size="small"
              avatar={<InfoIcon />}
              label={
                selectedResource?.visibility === 'published'
                  ? 'Clone the design to edit'
                  : 'Public designs cannot be cloned'
              }
            />
          )}
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default InfoModal;