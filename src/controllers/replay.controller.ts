import { Request, Response } from "express";
import { CreateReplayInput, GetReplayInput } from "../schemas/replay.schema";
import { createReplay, getReplay } from "../services/replay.service";
import { CreateReplaySuccess, GetReplaySuccess, GetReplayFail } from "../responses/replay.response";
import { ErrorRes } from "../responses/index.response";

export const createReplayHandler = async (req: Request<Record<string, never>, Record<string, never>, CreateReplayInput["body"]>, res: Response) => {
	try {
		const replay = await createReplay(req.body);
		return new CreateReplaySuccess(replay._id).send(res);
	} catch (e: any) {
		// Temporary database error
		return new ErrorRes(e.message).send(res);
	}
};

export const getReplayHandler = async (req: Request<Record<string, never>, Record<string, never>, GetReplayInput["body"]>, res: Response) => {
	try {
		const replay = await getReplay(req.body);
		if (!replay) return new GetReplayFail("replay_not_found").send(res);
		return new GetReplaySuccess(replay).send(res);
	} catch (e: any) {
		// Temporary database error
		return new ErrorRes(e.message).send(res);
	}
};
