import React, { useState } from 'react';
import {
  Modal,
  FormControlLabel,
  FormControl,
  RadioGroup,
  Radio,
  importModelUiSchema,
  importModelSchema,
  Typography,
} from '@layer5/sistent';
import { RJSFModalWrapper } from '../Modal';
import CsvStepper, { StyledDocsRedirectLink } from './Stepper/CSVStepper';
import { MESHERY_DOCS_URL } from '@/constants/endpoints';
import { getUnit8ArrayDecodedFile } from '@/utils/utils';
import { updateProgress } from 'lib/store';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';

const ImportModelModal = React.memo(({ isImportModalOpen, setIsImportModalOpen }) => {
  const [importModalDescription, setImportModalDescription] = useState('');
  const [isCsvModalOpen, setIsCsvModalOpen] = useState(false);
  const [importModelReq] = useImportMeshModelMutation();
  const handleGenerateModal = async (data) => {
    const { component_csv, model_csv, relationship_csv, register } = data;
    let requestBody = {
      importBody: {
        model_csv: model_csv,
        component_csv: component_csv,
        relationship_csv: relationship_csv,
      },
      uploadType: 'csv',
      register: register,
    };

    updateProgress({ showProgress: true });
    await importModelReq({ importBody: requestBody });
    updateProgress({ showProgress: false });
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
        setIsImportModalOpen(false);
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
        closeModal={() => setIsImportModalOpen(false)}
        maxWidth="sm"
        title="Import Model"
      >
        <RJSFModalWrapper
          schema={importModelSchema}
          uiSchema={importModelUiSchema}
          handleSubmit={handleImportModelSubmit}
          submitBtnText="Import"
          handleClose={() => setIsImportModalOpen(false)}
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
      </Modal>
      <Modal
        open={isCsvModalOpen}
        closeModal={() => setIsCsvModalOpen(false)}
        maxWidth="sm"
        title="Import CSV"
      >
        <CsvStepper
          handleGenerateModal={handleGenerateModal}
          handleClose={() => setIsCsvModalOpen(false)}
        />
      </Modal>
    </>
  );
});

ImportModelModal.displayName = 'ImportModelModal';

export default ImportModelModal;
