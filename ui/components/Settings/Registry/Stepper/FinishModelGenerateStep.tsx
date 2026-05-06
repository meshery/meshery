import React, { useContext, useEffect, useRef } from 'react';
import { capitalize } from 'lodash';
import { Loading } from '@/components/DesignLifeCycle/common';
import { NotificationCenterContext } from '@/components/NotificationCenter';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import {
  ModelImportedSection,
  ModelImportMessages,
} from '@/components/NotificationCenter/formatters/model_registration';
import { ErrorMetadataFormatter } from '@/components/NotificationCenter/formatters/error';

const isImportRequestReady = (requestBody: any) => {
  const importBody = requestBody?.importBody;

  if (!requestBody?.uploadType || !importBody) {
    return false;
  }

  switch (requestBody.uploadType) {
    case 'csv':
      return Boolean(
        importBody.model_csv && importBody.component_csv && importBody.relationship_csv,
      );
    case 'file':
      return Boolean(importBody.model_file);
    case 'url':
    case 'urlImport':
      return Boolean(importBody.url);
    default:
      return true;
  }
};

const FinishModelGenerateStep = ({
  requestBody,
  generateType,
}: {
  requestBody: any;
  generateType: string;
}) => {
  const [generateEvent, setGenerateEvent] = React.useState<any>();
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const [importMeshModel, { isLoading, error }] = useImportMeshModelMutation();
  const hasStartedImportRef = useRef(false);

  useEffect(() => {
    const subscription = operationsCenterActorRef.on(
      OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      (event) => {
        const serverEvent = event.data.event;
        if (serverEvent.action === generateType) {
          setGenerateEvent(serverEvent);
        }
      },
    );

    return () => subscription.unsubscribe();
  }, [operationsCenterActorRef, generateType]);

  useEffect(() => {
    if (hasStartedImportRef.current || !isImportRequestReady(requestBody)) {
      return;
    }

    hasStartedImportRef.current = true;

    const performImport = async () => {
      await importMeshModel({ importBody: requestBody });
    };

    void performImport();
  }, [importMeshModel, requestBody]);

  const progressMessage = `${capitalize(generateType)}ing model`;

  if (isLoading) {
    return <Loading message={progressMessage} />;
  }
  if (error) {
    return (
      <ErrorMetadataFormatter metadata={generateEvent?.metadata.error} event={generateEvent} />
    );
  }

  return (
    <>
      <ModelImportMessages message={generateEvent?.metadata?.ModelImportMessage} />
      <ModelImportedSection modelDetails={generateEvent?.metadata?.ModelDetails} />
    </>
  );
};

export default FinishModelGenerateStep;
