import React from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { CircularProgress } from '@layer5/sistent';
import { COMPONENTS } from '../../../constants/navigator';
import MinusSquare from '../../../assets/icons/MinusSquare';
import PlusSquare from '../../../assets/icons/PlusSquare';
import DotSquare from '../../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';

const ComponentTree = ({
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  setShowDetailsData,
  lastComponentRef,
  isComponentFetching,
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
      <div ref={lastComponentRef} style={{ height: '3rem' }}></div>
      {isComponentFetching ? <CircularProgress /> : null}
    </TreeView>
  );
};

export default ComponentTree;
