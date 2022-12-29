import { Request, Response } from "express";
import { cacheCreate } from "../middlewares/cache";
import { CreateReplayInput, GetReplayInput } from "../schemas/replay.schema";
import { createReplay, getReplay } from "../services/replay.service";
import logger from "../utils/logger";
import { failModel } from "../utils/responses";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return res.jsend.success(replay._id);
  } catch (e: any) {
    logger.error(e);
    return res.status(500).jsend.error(e.message);
  }
};

export const getReplayHandler = async (req: Request<{}, {}, GetReplayInput["body"]>, res: Response) => {
  try {
    const replay = await getReplay(req.body);
    if (!replay) return res.status(404).jsend.fail(failModel("not_found", ["body", "_id"], "Data not found"));
    await cacheCreate(replay._id.toString(), replay);
    return res.jsend.success(replay);
  } catch (e: any) {
    logger.error(e);
    return res.status(500).jsend.error(e.message);
  }
};
