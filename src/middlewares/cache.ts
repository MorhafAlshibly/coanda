import { NextFunction, Request, Response } from "express";
import { createClient } from "redis";
import { GetReplaySuccess } from "../responses/replay.response";
const redisClient = createClient({ url: process.env.REDISURI });

export const cacheMiddleware = (key: string) => async (req: Request, res: Response, next: NextFunction) => {
  try {
    await redisClient.connect();
    const cachedData = await redisClient.get(req.body[key].toString());
    await redisClient.disconnect();
    if (cachedData) return new GetReplaySuccess(JSON.parse(cachedData)).send(res);
    next();
  } catch (e: any) {
    await redisClient.disconnect();
    next();
  }
};

export const cacheCreate = async (key: string, data: Object) => {
  try {
    await redisClient.connect();
    await redisClient.set(key, JSON.stringify(data));
    await redisClient.disconnect();
  } catch (e: any) {
    await redisClient.disconnect();
  }
};
