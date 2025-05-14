import React from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { CircularProgress } from '@layer5/sistent';
import { RELATIONSHIPS } from '../../constants/navigator';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';

const RelationshipTree = ({
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  setShowDetailsData,
  view = RELATIONSHIPS,
  idForKindAsProp,
  lastRegistrantRef,
  isRelationshipFetching,
}) => {
  return (
    <TreeView
      aria-label="controlled"
      defaultExpanded={['3']}
      defaultCollapseIcon={<MinusSquare />}
      defaultExpandIcon={<PlusSquare />}
      defaultEndIcon={<DotSquare />}
      onNodeToggle={handleToggle}
      onNodeSelect={handleSelect}
      multiSelect
      expanded={expanded}
      selected={selected}
    >
      {data.map((relationshipByKind, index) => {
        const idForKind =
          view === RELATIONSHIPS
            ? `${relationshipByKind.relationships[0].id}`
            : `${idForKindAsProp}.${relationshipByKind.relationships[0].id}`;
        return (
          <StyledTreeItem
            key={index}
            nodeId={idForKind}
            data-id={idForKind}
            labelText={`${relationshipByKind.kind} (${relationshipByKind.relationships.length})`}
            onClick={() => {
              setShowDetailsData({
                type: 'none',
                data: {
                  id: relationshipByKind.relationships[0].id,
                },
              });
            }}
          >
            {relationshipByKind.relationships.map((relationship) => (
              <StyledTreeItem
                key={index}
                nodeId={`${idForKind}.${relationship.id}`}
                data-id={`${idForKind}.${relationship.id}`}
                labelText={`${relationship.subType} (${relationship.model.name})`}
                onClick={() => {
                  setShowDetailsData({
                    type: RELATIONSHIPS,
                    data: relationship,
                  });
                }}
              />
            ))}
          </StyledTreeItem>
        );
      })}
      <div ref={lastRegistrantRef} style={{ height: '3rem' }}></div>
      {isRelationshipFetching ? <CircularProgress color="inherit" /> : null}
    </TreeView>
  );
};

export default RelationshipTree;
