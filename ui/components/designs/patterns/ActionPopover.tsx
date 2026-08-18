import * as React from 'react';
import {
  ClickAwayListener,
  IconButton,
  MenuItem,
  MenuList,
  MoreVertIcon,
  Paper,
  Popper,
  Tooltip,
} from '@sistent/sistent';

const ActionPopover = ({ actions = [] }) => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);

  const handleToggle = (event) => {
    event.stopPropagation();
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }
    setOpen(false);
  };

  return (
    <>
      <div ref={anchorRef} style={{ display: 'inline-flex' }}>
        <Tooltip title="Actions" placement="top">
          <IconButton
            size="small"
            onClick={handleToggle}
            aria-label="Open actions menu"
            aria-haspopup="menu"
            aria-expanded={open ? 'true' : undefined}
          >
            <MoreVertIcon />
          </IconButton>
        </Tooltip>
      </div>

      <ClickAwayListener mouseEvent="onMouseDown" onClickAway={handleClose}>
        <Popper
          sx={{ zIndex: 1 }}
          open={open}
          anchorEl={anchorRef.current}
          placement="bottom-start"
        >
          <Paper sx={{ minWidth: 160 }}>
            <MenuList autoFocusItem>
              {actions.map((action, index) => (
                <MenuItem
                  disabled={action.disabled}
                  permissionKey={action.permissionKey}
                  key={index}
                  onClick={(event) => {
                    event.stopPropagation();
                    action.onClick(event);
                    setOpen(false);
                  }}
                  sx={{
                    minWidth: 160,
                    pr: 4.5,
                    display: 'flex',
                    alignItems: 'center',
                  }}
                >
                  <div style={{ marginRight: '0.5rem', display: 'flex', alignItems: 'center' }}>
                    {action.icon}
                  </div>
                  {action.label}
                </MenuItem>
              ))}
            </MenuList>
          </Paper>
        </Popper>
      </ClickAwayListener>
    </>
  );
};

export default ActionPopover;
