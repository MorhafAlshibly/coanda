import { Request, Response } from "express";
import { cacheCreate } from "../middlewares/cache";
import { CreateReplayInput, GetReplayInput } from "../schemas/replay.schema";
import { createReplay, getReplay } from "../services/replay.service";
import { success, fail, error } from "../utils/responder";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return success(res, replay._id);
  } catch (e: any) {
    return error(res, e.message);
  }
};

export const getReplayHandler = async (req: Request<{}, {}, GetReplayInput["body"]>, res: Response) => {
  try {
    const replay = await getReplay(req.body);
    if (!replay)
      return fail(res, {
        code: "not_found",
        path: ["body", "_id"],
        message: "Data not found",
      });
    await cacheCreate(replay._id.toString(), replay);
    return success(res, replay);
  } catch (e: any) {
    return error(res, e.message);
  }
};
