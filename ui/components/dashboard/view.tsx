import React, { useMemo, useState } from 'react';
import { TooltipIconButton } from '../../utils/TooltipButton';
import {
  ArrowBackIcon as ArrowBack,
  Box,
  ErrorBoundary,
  OperatorDataFormatter,
  Paper,
  styled,
  Tab,
  Tabs,
  Typography,
  useResourceCleanData,
} from '@sistent/sistent';
import { ALL_VIEW } from './resources/config';
import { FALLBACK_MESHERY_IMAGE_PATH } from '@/constants/common';
import { normalizeStaticImagePath } from '@/utils/fallback';
import { iconXLarge } from 'css/icons.styles';
import { getK8sContextFromClusterId } from '@/utils/multi-ctx';
import useKubernetesHook from '@/utils/hooks/useKubernetesHook';
import { TooltipWrappedConnectionChip } from '../connections/ConnectionChip';
import ResourceDetailFormatData, { JSONViewFormatter } from './view-component';
import { useRouter } from 'next/router';
import GetKubernetesNodeIcon from './utils';
import { CONNECTION_STATES } from '@/utils/Enum';
import { useGetConnectionsQuery } from '@/rtk-query/connection';
import SessionPanel from '@/components/sessions/SessionPanel';
import type { SessionKind } from 'lib/sessions/protocol';

const Container = styled('div')({
  margin: '1rem auto',
});

const Header = styled('div')({
  display: 'flex',
  flexDirection: 'row',
  gap: 10,
  justifyItems: 'center',
  textTransform: 'uppercase',
  fontSize: '.9rem',
  flexWrap: 'wrap',
  alignItems: 'center',
  justifyContent: 'space-between',
});

const HeaderLeft = styled('div')({
  display: 'flex',
  gap: '1rem',
  alignItems: 'center',
});

const TitleContainer = styled('div')({
  color: 'inherit',
  cursor: 'pointer',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
});

const TitleContent = styled('div')({
  display: 'inherit',
  alignItems: 'center',
});

/** Fixed-height host for a session pane, so xterm's fit addon has something to measure. */
const SessionPane = styled('div')({
  height: '32rem',
  flexDirection: 'column',
});

type DashboardViewProps = {
  setView: (view: string) => void;
  resource?: Record<string, any> | null;
  k8sConfig: unknown;
};

type DashboardTitleProps = {
  onClick: () => void;
  value: React.ReactNode;
  kind?: string;
  model?: string;
};

/** Resource kinds that host interactive sessions, as MeshSync spells them. */
const SESSION_KINDS: Record<string, string> = { Pod: 'pod' };

type DetailTab = 'details' | SessionKind;

const View = ({ setView, resource, k8sConfig }: DashboardViewProps) => {
  const ping = useKubernetesHook();
  const { getResourceCleanData } = useResourceCleanData();
  const router = useRouter();
  const cleanData = useMemo(
    () => getResourceCleanData({ resource, router }),
    [getResourceCleanData, resource, router],
  );

  const [tab, setTab] = useState<DetailTab>('details');
  // A tab's panel stays mounted once visited, so switching back to Details does
  // not kill a live shell or restart a log stream.
  const [visited, setVisited] = useState<Set<DetailTab>>(() => new Set(['details']));

  const selectTab = (next: DetailTab) => {
    setTab(next);
    setVisited((current) => (current.has(next) ? current : new Set(current).add(next)));
  };

  const { data: connections = [] } = useGetConnectionsQuery({
    page: 0,
    pagesize: 100,
    search: '',
    order: '',
    status: '',
    kind: JSON.stringify(['kubernetes']),
  });

  if (!resource) return null;

  const context = getK8sContextFromClusterId(resource.cluster_id, k8sConfig);
  const connection = connections?.connections.find((conn) => conn.id === context?.connectionId);
  const connectionStatus = connection?.status || CONNECTION_STATES.DISCONNECTED;
  const iconSrc = normalizeStaticImagePath(resource.component_metadata?.styles?.svgColor);

  // Session tabs appear only for resource kinds that can host one. Whether this
  // particular resource currently admits a terminal is resolved server-side by
  // SessionPanel, against live state rather than this cached MeshSync row.
  const connectionId = context?.connectionId;
  const sessionResource = SESSION_KINDS[resource.kind];
  const sessionTarget =
    sessionResource && resource.metadata?.name
      ? {
          resource: sessionResource,
          namespace: resource.metadata.namespace,
          name: resource.metadata.name,
        }
      : null;

  return (
    <Container>
      <Paper>
        <Box sx={{ padding: '1rem 1.5rem' }}>
          <Header>
            <HeaderLeft>
              <TooltipIconButton
                title="Back"
                placement="left"
                onClick={() => {
                  router.back();
                  setView(ALL_VIEW);
                }}
              >
                <ArrowBack />
              </TooltipIconButton>
              <img
                src={iconSrc || FALLBACK_MESHERY_IMAGE_PATH}
                alt={resource?.kind}
                onError={(event: React.SyntheticEvent<HTMLImageElement>) => {
                  event.currentTarget.src = FALLBACK_MESHERY_IMAGE_PATH;
                }}
                {...iconXLarge}
              />
              <Typography variant="h6">{resource?.metadata?.name}</Typography>
            </HeaderLeft>
            <TooltipWrappedConnectionChip
              title={context.name}
              width="100%"
              handlePing={() => ping(context.name, context.server, context.connectionId)}
              status={connectionStatus}
              iconSrc={'/static/img/integrations/kubernetes.svg'}
            />
          </Header>
          {sessionTarget && connectionId ? (
            <Tabs
              value={tab}
              onChange={(_, next: DetailTab) => selectTab(next)}
              sx={{ borderBottom: 1, borderColor: 'divider', marginBottom: '0.5rem' }}
            >
              <Tab value="details" label="Details" sx={{ textTransform: 'none' }} />
              <Tab value="logs" label="Logs" sx={{ textTransform: 'none' }} />
              <Tab value="terminal" label="Terminal" sx={{ textTransform: 'none' }} />
            </Tabs>
          ) : null}

          <Box hidden={tab !== 'details'}>
            <ErrorBoundary>
              <OperatorDataFormatter
                data={cleanData}
                FormatStructuredData={ResourceDetailFormatData}
                ReactJsonFormatter={JSONViewFormatter}
              />
            </ErrorBoundary>
          </Box>

          {sessionTarget && connectionId
            ? (['logs', 'terminal'] as SessionKind[]).map((kind) =>
                visited.has(kind) ? (
                  <SessionPane
                    key={kind}
                    // Hidden, not unmounted: see `visited` above.
                    style={{ display: tab === kind ? 'flex' : 'none' }}
                  >
                    <ErrorBoundary>
                      <SessionPanel
                        connectionId={connectionId}
                        target={sessionTarget}
                        kind={kind}
                      />
                    </ErrorBoundary>
                  </SessionPane>
                ) : null,
              )
            : null}
        </Box>
      </Paper>
    </Container>
  );
};

export default View;

export const Title = ({ onClick, value, kind, model }: DashboardTitleProps) => {
  const [isHovered, setHovered] = useState(false);

  return (
    <TitleContainer
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
      style={{ textDecoration: isHovered ? 'underline' : 'none' }}
    >
      <TitleContent onClick={onClick}>
        <div>
          <GetKubernetesNodeIcon kind={kind} model={model} />
        </div>
        <Typography style={{ marginLeft: '0.50rem' }} variant="body2">
          {value}
        </Typography>
      </TitleContent>
    </TitleContainer>
  );
};
