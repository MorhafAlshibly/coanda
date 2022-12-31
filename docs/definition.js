const { name, version } = require("../package.json");

module.exports = {
  openapi: "3.1.0",
  info: {
    title: name,
    version: version,
  },
};
