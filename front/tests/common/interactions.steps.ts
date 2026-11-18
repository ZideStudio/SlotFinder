import { When } from "./fixtures";

When("I click on the button {string}", async ({ page }, name) => {
  await page.getByRole("button", { name }).click();
});

When("I fill {string} with {string}", async ({ page }, label, value) => {
  await page.getByLabel(label).fill(value);
});
