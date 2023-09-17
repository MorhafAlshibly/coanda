import { itemsArchiver } from "./itemsArchiver";
import TableStorage from "../storage/table";
import { TableEntityResult } from "@azure/data-tables";
import { InvocationContext, Timer } from "@azure/functions";
import { jest, describe, expect, it, beforeEach } from "@jest/globals";

jest.mock("../storage/table");

const context = {
	log: jest.fn(),
} as unknown as InvocationContext;

const timer = jest.fn() as unknown as Timer;

beforeEach(() => {
	jest.clearAllMocks();
});

describe("itemsArchiver", () => {
	it("should delete expired items", async () => {
		TableStorage.prototype.getAllItems = jest.fn().mockReturnValue([
			{
				Expire: "2021-01-01T00:00:00.000Z",
				partitionKey: "1",
				rowKey: "1",
			} as unknown as TableEntityResult<Record<string, unknown>>,
		]) as jest.MockedFunction<typeof TableStorage.prototype.getAllItems>;
		// Run the function
		await itemsArchiver(timer, context);
		// Verify that the item was deleted
		expect(TableStorage.prototype.deleteItem).toHaveBeenCalledTimes(1);
		expect(TableStorage.prototype.deleteItem).toHaveBeenCalledWith("1", "1");
	});

	it("should not delete items that are not expired", async () => {
		TableStorage.prototype.getAllItems = jest.fn().mockReturnValue([
			{
				Expire: "9999-01-01T00:00:00.000Z",
				partitionKey: "1",
				rowKey: "1",
			} as unknown as TableEntityResult<Record<string, unknown>>,
		]) as jest.MockedFunction<typeof TableStorage.prototype.getAllItems>;
		// Run the function
		await itemsArchiver(timer, context);
		// Verify that the item was not deleted
		expect(TableStorage.prototype.deleteItem).toHaveBeenCalledTimes(0);
	});

	it("should not delete items that do not have an expire timestamp", async () => {
		TableStorage.prototype.getAllItems = jest.fn().mockReturnValue([
			{
				partitionKey: "1",
				rowKey: "1",
			} as unknown as TableEntityResult<Record<string, unknown>>,
		]) as jest.MockedFunction<typeof TableStorage.prototype.getAllItems>;
		// Run the function
		await itemsArchiver(timer, context);
		// Verify that the item was not deleted
		expect(TableStorage.prototype.deleteItem).toHaveBeenCalledTimes(0);
	});
});
