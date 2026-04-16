// @ts-nocheck
import React, { useState, useCallback } from 'react';
import ResponsiveHoneycomb from './ResponsiveHoneycomb';
import Hexagon from './Hexagon';
import {
  CustomTooltip,
  ErrorBoundary,
  SvgIcon,
  Skeleton,
  Typography,
  Select,
  MenuItem,
  IconButton,
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
import { useResourceOptions, useResourceFiltering, SORT_DIRECTIONS } from './useResourceOptions';
import GetKubernetesNodeIcon from '../../utils';

const ArrowUpwardIcon = React.forwardRef((props, ref) => (
  <SvgIcon {...props} ref={ref}>
    <path d="m4 12 1.41 1.41L11 7.83V20h2V7.83l5.58 5.59L20 12l-8-8z" />
  </SvgIcon>
));

ArrowUpwardIcon.displayName = 'ArrowUpwardIcon';

const ArrowDownwardIcon = React.forwardRef((props, ref) => (
  <SvgIcon {...props} ref={ref}>
    <path d="m20 12-1.41-1.41L13 16.17V4h-2v12.17l-5.58-5.59L4 12l8 8z" />
  </SvgIcon>
));

ArrowDownwardIcon.displayName = 'ArrowDownwardIcon';

const HoneycombComponent = (props) => {
  const { kinds, isClusterLoading, isClusterIdsEmpty, isEditMode } = props;
  const router = useRouter();
  const [groupBy, setGroupBy] = useState('all');
  const [sortDirection, setSortDirection] = useState(null);

  const groupOptions = useResourceOptions();
  const filteredKinds = useResourceFiltering(kinds, groupBy, sortDirection);
  const handleGroupChange = useCallback((e) => {
    setGroupBy(e.target.value);
  }, []);

  const handleSortChange = useCallback(() => {
    setSortDirection((prev) =>
      prev === SORT_DIRECTIONS.ASC ? SORT_DIRECTIONS.DESC : SORT_DIRECTIONS.ASC,
    );
  }, []);

  const renderLoadingSkeleton = () => {
    const loadingItems = Array(40).fill({ Kind: 'loading' });
    return (
      <ResponsiveHoneycomb
        defaultWidth={1024}
        size={47}
        items={loadingItems}
        renderItem={() => (
          <Hexagon>
            <SkeletonHexagon>
              <Skeleton variant="circular" width={50} height={50} />
            </SkeletonHexagon>
          </Hexagon>
        )}
      />
    );
  };

  return (
    <ErrorBoundary>
      <HoneycombRoot isEditMode={isEditMode}>
        <HeaderContainer>
          <Typography variant="h6" fontWeight={'700'}>
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
              <CustomTooltip title={`Sort by Count`}>
                {sortDirection === SORT_DIRECTIONS.ASC ? (
                  <ArrowUpwardIcon />
                ) : (
                  <ArrowDownwardIcon />
                )}
              </CustomTooltip>
            </IconButton>
          </ControlsContainer>
        </HeaderContainer>
        {isClusterLoading || isClusterIdsEmpty ? (
          renderLoadingSkeleton()
        ) : !kinds ? (
          <ConnectCluster message="No workloads found in your cluster(s)." />
        ) : (
          <>
            {Array.isArray(filteredKinds) && filteredKinds.length > 0 ? (
              <ResponsiveHoneycomb
                defaultWidth={1024}
                size={47}
                items={filteredKinds}
                renderItem={(item) => {
                  return (
                    <Hexagon
                      onClick={() => {
                        router.push(generateDynamicURL(item?.Kind));
                      }}
                    >
                      <SelectedHexagon>
                        <CustomTooltip title={item?.Kind || ''} placement="top">
                          <IconWrapper>
                            <GetKubernetesNodeIcon kind={item?.Kind} model={item?.Model} />
                            <ResourceCount variant="subtitle1">{item.Count}</ResourceCount>
                          </IconWrapper>
                        </CustomTooltip>
                      </SelectedHexagon>
                    </Hexagon>
                  );
                }}
              />
            ) : (
              <NoResourcesText variant="body1">
                No resources found for the selected group
              </NoResourcesText>
            )}
          </>
        )}
      </HoneycombRoot>
    </ErrorBoundary>
  );
};

export default HoneycombComponent;
