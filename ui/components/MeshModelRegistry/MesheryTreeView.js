import React, { useEffect, useState } from 'react';
import { TreeView } from '@mui/x-tree-view/TreeView';
import {
  Box,
  Typography,
  IconButton,
  FormControlLabel,
  Switch,
  SvgIcon,
  useTheme,
} from '@material-ui/core';
import Checkbox from '@mui/material/Checkbox';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import SearchBar from '../../utils/custom-search';
import { StyledTreeItemRoot } from './MeshModel.style';

function MinusSquare(props) {
  return (
    <SvgIcon fontSize="inherit" style={{ width: 14, height: 14 }} {...props}>
      {/* tslint:disable-next-line: max-line-length */}
      <path d="M22.047 22.074v0 0-20.147 0h-20.12v0 20.147 0h20.12zM22.047 24h-20.12q-.803 0-1.365-.562t-.562-1.365v-20.147q0-.776.562-1.351t1.365-.575h20.147q.776 0 1.351.575t.575 1.351v20.147q0 .803-.575 1.365t-1.378.562v0zM17.873 11.023h-11.826q-.375 0-.669.281t-.294.682v0q0 .401.294 .682t.669.281h11.826q.375 0 .669-.281t.294-.682v0q0-.401-.294-.682t-.669-.281z" />
    </SvgIcon>
  );
}

function PlusSquare(props) {
  return (
    <SvgIcon fontSize="inherit" style={{ width: 14, height: 14 }} {...props}>
      {/* tslint:disable-next-line: max-line-length */}
      <path d="M22.047 22.074v0 0-20.147 0h-20.12v0 20.147 0h20.12zM22.047 24h-20.12q-.803 0-1.365-.562t-.562-1.365v-20.147q0-.776.562-1.351t1.365-.575h20.147q.776 0 1.351.575t.575 1.351v20.147q0 .803-.575 1.365t-1.378.562v0zM17.873 12.977h-4.923v4.896q0 .401-.281.682t-.682.281v0q-.375 0-.669-.281t-.294-.682v-4.896h-4.923q-.401 0-.682-.294t-.281-.669v0q0-.401.281-.682t.682-.281h4.923v-4.896q0-.401.294-.682t.669-.281v0q.401 0 .682.281t.281.682v4.896h4.923q.401 0 .682.281t.281.682v0q0 .375-.281.669t-.682.294z" />
    </SvgIcon>
  );
}

function CloseSquare(props) {
  return (
    <SvgIcon
      className="close"
      fontSize="inherit"
      style={{ width: 14, height: 14, fillOpacity: 0.5 }}
      {...props}
    >
      {/* tslint:disable-next-line: max-line-length */}
      <path d="M17.485 17.512q-.281.281-.682.281t-.696-.268l-4.12-4.147-4.12 4.147q-.294.268-.696.268t-.682-.281-.281-.682.294-.669l4.12-4.147-4.12-4.147q-.294-.268-.294-.669t.281-.682.682-.281.696 .268l4.12 4.147 4.12-4.147q.294-.268.696-.268t.682.281 .281.669-.294.682l-4.12 4.147 4.12 4.147q.294.268 .294.669t-.281.682zM22.047 22.074v0 0-20.147 0h-20.12v0 20.147 0h20.12zM22.047 24h-20.12q-.803 0-1.365-.562t-.562-1.365v-20.147q0-.776.562-1.351t1.365-.575h20.147q.776 0 1.351.575t.575 1.351v20.147q0 .803-.575 1.365t-1.378.562v0z" />
    </SvgIcon>
  );
}
const StyledTreeItem = React.forwardRef(function StyledTreeItem(props, ref) {
  const [checked, setChecked] = useState(false);
  const [hover, setHover] = useState(false);
  const { check, labelText, root, setSearchText, ...other } = props;
  const theme = useTheme();

  return (
    <StyledTreeItemRoot
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      // className={`${!top ? StyleClass.line : ''}`}
      root={root}
      lineColor={theme.palette.secondary.text}
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
            <Typography variant={'body'} style={{ color: `${root}` }}>
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
    const div = event.target;
    if (div.scrollTop >= div.scrollHeight - div.clientHeight - 2) {
      setPage((prevPage) => ({
        ...prevPage, // Keep the current values for other keys
        [scrollingView]: prevPage[scrollingView] + 1, // Increment the specific key based on the view
      }));
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
                <PlusSquare />
              </IconButton>

              <IconButton
                onClick={() => setExpanded([])}
                style={{ marginRight: '4px' }}
                size="large"
              >
                <MinusSquare />
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
              defaultCollapseIcon={<MinusSquare />}
              defaultExpandIcon={<PlusSquare />}
              defaultEndIcon={<CloseSquare />}
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
              marginRight: '0.9rem',
            }}
          >
            <div style={{ display: 'flex', flexDirection: 'row' }}>
              <IconButton onClick={expandAll} size="large">
                <PlusSquare />
              </IconButton>

              <IconButton
                onClick={() => setExpanded([])}
                style={{ marginRight: '4px' }}
                size="large"
              >
                <MinusSquare />
              </IconButton>
            </div>
          </div>
          <div style={{ overflowY: 'auto', height: '27rem' }}>
            <TreeView
              aria-label="controlled"
              defaultExpanded={['3']}
              defaultCollapseIcon={<MinusSquare />}
              defaultExpandIcon={<PlusSquare />}
              defaultEndIcon={<CloseSquare />}
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
          defaultCollapseIcon={<MinusSquare />}
          defaultExpandIcon={<PlusSquare />}
          defaultEndIcon={<CloseSquare />}
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
          defaultCollapseIcon={<MinusSquare />}
          defaultExpandIcon={<PlusSquare />}
          defaultEndIcon={<CloseSquare />}
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
