import React, { useMemo, useState } from 'react';
import { TooltipIconButton } from '../../utils/TooltipButton';
import {
  ArrowBackIcon as ArrowBack,
  Box,
  ErrorBoundary,
  OperatorDataFormatter,
  Paper,
  styled,
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

const View = ({ setView, resource, k8sConfig }: DashboardViewProps) => {
  const ping = useKubernetesHook();
  const { getResourceCleanData } = useResourceCleanData();
  const router = useRouter();
  const cleanData = useMemo(
    () => getResourceCleanData({ resource, router }),
    [getResourceCleanData, resource, router],
  );

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
  const connection = connections?.connections?.find((conn) => conn.id === context?.connectionId);
  const connectionStatus = connection?.status || CONNECTION_STATES.DISCONNECTED;
  const iconSrc = normalizeStaticImagePath(resource.component_metadata?.styles?.svgColor);

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
          <ErrorBoundary>
            <OperatorDataFormatter
              data={cleanData}
              FormatStructuredData={ResourceDetailFormatData}
              ReactJsonFormatter={JSONViewFormatter}
            />
          </ErrorBoundary>
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
