import { Schema, model, Document } from "mongoose";

export interface ReplayDocument extends Document {
	data: object;
	userId: number;
}

const replaySchema = new Schema({
	data: {
		type: Object,
	},
	userId: {
		type: Number,
	},
});

const ReplayModel = model("Replay", replaySchema);

export default ReplayModel;
