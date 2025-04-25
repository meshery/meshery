import React from 'react';
import { MODELS } from '../../constants/navigator';
import StyledTreeItem from './StyledTreeItem';

import { StyledTreeItemDiv, StyledTreeItemNameDiv } from './MeshModel.style';
import VersionedModelComponentTree from './VersionedModelComponentTree';
import VersionedModelRelationshipTree from './VersionedModelRelationshipTree';

const MesheryTreeViewItem = ({
  modelDef,
  registrantID,
  setShowDetailsData,
  showDetailsData,
  handleToggle,
  handleSelect,
  selected,
  expanded,
}) => {
  const imgSrc = modelDef?.metadata?.svgColor;
  return (
    <StyledTreeItem
      key={modelDef.id}
      nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}`}
      data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}`}
      top
      labelText={
        <StyledTreeItemDiv>
          {imgSrc ? <img src={imgSrc} style={{ height: '1.5rem', width: '1.5rem' }} /> : null}
          <StyledTreeItemNameDiv>
            {modelDef.displayName ? modelDef.displayName : modelDef.name}
          </StyledTreeItemNameDiv>
        </StyledTreeItemDiv>
      }
      onClick={() => {
        setShowDetailsData({
          type: MODELS,
          data: modelDef,
        });
      }}
    >
      {modelDef.versionBasedData &&
        modelDef.versionBasedData.map((versionedModelDef) => (
          <StyledTreeItem
            key={versionedModelDef.id}
            nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
              versionedModelDef.id
            }`}
            data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
              versionedModelDef.id
            }`}
            labelText={
              versionedModelDef?.model?.version?.[0] == 'v'
                ? versionedModelDef?.model?.version
                : `v${versionedModelDef?.model?.version}`
            }
            check={true}
            onClick={() => {
              setShowDetailsData({
                type: MODELS,
                data: versionedModelDef,
              });
            }}
          >
            <VersionedModelComponentTree
              registrantID={registrantID}
              modelDef={modelDef}
              versionedModelDef={versionedModelDef}
              setShowDetailsData={setShowDetailsData}
              showDetailsData={showDetailsData}
            />
            <VersionedModelRelationshipTree
              registrantID={registrantID}
              modelDef={modelDef}
              versionedModelDef={versionedModelDef}
              setShowDetailsData={setShowDetailsData}
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              selected={selected}
              expanded={expanded}
            />
          </StyledTreeItem>
        ))}
    </StyledTreeItem>
  );
};

export default MesheryTreeViewItem;
