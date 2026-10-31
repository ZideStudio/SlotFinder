import { Given } from "@Tests/common/fixtures";

Given("I open the home page", async ({ page }) => {
  await page.goto("/");
});
