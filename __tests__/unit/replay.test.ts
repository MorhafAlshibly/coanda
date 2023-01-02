import mongoose from "mongoose";
import supertest from "supertest";
import { server } from "../../src/utils/server";
import { jest } from "@jest/globals";
import * as ReplayService from "../../src/services/replay.service";
import * as CacheService from "../../src/middlewares/cache";

const app = server();

describe("Replay", () => {
  describe("Create Replay", () => {
    (ReplayService.createReplay as jest.Mock) = jest.fn();

    describe("Given that the data field is not given", () => {
      it("Should return code 400 and 'required' message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay")
          .send({
            expireAt: new Date(9999999999999),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].message).toEqual("Required");
      });
    });

    describe("Given that the data field is empty", () => {
      it("Should return code 400 and 'required' message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay")
          .send({
            data: {},
            expireAt: new Date(9999999999999),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].message).toEqual("Required");
      });
    });

    describe("Given that expireAt is invalid", () => {
      it("Should return code 400 and 'invalid' message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay")
          .send({
            data: { replay: true },
            expireAt: "invaliddate",
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].code).toEqual("invalid_date");
      });
    });

    describe("Given that expireAt is not in the future", () => {
      it("Should return code 400 and 'too_small' message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay")
          .send({
            data: { replay: true },
            expireAt: 1000,
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].code).toEqual("too_small");
      });
    });

    describe("Given that expireAt is not given", () => {
      it("Should return code 200 and replay _id", async () => {
        const replayId = new mongoose.Types.ObjectId().toString();
        (ReplayService.createReplay as jest.Mock).mockReturnValueOnce({ _id: replayId });
        const { statusCode, body } = await supertest(app)
          .post("/replay")
          .send({
            data: { replay: true },
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(200);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(1);
        expect(body.data).toEqual(replayId);
      });
    });

    describe("Given that valid data and expireAt is given", () => {
      it("Should return code 200 and replay _id", async () => {
        const replayId = new mongoose.Types.ObjectId().toString();
        (ReplayService.createReplay as jest.Mock).mockReturnValueOnce({ _id: replayId });
        const { statusCode, body } = await supertest(app)
          .post("/replay")
          .send({
            data: { replay: true },
            expireAt: new Date(9999999999999),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(200);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(1);
        expect(body.data).toEqual(replayId);
      });
    });
  });

  describe("Get Replay", () => {
    (ReplayService.getReplay as jest.Mock) = jest.fn();

    describe("Given that the _id field is not given", () => {
      it("Should return code 400 and 'required' message", async () => {
        const { statusCode, body } = await supertest(app).get("/replay").send({}).set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.getReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].message).toEqual("Required");
      });
    });

    describe("Given that the _id field is not a valid id", () => {
      it("Should return code 400 and 'invalid' message", async () => {
        const { statusCode, body } = await supertest(app).get("/replay").send({ _id: "invalidid" }).set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.getReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].code).toEqual("invalid_type");
      });
    });

    describe("Given that the _id field is the wrong data type", () => {
      it("Should return code 400 and 'invalid' message", async () => {
        const { statusCode, body } = await supertest(app).get("/replay").send({ _id: 890476425 }).set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.getReplay).toHaveBeenCalledTimes(0);
        expect(body.data[0].code).toEqual("invalid_type");
      });
    });

    describe("Given that the _id field is valid", () => {
      it("Should return code 200 and replay data", async () => {
        (ReplayService.getReplay as jest.Mock).mockReturnValueOnce({ _id: new mongoose.Types.ObjectId().toString() });
        (CacheService.cacheCreate as jest.Mock) = jest.fn().mockReturnValueOnce(true);
        const { statusCode, body } = await supertest(app).get("/replay").send({ _id: new mongoose.Types.ObjectId().toString() }).set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(200);
        expect(ReplayService.getReplay).toHaveBeenCalledTimes(1);
      });
    });
  });
});
