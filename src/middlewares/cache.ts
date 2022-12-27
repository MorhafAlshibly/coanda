import { NextFunction, Request, Response } from "express";
import { createClient } from "redis";
const redisClient = createClient();

export const cacheMiddleware = (key: string) => async (req: Request, res: Response, next: NextFunction) => {
  try {
    await redisClient.connect();
    const cachedData = await redisClient.get(req.body[key].toString());
    await redisClient.disconnect();
    if (cachedData) return res.jsend.success(JSON.parse(cachedData));
    next();
  } catch (e: any) {
    return res.status(500).jsend.error(e.message);
  }
};

export const cacheCreate = async (key: string, data: Object) => {
  try {
    await redisClient.connect();
    await redisClient.set(key, JSON.stringify(data));
    await redisClient.disconnect();
  } catch (e: any) {
    throw new Error(e);
  }
};
