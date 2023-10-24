import React, { useEffect, useState } from 'react';
import { styled, useTheme } from '@mui/material/styles';
import { TreeView } from '@mui/x-tree-view/TreeView';
import { TreeItem, treeItemClasses } from '@mui/x-tree-view/TreeItem';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Box, Typography, Button } from '@material-ui/core';
import Checkbox from '@mui/material/Checkbox';
import {
  OVERVIEW,
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  POLICIES,
  REGISTRANTS,
} from '../constants/navigator';

const StyledTreeItemRoot = styled(TreeItem)(({ theme }) => ({
  color: theme.palette.type === 'dark' ? 'white' : 'black',
  [`& .${treeItemClasses.content}`]: {
    fontWeight: theme.typography.fontWeightMedium,
    borderLeft: '10px solid transparent',
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
    },
    '&:hover': {
      backgroundColor: theme.palette.action.hover,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: `var(--tree-view-bg-color, ${theme.palette.action.selected})`,
      borderLeft: '10px solid #00B39F',
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
  const { labelIcon: LabelIcon, check, root, labelText, ...other } = props;

  return (
    <StyledTreeItemRoot
      label={
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            p: !check && !root ? 1.2 : 0.2,
            pr: 0,
            borderBottom: check && '1px solid #DDDDDD',
          }}
          onClick={() => setChecked((prevcheck) => !prevcheck)}
        >
          <Typography variant="body2" sx={{ fontWeight: 'inherit', flexGrow: 1 }}>
            {labelText}
          </Typography>
          {check && (
            <Checkbox
              size="small"
              checked={checked}
              sx={{
                color: '#00B39F',
                '&.Mui-checked': {
                  color: '#00B39F',
                },
              }}
            />
          )}
          {root && <Button variant="text">Collapse All</Button>}
        </Box>
      }
      {...other}
      ref={ref}
    />
  );
});

const MesheryTreeView = ({ data, view, show, setShow, comp, setComp, rela, setRela }) => {
  const addComponent = (component) => {
    // Check if the model is already in the comp array
    if (comp.some((item) => item === component)) {
      // If it exists, remove it from comp
      setComp((prevComp) => prevComp.filter((item) => item !== component));
    } else {
      // If it doesn't exist, add it to comp
      setComp((prevComp) => [...prevComp, component]);
    }
  };
  const addRelationship = (relationship) => {
    // Check if the model is already in the comp array
    if (rela.some((item) => item === relationship)) {
      // If it exists, remove it from comp
      setRela((prevRela) => prevRela.filter((item) => item !== relationship));
    } else {
      // If it doesn't exist, add it to comp
      setRela((prevRela) => [...prevRela, relationship]);
    }
  };

  return (
    <TreeView
      aria-label="gmail"
      defaultExpanded={['3']}
      defaultCollapseIcon={<ExpandMoreIcon />}
      defaultExpandIcon={<ChevronRightIcon />}
      defaultEndIcon={<div style={{ width: 24 }} />}
    >
      {view === MODELS &&
        data.map((model, index) => (
          <StyledTreeItem
            key={index}
            nodeId={index}
            root
            labelText={model.displayName}
            onClick={() => {
              setShow(data[index]), setComp((prevC) => []);
              setRela((prevR) => []);
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
                    onClick={() => addComponent(component)}
                  />
                ))}
            </StyledTreeItem>
            <StyledTreeItem
              nodeId={`${index}.2`}
              labelText={`Relationships (${model.relationships ? model.relationships.length : 0})`}
            >
              {model.relationships &&
                model.relationships.map((relationship, subIndex) => (
                  <StyledTreeItem
                    key={subIndex}
                    nodeId={`${index}.2.${subIndex}`}
                    check
                    labelText={relationship.displayhostname}
                    onClick={() => addRelationship(relationship)}
                  />
                ))}
            </StyledTreeItem>
          </StyledTreeItem>
        ))}
      {view === COMPONENTS &&
        data.map((component, index) => (
          <StyledTreeItem
            key={index}
            check
            nodeId={index}
            labelText={component.displayName}
            onClick={() => {
              addComponent(data[index]);
              setShow(data[index]?.model);
            }}
          />
        ))}
      {view === RELATIONSHIPS &&
        data.map((relationship, index) => (
          <StyledTreeItem
            key={index}
            check
            nodeId={index}
            labelText={relationship.displayhostname}
            onClick={() => {
              addRelationship(data[index]);
              setShow(data[index]?.model);
            }}
          />
        ))}
    </TreeView>
  );
};

export default MesheryTreeView;
