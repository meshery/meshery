import React, { useEffect, useState } from 'react';
import { styled } from '@mui/material/styles';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Box, Typography, IconButton, FormControlLabel, Switch } from '@material-ui/core';
import Checkbox from '@mui/material/Checkbox';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../constants/navigator';
import UnfoldMoreIcon from '@mui/icons-material/UnfoldMore';
import UnfoldLessIcon from '@mui/icons-material/UnfoldLess';
import SearchBar from '../utils/custom-search';
import ExpandAll from '../public/static/img/expand_all.svg';

const StyledTreeItemRoot = styled(TreeItem)(({ theme }) => ({
  color: theme.palette.type === 'dark' ? 'white' : 'black',
  [`& .${treeItemClasses.content}`]: {
    fontWeight: theme.typography.fontWeightMedium,
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
    },
    '&:hover': {
      backgroundColor: `transparent`,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: `transparent`,
    },
    [`& .${treeItemClasses.label}`]: {
      fontWeight: 'inherit',
      color: 'inherit',
    },
  },
  [`& .${treeItemClasses.group}`]: {
    marginLeft: 0,
    [`& .${treeItemClasses.content}`]: {
      paddingLeft: theme.spacing(2),
    },
  },
}));

const StyledTreeItem = React.forwardRef(function StyledTreeItem(props, ref) {
  const [checked, setChecked] = useState(false);
  const [hover, setHover] = useState(false);
  const { check, labelText, ...other } = props;

  return (
    <StyledTreeItemRoot
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      sx={{
        backgroundColor: checked && '#ebebeb',
      }}
      label={
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            p: !check ? 1.5 : 0.5,
            pr: 0,
          }}
        >
          <div
            style={{
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'center',
              marginLeft: check && '1rem',
            }}
          >
            <Typography variant="body2">{labelText}</Typography>
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

  const handleScroll = () => {
    const div = event.target;
    console.log(data);
    if (div.scrollTop >= div.scrollHeight - div.clientHeight - 2) {
      setPage((prev) => prev + 1);
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

  const handleToggle = (event, nodeIds) => {
    setExpanded(nodeIds);
  };

  return (
    <div>
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: view === MODELS || view === COMPONENTS ? 'space-between' : 'flex-end',
          borderBottom: '1px solid #d2d3d4',
          marginRight: '0.9rem',
        }}
      >
        {view === MODELS && (
          <div style={{ display: 'flex', flexDirection: 'row' }}>
            <IconButton onClick={expandAll} size="large">
              <img src="static/img/expand_all.svg" />
            </IconButton>
            <div
              style={{ backgroundColor: '#d2d3d4', height: '33px', width: '1px', margin: '0 2px' }}
            ></div>
            <IconButton onClick={() => setExpanded([])} style={{ marginRight: '4px' }} size="large">
              <img src="static/img/collapse_all.svg" />
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
        {view === COMPONENTS && (
          <div style={{ display: 'flex', flexDirection: 'row' }}>
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
        <div style={{ display: 'flex' }}>
          <SearchBar
            onSearch={(value) => {
              setSearchText(value);
            }}
            placeholder="Search"
          />
        </div>
      </div>
      <div style={{ overflowY: 'auto', height: '27rem' }} onScroll={handleScroll}>
        <TreeView
          aria-label="controlled"
          defaultExpanded={['3']}
          defaultCollapseIcon={<img src="static/img/expand_more.svg" />}
          defaultExpandIcon={<img src="static/img/chevron_right.svg" />}
          onNodeToggle={handleToggle}
          multiSelect
          expanded={expanded}
        >
          {view === MODELS &&
            data.map((model, index) => (
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
          {view === COMPONENTS &&
            data.map((component, index) => (
              <StyledTreeItem
                key={index}
                nodeId={index}
                check
                labelText={component.displayName}
                onClick={() => {
                  setComp(component);
                }}
              />
            ))}
          {view === RELATIONSHIPS &&
            data.map((relationship, index) => (
              <StyledTreeItem
                key={index}
                nodeId={index}
                check
                labelText={relationship.kind}
                onClick={() => {
                  setRela(relationship);
                }}
              />
            ))}
          {view === REGISTRANTS &&
            data.map((registrant, index) => (
              <StyledTreeItem
                key={index}
                nodeId={index}
                check
                labelText={registrant.hostname}
                onClick={() => {
                  setRegi(registrant);
                }}
              />
            ))}
        </TreeView>
      </div>
    </div>
  );
};

export default MesheryTreeView;
