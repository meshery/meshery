import React from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { CircularProgress } from '@sistent/sistent';
import { REGISTRANTS } from '@/constants/navigator';
import MinusSquare from '../../../assets/icons/MinusSquare';
import PlusSquare from '../../../assets/icons/PlusSquare';
import DotSquare from '../../../assets/icons/DotSquare';
import StyledTreeItem from './StyledTreeItem';
import MesheryTreeViewItem from './MesheryTreeViewItem';

type MesheryTreeViewRegistrantsProps = {
  data: any[];
  setShow?: (_value: any) => void;
  handleToggle: (_event: unknown, _nodeIds: string[]) => void;
  handleSelect: (_event: unknown, _nodeIds: string[]) => void;
  expanded: string[];
  selected: string[];
  setShowDetailsData: (_data: { type: string; data: any }) => void;
  lastRegistrantRef: React.MutableRefObject<any>;
  isRegistrantFetching: boolean;
  showDetailsData: { type: string; data: any };
};

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
}: MesheryTreeViewRegistrantsProps) => {
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
      {data
        ?.filter((item) => item?.summary || item?.models)
        ?.map((registrant) => (
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
                nodeId={`${registrant?.id || 'unknown'}.1`}
                data-id={`${registrant?.id || 'unknown'}.1`}
                labelText={
                  registrant?.summary?.models !== undefined
                    ? `Models (${registrant.summary.models})`
                    : 'Models (0)'
                }
              >
                {(registrant?.models || [])
                  .filter((model) => model?.registrant?.kind === registrant?.kind)
                  .map((modelDef, index) => (
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
      <div ref={lastRegistrantRef} style={{ height: '48px' }}></div>
      {isRegistrantFetching && <CircularProgress />}
    </TreeView>
  );
};

export default MesheryTreeViewRegistrants;
