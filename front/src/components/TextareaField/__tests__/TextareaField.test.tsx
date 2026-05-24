import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { useForm, FormProvider } from "react-hook-form";
import { TextareaField } from "../TextareaField";
import { type ReactNode } from "react";

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
  it("renders without crashing", () => {
    render(
      <FormWrapper>
        <TextareaField
          name="description"
          label="Description"
          placeholder="Enter your description"
        />
      </FormWrapper>,
    );

    expect(screen.getByLabelText("Description")).toBeInTheDocument();
    expect(screen.getByLabelText("Description")).toHaveAttribute(
      "name",
      "description",
    );
    expect(
      screen.getByPlaceholderText("Enter your description"),
    ).toBeInTheDocument();
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

    render(<WrapperWithError />);

    await userEvent.click(
      screen.getByRole("button", { name: "Trigger error" }),
    );

    await expect(
      screen.findByText("This field is required"),
    ).resolves.toBeInTheDocument();
  });

  it("updates the input value on user typing", async () => {
    render(
      <FormWrapper>
        <TextareaField name="description" label="Description" />
      </FormWrapper>,
    );

    const input = screen.getByLabelText("Description");
    await userEvent.type(input, "Test");

    expect(input).toHaveValue("Test");
  });
});
