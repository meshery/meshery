import React from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { CircularProgress } from '@sistent/sistent';
import { COMPONENTS } from '../../../constants/navigator';
import MinusSquare from '../../../assets/icons/MinusSquare';
import PlusSquare from '../../../assets/icons/PlusSquare';
import DotSquare from '../../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';

type ComponentTreeProps = {
  expanded: string[];
  selected: string[];
  handleToggle: (_event: unknown, _nodeIds: string[]) => void;
  handleSelect: (_event: unknown, _nodeIds: string[]) => void;
  data: any[];
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  lastComponentRef: React.MutableRefObject<any>;
  isComponentFetching: boolean;
};

const ComponentTree = ({
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  setShowDetailsData,
  lastComponentRef,
  isComponentFetching,
}: ComponentTreeProps) => {
  return (
    <TreeView
      aria-label="controlled"
      defaultExpandedItems={['3']}
      slots={{
        collapseIcon: MinusSquare,
        expandIcon: PlusSquare,
        endIcon: DotSquare,
      }}
      onExpandedItemsChange={handleToggle}
      onSelectedItemsChange={handleSelect}
      multiSelect
      expandedItems={expanded}
      selectedItems={selected}
    >
      {data.map((component, index) => (
        <StyledTreeItem
          key={index}
          nodeId={`${component.id}`}
          data-id={`${component.id}`}
          labelText={component.displayName}
          onClick={() => {
            setShowDetailsData({
              type: COMPONENTS,
              data: component,
            });
          }}
        />
      ))}
      <div ref={lastComponentRef} style={{ height: '48px' }}></div>
      {isComponentFetching ? <CircularProgress /> : null}
    </TreeView>
  );
};

export default ComponentTree;
