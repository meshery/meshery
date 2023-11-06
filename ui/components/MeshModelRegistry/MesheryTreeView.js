import React, { useEffect, useState } from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';
import { Box, Typography, IconButton, FormControlLabel, Switch } from '@material-ui/core';
import Checkbox from '@mui/material/Checkbox';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import SearchBar from '../../utils/custom-search';
import ExpandAllIcon from '../../assets/icons/expand_all';
import CollapseAllIcon from '../../assets/icons/collapse_all';
import ExpandMoreIcon from '../../assets/icons/expand_more';
import ChevronRightIcon from '../../assets/icons/chevron_right';
import { alpha, styled } from '@mui/material/styles';

const StyledTreeItemRoot = styled(TreeItem)(({ theme }) => ({
  [`& .${treeItemClasses.content}`]: {
    fontWeight: theme.typography.fontWeightMedium,
    borderRadius: '4px',
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
    },
    '&:hover': {
      backgroundColor: `transparent`,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: `#00bfa030`,
      borderLeft: '2px solid #00bfa0',
    },
    [`& .${treeItemClasses.label}`]: {
      fontWeight: 'inherit',
    },
  },
  [`& .${treeItemClasses.group}`]: {
    paddingRight: '0',
    borderLeft: `1px solid ${alpha(theme.palette.text.primary, 0.4)}`,
    // [`& .${treeItemClasses.content}`]: {
    //   paddingLeft: theme.spacing(2),
    // },
  },
}));

const StyledTreeItem = React.forwardRef(function StyledTreeItem(props, ref) {
  const [checked, setChecked] = useState(false);
  const [hover, setHover] = useState(false);
  const { check, labelText, root, setSearchText, ...other } = props;

  return (
    <StyledTreeItemRoot
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      sx={
        {
          // borderLeft:!root && '2px solid #00bfa0'
        }
      }
      label={
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            py: check ? 0.5 : root ? 0.2 : 1.5,
            px: 0,
          }}
        >
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'center',
            }}
          >
            <Typography
              variant={hover || root ? 'body' : 'body2'}
              style={{ color: `${root && '#005711'}` }}
            >
              {labelText}
            </Typography>
          </div>
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
          {root && (
            <SearchBar
              onSearch={(value) => {
                setSearchText(value);
              }}
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

  useEffect(() => {
    // console.log(data);
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

  const handleScroll = (scrollingView) => () => {
    console.log('event', event);
    const div = event.target;
    // console.log(data);
    console.log('view: ', scrollingView);
    if (div.scrollTop >= div.scrollHeight - div.clientHeight - 2) {
      setPage((prevPage) => ({
        ...prevPage, // Keep the current values for other keys
        [scrollingView]: prevPage[scrollingView] + 1, // Increment the specific key based on the view
      }));
      // setPage((prev) => prev + 1);
    }
  };

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
    <div>
      {view === MODELS && (
        <div>
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              justifyContent: 'space-between',
              borderBottom: '1px solid #d2d3d4',
              marginRight: '0.9rem',
            }}
          >
            <div style={{ display: 'flex', flexDirection: 'row' }}>
              <IconButton onClick={expandAll} size="large">
                <ExpandAllIcon />
              </IconButton>
              <div
                style={{
                  backgroundColor: '#d2d3d4',
                  height: '33px',
                  width: '1px',
                  margin: '0 2px',
                }}
              ></div>
              <IconButton
                onClick={() => setExpanded([])}
                style={{ marginRight: '4px' }}
                size="large"
              >
                <CollapseAllIcon />
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
            <div style={{ display: 'flex' }}>
              <SearchBar
                onSearch={(value) => {
                  setSearchText(value);
                }}
                placeholder="Search"
              />
            </div>
          </div>
          <div style={{ overflowY: 'auto', height: '27rem' }} onScroll={handleScroll(MODELS)}>
            <TreeView
              aria-label="controlled"
              defaultExpanded={['3']}
              defaultCollapseIcon={<ExpandMoreIcon />}
              defaultExpandIcon={<ChevronRightIcon />}
              onNodeToggle={handleToggle}
              multiSelect
              expanded={expanded}
            >
              {data.map((model, index) => (
                <StyledTreeItem
                  key={index}
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
              marginRight: '0.9rem',
            }}
          >
            <div style={{ display: 'flex', flexDirection: 'row' }}>
              <IconButton onClick={expandAll} size="large">
                <ExpandAllIcon />
              </IconButton>
              <div
                style={{
                  backgroundColor: '#d2d3d4',
                  height: '33px',
                  width: '1px',
                  margin: '0 2px',
                }}
              ></div>
              <IconButton
                onClick={() => setExpanded([])}
                style={{ marginRight: '4px' }}
                size="large"
              >
                <CollapseAllIcon />
              </IconButton>
            </div>
          </div>
          <div style={{ overflowY: 'auto', height: '27rem' }}>
            <TreeView
              aria-label="controlled"
              defaultExpanded={['3']}
              defaultCollapseIcon={<ExpandMoreIcon />}
              defaultExpandIcon={<ChevronRightIcon />}
              onNodeToggle={handleToggle}
              multiSelect
              expanded={expanded}
            >
              {data?.map((registrant) => (
                <StyledTreeItem
                  nodeId={0}
                  labelText={registrant.hostname}
                  onClick={() => setRegi(registrant)}
                >
                  <div
                    style={{ overflowY: 'auto', height: '27rem' }}
                    onScroll={handleScroll(MODELS)}
                  >
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
                          <div onScroll={handleScroll(COMPONENTS)}>
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

                          <div onScroll={handleScroll(RELATIONSHIPS)}>
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
          defaultCollapseIcon={<ExpandMoreIcon />}
          defaultExpandIcon={<ChevronRightIcon />}
          onNodeToggle={handleToggle}
          onNodeSelect={handleSelect}
          multiSelect
          expanded={expanded}
          selected={selected}
        >
          <StyledTreeItem
            nodeId={0}
            root
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
            <div style={{ overflowY: 'auto', height: '27rem' }} onScroll={handleScroll(COMPONENTS)}>
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
          defaultCollapseIcon={<ExpandMoreIcon />}
          defaultExpandIcon={<ChevronRightIcon />}
          onNodeToggle={handleToggle}
          onNodeSelect={handleSelect}
          multiSelect
          expanded={expanded}
          selected={selected}
        >
          <StyledTreeItem
            nodeId={0}
            root
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
              style={{ overflowY: 'auto', height: '27rem' }}
              onScroll={handleScroll(RELATIONSHIPS)}
            >
              {data.map((relationship, index) => (
                <StyledTreeItem
                  key={index}
                  nodeId={index + 1}
                  check
                  labelText={relationship.kind}
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
