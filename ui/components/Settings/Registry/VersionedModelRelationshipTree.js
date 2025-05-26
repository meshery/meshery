import React, { useEffect } from 'react';
import { CircularProgress } from '@layer5/sistent';
import { MODELS } from '../../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';
import { groupRelationshipsByKind } from './helper';
import { useGetRelationshipsFromModalQuery } from '@/rtk-query/meshModel';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import RelationshipTree from './RelationshipTree';

const VersionedModelRelationshipTree = ({
  registrantID,
  modelDef,
  versionedModelDef,
  setShowDetailsData,
  handleToggle,
  handleSelect,
  selected,
  expanded,
}) => {
  const { notify } = useNotification();
  const {
    data: relationshipsData,
    isLoading,
    isError,
    error,
  } = useGetRelationshipsFromModalQuery({
    model: versionedModelDef.name,
    params: { version: versionedModelDef.model.version },
  });

  useEffect(() => {
    if (isError) {
      notify({
        message: `There was an error fetching relationships data: ${error?.data}`,
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
          }.2`}
          data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.2`}
          labelText={`Relationships (${
            relationshipsData.relationships ? relationshipsData.relationships.length : 0
          })`}
        >
          {relationshipsData.relationships.length > 0 && (
            <RelationshipTree
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              data={groupRelationshipsByKind(relationshipsData.relationships)}
              view={MODELS}
              setShowDetailsData={setShowDetailsData}
              idForKindAsProp={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                versionedModelDef.id
              }.2`}
            />
          )}
        </StyledTreeItem>
      )}
    </>
  );
};

export default VersionedModelRelationshipTree;
