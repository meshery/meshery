import React, { useState } from 'react';
import { ArrowBack } from '@material-ui/icons';
import { TooltipIconButton } from '../../utils/TooltipButton';
import { Paper, Typography } from '@material-ui/core';
import { Box, ErrorBoundary, OperatorDataFormatter, useResourceCleanData } from '@layer5/sistent';
import { ALL_VIEW } from './resources/config';
import _ from 'lodash';
import GetNodeIcon from '../configuratorComponents/MeshModel/NodeIcon';
import { JsonParse } from '../../utils/utils';
import { FALLBACK_MESHERY_IMAGE_PATH } from '@/constants/common';
import { iconXLarge } from 'css/icons.styles';
import { getK8sContextFromClusterId } from '@/utils/multi-ctx';
import useKubernetesHook from '../hooks/useKubernetesHook';
import { TootltipWrappedConnectionChip } from '../connections/ConnectionChip';
import ResourceDetailFormatData, { JSONViewFormatter } from './view-component';
import { styled } from '@mui/system';

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
  const { getResourceCleanData } = useResourceCleanData();
  const cleanData = getResourceCleanData({ resource: resource });
  console.log('amit cleanData is here', cleanData);

  const ping = useKubernetesHook();
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
                onClick={() => setView(ALL_VIEW)}
                style={{ color: 'inherit' }}
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

export const Title = ({ onClick, data, value }) => {
  const [isHovered, setHovered] = useState(false);
  return (
    <TitleContainer
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
      style={{ textDecoration: isHovered ? 'underline' : 'none' }}
    >
      <TitleContent onClick={onClick}>
        <div>
          <GetNodeIcon metadata={JsonParse(data)} />
        </div>
        <Typography style={{ marginLeft: '0.50rem' }} variant="body2">
          {value}
        </Typography>
      </TitleContent>
    </TitleContainer>
  );
};
