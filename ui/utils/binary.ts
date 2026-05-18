/**
 * Converts an ArrayBuffer to a Base64 string.
 * Uses a chunked approach to prevent stack overflows with large buffers.
 */
export const arrayBufferToBase64 = (buffer: ArrayBuffer): Promise<string> => {
  return new Promise((resolve, reject) => {
    try {
      const bytes = new Uint8Array(buffer);
      let binary = '';
      const len = bytes.byteLength;

      // Process in chunks to avoid call stack limits in some environments
      const CHUNK_SIZE = 8192;
      for (let i = 0; i < len; i += CHUNK_SIZE) {
        const chunk = bytes.subarray(i, i + CHUNK_SIZE);

        binary += String.fromCharCode.apply(null, chunk as unknown as number[]);
      }
      resolve(btoa(binary));
    } catch (e) {
      reject(e);
    }
  });
};

/**
 * Converts a Base64 string back to a Uint8Array.
 */
export const base64ToUint8Array = (base64: string): Uint8Array => {
  const binaryString = atob(base64);
  const len = binaryString.length;
  const bytes = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  return bytes;
};
