import { screen } from "@testing-library/react";
import { userEvent } from "@vitest/browser/context";
import { render } from "vitest-browser-react";
import type { SVGProps } from "react";
import { ClickIcon } from "../ClickIcon";

const TestIcon = (props: SVGProps<SVGSVGElement>) => (
  <svg {...props}>
    <rect width="100" height="100" fill="blue" />
  </svg>
);

describe("ClickIcon", () => {
  it("applies the custom class name", async () => {
    await render(<ClickIcon icon={TestIcon} className="custom-class" />);
    const button = screen.getByRole("button");
    expect(button).toHaveClass("custom-class");
  });

  it("applies props to the button element", async () => {
    const onClick = vi.fn();
    await render(<ClickIcon icon={TestIcon} onClick={onClick} />);
    const button = screen.getByRole("button");
    await userEvent.click(button);
    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('renders as anchor when as="a"', async () => {
    await render(<ClickIcon as="a" icon={TestIcon} href="/test" />);
    const anchor = screen.getByRole("link");
    expect(anchor).toHaveAttribute("href", "/test");
  });

  it("applies the bordered variant class", async () => {
    await render(<ClickIcon icon={TestIcon} variant="bordered" />);
    expect(screen.getByRole("button")).toHaveClass("ds-click-icon--bordered");
  });

  describe("type attribute", () => {
    it('has type="button" by default', async () => {
      await render(<ClickIcon icon={TestIcon} />);
      expect(screen.getByRole("button")).toHaveAttribute("type", "button");
    });

    it('has type="button" when as="button"', async () => {
      await render(<ClickIcon as="button" icon={TestIcon} />);
      expect(screen.getByRole("button")).toHaveAttribute("type", "button");
    });

    it("respects an explicit type prop", async () => {
      await render(<ClickIcon type="submit" icon={TestIcon} />);
      expect(screen.getByRole("button")).toHaveAttribute("type", "submit");
    });

    it("does not set type when rendered as an anchor", async () => {
      await render(<ClickIcon as="a" icon={TestIcon} href="/test" />);
      expect(screen.getByRole("link")).not.toHaveAttribute("type");
    });
  });
});
