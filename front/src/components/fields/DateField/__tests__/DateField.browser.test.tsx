import { type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { DateField } from "../DateField";
import { userEvent } from "@vitest/browser/context";
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

describe("DateField", () => {
  it("renders without crashing", async () => {
    await render(
      <FormWrapper>
        <DateField
          name="birthDate"
          label="Birth Date"
          placeholder="Select your birth date"
        />
      </FormWrapper>,
    );

    await expect.element(page.getByLabelText("Birth Date")).toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Birth Date"))
      .toHaveAttribute("name", "birthDate");
    await expect
      .element(page.getByPlaceholder("Select your birth date"))
      .toBeInTheDocument();
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { birthDate: "" } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <DateField name="birthDate" label="Birth Date" />
          <button
            onClick={() =>
              setError("birthDate", { message: "This field is required" })
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
        <DateField name="eventDate" label="Event Date" />
      </FormWrapper>,
    );

    const input = page.getByLabelText("Event Date");
    await userEvent.type(input, "2024-01-01");
    await expect.element(input).toHaveValue("2024-01-01");
  });
});
