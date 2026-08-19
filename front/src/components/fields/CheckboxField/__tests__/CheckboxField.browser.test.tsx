import { screen } from "@testing-library/react";
import { userEvent } from "@vitest/browser/context";
import { useForm, FormProvider } from "react-hook-form";
import { CheckboxField } from "../CheckboxField";
import { type ReactNode } from "react";
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

describe("CheckboxField", () => {
  it("should render checkbox input with correct label and name attribute", async () => {
    await render(
      <FormWrapper>
        <CheckboxField name="acceptTerms" label="Accept Terms" />
      </FormWrapper>,
    );

    expect(screen.getByLabelText("Accept Terms")).toBeInTheDocument();
    expect(screen.getByLabelText("Accept Terms")).toHaveAttribute(
      "name",
      "acceptTerms",
    );
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

    await userEvent.click(screen.getByRole("button", { name: "Trigger error" }));

    const errorMessage = await screen.findByText("This field is required");
    const input = screen.getByRole("checkbox", { name: "Accept Terms" });

    expect(errorMessage).toBeInTheDocument();
    expect(input).toHaveAttribute("aria-invalid", "true");
    expect(input).toHaveAttribute("aria-describedby", errorMessage.id);
  });

  it("updates the input checked state on user click", async () => {
    await render(
      <FormWrapper>
        <CheckboxField name="acceptTerms" label="Accept Terms" />
      </FormWrapper>,
    );

    const checkbox = screen.getByLabelText("Accept Terms");
    await userEvent.click(checkbox);

    expect(checkbox).toBeChecked();
  });
});
