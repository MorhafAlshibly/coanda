import config from "config";
import mongoose from "mongoose";
import { object, date, string, preprocess, ZodIssueCode, TypeOf } from "zod";

export const createReplaySchema = object({
	body: object({
		data: object({})
			.passthrough()
			.superRefine((val, ctx) => {
				if (Object.keys(val).length == 0)
					ctx.addIssue({
						code: ZodIssueCode.invalid_type,
						expected: "object",
						received: "undefined",
					});
			})
			.describe(config.get<string>("replay.createReplay.properties.data")),
		expireAt: preprocess(
			(arg) => {
				if (typeof arg == "number" || typeof arg == "string" || arg instanceof Date) return new Date(arg);
			},
			date().min(new Date(), {
				message: config.get<string>("replay.createReplay.minDate"),
			})
		)
			.default(new Date(9999999999999))
			.describe(config.get<string>("replay.createReplay.properties.expireAt")),
	}),
});
export type CreateReplayInput = TypeOf<typeof createReplaySchema>;

export const getReplaySchema = object({
	body: object({
		_id: string()
			.transform((val, ctx) => {
				if (mongoose.Types.ObjectId.isValid(val)) return new mongoose.Types.ObjectId(val);
				else
					ctx.addIssue({
						code: ZodIssueCode.invalid_type,
						expected: "string",
						received: "unknown",
					});
			})
			.describe(config.get<string>("replay.getReplay.properties._id")),
	}),
});

export type GetReplayInput = TypeOf<typeof getReplaySchema>;
