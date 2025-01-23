import React, { useState, useCallback } from 'react';
import ResponsiveHoneycomb from './ResponsiveHoneycomb';
import Hexagon from './Hexagon';
import {
  CustomTooltip,
  ErrorBoundary,
  Skeleton,
  Typography,
  Select,
  MenuItem,
} from '@layer5/sistent';
import { useRouter } from 'next/router';
import ConnectCluster from '../charts/ConnectCluster';
import { generateDynamicURL } from '../resources/config';
import {
  HoneycombRoot,
  IconWrapper,
  ResourceCount,
  SelectedHexagon,
  SkeletonHexagon,
  HeaderContainer,
  ControlsContainer,
  NoResourcesText,
  StyledIconButton,
} from '../style';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import { useResourceOptions, useResourceFiltering, SORT_DIRECTIONS } from './useResourceOptions';
import GetKubernetesNodeIcon from '../utils';

const HoneycombComponent = (props) => {
  const { kinds, isClusterLoading, isClusterIdsEmpty } = props;
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
      <HoneycombRoot>
        <HeaderContainer>
          <Typography variant="h6">Cluster Resource Overview</Typography>
          <ControlsContainer>
            <Select value={groupBy} onChange={handleGroupChange} size="small">
              {groupOptions.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </Select>
            <StyledIconButton size="small" onClick={handleSortChange}>
              <CustomTooltip title={`Sort by Count`}>
                {sortDirection === SORT_DIRECTIONS.ASC ? (
                  <ArrowUpwardIcon />
                ) : (
                  <ArrowDownwardIcon />
                )}
              </CustomTooltip>
            </StyledIconButton>
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
