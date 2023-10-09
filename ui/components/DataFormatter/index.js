import * as React from 'react';
import { Typography } from '@material-ui/core';
import { Launch as LaunchIcon } from '@material-ui/icons';
import _ from 'lodash';
import { useContext } from 'react';
import { isEmptyAtAllDepths } from '../../utils/objects';

const FormatterContext = React.createContext({
  propertyFormatters: {},
});
const LevelContext = React.createContext(0);

const Level = ({ children }) => {
  const level = useContext(LevelContext);
  return <LevelContext.Provider value={level + 1}> {children} </LevelContext.Provider>;
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
  const parts = text.split(linkRegex);

  // Map the parts to React elements
  const elements = parts.map((part) => {
    if (part.match(linkRegex)) {
      // If the part is a link, wrap it in a Link component
      return getFormattedLink(part);
    } else {
      return <span>{part}</span>;
    }
  });

  return <Typography {...typographyProps}>{elements}</Typography>;
};

export const KeyValue = ({ Key, Value }) => {
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'baseline',
        gap: '0.25rem',
        flexWrap: 'wrap',
        marginBlock: '0.5rem',
      }}
    >
      <SectionBody body={Key + ':'} style={{ fontWeight: 'bold' }} /> <SectionBody body={Value} />
    </div>
  );
};

export const SectionHeading = ({ children, ...props }) => {
  const level = useContext(LevelContext);
  const fontSize = Math.max(0.9, 1.3 - 0.1 * level) + 'rem';
  const margin = Math.max(0.25, 0.55 - 0.15 * level) + 'rem';

  return (
    <Typography
      variant="h5"
      style={{
        fontWeight: 'bold !important',
        textTransform: 'capitalize',
        marginBottom: margin,
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

const ArrayFormatter = ({ items }) => {
  return (
    <ol style={{ paddingInline: '0.75rem', paddingBlock: '0.25rem', margin: '0rem' }}>
      {items.map((item) => (
        <li key={item}>
          <Level>
            <DynamicFormatter data={item} />
          </Level>
        </li>
      ))}
    </ol>
  );
};

const DynamicFormatter = ({ data }) => {
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
        return propertyFormatters[title](data, data);
      }
      if (typeof data == 'string') {
        return <KeyValue key={title} Key={title} Value={data} />;
      }

      return (
        <div
          key={title}
          style={{
            marginBlock: '0.25rem',
          }}
        >
          <SectionHeading level={level}>{title}</SectionHeading>
          <Level>
            <DynamicFormatter level={level + 1} data={data} />
          </Level>
        </div>
      );
    });
  }

  return null;
};

export const FormatStructuredData = ({ propertyFormatters = {}, data }) => {
  if (!data || isEmptyAtAllDepths(data)) {
    console.log('data is empty', data);
    return null;
  }
  return (
    <>
      <FormatterContext.Provider
        value={{
          propertyFormatters: propertyFormatters,
        }}
      >
        <DynamicFormatter data={data} />
      </FormatterContext.Provider>
    </>
  );
};
