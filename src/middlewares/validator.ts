import config from "config";
import { Request, Response, NextFunction } from "express";
import { AnyZodObject, string, ZodIssueCode, ZodError } from "zod";
import { InvalidRes } from "../responses/index.response";

// Middleware to validate request body based on Zod schema
const validate = (schema: AnyZodObject) => (req: Request, res: Response, next: NextFunction) => {
	try {
		// Validating API key
		const apiSchema = schema.extend({
			apikey: string().superRefine((val, ctx) => {
				if (val != process.env.APIKEY) {
					ctx.addIssue({
						code: ZodIssueCode.custom,
						message: config.get<string>("auth.message"),
					});
				}
			}),
		});
		// Parsing request body into validation schema
		req.body = apiSchema.parse({
			body: req.body,
			apikey: req.headers.apikey,
		}).body;
		next();
	} catch (e: unknown) {
		if (e instanceof ZodError) return new InvalidRes(e.errors).send(res);
	}
};

export default validate;
