import { Request, Response } from "express";
import { cacheCreate } from "../middlewares/cache";
import { CreateReplayInput, GetReplayInput } from "../schemas/replay.schema";
import { createReplay, getReplay } from "../services/replay.service";
import { Success, Fail, Error } from "../utils/responder";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return new Success(replay._id).send(res);
  } catch (e: any) {
    return new Error(e.message).send(res);
  }
};

export const getReplayHandler = async (req: Request<{}, {}, GetReplayInput["body"]>, res: Response) => {
  try {
    const replay = await getReplay(req.body);
    if (!replay) return new Fail("not_found", "Data not found").send(res);
    await cacheCreate(replay._id.toString(), replay);
    return new Success(replay).send(res);
  } catch (e: any) {
    return new Error(e.message).send(res);
  }
};
