import { Request, Response } from "express";
import { CreateReplayInput, GetReplayInput } from "../schemas/replay.schema";
import { createReplay, getReplay } from "../services/replay.service";
import logger from "../utils/logger";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return res.jsend.success(replay._id);
  } catch (e: any) {
    logger.error(e);
    return res.jsend.error(e.message);
  }
};

export const getReplayHandler = async (req: Request<{}, {}, GetReplayInput["body"]>, res: Response) => {
  try {
    const replay = await getReplay(req.body);
    return res.jsend.success(replay);
  } catch (e: any) {
    logger.error(e);
    return res.jsend.error(e.message);
  }
};
