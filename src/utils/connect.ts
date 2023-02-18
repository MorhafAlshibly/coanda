import config from "config";
import mongoose from "mongoose";
import logger from "./logger";
import { cosmosUri } from "./secrets";

// Connecting to MongoDB database
const connect = async () => {
	try {
		mongoose.set("strictQuery", false);
		// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
		await mongoose.connect((await cosmosUri())!);

		logger.info(config.get<string>("mongodb.message"));
	} catch (error) {
		process.exit(1);
	}
};

export default connect;
