import React, { useCallback, useMemo, useState } from 'react';
import ResponsiveHoneycomb from './ResponsiveHoneycomb';
import Hexagon from './Hexagon';
import {
  ArrowDownwardIcon,
  ArrowUpwardIcon,
  CustomTooltip,
  ErrorBoundary,
  IconButton,
  MenuItem,
  Select,
  Skeleton,
  Typography,
} from '@sistent/sistent';
import { useRouter } from 'next/router';
import ConnectCluster from '../../charts/ConnectCluster';
import { generateDynamicURL } from '../../resources/config';
import {
  HoneycombRoot,
  IconWrapper,
  ResourceCount,
  SelectedHexagon,
  SkeletonHexagon,
  HeaderContainer,
  ControlsContainer,
  NoResourcesText,
} from '../../style';
import {
  DEFAULT_GROUP_BY,
  SORT_DIRECTIONS,
  type ResourceKind,
  type SortDirection,
  useResourceFiltering,
  useResourceOptions,
} from './useResourceOptions';
import GetKubernetesNodeIcon from '../../utils';

type HoneycombComponentProps = {
  kinds?: ResourceKind[];
  isClusterLoading?: boolean;
  isClusterIdsEmpty?: boolean;
  isEditMode?: boolean;
};

type GroupChangeHandler = NonNullable<React.ComponentProps<typeof Select>['onChange']>;

const HONEYCOMB_DEFAULT_WIDTH = 1024;
const HONEYCOMB_SIZE = 47;
const LOADING_SKELETON_COUNT = 40;

const HoneycombComponent = ({
  kinds,
  isClusterLoading,
  isClusterIdsEmpty,
  isEditMode,
}: HoneycombComponentProps) => {
  const router = useRouter();
  const [groupBy, setGroupBy] = useState(DEFAULT_GROUP_BY);
  const [sortDirection, setSortDirection] = useState<SortDirection | null>(null);

  const groupOptions = useResourceOptions();
  const filteredKinds = useResourceFiltering(kinds, groupBy, sortDirection);
  const loadingItems = useMemo<ResourceKind[]>(
    () => Array.from({ length: LOADING_SKELETON_COUNT }, () => ({ Kind: 'loading' })),
    [],
  );

  const handleGroupChange = useCallback<GroupChangeHandler>((event) => {
    setGroupBy(String(event.target.value));
  }, []);

  const handleSortChange = useCallback(() => {
    setSortDirection((previousDirection) =>
      previousDirection === SORT_DIRECTIONS.ASC ? SORT_DIRECTIONS.DESC : SORT_DIRECTIONS.ASC,
    );
  }, []);

  const handleKindClick = useCallback(
    (kind: string) => {
      router.push(generateDynamicURL(kind));
    },
    [router],
  );

  const renderLoadingItem = useCallback(() => {
    return (
      <Hexagon>
        <SkeletonHexagon>
          <Skeleton variant="circular" width={50} height={50} />
        </SkeletonHexagon>
      </Hexagon>
    );
  }, []);

  const renderKind = useCallback(
    (item: ResourceKind) => {
      return (
        <Hexagon onClick={() => handleKindClick(item.Kind)}>
          <SelectedHexagon>
            <CustomTooltip title={item.Kind || ''} placement="top">
              <IconWrapper>
                <GetKubernetesNodeIcon kind={item.Kind} model={item.Model} />
                <ResourceCount variant="subtitle1">{item.Count}</ResourceCount>
              </IconWrapper>
            </CustomTooltip>
          </SelectedHexagon>
        </Hexagon>
      );
    },
    [handleKindClick],
  );

  const hasFilteredKinds = filteredKinds.length > 0;

  return (
    <ErrorBoundary>
      <HoneycombRoot isEditMode={isEditMode}>
        <HeaderContainer>
          <Typography variant="h6" fontWeight="700">
            Cluster Resource Overview
          </Typography>
          <ControlsContainer>
            <Select value={groupBy} onChange={handleGroupChange} variant="standard">
              {groupOptions.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </Select>
            <IconButton size="small" onClick={handleSortChange}>
              <CustomTooltip title="Sort by Count">
                {/*
                 * Sistent's typed-SVG icons are functional components without
                 * forwardRef, so Tooltip's cloneElement warns when used as a
                 * direct child. Wrap them in a span to attach the ref safely.
                 */}
                <span>
                  {sortDirection === SORT_DIRECTIONS.ASC ? (
                    <ArrowUpwardIcon />
                  ) : (
                    <ArrowDownwardIcon />
                  )}
                </span>
              </CustomTooltip>
            </IconButton>
          </ControlsContainer>
        </HeaderContainer>
        {isClusterLoading || isClusterIdsEmpty ? (
          <ResponsiveHoneycomb
            defaultWidth={HONEYCOMB_DEFAULT_WIDTH}
            size={HONEYCOMB_SIZE}
            items={loadingItems}
            renderItem={renderLoadingItem}
          />
        ) : !kinds ? (
          <ConnectCluster message="No workloads found in your cluster(s)." />
        ) : hasFilteredKinds ? (
          <ResponsiveHoneycomb
            defaultWidth={HONEYCOMB_DEFAULT_WIDTH}
            size={HONEYCOMB_SIZE}
            items={filteredKinds}
            renderItem={renderKind}
          />
        ) : (
          <NoResourcesText variant="body1">
            No resources found for the selected group
          </NoResourcesText>
        )}
      </HoneycombRoot>
    </ErrorBoundary>
  );
};

export default HoneycombComponent;
