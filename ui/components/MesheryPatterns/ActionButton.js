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
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';
import { UsesSistent } from '../SistentWrapper';

export default function ActionButton({ options }) {
  const [open, setOpen] = React.useState(false);
  const [dropdownMode, setDropdownMode] = React.useState('withLabel');
  const anchorRef = React.useRef(null);

  const handleMenuItemClick = () => {
    setOpen(false);
  };

  const handleToggle = (event, mode) => {
    event.stopPropagation();
    setDropdownMode(mode);
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }

    setOpen(false);
  };

  return (
    <UsesSistent>
      <React.Fragment>
        <ButtonGroup
          variant="outlined"
          style={{ boxShadow: 'none' }}
          ref={anchorRef}
          aria-label="Button group with a nested menu"
        >
          <Button
            sx={{
              padding: '6px 9px',
              borderRadius: '8px',
            }}
            onClick={(e) => handleToggle(e, 'withLabel')} // Show names + icons
            variant="outlined"
          >
            Action
          </Button>
          <Button
            sx={{
              padding: '6px 9px',
              borderRadius: '8px',
            }}
            size="small"
            onClick={(e) => handleToggle(e, 'icon-only')} // Show icons only
            variant="outlined"
          >
            <ArrowDropDownIcon />
          </Button>
        </ButtonGroup>
        <Popper
          sx={{
            zIndex: 1,
          }}
          open={open}
          anchorEl={anchorRef.current}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'left',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'left',
          }}
        >
          <Paper>
            <ClickAwayListener onClickAway={handleClose}>
              <MenuList id="split-button-menu" autoFocusItem>
                {options.map((option, index) => (
                  <MenuItem
                    disabled={option.disabled}
                    key={option}
                    onClick={(event) => {
                      handleMenuItemClick(event);
                      option.onClick(event, index);
                    }}
                  >
                    {dropdownMode === 'icon-only' ? (
                      <div>{option.icon}</div>
                    ) : (
                      <div style={{ display: 'flex' }}>
                        <div style={{ marginRight: '0.5rem' }}>{option.icon}</div>
                        {option.label}
                      </div>
                    )}
                  </MenuItem>
                ))}
              </MenuList>
            </ClickAwayListener>
          </Paper>
        </Popper>
      </React.Fragment>
    </UsesSistent>
  );
}
