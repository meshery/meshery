import React from 'react';
import { SimpleTreeView } from '../shared/TreeView';
import { CircularProgress } from '@sistent/sistent';
import { RELATIONSHIPS } from '@/constants/navigator';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';

type RelationshipTreeProps = {
  expanded: string[];
  selected: string[];
  handleToggle: (_event: unknown, _nodeIds: string[]) => void;
  handleSelect: (_event: unknown, _nodeIds: string[]) => void;
  data: any[];
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  view?: string;
  idForKindAsProp?: string;
  lastRegistrantRef?: React.MutableRefObject<any>;
  isRelationshipFetching?: boolean;
};

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
}: RelationshipTreeProps) => {
  // Build the tree items — shared between both rendering modes below.
  const treeItems = (
    <>
      {data.map((relationshipByKind, index) => {
        const idForKind =
          view === RELATIONSHIPS
            ? `${relationshipByKind.relationships[0].id}`
            : `${idForKindAsProp}.${relationshipByKind.relationships[0].id}`;
        return (
          <StyledTreeItem
            key={idForKind}
            itemId={idForKind}
            data-id={idForKind}
            labelText={`${relationshipByKind.kind} (${relationshipByKind.relationships.length})`}
            onClick={(e) => {
              const target = e.target;
              if (!(target instanceof Element)) {
                return;
              }
              const treeItem = target.closest('[data-id]');
              if (treeItem?.getAttribute('data-id') !== idForKind) {
                return;
              }

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
                key={relationship.id}
                itemId={`${idForKind}.${relationship.id}`}
                data-id={`${idForKind}.${relationship.id}`}
                labelText={`${relationship.subType} (${relationship.model.name})`}
                onClick={(e) => {
                  e.stopPropagation();
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
      <div ref={lastRegistrantRef} style={{ height: '48px' }}></div>
      {isRelationshipFetching ? <CircularProgress color="inherit" /> : null}
    </>
  );

  // When RelationshipTree is nested inside a Models tree item (view !== RELATIONSHIPS),
  // do NOT create a new SimpleTreeView. A SimpleTreeView rendered inside a StyledTreeItem
  // of an outer SimpleTreeView causes MUI's traverseDescendants plugin to produce circular
  // itemId references in its internal item map → RangeError: Maximum call stack size exceeded.
  // Instead, render items directly so they register with the outer SimpleTreeView context,
  // exactly as VersionedModelComponentTree does.
  if (view !== RELATIONSHIPS) {
    return treeItems;
  }

  // Top-level Relationships view — needs its own SimpleTreeView root.
  return (
    <SimpleTreeView
      aria-label="controlled"
      slots={{ collapseIcon: MinusSquare, expandIcon: PlusSquare, endIcon: DotSquare }}
      onExpandedItemsChange={handleToggle}
      onSelectedItemsChange={handleSelect}
      multiSelect
      expandedItems={expanded}
      selectedItems={selected}
    >
      {treeItems}
    </SimpleTreeView>
  );
};

export default RelationshipTree;
