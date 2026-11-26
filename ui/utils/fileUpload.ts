const findSelectedFile = (inputs?: Iterable<HTMLInputElement> | null): File | undefined => {
  if (!inputs || typeof (inputs as Iterable<HTMLInputElement>)[Symbol.iterator] !== 'function') {
    return undefined;
  }

  for (const input of inputs) {
    const file = input?.files?.[0];
    if (file) return file;
  }

  return undefined;
};

export const findSelectedFileInDialog = (selector = 'input[type="file"]'): File | undefined => {
  if (typeof document === 'undefined') return undefined;

  try {
    const dialogs = Array.from(document.querySelectorAll<HTMLElement>('[role="dialog"]')).reverse();

    for (const dialog of dialogs) {
      const file = findSelectedFile(dialog.querySelectorAll<HTMLInputElement>(selector));
      if (file) {
        return file;
      }
    }

    return findSelectedFile(document.querySelectorAll<HTMLInputElement>(selector));
  } catch {
    return undefined;
  }
};

export const readFileAsBytes = (file?: File | Blob | null): Promise<number[]> =>
  new Promise((resolve, reject) => {
    if (!file || typeof file.slice !== 'function') {
      return reject(new TypeError('A valid File or Blob must be provided'));
    }

    if (typeof FileReader === 'undefined') {
      return reject(new Error('FileReader is not supported in this environment'));
    }

    const reader = new FileReader();
    reader.onload = () => {
      const buffer = reader.result as ArrayBuffer | null;
      if (!buffer) {
        return resolve([]);
      }
      resolve(Array.from(new Uint8Array(buffer)));
    };
    reader.onerror = () => reject(reader.error || new Error('Failed to read file'));
    reader.onabort = () => reject(new Error('File reading was aborted'));

    try {
      reader.readAsArrayBuffer(file);
    } catch (err) {
      reject(err);
    }
  });
