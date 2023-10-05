import * as React from 'react';
import { Divider, Grid, Tooltip, Typography } from '@material-ui/core';
import { Launch as LaunchIcon } from '@material-ui/icons';
import _ from 'lodash';

export const Link = ({ href, title }) => {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      style={{
        color: 'inherit',
        textDecorationLine: 'underline',
        cursor: 'pointer',
        marginBottom: '0.5rem',
      }}
    >
      {title}
      <sup>
        <LaunchIcon style={{ width: '1rem', height: '1rem' }} />
      </sup>
    </a>
  );
};
function truncateString(str, maxLength) {
  if (str.length > maxLength) {
    return str.slice(0, maxLength) + '...';
  }
  return str;
}
function captureLinkData(url) {
  for (const [title, siteUrl] of Object.entries(SITE_URLS)) {
    if (url.startsWith(siteUrl)) {
      return {
        href: url,
        title,
      };
    }
  }

  return {
    href: url,
    title: truncateString(url, 50),
  };
}

const SITE_URLS = {
  Doc: 'https://docs.meshery.io',
};

export const TextWithLinks = ({ text, ...typographyProps }) => {
  // Regular expression to find HTTP links in the text
  const linkRegex = /(https?:\/\/[^\s]+)/g;

  // Split the text into parts, alternating between text and link components
  const parts = text.split(linkRegex);

  // Map the parts to React elements
  const elements = parts.map((part) => {
    if (part.match(linkRegex)) {
      // If the part is a link, wrap it in a Link component
      const link_data = captureLinkData(part);
      return <Link href={link_data.href} title={link_data.title} />;
    } else {
      return <span>{part}</span>;
    }
  });

  return <Typography {...typographyProps}>{elements}</Typography>;
};

function isEmptyAtAllDepths(input) {
  if (_.isArray(input)) {
    // If the input is an array, check if all items are empty at all depths
    return input.every(isEmptyAtAllDepths);
  } else if (_.isObject(input)) {
    // If the input is an object, check if all properties are empty at all depths
    return _.every(input, isEmptyAtAllDepths);
  } else {
    // If the input is not an array or object, check if it's empty
    return _.isEmpty(input);
  }
}

export const KeyValue = ({ Key, Value }) => {
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'end',
        gap: '0.25rem',
        flexWrap: 'wrap',
        marginBlock: '0.25rem',
      }}
    >
      <SectionBody body={Key} style={{ fontWeight: 'bold' }} /> <SectionBody body={Value} />
    </div>
  );
};

const NestedData = ({ heading, data, classes }) => {
  if (!data || data?.length == 0) return null;
  return (
    <>
      <Typography variant="h6" className={classes.descriptionHeading}>
        {heading}
      </Typography>
      {typeof data === 'string' ? (
        <SectionBody body={data}></SectionBody>
      ) : (
        <BulletList items={data} />
      )}
    </>
  );
};
export const MetaDataSectionHeading = ({ children, level, ...props }) => {
  const fontSize = Math.max(0.9, 1.5 - 0.2 * level) + 'rem';
  return (
    <Typography
      variant="h5"
      style={{
        fontWeight: 'bolder !important',
        textTransform: 'uppercase',
        marginBlock: '0.2rem',
        wordBreak: 'break-all',
        fontSize,
      }}
      {...props}
    >
      {children}
    </Typography>
  );
};

export const SectionBody = ({ body, style = {} }) => {
  return (
    <TextWithLinks
      variant="body1"
      style={{
        textTransform: 'capitalize',
        wordWrap: 'break-word',
        ...style,
      }}
      text={body}
    ></TextWithLinks>
  );
};

const DryRunResponse = ({ response }) => {
  const cleanedResponse = [];
  Object.entries(response).forEach(([componentKind, components]) => {
    Object.entries(components).forEach(([id, data]) => {
      cleanedResponse.push({
        name: componentKind,
        error: data.error,
        id: id,
      });
    });
  });

  const componentHasErrors = (component) => {
    return component?.error?.Causes.length > 0 || component.error.Status !== '';
  };

  const formattedResponse = cleanedResponse.filter(componentHasErrors).map((component) => (
    <div style={{ marginBottom: '1rem' }} key={component.id}>
      <Tooltip title="dkje" placement="top">
        <MetaDataSectionHeading level={2}>{component.name}</MetaDataSectionHeading>
      </Tooltip>
      <DynamicMetadataFormatter metadata={component.error} />
      <Divider light />
    </div>
  ));

  if (formattedResponse.length == 0) {
    return (
      <Typography variant="h6" style={{ textAlign: 'center', marginBlock: '1rem' }}>
        No Errors Found
      </Typography>
    );
  }
  return formattedResponse;
};

