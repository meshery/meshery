import React, { useState, useMemo } from 'react';
import ResponsiveHoneycomb from './ResponsiveHoneycomb';
import Hexagon from './Hexagon';
import { makeStyles } from '@material-ui/core';
import { CustomTooltip, ErrorBoundary, Skeleton, Typography } from '@layer5/sistent';
import { ResourceSelector } from '../charts/ResourceSelector';
import { useRouter } from 'next/router';
import { componentIcon } from '../charts/utils';
import ConnectCluster from '../charts/ConnectCluster';
import { generateDynamicURL } from '../resources/config';

const useStyles = makeStyles((theme) => ({
  rootContainer: {
    backgroundColor: theme.palette.secondary.elevatedComponents,
    padding: theme.spacing(2),
    borderRadius: 4,
    width: '100%',
  },
  selected: {
    display: 'flex',
    height: '95%',
    background: theme.palette.secondary.honeyComb,
    justifyContent: 'center',
    alignItems: 'center',
    '&:hover': {
      cursor: 'pointer',
    },
  },
  unselected: {
    display: 'flex',
    height: '95%',
    background: theme.palette.secondary.honeyComb,
    justifyContent: 'center',
    alignItems: 'center',
    opacity: 0.5,
    '&:hover': {
      cursor: 'pointer',
    },
  },
  sectionHeading: {
    backgroundColor: theme.palette.type === 'dark' ? '#252E31' : 'rgba(57, 102, 121, .1)',
  },
  sectionHeadingText: {
    color: theme.palette.secondary.titleText,
  },
  mainSection: {
    backgroundColor: theme.palette.secondary.sidebar,
  },
  skeletonCell: {
    display: 'flex',
    height: '95%',
    background: theme.palette.secondary.honeyComb,
    justifyContent: 'center',
    alignItems: 'center',
    opacity: 0.5,
  },
}));

const HoneycombComponent = (props) => {
  const classes = useStyles();
  const { kinds, isClusterLoading } = props;
  const [selectedKind, setSelectedKind] = useState('all');

  const onKindChange = (kind) => {
    setSelectedKind(kind);
  };

  const filteredKinds = useMemo(() => {
    if (selectedKind === 'all') return kinds;
    return kinds.filter((kind) => kind.Kind === selectedKind);
  }, [selectedKind, kinds]);

  const router = useRouter();

  const renderLoadingSkeleton = () => {
    const loadingItems = Array(40).fill({ Kind: 'loading' });
    return (
      <ResponsiveHoneycomb
        defaultWidth={1024}
        size={47}
        items={loadingItems}
        renderItem={() => (
          <Hexagon className={classes.skeletonCell}>
            <Skeleton variant="circular" width={50} height={50} />
          </Hexagon>
        )}
      />
    );
  };
  return (
    <ErrorBoundary>
      <div className={classes.rootContainer}>
        <div className={classes.mainSection}>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <Typography variant="h6">Cluster Resource Overview</Typography>
            {!isClusterLoading && kinds && (
              <ResourceSelector
                kinds={kinds?.map((kind) => kind?.Kind)}
                selectedkind={selectedKind}
                onKindChange={onKindChange}
              />
            )}
          </div>
          {isClusterLoading ? (
            renderLoadingSkeleton()
          ) : !kinds ? (
            <ConnectCluster />
          ) : (
            <>
              {Array.isArray(kinds) && kinds.length > 0 && (
                <>
                  {!isClusterLoading && (
                    <ResponsiveHoneycomb
                      defaultWidth={1024}
                      size={47}
                      items={filteredKinds}
                      renderItem={(item) => (
                        <Hexagon
                          className={classes.selected}
                          onClick={() => {
                            router.push(generateDynamicURL(item?.Kind));
                          }}
                        >
                          <CustomTooltip title={item?.Kind || ''} placement="top">
                            <div
                              style={{
                                marginTop: '20px',
                                display: 'flex',
                                flexDirection: 'column',
                              }}
                            >
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
                              <Typography
                                style={{
                                  margin: '0px auto',
                                }}
                                variant="subtitle1"
                              >
                                {item.Count}
                              </Typography>
                            </div>
                          </CustomTooltip>
                        </Hexagon>
                      )}
                    />
                  )}
                </>
              )}
            </>
          )}
        </div>
      </div>
    </ErrorBoundary>
  );
};

export default HoneycombComponent;
