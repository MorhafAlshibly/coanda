import config from "config";
import mongoose from "mongoose";
import { object, number, string, ZodIssueCode, TypeOf } from "zod";

// Schema to validate a createReplay request
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
			.describe(config.get<string>("microservices.replays.createReplay.properties.data")),
		userId: number().nonnegative().describe(config.get<string>("microservices.replays.createReplay.properties.userId")),
	}),
});
export type CreateReplayInput = TypeOf<typeof createReplaySchema>;

// Schema to validate a getReplay request
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
			.describe(config.get<string>("microservices.replays.getReplay.properties._id")),
	}),
});

export type GetReplayInput = TypeOf<typeof getReplaySchema>;
