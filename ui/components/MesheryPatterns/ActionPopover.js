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
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import EditIcon from '@mui/icons-material/Edit';
import CloneIcon from '../../public/static/img/CloneIcon';
import CheckIcon from '@mui/icons-material/Check';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import UndeployIcon from '../../public/static/img/UndeployIcon';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import GetAppIcon from '@mui/icons-material/GetApp';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import PublicIcon from '@mui/icons-material/Public';
import { VISIBILITY } from '@/utils/Enum';
import PatternConfigureIcon from '@/assets/icons/PatternConfigure';

const ActionPopover = ({
  rowData,
  visibility,
  patterns,
  tableMeta,
  handleOpenInConfigurator,
  handleClone,
  openValidateModal,
  openDryRunModal,
  openUndeployModal,
  openDeployModal,
  handleDesignDownloadModal,
  handleInfoModal,
  handleUnpublishModal,
  userCanEdit,
  CAN,
  keys,
}) => {
  const [open, setOpen] = React.useState(false);
  const anchorRef = React.useRef(null);

  const handleMenuItemClick = () => {
    setOpen(false);
  };

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

  const genericClickHandler = (ev, fn) => {
    ev.stopPropagation();
    fn(ev);
  };

  console.log(visibility);

  const options = [
    {
      label: 'Edit',
      icon: <EditIcon fill="currentColor" />,
      onClick: () => handleOpenInConfigurator(rowData.id),
      disabled: !userCanEdit(rowData),
    },
    {
      label: 'Clone',
      icon: <CloneIcon fill="currentColor" />,
      onClick: () => handleClone(rowData.id, rowData.name),
      disabled:
        visibility !== VISIBILITY.PUBLISHED ||
        !CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject),
    },
    {
      label: 'Design',
      icon: <PatternConfigureIcon />,
      onClick: () => handleOpenInConfigurator(patterns[tableMeta.rowIndex].id),
      disabled:
        visibility == VISIBILITY.PUBLISHED ||
        !CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject),
    },
    {
      label: 'Validate Design',
      icon: <CheckIcon />,
      onClick: (event) => openValidateModal(event, rowData.pattern_file, rowData.name, rowData.id),
      disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
    },
    {
      label: 'Dry Run',
      icon: <DryRunIcon data-cy="verify-button" />,
      onClick: (event) => openDryRunModal(event, rowData.pattern_file, rowData.name, rowData.id),
      disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
    },
    {
      label: 'Undeploy',
      icon: <UndeployIcon fill="#F91313" data-cy="undeploy-button" />,
      onClick: (event) => openUndeployModal(event, rowData.pattern_file, rowData.name, rowData.id),
      disabled: !CAN(keys.UNDEPLOY_DESIGN.action, keys.UNDEPLOY_DESIGN.subject),
    },
    {
      label: 'Deploy',
      icon: <DoneAllIcon />,
      onClick: (event) => openDeployModal(event, rowData.pattern_file, rowData.name, rowData.id),
      disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
    },
    {
      label: 'Download',
      icon: <GetAppIcon />,
      onClick: (event) => handleDesignDownloadModal(event, rowData),
      disabled: !CAN(keys.DOWNLOAD_A_DESIGN.action, keys.DOWNLOAD_A_DESIGN.subject),
    },
    {
      label: 'Design Information',
      icon: <InfoOutlinedIcon />,
      onClick: (event) => genericClickHandler(event, () => handleInfoModal(rowData)),
      disabled: !CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject),
    },
    visibility === VISIBILITY.PUBLISHED && {
      label: 'Unpublish',
      icon: <PublicIcon fill="#F91313" />,
      onClick: (event) => handleUnpublishModal(event, rowData),
      disabled: !CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject),
    },
  ].filter(Boolean);

  return (
    <>
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
          size="small"
          onClick={handleToggle}
          variant="outlined"
        >
          <MoreHorizIcon />
        </Button>
      </ButtonGroup>

      <ClickAwayListener onClickAway={handleClose}>
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
            <MenuList id="split-button-menu" autoFocusItem>
              {options.map((option, index) => (
                <MenuItem
                  disabled={option.disabled}
                  key={index}
                  onClick={(event) => {
                    handleMenuItemClick(event);
                    option.onClick(event);
                  }}
                >
                  <div style={{ marginRight: '0.5rem' }}>{option.icon}</div>
                  {option.label}
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
