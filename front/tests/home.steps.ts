import { expect } from "@playwright/test";
import { createBdd } from "playwright-bdd";

const { Given, When, Then } = createBdd();

Given("the application is started with the mocked API", async ({ page }) => {
  await page.goto("/");
});

When("the user opens the home page", async ({ page }) => {
  await page.goto("/");
  await page.waitForLoadState("networkidle");
});

Then("the page contains the title {string}", async ({ page }, title) => {
  await expect(
    page.getByRole("heading", { name: title, level: 1 }),
  ).toBeVisible();
});
