import { expect } from "@playwright/test";
import { Then } from "./fixtures";

Then("I see the main heading {string}", async ({ page }, title) => {
  await expect(
    page.getByRole("heading", { name: title, level: 1 }),
  ).toBeVisible();
});

Then("I see the text {string}", async ({ page }, text) => {
  await expect(page.getByText(text).first()).toBeVisible();
});
