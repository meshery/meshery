import React from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { CircularProgress } from '@layer5/sistent';
import { REGISTRANTS } from '../../constants/navigator';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';
import MesheryTreeViewItem from './MesheryTreeViewItem';

const MesheryTreeViewRegistrants = ({
  data,
  setShow,
  handleToggle,
  handleSelect,
  expanded,
  selected,
  setShowDetailsData,
  lastRegistrantRef,
  isRegistrantFetching,
  showDetailsData,
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
      {data?.map((registrant) => (
        <StyledTreeItem
          key={registrant.id}
          nodeId={registrant.id}
          data-id={registrant.id}
          top
          labelText={registrant?.name}
          newParentId={registrant.id}
          onClick={() => {
            setShowDetailsData({
              type: REGISTRANTS,
              data: registrant,
            });
          }}
        >
          <div>
            <StyledTreeItem
              nodeId={`${registrant.id}.1`}
              data-id={`${registrant.id}.1`}
              labelText={`Models (${registrant?.summary?.models})`}
            >
              {registrant.models
                .filter((model) => model?.registrant?.kind == registrant?.kind)
                ?.map((modelDef, index) => (
                  <MesheryTreeViewItem
                    key={index}
                    modelDef={modelDef}
                    handleToggle={handleToggle}
                    handleSelect={handleSelect}
                    expanded={expanded}
                    selected={selected}
                    setShow={setShow}
                    registrantID={registrant.id}
                    setShowDetailsData={setShowDetailsData}
                    showDetailsData={showDetailsData}
                  />
                ))}
            </StyledTreeItem>
          </div>
        </StyledTreeItem>
      ))}
      <div ref={lastRegistrantRef} style={{ height: '3rem' }}></div>
      {isRegistrantFetching && <CircularProgress />}
    </TreeView>
  );
};

export default MesheryTreeViewRegistrants;
