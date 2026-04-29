import { useMemo } from 'react';

const useTestIDsGenerator = (prefix: string): ((_identifier: string) => string) => {
  return useMemo(() => (identifier: string) => `${prefix}-${identifier}`, [prefix]);
};

export default useTestIDsGenerator;