const DynamicMetadataFormatter = ({ metadata, level = 1 }) => {
  if (_.isString(metadata)) {
    return <SectionBody body={metadata}></SectionBody>;
  }
  if (_.isArray(metadata)) {
    return <BulletList items={metadata} />;
  }

  if (_.isObject(metadata)) {
    return Object.entries(metadata).map(([title, data]) => {
      if (!title.trim() || !data || _.isEmpty(data)) {
        return null;
      }
      if (typeof data == 'string') {
        if (title == 'doc') {
          return <Link title="Doc" href={data} />;
        }
        return <KeyValue key={title} Key={title} Value={data} />;
      }
      if (title === 'dryRunResponse') {
        return <DryRunResponse key={title} response={data} />;
      }
      return (
        <div
          key={title}
          style={{
            paddingLeft: level > 1 ? '0.25rem' : '0rem',
            marginBlock: level > 1 ? '0.25rem' : '0rem',
          }}
        >
          <MetaDataSectionHeading level={level}>{title}</MetaDataSectionHeading>
          <DynamicMetadataFormatter level={level + 1} metadata={data} />
        </div>
      );
    });
  }

  return <EmptyState />;
};

const BulletList = ({ items }) => {
  return (
    <ol style={{ paddingInline: '0.75rem', paddingBlock: '0.3rem', margin: '0rem' }}>
      {items.map((item) => (
        <li key={item}>
          <DynamicMetadataFormatter metadata={item} />
        </li>
      ))}
    </ol>
  );
};

export const ErrorMetadataFormatter = ({ metadata, event, classes }) => {
  const longDescription = metadata?.LongDescription || [];
  const probableCause = metadata?.ProbableCause || [];
  const suggestedRemediation = metadata?.SuggestedRemediation || [];
  const errorCode = metadata?.error_code || '';
  const code = metadata?.Code || '';
  const formattedErrorCode = errorCode ? `${errorCode}-${code}` : code;
  const errorLink = `https://docs.meshery.io/reference/error-codes#${formattedErrorCode}`;
  return (
    <Grid container>
      <div>
        <a href={errorLink} target="_blank" rel="noopener noreferrer" style={{ color: 'inherit' }}>
          <Typography
            variant="h5"
            className={classes.descriptionHeading}
            style={{ textDecorationLine: 'underline', cursor: 'pointer', marginBottom: '0.5rem' }}
          >
            {formattedErrorCode}
            <sup>
              <LaunchIcon style={{ width: '1rem', height: '1rem' }} />
            </sup>
          </Typography>
        </a>
        <NestedData classes={classes} data={event.description} />
        <div style={{ marginTop: '1rem' }}>
          <NestedData classes={classes} heading="Details" data={longDescription} />
        </div>
      </div>
      <Grid container spacing={1} style={{ marginTop: '0.5rem' }}>
        <Grid item sm={suggestedRemediation?.length > 0 ? 6 : 12}>
          <NestedData classes={classes} heading="Probable Cause" data={probableCause} />
        </Grid>
        <Grid item sm={probableCause?.length > 0 ? 6 : 12}>
          <NestedData
            classes={classes}
            heading="Suggested Remediation"
            data={suggestedRemediation}
          />
        </Grid>
      </Grid>
    </Grid>
  );
};

const EmptyState = () => {
  return (
    <div>
      <Typography
        variant="h6"
        style={{
          textAlign: 'center',
          marginBlock: '1rem',
        }}
      >
        No Further Details Available{' '}
      </Typography>
    </div>
  );
};

const METADATA_FORMATTER = {
  error: ErrorMetadataFormatter,
  default: DynamicMetadataFormatter, // for all other metadata types
  empty: EmptyState,
};

// Maps the metadata to the appropriate formatter component
export const FormattedMetadata = ({ event, classes }) => {
  if (!event || !event.metadata || isEmptyAtAllDepths(event.metadata)) {
    return <EmptyState event={event} />;
  }
  const metdataKeys = Object.keys(event.metadata);
  return metdataKeys.map((key) => {
    const Formatter = METADATA_FORMATTER[key] || METADATA_FORMATTER.default;
    return <Formatter key={key} metadata={event.metadata[key]} event={event} classes={classes} />;
  });
};
