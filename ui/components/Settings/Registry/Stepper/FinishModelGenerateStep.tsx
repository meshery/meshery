import React from 'react';
import { useContext } from 'react';
import { capitalize } from 'lodash';
import { Loading } from '@/components/DesignLifeCycle/common';
import { NotificationCenterContext } from '@/components/NotificationCenter';
import { useEffect } from 'react';
import { OPERATION_CENTER_EVENTS } from 'machines/operationsCenter';
import { useImportMeshModelMutation } from '@/rtk-query/meshModel';
import {
  ModelImportedSection,
  ModelImportMessages,
} from '@/components/NotificationCenter/formatters/model_registration';
import { ErrorMetadataFormatter } from '@/components/NotificationCenter/formatters/error';

const FinishModelGenerateStep = ({
  requestBody,
  generateType,
}: {
  requestBody: any;
  generateType: string;
}) => {
  const [generateEvent, setGenerateEvent] = React.useState<any>();
  const { operationsCenterActorRef } = useContext(NotificationCenterContext);
  const [registerMeshmodels, { isLoading, error }] = useImportMeshModelMutation();

  useEffect(() => {
    const performImport = async () => {
      await registerMeshmodels({ body: requestBody });
    };

    performImport();
  }, [registerMeshmodels, requestBody]);

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
