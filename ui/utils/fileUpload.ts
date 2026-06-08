export const findSelectedFileInDialog = (selector = 'input[type="file"]'): File | undefined => {
  if (typeof document === 'undefined') return undefined;

  const dialog = document.querySelector<HTMLElement>('[role="dialog"]');

  if (dialog) {
    const inputs = dialog.querySelectorAll<HTMLInputElement>(selector);
    for (const input of Array.from(inputs)) {
      const file = input.files?.[0];
      if (file) return file;
    }
  }

  const inputs = document.querySelectorAll<HTMLInputElement>(selector);
  for (const input of Array.from(inputs)) {
    const file = input.files?.[0];
    if (file) return file;
  }

  return undefined;
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
