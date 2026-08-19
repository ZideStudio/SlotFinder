import { userEvent } from "@vitest/browser/context";
import { useForm, FormProvider } from "react-hook-form";
import { TextField } from "../TextField";
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

describe("TextField", () => {
  it("renders without crashing", async () => {
    await render(
      <FormWrapper>
        <TextField name="email" label="Email" placeholder="Enter your email" />
      </FormWrapper>,
    );

    await expect.element(page.getByLabelText("Email")).toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Email"))
      .toHaveAttribute("name", "email");
    await expect
      .element(page.getByPlaceholder("Enter your email"))
      .toBeInTheDocument();
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { username: "" } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <TextField name="username" label="Username" />
          <button
            onClick={() =>
              setError("username", { message: "This field is required" })
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
        <TextField name="firstName" label="First Name" />
      </FormWrapper>,
    );

    const input = page.getByLabelText("First Name");
    await userEvent.type(input, "Test");

    await expect.element(input).toHaveValue("Test");
  });
});
