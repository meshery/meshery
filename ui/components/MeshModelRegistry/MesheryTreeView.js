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
import { getFilteredDataForDetailsComponent } from './helper';
import { CustomTextTooltip } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import _ from 'lodash';
import CollapseAllIcon from '@/assets/icons/CollapseAll';
import ExpandAllIcon from '@/assets/icons/ExpandAll';

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
          check
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
      {data.map((relationship, index) => (
        <StyledTreeItem
          key={index}
          nodeId={`${relationship.id}`}
          data-id={`${relationship.id}`}
          check
          labelText={relationship.subType}
          onClick={() => {
            setShowDetailsData({
              type: RELATIONSHIPS,
              data: relationship,
            });
          }}
        />
      ))}
    </TreeView>
  );
};

const MesheryTreeViewItem = ({ model, registrantID, setShowDetailsData }) => {
  return (
    <StyledTreeItem
      key={model.id}
      nodeId={`${registrantID ? `${registrantID}.1.` : ''}${model.id}`}
      data-id={`${registrantID ? `${registrantID}.1.` : ''}${model.id}`}
      top
      labelText={model.displayName}
      onClick={() => {
        setShowDetailsData({
          type: MODELS,
          data: model,
        });
      }}
    >
      {model.versionBasedData &&
        model.versionBasedData.map((versionedModel) => (
          <StyledTreeItem
            key={versionedModel.id}
            nodeId={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${versionedModel.id}`}
            labelText={versionedModel.version}
            onClick={() => {
              setShowDetailsData({
                type: MODELS,
                data: versionedModel,
              });
            }}
          >
            <StyledTreeItem
              nodeId={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                versionedModel.id
              }.1`}
              data-id={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                versionedModel.id
              }.1`}
              labelText={`Components (${
                versionedModel.components ? versionedModel.components.length : 0
              })`}
            >
              {versionedModel.components &&
                versionedModel.components.map((component, subIndex) => (
                  <StyledTreeItem
                    key={subIndex}
                    nodeId={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                      versionedModel.id
                    }.1.${component.id}`}
                    data-id={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                      versionedModel.id
                    }.1.${component.id}`}
                    check
                    labelText={component.displayName}
                    onClick={() => {
                      setShowDetailsData({
                        type: COMPONENTS,
                        data: component,
                      });
                    }}
                  />
                ))}
            </StyledTreeItem>
            <StyledTreeItem
              nodeId={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                versionedModel.id
              }.2`}
              data-id={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                versionedModel.id
              }.2`}
              labelText={`Relationships (${
                versionedModel.relationships ? versionedModel.relationships.length : 0
              })`}
            >
              {versionedModel.relationships &&
                versionedModel.relationships.map((relationship, subIndex) => (
                  <StyledTreeItem
                    key={subIndex}
                    nodeId={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                      versionedModel.id
                    }.2.${relationship.id}`}
                    data-id={`${registrantID ? `${registrantID}.1.` : ''}${model.id}.${
                      versionedModel.id
                    }.2.${relationship.id}`}
                    check
                    labelText={relationship.subType}
                    onClick={() => {
                      setShowDetailsData({
                        type: COMPONENTS,
                        data: relationship,
                      });
                    }}
                  />
                ))}
            </StyledTreeItem>
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
      {data.map((model, index) => (
        <MesheryTreeViewItem
          key={index}
          model={model}
          handleToggle={handleToggle}
          handleSelect={handleSelect}
          expanded={expanded}
          selected={selected}
          setShowDetailsData={setShowDetailsData}
        />
      ))}
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
          labelText={registrant.hostname}
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
              {registrant?.models.map((model, index) => (
                <MesheryTreeViewItem
                  key={index}
                  model={model}
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

const MesheryTreeView = ({
  data,
  view,
  setSearchText,
  searchText,
  setPage,
  checked,
  setChecked,
  setShowDetailsData,
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

    scrollRef.current = div.scrollTop;
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

  const expandAll = () => {
    const arr = [];
    data.map((parent) => {
      arr.push(parent.id);
      arr.push(`${parent.id}.1`);
      arr.push(`${parent.id}.2`);
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
      setExpanded(
        selectedIdArr.reduce(
          (acc, id, index) => [...acc, index > 0 ? `${acc[index - 1]}.${id}` : id],
          [],
        ),
      );
      setSelected([selectedItemUUID]);

      const showData = getFilteredDataForDetailsComponent(data, selectedItemUUID);
      setShowDetailsData(showData);
    } else {
      setExpanded([]);
      setSelected([]);
      setShowDetailsData({
        type: '',
        data: {},
      });
    }
  }, [view]);

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
      }, 200);
    }
  }, [view]);

  const disabledExpand = () => {
    return view === RELATIONSHIPS || view === COMPONENTS;
  };

  const renderHeader = (type) => (
    <div
      style={{
        display: 'flex',
        flexDirection: 'row',
        justifyContent: 'space-between',
        borderBottom: '1px solid #d2d3d4',
      }}
    >
      <div>
        {width < 1370 && isSearchExpanded ? null : (
          <div style={{ display: 'flex', flexDirection: 'row' }}>
            <CustomTextTooltip title="Expand All" placement="top">
              <IconButton
                onClick={expandAll}
                size="large"
                disableRipple
                disabled={disabledExpand()}
              >
                <ExpandAllIcon height={15} width={15} />
              </IconButton>
            </CustomTextTooltip>

            <CustomTextTooltip title="Collapse All" placement="top">
              <IconButton
                onClick={() => setExpanded([])}
                style={{ marginRight: '4px' }}
                size="large"
                disableRipple
                disabled={disabledExpand()}
              >
                <CollapseAllIcon height={15} width={15} />
              </IconButton>
            </CustomTextTooltip>
            {type === MODELS && (
              <>
                <FormControlLabel
                  control={
                    <Switch
                      color="primary"
                      checked={checked}
                      onClick={handleChecked}
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
                  )}. Entries with identical name and version attributes are considered duplicates.`}
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
          setExpanded={setIsSearchExpanded}
          placeholder="Search"
          value={searchText}
        />
      </div>
    </div>
  );

  const renderTree = (treeComponent, type) => (
    <div>
      {renderHeader(type)}
      <div
        className="scrollElement"
        style={{ overflowY: 'auto', height: '27rem' }}
        onScroll={handleScroll(type)}
      >
        {treeComponent}
      </div>
    </div>
  );

  return (
    <div style={{ width: '100%', height: '28.86rem' }}>
      {view === MODELS &&
        renderTree(
          <MesheryTreeViewModel
            data={data}
            handleToggle={handleToggle}
            handleSelect={handleSelect}
            expanded={expanded}
            selected={selected}
            setShowDetailsData={setShowDetailsData}
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
            setExpanded={setExpanded}
            setSelected={setSelected}
            handleScroll={handleScroll}
            setSearchText={setSearchText}
            setShowDetailsData={setShowDetailsData}
          />,
          RELATIONSHIPS,
        )}
    </div>
  );
};

export default MesheryTreeView;
