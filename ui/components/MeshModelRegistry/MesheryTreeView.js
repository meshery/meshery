import React, { useEffect, useState, useRef } from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { IconButton, FormControlLabel, Switch, Tooltip } from '@material-ui/core';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import SearchBar from '../../utils/custom-search';
import debounce from '../../utils/debounce';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { useWindowDimensions } from '../../utils/dimension';
import StyledTreeItem from './StyledTreeItem';
import { useRouter } from 'next/router';
import { getFilteredDataForDetailsComponent } from './helper';

const ComponentTree = ({
  setComp,
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  handleScroll,
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
      <div
        id="scrollElement"
        style={{ overflowY: 'auto', height: '27rem' }}
        onScroll={handleScroll(COMPONENTS)}
      >
        {data.map((component, index) => (
          <StyledTreeItem
            key={index}
            nodeId={`${component.id}`}
            check
            labelText={component.displayName}
            onClick={() => {
              setComp(component);
            }}
          />
        ))}
      </div>
    </TreeView>
  );
};

const RelationshipTree = ({
  setRela,
  expanded,
  selected,
  handleToggle,
  handleSelect,
  data,
  handleScroll,
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
      <div
        id="scrollElement"
        style={{ overflowY: 'auto', maxHeight: '27rem' }}
        onScroll={handleScroll(RELATIONSHIPS)}
      >
        {data.map((relationship, index) => (
          <StyledTreeItem
            key={index}
            nodeId={`${relationship.id}`}
            check
            labelText={relationship.subType}
            onClick={() => {
              setRela(relationship);
            }}
          />
        ))}
      </div>
    </TreeView>
  );
};

const MesheryTreeViewItem = ({
  model,
  setShow,
}) => {
  return (
    <StyledTreeItem
      key={model.id}
      nodeId={model.id}
      top
      labelText={model.displayName}
      onClick={() => {
        setShow({
          model: model,
          components: [],
          relationships: [],
        });
      }}
    >
      <StyledTreeItem
        nodeId={`${model.id}.1`}
        labelText={`Components (${model.components ? model.components.length : 0})`}
      >
        {model.components &&
          model.components.map((component, subIndex) => (
            <StyledTreeItem
              key={subIndex}
              nodeId={`${model.id}.1.${component.id}`}
              check
              labelText={component.displayName}
              data-id={`${model.id}.1.${component.id}`}
              onClick={() => {
                setShow((prevShow) => {
                  const { components } = prevShow;
                  const compIndex = components.findIndex((item) => item === component);
                  if (compIndex !== -1) {
                    return {
                      ...prevShow,
                      model: model,
                      components: components.filter((item) => item !== component),
                    };
                  } else {
                    return {
                      ...prevShow,
                      model: model,
                      components: [...components, component],
                    };
                  }
                });
              }}
            />
          ))}
      </StyledTreeItem>
      <StyledTreeItem
        nodeId={`${model.id}.2`}
        labelText={`Relationships (${model.relationships ? model.relationships.length : 0})`}
      >
        {model.relationships &&
          model.relationships.map((relationship, subIndex) => (
            <StyledTreeItem
              key={subIndex}
              nodeId={`${model.id}.2.${relationship.id}`}
              check
              labelText={relationship.displayhostname}
              onClick={() => {
                setShow((prevShow) => {
                  const { relationships } = prevShow;
                  const relaIndex = relationships.findIndex((item) => item === relationship);
                  if (relaIndex !== -1) {
                    return {
                      ...prevShow,
                      model: model,
                      relationships: relationships.filter((item) => item !== relationship),
                    };
                  } else {
                    return {
                      ...prevShow,
                      model: model,
                      relationships: [...relationships, relationship],
                    };
                  }
                });
              }}
            />
          ))}
      </StyledTreeItem>
    </StyledTreeItem>
  );
};

