import React, { type ReactNode } from 'react';

export type AutocompleteInputProps = {
  startAdornment?: ReactNode;
  [key: string]: unknown;
};

export const buildCustomInputProps = (
  paramsInputProps: AutocompleteInputProps | undefined,
  leadingAdornment: ReactNode,
): AutocompleteInputProps => {
  const inputProps = paramsInputProps ?? {};
  return {
    ...inputProps,
    startAdornment: (
      <>
        {leadingAdornment}
        {inputProps.startAdornment}
      </>
    ),
  };
};
