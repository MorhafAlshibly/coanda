import { Schema, model, Document } from "mongoose";

export interface ReplayDocument extends Document {
  data: Schema.Types.Mixed;
  expireAt: Date;
}

const replaySchema = new Schema({
  data: {
    type: Schema.Types.Mixed,
    required: true,
  },
  expireAt: {
    type: Date,
  },
});

const ReplayModel = model("Replay", replaySchema);

export default ReplayModel;
