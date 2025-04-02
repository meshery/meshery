export const waitFor = async (func, timeout = 5000) => {
  return new Promise((resolve, reject) => {
    const pollFn = async (pollingInterval) => {
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

    const pollingInterval = setInterval(function () {
      pollFn(this);
    }, 200);
    setTimeout(() => {
      clearInterval(pollingInterval);
      reject(`Timeout after ${timeout} ms for condition ${func}`);
    }, timeout);
  });
};
