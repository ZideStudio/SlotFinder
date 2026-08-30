import { render, screen } from "@testing-library/react";

import type { DurationUnit } from "@Front/utils/constants/units";
import { DurationInput } from "../DurationInput";

const UNITS: DurationUnit[] = ["days", "hours", "minutes"];

const buildFields = (errors: Partial<Record<DurationUnit, string>> = {}) =>
  UNITS.map((unit) => ({
    unit,
    name: `duration.${unit}`,
    onChange: vi.fn(),
    onBlur: vi.fn(),
    ref: vi.fn(),
    error: errors[unit],
  }));

describe("DurationInput", () => {
  it("shows the legend and labels for each unit", () => {
    render(<DurationInput legend="Duration" fields={buildFields()} />);

    expect(screen.getByText("Duration")).toBeInTheDocument();
    expect(screen.getByText("duration.days")).toBeInTheDocument();
    expect(screen.getByText("duration.hours")).toBeInTheDocument();
    expect(screen.getByText("duration.minutes")).toBeInTheDocument();
  });

  it("shows the asterisk when required is true", () => {
    render(<DurationInput legend="Duration" required fields={buildFields()} />);
    expect(screen.getByText("*")).toBeInTheDocument();
  });

  it("does not show the asterisk when required is absent", () => {
    render(<DurationInput legend="Duration" fields={buildFields()} />);
    expect(screen.queryByText("*")).not.toBeInTheDocument();
  });

  it("applies min/max based on the unit", () => {
    render(<DurationInput legend="Duration" fields={buildFields()} />);
    const inputs = screen.getAllByRole("spinbutton");

    expect(inputs[0]).toHaveAttribute("min", "0");
    expect(inputs[0]).toHaveAttribute("max", "21");
    expect(inputs[1]).toHaveAttribute("max", "23");
    expect(inputs[2]).toHaveAttribute("max", "59");
  });

  it("shows the error message and aria-invalid when an error is present", () => {
    render(
      <DurationInput
        legend="Duration"
        fields={buildFields({ days: "Required field" })}
      />,
    );

    const inputs = screen.getAllByRole("spinbutton");
    expect(inputs[0]).toHaveAttribute("aria-invalid", "true");
    expect(inputs[0]).toHaveAttribute(
      "aria-describedby",
      expect.stringContaining("days-error"),
    );
    expect(screen.getByText("Required field")).toBeInTheDocument();
  });

  it("does not have aria-invalid when there is no error", () => {
    render(<DurationInput legend="Duration" fields={buildFields()} />);
    screen.getAllByRole("spinbutton").forEach((input) => {
      expect(input).toHaveAttribute("aria-invalid", "false");
    });
  });
});
