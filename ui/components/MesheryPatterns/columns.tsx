import React from 'react';
import { Box } from '@sistent/sistent';
import Moment from 'react-moment';
import { SortableTableCell, DefaultTableCell } from '../connections/common';
import { VISIBILITY } from '../../utils/Enum';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import ActionPopover from './ActionPopover';
import { Edit as EditIcon, Done as CheckIcon, DoneAll as DoneAllIcon, Download as GetAppIcon, InfoOutlined as InfoOutlinedIcon, Public as PublicIcon } from '@/components/icons';
import CloneIcon from '../../public/static/img/CloneIcon';
import PatternConfigureIcon from '@/assets/icons/PatternConfigure';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import UndeployIcon from '../../public/static/img/UndeployIcon';

export const getColumns = ({
  patterns,
  sortOrder,
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
}) => [
  {
    name: 'name',
    label: 'Design Name',
    options: {
      filter: false,
      sort: true,
      searchable: true,
      customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
        return (
          <SortableTableCell
            index={index}
            columnData={column}
            columnMeta={columnMeta}
            onSort={() => sortColumn(index)}
          />
        );
      },
    },
  },
  {
    name: 'created_at',
    label: 'Created At',
    options: {
      filter: false,
      sort: true,
      searchable: true,
      customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
        return (
          <SortableTableCell
            index={index}
            columnData={column}
            columnMeta={columnMeta}
            onSort={() => sortColumn(index)}
          />
        );
      },
      customBodyRender: function CustomBody(value) {
        return <Moment format="LLLL">{value}</Moment>;
      },
    },
  },
  {
    name: 'updated_at',
    label: 'Update At',
    options: {
      filter: false,
      sort: true,
      searchable: true,
      customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
        return (
          <SortableTableCell
            index={index}
            columnData={column}
            columnMeta={columnMeta}
            onSort={() => sortColumn(index)}
          />
        );
      },
      customBodyRender: function CustomBody(value) {
        return <Moment format="LLLL">{value}</Moment>;
      },
    },
  },
  {
    name: 'visibility',
    label: 'Visibility',
    options: {
      filter: false,
      sort: true,
      searchable: true,
      customHeadRender: function CustomHead({ ...column }) {
        return <DefaultTableCell columnData={column} />;
      },
    },
  },
  {
    name: 'Actions',
    label: 'Actions',
    options: {
      filter: false,
      sort: false,
      searchable: false,
      customHeadRender: function CustomHead({ ...column }) {
        return <DefaultTableCell columnData={column} />;
      },
      customBodyRender: function CustomBody(_, tableMeta) {
        const rowData = patterns[tableMeta.rowIndex];
        const visibility = rowData?.visibility;
        const actions = [
          {
            label: 'Edit',
            icon: <EditIcon fill="currentColor" />,
            onClick: (e) => {
              e.stopPropagation();
              handleOpenInConfigurator(rowData.id);
            },
            disabled: !CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject),
            condition: userCanEdit(rowData),
          },
          {
            label: 'Clone',
            icon: <CloneIcon fill="currentColor" />,
            onClick: (e) => {
              e.stopPropagation();
              handleClone(rowData.id, rowData.name);
            },
            disabled: !CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject),
            condition: visibility === VISIBILITY.PUBLISHED,
          },
          {
            label: 'Design',
            icon: <PatternConfigureIcon />,
            onClick: (e) => {
              e.stopPropagation();
              handleOpenInConfigurator(rowData.id);
            },
            disabled: !CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject),
            condition: visibility !== VISIBILITY.PUBLISHED,
          },
          {
            label: 'Validate Design',
            icon: <CheckIcon data-cy="verify-button" />,
            onClick: (e) => {
              openValidateModal(e, rowData.pattern_file, rowData.name, rowData.id);
            },
            disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
          },
          {
            label: 'Dry Run',
            icon: <DryRunIcon data-cy="verify-button" />,
            onClick: (e) => {
              openDryRunModal(e, rowData.pattern_file, rowData.name, rowData.id);
            },
            disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
          },
          {
            label: 'Undeploy',
            icon: <UndeployIcon fill="#F91313" data-cy="undeploy-button" />,
            onClick: (e) => {
              openUndeployModal(e, rowData.pattern_file, rowData.name, rowData.id);
            },
            disabled: !CAN(keys.UNDEPLOY_DESIGN.action, keys.UNDEPLOY_DESIGN.subject),
          },
          {
            label: 'Deploy',
            icon: <DoneAllIcon data-cy="deploy-button" />,
            onClick: (e) => {
              openDeployModal(e, rowData.pattern_file, rowData.name, rowData.id);
            },
            disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
          },
          {
            label: 'Download',
            icon: <GetAppIcon data-cy="download-button" />,
            onClick: (e) => {
              handleDesignDownloadModal(e, rowData);
            },
            disabled: !CAN(keys.DOWNLOAD_A_DESIGN.action, keys.DOWNLOAD_A_DESIGN.subject),
          },
          {
            label: 'Design Information',
            icon: <InfoOutlinedIcon data-cy="information-button" />,
            onClick: (e) => {
              e.stopPropagation();
              handleInfoModal(rowData);
            },
            disabled: !CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject),
          },
          {
            label: 'Unpublish',
            icon: <PublicIcon fill="#F91313" data-cy="unpublish-button" />,
            onClick: (e) => {
              handleUnpublishModal(e, rowData)();
            },
            disabled: !CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject),
            condition: visibility === VISIBILITY.PUBLISHED,
          },
        ].filter((action) => action.condition === undefined || action.condition);

        return (
          <Box sx={{ display: 'flex', justifyContent: 'center' }}>
            <ActionPopover actions={actions} />
          </Box>
        );
      },
    },
  },
];
