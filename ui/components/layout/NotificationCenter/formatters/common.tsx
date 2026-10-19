import React from 'react';
import { LaunchIcon, Typography } from '@sistent/sistent';
import { TextWithLinks } from '../../../data-formatter';

type TitleLinkProps = {
  href: string;
  children: React.ReactNode;
} & React.AnchorHTMLAttributes<HTMLAnchorElement>;

export const TitleLink = ({ href, children, ...props }: TitleLinkProps) => {
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
          <LaunchIcon style={{ width: '16px', height: '16px' }} />
        </sup>
      </Typography>
    </a>
  );
};

/**
 * Formats an error code into a globally unique universal identifier (mnemonic + numeric code)
 * per the Meshery error conventions (e.g., '1562' -> 'meshery-server-1562').
 * If the code already contains a component mnemonic prefix (e.g. 'meshery-server-1133' or 'meshkit-1001'),
 * it is returned as is.
 */
export const formatErrorCode = (code?: unknown, component = 'meshery-server'): string => {
  if (code === undefined || code === null) return '';
  const trimmed = String(code).trim();
  if (!trimmed) return '';
  if (trimmed.includes('-')) {
    return trimmed;
  }
  if (/^\d+$/.test(trimmed)) {
    const prefix = component ? component.replace(/_/g, '-') : 'meshery-server';
    return `${prefix}-${trimmed}`;
  }
  return trimmed;
};

export const EmptyState = ({ event }: { event: { description?: string } }) => {
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

export const DataToFileLink = ({ data }: { data: unknown }) => {
  // convert the trace to a file
  const dataString = typeof data === 'string' ? data : JSON.stringify(data, null, 2);
  const file = new File([dataString], 'trace.txt', { type: 'text/plain' });

  return (
    <TitleLink href={URL.createObjectURL(file)} download="trace.txt">
      Download Trace
    </TitleLink>
  );
};
