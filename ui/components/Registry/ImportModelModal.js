import {
  FormControlLabel,
  Button,
  FormControl,
  RadioGroup,
  Radio,
  importModelUiSchema,
  importModelSchema,
  Typography,
  ModalFooter,
  Box,
  Modal,
  useTheme,
} from '@layer5/sistent';
import { RJSFModalWrapper } from '../Modal';
import CsvStepper from './Stepper/CSVStepper';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { updateProgress } from 'lib/store';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import React, { useState, useEffect, useContext } from 'react';
import { capitalize } from 'lodash';
import { Loading } from '@/components/DesignLifeCycle/common';
import { NotificationCenterContext } from '../NotificationCenter';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import {
  ModelImportedSection,
  ModelImportMessages,
} from '../NotificationCenter/formatters/model_registration';
import { StyledDocsRedirectLink } from './Stepper/style';

const FinishDeploymentStep = ({ deploymentType, handleClose }) => {
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const [isDeploying, setIsDeploying] = useState(true);
  const [deployEvent, setDeployEvent] = useState();

  useEffect(() => {
    const subscription = operationsCenterActorRef.on(
      OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      (event) => {
        const serverEvent = event.data.event;
        if (serverEvent.action === deploymentType) {
          setIsDeploying(false);
          setDeployEvent(serverEvent);
        }
      },
    );

    return () => subscription.unsubscribe();
  }, []);

  const progressMessage = `${capitalize(deploymentType)}ing model`;

  const theme = useTheme();
  return (
    <>
      <Box
        sx={{
          padding: '1rem',
          overflow: 'hidden',
          backgroundColor: theme.palette.background.surfaces,
        }}
      >
        {isDeploying ? (
          <Box style={{ padding: '1rem' }}>
            <Loading message={progressMessage} />
          </Box>
        ) : (
          <>
            <ModelImportMessages message={deployEvent?.metadata?.ModelImportMessage} />
            <ModelImportedSection modelDetails={deployEvent?.metadata?.ModelDetails} />
          </>
        )}
      </Box>
      <ModalFooter
        variant="filled"
        helpText={
          'Click Finish to complete the model import process. The imported model will be available in your Model Registry.'
        }
      >
        <Button variant="contained" color="primary" onClick={handleClose}>
          Finish
        </Button>
      </ModalFooter>
    </>
  );
};

const ImportModelModal = React.memo(({ isImportModalOpen, setIsImportModalOpen }) => {
  const [importModalDescription, setImportModalDescription] = useState('');
  const [isCsvModalOpen, setIsCsvModalOpen] = useState(false);
  const [importModelReq] = useImportMeshModelMutation();
  const [activeStep, setActiveStep] = useState(0);

  const handleNext = () => {
    if (activeStep === 0) {
      setActiveStep(1);
    }
  };

  const handleClose = () => {
    setIsImportModalOpen(false);
    setIsCsvModalOpen(false);
    setActiveStep(0);
  };

  const handleImportModelSubmit = async (data) => {
    const { uploadType, url, file } = data;
    let requestBody = null;

    const fileElement = document.getElementById('root_file');

    switch (uploadType) {
      case 'File Import': {
        const fileName = fileElement.files[0].name;
        const fileData = getUnit8ArrayDecodedFile(file);
        if (fileData) {
          requestBody = {
            importBody: {
              model_file: fileData,
              url: '',
              filename: fileName,
            },
            uploadType: 'file',
            register: true,
          };
        } else {
          console.error('Error: File data is empty or invalid');
          return;
        }
        break;
      }
      case 'URL Import': {
        if (url) {
          requestBody = {
            importBody: {
              url: url,
            },
            uploadType: 'urlImport',
            register: true,
          };
        } else {
          console.error('Error: URL is empty');
          return;
        }
        break;
      }
      case 'CSV Import': {
        handleClose();
        setIsCsvModalOpen(true);
        return;
      }
      default: {
        console.error('Error: Invalid upload type');
        return;
      }
    }
    updateProgress({ showProgress: true });
    await importModelReq({ importBody: requestBody });
    updateProgress({ showProgress: false });
  };

  const CustomRadioWidget = (props) => {
    const { options, value, onChange, label, schema } = props;
    const { enumOptions } = options;

    setImportModalDescription(schema.description);

    return (
      <FormControl component="fieldset">
        <RadioGroup
          value={value}
          onChange={(e) => onChange(e.target.value)}
          style={{ marginTop: '-1.7rem', marginLeft: '-1rem' }}
        >
          <Typography fontWeight={'bold'} fontSize={'1rem'}>
            {label}
          </Typography>

          {enumOptions.map((option, index) => (
            <FormControlLabel
              key={option.value}
              value={option.value}
              control={<Radio />}
              label={
                <div>
                  <Typography variant="subtitle1">{option.label}</Typography>
                  <Typography variant="body2" color="textSecondary" textTransform={'none'}>
                    {schema.enumDescriptions[index]}
                  </Typography>
                </div>
              }
            />
          ))}
        </RadioGroup>
      </FormControl>
    );
  };

  const widgets = {
    RadioWidget: CustomRadioWidget,
  };

  return (
    <>
      <Modal
        open={isImportModalOpen}
        closeModal={handleClose}
        maxWidth="sm"
        title="Import Model"
        style={{
          zIndex: 1500,
        }}
      >
        {activeStep === 0 ? (
          <RJSFModalWrapper
            schema={importModelSchema}
            uiSchema={importModelUiSchema}
            handleSubmit={handleImportModelSubmit}
            handleNext={handleNext}
            submitBtnText="Next"
            handleClose={handleClose}
            widgets={widgets}
            helpText={
              <p>
                {importModalDescription} <br />
                Learn more about importing Models in our{' '}
                <StyledDocsRedirectLink
                  href={`${MESHERY_DOCS_URL}/guides/configuration-management/importing-models`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  documentation
                </StyledDocsRedirectLink>
                .
              </p>
            }
          />
        ) : (
          <FinishDeploymentStep deploymentType="register" handleClose={handleClose} />
        )}
      </Modal>
      <Modal
        open={isCsvModalOpen}
        closeModal={() => setIsCsvModalOpen(false)}
        maxWidth="sm"
        title="Import CSV"
        style={{
          zIndex: 1500,
        }}
      >
        <CsvStepper handleClose={handleClose} />
      </Modal>
    </>
  );
});

ImportModelModal.displayName = 'ImportModelModal';

export default ImportModelModal;
