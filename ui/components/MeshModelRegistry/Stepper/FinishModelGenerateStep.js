import React from 'react';
import { useContext } from 'react';
import { capitalize } from 'lodash';
import { Loading } from '@/components/DesignLifeCycle/common';
import { NotificationCenterContext } from '../../NotificationCenter';
import { useEffect } from 'react';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import {
  ModelImportedSection,
  ModelImportMessages,
} from '@/components/NotificationCenter/formatters/model_registration';
import { ErrorMetadataFormatter } from '@/components/NotificationCenter/formatters/error';

const FinishModelGenerateStep = ({ requestBody, generateType }) => {
  const [generateEvent, setGenerateEvent] = React.useState();
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const [importMeshModel, { isLoading, error }] = useImportMeshModelMutation();

  useEffect(() => {
    const performImport = async () => {
      await importMeshModel({ importBody: requestBody });
    };

    performImport();
  }, []);

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
