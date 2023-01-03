import type { Config } from "@jest/types";

const jestConfig: Config.InitialOptions = {
	preset: "ts-jest",
	testEnvironment: "node",
	testMatch: ["**/__tests__/**/*.[jt]s?(x)"],
	setupFiles: ["dotenv/config"],
	verbose: true,
	forceExit: true,
	clearMocks: true,
	resetMocks: true,
	restoreMocks: true,
};

export default jestConfig;
