const findSelectedFile = (inputs: Iterable<HTMLInputElement>): File | undefined => {
  for (const input of inputs) {
    const file = input.files?.[0];
    if (file) return file;
  }

  return undefined;
};

export const findSelectedFileInDialog = (selector = 'input[type="file"]'): File | undefined => {
  if (typeof document === 'undefined') return undefined;

  const dialogs = Array.from(document.querySelectorAll<HTMLElement>('[role="dialog"]')).reverse();

  for (const dialog of dialogs) {
    const file = findSelectedFile(dialog.querySelectorAll<HTMLInputElement>(selector));
    if (file) {
      return file;
    }
  }

  return findSelectedFile(document.querySelectorAll<HTMLInputElement>(selector));
};

export const readFileAsBytes = (file: File): Promise<number[]> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const buffer = reader.result as ArrayBuffer;
      resolve(Array.from(new Uint8Array(buffer)));
    };
    reader.onerror = () => reject(reader.error);
    reader.readAsArrayBuffer(file);
  });
