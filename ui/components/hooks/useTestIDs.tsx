import { useMemo } from 'react';

const useTestIDsGenerator = (prefix: string): ((_identifier: string) => string) => {
  return useMemo(() => (_identifier: string) => `${prefix}-${_identifier}`, [prefix]);
};

export default useTestIDsGenerator;
