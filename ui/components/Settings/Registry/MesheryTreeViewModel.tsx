import React from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { CircularProgress, Box } from '@sistent/sistent';
import MinusSquare from '../../../assets/icons/MinusSquare';
import PlusSquare from '../../../assets/icons/PlusSquare';
import DotSquare from '../../../assets/icons/DotSquare';
import MesheryTreeViewItem from './MesheryTreeViewItem';

type MesheryTreeViewModelProps = {
  data: any[];
  handleToggle: (_event: unknown, _nodeIds: string[]) => void;
  handleSelect: (_event: unknown, _nodeIds: string[]) => void;
  expanded: string[];
  selected: string[];
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  lastModelRef: React.MutableRefObject<any>;
  isModelFetching: boolean;
  showDetailsData: { type: string; data: any };
};

const MesheryTreeViewModel = ({
  data,
  handleToggle,
  handleSelect,
  expanded,
  selected,
  setShowDetailsData,
  lastModelRef,
  isModelFetching,
  showDetailsData,
}: MesheryTreeViewModelProps) => {
  return (
    <TreeView
      aria-label="controlled"
      defaultCollapseIcon={<MinusSquare />}
      defaultExpandIcon={<PlusSquare />}
      defaultEndIcon={<DotSquare />}
      onNodeToggle={handleToggle}
      onNodeSelect={handleSelect}
      multiSelect
      expanded={expanded}
      selected={selected}
    >
      {data.map((modelDef, index) => (
        <MesheryTreeViewItem
          key={index}
          modelDef={modelDef}
          handleToggle={handleToggle}
          handleSelect={handleSelect}
          expanded={expanded}
          selected={selected}
          setShowDetailsData={setShowDetailsData}
          showDetailsData={showDetailsData}
        />
      ))}
      <div ref={lastModelRef} style={{ height: '3rem' }}>
        {isModelFetching && (
          <Box display="flex" justifyContent="center" alignItems="center" padding={2}>
            <CircularProgress size={24} />
          </Box>
        )}
      </div>
    </TreeView>
  );
};

export default MesheryTreeViewModel;
