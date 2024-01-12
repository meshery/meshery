import React, { useEffect, useState, useRef } from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import {
  Box,
  Typography,
  IconButton,
  FormControlLabel,
  Switch,
  useTheme,
  Tooltip,
} from '@material-ui/core';
import Checkbox from '@mui/material/Checkbox';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import SearchBar from '../../utils/custom-search';
import debounce from '../../utils/debounce';
import { StyledTreeItemRoot } from './MeshModel.style';
import MinusSquare from '../../assets/icons/MinusSquare';
import PlusSquare from '../../assets/icons/PlusSquare';
import DotSquare from '../../assets/icons/DotSquare';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';
import { useWindowDimensions } from '../../utils/dimension';
import StyledTreeItem from './StyledTreeItem';

const ComponentTree = ({
  setExpanded,
  setComp,
  setSelected,
  expanded,
  selected,
  handleToggle,
  handleSelect,
  comp,
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
      <StyledTreeItem
        nodeId={0}
        top
        root
        search
        setSearchText={setSearchText}
        labelText={
          comp.model?.name ? `Model: ${comp.model?.displayName}` : 'Select a component node'
        }
        onClick={() => {
          setExpanded([0]);
          setComp({});
          setSelected([]);
        }}
      >
        <div
          id="scrollElement"
          style={{ overflowY: 'auto', height: '27rem' }}
          onScroll={handleScroll(COMPONENTS)}
        >
          {data.map((component, index) => (
            <StyledTreeItem
              key={index}
              nodeId={index + 1}
              check
              labelText={component.displayName}
              onClick={() => {
                setComp(component);
              }}
            />
          ))}
        </div>
      </StyledTreeItem>
    </TreeView>
  );
};

