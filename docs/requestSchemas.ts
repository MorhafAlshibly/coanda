import zodToJsonSchema from "zod-to-json-schema";

import { createReplaySchema, getReplaySchema } from "../src/microservices/replays/schema";

export default {
	...zodToJsonSchema(createReplaySchema.shape.body, "CreateReplayInput").definitions,
	...zodToJsonSchema(getReplaySchema.shape.body, "GetReplayInput").definitions,
};
