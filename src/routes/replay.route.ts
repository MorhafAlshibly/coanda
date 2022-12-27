import express from "express";
import { cacheMiddleware } from "../middlewares/cache";
import validator from "../middlewares/validator";
import { createReplayHandler, getReplayHandler } from "../controllers/replay.controller";
import { createReplaySchema, getReplaySchema } from "../schemas/replay.schema";
let router = express.Router();

router.post("/create", validator(createReplaySchema), createReplayHandler);
router.get("/get", validator(getReplaySchema), cacheMiddleware("_id"), getReplayHandler);

export default router;