const MesheryTreeViewModel = ({
  data,
  setShow,
  handleToggle,
  handleSelect,
  expanded,
  selected,
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
      {data.map((model, index) => (
        <MesheryTreeViewItem
          key={index}
          model={model}
          handleToggle={handleToggle}
          handleSelect={handleSelect}
          expanded={expanded}
          selected={selected}
          setShow={setShow}
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
          top
          labelText={registrant.hostname}
          newParentId={registrant.id}
          onClick={() => {
            setShow({
              model: {},
              components: [],
              relationships: [],
            });
          }}
        >
          <div>
            <StyledTreeItem
              nodeId={`${registrant.id}.1`}
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
  comp,
  rela,
  setShow,
  setComp,
  setRela,
  setRegi,
  setSearchText,
  searchText,
  setPage,
  checked,
  setChecked,
}) => {
  const { handleUpdateSelectedRoute, selectedItemUUID } = useRegistryRouter();
  const [expanded, setExpanded] = React.useState([]);
  const [selected, setSelected] = React.useState([]);
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
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
      const div = document.getElementById('scrollElement');
      div.scrollTop = scrollRef.current;
    }
  }, [data]);

  const handleChecked = () => {
    setChecked(!checked);
  };

  const expandAll = () => {
    const arr = [];
    data.map((model, index) => {
      arr.push(index);
      arr.push(`${model.id}.1`);
      arr.push(`${model.id}.2`);
    });
    setExpanded(arr);
  };

  const handleSelect = (event, nodeIds) => {
    if (nodeIds.length >= 0) {
      let selectedIdArr = selectedItemUUID.split('.')
      let indx = data.findIndex((item) => item.id === selectedIdArr[0]);

      // Filter object contains current filter applied to data
      // Route will contain filters to support deeplink
      const filter = {
        ...(searchText && { searchText }),
        pagesize: indx + 5,
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
    let selectedIdArr = selectedItemUUID.split('.')
    if (selectedIdArr.length >= 0) {
      setExpanded(
        selectedIdArr.reduce(
          (acc, id, index) => [...acc, index > 0 ? `${acc[index - 1]}.${id}` : id],
          [],
        ),
      );
      setSelected([selectedItemUUID]);

      const { selectedComponent, selectedModel, selectedRelationship } =
        getFilteredDataForDetailsComponent(data, selectedItemUUID, view);
      setShow({
        model: selectedModel || {},
        components: selectedComponent ? [selectedComponent] : [],
        relationships: selectedRelationship ? [selectedRelationship] : [],
      });

      const selectedNode = document.querySelector(`[data-id="${selectedItemUUID}"]`);
      if (selectedNode) {
        selectedNode.scrollIntoView({ behavior: 'smooth' });
      }
    } else {
      setExpanded([]);
      setSelected([]);
      setShow({
        model: {},
        components: [],
        relationships: [],
      });
    }
  }, [view, data]);

  useEffect(() => {
    setComp({});
    setRela({});
    setRegi({});
  }, [view]);

  console.log('data-->', data);

  return (
    <div style={{ width: '100%' }}>
      {view === MODELS && (
        <div>
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
                  <Tooltip title="Expand All" placement="top">
                    <IconButton onClick={expandAll} size="large" disableRipple>
                      {/* <PlusSquare /> */}
                      <KeyboardArrowDownIcon />
                    </IconButton>
                  </Tooltip>

                  <Tooltip title="Collapse All" placement="top">
                    <IconButton
                      onClick={() => setExpanded([])}
                      style={{ marginRight: '4px' }}
                      size="large"
                      disableRipple
                    >
                      {/* <MinusSquare /> */}
                      <KeyboardArrowUpIcon />
                    </IconButton>
                  </Tooltip>
                  <FormControlLabel
                    control={
                      <Switch
                        color="primary"
                        checked={checked}
                        onClick={handleChecked}
                        inputProps={{ 'aria-label': 'controlled' }}
                      />
                    }
                    label="Duplicates"
                  />
                </div>
              )}
            </div>
            <div style={{ display: 'flex' }}>
              <SearchBar
                onSearch={debounce((value) => setSearchText(value), 200)}
                expanded={isSearchExpanded}
                setExpanded={setIsSearchExpanded}
                placeholder="Search"
              />
            </div>
          </div>
          <div
            id="scrollElement"
            style={{
              overflowY: 'auto',
              height: '27rem',
            }}
            onScroll={handleScroll(MODELS)}
          >
            <MesheryTreeViewModel
              data={data}
              setShow={setShow}
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
            />
          </div>
        </div>
      )}
      {view === REGISTRANTS && (
        <div>
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              justifyContent: 'space-between',
              borderBottom: '1px solid #d2d3d4',
            }}
          >
            <div style={{ display: 'flex', flexDirection: 'row' }}>
              <IconButton onClick={expandAll} size="large">
                {/* <PlusSquare /> */}
                <KeyboardArrowDownIcon />
              </IconButton>

              <IconButton
                onClick={() => setExpanded([])}
                style={{ marginRight: '4px' }}
                size="large"
              >
                {/* <MinusSquare /> */}
                <KeyboardArrowUpIcon />
              </IconButton>
            </div>
          </div>
          <div
            id="scrollElement"
            style={{ overflowY: 'auto', height: '27rem' }}
            onScroll={handleScroll(REGISTRANTS)}
          >
            <MesheryTreeViewRegistrants
              data={data}
              setShow={setShow}
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
            />
          </div>
        </div>
      )}
      {view === COMPONENTS && (
        <ComponentTree
          handleToggle={handleToggle}
          handleSelect={handleSelect}
          expanded={expanded}
          selected={selected}
          data={data}
          comp={comp}
          setExpanded={setExpanded}
          setComp={setComp}
          setSelected={setSelected}
          handleScroll={handleScroll}
          setSearchText={setSearchText}
        />
      )}
      {view === RELATIONSHIPS && (
        <RelationshipTree
          handleToggle={handleToggle}
          handleSelect={handleSelect}
          expanded={expanded}
          selected={selected}
          data={data}
          rela={rela}
          setExpanded={setExpanded}
          setRela={setRela}
          setSelected={setSelected}
          handleScroll={handleScroll}
          setSearchText={setSearchText}
        />
      )}
    </div>
  );
};

export default MesheryTreeView;
