import * as React from 'react';
import { CustomTooltip, Typography, Grid, Box, IconButton, useTheme } from '@layer5/sistent';
import { Launch as LaunchIcon } from '@mui/icons-material';
import _ from 'lodash';
import { useContext } from 'react';
import { isEmptyAtAllDepths } from '../../utils/objects';
import CopyIcon from '../../assets/icons/CopyIcon';
import { UsesSistent } from './../SistentWrapper';
const FormatterContext = React.createContext({
  propertyFormatters: {},
});

/**
 * Context to store the level / depth of content in the formatter
 * @type {React.Context<number>}
 */
const LevelContext = React.createContext(0);

/**
 * Level context provider to autoincrement the level of content in the formatter
 * @param {React.PropsWithChildren<{}>} param0
 * @returns {React.ReactElement}
 */
const Level = ({ children }) => {
  const level = useContext(LevelContext);
  return <LevelContext.Provider value={level + 1}> {children} </LevelContext.Provider>;
};

/**
 * Pure function to format data
 * @returns {string}
 */
export const formatDate = (date) => {
  const options = { year: 'numeric', month: 'short', day: 'numeric' };
  const formattedDate = new Date(date).toLocaleDateString('en-US', options);
  return formattedDate;
};

export const formatTime = (date) => {
  const options = { hour: 'numeric', minute: 'numeric', second: 'numeric' };
  const formattedTime = new Date(date).toLocaleTimeString('en-US', options);
  return formattedTime;
};

export const formatDateTime = (date) => {
  const formattedDate = formatDate(date);
  const formattedTime = formatTime(date);
  return `${formattedDate} ${formattedTime || ''}`;
};

/**
 * Component to format the date
 * @param {object} param0
 * @returns {React.ReactElement}
 */
export const FormattedDate = ({ date }) => {
  return (
    <CustomTooltip title={formatDateTime(date)} placement="top">
      <div>
        <SectionBody
          body={formatDate(date)}
          style={{
            textTransform: 'capitalize',
          }}
        ></SectionBody>
      </div>
    </CustomTooltip>
  );
};

export const FormatId = ({ id }) => {
  const [copied, setCopied] = React.useState(false);
  const theme = useTheme();
  // truncates the id to 15 characters and adds an ellipsis and adds a clipboard copy button
  const copyToClipboard = (e) => {
    e.stopPropagation();
    navigator.clipboard.writeText(id);
    setCopied(true);

    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  const truncatedId = _.truncate(id, { length: 15 });
  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: '0.25rem' }}>
      <CustomTooltip title={id} placement="top">
        <Typography
          variant="body2"
          style={{
            cursor: 'pointer',
            color: theme.palette.text.tertiary,
          }}
        >
          {truncatedId}
        </Typography>
      </CustomTooltip>
      <CustomTooltip title={copied ? 'Copied!' : 'Copy'} placement="top">
        <IconButton onClick={copyToClipboard} style={{ padding: '0.25rem' }}>
          <CopyIcon width="1rem" height="1rem" />
        </IconButton>
      </CustomTooltip>
    </Box>
  );
};

export const createColumnUiSchema = ({ metadata, numCols }) => {
  return Object.keys(metadata).reduce((schema, key) => {
    schema[key] = Object.keys(numCols).reduce(
      (colSpan, key) => ({
        ...colSpan,
        [key]: Math.floor(12 / numCols[key]),
      }),
      {},
    );
    return schema;
  }, {});
};
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
        marginBottom: '0.75rem',
      }}
    >
      {title}
      <sup>
        <LaunchIcon style={{ width: '1rem', height: '1rem' }} />
      </sup>
    </a>
  );
};

function getFormattedLink(url) {
  for (const formatter of Object.values(LinkFormatters)) {
    if (url.startsWith(formatter.base_url)) {
      return formatter.formatter(url);
    }
  }

  return LinkFormatters.DEFAULT.formatter(url);
}

export const LinkFormatters = {
  DOC: {
    base_url: 'https://docs.meshery.io',
    formatter: (link) => <Link title="Doc" href={link} />,
  },
  DEFAULT: {
    base_url: '',
    formatter: (link) => <Link title={_.truncate(link, 30)} href={link} />,
  },
};

export const TextWithLinks = ({ text, ...typographyProps }) => {
  // Regular expression to find HTTP links in the text
  const linkRegex = /(https?:\/\/[^\s]+)/g;

  // Split the text into parts, alternating between text and link components
  const parts = text?.split?.(linkRegex) || [];

  // Map the parts to React elements
  const elements = parts.map((part, idx) => {
    if (part.match(linkRegex)) {
      // If the part is a link, wrap it in a Link component
      return getFormattedLink(part);
    } else {
      return <span key={idx}>{part}</span>;
    }
  });

  return <Typography {...typographyProps}> {elements}</Typography>;
};

