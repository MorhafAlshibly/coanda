import dotenv from "dotenv";
dotenv.config();

import config from "config";
import connect from "./src/utils/connect";
import logger from "./src/utils/logger";
import { server } from "./src/utils/server";

const port = config.get<number>("express.port");

const app = server();

const listener = app.listen(port, async () => {
  logger.info(config.get<string>("express.message"));
  await connect();
});
listener.setTimeout(config.get<number>("express.timeout"));
