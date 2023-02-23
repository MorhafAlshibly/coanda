/* eslint-disable indent */
import fs from "fs";
import config from "config";

const services = Object.keys(config.get<string>("microservices")).map((name) => ({
	name,
	port: config.get<string>("microservices." + name + ".port"),
}));

const dockerCompose = `
version: '3.9'

services:
${services
	.map(
		({ name, port }) => `
  ${name}:
    env_file: .env
    build:
      context: .
      args:
        RUN_PATH: "build/src/microservices/${name}/index.js"
    ports:
      - "${port}:${port}"
`
	)
	.join("")}
`;

fs.writeFileSync("docker-compose.yml", dockerCompose);
