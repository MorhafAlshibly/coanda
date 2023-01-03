import config from "config";
import mongoose from "mongoose";
import logger from "./logger";

// Connecting to MongoDB database
const connect = async () => {
	try {
		mongoose.set("strictQuery", false);
		// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
		await mongoose.connect(process.env.MONGOURI!);
		logger.info(config.get<string>("mongodb.message"));
	} catch (error) {
		process.exit(1);
	}
};

export default connect;
