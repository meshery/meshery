import React, { useState, useEffect } from 'react';
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
  const [downloadUrl, setDownloadUrl] = useState<string>('');

  useEffect(() => {
    // convert the trace to a file
    const dataString = typeof data === 'string' ? data : JSON.stringify(data ?? '', null, 2) || '';
    const file = new File([dataString], 'trace.txt', { type: 'text/plain' });
    const url = URL.createObjectURL(file);
    setDownloadUrl(url);

    return () => {
      URL.revokeObjectURL(url);
    };
  }, [data]);

  return (
    <TitleLink href={downloadUrl} download="trace.txt">
      Download Trace
    </TitleLink>
  );
};
