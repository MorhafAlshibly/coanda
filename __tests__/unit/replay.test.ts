import mongoose from "mongoose";
import supertest from "supertest";
import { server } from "../../src/utils/server";
import { jest } from "@jest/globals";
import * as ReplayService from "../../src/services/replay.service";

const app = server();

describe("Replay", () => {
  describe("Create Replay", () => {
    (ReplayService.createReplay as jest.Mock) = jest.fn();

    describe("Given that the data field is not sent", () => {
      it("Should return code 400 and 'required' message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            expireAt: new Date(9999999999999),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).not.toHaveBeenCalled();
        expect(body.errors[0].message).toEqual("Required");
      });
    });

    describe("Given that the data field is empty", () => {
      it("Should return code 400 and 'required' message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            data: {},
            expireAt: new Date(9999999999999),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).not.toHaveBeenCalled();
        expect(body.errors[0].message).toEqual("Required");
      });
    });

    describe("Given that expireAt is invalid", () => {
      it("Should return code 400 and 'invalid' error message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            data: { replay: true },
            expireAt: "invaliddate",
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).not.toHaveBeenCalled();
        expect(body.errors[0].code).toEqual("invalid_date");
      });
    });

    describe("Given that expireAt is not in the future", () => {
      it("Should return code 400 and 'too_small' error message", async () => {
        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            data: { replay: true },
            expireAt: 1000,
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(ReplayService.createReplay).not.toHaveBeenCalled();
        expect(body.errors[0].code).toEqual("too_small");
      });
    });

    describe("Given that expireAt is not given", () => {
      it("Should return code 201 and replay id", async () => {
        (ReplayService.createReplay as jest.Mock).mockReturnValue({ data: { replay: true }, expireAt: new Date(), _id: new mongoose.Types.ObjectId().toString() });
        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            data: { replay: true },
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(201);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(1);
        expect(mongoose.Types.ObjectId.isValid(body._id)).toBeTruthy();
      });
    });

    describe("Given that valid data and expireAt is given", () => {
      it("Should return code 201 and replay id", async () => {
        (ReplayService.createReplay as jest.Mock).mockReturnValue({ data: { replay: true }, expireAt: new Date(), _id: new mongoose.Types.ObjectId().toString() });
        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            data: { replay: true },
            expireAt: new Date(9999999999999),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(201);
        expect(ReplayService.createReplay).toHaveBeenCalledTimes(1);
        expect(mongoose.Types.ObjectId.isValid(body._id)).toBeTruthy();
      });
    });
  });
});
