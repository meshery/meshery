import React from 'react';
import { Typography } from '@sistent/sistent';
import { Launch as LaunchIcon } from '@mui/icons-material';
import { TextWithLinks } from '../../DataFormatter';

export const TitleLink = ({ href, children, ...props }) => {
  return (
    <a
      href={href}
      target={'_blank'}
      rel="noopener noreferrer"
      style={{ color: 'inherit' }}
      {...props}
    >
      <Typography
        variant="h5"
        style={{
          textDecorationLine: 'underline',
          cursor: 'pointer',
          marginBottom: '0.5rem',
          fontWeight: 'bolder !important',
          textTransform: 'uppercase',
          fontSize: '0.9rem',
        }}
      >
        {children}
        <sup>
          <LaunchIcon style={{ width: '1rem', height: '1rem' }} />
        </sup>
      </Typography>
    </a>
  );
};

export const EmptyState = ({ event }) => {
  return (
    <Typography
      variant="body1"
      style={{
        marginBlock: '0.5rem',
      }}
    >
      {' '}
      {
        <TextWithLinks
          variant="body1"
          style={{
            wordWrap: 'break-word',
          }}
          text={event.description || ''}
        ></TextWithLinks>
      }{' '}
    </Typography>
  );
};

export const DataToFileLink = ({ data }) => {
  // convert the trace to a file
  const dataString = typeof data === 'string' ? data : JSON.stringify(data, null, 2);
  const file = new File([dataString], 'trace.txt', { type: 'text/plain' });

  return (
    <TitleLink href={URL.createObjectURL(file)} download="trace.txt">
      Download Trace
    </TitleLink>
  );
};
