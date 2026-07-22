import { render, screen } from "@testing-library/react";
import { Card } from "../Card";

describe("Card", () => {
  it("renders children", () => {
    render(<Card>Card</Card>);

    expect(screen.getByText("Card")).toBeInTheDocument();
  });

  it("applies additional class names", () => {
    render(<Card className="custom-class">Card</Card>);

    expect(screen.getByText("Card")).toHaveClass("ds-card");
    expect(screen.getByText("Card")).toHaveClass("custom-class");
  });

  it('renders as a different element when "as" prop is provided', () => {
    render(<Card as="section">Card</Card>);

    expect(screen.getByText("Card").tagName).toBe("SECTION");
  });

  it("renders with div by default", () => {
    render(<Card>Card</Card>);

    expect(screen.getByText("Card").tagName).toBe("DIV");
  });

  it("applies border weight class when borderWeight prop is provided", () => {
    render(<Card borderWeight="bold">Card</Card>);

    expect(screen.getByText("Card")).toHaveClass("ds-card--bold");
  });
});
