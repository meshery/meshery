import React, { useEffect, useState, useRef, useCallback, useMemo } from 'react';
import {
  IconButton,
  FormControlLabel,
  Switch,
  CircularProgress,
  Typography,
} from '@layer5/sistent';
import { MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../../constants/navigator';
import SearchBar from '../../utils/custom-search';
import debounce from '../../utils/debounce';
import { useWindowDimensions } from '../../utils/dimension';
import { getFilteredDataForDetailsComponent } from './helper';
import { CustomTextTooltip } from '../MesheryMeshInterface/PatternService/CustomTextTooltip';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import _ from 'lodash';
import CollapseAllIcon from '@/assets/icons/CollapseAll';
import ExpandAllIcon from '@/assets/icons/ExpandAll';
import { Colors } from '../../themes/app';
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
    lastItemRef,
    isFetching,
    isLoading,
  }) => {
    const { handleUpdateSelectedRoute, selectedItemUUID } = useRegistryRouter();
    const [treeState, setTreeState] = useState({
      expanded: [],
      selected: [],
      isSearchExpanded: !!searchText,
    });
    const { width } = useWindowDimensions();
    const [prevState, setPrevState] = useState({ data: [], uuid: '' });
    const scrollRef = useRef();

    const treeData = useMemo(() => {
      return data;
    }, [data]);

    // Optimized scroll handler with throttling
    const handleScroll = useCallback(
      _.throttle(
        (scrollingView) => (event) => {
          const { scrollTop, scrollHeight, clientHeight } = event.target;
          if (scrollHeight - scrollTop <= clientHeight + 1) {
            setPage((prev) => ({
              ...prev,
              [scrollingView]: prev[scrollingView] + 1,
            }));
          }
        },
        150,
      ),
      [setPage],
    );

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
      setTreeState((prev) => ({ ...prev, expanded: arr }));
    };

    const handleSelect = useCallback(
      (event, nodeIds) => {
        if (!nodeIds.length) {
          setTreeState((prev) => ({ ...prev, selected: [] }));
          return;
        }
        const selectedId = nodeIds[0];
        // Filter object contains current filter applied to data
        // Route will contain filters to support deeplink
        const filter = {
          ...(searchText && { searchText }),
          pagesize: 25,
        };
        handleUpdateSelectedRoute(nodeIds, filter);

        updateSelection(selectedId, data, setShowDetailsData);
        setTreeState((prev) => ({ ...prev, selected: [selectedId] }));
      },
      [data, setShowDetailsData, handleUpdateSelectedRoute, searchText],
    );

    const handleToggle = useCallback((event, nodeIds) => {
      setTreeState((prev) => ({ ...prev, expanded: nodeIds }));
    }, []);

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
          if (JSON.stringify(newExpanded) !== JSON.stringify(treeState.expanded)) {
            setTreeState((prev) => ({ ...prev, expanded: newExpanded }));
          }
          if (JSON.stringify([selectedItemUUID]) !== JSON.stringify(treeState.selected)) {
            setTreeState((prev) => ({ ...prev, selected: [selectedItemUUID] }));
          }
          const showData = getFilteredDataForDetailsComponent(data, selectedItemUUID);
          if (JSON.stringify(showData) !== JSON.stringify(showDetailsData)) {
            setShowDetailsData(showData);
          }
        }
      } else {
        setTreeState({
          expanded: [],
          selected: [],
          isSearchExpanded: !!searchText,
        });
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
          {width < 1370 && treeState.isSearchExpanded ? null : (
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
                    onClick={() => setTreeState((prev) => ({ ...prev, expanded: [] }))}
                    style={{ marginRight: '4px' }}
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
                    label="Show Duplicates"
                    style={{ margin: 0 }}
                  />
                  <CustomTextTooltip
                    placement="right"
                    interactive={true}
                    title={`View all duplicate entries of ${_.toLower(
                      view,
                    )}. Entries with identical name and version attributes are considered duplicates. [Learn More](https://docs.meshery.io/concepts/logical/models#models)`}
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
            onSearch={handleSearch}
            expanded={treeState.isSearchExpanded}
            setExpanded={setSearchExpand}
            placeholder="Search"
            value={searchText || ''}
          />
        </div>
      </div>
    );

    const setSearchExpand = useCallback(
      (isExpand) => {
        if (!isExpand) {
          setSearchText('');
          setResourcesDetail([]);
          setPage({
            Models: 0,
            Components: 0,
            Relationships: 0,
            Registrants: 0,
          });
        }
        setTreeState((prev) => ({ ...prev, isSearchExpanded: isExpand }));
      },
      [setSearchText, setResourcesDetail, setPage],
    );

    const handleSearch = useCallback(
      debounce((value) => {
        const searchValue = value || '';
        setSearchText(searchValue);
        if (searchValue) {
          setTreeState((prev) => ({
            ...prev,
            selected: [],
            expanded: [],
          }));
        }
      }, 300),
      [],
    );

    useEffect(() => {
      if (!selectedItemUUID || data.length === 0) {
        return;
      }

      const selectedIdArr = selectedItemUUID.split('.');
      if (selectedIdArr.length > 0) {
        const newExpanded = selectedIdArr.reduce(
          (acc, id, index) => [...acc, index > 0 ? `${acc[index - 1]}.${id}` : id],
          [],
        );

        setTreeState((prev) => ({
          ...prev,
          expanded: newExpanded,
          selected: [selectedItemUUID],
        }));

        const showData = getFilteredDataForDetailsComponent(data, selectedItemUUID);
        if (JSON.stringify(showData) !== JSON.stringify(showDetailsData)) {
          setShowDetailsData(showData);
        }
      }
    }, [selectedItemUUID, data]);

    const renderTree = (type, { data, handlers, state, refs, loading }) => (
      <>
        {renderHeader(type, !!data.length)}
        {data.length === 0 && !searchText ? (
          <JustifyAndAlignCenter style={{ height: '27rem' }}>
            {loading.isLoading || (data.length === 0 && !searchText) ? (
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
            onScroll={handlers.handleScroll}
          >
            {type === MODELS && (
              <MesheryTreeViewModel
                data={data}
                handleToggle={handlers.handleToggle}
                handleSelect={handlers.handleSelect}
                expanded={state.expanded}
                selected={state.selected}
                setShowDetailsData={setShowDetailsData}
                showDetailsData={showDetailsData}
                lastModelRef={refs[MODELS]}
                isModelFetching={loading.isFetching}
              />
            )}
            {type === REGISTRANTS && (
              <MesheryTreeViewRegistrants
                data={data}
                handleToggle={handlers.handleToggle}
                handleSelect={handlers.handleSelect}
                expanded={state.expanded}
                selected={state.selected}
                setShowDetailsData={setShowDetailsData}
                showDetailsData={showDetailsData}
                lastRegistrantRef={refs[REGISTRANTS]}
                isRegistrantFetching={loading.isFetching}
              />
            )}
            {type === COMPONENTS && (
              <ComponentTree
                handleToggle={handlers.handleToggle}
                handleSelect={handlers.handleSelect}
                expanded={state.expanded}
                selected={state.selected}
                data={data}
                setExpanded={setTreeState}
                setSelected={setTreeState}
                setSearchText={setSearchText}
                setShowDetailsData={setShowDetailsData}
                lastComponentRef={refs[COMPONENTS]}
                isComponentFetching={loading.isFetching}
              />
            )}
            {type === RELATIONSHIPS && (
              <RelationshipTree
                handleToggle={handlers.handleToggle}
                handleSelect={handlers.handleSelect}
                expanded={state.expanded}
                selected={state.selected}
                data={data}
                setShowDetailsData={setShowDetailsData}
                lastRegistrantRef={refs[RELATIONSHIPS]}
                isRelationshipFetching={loading.isFetching}
              />
            )}
          </div>
        )}
      </>
    );

    return (
      <MesheryTreeViewWrapper>
        {renderTree(view, {
          data: treeData,
          handlers: {
            handleToggle,
            handleSelect,
            handleScroll: handleScroll(view),
          },
          state: treeState,
          refs: lastItemRef,
          loading: {
            isFetching: isFetching[view],
            isLoading: isLoading[view],
          },
        })}
      </MesheryTreeViewWrapper>
    );
  },
);

const updateSelection = (selectedId, data, setShowDetailsData) => {
  const showData = getFilteredDataForDetailsComponent(data, selectedId);
  setShowDetailsData(showData);
};

MesheryTreeView.displayName = 'MesheryTreeView';

export default MesheryTreeView;
