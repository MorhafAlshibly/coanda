import { NextFunction, Request, Response } from "express";
import { UnauthorizedRes } from "../utils/responses";

// Middleware to authenticate the request
const auth = (req: Request, res: Response, next: NextFunction) => {
	if (req.headers.apikey !== process.env.APIKEY) {
		return new UnauthorizedRes().send(res);
	}
	next();
};

export default auth;
