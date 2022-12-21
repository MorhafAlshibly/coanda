import { Request, Response } from "express";
import { TypeOf } from "zod";
import createReplaySchema from "../schema/replay.schema";
import createReplay from "../service/replay.service";
import logger from "../utils/logger";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return res.send(replay);
  } catch (e: any) {
    logger.error(e);
    return res.status(400).send(e.message);
  }
};

export type CreateReplayInput = TypeOf<typeof createReplaySchema>;