export const KeyValue = ({ Key, Value }) => {
  const theme = useTheme();
  return (
    <div
      style={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'flex-start',
        gap: '0.25rem',
        flexWrap: 'wrap',
        fontFamily: 'Qanelas Soft, sans-serif',
      }}
    >
      <SectionBody
        body={Key.replaceAll('_', ' ')}
        style={{
          textTransform: 'capitalize',
          overflow: 'hidden',
          whiteSpace: 'nowrap',
          color: theme.palette.text.default,
        }}
      />

      {React.isValidElement(Value) ? (
        Value
      ) : (
        <SectionBody
          body={Value}
          style={{
            color: theme.palette.text.tertiary,
            textOverflow: 'ellipsis',
            wordBreak: 'break-all',
            overflow: 'hidden',
            whiteSpace: 'nowrap',
          }}
        />
      )}
    </div>
  );
};

export const SectionHeading = ({ children, isLevel, ...props }) => {
  const level = useContext(LevelContext);
  const fontSize = isLevel ? Math.max(0.9, 1.3 - 0.1 * level) + 'rem' : '1rem';
  const margin = isLevel ? Math.max(0.25, 0.55 - 0.15 * level) + 'rem' : 'inherit';

  return (
    <div style={{ marginBlock: margin }}>
      <Typography
        variant="body1"
        style={{
          fontWeight: 'bold !important',
          textTransform: 'capitalize',
          wordBreak: 'break-all',
          fontSize,
        }}
        {...props}
      >
        {children}
      </Typography>
    </div>
  );
};
export const SectionBody = ({ body, style = {} }) => {
  const theme = useTheme();
  return (
    <TextWithLinks
      variant="body1"
      style={{
        wordWrap: 'break-word',
        color: theme.palette.text.tertiary,
        ...style,
      }}
      text={body}
    ></TextWithLinks>
  );
};
export const ArrayFormatter = ({ items }) => {
  const theme = useTheme();
  return (
    <ol
      style={{
        paddingInline: '1rem',
        paddingBlock: '0.25rem',
        margin: '0rem',
      }}
    >
      {items.map((item) => (
        <li key={item} style={{ color: theme.palette.text.tertiary }}>
          <Level>
            <DynamicFormatter data={item} />
          </Level>
        </li>
      ))}
    </ol>
  );
};

export function reorderObjectProperties(obj, order) {
  if (!_.isObject(obj) || obj == null) {
    return obj;
  }

  const orderedProperties = _.pick(obj, order);
  const remainingProperties = _.omit(obj, order);
  return { ...orderedProperties, ...remainingProperties };
}

const DynamicFormatter = ({ data, uiSchema, isLevel = true }) => {
  const { propertyFormatters } = useContext(FormatterContext);
  const level = useContext(LevelContext);

  if (_.isString(data)) {
    return <SectionBody body={data}></SectionBody>;
  }

  if (_.isArray(data)) {
    return <ArrayFormatter items={data} />;
  }

  if (_.isObject(data)) {
    return Object.entries(data).map(([title, data]) => {
      if (!title.trim() || !data || _.isEmpty(data)) {
        return null;
      }

      if (propertyFormatters?.[title]) {
        return (
          <Grid key={title} sm={12} {...(uiSchema?.[title] || {})}>
            {propertyFormatters[title](data, data)}
          </Grid>
        );
      }

      if (typeof data == 'string') {
        return (
          <Grid
            item
            key={title}
            sm={12}
            {...(uiSchema?.[title] || {})}
            spacing={3}
            style={{
              marginBlock: '0.4rem',
              maxWidth: title !== 'age' && 'fit-content',
              marginRight: '1rem',
            }}
          >
            <KeyValue key={title} Key={title} Value={data} />
          </Grid>
        );
      }

      return (
        <Grid
          item
          key={title}
          sm={12}
          {...(uiSchema?.[title] || {})}
          style={{
            marginBlock: '0.25rem',
          }}
        >
          <SectionHeading level={level} isLevel={isLevel}>
            {title}
          </SectionHeading>
          <Level>
            <DynamicFormatter level={level + 1} data={data} />
          </Level>
        </Grid>
      );
    });
  }

  return null;
};

export const FormatStructuredData = ({ propertyFormatters = {}, data, uiSchema, isLevel }) => {
  if (!data || isEmptyAtAllDepths(data)) {
    return null;
  }

  return (
    <UsesSistent>
      <FormatterContext.Provider
        value={{
          propertyFormatters: propertyFormatters,
        }}
      >
        <Grid
          container
          style={{
            wordBreak: 'break-word',
            overflowWrap: 'break-word',
            gap: '0.3rem 1rem',
          }}
        >
          <DynamicFormatter data={data} uiSchema={uiSchema} isLevel={isLevel} />
        </Grid>
      </FormatterContext.Provider>
    </UsesSistent>
  );
};
