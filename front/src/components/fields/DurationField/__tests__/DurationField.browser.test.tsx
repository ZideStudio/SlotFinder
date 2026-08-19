import { type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { userEvent } from "@vitest/browser/context";
import { DurationField } from "../DurationField";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";

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
  it("renders without crashing", async () => {
    await render(
      <FormWrapper>
        <DurationField name="duration" legend="Select Date Range" required />
      </FormWrapper>,
    );

    await expect
      .element(page.getByText("Select Date Range"))
      .toBeInTheDocument();
    await expect.element(page.getByText("duration.days")).toBeInTheDocument();
    await expect.element(page.getByText("duration.hours")).toBeInTheDocument();
    await expect
      .element(page.getByText("duration.minutes"))
      .toBeInTheDocument();
  });

  it("displays the error message from form state when validation fails", async () => {
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

    await render(<WrapperWithError />);

    await userEvent.click(page.getByRole("button", { name: "Trigger error" }));

    await expect
      .element(page.getByText("This field is required"))
      .toBeInTheDocument();
  });

  it("updates the input value on user typing", async () => {
    await render(
      <FormWrapper>
        <DurationField name="duration" legend="Duration" />
      </FormWrapper>,
    );

    const daysInput = page.getByRole("spinbutton", { name: "duration.days" });

    await userEvent.type(daysInput, "5");

    await expect.element(daysInput).toHaveValue(5);
  });
});
