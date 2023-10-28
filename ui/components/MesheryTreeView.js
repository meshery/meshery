import React, { useEffect, useState } from 'react';
import { styled, useTheme } from '@mui/material/styles';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Box, Typography, Button, Avatar } from '@material-ui/core';
import Checkbox from '@mui/material/Checkbox';
import {
  OVERVIEW,
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  POLICIES,
  REGISTRANTS,
} from '../constants/navigator';
import UnfoldMoreIcon from '@mui/icons-material/UnfoldMore';
import UnfoldLessIcon from '@mui/icons-material/UnfoldLess';
import SearchBar from '../utils/custom-search';

const StyledTreeItemRoot = styled(TreeItem)(({ theme }) => ({
  color: theme.palette.type === 'dark' ? 'white' : 'black',
  [`& .${treeItemClasses.content}`]: {
    fontWeight: theme.typography.fontWeightMedium,
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
      borderBottom: '1px solid #DDDDDD',
    },
    '&:hover': {
      backgroundColor: theme.palette.action.hover,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: `var(--tree-view-bg-color, ${theme.palette.action.selected})`,
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
  const theme = useTheme();
  const [checked, setChecked] = useState(false);
  const [hover, setHover] = useState(false);
  const { labelIcon: LabelIcon, check, labelText, ...other } = props;

  return (
    <StyledTreeItemRoot
      onMouseEnter={() => setHover(true)}
      onMouseLeave={() => setHover(false)}
      sx={{
        backgroundColor: checked && '#ebebeb',
        borderBottom: '1px solid #DDDDDD',
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
                visibility: hover ? 'visible' : 'hidden',
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

const MesheryTreeView = ({ data, view, setShow, setComp, setRela, setRegi, setSearchText }) => {
  const [expanded, setExpanded] = React.useState([]);

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
          justifyContent: view === MODELS ? 'space-between' : 'flex-end',
          borderBottom: '1px solid #d2d3d4',
        }}
      >
        {view === MODELS && (
          <div style={{ display: 'flex', flexDirection: 'row' }}>
            <Button startIcon={<UnfoldMoreIcon />} onClick={expandAll}>
              Expand All
            </Button>
            <div style={{ backgroundColor: '#d2d3d4', height: '33px', width: '1px' }}></div>
            <Button startIcon={<UnfoldLessIcon />} onClick={() => setExpanded([])}>
              Collapse All
            </Button>
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
      <div style={{ overflowY: 'scroll', height: '27rem' }}>
        <TreeView
          aria-label="controlled"
          defaultExpanded={['3']}
          defaultCollapseIcon={<ExpandMoreIcon />}
          defaultExpandIcon={<ChevronRightIcon />}
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
