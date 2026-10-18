import { expect } from "@playwright/test";
import { Given, When, Then } from "./bdd.steps";
import { getAuthStatus200, getAuthStatus401 } from "../mocks/handlers/authStatusHandlers";

Given("the application is started with the mocked API", async ({ page }) => {
  await page.goto("/");
});

Given("the application is started with auth status {int}", async ({ page, network }, status) => {
  network.use(status === 200 ? getAuthStatus200 : getAuthStatus401);
  await page.goto("/");
  await page.waitForLoadState("networkidle");
});

When("the user opens the home page", async ({ page, network }) => {
  network.use(getAuthStatus401);
  await page.goto("/");
  await page.waitForLoadState("networkidle");
});

When("the user opens the home page with auth status {int}", async ({ page, network }, status) => {
  if (status === 200) {
    network.use(getAuthStatus200);
  } else {
    network.use(getAuthStatus401);
  }

  await page.goto("/");
  await page.waitForLoadState("networkidle");
});

Then("the page contains the title {string}", async ({ page }, title) => {
  await expect(
    page.getByRole("heading", { name: title, level: 1 }),
  ).toBeVisible();
});

Then("the page contains the dashboard heading {string}", async ({ page }, title) => {
  await expect(
    page.getByRole("heading", { name: title, level: 1 }),
  ).toBeVisible();
});
