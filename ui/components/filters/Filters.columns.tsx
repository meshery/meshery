import React from 'react';
import { GetApp as GetAppIcon, Public as PublicIcon } from '@/assets/icons';
import { InfoOutlinedIcon, EditIcon } from '@sistent/sistent';
import Moment from 'react-moment';
import CloneIcon from '../../public/static/img/CloneIcon';
import { iconMedium } from '../../css/icons.styles';
import { VISIBILITY } from '../../utils/Enum';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { DefaultTableCell, SortableTableCell } from '../connections/common/index';
import TooltipIcon from './TooltipIcon';
import { ActionsBox } from './Filters.styled';

type BuildFiltersColumnsArgs = {
  filters: any[];
  canPublishFilter: boolean;
  handleClone: (_id: string, _name: string) => void;
  handleDownload: (_e: React.MouseEvent, _id: string, _name: string) => void;
  handleInfoModal: (_filter: any) => void;
  handlePublishModal: (_ev: React.MouseEvent, _filter: any) => void;
  handleUnpublishModal: (_ev: React.MouseEvent, _filter: any) => () => Promise<void>;
  setSelectedRowData: (_row: any) => void;
  sortOrder: string;
};

export function buildFiltersColumns({
  filters,
  canPublishFilter,
  handleClone,
  handleDownload,
  handleInfoModal,
  handlePublishModal,
  handleUnpublishModal,
  setSelectedRowData,
  sortOrder,
}: BuildFiltersColumnsArgs) {
  const columns: any[] = [
    {
      name: 'name',
      label: 'Filter Name',
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
          const rowData = filters[tableMeta.rowIndex];
          const visibility = filters[tableMeta.rowIndex]?.visibility;
          return (
            <ActionsBox
              sx={{
                display: 'flex',
              }}
            >
              {visibility === VISIBILITY.PUBLISHED ? (
                <TooltipIcon
                  placement="top"
                  title={'Clone'}
                  onClick={(e) => {
                    e.stopPropagation();
                    handleClone(rowData.id, rowData.name);
                  }}
                  disabled={
                    !CAN(
                      Keys.CatalogManagementCloneWasmFilter.id,
                      Keys.CatalogManagementCloneWasmFilter.function,
                    )
                  }
                >
                  <CloneIcon fill="currentColor" />
                </TooltipIcon>
              ) : (
                <TooltipIcon
                  title="Config"
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedRowData(filters[tableMeta.rowIndex]);
                  }}
                  disabled={
                    !CAN(
                      Keys.CatalogManagementEditWasmFilter.id,
                      Keys.CatalogManagementEditWasmFilter.function,
                    )
                  }
                >
                  <EditIcon aria-label="config" color="inherit" style={iconMedium} />
                </TooltipIcon>
              )}
              <TooltipIcon
                title="Download"
                onClick={(e) => handleDownload(e, rowData.id, rowData.name)}
                disabled={
                  !CAN(
                    Keys.CatalogManagementDownloadAWasmFilter.id,
                    Keys.CatalogManagementDownloadAWasmFilter.function,
                  )
                }
              >
                <GetAppIcon data-cy="download-button" />
              </TooltipIcon>
              <TooltipIcon
                title="Filter Information"
                onClick={() => handleInfoModal(rowData)}
                disabled={
                  !CAN(
                    Keys.CatalogManagementDetailsOfWasmFilter.id,
                    Keys.CatalogManagementDetailsOfWasmFilter.function,
                  )
                }
              >
                <InfoOutlinedIcon data-cy="information-button" />
              </TooltipIcon>
              {canPublishFilter && visibility !== VISIBILITY.PUBLISHED ? (
                <TooltipIcon
                  title="Publish"
                  onClick={(ev) => handlePublishModal(ev, rowData)}
                  disabled={
                    !CAN(
                      Keys.CatalogManagementPublishWasmFilter.id,
                      Keys.CatalogManagementPublishWasmFilter.function,
                    )
                  }
                >
                  <PublicIcon fill="#F91313" data-cy="publish-button" />
                </TooltipIcon>
              ) : (
                <TooltipIcon
                  title="Unpublish"
                  onClick={(ev) => handleUnpublishModal(ev, rowData)?.()}
                  disabled={
                    !CAN(
                      Keys.CatalogManagementUnpublishWasmFilter.id,
                      Keys.CatalogManagementUnpublishWasmFilter.function,
                    )
                  }
                >
                  <PublicIcon fill="#F91313" data-cy="unpublish-button" />
                </TooltipIcon>
              )}
            </ActionsBox>
          );
        },
      },
    },
  ];

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(' ')[0]) {
      columns[idx].options.sortDirection = sortOrder.split(' ')[1];
    }
  });

  return columns;
}
