import React from 'react';
import { styled, Menu, MenuItem, Button } from '@material-ui/core';
import { useState } from 'react';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { iconMedium } from '../../css/icons.styles';

export const VIEW_VISIBILITY = {
  PUBLIC: 'public',
  PRIVATE: 'private',
};

const StyledMenu = styled(Menu)(({ theme }) => ({
  '& .MuiPaper-root': {
    backgroundColor: theme.palette.secondary.menuBg,
    color: theme.palette.secondary.titleText,
    border: `1px solid ${theme.palette.secondary.border}`,
    borderRadius: '0.25rem',
    padding: '0rem',
  },
  '& .MuiMenuItem-root': {
    fontSize: '.9rem',
    padding: '0.5rem',
    '&:hover': {
      backgroundColor: theme.palette.secondary.menuSelectedBg,
    },
  },
  //selected
  '& .Mui-selected': {
    backgroundColor: theme.palette.secondary.menuSelectedBg,
  },
  '& .MuiList-padding': {
    padding: '0px',
  },
}));
const StyledMenuWrapper = styled('div')(({ theme }) => ({
  paddingBlock: '0.15rem',
  paddingInline: '0.5rem',
  borderRadius: '0.25rem',
  border: '1px solid #666',
  background: theme.palette.secondary.menuBg,
  textTransform: 'uppercase',
  color: theme.palette.secondary.text,
  display: 'flex',
  gap: '0.5rem',
}));
const StyledMenuItem = styled(MenuItem)(({ theme }) => ({
  textTransform: 'capitalize',
  color: theme.palette.secondary.iconMain,
}));
export const VisibilityMenu = ({ value, onChange, options, enabled }) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);
  const close = () => {
    setAnchorEl(null);
  };
  const handleOpen = (e) => {
    if (!enabled) return;
    setAnchorEl(e.currentTarget);
  };
  const handleChange = (value) => {
    onChange(value);
    close();
  };
  return (
    <>
      <Button
        disabled={!enabled}
        onClick={handleOpen}
        style={{
          padding: '0px',
        }}
      >
        <StyledMenuWrapper>
          <span>{value}</span>
          <ArrowDropDownIcon {...iconMedium} />
        </StyledMenuWrapper>
      </Button>

      <StyledMenu
        open={open}
        onClose={close}
        anchorEl={anchorEl}
        anchorReference="anchorPosition"
        anchorPosition={{
          top: anchorEl?.getBoundingClientRect().bottom + 5 ?? 0,
          left: anchorEl?.getBoundingClientRect().left + 5 ?? 0,
        }}
      >
        {options.map(([visibility, Icon], index) => (
          <StyledMenuItem key={index} onClick={() => handleChange(visibility)}>
            <Icon
              width={16}
              height={16}
              style={{
                marginRight: '0.5rem',
              }}
            />
            {visibility}
          </StyledMenuItem>
        ))}
      </StyledMenu>
    </>
  );
};
