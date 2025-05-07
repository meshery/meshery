import { Grid, Button, Pagination } from '@layer5/sistent';
import React, { useState } from 'react';
import FiltersCard from './FiltersCard';
import { FILE_OPS } from '../../utils/Enum';
import {
  GridAddIconStyles,
  GridNoContainerStyles,
  GridNoPapperStyles,
  GridNoTextStyles,
  GridPaginationStyles,
} from '../MesheryPatterns/Grid.styles';
import { RJSFModalWrapper } from '../Modal';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { Modal as SistentModal } from '@layer5/sistent';

import Filter from '../../public/static/img/drawer-icons/filter_svg.js';

const INITIAL_GRID_SIZE = { xl: 6, md: 6, xs: 12 };

function FilterCardGridItem({
  filter,
  yamlConfig,
  handleDownload,
  handleSubmit,
  setSelectedFilters,
  handleClone,
  handlePublishModal,
  handleUnpublishModal,
  canPublishFilter,
  handleInfoModal,
}) {
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [yaml, setYaml] = useState(yamlConfig);

  return (
    <Grid item {...gridProps}>
      <FiltersCard
        name={filter.name}
        updated_at={filter.updated_at}
        created_at={filter.created_at}
        ownerId={filter.user_id}
        filter_resource={yaml}
        canPublishFilter={canPublishFilter}
        handlePublishModal={handlePublishModal}
        handleUnpublishModal={handleUnpublishModal}
        requestFullSize={() => setGridProps({ xl: 12, md: 12, xs: 12 })}
        requestSizeRestore={() => setGridProps(INITIAL_GRID_SIZE)}
        handleClone={handleClone}
        handleDownload={(ev) => handleDownload(ev, filter.id, filter.name)}
        deleteHandler={() =>
          handleSubmit({
            data: yaml,
            id: filter.id,
            type: FILE_OPS.DELETE,
            name: filter.name,
            catalog_data: filter.catalog_data,
          })
        }
        updateHandler={() =>
          handleSubmit({
            data: yaml,
            id: filter.id,
            type: FILE_OPS.UPDATE,
            name: filter.name,
            catalog_data: filter.catalog_data,
          })
        }
        setSelectedFilters={() => setSelectedFilters({ filter: filter, show: true })}
        setYaml={setYaml}
        description={filter.desciption}
        visibility={filter.visibility}
        handleInfoModal={handleInfoModal}
      />
    </Grid>
  );
}

/**
 * FilterGrid is the react component for rendering grid
 * @param {{
 *  filters:Array<{
 *  id:string,
 *  created_at: string,
 *  updated_at: string,
 *  filter_file: string,
 * }>,
 *  handleVerify: (e: Event, filter_file: any, filter_id: string) => void,
 *  handlePublish: (catalog_data : any) => void,
 *  handleUnpublishModal: (ev: Event, filter: any) => (() => Promise<void>),
 *  handleDeploy: (filter_file: any) => void,
 *  handleUnDeploy: (filter_file: any) => void,
 *  handleDownload: (e : Event, id : string , name : string ) => void,
 *  handleSubmit: (data: any, id: string, name: string, type: string) => void,
 *  setSelectedFilter : ({show: boolean, filter:any}) => void,
 *  selectedFilter: {show : boolean, filter : any},
 *  pages?: number,
 *  selectedPage?: number,
 *  setPage: (page: number) => void,
 *  filterErrors: Map,
 *  canPublishFilter: boolean,
 *  publishModal: {
 *   open: boolean,
 *   filter: any,
 *   name: string,
 *  },
 *  setPublishModal: (publishModal: { open: boolean, filter: any, name: string }) => void,
 *  publishSchema: object,
 *  handleInfoModal: (filter: object) => void
 * }} props props
 */

function FiltersGrid({
  filters = [],
  handleClone,
  handleDownload,
  handleSubmit,
  setSelectedFilter,
  selectedFilter,
  pages = 1,
  setPage,
  selectedPage,
  canPublishFilter,
  handleUploadImport,
  handlePublish,
  handleUnpublishModal,
  publishModal,
  setPublishModal,
  publishSchema,
  handleInfoModal,
}) {
  const handlePublishModal = (filter) => {
    if (canPublishFilter) {
      setPublishModal({
        open: true,
        filter: filter,
        name: '',
      });
    }
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open: false,
      filter: {},
      name: '',
    });
  };

  const getYamlConfig = (filter_resource) => {
    if (filter_resource) {
      return JSON.parse(filter_resource).settings.config;
    }

    return '';
  };

  return (
    <div>
      {!selectedFilter.show && (
        <Grid container spacing={3}>
          {filters.map((filter) => (
            <FilterCardGridItem
              key={filter.id}
              filter={filter}
              yamlConfig={getYamlConfig(filter.filter_resource)}
              handleClone={() => handleClone(filter.id, filter.name)}
              handleDownload={handleDownload}
              handleSubmit={handleSubmit}
              setSelectedFilters={setSelectedFilter}
              canPublishFilter={canPublishFilter}
              handlePublishModal={() => handlePublishModal(filter)}
              handleUnpublishModal={(e) => handleUnpublishModal(e, filter)()}
              handleInfoModal={() => handleInfoModal(filter)}
            />
          ))}
        </Grid>
      )}
      {!selectedFilter.show && filters.length === 0 && (
        <GridNoPapperStyles>
          <GridNoContainerStyles>
            <GridNoTextStyles align="center" color="textSecondary">
              No Filters Found
            </GridNoTextStyles>
            <div>
              <Button
                aria-label="Add Application"
                variant="contained"
                color="primary"
                disabled={!CAN(keys.IMPORT_FILTER.action, keys.IMPORT_FILTER.subject)}
                size="large"
                // @ts-ignore
                onClick={handleUploadImport}
                style={{ marginRight: '2rem' }}
              >
                <GridAddIconStyles />
                Import Filter
              </Button>
            </div>
          </GridNoContainerStyles>
        </GridNoPapperStyles>
      )}
      {filters.length ? (
        <GridPaginationStyles>
          <Pagination
            count={pages}
            page={selectedPage + 1}
            onChange={(_, page) => setPage(page - 1)}
          />
        </GridPaginationStyles>
      ) : null}
      {canPublishFilter && publishModal.open && (
        <SistentModal
          open={true}
          title={publishModal.filter?.name}
          closeModal={handlePublishModalClose}
          maxWidth="sm"
          headerIcon={
            <Filter
              fill="#fff"
              style={{ height: '24px', width: '24px', fonSize: '1.45rem' }}
              className={undefined}
            />
          }
        >
          <RJSFModalWrapper
            schema={publishSchema.rjsfSchema}
            uiSchema={publishSchema.uiSchema}
            handleSubmit={handlePublish}
            handleClose={handlePublishModalClose}
            submitBtnText="Submit for Approval"
            helpText="Upon submitting your catalog item, an approval flow will be initiated.[Learn more](https://docs.meshery.io/concepts/catalog)"
          />
        </SistentModal>
      )}
    </div>
  );
}

export default FiltersGrid;
