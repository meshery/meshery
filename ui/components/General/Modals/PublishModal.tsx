import React, { useEffect, useState } from 'react';
import Modal from './Modal';
import _ from 'lodash';
import { getMeshModels } from '../../api/meshmodel';
import { modifyRJSFSchema } from '../../utils/utils';
import { useGetSchemaQuery } from '@/rtk-query/schema';
import { PublicIcon } from '@sistent/sistent';

type PublishSchemaState = {
  rjsfSchema?: Record<string, unknown>;
  uiSchema?: Record<string, unknown>;
};

type PublishModalProps = {
  open: boolean;
  title?: string;
  handleClose: () => void;
  handleSubmit: (formData: Record<string, unknown>) => void;
};

// This modal is used in Meshery Extensions also
export default function PublishModal(props: PublishModalProps) {
  const { open, title, handleClose, handleSubmit } = props;
  const [publishSchema, setPublishSchema] = useState<PublishSchemaState>({});
  const { data: schemaData, isSuccess } = useGetSchemaQuery({ schemaName: 'publish' });

  useEffect(() => {
    let cancelled = false;

    const processSchema = async () => {
      if (!isSuccess || !schemaData) {
        return;
      }

      try {
        const { models } = await getMeshModels();
        const modelNames = _.uniq(models?.map((model) => model.displayName));
        modelNames.sort();

        const modifiedSchema = modifyRJSFSchema(
          schemaData.rjsfSchema,
          'properties.compatibility.items.enum',
          modelNames,
        );

        if (!cancelled) {
          setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: schemaData.uiSchema });
        }
      } catch (error) {
        console.error(error);
        if (!cancelled) {
          setPublishSchema(schemaData);
        }
      }
    };

    processSchema();

    return () => {
      cancelled = true;
    };
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
