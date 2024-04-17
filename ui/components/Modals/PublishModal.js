import React, { useState, useEffect } from 'react';
import Modal from '../Modal';
import PublicIcon from '@material-ui/icons/Public';
import { getMeshModels } from '../../api/meshmodel';
import { modifyRJSFSchema } from '../../utils/utils';
import { publishSchema, publishUiSchema } from '@layer5/sistent';

// This modal is used in MeshMap also
export default function PublishModal(props) {
  const { open, title, handleClose, handleSubmit } = props;
  const [setPublishSchema] = useState({});

  useEffect(() => {
    async function fetchMeshModels() {
      try {
        const { models } = await getMeshModels();
        const modelNames = models?.map((model) => model.displayName) || [];
        modelNames.sort(); // Sort model names

        // Modify the schema to include mesh models
        const modifiedSchema = modifyRJSFSchema(
          publishSchema, // Use publishSchema as the base schema
          'properties.compatibility.items.enum',
          modelNames,
        );

        // Set the modified schema and UI schema
        setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: publishUiSchema });
      } catch (error) {
        console.error('Error fetching mesh models:', error);
      }
    }

    fetchMeshModels();
  }, []);

  return (
    <Modal
      open={open}
      schema={publishSchema}
      uiSchema={publishUiSchema}
      title={title}
      handleClose={handleClose}
      handleSubmit={handleSubmit}
      submitBtnText="Submit for Approval"
      submitBtnIcon={<PublicIcon data-cy="import-button" />}
      showInfoIcon={{
        text: 'Upon submitting your catalog item, an approval flow will be initiated.',
        link: 'https://docs.meshery.io/concepts/catalog',
      }}
    />
  );
}
