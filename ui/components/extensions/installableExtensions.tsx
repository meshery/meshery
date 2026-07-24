// Meshery Extension Point
// Add your installable plugins  as a component with the Extension Interface into this extension point.
// Learn more: https://docs.meshery.io/extensibility/ui

import React, { useEffect, useMemo, useState } from 'react';
import { Box, Button, Grid2, Typography, styled } from '@sistent/sistent';
import {
  useGetProviderCapabilitiesQuery,
  useGetSystemVersionQuery,
  useInstallProviderExtensionMutation,
  useRemoveProviderExtensionMutation,
} from '@/rtk-query/user';
import { CardContainer, FrontSideDescription } from 'css/icons.styles';
import { EVENT_TYPES } from '../../lib/event-types';
import { useNotification } from '@/utils/hooks';
import { formatApiError } from '@/utils/helpers/meshkitError';
import { isLocalProvider } from '@/utils/provider';

type ChildrenProps = {
  children: React.ReactNode;
};

type ExtensionMetadata = {
  title?: string;
  href?: { uri?: string; external?: boolean } | string;
  icon?: string;
  link?: boolean;
  [k: string]: any;
};

type Extension = {
  type?: string;
  version?: string;
  packagePath?: string;
  title: string;
  description?: string;
  icon?: string;
  metadata?: ExtensionMetadata;
};

// Meshery Extension Point
// Add your installable plugins  as a component with the Extension Interface into this extension point.
// Learn more: https://docs.meshery.io/extensibility/ui
type InstallableExtensionProps = {
  extension: Extension;
};

const GRID_SIZE = { xs: 12, sm: 12, md: 12, lg: 6, xl: 4 };

const UnifiedCardContainer = ({ children, sx = {} }: ChildrenProps & { sx?: object }) => (
  <CardContainer
    sx={{
      height: '100%',
      minHeight: { xs: '280px', sm: '260px', lg: '280px' },
      display: 'flex',
      flexDirection: 'column',
      position: 'relative',
      ...sx,
    }}
  >
    {children}
  </CardContainer>
);

const UnifiedDescription = ({
  children,
  hasIcon = false,
  ...props
}: ChildrenProps & { hasIcon?: boolean; 'data-testid'?: string }) => (
  <FrontSideDescription
    {...props}
    sx={{
      flex: 1,
      display: 'flex',
      flexDirection: { xs: hasIcon ? 'column' : 'row', sm: 'row' },
      alignItems: { xs: hasIcon ? 'center' : 'flex-start', sm: 'flex-start' },
      textAlign: { xs: hasIcon ? 'center' : 'left', sm: 'left' },
      gap: { xs: hasIcon ? '12px' : '8px', sm: '8px' },
      marginBottom: { xs: '50px', sm: '45px', lg: '40px' },
      overflow: 'hidden',
      wordWrap: 'break-word',
      hyphens: 'auto',
      '& > *:not(img)': {
        maxWidth: '100%',
        overflow: 'hidden',
        wordWrap: 'break-word',
        hyphens: 'auto',
      },
    }}
  >
    {children}
  </FrontSideDescription>
);

const UnifiedButtonContainer = ({ children }: ChildrenProps) => (
  <Box sx={{ position: 'absolute', bottom: 12, left: 12, right: 12, textAlign: 'right' }}>
    {children}
  </Box>
);

const StyledResponsiveImage = styled('img')({
  height: 'auto',
  width: 'auto',
  maxWidth: '140px',
  maxHeight: '85px',
  flexShrink: 0,
});

const ResponsiveImage = ({ src, alt, testId }: { src: string; alt?: string; testId?: string }) => (
  <StyledResponsiveImage data-testid={testId} src={src} alt={alt || ''} />
);

const resolveExtensionHref = (href?: ExtensionMetadata['href']) => {
  if (!href) {
    return '';
  }

  // Guarantee exactly one separating slash so a relative href such as 'meshmap'
  // resolves to '/extension/meshmap' rather than the malformed '/extensionmeshmap'.
  const ensureLeadingSlash = (path: string) => (path.startsWith('/') ? path : '/' + path);

  if (typeof href === 'string') {
    if (
      href.startsWith('http://') ||
      href.startsWith('https://') ||
      href.startsWith('/extension')
    ) {
      return href;
    }

    return '/extension' + ensureLeadingSlash(href);
  }

  if (!href.uri) {
    return '';
  }

  if (href.external || href.uri.startsWith('http://') || href.uri.startsWith('https://')) {
    return href.uri;
  }

  if (href.uri.startsWith('/extension')) {
    return href.uri;
  }

  return '/extension' + ensureLeadingSlash(href.uri);
};

