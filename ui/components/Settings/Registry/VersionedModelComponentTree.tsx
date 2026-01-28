import React, { useEffect } from 'react';
import { CircularProgress } from '@sistent/sistent';
import { COMPONENTS } from '../../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';
import { getFilteredDataForDetailsComponent } from './helper';
import { useGetComponentsFromModalQuery } from '@/rtk-query/meshModel';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useRegistryRouter } from './hooks';

type VersionedModelComponentTreeProps = {
  registrantID?: string;
  modelDef: any;
  versionedModelDef: any;
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  showDetailsData: { type: string; data: any };
};

const getErrorMessage = (err: unknown): string => {
  if (!err) return 'Unknown error';
  if (typeof err === 'string') return err;
  if (typeof err !== 'object') return String(err);

  // RTK Query fetch errors often include `data` and/or `error`
  if ('data' in err) return JSON.stringify((err as { data?: unknown }).data);
  if ('error' in err && typeof (err as { error?: unknown }).error === 'string') {
    return (err as { error: string }).error;
  }
  if ('message' in err && typeof (err as { message?: unknown }).message === 'string') {
    return (err as { message: string }).message;
  }

  return 'Unknown error';
};

const VersionedModelComponentTree = ({
  registrantID,
  modelDef,
  versionedModelDef,
  setShowDetailsData,
  showDetailsData,
}: VersionedModelComponentTreeProps) => {
  const { notify } = useNotification();
  const {
    data: componentsData,
    isLoading,
    isError,
    error,
  } = useGetComponentsFromModalQuery({
    model: versionedModelDef.name,
    params: { version: versionedModelDef.model.version },
  });
  const { selectedItemUUID } = useRegistryRouter();

  useEffect(() => {
    if (componentsData && componentsData.components && componentsData.components.length > 0) {
      const selectedItemUUIDStr = Array.isArray(selectedItemUUID)
        ? (selectedItemUUID[0] ?? '')
        : (selectedItemUUID ?? '');
      const selectedIdArr = selectedItemUUIDStr.split('.');
      // Check if selected item is a component
      if (
        selectedIdArr.length > 0 &&
        selectedIdArr[0] === modelDef.id &&
        selectedIdArr[1] === versionedModelDef.id
      ) {
        const selectedComponentId = selectedIdArr[selectedIdArr.length - 1] ?? '';
        const showData = getFilteredDataForDetailsComponent(
          componentsData.components,
          selectedComponentId,
        );
        if (JSON.stringify(showData) !== JSON.stringify(showDetailsData)) {
          setShowDetailsData(showData);
        }
      }
    }
  }, [componentsData, showDetailsData]);

  useEffect(() => {
    if (isError) {
      notify({
        message: `There was an error fetching components data: ${getErrorMessage(error)}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [error, isError]);

  return (
    <>
      {isLoading ? (
        <CircularProgress />
      ) : (
        <StyledTreeItem
          nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.1`}
          data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.1`}
          labelText={`Components (${
            componentsData.components ? componentsData.components.length : 0
          })`}
        >
          {componentsData.components &&
            componentsData.components.map((component, subIndex) => {
              return (
                <StyledTreeItem
                  key={subIndex}
                  nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                    versionedModelDef.id
                  }.1.${component.id}`}
                  data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                    versionedModelDef.id
                  }.1.${component.id}`}
                  labelText={component.displayName}
                  onClick={() => {
                    setShowDetailsData({
                      type: COMPONENTS,
                      data: component,
                    });
                  }}
                />
              );
            })}
        </StyledTreeItem>
      )}
    </>
  );
};

export default VersionedModelComponentTree;
