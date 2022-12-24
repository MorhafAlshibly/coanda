import { object, date, preprocess } from "zod";

const createReplaySchema = object({
  body: object({
    data: object({}).passthrough(),
    expireAt: preprocess(
      (arg) => {
        if (typeof arg == "number" || typeof arg == "string" || arg instanceof Date) return new Date(arg);
      },
      date().min(new Date(), {
        message: "Must be a date in the future",
      })
    ).default(new Date(9999999999999)),
  }),
});

export default createReplaySchema;
