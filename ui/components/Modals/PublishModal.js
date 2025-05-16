import React, { useState, useEffect } from 'react';
import Modal from '../Modal';
import PublicIcon from '@mui/icons-material/Public';
import _ from 'lodash';
import { getMeshModels } from '../../api/meshmodel';
import { modifyRJSFSchema } from '../../utils/utils';
import { useGetSchemaQuery } from '@/rtk-query/schema';

// This modal is used in Meshery Extensions also
export default function PublishModal(props) {
  const { open, title, handleClose, handleSubmit } = props;
  const [publishSchema, setPublishSchema] = useState({});
  const { data: schemaData, isSuccess } = useGetSchemaQuery({ schemaName: 'publish' });

  const processSchema = async () => {
    if (isSuccess && schemaData) {
      try {
        const { models } = await getMeshModels();
        const modelNames = _.uniq(models?.map((model) => model.displayName));
        modelNames.sort();

        const modifiedSchema = modifyRJSFSchema(
          schemaData.rjsfSchema,
          'properties.compatibility.items.enum',
          modelNames,
        );

        setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: schemaData.uiSchema });
      } catch (err) {
        console.error(err);
        setPublishSchema(schemaData);
      }
    }
  };

  useEffect(() => {
    processSchema();
  }, [isSuccess, schemaData]);

  return (
    <Modal
      open={open}
      schema={publishSchema.rjsfSchema}
      uiSchema={publishSchema.uiSchema}
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
