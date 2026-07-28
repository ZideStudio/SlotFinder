import { render, screen } from "@testing-library/react";

import { appRoutes } from "@Front/routing/appRoutes";
import { BrowserRouter } from "react-router";
import { CardPage } from "../CardPage";

describe("CardPage", () => {
  it("should render logo link, title as level 1 heading and children", () => {
    render(
      <BrowserRouter>
        <CardPage title="Authentication">
          <p>Form content</p>
        </CardPage>
      </BrowserRouter>,
    );

    const logoLink = screen.getByRole("link", { name: "Slot Finder logo" });
    expect(logoLink).toBeInTheDocument();
    expect(logoLink).toHaveAttribute("href", appRoutes.home());

    expect(
      screen.getByRole("heading", { level: 1, name: "Authentication" }),
    ).toBeInTheDocument();

    expect(screen.getByText("Form content")).toBeInTheDocument();
  });

  it("should apply default and custom class names on section", () => {
    render(
      <BrowserRouter>
        <CardPage
          title="Title"
          className="custom-class"
          data-testid="card-page"
        >
          <span>Child</span>
        </CardPage>
      </BrowserRouter>,
    );

    expect(screen.getByTestId("card-page")).toHaveClass(
      "card-page",
      "custom-class",
    );
    expect(
      screen.getByRole("heading", { level: 1, name: "Title" }),
    ).toHaveClass("card-page__title");
  });

  it("should not render logo when isIconDisplayed is false", () => {
    render(
      <CardPage title="Title" isIconDisplayed={false}>
        <span>Child</span>
      </CardPage>,
    );

    expect(
      screen.queryByRole("link", { name: "Slot Finder logo" }),
    ).not.toBeInTheDocument();
  });

  it("should render children inside a bold Card wrapper", () => {
    render(
      <BrowserRouter>
        <CardPage title="Title">
          <span>Child content</span>
        </CardPage>
      </BrowserRouter>,
    );

    const cardContainer = screen.getByText("Child content").closest(".ds-card");
    expect(cardContainer).toBeInTheDocument();
    expect(cardContainer).toHaveClass("ds-card--bold");
  });
});
