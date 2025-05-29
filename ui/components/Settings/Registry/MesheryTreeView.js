import React, { useEffect, useState, useRef, useCallback } from 'react';
import {
  IconButton,
  FormControlLabel,
  Switch,
  CircularProgress,
  Typography,
} from '@layer5/sistent';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../../constants/navigator';
import SearchBar from '@/utils/custom-search';
import debounce from '@/utils/debounce';
import { useWindowDimensions } from '@/utils/dimension';
import { getFilteredDataForDetailsComponent } from './helper';
import { CustomTextTooltip } from '@/components/MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import _ from 'lodash';
import CollapseAllIcon from '@/assets/icons/CollapseAll';
import ExpandAllIcon from '@/assets/icons/ExpandAll';
import { Colors } from '@/themes/app';
import { JustifyAndAlignCenter, MesheryTreeViewWrapper } from './MeshModel.style';
import { useRegistryRouter } from './hooks';
import MesheryTreeViewModel from './MesheryTreeViewModel';
import MesheryTreeViewRegistrants from './MesheryTreeViewRegistrants';
import ComponentTree from './ComponentTree';
import RelationshipTree from './RelationshipTree';

const MesheryTreeView = React.memo(
  ({
    data,
    view,
    setSearchText,
    searchText,
    setPage,
    checked,
    setChecked,
    setShowDetailsData,
    showDetailsData,
    setResourcesDetail,
    setModelsFilters,
    lastItemRef,
    isFetching,
    isLoading,
  }) => {
    const { handleUpdateSelectedRoute, selectedItemUUID } = useRegistryRouter();
    const [expanded, setExpanded] = React.useState([]);
    const [selected, setSelected] = React.useState([]);
    const { width } = useWindowDimensions();
    const [isSearchExpanded, setIsSearchExpanded] = useState(searchText ? true : false);
    const [prevState, setPrevState] = useState({ data: [], uuid: '' });
    const scrollRef = useRef();

    const handleScroll = (scrollingView) => (event) => {
      const div = event.target;
      if (div.scrollTop >= div.scrollHeight - div.clientHeight - 1) {
        setPage((prevPage) => ({
          ...prevPage,
          [scrollingView]: Number(prevPage[scrollingView]) + 1,
        }));
      }
      if (!data.length === 0) {
        scrollRef.current = div.scrollTop;
      }
    };

    useEffect(() => {
      if (scrollRef.current) {
        const div = document.querySelector('.scrollElement');
        div.scrollTop = scrollRef.current;
      }
    }, [data]);

    const handleChecked = useCallback(() => {
      setChecked((prevChecked) => !prevChecked);
    }, [setChecked]);

    // Expand first level tree
    const expandAll = () => {
      const arr = [];
      data.map((parent) => {
        if (view === RELATIONSHIPS) {
          // Hard coded for relationships
          // parent id will be same as relationships[0].id
          // so we can use that id for expanding first level tree for relationships
          arr.push(parent.relationships[0].id);
        } else if (view === REGISTRANTS) {
          arr.push(parent.id);
          arr.push(`${parent.id}.1`);
        } else {
          arr.push(parent.id);
          parent.versionBasedData.map((child) => {
            arr.push(`${parent.id}.${child.id}`);
            arr.push(`${parent.id}.${child.id}.1`);
          });
        }
      });
      setExpanded(arr);
    };

    const handleSelect = (event, nodeIds) => {
      if (nodeIds.length >= 0) {
        let selectedIdArr = nodeIds[0].split('.');
        let indx = data.findIndex((item) => item.id === selectedIdArr[0]);

        // Filter object contains current filter applied to data
        // Route will contain filters to support deeplink
        const filter = {
          ...(searchText && { searchText }),
          pagesize: indx + 14,
        };
        handleUpdateSelectedRoute(nodeIds, filter);
        setSelected([0, nodeIds]);
      } else {
        setSelected([]);
      }
    };

    const handleToggle = (event, nodeIds) => {
      setExpanded(nodeIds);
    };

    useEffect(() => {
      // Compare previous data with current data and uuid
      if (
        prevState.data &&
        JSON.stringify(prevState.data) === JSON.stringify(data) &&
        prevState.uuid === selectedItemUUID
      ) {
        return;
      }

      // Update the state with the current data and uuid
      setPrevState({ data, uuid: selectedItemUUID });

      // No data present then return
      if (data.length === 0) {
        return;
      }

      let selectedIdArr = selectedItemUUID.split('.');
      if (selectedIdArr.length >= 0) {
        // Check if showDetailsData data matches with item from route
        // This will prevent unnecessary state update
        if (showDetailsData.data.id !== selectedIdArr[selectedIdArr.length - 1]) {
          const newExpanded = selectedIdArr.reduce(
            (acc, id, index) => [...acc, index > 0 ? `${acc[index - 1]}.${id}` : id],
            [],
          );
          if (JSON.stringify(newExpanded) !== JSON.stringify(expanded)) {
            setExpanded(newExpanded);
          }
          if (JSON.stringify([selectedItemUUID]) !== JSON.stringify(selected)) {
            setSelected([selectedItemUUID]);
          }
          const showData = getFilteredDataForDetailsComponent(data, selectedItemUUID);
          if (JSON.stringify(showData) !== JSON.stringify(showDetailsData)) {
            setShowDetailsData(showData);
          }
        }
      } else {
        setExpanded([]);
        setSelected([]);
        setShowDetailsData({
          type: '',
          data: {},
        });
      }
    }, [selectedItemUUID, data, showDetailsData]);

    useEffect(() => {
      let selectedIdArr = selectedItemUUID.split('.');
      if (selectedIdArr.length >= 0) {
        setTimeout(() => {
          requestAnimationFrame(() => {
            const selectedNode = document.querySelector(`[data-id="${selectedItemUUID}"]`);
            if (selectedNode) {
              selectedNode.scrollIntoView({ behavior: 'smooth' });
            }
          });
        }, 1000);
      }
    }, [view]);

    const disabledExpand = () => {
      return view === COMPONENTS;
    };

    const renderHeader = (type, hasRecords) => (
      <div
        style={{
          display: 'flex',
          flexDirection: 'row',
          justifyContent: 'space-between',
          borderBottom: '1px solid #d2d3d4',
          width: '100%',
          height: '2.55rem',
        }}
      >
        <div>
          {width < 1370 && isSearchExpanded ? null : (
            <div style={{ display: 'flex', flexDirection: 'row' }}>
              <CustomTextTooltip title="Expand All" placement="top">
                {/* span is added to make sure tooltip is not listening to disabled elements to prevent MUI error */}
                <span>
                  <IconButton
                    onClick={expandAll}
                    size="large"
                    disableRipple
                    disabled={disabledExpand()}
                  >
                    <ExpandAllIcon height={17} width={17} />
                  </IconButton>
                </span>
              </CustomTextTooltip>

              <CustomTextTooltip title="Collapse All" placement="top">
                <span>
                  <IconButton
                    onClick={() => setExpanded([])}
                    style={{ paddingLeft: '0rem' }}
                    size="large"
                    disableRipple
                    disabled={disabledExpand()}
                  >
                    <CollapseAllIcon height={17} width={17} />
                  </IconButton>
                </span>
              </CustomTextTooltip>
              {type === MODELS && (
                <>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={checked}
                        onClick={handleChecked}
                        disabled={!hasRecords}
                        inputProps={{ 'aria-label': 'controlled' }}
                        size="small"
                      />
                    }
                    label="Duplicates"
                    sx={{ marginLeft: '2rem' }}
                  />
                  <CustomTextTooltip
                    placement="right"
                    interactive={true}
                    title={`View all duplicate entries of ${_.toLower(
                      view,
                    )}. Entries with identical name and version attributes are considered duplicates. [Learn More](https://docs.meshery.io/concepts/logical/models#models)`}
                    sx={{ margin: '0rem', padding: '0rem' }}
                  >
                    <IconButton>
                      <InfoOutlinedIcon height={20} width={20} />
                    </IconButton>
                  </CustomTextTooltip>
                </>
              )}
            </div>
          )}
        </div>
        <div style={{ display: 'flex' }}>
          <SearchBar
            onSearch={debounce((value) => setSearchText(value), 300)}
            expanded={isSearchExpanded}
            setExpanded={setSearchExpand}
            placeholder="Search"
            value={searchText}
            setModelsFilters={setModelsFilters}
          />
        </div>
      </div>
    );

    const setSearchExpand = (isExpand) => {
      if (!isExpand) {
        setSearchText(() => null);
        setResourcesDetail(() => []);
        setPage({
          Models: 0,
          Components: 0,
          Relationships: 0,
          Registrants: 0,
        });
      }
      setIsSearchExpanded(isExpand);
    };

    const renderTree = (treeComponent, type, isLoading) => (
      <>
        {renderHeader(type, !!data.length)}
        {data.length === 0 && !searchText ? (
          <JustifyAndAlignCenter style={{ height: '27rem' }}>
            {isLoading || (data.length === 0 && !searchText) ? (
              <CircularProgress sx={{ color: Colors.keppelGreen }} />
            ) : (
              <Typography>No {type.toLowerCase()} found</Typography>
            )}
          </JustifyAndAlignCenter>
        ) : data.length === 0 && searchText ? (
          <JustifyAndAlignCenter style={{ height: '27rem' }}>
            <p>No result found</p>
          </JustifyAndAlignCenter>
        ) : (
          <div
            className="scrollElement"
            style={{ overflowY: 'auto', height: '55vh' }}
            onScroll={handleScroll(type)}
          >
            {treeComponent}
          </div>
        )}
      </>
    );

    return (
      <MesheryTreeViewWrapper style={{ width: '100%', height: '100%' }}>
        {view === MODELS &&
          renderTree(
            <MesheryTreeViewModel
              data={data}
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              setShowDetailsData={setShowDetailsData}
              showDetailsData={showDetailsData}
              lastModelRef={lastItemRef[MODELS]}
              isModelFetching={isFetching[MODELS]}
            />,
            MODELS,
            isLoading[view],
          )}
        {view === REGISTRANTS &&
          renderTree(
            <MesheryTreeViewRegistrants
              data={data}
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              setShowDetailsData={setShowDetailsData}
              showDetailsData={showDetailsData}
              lastRegistrantRef={lastItemRef[REGISTRANTS]}
              isRegistrantFetching={isFetching[REGISTRANTS]}
            />,
            REGISTRANTS,
            isLoading[view],
          )}
        {view === COMPONENTS &&
          renderTree(
            <ComponentTree
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              data={data}
              setExpanded={setExpanded}
              setSelected={setSelected}
              setSearchText={setSearchText}
              setShowDetailsData={setShowDetailsData}
              lastComponentRef={lastItemRef[COMPONENTS]}
              isComponentFetching={isFetching[COMPONENTS]}
            />,
            COMPONENTS,
            isLoading[view],
          )}
        {view === RELATIONSHIPS &&
          renderTree(
            <RelationshipTree
              handleToggle={handleToggle}
              handleSelect={handleSelect}
              expanded={expanded}
              selected={selected}
              data={data}
              setShowDetailsData={setShowDetailsData}
              lastRegistrantRef={lastItemRef[RELATIONSHIPS]}
              isRelationshipFetching={isFetching[RELATIONSHIPS]}
            />,
            RELATIONSHIPS,
            isLoading[view],
          )}
      </MesheryTreeViewWrapper>
    );
  },
);
MesheryTreeView.displayName = 'MesheryTreeView';

export default MesheryTreeView;
