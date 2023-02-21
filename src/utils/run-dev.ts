import concurrently, { ConcurrentlyCommandInput } from "concurrently";
import { fdir } from "fdir";

const paths = new fdir().withFullPaths().glob("./**/index.ts").crawl("./src/microservices").sync() as string[];

const commands = [] as ConcurrentlyCommandInput[];

for (const path of paths) {
	const directories = path.split("\\");
	const name = directories[directories.length - 2];
	commands.push({
		command: "ts-node-dev -r dotenv/config " + path,
		name: name.charAt(0).toUpperCase() + name.slice(1),
	});
}

concurrently(commands, {
	prefix: "name",
	killOthers: ["failure", "success"],
});
