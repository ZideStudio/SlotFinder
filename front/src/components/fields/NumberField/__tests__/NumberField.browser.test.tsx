import { screen } from "@testing-library/react";
import { userEvent } from "@vitest/browser/context";
import { useForm, FormProvider } from "react-hook-form";
import { type ReactNode } from "react";
import { NumberField } from "../NumberField";
import { render } from "vitest-browser-react";

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

describe("NumberField", () => {
  it("should render number input with correct label and name attribute", async () => {
    await render(
      <FormWrapper>
        <NumberField name="number" label="Number" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText("Number")).toBeInTheDocument();
    expect(screen.getByLabelText("Number")).toHaveAttribute("name", "number");
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { number: 0 } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <NumberField name="number" label="Number" />
          <button
            onClick={() =>
              setError("number", { message: "This field is required" })
            }
          >
            Trigger error
          </button>
        </FormProvider>
      );
    };

    await render(<WrapperWithError />);

    await userEvent.click(screen.getByRole("button", { name: "Trigger error" }));

    const errorMessage = await screen.findByText("This field is required");
    const input = screen.getByRole("spinbutton", { name: "Number" });

    expect(errorMessage).toBeInTheDocument();
    expect(input).toHaveAttribute("aria-invalid", "true");
    expect(input).toHaveAttribute("aria-describedby", errorMessage.id);
  });

  it("updates the input value on user typing", async () => {
    await render(
      <FormWrapper>
        <NumberField name="number" label="Number" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText("Number");
    await userEvent.type(input, "1");

    expect(input).toHaveValue(1);
  });
});
