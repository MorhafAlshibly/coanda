import logger from "pino";

// Pino logger
const log = logger({
	transport: {
		target: "pino-pretty",
	},
	base: {
		pid: false,
	},
	timestamp: () => `,"time":"${new Date()}"`,
});

export default log;
