import React from 'react';
import { Box, crimson, InfoOutlinedIcon, AccountTreeIcon, EditIcon } from '@sistent/sistent';
import Moment from 'react-moment';
import { GetApp as GetAppIcon } from '@/assets/icons';
import { DoneAll as DoneAllIcon, Public as PublicIcon } from '@/assets/icons';
import UndeployIcon from '../../../public/static/img/UndeployIcon';
import CloneIcon from '../../../public/static/img/CloneIcon';

import { DefaultTableCell, SortableTableCell } from '../../connections/common';
import CheckIcon from '@/assets/icons/CheckIcon';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import PatternConfigureIcon from '@/assets/icons/PatternConfigure';
import ActionPopover from './ActionPopover';
import CustomToolbarSelect from './CustomToolbarSelect';
import { VISIBILITY } from '../../../utils/Enum';
import { genericClickHandler } from './MesheryPatterns.constants';

/**
 * Builds the per-row action items list for the Designs table.
 * Extracted from the original Actions-column customBodyRender so the
 * entry-point file stays under the line-budget. Behavior is identical:
 * the same set of actions are returned in the same order, and each item's
 * onClick / disabled / condition values are wired to the same handlers.
 */
export function buildPatternActions({
  rowData,
  visibility,
  patterns,
  tableMeta,
  handlers,
  permissions,
}) {
  const {
    handleOpenInConfigurator,
    handleClone,
    openValidateModal,
    openDryRunModal,
    openUndeployModal,
    openDeployModal,
    handleDesignDownloadModal,
    handleInfoModal,
    handleUnpublishModal,
    handleEvaluateRelationship,
    userCanEdit,
  } = handlers;

  return [
    {
      label: 'Edit',
      icon: <EditIcon fill="currentColor" />,
      onClick: (e) => {
        e.stopPropagation();
        handleOpenInConfigurator(rowData.id);
      },
      disabled: !permissions.editDesign,
      condition: userCanEdit(rowData),
    },
    {
      label: 'Clone',
      icon: <CloneIcon fill="currentColor" />,
      onClick: (e) => {
        e.stopPropagation();
        handleClone(rowData.id, rowData.name);
      },
      disabled: !permissions.cloneDesign,
      condition: visibility === VISIBILITY.PUBLISHED,
    },
    {
      label: 'Edit',
      icon: <PatternConfigureIcon />,
      onClick: (e) => {
        e.stopPropagation();
        handleOpenInConfigurator(patterns[tableMeta.rowIndex].id);
      },
      disabled: !permissions.editDesign,
      condition: visibility !== VISIBILITY.PUBLISHED,
    },
    {
      label: 'Validate Design',
      icon: <CheckIcon data-cy="verify-button" />,
      onClick: (e) => {
        openValidateModal(e, rowData.patternFile, rowData.name, rowData.id);
      },
      disabled: !permissions.validateDesign,
    },
    {
      label: 'Dry Run',
      icon: <DryRunIcon data-cy="verify-button" />,
      onClick: (e) => {
        openDryRunModal(e, rowData.patternFile, rowData.name, rowData.id);
      },
      disabled: !permissions.validateDesign,
    },
    {
      label: 'Evaluate',
      icon: <AccountTreeIcon fill={'currentColor'} data-cy="evaluate-button" />,
      onClick: (e) => {
        e.stopPropagation();
        handleEvaluateRelationship(rowData);
      },
      disabled: !permissions.evaluateRelationships,
    },
    {
      label: 'Undeploy',
      icon: <UndeployIcon fill={crimson[40]} data-cy="undeploy-button" />,
      onClick: (e) => {
        openUndeployModal(e, rowData.patternFile, rowData.name, rowData.id);
      },
      disabled: !permissions.undeployDesign,
    },
    {
      label: 'Deploy',
      icon: <DoneAllIcon fill="currentColor" data-cy="deploy-button" />,
      onClick: (e) => {
        openDeployModal(e, rowData.patternFile, rowData.name, rowData.id);
      },
      disabled: !permissions.deployDesign,
    },
    {
      label: 'Download',
      icon: <GetAppIcon data-cy="download-button" />,
      onClick: (e) => {
        handleDesignDownloadModal(e, rowData);
      },
      disabled: !permissions.downloadDesign,
    },
    {
      label: 'Design Information',
      icon: <InfoOutlinedIcon data-cy="information-button" />,
      onClick: (e) => {
        genericClickHandler(e, () => handleInfoModal(rowData));
      },
      disabled: !permissions.detailsOfDesign,
    },

    /* Publish action can be done through Info modal so we might not need separate publish action */
    /*{
      label="Publish",
      icon: <PublicIcon fill="#F91313" data-cy="publish-button" />,
      onClick: (e) => handlePublishModal(e, rowData)(),
      disabled: !permissions.publishDesign,
      condition: canPublishPattern && visibility !== VISIBILITY.PUBLISHED,
    },*/

    {
      label: 'Unpublish',
      icon: <PublicIcon fill={crimson[40]} data-cy="unpublish-button" />,
      onClick: (e) => {
        handleUnpublishModal(e, rowData)();
      },
      disabled: !permissions.unpublishDesign,
      condition: visibility === VISIBILITY.PUBLISHED,
    },
  ].filter((action) => action.condition === undefined || action.condition);
}

