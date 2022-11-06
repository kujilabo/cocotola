export const extractErrorMessage = (err: any): string => {
  console.log('err', err);
  if (err.response) {
    if (err.response.data) {
      if (err.response.data.message) {
        return err.response.data.message;
      }
    }
    if (err.response.statusText) {
      return err.response.statusText;
    }
  }
  return 'Error';
};
