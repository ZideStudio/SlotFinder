import { render, screen } from "@testing-library/react";
import { type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import userEvent from "@testing-library/user-event";
import { DurationField } from "../DurationField";

const FormWrapper = ({
  children,
  defaultValues = {},
}: {
  children: ReactNode;
  defaultValues?: Record<string, unknown>;
}) => {
  const methods = useForm({ defaultValues });
  return <FormProvider {...methods}>{children}</FormProvider>;
};

describe("DurationField", () => {
  it("renders without crashing", () => {
    render(
      <FormWrapper>
        <DurationField name="duration" legend="Select Date Range" required />
      </FormWrapper>,
    );

    expect(screen.getByText("Select Date Range")).toBeInTheDocument();
    expect(screen.getByText("duration.days")).toBeInTheDocument();
    expect(screen.getByText("duration.hours")).toBeInTheDocument();
    expect(screen.getByText("duration.minutes")).toBeInTheDocument();
  });

  it("displays the error message from form state when validation fails", async () => {
    const user = userEvent.setup();

    type DurationFormValues = {
      duration: {
        days?: number;
        hours?: number;
        minutes?: number;
      };
    };

    const WrapperWithError = () => {
      const methods = useForm<DurationFormValues>({
        defaultValues: { duration: {} },
      });
      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <DurationField name="duration" legend="Duration" />
          <button
            type="button"
            onClick={() =>
              setError("duration.days", { message: "This field is required" })
            }
          >
            Trigger error
          </button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    await user.click(screen.getByRole("button", { name: "Trigger error" }));

    await expect(
      screen.findByText("This field is required"),
    ).resolves.toBeInTheDocument();
  });

  it("updates the input value on user typing", async () => {
    const user = userEvent.setup();

    render(
      <FormWrapper>
        <DurationField name="duration" legend="Duration" />
      </FormWrapper>,
    );

    const daysInput = screen.getByRole("spinbutton", { name: "duration.days" });

    await user.type(daysInput, "5");

    expect(daysInput).toHaveValue(5);
  });
});
