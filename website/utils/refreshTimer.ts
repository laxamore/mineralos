let refreshTimerID: NodeJS.Timer;

export const getRefreshTimerID = () => {
  return refreshTimerID;
};

export const setRefreshTimerCallback = (callback: Function) => {
  clearRefreshTimer();
  refreshTimerID = setInterval(() => {
    callback();
  }, 5000);
};

export const clearRefreshTimer = () => {
  if (refreshTimerID !== undefined) {
    clearInterval(refreshTimerID);
  }
};
