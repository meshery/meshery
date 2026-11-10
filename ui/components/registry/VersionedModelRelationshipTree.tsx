import React, { useEffect, useRef } from 'react';
import { CircularProgress } from '@sistent/sistent';
import { MODELS, RELATIONSHIPS } from '../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';
import { groupRelationshipsByKind } from './helper';
import { useGetRelationshipsFromModalQuery } from '@/rtk-query/meshModel';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import RelationshipTree from './RelationshipTree';
import { useRegistryRouter } from './hooks';

type VersionedModelRelationshipTreeProps = {
  registrantID?: string;
  modelDef: any;
  versionedModelDef: any;
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  showDetailsData: { type: string; data: any };
  handleToggle: (_event: unknown, _nodeIds: string[]) => void;
  handleSelect: (_event: unknown, _nodeIds: string[]) => void;
  selected: string[];
  expanded: string[];
};

const VersionedModelRelationshipTree = ({
  registrantID,
  modelDef,
  versionedModelDef,
  setShowDetailsData,
  showDetailsData,
  handleToggle,
  handleSelect,
  selected,
  expanded,
}: VersionedModelRelationshipTreeProps) => {
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

  const { selectedItemUUID } = useRegistryRouter();

  // Stable ref so the deeplink effect can read showDetailsData without
  // subscribing to it as a dependency (which would cause an infinite loop).
  const showDetailsDataRef = useRef(showDetailsData);
  useEffect(() => {
    showDetailsDataRef.current = showDetailsData;
  }, [showDetailsData]);

  useEffect(() => {
    if (relationshipsData?.relationships?.length > 0) {
      const selectedIdArr = selectedItemUUID.split('.');
      // Only act when the deeplink targets a relationship inside this model+version.
      // selectedIdArr[2] === '2' is the Relationships group node constant from itemId.
      if (
        selectedIdArr.length > 3 &&
        selectedIdArr[0] === modelDef.id &&
        selectedIdArr[1] === versionedModelDef.id &&
        selectedIdArr[2] === '2'
      ) {
        const relationshipId = selectedIdArr[selectedIdArr.length - 1];
        const relationship = relationshipsData.relationships.find((r) => r.id === relationshipId);
        if (relationship && relationship.id !== showDetailsDataRef.current?.data?.id) {
          setShowDetailsData({ type: RELATIONSHIPS, data: relationship });
        }
      }
    }
    // showDetailsData intentionally omitted — read via showDetailsDataRef to
    // prevent the infinite loop described above.
  }, [relationshipsData, selectedItemUUID, modelDef.id, versionedModelDef.id]);

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
          itemId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.2`}
          data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.2`}
          labelText={`Relationships (${
            relationshipsData?.relationships ? relationshipsData.relationships.length : 0
          })`}
        >
          {relationshipsData?.relationships?.length > 0 && (
            <RelationshipTree
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              data={groupRelationshipsByKind(relationshipsData.relationships || [])}
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
