import React, { useEffect, useRef } from 'react';
import { CircularProgress } from '@sistent/sistent';
import { COMPONENTS } from '../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';
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

  // Stable ref so the deeplink effect can read showDetailsData without
  // subscribing to it as a dependency (which would cause an infinite loop).
  const showDetailsDataRef = useRef(showDetailsData);
  useEffect(() => {
    showDetailsDataRef.current = showDetailsData;
  }, [showDetailsData]);

  useEffect(() => {
    if (componentsData?.components?.length > 0) {
      const selectedIdArr = selectedItemUUID.split('.');
      // Only act when the deeplink targets a component inside this model+version.
      // selectedIdArr[2] === '1' is the Components group node constant from itemId.
      if (
        selectedIdArr.length > 3 &&
        selectedIdArr[0] === modelDef.id &&
        selectedIdArr[1] === versionedModelDef.id &&
        selectedIdArr[2] === '1'
      ) {
        const componentId = selectedIdArr[selectedIdArr.length - 1];
        const component = componentsData.components.find((c) => c.id === componentId);
        if (component && component.id !== showDetailsDataRef.current?.data?.id) {
          setShowDetailsData({ type: COMPONENTS, data: component });
        }
      }
    }
    // showDetailsData intentionally omitted — read via showDetailsDataRef to
    // prevent the infinite loop described above.
  }, [componentsData, selectedItemUUID, modelDef.id, versionedModelDef.id]);

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
          itemId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
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
                  itemId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                    versionedModelDef.id
                  }.1.${component.id}`}
                  data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                    versionedModelDef.id
                  }.1.${component.id}`}
                  labelText={component.displayName}
                  onClick={(e) => {
                    e.stopPropagation();
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
