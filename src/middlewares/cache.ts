import { NextFunction, Request, Response } from "express";
import { createClient } from "redis";
import { success } from "../utils/responder";
const redisClient = createClient({ url: process.env.REDISURI });

export const cacheMiddleware = (key: string) => async (req: Request, res: Response, next: NextFunction) => {
  try {
    await redisClient.connect();
    const cachedData = await redisClient.get(req.body[key].toString());
    await redisClient.disconnect();
    if (cachedData) return success(res, JSON.parse(cachedData));
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
