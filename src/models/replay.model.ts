import { Schema, model, Document } from "mongoose";

export interface ReplayDocument extends Document {
  data: Object;
  expireAt: Date;
}

const replaySchema = new Schema({
  data: {
    type: Object,
  },
  expireAt: {
    type: Date,
  },
});

const ReplayModel = model("Replay", replaySchema);

export default ReplayModel;
