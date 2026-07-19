import React from 'react';
import {
  CustomColumnVisibilityControl,
  DataTableToolbar,
  SearchBar,
  UniversalFilter,
} from '@sistent/sistent';
import { Publish as PublishIcon } from '@/assets/icons';
import ViewSwitch from '../../ViewSwitch';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import TooltipButton from '@/utils/TooltipButton';
import { AddIconStyled, BtnText } from './MesheryPatterns.styled';

/**
 * Header toolbar for the Designs page.
 *
 * Migrated to Sistent's `DataTableToolbar` (meshery/meshery#20658): the
 * create/import buttons, search bar, universal filter, column-visibility
 * control and grid/table switch are now passed in as slots instead of being
 * laid out by hand with ToolWrapper/CreateButton/SearchWrapper, so this
 * toolbar shares layout, spacing and responsive behavior with every other
 * migrated table toolbar in the app. Slot content/behavior is unchanged
 * from the previous implementation.
 */
function MesheryPatternsToolbar({
  width,
  isSearchExpanded,
  setIsSearchExpanded,
  selectedPattern,
  viewType,
  setViewType,
  disableCreateImportDesignButton,
  disableUniversalFilter,
  pageTitle,
  router,
  handleUploadImport,
  setSearch,
  filter,
  selectedFilters,
  setSelectedFilters,
  handleApplyFilter,
  columns,
  columnVisibility,
  setColumnVisibility,
}) {
  // On narrow viewports, hide the create/import buttons while the search
  // bar is expanded so the search bar can take the full width.
  const hideActionsRow = width < 600 && isSearchExpanded;

  const primaryActions =
    hideActionsRow || selectedPattern.show || disableCreateImportDesignButton ? null : (
      <div style={{ display: 'flex', order: 1 }}>
        <TooltipButton
          title="Create Design"
          data-testid="meshery-patterns-create-design-btn"
          aria-label="Create Design"
          variant="contained"
          color="primary"
          size="large"
          onClick={() => router.push('/configuration/designs/configurator')}
          style={{ display: 'flex', marginRight: '2rem' }}
          disabled={
            !CAN(
              Keys.CatalogManagementCreateNewDesign.id,
              Keys.CatalogManagementCreateNewDesign.function,
            )
          }
        >
          <AddIconStyled />
          <BtnText> Create Design </BtnText>
        </TooltipButton>
        <TooltipButton
          title="Import Design"
          data-testid="meshery-patterns-import-design-btn"
          aria-label="Import Design"
          variant="contained"
          color="primary"
          size="large"
          onClick={handleUploadImport}
          style={{ display: 'flex', marginRight: '2rem', marginLeft: '-0.6rem' }}
          disabled={
            !CAN(Keys.CatalogManagementImportDesign.id, Keys.CatalogManagementImportDesign.function)
          }
        >
          <AddIconStyled>
            <PublishIcon />
          </AddIconStyled>
          <BtnText> Import Design </BtnText>
        </TooltipButton>
      </div>
    );

  return (
    <DataTableToolbar
      primaryActions={primaryActions}
      search={
        <SearchBar
          onSearch={(value) => {
            setSearch(value);
          }}
          expanded={isSearchExpanded}
          setExpanded={setIsSearchExpanded}
          placeholder={`Search ${pageTitle.toLowerCase()}...`}
          data-testid="meshery-patterns-search-bar"
        />
      }
      filter={
        disableUniversalFilter ? null : (
          <UniversalFilter
            id="designs-universal-filter-ref"
            filters={filter}
            selectedFilters={selectedFilters}
            setSelectedFilters={setSelectedFilters}
            handleApplyFilter={handleApplyFilter}
            data-testid="meshery-patterns-universal-filter"
          />
        )
      }
      columnVisibility={
        viewType === 'table' ? (
          <CustomColumnVisibilityControl
            data-testid="meshery-patterns-column-visibility-control"
            id="designs-column-visibility-ref"
            columns={columns}
            customToolsProps={{ columnVisibility, setColumnVisibility }}
          />
        ) : null
      }
      viewSwitch={
        !selectedPattern.show ? <ViewSwitch view={viewType} changeView={setViewType} /> : null
      }
    />
  );
}

export default MesheryPatternsToolbar;
