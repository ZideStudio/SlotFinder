import { render, screen } from "@testing-library/react";

import { CardPage } from "../CardPage";

describe("CardPage", () => {
  it("should render title as level 1 heading and children", () => {
    render(
      <CardPage title="Authentication">
        <p>Form content</p>
      </CardPage>,
    );

    expect(
      screen.getByRole("heading", { level: 1, name: "Authentication" }),
    ).toBeInTheDocument();
    expect(screen.getByText("Form content")).toBeInTheDocument();
  });

  it("should apply default and custom class names on section", () => {
    render(
      <CardPage title="Title" className="custom-class" data-testid="card-page">
        <span>Child</span>
      </CardPage>,
    );

    expect(screen.getByTestId("card-page")).toHaveClass(
      "card-page",
      "custom-class",
    );
    expect(
      screen.getByRole("heading", { level: 1, name: "Title" }),
    ).toHaveClass("card-page__title");
  });

  it("should render children inside a bold Card wrapper", () => {
    render(
      <CardPage title="Title">
        <span>Child content</span>
      </CardPage>,
    );

    const cardContainer = screen.getByText("Child content").closest(".ds-card");
    expect(cardContainer).toBeInTheDocument();
    expect(cardContainer).toHaveClass("ds-card--bold");
  });
});
