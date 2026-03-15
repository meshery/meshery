import * as React from 'react';
import {
  ClickAwayListener,
  IconButton,
  MenuItem,
  MenuList,
  Paper,
  Popper,
  Tooltip,
} from '@sistent/sistent';
import MoreVertIcon from '@mui/icons-material/MoreVert';

const ActionPopover = ({
  actions = [],
}: {
  actions?: Array<{
    disabled?: boolean;
    onClick: (_event: any) => void;
    icon?: React.ReactNode;
    label?: string;
  }>;
}) => {
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

      <ClickAwayListener onClickAway={handleClose}>
        <Popper
          sx={{ zIndex: 1 }}
          open={open}
          anchorEl={anchorRef.current}
          {...({
            anchorOrigin: { vertical: 'bottom', horizontal: 'left' },
            transformOrigin: { vertical: 'top', horizontal: 'left' },
          } as any)}
          placement="bottom-start"
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
