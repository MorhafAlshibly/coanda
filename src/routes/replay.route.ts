import express from "express";
import { createReplayHandler, getReplayHandler } from "../controllers/replay.controller";
import validator from "../middlewares/validator";
import { createReplaySchema, getReplaySchema } from "../schemas/replay.schema";

let router = express.Router();

router.post("/create", validator(createReplaySchema), createReplayHandler);
router.get("/get", validator(getReplaySchema), getReplayHandler);

export default router;
