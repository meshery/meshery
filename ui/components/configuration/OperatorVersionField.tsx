import { useEffect, useState } from 'react';
import { Autocomplete, TextField } from '@sistent/sistent';
import { fitWidth } from './controllersConfigForm.shared';

const CHART_PACKAGE = 'https://artifacthub.io/api/v1/packages/helm/meshery/meshery-operator';

export const useOperatorChartVersions = (): string[] => {
  const [versions, setVersions] = useState<string[]>([]);
  useEffect(() => {
    if (typeof fetch !== 'function') return;
    const ac = new AbortController();
    fetch(CHART_PACKAGE, { signal: ac.signal })
      .then((res) => (res.ok ? res.json() : Promise.reject()))
      .then((data) => {
        const tags = (data?.available_versions ?? [])
          .map((entry: { version?: string }) => entry.version)
          .filter((version: string | undefined): version is string => Boolean(version));
        setVersions(tags);
      })
      .catch(() => undefined);
    return () => ac.abort();
  }, []);
  return versions;
};

type OperatorVersionFieldProps = {
  value: string | undefined;
  placeholder: string;
  disabled: boolean;
  onChange: (next: string | undefined) => void;
};

export default function OperatorVersionField({
  value,
  placeholder,
  disabled,
  onChange,
}: OperatorVersionFieldProps) {
  const versions = useOperatorChartVersions();
  return (
    <Autocomplete
      freeSolo
      disabled={disabled}
      options={versions}
      value={value ?? ''}
      onChange={(_, next) => onChange(next === '' || next == null ? undefined : String(next))}
      onInputChange={(_, next, reason) => {
        if (reason === 'reset') return;
        onChange(next.trim() === '' ? undefined : next);
      }}
      renderInput={(params) => <TextField {...params} size="small" placeholder={placeholder} />}
      sx={{
        width: fitWidth(value, placeholder, 'Inherit (server release)'),
        maxWidth: '100%',
        '& input::placeholder': {
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
        },
      }}
    />
  );
}
