import React from 'react';
import { SimpleTreeView } from '../shared/TreeView';
import { CircularProgress } from '@sistent/sistent';
import { styled } from '@/theme';
import { CONNECTIONS } from '../../constants/navigator';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';

// Sentinel element observed by the infinite-scroll ref to trigger the next page.
const InfiniteScrollSentinel = styled('div')({
  height: '48px',
});

type ConnectionDefinitionTreeProps = {
  expanded: string[];
  selected: string[];
  handleToggle: (_event: unknown, _nodeIds: string[]) => void;
  handleSelect: (_event: unknown, _nodeIds: string[]) => void;
  data: any[];
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  lastConnectionRef: React.MutableRefObject<any>;
  isConnectionFetching: boolean;
};

const ConnectionDefinitionTree = ({
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  setShowDetailsData,
  lastConnectionRef,
  isConnectionFetching,
}: ConnectionDefinitionTreeProps) => {
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
      {data.map((connection, index) => (
        <StyledTreeItem
          key={connection.id || index}
          itemId={`${connection.id || index}`}
          data-id={`${connection.id || index}`}
          labelText={`${connection.name || connection.kind}${
            connection.kind ? ` (${connection.kind})` : ''
          }`}
          onClick={() => {
            setShowDetailsData({
              type: CONNECTIONS,
              data: connection,
            });
          }}
        />
      ))}
      <InfiniteScrollSentinel ref={lastConnectionRef} />
      {isConnectionFetching ? <CircularProgress /> : null}
    </SimpleTreeView>
  );
};

export default ConnectionDefinitionTree;
