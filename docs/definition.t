const { name, version } = require("../package.json");

module.exports = {
  openapi: "3.0.0",
  info: {
    title: name,
    version: version,
  },
};
