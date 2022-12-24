import supertest from "supertest";
import { server } from "../../src/utils/server";
import * as ReplayService from "../../src/services/replay.service";

const app = server();

describe("Replay", () => {
  describe("Create Replay", () => {
    describe("Given that the data field is not sent", () => {
      it("Should return code 400 and error message", async () => {
        const replayServiceMock = jest.spyOn(ReplayService, "createReplay");

        const { statusCode, body } = await supertest(app)
          .post("/replay/create")
          .send({
            expireAt: new Date(),
          })
          .set({ apikey: process.env.APIKEY });

        expect(statusCode).toBe(400);
        expect(replayServiceMock).not.toHaveBeenCalled();
      });
    });
  });
});
