// @ts-nocheck
import React from 'react';
import { cursorNotAllowed } from '../../../css/disableComponent.styles';
import { DASHBOARD } from '../../../constants/navigator';
import {
  MainLogo,
  MainLogoCollapsed,
  MainLogoText,
  MainLogoTextCollapsed,
  StyledListItem,
} from '../../General/style';

const NavigatorHeader = ({ capabilitiesRegistryObj, isDrawerCollapsed, handleTitleClick }) => {
  return (
    <div
      style={
        !capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD]) ? cursorNotAllowed : {}
      }
    >
      <>
        <StyledListItem
          component="a"
          onClick={handleTitleClick}
          disableLogo={!capabilitiesRegistryObj?.isNavigatorComponentEnabled([DASHBOARD])}
        >
          {isDrawerCollapsed ? (
            <>
              <MainLogoCollapsed src="/static/img/meshery-logo.png" onClick={handleTitleClick} />
              <MainLogoTextCollapsed
                src="/static/img/meshery-logo-text.png"
                onClick={handleTitleClick}
              />
            </>
          ) : (
            <>
              <MainLogo src="/static/img/meshery-logo.png" onClick={handleTitleClick} />
              <MainLogoText src="/static/img/meshery-logo-text.png" onClick={handleTitleClick} />
            </>
          )}
        </StyledListItem>
      </>
    </div>
  );
};

export default NavigatorHeader;
