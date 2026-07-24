import { useCallback } from 'react';
import {
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
  useSaveEnvironmentMutation,
} from '../../rtk-query/environments';
import { useUpdateConnectionByIdMutation } from '@/rtk-query/connection';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import { ACTION_TYPES, getErrorMessage } from './ConnectionTable.constants';

type UseConnectionActionsArgs = {
  organizationId?: string;
};

export const useConnectionActions = ({ organizationId }: UseConnectionActionsArgs) => {
  const { notify } = useNotification();
  const [updateConnectionByIdMutator] = useUpdateConnectionByIdMutation();
  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();
  const [saveEnvironmentMutator] = useSaveEnvironmentMutation();

  const addConnectionToEnvironment = useCallback(
    async (
      environmentId: string,
      environmentName: string,
      connectionId: string,
      connectionName: string,
    ) => {
      try {
        await addConnectionToEnvironmentMutator({ environmentId, connectionId }).unwrap();
        notify({
          message: `Connection: ${connectionName} assigned to environment: ${environmentName}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [addConnectionToEnvironmentMutator, notify],
  );

  const removeConnectionFromEnvironment = useCallback(
    async (
      environmentId: string,
      environmentName: string,
      connectionId: string,
      connectionName: string,
    ) => {
      try {
        await removeConnectionFromEnvMutator({ environmentId, connectionId }).unwrap();
        notify({
          message: `Connection: ${connectionName} removed from environment: ${environmentName}`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [notify, removeConnectionFromEnvMutator],
  );

  const saveEnvironment = useCallback(
    async (connectionId: string, connectionName: string, environmentName: string) => {
      try {
        const response = await saveEnvironmentMutator({
          body: {
            name: environmentName,
            organization_id: organizationId,
          },
        }).unwrap();

        notify({
          message: `Environment "${response.name}" created`,
          event_type: EVENT_TYPES.SUCCESS,
        });

        await addConnectionToEnvironment(response.id, response.name, connectionId, connectionName);
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.WorkspaceManagementCreateEnvironment.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [addConnectionToEnvironment, notify, organizationId, saveEnvironmentMutator],
  );

  const updateConnectionStatus = useCallback(
    async (connectionId: string, newStatus: string) => {
      try {
        await updateConnectionByIdMutator({
          connectionId,
          body: { status: newStatus },
        }).unwrap();

        notify({
          message: `Connection status updated`,
          event_type: EVENT_TYPES.SUCCESS,
        });
      } catch (error) {
        notify({
          message: `${ACTION_TYPES.UPDATE_CONNECTION.error_msg}: ${getErrorMessage(error)}`,
          event_type: EVENT_TYPES.ERROR,
          details: String(error),
        });
      }
    },
    [notify, updateConnectionByIdMutator],
  );

  return {
    notify,
    updateConnectionByIdMutator,
    addConnectionToEnvironment,
    removeConnectionFromEnvironment,
    saveEnvironment,
    updateConnectionStatus,
  };
};
