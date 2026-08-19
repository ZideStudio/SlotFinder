import { userEvent } from "@vitest/browser/context";
import { useForm, FormProvider } from "react-hook-form";
import { CheckboxField } from "../CheckboxField";
import { type ReactNode } from "react";
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

describe("CheckboxField", () => {
  it("should render checkbox input with correct label and name attribute", async () => {
    await render(
      <FormWrapper>
        <CheckboxField name="acceptTerms" label="Accept Terms" />
      </FormWrapper>,
    );

    await expect
      .element(page.getByLabelText("Accept Terms"))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Accept Terms"))
      .toHaveAttribute("name", "acceptTerms");
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { acceptTerms: false } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <CheckboxField name="acceptTerms" label="Accept Terms" />
          <button
            onClick={() =>
              setError("acceptTerms", { message: "This field is required" })
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
    const input = page.getByRole("checkbox", { name: "Accept Terms" });

    await expect.element(errorMessage).toBeInTheDocument();
    await expect.element(input).toHaveAttribute("aria-invalid", "true");
    await expect
      .element(input)
      .toHaveAttribute("aria-describedby", (await errorMessage.element()).id);
  });

  it("updates the input checked state on user click", async () => {
    await render(
      <FormWrapper>
        <CheckboxField name="acceptTerms" label="Accept Terms" />
      </FormWrapper>,
    );

    const checkbox = page.getByLabelText("Accept Terms");
    await userEvent.click(checkbox);

    await expect.element(checkbox).toBeChecked();
  });
});
