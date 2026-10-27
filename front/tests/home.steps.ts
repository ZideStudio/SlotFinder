import { expect } from "@playwright/test";
import {
  getAuthStatus200,
  getAuthStatus401,
} from "../mocks/handlers/authStatusHandlers";
import { Given, Then, When } from "./bdd.steps";

Given("I am signed in", async ({ network }) => {
  network.use(getAuthStatus200);
});

Given("I am not signed in", async ({ network }) => {
  network.use(getAuthStatus401);
});

When("I open the home page", async ({ page }) => {
  await page.goto("/");
});

Then("I see the main heading {string}", async ({ page }, title) => {
  await expect(
    page.getByRole("heading", { name: title, level: 1 }),
  ).toBeVisible();
});
