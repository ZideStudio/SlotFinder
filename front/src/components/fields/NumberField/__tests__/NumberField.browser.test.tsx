import { userEvent } from "@vitest/browser/context";
import { useForm, FormProvider } from "react-hook-form";
import { type ReactNode } from "react";
import { NumberField } from "../NumberField";
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

describe("NumberField", () => {
  it("should render number input with correct label and name attribute", async () => {
    await render(
      <FormWrapper>
        <NumberField name="number" label="Number" />
      </FormWrapper>,
    );

    await expect.element(page.getByLabelText("Number")).toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Number"))
      .toHaveAttribute("name", "number");
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

    await userEvent.click(page.getByRole("button", { name: "Trigger error" }));

    const errorMessage = page.getByText("This field is required");
    const input = page.getByRole("spinbutton", { name: "Number" });

    await expect.element(errorMessage).toBeInTheDocument();
    await expect.element(input).toHaveAttribute("aria-invalid", "true");
    await expect
      .element(input)
      .toHaveAttribute("aria-describedby", (await errorMessage.element()).id);
  });

  it("updates the input value on user typing", async () => {
    await render(
      <FormWrapper>
        <NumberField name="number" label="Number" />
      </FormWrapper>,
    );

    const input = page.getByLabelText("Number");
    await userEvent.type(input, "1");

    await expect.element(input).toHaveValue(1);
  });
});
