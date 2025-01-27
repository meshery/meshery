import React, { useState } from 'react';
import { ArrowBack } from '@mui/icons-material';
import { TooltipIconButton } from '../../utils/TooltipButton';
import { Paper, Typography } from '@material-ui/core';
import { Box, ErrorBoundary, OperatorDataFormatter, useResourceCleanData } from '@layer5/sistent';
import { ALL_VIEW } from './resources/config';
import { FALLBACK_MESHERY_IMAGE_PATH } from '@/constants/common';
import { iconXLarge } from 'css/icons.styles';
import { getK8sContextFromClusterId } from '@/utils/multi-ctx';
import useKubernetesHook from '../hooks/useKubernetesHook';
import { TootltipWrappedConnectionChip } from '../connections/ConnectionChip';
import ResourceDetailFormatData, { JSONViewFormatter } from './view-component';
import { styled } from '@mui/system';
import { useRouter } from 'next/router';
import GetKubernetesNodeIcon from './utils';

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

const View = (props) => {
  const { setView, resource, k8sConfig } = props;
  const ping = useKubernetesHook();
  const { getResourceCleanData } = useResourceCleanData();
  const router = useRouter();
  const cleanData = getResourceCleanData({ resource: resource, router: router });
  if (!resource) return null;
  const context = getK8sContextFromClusterId(resource.cluster_id, k8sConfig);

  return (
    <Container>
      <Paper>
        <Box padding={'1rem 1.5rem'}>
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
                src={`/${resource.component_metadata?.styles?.svgColor}`}
                alt={resource?.kind}
                onError={(e) => {
                  e.currentTarget.src = FALLBACK_MESHERY_IMAGE_PATH;
                }}
                {...iconXLarge}
              />
              <Typography variant="h6">{resource?.metadata?.name}</Typography>
            </HeaderLeft>
            <TootltipWrappedConnectionChip
              title={context.name}
              width="100%"
              handlePing={() => ping(context.name, context.server, context.connection_id)}
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

export const Title = ({ onClick, value, kind, model }) => {
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
