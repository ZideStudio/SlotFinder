import { LinkAtom } from "../LinkAtom";

describe("LinkAtom", () => {
  it("should render a link with required href and children props", () => {
    render(<LinkAtom href="https://example.com">Example Link</LinkAtom>);
    const link = screen.getByRole("link", { name: "Example Link" });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute("href", "https://example.com");
    expect(link).toHaveClass("ds-link-atom");
  });

  it("should open link in a new tab when openInNewTab is true", () => {
    render(
      <LinkAtom href="https://example.com" openInNewTab>
        New Tab Link
      </LinkAtom>
    );
    const link = screen.getByRole("link", { name: "New Tab Link" });
    expect(link).toHaveAttribute("target", "_blank");
    expect(link).toHaveAttribute("rel", "noopener noreferrer");
  });

  it("should apply custom className", () => {
    render(
      <LinkAtom href="https://example.com" className="custom-class">
        Custom Class Link
      </LinkAtom>
    );
    const link = screen.getByRole("link", { name: "Custom Class Link" });
    expect(link).toHaveClass("ds-link-atom custom-class");
  });
});import { render, screen } from "@testing-library/react";