const InstallableExtension: React.FC<InstallableExtensionProps> = ({ extension }) => {
  const { notify } = useNotification();
  const { data: providerCaps } = useGetProviderCapabilitiesQuery();
  const [installProviderExtension, { isLoading: isInstalling }] =
    useInstallProviderExtensionMutation();
  const [removeProviderExtension, { isLoading: isRemovingFromProvider }] =
    useRemoveProviderExtensionMutation();

  const installedFromProvider = useMemo(() => {
    try {
      const extType = (extension.type || extension.metadata?.type || 'navigator').toLowerCase();
      if (!providerCaps || !providerCaps.extensions) return false;
      const list = (providerCaps.extensions as any)[extType];
      if (!Array.isArray(list)) return false;
      return list.some((e: any) => {
        if (!e) return false;
        const title = (e.title || e.Title || e.name || '').toString();
        return (
          title &&
          title.toLowerCase() === extension.title.toLowerCase() &&
          e.component === extension.metadata?.component
        );
      });
    } catch {
      return false;
    }
  }, [providerCaps, extension]);

  const [installed, setInstalled] = useState<boolean>(installedFromProvider);
  useEffect(() => setInstalled(installedFromProvider), [installedFromProvider]);

  const [isRemoving, setIsRemoving] = useState<boolean>(false);
  const isLocal = isLocalProvider(providerCaps);
  const isMutating = isInstalling || isRemoving || isRemovingFromProvider;
  const installReady = Boolean(extension.packagePath);

  const iconSrc =
    extension.icon || extension.metadata?.icon || '/static/img/extensions/default-extension.svg';
  const description =
    extension.description || extension.metadata?.description || 'No description available.';

  const handleInstall = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    try {
      await installProviderExtension({
        extType: extension.type || extension.metadata?.type || 'navigator',
        packageUrl: extension.packagePath || '',
        extensionMetadata: {
          title: extension.title,
          description: extension.description,
          icon: extension.icon,
          ...(extension.metadata || {}),
        },
      }).unwrap();

      setInstalled(true);
      notify({ message: `${extension.title} installed`, event_type: EVENT_TYPES.SUCCESS });
    } catch (error) {
      const { message } = formatApiError(error, `Failed to install ${extension.title}`);
      notify({ message, event_type: EVENT_TYPES.ERROR });
    }
  };

  const handleUninstall = async (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    setIsRemoving(true);
    try {
      await removeProviderExtension({
        extType: extension.type || extension.metadata?.type || 'navigator',
        title: extension.title,
      }).unwrap();
      setInstalled(false);
      notify({ message: `${extension.title} removed`, event_type: EVENT_TYPES.SUCCESS });
    } catch (error) {
      const { message } = formatApiError(error, `Failed to remove ${extension.title}`);
      notify({ message, event_type: EVENT_TYPES.ERROR });
    } finally {
      setIsRemoving(false);
    }
  };

  const openLink = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    const url = resolveExtensionHref(extension.metadata?.href);
    if (url) window.open(url, '_blank', 'noopener,noreferrer');
  };

  return (
    <Grid2 size={GRID_SIZE}>
      <UnifiedCardContainer>
        <Typography variant="h5" component="div" data-testid="installable-extension-title">
          {extension.title}
        </Typography>

        <UnifiedDescription hasIcon data-testid="installable-extension-desc">
          <ResponsiveImage
            src={iconSrc}
            alt={`${extension.title} icon`}
            testId="installable-extension-icon"
          />

          <div>
            <Box sx={{ marginBottom: 2 }}>{description}</Box>
            {extension.metadata?.isBeta && (
              <Typography variant="caption" sx={{ fontStyle: 'italic' }}>
                Beta
              </Typography>
            )}
          </div>
        </UnifiedDescription>

        <UnifiedButtonContainer>
          {!installed && isLocal ? (
            <Button
              variant={isLocal ? 'contained' : 'outlined'}
              color="primary"
              onClick={handleInstall}
              data-testid="install-btn"
              disabled={!isLocal || !installReady || isMutating}
            >
              {isInstalling ? 'Installing...' : installReady ? 'Install' : 'Preparing...'}
            </Button>
          ) : (
            <>
              {isLocal ? (
                <Button
                  variant="outlined"
                  color="secondary"
                  onClick={handleUninstall}
                  data-testid="uninstall-btn"
                  disabled={isMutating}
                  sx={{ marginRight: 2 }}
                >
                  {isRemoving || isRemovingFromProvider ? 'Removing...' : 'Remove'}
                </Button>
              ) : null}

              {extension.metadata?.link || extension.metadata?.href || extension.packagePath ? (
                <Button
                  variant="contained"
                  color="primary"
                  onClick={openLink}
                  data-testid="open-extension-btn"
                >
                  Open
                </Button>
              ) : null}
            </>
          )}
        </UnifiedButtonContainer>
      </UnifiedCardContainer>
    </Grid2>
  );
};

export const VisualDesignerExtension: React.FC = () => {
  const { data: mesheryVersionData } = useGetSystemVersionQuery();
  const version = mesheryVersionData?.build ? `${mesheryVersionData.build}-1` : null;

  const data: Extension = {
    type: 'navigator',
    packagePath: version
      ? `https://github.com/meshery-extensions/meshery-extensions-packages/releases/download/${version}/provider-meshery.tar.gz`
      : '',
    title: 'Kanvas',
    icon: '/static/img/extensions/visual-designer-extension-logo.svg',
    description:
      'Collaboratively design and manage your Kubernetes clusters and Cloud services. Kanvas is a visual interface for managing your infrastructure. It allows you to create, edit, and share your infrastructure as code.',
    metadata: {
      title: 'Kanvas',
      onClickCallback: 1,
      href: {
        uri: '/meshmap',
        external: false,
      },
      component: version
        ? `/provider/navigator/meshmap/index.js?packageVersion=${encodeURIComponent(version)}`
        : '/provider/navigator/meshmap/index.js',
      icon: '/provider/navigator/img/kanvas-icon.svg',
      link: true,
      show: true,
      type: 'full_page',
      allowedTo: {
        designer: {},
      },
      isBeta: true,
      version,
    },
  };

  return <InstallableExtension extension={data} />;
};
