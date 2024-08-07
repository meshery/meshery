import React, { useEffect, useState, useRef, useCallback } from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { IconButton, FormControlLabel, Switch } from '@material-ui/core';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import SearchBar from '../../utils/custom-search';
import debounce from '../../utils/debounce';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import { useWindowDimensions } from '../../utils/dimension';
import StyledTreeItem from './StyledTreeItem';
import { useRouter } from 'next/router';
import { getFilteredDataForDetailsComponent, groupRelationshipsByKind } from './helper';
import { CustomTextTooltip } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import _ from 'lodash';
import CollapseAllIcon from '@/assets/icons/CollapseAll';
import ExpandAllIcon from '@/assets/icons/ExpandAll';
import CircularProgress from '@mui/material/CircularProgress';
import { Colors } from '../../themes/app';
import { JustifyAndAlignCenter } from './MeshModel.style';
import { styled } from '@mui/styles';
import {
  useGetComponentsFromModalQuery,
  useGetRelationshipsFromModalQuery,
} from '@/rtk-query/meshModel';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';

const VersionedModelComponentTree = ({
  registrantID,
  modelDef,
  versionedModelDef,
  setShowDetailsData,
}) => {
  const { notify } = useNotification();
  const {
    data: componentsData,
    isLoading,
    isError,
    error,
  } = useGetComponentsFromModalQuery({
    model: versionedModelDef.name,
    params: { version: versionedModelDef.model.version },
  });

  useEffect(() => {
    if (isError) {
      notify({
        message: `There was an error fetching components data: ${error?.data}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [error, isError]);

  return (
    <>
      {isLoading ? (
        <CircularProgress color="inherit" />
      ) : (
        <StyledTreeItem
          nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.1`}
          data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.1`}
          labelText={`Components (${
            componentsData.components ? componentsData.components.length : 0
          })`}
        >
          {componentsData.components &&
            componentsData.components.map((component, subIndex) => {
              return (
                <StyledTreeItem
                  key={subIndex}
                  nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                    versionedModelDef.id
                  }.1.${component.id}`}
                  data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
                    versionedModelDef.id
                  }.1.${component.id}`}
                  labelText={component.displayName}
                  onClick={() => {
                    setShowDetailsData({
                      type: COMPONENTS,
                      data: component,
                    });
                  }}
                />
              );
            })}
        </StyledTreeItem>
      )}
    </>
  );
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
}) => {
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
        message: `There was an error fetching relationships data: ${error?.data}`,
        event_type: EVENT_TYPES.ERROR,
      });
    }
  }, [error, isError]);

  return (
    <>
      {isLoading ? (
        <CircularProgress color="inherit" />
      ) : (
        <StyledTreeItem
          nodeId={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.2`}
          data-id={`${registrantID ? `${registrantID}.1.` : ''}${modelDef.id}.${
            versionedModelDef.id
          }.2`}
          labelText={`Relationships (${
            relationshipsData.relationships ? relationshipsData.relationships.length : 0
          })`}
        >
          {relationshipsData.relationships.length > 0 && (
            <RelationshipTree
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              data={groupRelationshipsByKind(relationshipsData.relationships)}
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

const ComponentTree = ({
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  setShowDetailsData,
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
    </TreeView>
  );
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
                labelText={relationship.subType}
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
    </TreeView>
  );
};

const MesheryTreeViewItem = ({
  modelDef,
  registrantID,
  setShowDetailsData,
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
        <div
          style={{
            display: 'flex',
            alignItems: 'center',
            gap: '1rem',
          }}
        >
          {imgSrc ? <img src={imgSrc} style={{ height: '1.5rem', width: '1.5rem' }} /> : null}
          <span>{modelDef.displayName ? modelDef.displayName : modelDef.name}</span>
        </div>
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

const MesheryTreeViewModel = ({
  data,
  handleToggle,
  handleSelect,
  expanded,
  selected,
  setShowDetailsData,
  lastModelRef,
  isModelFetching,
}) => {
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
        />
      ))}
      <div ref={lastModelRef} style={{ height: '3rem' }}></div>
      {isModelFetching ? <CircularProgress color="inherit" /> : null}
    </TreeView>
  );
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
          labelText={registrant?.hostname}
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
              labelText={`Models (${registrant?.models?.length})`}
            >
              {registrant?.models.map((modelDef, index) => (
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
                />
              ))}
            </StyledTreeItem>
          </div>
        </StyledTreeItem>
      ))}
      <div ref={lastRegistrantRef} style={{ height: '3rem' }}></div>
      {isRegistrantFetching && <CircularProgress color="inherit" />}
    </TreeView>
  );
};

