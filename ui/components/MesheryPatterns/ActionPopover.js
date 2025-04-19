import * as React from 'react';
import {
  Button,
  ButtonGroup,
  Paper,
  Popper,
  MenuItem,
  MenuList,
  ClickAwayListener,
} from '@layer5/sistent';
import MoreVertIcon from '@mui/icons-material/MoreVert';

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
      <ButtonGroup
        variant="text"
        style={{ boxShadow: 'none' }}
        ref={anchorRef}
        aria-label="Button group with a nested menu"
      >
        <Button size="small" onClick={handleToggle} variant="text">
          <MoreVertIcon />
        </Button>
      </ButtonGroup>

      <ClickAwayListener onClickAway={handleClose}>
        <Popper
          sx={{ zIndex: 1 }}
          open={open}
          anchorEl={anchorRef.current}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
          transformOrigin={{ vertical: 'top', horizontal: 'left' }}
        >
          <Paper>
            <MenuList autoFocusItem>
              {actions.map((action, index) => (
                <MenuItem
                  disabled={action.disabled}
                  key={index}
                  onClick={(event) => {
                    event.stopPropagation();
                    action.onClick(event);
                    setOpen(false);
                  }}
                >
                  <div style={{ marginRight: '0.5rem' }}>{action.icon}</div>
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
