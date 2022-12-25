import express from "express";
import cors from "cors";
import config from "config";
import replay from "../routes/replay.route";
import parse from "../middlewares/parsing";
import jsend from "jsend";

export const server = () => {
  const app = express();
  app.use(jsend.middleware);
  app.use(cors());
  app.use(
    express.json({
      limit: config.get<string>("sizeLimit"),
    })
  );
  app.use(parse);

  // Routes
  app.use("/replay", replay);

  return app;
};
