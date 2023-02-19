import express from "express";
import cors from "cors";
import config from "config";
import helmet from "helmet";
import parse from "../middlewares/parsing";
import auth from "../middlewares/auth";
import connect from "../utils/connect";
import logger from "../utils/logger";

export const server = () => {
	const app = express();
	const port = config.get<number>("express.port");

	// Middleware
	app.use(cors());
	app.use(helmet());
	app.use(
		express.json({
			limit: config.get<string>("express.sizeLimit"),
		})
	);
	app.use(parse);
	app.use(auth);

	// Start express app on a port
	const listener = app.listen(port, async () => {
		logger.info(config.get<string>("express.message"));
		await connect();
	});
	listener.setTimeout(config.get<number>("express.timeout"));

	return app;
};

export default server;
