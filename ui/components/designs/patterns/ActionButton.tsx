import * as React from 'react';
import {
  ArrowDropDownIcon,
  Button,
  ButtonGroup,
  ClickAwayListener,
  MenuItem,
  MenuList,
  Paper,
  Popper,
  CustomTooltip,
  Box,
} from '@sistent/sistent';

export default function ActionButton({
  defaultActionClick,
  options = [],
  permissionKey,
  permissionAction,
  label = 'Action',
}) {
  const [open, setOpen] = React.useState(false);
  const [interactiveMode, setInteractiveMode] = React.useState(true);
  const anchorRef = React.useRef(null);

  const handleMenuItemClick = () => {
    setOpen(false);
  };

  const handleActionButtonClick = (event) => {
    event.stopPropagation();
    if (defaultActionClick) {
      defaultActionClick(event);
      return;
    }
    if (open && interactiveMode) {
      setOpen(false);
    } else {
      setInteractiveMode(true);
      setOpen(true);
    }
  };

  const handleArrowClick = (event) => {
    event.stopPropagation();
    if (open && !interactiveMode) {
      setOpen(false);
    } else {
      setInteractiveMode(false);
      setOpen(true);
    }
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
            permissionKey={permissionKey}
            permissionAction={permissionAction}
          >
            {label}
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
          zIndex: 1300,
        }}
        open={open}
        anchorEl={anchorRef.current}
        placement={interactiveMode ? 'bottom-start' : 'bottom-end'}
      >
        <Paper>
          <ClickAwayListener onClickAway={handleClose}>
            <MenuList id="split-button-menu" autoFocusItem>
              {options.map((option, index) => (
                <MenuItem
                  data-testid={option['data-testid'] || `action-btn-option-${option.label}`}
                  disabled={option.disabled}
                  permissionKey={option.permissionKey}
                  permissionAction={option.permissionAction}
                  key={option.label || index}
                  onClick={(event) => {
                    handleMenuItemClick();
                    if (interactiveMode) {
                      option.onClick?.(event, index);
                    } else if (option.onDirectClick) {
                      option.onDirectClick(event, index);
                    } else {
                      option.onClick?.(event, index);
                    }
                  }}
                  sx={{
                    justifyContent: interactiveMode ? 'flex-start' : 'center',
                    minWidth: interactiveMode ? '140px' : 'auto',
                    padding: interactiveMode ? '6px 16px' : '6px 12px',
                  }}
                >
                  <CustomTooltip
                    title={option.label}
                    placement="left"
                    disableHoverListener={interactiveMode}
                  >
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: interactiveMode ? 'flex-start' : 'center',
                        width: '100%',
                      }}
                    >
                      <Box
                        sx={{
                          display: 'flex',
                          alignItems: 'center',
                          marginRight: interactiveMode ? '0.5rem' : 0,
                        }}
                      >
                        {option.icon}
                      </Box>
                      {interactiveMode && option.label}
                    </Box>
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