const useRegistryRouter = () => {
  const router = useRouter();
  const { query, push: pushRoute, route } = router;

  const settingsCategory = query.settingsCategory;
  const tab = query.tab;
  const selectedItemUUID = query.selectedItemUUID || '';
  const searchText = query.searchText || null;
  let filters = {
    searchText: searchText,
  };

  const handleUpdateSelectedRoute = (nodeIds, filters) => {
    const id = nodeIds[0];
    const queryString = new URLSearchParams({
      settingsCategory,
      tab,
      selectedItemUUID: id,
      ...filters,
    }).toString();
    pushRoute(`${route}?${queryString}`, undefined, { shallow: true });
  };

  return {
    settingsCategory,
    tab,
    handleUpdateSelectedRoute,
    selectedItemUUID,
    filters,
  };
};

const MesheryTreeViewWrapper = styled('div')(() => ({
  display: 'flex',
  flexDirection: 'column',
}));
const MesheryTreeView = ({
  data,
  view,
  setSearchText,
  searchText,
  setPage,
  checked,
  setChecked,
  setShowDetailsData,
  showDetailsData,
  setResourcesDetail,
  lastItemRef,
  isFetching,
}) => {
  const { handleUpdateSelectedRoute, selectedItemUUID } = useRegistryRouter();
  const [expanded, setExpanded] = React.useState([]);
  const [selected, setSelected] = React.useState([]);
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(searchText ? true : false);
  const scrollRef = useRef();

  const handleScroll = (scrollingView) => (event) => {
    const div = event.target;
    if (div.scrollTop >= div.scrollHeight - div.clientHeight - 1) {
      setPage((prevPage) => ({
        ...prevPage,
        [scrollingView]: Number(prevPage[scrollingView]) + 1,
      }));
    }
    if (!data.length === 0) {
      scrollRef.current = div.scrollTop;
    }
  };

  useEffect(() => {
    if (scrollRef.current) {
      const div = document.querySelector('.scrollElement');
      div.scrollTop = scrollRef.current;
    }
  }, [data]);

  const handleChecked = useCallback(() => {
    setChecked((prevChecked) => !prevChecked);
  }, [setChecked]);

  // Expand first level tree
  const expandAll = () => {
    const arr = [];
    data.map((parent) => {
      if (view === RELATIONSHIPS) {
        // Hard coded for relationships
        // parent id will be same as relationships[0].id
        // so we can use that id for expanding first level tree for relationships
        arr.push(parent.relationships[0].id);
      } else {
        arr.push(parent.id);
      }
    });
    setExpanded(arr);
  };

  const handleSelect = (event, nodeIds) => {
    if (nodeIds.length >= 0) {
      let selectedIdArr = nodeIds[0].split('.');
      let indx = data.findIndex((item) => item.id === selectedIdArr[0]);

      // Filter object contains current filter applied to data
      // Route will contain filters to support deeplink
      const filter = {
        ...(searchText && { searchText }),
        pagesize: indx + 14,
      };
      handleUpdateSelectedRoute(nodeIds, filter);
      setSelected([0, nodeIds]);
    } else {
      setSelected([]);
    }
  };

  const handleToggle = (event, nodeIds) => {
    setExpanded(nodeIds);
  };

  useEffect(() => {
    let selectedIdArr = selectedItemUUID.split('.');
    if (selectedIdArr.length >= 0) {
      // Check if showDetailsData data matches with item from route
      // This will prevent unnecessary state update
      if (showDetailsData.data.id !== selectedIdArr[selectedIdArr.length - 1]) {
        setExpanded(
          selectedIdArr.reduce(
            (acc, id, index) => [...acc, index > 0 ? `${acc[index - 1]}.${id}` : id],
            [],
          ),
        );
        setSelected([selectedItemUUID]);
        const showData = getFilteredDataForDetailsComponent(data, selectedItemUUID);
        setShowDetailsData(showData);
      }
    } else {
      setExpanded([]);
      setSelected([]);
      setShowDetailsData({
        type: '',
        data: {},
      });
    }
  }, [view, selectedItemUUID]);

  useEffect(() => {
    let selectedIdArr = selectedItemUUID.split('.');
    if (selectedIdArr.length >= 0) {
      setTimeout(() => {
        requestAnimationFrame(() => {
          const selectedNode = document.querySelector(`[data-id="${selectedItemUUID}"]`);
          if (selectedNode) {
            selectedNode.scrollIntoView({ behavior: 'smooth' });
          }
        });
      }, 1000);
    }
  }, [view]);

  const disabledExpand = () => {
    return view === COMPONENTS;
  };

  const renderHeader = (type, hasRecords) => (
    <div
      style={{
        display: 'flex',
        flexDirection: 'row',
        justifyContent: 'space-between',
        borderBottom: '1px solid #d2d3d4',
        width: '100%',
        height: '2.55rem',
      }}
    >
      <div>
        {width < 1370 && isSearchExpanded ? null : (
          <div style={{ display: 'flex', flexDirection: 'row' }}>
            <CustomTextTooltip title="Expand All" placement="top">
              {/* span is added to make sure tooltip is not listening to disabled elements to prevent MUI error */}
              <span>
                <IconButton
                  onClick={expandAll}
                  size="large"
                  disableRipple
                  disabled={disabledExpand()}
                >
                  <ExpandAllIcon height={17} width={17} />
                </IconButton>
              </span>
            </CustomTextTooltip>

            <CustomTextTooltip title="Collapse All" placement="top">
              <span>
                <IconButton
                  onClick={() => setExpanded([])}
                  style={{ marginRight: '4px' }}
                  size="large"
                  disableRipple
                  disabled={disabledExpand()}
                >
                  <CollapseAllIcon height={17} width={17} />
                </IconButton>
              </span>
            </CustomTextTooltip>
            {type === MODELS && (
              <>
                <FormControlLabel
                  control={
                    <Switch
                      color="primary"
                      checked={checked}
                      onClick={handleChecked}
                      disabled={!hasRecords}
                      inputProps={{ 'aria-label': 'controlled' }}
                    />
                  }
                  label="Show Duplicates"
                  style={{ margin: 0 }}
                />
                <CustomTextTooltip
                  placement="right"
                  interactive={true}
                  title={`View all duplicate entries of ${_.toLower(
                    view,
                  )}. Entries with identical name and version attributes are considered duplicates. [Learn More](https://docs.meshery.io/concepts/logical/models#models)`}
                >
                  <IconButton color="primary">
                    <InfoOutlinedIcon height={20} width={20} />
                  </IconButton>
                </CustomTextTooltip>
              </>
            )}
          </div>
        )}
      </div>
      <div style={{ display: 'flex' }}>
        <SearchBar
          onSearch={debounce((value) => setSearchText(value), 200)}
          expanded={isSearchExpanded}
          setExpanded={setSearchExpand}
          placeholder="Search"
          value={searchText}
        />
      </div>
    </div>
  );

  const setSearchExpand = (isExpand) => {
    if (!isExpand) {
      setSearchText(() => null);
      setResourcesDetail(() => []);
    }
    setIsSearchExpanded(isExpand);
  };

  const renderTree = (treeComponent, type) => (
    <>
      {renderHeader(type, !!data.length)}
      {data.length === 0 && !searchText ? (
        <JustifyAndAlignCenter style={{ height: '27rem' }}>
          <CircularProgress sx={{ color: Colors.keppelGreen }} />
        </JustifyAndAlignCenter>
      ) : data.length === 0 && searchText ? (
        <JustifyAndAlignCenter style={{ height: '27rem' }}>
          <p>No result found</p>
        </JustifyAndAlignCenter>
      ) : (
        <div
          className="scrollElement"
          style={{ overflowY: 'auto', height: '55vh' }}
          onScroll={handleScroll(type)}
        >
          {treeComponent}
        </div>
      )}
    </>
  );

  return (
    <MesheryTreeViewWrapper style={{ width: '100%', height: '100%' }}>
      {view === MODELS &&
        renderTree(
          <MesheryTreeViewModel
            data={data}
            handleToggle={handleToggle}
            handleSelect={handleSelect}
            expanded={expanded}
            selected={selected}
            setShowDetailsData={setShowDetailsData}
            lastModelRef={lastItemRef[MODELS]}
            isModelFetching={isFetching[MODELS]}
          />,
          MODELS,
        )}
      {view === REGISTRANTS &&
        renderTree(
          <MesheryTreeViewRegistrants
            data={data}
            handleToggle={handleToggle}
            handleSelect={handleSelect}
            expanded={expanded}
            selected={selected}
            setShowDetailsData={setShowDetailsData}
            lastRegistrantRef={lastItemRef[REGISTRANTS]}
            isRegistrantFetching={isFetching[REGISTRANTS]}
          />,
          REGISTRANTS,
        )}
      {view === COMPONENTS &&
        renderTree(
          <ComponentTree
            handleToggle={handleToggle}
            handleSelect={handleSelect}
            expanded={expanded}
            selected={selected}
            data={data}
            setExpanded={setExpanded}
            setSelected={setSelected}
            handleScroll={handleScroll}
            setSearchText={setSearchText}
            setShowDetailsData={setShowDetailsData}
          />,
          COMPONENTS,
        )}
      {view === RELATIONSHIPS &&
        renderTree(
          <RelationshipTree
            handleToggle={handleToggle}
            handleSelect={handleSelect}
            expanded={expanded}
            selected={selected}
            data={data}
            setShowDetailsData={setShowDetailsData}
          />,
          RELATIONSHIPS,
        )}
    </MesheryTreeViewWrapper>
  );
};

export default MesheryTreeView;
