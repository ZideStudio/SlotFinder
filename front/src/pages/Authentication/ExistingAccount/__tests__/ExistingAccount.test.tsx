import { renderRoute } from "@Front/utils/testsUtils/customRender/customRender";
import { screen } from "@testing-library/react";
import type { RouteObject } from "react-router";
import { ExistingAccount } from "../ExistingAccount";

const routes: RouteObject[] = [
  { path: "/sign-up", element: <ExistingAccount /> },
  { path: "/login", element: <ExistingAccount /> },
  { path: "/other", element: <ExistingAccount /> },
];

describe("ExistingAccount", () => {
  it("should show a login link on the sign-up page", () => {
    renderRoute({ routes, initialEntry: "/sign-up" });

    expect(screen.getByRole("link", { name: "signUp.login" })).toHaveAttribute(
      "href",
      "/login",
    );
  });

  it("should show a sign-up link on the login page", () => {
    renderRoute({ routes, initialEntry: "/login" });

    expect(
      screen.getByRole("link", { name: "signUp.signUpLink" }),
    ).toHaveAttribute("href", "/sign-up");
  });

  it("should not render outside authentication pages", () => {
    renderRoute({ routes, initialEntry: "/other" });

    expect(screen.queryByRole("link")).not.toBeInTheDocument();
  });
});