/**
 * Builds the responsive-column view configuration used to determine which
 * columns are visible at different viewport widths.
 */
export const PATTERN_COL_VIEWS = [
  ['name', 'xs'],
  ['created_at', 'm'],
  ['updated_at', 'm'],
  ['visibility', 's'],
  ['Actions', 'xs'],
];

/**
 * Builds the full column definition array for the Designs table.
 *
 * The Actions column's customBodyRender closes over `patterns` and the
 * caller-supplied `handlers` so it picks up fresh state on every render;
 * we therefore rebuild the columns inside the parent component on each
 * render (matching the original behavior).
 */
export function buildPatternColumns({ patterns, handlers, permissions }) {
  return [
    {
      name: 'name',
      label: 'Name',
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
      label: 'Updated At',
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
          const visibility = patterns[tableMeta.rowIndex]?.visibility;
          const actions = buildPatternActions({
            rowData,
            visibility,
            patterns,
            tableMeta,
            handlers,
            permissions,
          });

          return (
            <Box sx={{ display: 'flex', justifyContent: 'center' }}>
              <ActionPopover actions={actions} />
            </Box>
          );
        },
      },
    },
  ];
}

/**
 * Builds the `options` config for the Designs ResponsiveDataTable.
 *
 * Closes over caller-supplied state (patterns array, paging info) and
 * caller-supplied callbacks (deletePatterns, showModal, setters,
 * subscription init, etc.) so the table picks up fresh data on every
 * render. Behavior matches the original inline `options` block one-for-one.
 */
export function buildPatternsTableOptions({
  patterns,
  columns,
  count,
  pageSize,
  page,
  search,
  sortOrder,
  isLocalProvider,
  searchTimeout,
  setPage,
  setPageSize,
  setSearch,
  setSortOrder,
  setSelectedRowData,
  deletePatterns,
  showModal,
}) {
  return {
    customToolbarSelect: (selectedRows, displayData, setSelectedRows) => (
      <CustomToolbarSelect
        selectedRows={selectedRows}
        displayData={displayData}
        setSelectedRows={setSelectedRows}
        patterns={patterns}
        deletePatterns={deletePatterns}
        showModal={showModal}
      />
    ),
    filter: false,
    search: false,
    viewColumns: false,
    sort: !isLocalProvider,
    filterType: 'textField',
    responsive: 'standard',
    resizableColumns: true,
    serverSide: true,
    count,
    rowsPerPage: pageSize,
    fixedHeader: true,
    page,
    print: false,
    download: false,
    sortOrder: {
      name: sortOrder.split(' ')[0],
      direction: sortOrder.split(' ')[1],
    },
    textLabels: {
      selectedRows: {
        text: 'pattern(s) selected',
      },
    },

    onCellClick: (_, meta) =>
      meta.colIndex !== 3 && meta.colIndex !== 4 && setSelectedRowData(patterns[meta.rowIndex]),

    onRowsDelete: async function handleDelete(row) {
      const toBeDeleted = Object.keys(row.lookup).map((idx) => ({
        id: patterns[idx]?.id,
        name: patterns[idx]?.name,
      }));
      let response = await showModal(
        toBeDeleted.length,
        toBeDeleted.map((p) => ' ' + p.name),
      );
      if (response && response.trim().toUpperCase() === 'DELETE') {
        deletePatterns({ patterns: toBeDeleted });
      }
      // if (response.toLowerCase() === "no")
      // fetchPatterns(page, pageSize, search, sortOrder);
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn != null) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              setSearch(tableState.searchText);
            }
          }, 500);
          break;
        case 'sort':
          if (sortInfo.length === 2) {
            if (sortInfo[1] === 'ascending') {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            setSortOrder(order);
          }
          break;
      }
    },
    setRowProps: (row, dataIndex, rowIndex) => {
      return {
        'data-cy': `config-row-${rowIndex}`,
      };
    },
    setTableProps: () => {
      return {
        'data-cy': 'filters-grid',
      };
    },
  };
}
