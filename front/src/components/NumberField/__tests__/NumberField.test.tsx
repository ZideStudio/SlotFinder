import { render, screen } from "@testing-library/react";
import { userEvent } from "@testing-library/user-event";
import { useForm, FormProvider } from "react-hook-form";
import { type ReactNode } from "react";
import { NumberField } from "../NumberField";

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
  it("should render number input with correct label and name attribute", () => {
    render(
      <FormWrapper>
        <NumberField name="number" label="Number" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText("Number")).toBeInTheDocument();
    expect(screen.getByLabelText("Number")).toHaveAttribute("name", "number");
  });

  it("displays the error message from form state when validation fails", async () => {
    const user = userEvent.setup();

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

    render(<WrapperWithError />);

    await user.click(screen.getByRole("button", { name: "Trigger error" }));

    const errorMessage = await screen.findByText("This field is required");
    const input = screen.getByRole("spinbutton", { name: "Number" });

    expect(errorMessage).toBeInTheDocument();
    expect(input).toHaveAttribute("aria-invalid", "true");
    expect(input).toHaveAttribute("aria-describedby", errorMessage.id);
  });

  it("updates the input value on user typing", async () => {
    const user = userEvent.setup();

    render(
      <FormWrapper>
        <NumberField name="number" label="Number" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText("Number");
    await user.type(input, "1");

    expect(input).toHaveValue(1);
  });
});
