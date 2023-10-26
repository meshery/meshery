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
  ENDPOINTURL,
} from '../constants/navigator';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import KeyboardArrowUpIcon from '@mui/icons-material/KeyboardArrowUp';

const StyledTreeItemRoot = styled(TreeItem)(({ theme }) => ({
  color: theme.palette.type === 'dark' ? 'white' : 'black',
  [`& .${treeItemClasses.content}`]: {
    flexDirection: 'row-reverse',
    fontWeight: theme.typography.fontWeightMedium,
    borderLeft: '10px solid transparent',
    '&.Mui-expanded': {
      fontWeight: theme.typography.fontWeightRegular,
    },
    '&:hover': {
      backgroundColor: theme.palette.action.hover,
    },
    '&.Mui-focused, &.Mui-selected, &.Mui-selected.Mui-focused': {
      backgroundColor: '#ebebeb',
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
  const [hover, setHover] = useState(false);
  const { labelIcon: LabelIcon, check, root, labelText, avatar, ...other } = props;

  return (
    <StyledTreeItemRoot
      sx={{
        backgroundColor: checked && check && '#ebebeb',
      }}
      label={
        <Box
          onMouseEnter={() => setHover(true)}
          onMouseLeave={() => setHover(false)}
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            p: !check && !root ? 1.5 : 0.5,
            pr: 0,
            borderBottom: check && '1px solid #DDDDDD',
          }}
        >
          <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'center' }}>
            {avatar && (
              <Avatar alt="Image" src={avatar} style={{ marginRight: '6px' }}>
                IM
              </Avatar>
            )}
            <Typography variant="body2" sx={{ fontWeight: 'inherit', flexGrow: 1 }}>
              {labelText}
            </Typography>
          </div>
          {check && (
            <Checkbox
              size="small"
              checked={checked}
              onClick={() => setChecked((prevcheck) => !prevcheck)}
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

const MesheryTreeView = ({
  data,
  view,
  show,
  setShow,
  comp,
  setComp,
  rela,
  setRela,
  regi,
  setRegi,
}) => {
  const [expanded, setExpanded] = React.useState([]);

  useEffect(() => {
    if (view === RELATIONSHIPS || view === COMPONENTS) {
      expandAll();
    } else if (view === MODELS) {
      setExpanded([]);
    }
    console.log(data);
  }, [view]);

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

  const addComponent = (component) => {
    if (comp.some((item) => item === component)) {
      setComp((prevComp) => prevComp.filter((item) => item !== component));
    } else {
      setComp((prevComp) => [...prevComp, component]);
    }
  };
  const addRelationship = (relationship) => {
    if (rela.some((item) => item === relationship)) {
      setRela((prevRela) => prevRela.filter((item) => item !== relationship));
    } else {
      setRela((prevRela) => [...prevRela, relationship]);
    }
  };

  const addRegistrant = (registrant) => {
    if (regi.some((item) => item === registrant)) {
      setRegi((prevRegi) => prevRegi.filter((item) => item !== registrant));
    } else {
      setRegi((prevRegi) => [...prevRegi, registrant]);
    }
  };

  return (
    <div>
      <div
        style={{ borderBottom: '1px solid #d2d3d4', display: 'flex', justifyContent: 'flex-end' }}
      >
        <Button startIcon={<KeyboardArrowDownIcon />} onClick={expandAll}>
          Expand All
        </Button>
        <div style={{ backgroundColor: '#d2d3d4', height: '33px', width: '1px' }}></div>
        <Button startIcon={<KeyboardArrowUpIcon />} onClick={() => setExpanded([])}>
          Collapse All
        </Button>
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
          {view !== REGISTRANTS &&
            data.map((model, index) => (
              <StyledTreeItem
                key={index}
                nodeId={index}
                root
                avatar={ENDPOINTURL + model.metadata?.svgColor}
                labelText={model.displayName}
                onClick={() => {
                  setShow(data[index]), setComp((prevC) => []);
                  setRela((prevR) => []);
                }}
              >
                {view !== RELATIONSHIPS && (
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
                            addComponent(component);
                            setShow(data[index]);
                          }}
                        />
                      ))}
                  </StyledTreeItem>
                )}
                {view !== COMPONENTS && (
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
                            addRelationship(relationship);
                            setShow(data[index]);
                          }}
                        />
                      ))}
                  </StyledTreeItem>
                )}
              </StyledTreeItem>
            ))}
          {/* {view === COMPONENTS &&
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
          ))} */}
          {view === REGISTRANTS &&
            data.map((registrant, index) => (
              <StyledTreeItem
                key={index}
                check
                nodeId={index}
                labelText={registrant.hostname}
                onClick={() => {
                  addRegistrant(data[index]);
                }}
              />
            ))}
        </TreeView>
      </div>
    </div>
  );
};

export default MesheryTreeView;
