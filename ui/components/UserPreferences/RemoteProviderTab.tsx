import React, { useState } from 'react';
import {
  Typography,
  Grid2,
  CardContent,
  CardHeader,
  CustomTooltip,
  NoSsr,
  IconButton,
  ErrorBoundary,
} from '@sistent/sistent';
import CopyIcon from '../../assets/icons/CopyIcon';
import ExtensionSandbox from '../ExtensionSandbox';
import RemoteComponent from '../RemoteComponent';
import { SecondaryTab, SecondaryTabs } from '../Dashboard/style';
import { useWindowDimensions } from '@/utils/dimension';
import {
  StatsWrapper,
  ProviderCard,
  RootContainer,
  BoxWrapper,
  Divider,
  GridCapabilityHeader,
  GridExtensionHeader,
  GridExtensionItem,
  TabLabel,
  HideScrollbar,
} from './style';

function convertToTitleCase(str: string): string {
  const words = str.split('_');
  for (let i = 0; i < words.length; i++) {
    words[i] = words[i].charAt(0).toUpperCase() + words[i].slice(1);
  }
  return words.join(' ');
}

interface Capability {
  feature: string;
  endpoint: string;
}

interface ExtensionEntry {
  component: string;
  type: string;
}

interface ProviderInfoData {
  provider_description?: string[];
  capabilities?: Capability[];
  extensions?: Record<string, ExtensionEntry[]>;
  [key: string]: unknown;
}

interface ProviderInfoTabProps {
  providerInfo: ProviderInfoData;
  userData?: {
    remoteProviderPreferences?: {
      theme?: string;
    };
  };
}

