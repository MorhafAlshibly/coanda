import express from "express";
import cors from "cors";
import config from "config";
import helmet from "helmet";
import routes from "../routes/index.route";
import parse from "../middlewares/parsing";

export const server = () => {
	const app = express();

	// Middleware
	app.use(cors());
	app.use(helmet());
	app.use(
		express.json({
			limit: config.get<string>("express.sizeLimit"),
		})
	);
	app.use(parse);
	// Add routes
	routes(app);

	return app;
};
