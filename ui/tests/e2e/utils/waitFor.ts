export const waitFor = async <T>(func: () => T | Promise<T>, timeout = 5000): Promise<T> => {
  return new Promise((resolve, reject) => {
    let pollingInterval: ReturnType<typeof setInterval>;

    const pollFn = async () => {
      try {
        const result = await func();
        if (result) {
          clearInterval(pollingInterval);
          resolve(result);
        }
      } catch (error) {
        clearInterval(pollingInterval);
        reject(error);
      }
    };

    pollingInterval = setInterval(pollFn, 200);

    setTimeout(() => {
      clearInterval(pollingInterval);
      reject(`Timeout after ${timeout} ms for condition ${func.toString()}`);
    }, timeout);
  });
};
