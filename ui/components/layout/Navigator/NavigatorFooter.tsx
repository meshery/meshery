// @ts-nocheck
import React from 'react';
import {
  Grow,
  ListItem,
  LeftArrowIcon,
  useTheme,
  HelpOutlinedIcon,
  OpenInNewIcon,
} from '@sistent/sistent';
import { cursorNotAllowed, disabledStyle } from '../../../css/disableComponent.styles';
import { TOGGLER } from '../../../constants/navigator';
import { iconSmall } from '../../../css/icons.styles';
import { CustomTextTooltip } from '../../MesheryMeshInterface/PatternService/CustomTextTooltip';
import {
  ListIconSide,
  NavigatorHelpIcons,
  HelpListItem,
  HelpButton,
  ChevronButtonWrapper,
  FixedSidebarFooter,
} from '../../General/style';

const NavigatorFooter = ({
  isDrawerCollapsed,
  state,
  externlinks,
  toggleSpacing,
  toggleMiniDrawer,
}) => {
  const theme = useTheme();
  const [hoveredHelpIcon, setHoveredHelpIcon] = React.useState(null);

  const getMesheryVersionText = () => {
    const { build, outdated, release_channel } = state.versionDetail;
    if (release_channel === 'edge' && outdated) return `${build}`;
    if (release_channel === 'edge' && !outdated) return `${release_channel}-latest`;
    if (release_channel === 'stable') return `${release_channel}-${build}`;
    return `${build}`;
  };

  const versionUpdateMsg = () => {
    const { outdated, latest } = state.versionDetail;
    if (outdated)
      return (
        <span style={{ marginLeft: '15px' }}>
          {'Update available '}
          <a
            href={`https://docs.meshery.io/project/releases/${latest}`}
            target="_blank"
            rel="noreferrer"
            style={{ color: 'white' }}
          >
            <CustomTextTooltip
              title={`Newer version of Meshery available: ${latest}`}
              placement="right"
            >
              <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
            </CustomTextTooltip>
          </a>
        </span>
      );
    return <span style={{ marginLeft: '15px' }}>Running latest</span>;
  };

  const openReleaseNotesInNew = () => {
    const { release_channel, build } = state.versionDetail;
    if (release_channel === 'edge')
      return (
        <a
          href="https://docs.meshery.io/project/releases"
          target="_blank"
          rel="noreferrer"
          style={{ color: 'white' }}
        >
          <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
        </a>
      );
    return (
      <a
        href={`https://docs.meshery.io/project/releases/${build}`}
        target="_blank"
        rel="noreferrer"
        style={{ color: 'white' }}
      >
        <OpenInNewIcon style={{ width: '0.85rem', verticalAlign: 'middle' }} />
      </a>
    );
  };

  const HelpIcons = (
    <>
      <NavigatorHelpIcons
        isCollapsed={isDrawerCollapsed}
        size="large"
        orientation={isDrawerCollapsed ? 'vertical' : 'horizontal'}
      >
        {externlinks.map(({ id, icon, hovericon, title, href }, index) => {
          const isHovered = hoveredHelpIcon === id;
          return (
            <HelpListItem
              key={id}
              style={isDrawerCollapsed && !state.showHelperButton ? { display: 'none' } : {}}
              onMouseEnter={() => setHoveredHelpIcon(id)}
              onMouseLeave={() => setHoveredHelpIcon(null)}
            >
              <Grow
                in={state.showHelperButton || !isDrawerCollapsed}
                timeout={{ enter: 600 - index * 200, exit: 100 * index }}
              >
                <a
                  href={href}
                  target="_blank"
                  rel="noreferrer"
                  style={
                    isDrawerCollapsed
                      ? { display: 'flex', alignItems: 'center', justifyContent: 'center' }
                      : {}
                  }
                >
                  <CustomTextTooltip title={title} placement={isDrawerCollapsed ? 'right' : 'top'}>
                    <ListIconSide>{isHovered && hovericon ? hovericon : icon}</ListIconSide>
                  </CustomTextTooltip>
                </a>
              </Grow>
            </HelpListItem>
          );
        })}
        <ListItem key="help-button" style={{ display: isDrawerCollapsed ? 'inherit' : 'none' }}>
          <CustomTextTooltip title="Help" placement={isDrawerCollapsed ? 'right' : 'top'}>
            <HelpButton isCollapsed={isDrawerCollapsed} onClick={toggleSpacing}>
              <HelpOutlinedIcon
                style={{
                  fontSize: '1.45rem',
                  ...iconSmall,
                  color: theme.palette.background.constant.white,
                  opacity: '0.7',
                  transition: 'opacity 200ms linear',
                  '&:hover': {
                    opacity: 1,
                    background: 'transparent',
                  },
                  '&:focus': {
                    opacity: 1,
                    background: 'transparent',
                  },
                }}
              />
            </HelpButton>
          </CustomTextTooltip>
        </ListItem>
      </NavigatorHelpIcons>
    </>
  );

  const Version = (
    <ListItem
      style={{
        position: 'sticky',
        paddingLeft: 0,
        paddingRight: 0,
        color: theme.palette.background.constant.white,
        fontSize: '0.75rem',
      }}
    >
      {isDrawerCollapsed ? (
        <div style={{ textAlign: 'center', width: '100%' }}>{state.versionDetail.build}</div>
      ) : (
        <Grow
          in={!isDrawerCollapsed}
          timeout={{ enter: 800, exit: 100 }}
          style={{ textAlign: 'center', width: '100%' }}
        >
          <span>
            {getMesheryVersionText()} {'  '}
            <span style={{ cursor: 'pointer' }}>{openReleaseNotesInNew()}</span>
            {versionUpdateMsg()}
          </span>
        </Grow>
      )}
    </ListItem>
  );

  const Chevron = (
    <ChevronButtonWrapper
      isCollapsed={isDrawerCollapsed}
      style={
        state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
          ? {}
          : cursorNotAllowed
      }
    >
      <div
        style={
          state?.capabilitiesRegistryObj?.isNavigatorComponentEnabled?.([TOGGLER])
            ? {}
            : disabledStyle
        }
        onClick={toggleMiniDrawer}
      >
        <LeftArrowIcon
          alt="Sidebar collapse toggle"
          style={{
            cursor: 'pointer',
            verticalAlign: 'middle',
          }}
          fill={theme.palette.icon.default}
          stroke={theme.palette.icon.default}
          width="1.2rem"
          height="2.8rem"
        />
      </div>
    </ChevronButtonWrapper>
  );

  return (
    <FixedSidebarFooter>
      {Chevron}
      {HelpIcons}
      {Version}
    </FixedSidebarFooter>
  );
};

export default NavigatorFooter;
