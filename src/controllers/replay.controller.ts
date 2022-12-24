import { Request, Response } from "express";
import { TypeOf } from "zod";
import createReplaySchema from "../schemas/replay.schema";
import { createReplay } from "../services/replay.service";
import logger from "../utils/logger";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return res.status(201).send({ success: true, _id: replay._id });
  } catch (e: any) {
    logger.error(e);
    return res.status(500).send({ success: false, message: e.message });
  }
};

export type CreateReplayInput = TypeOf<typeof createReplaySchema>;
