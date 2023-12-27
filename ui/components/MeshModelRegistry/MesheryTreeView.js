import React, { useEffect, useState, useRef } from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { Box, Typography, IconButton, FormControlLabel, Switch, useTheme } from '@material-ui/core';
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

const StyledTreeItem = React.forwardRef(function StyledTreeItem(props, ref) {
  const [checked, setChecked] = useState(false);
  const [hover, setHover] = useState(false);
  const { check, labelText, root, search, setSearchText, ...other } = props;
  const theme = useTheme();
  const { width } = useWindowDimensions();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  return (
    <StyledTreeItemRoot
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      root={root}
      lineColor={theme.palette.secondary.text}
      label={
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            py: check ? 0.5 : search ? 0.2 : 1.5,
            px: 0,
          }}
        >
          {width < 1370 && isSearchExpanded ? null : (
            <div
              style={{
                display: 'flex',
                flexDirection: 'row',
                alignItems: 'center',
              }}
            >
              <Typography variant={'body'} style={{ color: `${root}` }}>
                {labelText}
              </Typography>
            </div>
          )}

          {check && (
            <Checkbox
              onClick={() => setChecked((prevcheck) => !prevcheck)}
              size="small"
              checked={checked}
              sx={{
                visibility: hover || checked ? 'visible' : 'hidden',
                color: '#00B39F',
                '&.Mui-checked': {
                  color: '#00B39F',
                },
              }}
            />
          )}
          {search && (
            <SearchBar
              onSearch={debounce((value) => setSearchText(value), 200)}
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
              placeholder="Search"
            />
          )}
        </Box>
      }
      {...other}
      ref={ref}
    />
  );
});

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
    if (div.scrollTop >= div.scrollHeight - div.clientHeight - 10) {
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
            style={{ overflowY: 'auto', height: '27rem' }}
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
              {data.map((model, index) => (
                <StyledTreeItem
                  key={index}
                  top
                  nodeId={index}
                  check
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
                      model.components.map((component, subIndex) => (
                        <StyledTreeItem
                          key={subIndex}
                          nodeId={`${index}.1.${subIndex}`}
                          check
                          labelText={component.displayName}
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
              ))}
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
              {data?.map((registrant) => (
                <StyledTreeItem
                  key={registrant.id}
                  nodeId={0}
                  top
                  labelText={registrant.hostname}
                  onClick={() => setRegi(registrant)}
                >
                  <div>
                    <StyledTreeItem
                      nodeId={1}
                      labelText={`Models (${registrant?.summary?.models})`}
                    >
                      {registrant?.models.map((model, index) => (
                        <StyledTreeItem
                          key={index}
                          nodeId={index + 2}
                          check
                          labelText={model.displayName}
                          onClick={() => {
                            setRegi(registrant);
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
                                      nodeId={`${index + 2}.1.${subIndex}`}
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
      )}
      {view === RELATIONSHIPS && (
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
      )}
    </div>
  );
};

export default MesheryTreeView;
