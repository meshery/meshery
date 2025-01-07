import React from 'react';
import ResponsiveHoneycomb from './ResponsiveHoneycomb';
import Hexagon from './Hexagon';
import { CustomTooltip, ErrorBoundary, Skeleton, Typography } from '@layer5/sistent';
import { useRouter } from 'next/router';
import { componentIcon } from '../charts/utils';
import ConnectCluster from '../charts/ConnectCluster';
import { generateDynamicURL } from '../resources/config';
import {
  HoneycombRoot,
  IconWrapper,
  ResourceCount,
  SelectedHexagon,
  SkeletonHexagon,
} from '../style';

const HoneycombComponent = (props) => {
  const { kinds, isClusterLoading, isClusterIdsEmpty } = props;
  const router = useRouter();

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
        <Typography variant="h6">Cluster Resource Overview</Typography>
        {isClusterLoading || isClusterIdsEmpty ? (
          renderLoadingSkeleton()
        ) : !kinds ? (
          <ConnectCluster message="No workloads found in your cluster(s)." />
        ) : (
          <>
            {Array.isArray(kinds) && kinds.length > 0 && (
              <>
                {!isClusterLoading && (
                  <ResponsiveHoneycomb
                    defaultWidth={1024}
                    size={47}
                    items={kinds}
                    renderItem={(item) => (
                      <Hexagon
                        onClick={() => {
                          router.push(generateDynamicURL(item?.Kind));
                        }}
                      >
                        <SelectedHexagon>
                          <CustomTooltip title={item?.Kind || ''} placement="top">
                            <IconWrapper>
                              <img
                                src={componentIcon({
                                  kind: item?.Kind?.toLowerCase(),
                                  color: 'color',
                                  model: 'kubernetes',
                                })}
                                width="40"
                                height="40"
                                onError={(event) => {
                                  event.target.src = '/static/img/kubernetes.svg';
                                }}
                                alt={item?.Kind || 'Resource Icon'}
                              />
                              <ResourceCount variant="subtitle1">{item.Count}</ResourceCount>
                            </IconWrapper>
                          </CustomTooltip>
                        </SelectedHexagon>
                      </Hexagon>
                    )}
                  />
                )}
              </>
            )}
          </>
        )}
      </HoneycombRoot>
    </ErrorBoundary>
  );
};

export default HoneycombComponent;