const RelationshipTree = ({
  setExpanded,
  setRela,
  setSelected,
  expanded,
  selected,
  handleToggle,
  handleSelect,
  rela,
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
      <StyledTreeItem
        nodeId={0}
        root
        top
        setSearchText={setSearchText}
        labelText={
          rela.model?.name ? `Model: ${rela.model?.displayName}` : 'Select a relationship node'
        }
        onClick={() => {
          setExpanded([0]);
          setRela({});
          setSelected([]);
        }}
      >
        <div
          id="scrollElement"
          style={{ overflowY: 'auto', maxHeight: '27rem' }}
          onScroll={handleScroll(RELATIONSHIPS)}
        >
          {data.map((relationship, index) => (
            <StyledTreeItem
              key={index}
              nodeId={index + 1}
              check
              labelText={relationship.subType}
              onClick={() => {
                setRela(relationship);
              }}
            />
          ))}
        </div>
      </StyledTreeItem>
    </TreeView>
  );
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
  setPage,
  checked,
  setChecked,
}) => {
  const [expanded, setExpanded] = React.useState([]);
  const [selected, setSelected] = React.useState([]);
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  useEffect(() => {
    setSelected([]);
    if (view === COMPONENTS || view === RELATIONSHIPS) {
      setExpanded([0]);
    } else {
      setExpanded([0]);
    }
    setComp({});
    setRela({});
    setRegi({});
    setShow({
      model: {},
      components: [],
      relationships: [],
    });
  }, [view]);

  const scrollRef = useRef();

  const handleScroll = (scrollingView) => (event) => {
    const div = event.target;
    if (div.scrollTop >= div.scrollHeight - div.clientHeight - 1) {
      setPage((prevPage) => ({
        ...prevPage, // Keep the current values for other keys
        [scrollingView]: prevPage[scrollingView] + 1, // Increment the specific key based on the view
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
      arr.push(`${index}.1`);
      arr.push(`${index}.2`);
    });
    setExpanded(arr);
  };

  const handleSelect = (event, nodeIds) => {
    if (nodeIds !== 0) {
      setSelected([0, nodeIds]);
    } else {
      setSelected([]);
    }
  };

  const handleToggle = (event, nodeIds) => {
    setExpanded(nodeIds);
  };

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
            <TreeView
              aria-label="controlled"
              defaultExpanded={['3']}
              defaultCollapseIcon={<MinusSquare />}
              defaultExpandIcon={<PlusSquare />}
              defaultEndIcon={<DotSquare />}
              onNodeToggle={handleToggle}
              multiSelect
              expanded={expanded}
            >
              {data.map((model, index) => {
                return (
                  <StyledTreeItem
                    key={index}
                    top
                    nodeId={index}
                    check
                    newParentId={`${model.name}`}
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
                      nodeId={`${index}.1`}
                      labelText={`Components (${model.components ? model.components.length : 0})`}
                    >
                      {model.components &&
                        model.components.map((component, subIndex) => {
                          return (
                            <StyledTreeItem
                              key={subIndex}
                              nodeId={`${index}.1.${subIndex}`}
                              check
                              labelText={component.displayName}
                              newId={`${index}.1.${subIndex}`}
                              onClick={() => {
                                setShow((prevShow) => {
                                  const { components } = prevShow;
                                  const compIndex = components.findIndex(
                                    (item) => item === component,
                                  );
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
                          );
                        })}
                    </StyledTreeItem>
                    <StyledTreeItem
                      nodeId={`${index}.2`}
                      labelText={`Relationships (${
                        model.relationships ? model.relationships.length : 0
                      })`}
                    >
                      {model.relationships &&
                        model.relationships.map((relationship, subIndex) => (
                          <StyledTreeItem
                            key={subIndex}
                            nodeId={`${index}.2.${subIndex}`}
                            check
                            labelText={relationship.displayhostname}
                            onClick={() => {
                              setShow((prevShow) => {
                                const { relationships } = prevShow;
                                const relaIndex = relationships.findIndex(
                                  (item) => item === relationship,
                                );
                                if (relaIndex !== -1) {
                                  return {
                                    ...prevShow,
                                    model: model,
                                    relationships: relationships.filter(
                                      (item) => item !== relationship,
                                    ),
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
              })}
            </TreeView>
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
            <TreeView
              aria-label="controlled"
              defaultExpanded={['3']}
              defaultCollapseIcon={<MinusSquare />}
              defaultExpandIcon={<PlusSquare />}
              defaultEndIcon={<DotSquare />}
              onNodeToggle={handleToggle}
              multiSelect
              expanded={expanded}
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
                        <StyledTreeItem
                          key={index}
                          nodeId={`${registrant.id}.1.${index}`}
                          check
                          labelText={model.displayName}
                          newId={`${registrant.id}.1.${index}`}
                          onClick={() => {
                            setShow({
                              model: model,
                              components: [],
                              relationships: [],
                            });
                          }}
                        >
                          <div id="scrollElement" onScroll={handleScroll(COMPONENTS)}>
                            <StyledTreeItem
                              nodeId={`${index}.1`}
                              labelText={`Components (${
                                model.components ? model.components.length : 0
                              })`}
                            >
                              {model.components &&
                                model.components.map((component, subIndex) => {
                                  return (
                                    <StyledTreeItem
                                      key={subIndex}
                                      nodeId={`${index + 1}.1.${subIndex}`}
                                      check
                                      labelText={component.displayName}
                                      onClick={() => {
                                        setShow((prevShow) => {
                                          const { components } = prevShow;
                                          const compIndex = components.findIndex(
                                            (item) => item === component,
                                          );
                                          if (compIndex !== -1) {
                                            return {
                                              ...prevShow,
                                              model: model,
                                              components: components.filter(
                                                (item) => item !== component,
                                              ),
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
                                  );
                                })}
                            </StyledTreeItem>
                          </div>

                          <div id="scrollElement" onScroll={handleScroll(RELATIONSHIPS)}>
                            <StyledTreeItem
                              nodeId={`${index}.2`}
                              labelText={`Relationships (${
                                model.relationships ? model.relationships.length : 0
                              })`}
                            >
                              {model.relationships &&
                                model.relationships.map((relationship, subIndex) => (
                                  <StyledTreeItem
                                    key={subIndex}
                                    nodeId={`${index + 2}.2.${subIndex}`}
                                    check
                                    labelText={relationship.displayhostname}
                                    onClick={() => {
                                      setShow((prevShow) => {
                                        const { relationships } = prevShow;
                                        const relaIndex = relationships.findIndex(
                                          (item) => item === relationship,
                                        );
                                        if (relaIndex !== -1) {
                                          return {
                                            ...prevShow,
                                            model: model,
                                            relationships: relationships.filter(
                                              (item) => item !== relationship,
                                            ),
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
                          </div>
                        </StyledTreeItem>
                      ))}
                    </StyledTreeItem>
                  </div>
                </StyledTreeItem>
              ))}
            </TreeView>
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
        />
      )}
    </div>
  );
};

export default MesheryTreeView;