const ProviderInfoDetails: React.FC<ProviderInfoTabProps> = ({ providerInfo, userData }) => {
  const [copied, setCopied] = useState(false);
  const copyToClipboard = (text: string) => {
    navigator.clipboard
      .writeText(text)
      .then(() => {
        setCopied(true);
        setTimeout(() => {
          setCopied(false);
        }, 2000);
      })
      .catch((error) => {
        console.error('error copying to clipboard:', error);
      });
  };

  return (
    <NoSsr>
      <ErrorBoundary>
        <RootContainer>
          <Typography variant="h5">Provider Information</Typography>
          <ErrorBoundary>
            <Grid2 container spacing={2} size="grow">
              {providerInfo &&
                Object.entries(providerInfo).map(
                  ([providerName, provider], index) =>
                    (index < 2 || index === 3) && (
                      <Grid2 key={index} size={{ xs: 12, md: 4 }}>
                        <ProviderCard>
                          <CardHeader
                            title={
                              <Typography
                                variant="h6"
                                style={{
                                  textDecoration: 'underline',
                                  textDecorationColor: 'rgba(116,147,161,0.5)',
                                  textUnderlineOffset: 10,
                                }}
                              >
                                {convertToTitleCase(providerName)}
                              </Typography>
                            }
                          />
                          <CardContent>
                            {' '}
                            <BoxWrapper>
                              <Typography
                                variant="body1"
                                component={HideScrollbar}
                                style={{ marginRight: '20px' }}
                              >
                                {provider}
                              </Typography>
                            </BoxWrapper>
                          </CardContent>
                        </ProviderCard>
                      </Grid2>
                    ),
                )}
            </Grid2>
            {providerInfo &&
              Object.entries(providerInfo).map(
                ([providerName, provider], index) =>
                  (index === 2 || index === 5) && (
                    <ProviderCard key={index} sx={{ margin: '20px' }}>
                      <CardHeader
                        title={
                          <Typography
                            variant="h6"
                            style={{
                              textDecoration: 'underline',
                              textDecorationColor: 'rgba(116,147,161,0.5)',
                              textUnderlineOffset: 10,
                            }}
                          >
                            {convertToTitleCase(providerName)}
                          </Typography>
                        }
                      />
                      <CardContent>
                        {' '}
                        <BoxWrapper>
                          <Typography
                            variant="body1"
                            component={HideScrollbar}
                            style={{ marginRight: '20px' }}
                          >
                            {provider}
                          </Typography>

                          <CustomTooltip title={copied ? 'Copied!' : 'Copy'} placement="top">
                            <IconButton
                              onClick={() => copyToClipboard(provider)}
                              style={{ padding: '0.25rem', float: 'right' }}
                            >
                              <CopyIcon />
                            </IconButton>
                          </CustomTooltip>
                        </BoxWrapper>
                      </CardContent>
                    </ProviderCard>
                  ),
              )}
          </ErrorBoundary>
          <ProviderCard>
            <CardHeader
              title={
                <Typography
                  variant="h6"
                  style={{
                    textDecoration: 'underline',
                    textDecorationColor: 'rgba(116,147,161,0.5)',
                    textUnderlineOffset: 10,
                  }}
                >
                  Description
                </Typography>
              }
            />
            <CardContent>
              <Typography component="div">
                <ul>
                  {providerInfo.provider_description &&
                    providerInfo.provider_description.map((desc, index) => (
                      <li key={index}>
                        <Typography>{desc}</Typography>
                      </li>
                    ))}
                </ul>
              </Typography>
            </CardContent>
          </ProviderCard>

          <Divider />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Capabilities
          </Typography>

          <Grid2 container size="grow" style={{ display: 'flex', justifyContent: 'space-between' }}>
            <GridCapabilityHeader
              size={{ xs: 6 }}
              style={{ borderRadius: '10px 0 0 0', padding: '10px 20px' }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Feature
              </Typography>
            </GridCapabilityHeader>
            <GridCapabilityHeader
              size={{ xs: 6 }}
              style={{ borderRadius: '0 10px 0 0', padding: '10px 20px' }}
            >
              <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                Endpoint
              </Typography>
            </GridCapabilityHeader>
            {providerInfo.capabilities &&
              providerInfo.capabilities.map((capability, index) => (
                <React.Fragment key={`${index}-${capability.feature}`}>
                  <GridCapabilityHeader
                    size={{ xs: 6 }}
                    style={{
                      padding: '20px 20px',
                      backgroundColor:
                        userData?.remoteProviderPreferences?.theme === 'dark'
                          ? index % 2 === 0
                            ? '#3D4F57'
                            : '#293B43'
                          : index % 2 === 0
                            ? '#E7EFF3'
                            : '#C9DBE3',
                    }}
                  >
                    <Typography variant="body1">{capability.feature}</Typography>
                  </GridCapabilityHeader>
                  <GridCapabilityHeader
                    size={{ xs: 6 }}
                    style={{
                      padding: '20px 20px',
                      backgroundColor:
                        userData?.remoteProviderPreferences?.theme === 'dark'
                          ? index % 2 === 0
                            ? '#3D4F57'
                            : '#293B43'
                          : index % 2 === 0
                            ? '#E7EFF3'
                            : '#C9DBE3',
                    }}
                  >
                    <Typography variant="body1">{capability.endpoint}</Typography>
                  </GridCapabilityHeader>
                </React.Fragment>
              ))}
          </Grid2>
          <Divider />
          <Typography variant="h5" style={{ margin: '20px 0' }}>
            Extensions
          </Typography>
          {providerInfo.extensions &&
            Object.entries(providerInfo.extensions).map(([extensionName, extension], index) => (
              <div key={index} style={{ margin: '20px 0px' }}>
                <Typography variant="h6"> {convertToTitleCase(extensionName)}</Typography>
                <Grid2 container spacing={2} size="grow" style={{ margin: '10px 0 20px 0' }}>
                  <GridExtensionHeader
                    size={{ xs: 6 }}
                    style={{
                      borderRadius: '10px 0 0 0',
                      padding: '10px 20px',
                    }}
                  >
                    <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                      Component
                    </Typography>
                  </GridExtensionHeader>
                  <GridExtensionHeader
                    size={{ xs: 6 }}
                    style={{
                      borderRadius: '0 10px 0 0',
                      padding: '10px 20px',
                    }}
                  >
                    <Typography variant="body1" style={{ fontWeight: 'bold' }}>
                      Type
                    </Typography>
                  </GridExtensionHeader>

                  {extension[0] && (
                    <>
                      <GridExtensionItem
                        size={{ xs: 6 }}
                        style={{
                          borderRadius: '0 0 0 10px',
                          padding: '20px 20px',
                        }}
                      >
                        <Typography variant="body1">{extension[0].component}</Typography>
                      </GridExtensionItem>
                      <GridExtensionItem
                        size={{ xs: 6 }}
                        style={{
                          borderRadius: '0 0 10px 0',
                          padding: '20px 20px',
                        }}
                      >
                        <Typography variant="body1">
                          {convertToTitleCase(extension[0].type)}
                        </Typography>
                      </GridExtensionItem>
                    </>
                  )}
                </Grid2>
              </div>
            ))}
        </RootContainer>
      </ErrorBoundary>
    </NoSsr>
  );
};

interface RemoteProviderTabProps {
  providerInfo: ProviderInfoData;
  userData?: {
    remoteProviderPreferences?: {
      theme?: string;
    };
  };
}

const RemoteProviderTab: React.FC<RemoteProviderTabProps> = ({ providerInfo, userData }) => {
  const [value, setValue] = useState(0);
  const { width } = useWindowDimensions();

  const handleValChange = (_event: React.SyntheticEvent, newVal: number) => {
    setValue(newVal);
  };

  return (
    <>
      <SecondaryTabs
        value={value}
        onChange={handleValChange}
        variant={width < 600 ? 'scrollable' : 'fullWidth'}
        scrollButtons="on"
        allowScrollButtonsMobile={true}
        indicatorColor="primary"
        textColor="primary"
        centered
      >
        <CustomTooltip title="Details" placement="top">
          <SecondaryTab label={<TabLabel>Details</TabLabel>} />
        </CustomTooltip>
        <CustomTooltip title="Plugins" placement="top">
          <SecondaryTab label={<TabLabel>Plugins</TabLabel>} />
        </CustomTooltip>
      </SecondaryTabs>
      <StatsWrapper>
        {value === 0 && <ProviderInfoDetails providerInfo={providerInfo} userData={userData} />}
        {value === 1 && (
          <ExtensionSandbox type="user_prefs" Extension={(url) => RemoteComponent({ url })} />
        )}
      </StatsWrapper>
    </>
  );
};

export default RemoteProviderTab;
