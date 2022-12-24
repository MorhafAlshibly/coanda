import express from "express";
import { createReplayHandler } from "../controllers/replay.controller";
import validator from "../middlewares/validator";
import createReplaySchema from "../schemas/replay.schema";

let router = express.Router();

router.post("/create", validator(createReplaySchema), createReplayHandler);

export default router;
