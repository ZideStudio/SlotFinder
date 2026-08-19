import { userEvent } from "@vitest/browser/context";
import { useForm, FormProvider } from "react-hook-form";
import { TextareaField } from "../TextareaField";
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

describe("TextareaField", () => {
  it("renders without crashing", async () => {
    await render(
      <FormWrapper>
        <TextareaField
          name="description"
          label="Description"
          placeholder="Enter your description"
        />
      </FormWrapper>,
    );

    await expect
      .element(page.getByLabelText("Description"))
      .toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Description"))
      .toHaveAttribute("name", "description");
    await expect
      .element(page.getByPlaceholder("Enter your description"))
      .toBeInTheDocument();
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { description: "" } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <TextareaField name="description" label="Description" />
          <button
            onClick={() =>
              setError("description", { message: "This field is required" })
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
        <TextareaField name="description" label="Description" />
      </FormWrapper>,
    );

    const input = page.getByLabelText("Description");
    await userEvent.type(input, "Test");

    await expect.element(input).toHaveValue("Test");
  });
});
