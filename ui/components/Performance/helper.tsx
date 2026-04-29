import { v4 as uuid } from 'uuid';

/**
 * generateTestName takes in test name and service mesh name
 * and generates a random name (if test name is an empty string or is falsy) or
 * will return the given name
 */
export const generateTestName = (name: string, meshName: string): string => {
  if (!name || name.trim() === '') {
    const mesh = meshName === '' || meshName === 'None' ? 'No mesh' : meshName;
    return `${mesh}_${new Date().getTime()}`;
  }

  return name;
};

export function generateUUID(): string {
  return uuid();
}
