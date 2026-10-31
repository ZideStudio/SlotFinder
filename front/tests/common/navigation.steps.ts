import { When } from "./fixtures";

When("I open the page {string}", async ({ page }, path) => {
  await page.goto(path);
});
