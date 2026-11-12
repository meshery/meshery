import { Typography, useTheme, styled } from '@sistent/sistent';
import React from 'react';
import AnimatedMeshSync from './Animations/AnimatedMeshSync';
import AnimatedMeshPattern from './Animations/AnimatedMeshPattern';
import AnimatedMeshery from './Animations/AnimatedMeshery';
import AnimatedFilter from './Animations/AnimatedFilter';
import AnimatedLightMeshery from './Animations/AnimatedLightMeshery';

const LoadingContainer = styled('div')({
  textAlign: 'center',
  marginTop: 'calc(50vh - 141px)',
  transform: 'translateY(-50%)',
});

const animatedIconList = {
  AnimatedMeshPattern: <AnimatedMeshPattern style={{ height: '100px', margin: '4px 0px 8px' }} />,
  AnimatedMeshSync: (
    <>
      <AnimatedMeshSync style={{ height: '75px' }} />
      <img
        src="/static/img/meshery-logo/meshery-black.svg"
        alt="mehsery-logo"
        width="125px"
        style={{ margin: '4px 0px 8px' }}
      />
    </>
  ),
  AnimatedFilter: <AnimatedFilter style={{ height: '75px', margin: '4px 0px 8px' }} />,
  AnimatedMeshery: <AnimatedMeshery style={{ height: '100px', margin: '4px 0px 8px' }} />,
};

const animatedLightIconList = {
  AnimatedMeshPattern: <AnimatedMeshPattern style={{ height: '100px', margin: '4px 0px 8px' }} />,
  AnimatedMeshSync: (
    <>
      <AnimatedMeshSync style={{ height: '75px' }} />
      <img
        src="/static/img/meshery-logo/meshery-white.svg"
        alt="mehsery-logo"
        width="125px"
        style={{ margin: '4px 0px 8px' }}
      />
    </>
  ),
  AnimatedFilter: <AnimatedFilter style={{ height: '75px', margin: '4px 0px 8px' }} />,
  AnimatedMeshery: <AnimatedLightMeshery style={{ height: '100px', margin: '4px 0px 8px' }} />,
};

type AnimatedIconKey = keyof typeof animatedIconList;

type LoadingScreenProps = {
  message: string;
  className?: string;
  animatedIcon: AnimatedIconKey;
} & React.HTMLAttributes<HTMLDivElement>;

function LoadingScreen({ message, className, animatedIcon, ...other }: LoadingScreenProps) {
  const theme = useTheme();

  return (
    <LoadingContainer data-testid="page-loader" className={className} {...other}>
      {theme.palette.mode === 'light'
        ? animatedIconList[animatedIcon]
        : animatedLightIconList[animatedIcon]}
      <Typography variant="caption" component="div">
        {message}
      </Typography>
    </LoadingContainer>
  );
}

export default LoadingScreen;
