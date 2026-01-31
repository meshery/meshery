import * as React from 'react';
import {
  Button,
  ButtonGroup,
  Paper,
  Popper,
  MenuItem,
  MenuList,
  ClickAwayListener,
  CustomTooltip,
} from '@sistent/sistent';
import ArrowDropDownIcon from '@mui/icons-material/ArrowDropDown';

export default function ActionButton({ options }) {
  const [open, setOpen] = React.useState(false);
  const [interactiveMode, setInteractiveMode] = React.useState(true);
  const anchorRef = React.useRef(null);

  const handleMenuItemClick = (_event) => {
    setOpen(false);
  };

  const handleActionButtonClick = (event) => {
    event.stopPropagation();
    setInteractiveMode(true);
    setOpen((prevOpen) => !prevOpen);
  };

  const handleArrowClick = (event) => {
    event.stopPropagation();
    setInteractiveMode(false);
    setOpen((prevOpen) => !prevOpen);
  };

  const handleClose = (event) => {
    if (anchorRef.current && anchorRef.current.contains(event.target)) {
      return;
    }

    setOpen(false);
  };

  return (
    <React.Fragment>
      <ButtonGroup
        variant="outlined"
        style={{ boxShadow: 'none' }}
        ref={anchorRef}
        aria-label="Button group with a nested menu"
      >
        <CustomTooltip title="Invoke actions interactively" placement="top">
          <Button
            sx={{
              padding: '6px 9px',
              borderRadius: '8px',
            }}
            onClick={handleActionButtonClick}
            variant="outlined"
          >
            Action
          </Button>
        </CustomTooltip>
        <CustomTooltip title="Invoke actions in single click" placement="top">
          <Button
            sx={{
              padding: '6px 9px',
              borderRadius: '8px',
            }}
            size="small"
            onClick={handleArrowClick}
            variant="outlined"
            data-testid="action-btn-toggle"
          >
            <ArrowDropDownIcon />
          </Button>
        </CustomTooltip>
      </ButtonGroup>
      <Popper
        sx={{
          zIndex: 1,
        }}
        open={open}
        anchorEl={anchorRef.current}
        placement="bottom-end"
      >
        <Paper>
          <ClickAwayListener onClickAway={handleClose}>
            <MenuList id="split-button-menu" autoFocusItem>
              {options.map((option, index) => (
                <MenuItem
                  key={option.label}
                  data-testid={`action-btn-option-${option.label}`}
                  disabled={option.disabled}
                  onClick={(event) => {
                    handleMenuItemClick(event);
                    if (interactiveMode) {
                      option.onClick(event, index);
                    } else {
                      option.onDirectClick?.(event, index) || option.onClick(event, index);
                    }
                  }}
                >
                  <CustomTooltip
                    title={option.label}
                    placement="left"
                    disableHoverListener={interactiveMode}
                  >
                    <div style={{ display: 'flex', alignItems: 'center', width: '100%' }}>
                      <div style={{ marginRight: '0.5rem' }}>{option.icon}</div>
                      {interactiveMode && option.label}
                    </div>
                  </CustomTooltip>
                </MenuItem>
              ))}
            </MenuList>
          </ClickAwayListener>
        </Paper>
      </Popper>
    </React.Fragment>
  );
}
