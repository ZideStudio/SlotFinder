import {
  getAuthStatus200,
  getAuthStatus401,
} from "../../../mocks/handlers/authStatusHandlers";
import { Given, When } from "../../common/fixtures";

Given("I am signed in", async ({ network }) => {
  network.use(getAuthStatus200);
});

Given("I am not signed in", async ({ network }) => {
  network.use(getAuthStatus401);
});

When("I open the home page", async ({ page }) => {
  await page.goto("/");
});
