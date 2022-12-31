import { Request, Response } from "express";
import { cacheCreate } from "../middlewares/cache";
import { CreateReplayInput, GetReplayInput } from "../schemas/replay.schema";
import { createReplay, getReplay } from "../services/replay.service";
import { Invalid, Error } from "../utils/responder";
import { CreateReplaySuccess, GetReplaySuccess, GetReplayFail } from "../schemas/replay.schema";

export const createReplayHandler = async (req: Request<{}, {}, CreateReplayInput["body"]>, res: Response) => {
  try {
    const replay = await createReplay(req.body);
    return new CreateReplaySuccess(replay._id).send(res);
  } catch (e: any) {
    return new Error(e.message).send(res);
  }
};

export const getReplayHandler = async (req: Request<{}, {}, GetReplayInput["body"]>, res: Response) => {
  try {
    const replay = await getReplay(req.body);
    if (!replay) return new GetReplayFail("replay_not_found").send(res);
    await cacheCreate(replay._id.toString(), replay);
    return new GetReplaySuccess(replay).send(res);
  } catch (e: any) {
    return new Error(e.message).send(res);
  }
};
