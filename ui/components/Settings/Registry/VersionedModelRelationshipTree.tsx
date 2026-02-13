import React, { useEffect } from 'react';
import { CircularProgress } from '@sistent/sistent';
import { MODELS } from '../../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';
import { groupRelationshipsByKind } from './helper';
import { useGetRelationshipsFromModalQuery } from '@/rtk-query/meshModel';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import RelationshipTree from './RelationshipTree';

const getErrorMessage = (err: unknown): string => {
  if (!err) return 'Unknown error';
  if (typeof err === 'string') return err;
  if (typeof err !== 'object') return String(err);

  if ('data' in err) return JSON.stringify((err as { data?: unknown }).data);
  if ('error' in err && typeof (err as { error?: unknown }).error === 'string') {
    return (err as { error: string }).error;
  }
  if ('message' in err && typeof (err as { message?: unknown }).message === 'string') {
    return (err as { message: string }).message;
  }

  return 'Unknown error';
};

type VersionedModelRelationshipTreeProps = {
  registrantID?: string;
  modelDef: any;
  versionedModelDef: any;
  setShowDetailsData: (_data: { type: string; data: any }) => void;
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

  useEffect(() => {
    if (isError) {
      notify({
        message: `There was an error fetching relationships data: ${getErrorMessage(error)}`,
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
