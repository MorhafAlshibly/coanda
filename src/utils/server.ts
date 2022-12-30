import express from "express";
import cors from "cors";
import config from "config";
import routes from "../routes/index.route";
import parse from "../middlewares/parsing";

export const server = () => {
  const app = express();

  // Middleware
  app.use(cors());
  app.use(
    express.json({
      limit: config.get<string>("sizeLimit"),
    })
  );
  app.use(parse);
  routes(app);

  return app;
};
