import React, { useEffect } from 'react';
import { CircularProgress } from '@sistent/sistent';
import { COMPONENTS } from '../../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';
import { getFilteredDataForDetailsComponent } from './helper';
import { useGetComponentsFromModalQuery } from '@/rtk-query/meshModel';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { useRegistryRouter } from './hooks';

const VersionedModelComponentTree = ({
  registrantID,
  modelDef,
  versionedModelDef,
  setShowDetailsData,
  showDetailsData,
}) => {
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
      const selectedIdArr = selectedItemUUID.split('.');
      // Check if selected item is a component
      if (
        selectedIdArr.length > 0 &&
        selectedIdArr[0] === modelDef.id &&
        selectedIdArr[1] === versionedModelDef.id
      ) {
        const showData = getFilteredDataForDetailsComponent(
          componentsData.components,
          selectedIdArr[selectedIdArr.length - 1],
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
        message: `There was an error fetching components data: ${error?.data}`,
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
