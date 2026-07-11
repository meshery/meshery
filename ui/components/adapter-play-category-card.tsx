import React from 'react';
import {
  AddIcon,
  Box,
  CardActions,
  CardHeader,
  DeleteIcon,
  IconButton,
  Menu,
  MenuItem,
  PlayArrowIcon as PlayIcon,
} from '@sistent/sistent';
import { iconMedium } from '../css/icons.styles';
import { Keys } from '@meshery/schemas/permissions';
import CAN from '@/utils/can';
import { AdapterCard } from './adapter-play-styled';

interface AdapterOperation {
  key: string;
  value: string;
  category?: number;
}

interface MenuStateEntry {
  add: boolean;
  delete: boolean;
}

interface CategoryContent {
  content: string;
  description: string;
  permission: { action: string; subject: string };
}

interface AdapterCategoryCardProps {
  cat: number;
  adapterOps: AdapterOperation[];
  menuState: { [key: number]: MenuStateEntry };
  addIconEles: React.MutableRefObject<{ [key: number]: HTMLElement | null }>;
  delIconEles: React.MutableRefObject<{ [key: number]: HTMLElement | null }>;
  onMenuToggle: (cat: number, isDelete: boolean) => void;
  onMenuItemClick: (cat: number, opKey: string, isDelete: boolean) => () => void;
  renderYamlDialog: (cat: number, isDelete: boolean) => React.ReactNode;
}

const CATEGORY_DEFS: Record<number, CategoryContent> = {
  0: {
    content: 'Manage Cloud Native Infrastructure Lifecycle',
    description: 'Deploy cloud native infrastructure or SMI adapter on your cluster.',
    permission: {
      action: Keys.InfrastructureManagementManageCloudNativeInfrastructureLifeCycle.id,
      subject: Keys.InfrastructureManagementManageCloudNativeInfrastructureLifeCycle.function,
    },
  },
  1: {
    content: 'Manage Sample Application Lifecycle',
    description: 'Deploy sample applications on/off the service mesh.',
    permission: {
      action: Keys.InfrastructureManagementManageCloudNativeInfrastructureLifeCycle.id,
      subject: Keys.InfrastructureManagementManageCloudNativeInfrastructureLifeCycle.function,
    },
  },
  2: {
    content: 'Apply Cloud Native Infrastructure Configuration',
    description: 'Configure your cloud native infrastructure using some pre-defined options.',
    permission: {
      action: Keys.InfrastructureManagementApplyCloudNativeInfrastructureConfiguration.id,
      subject: Keys.InfrastructureManagementApplyCloudNativeInfrastructureConfiguration.function,
    },
  },
  3: {
    content: 'Validate Cloud Native Infrastructure Configuration',
    description: 'Validate your cloud native infrastructure configuration against best practices.',
    permission: {
      action: Keys.InfrastructureManagementValidateCloudNativeInfrastructureConfiguration.id,
      subject: Keys.InfrastructureManagementValidateCloudNativeInfrastructureConfiguration.function,
    },
  },
  4: {
    content: 'Apply Custom Configuration',
    description: 'Customize the configuration of your cloud native infrastructure.',
    permission: {
      action: Keys.InfrastructureManagementApplyCustomCloudNativeConfiguration.id,
      subject: Keys.InfrastructureManagementApplyCustomCloudNativeConfiguration.function,
    },
  },
};

const renderMenu = (
  cat: number,
  isDelete: boolean,
  anchorEl: HTMLElement | null,
  menuState: { [key: number]: MenuStateEntry },
  onMenuToggle: AdapterCategoryCardProps['onMenuToggle'],
  onMenuItemClick: AdapterCategoryCardProps['onMenuItemClick'],
  selectedAdapterOps: AdapterOperation[],
) => {
  return (
    <Menu
      id="long-menu"
      anchorEl={anchorEl}
      keepMounted
      open={menuState[cat][isDelete ? 'delete' : 'add']}
      onClose={() => onMenuToggle(cat, isDelete)}
    >
      {selectedAdapterOps
        .sort((adap1, adap2) => adap1.value.localeCompare(adap2.value))
        .map(({ key, value }) => (
          <MenuItem key={key} onClick={onMenuItemClick(cat, key, isDelete)}>
            {value}
          </MenuItem>
        ))}
    </Menu>
  );
};

/**
 * Renders a single adapter category card (Lifecycle, Configuration,
 * Validation, etc.) with its install/uninstall menu pair. Extracted
 * from MesheryAdapterPlayComponent's generateCardForCategory() — the
 * filtering, permission checks, and menu wiring are unchanged.
 */
const AdapterCategoryCard: React.FC<AdapterCategoryCardProps> = ({
  cat,
  adapterOps,
  menuState,
  addIconEles,
  delIconEles,
  onMenuToggle,
  onMenuItemClick,
  renderYamlDialog,
}) => {
  let selectedAdapterOps = adapterOps
    ? adapterOps.filter(
        ({ category }) => (typeof category === 'undefined' && cat === 0) || category === cat,
      )
    : [];
  if (cat === 2) {
    selectedAdapterOps = selectedAdapterOps.filter((ops) => !ops.value.startsWith('Add-on:'));
  }

  const def = CATEGORY_DEFS[cat];
  const { content, description, permission } = def;

  return (
    <AdapterCard>
      <CardHeader title={content} subheader={description} style={{ flexGrow: 1 }} />
      <CardActions disableSpacing>
        <IconButton
          aria-label="install"
          ref={(ch) => (addIconEles.current[cat] = ch)}
          onClick={() => onMenuToggle(cat, false)}
          disabled={!CAN(permission.action, permission.subject)}
        >
          {cat !== 4 ? <AddIcon style={iconMedium} /> : <PlayIcon style={iconMedium} />}
        </IconButton>
        {cat !== 4 &&
          renderMenu(
            cat,
            false,
            addIconEles.current[cat],
            menuState,
            onMenuToggle,
            onMenuItemClick,
            selectedAdapterOps,
          )}
        {cat === 4 && renderYamlDialog(cat, false)}
        {cat !== 3 && (
          <Box sx={{ width: '100%' }}>
            <IconButton
              aria-label="delete"
              ref={(ch) => (delIconEles.current[cat] = ch)}
              style={{ float: 'right' }}
              onClick={() => onMenuToggle(cat, true)}
              disabled={!CAN(permission.action, permission.subject)}
            >
              <DeleteIcon style={iconMedium} />
            </IconButton>
            {cat !== 4 &&
              renderMenu(
                cat,
                true,
                delIconEles.current[cat],
                menuState,
                onMenuToggle,
                onMenuItemClick,
                selectedAdapterOps,
              )}
            {cat === 4 && renderYamlDialog(cat, true)}
          </Box>
        )}
      </CardActions>
    </AdapterCard>
  );
};

export default AdapterCategoryCard;
