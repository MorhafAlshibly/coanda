import express from "express";
import cors from "cors";
import config from "config";
import replay from "../routes/replay.route";

export const server = () => {
  const app = express();
  app.use(cors());
  app.use(
    express.json({
      limit: config.get<string>("sizeLimit"),
    })
  );

  // Routes
  app.use("/replay", replay);

  return app;
};
