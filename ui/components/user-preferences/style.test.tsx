import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => {
  const styled = (Component: any) => () => {
    const StyledComponent = ({ children, ...props }: any) => {
      if (typeof Component === 'string') {
        return React.createElement(Component, props, children);
      }
      return React.createElement(Component, props, children);
    };
    StyledComponent.displayName = 'StyledSistentMock';
    return StyledComponent;
  };

  const make = (tag: string) => (props: any) => React.createElement(tag, props, props.children);

  return {
    styled,
    Card: make('div'),
    FormControl: make('div'),
    FormGroup: make('div'),
    Paper: make('div'),
    Select: make('select'),
    Grid2: make('div'),
  };
});

import {
  StyledSelect,
  SelectItem,
  OrgText,
  OrgIconContainer,
  FormGroupWrapper,
  FormContainerWrapper,
  StatsWrapper,
  PaperRoot,
  RootContainer,
  ProviderCard,
  StyledPaper,
  BoxWrapper,
  GridCapabilityHeader,
  GridExtensionHeader,
  GridExtensionItem,
  Divider,
  TabLabel,
  StyledFormGroup,
  FormLegend,
  FormLegendSmall,
  IconStyled,
  IconText,
  BackToPlay,
  LinkStyled,
  HideScrollbar,
  FormContainer,
  PreferenceLabel,
  PreferenceGroup,
  PreferenceValue,
  PreferenceSection,
  ListItemStyled,
  ContentWrapper,
} from './style';

describe('user-preferences/style', () => {
  it('exports all styled components', () => {
    const exports = [
      StyledSelect,
      SelectItem,
      OrgText,
      OrgIconContainer,
      FormGroupWrapper,
      FormContainerWrapper,
      StatsWrapper,
      PaperRoot,
      RootContainer,
      ProviderCard,
      StyledPaper,
      BoxWrapper,
      GridCapabilityHeader,
      GridExtensionHeader,
      GridExtensionItem,
      Divider,
      TabLabel,
      StyledFormGroup,
      FormLegend,
      FormLegendSmall,
      IconStyled,
      IconText,
      BackToPlay,
      LinkStyled,
      HideScrollbar,
      FormContainer,
      PreferenceLabel,
      PreferenceGroup,
      PreferenceValue,
      PreferenceSection,
      ListItemStyled,
      ContentWrapper,
    ];

    for (const c of exports) {
      expect(c).toBeDefined();
    }
  });

  it('renders a small subset of style components without crashing', () => {
    const { container } = render(
      <RootContainer>
        <PreferenceSection>
          <PreferenceLabel>label</PreferenceLabel>
          <PreferenceValue>value</PreferenceValue>
          <PreferenceGroup>
            <FormContainer>
              <Divider />
              <TabLabel>tab</TabLabel>
              <IconStyled>icon</IconStyled>
              <IconText>icon text</IconText>
              <BackToPlay>back</BackToPlay>
              <LinkStyled>link</LinkStyled>
              <HideScrollbar>hidden</HideScrollbar>
              <FormLegend>legend</FormLegend>
              <FormLegendSmall>legend small</FormLegendSmall>
              <ContentWrapper>content</ContentWrapper>
              <BoxWrapper>boxwrap</BoxWrapper>
              <ListItemStyled>listitem</ListItemStyled>
              <SelectItem>selitem</SelectItem>
              <OrgText>orgtext</OrgText>
              <OrgIconContainer>orgicon</OrgIconContainer>
            </FormContainer>
          </PreferenceGroup>
        </PreferenceSection>
      </RootContainer>,
    );

    const text = container.textContent || '';
    expect(text).toContain('label');
    expect(text).toContain('value');
    expect(text).toContain('tab');
    expect(text).toContain('legend');
    expect(text).toContain('content');
  });
});
