export default {
  express: {
    port: 5050,
    timeout: 5000,
    sizeLimit: "100mb",
    message: "Coanda API has started",
  },
  mongodb: {
    message: "Connected to Coanda DB",
  },
  replay: {
    createReplay: {
      minDate: "Must be a date in the future",
    },
  },
};
