import zodToJsonSchema from "zod-to-json-schema";

import { createReplaySchema, getReplaySchema } from "../src/schemas/replay.schema";

export default {
	...zodToJsonSchema(createReplaySchema.shape.body, "CreateReplayInput").definitions,
	...zodToJsonSchema(getReplaySchema.shape.body, "GetReplayInput").definitions,
};
