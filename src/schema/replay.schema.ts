import { object, any } from "zod";

const createReplaySchema = object({
  body: object({
    data: any({
      required_error: "Replay data is required",
    }),
  }),
});

export default createReplaySchema;
