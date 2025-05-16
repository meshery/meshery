import { useMemo } from 'react';

const useTestIDsGenerator = (prefix) => {
  return useMemo(() => (identifier) => `${prefix}-${identifier}`, [prefix]);
};

export default useTestIDsGenerator;
