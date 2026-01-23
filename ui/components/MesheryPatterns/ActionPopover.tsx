import * as React from 'react';
import {
  Button,
  ButtonGroup,
  Paper,
  Popper,
  MenuItem,
  MenuList,
  ClickAwayListener,
} from '@sistent/sistent';
import MoreVertIcon from '@mui/icons-material/MoreVert';

const ActionPopover = ({ actions = [] }: { actions?: Array<{ disabled?: boolean; onClick: (event: any) => void; icon?: React.ReactNode; label?: string }> }) => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef<HTMLDivElement>(null);

  const handleToggle = (_event: any) => {
    _event.stopPropagation();
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event: any) => {
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
          {...({
            anchorOrigin: { vertical: 'bottom', horizontal: 'left' },
            transformOrigin: { vertical: 'top', horizontal: 'left' },
          } as any)}
        >
          <Paper>
            <MenuList autoFocusItem>
              {actions.map((action, index) => (
                <MenuItem
                  {...(action.disabled !== undefined ? { disabled: action.disabled } : {})}
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
