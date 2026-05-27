import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { type ComponentProps, type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { Field } from "../Field";

type MockInputProps = ComponentProps<"input"> & {
  error?: string;
  label: string;
};

const MockInput = ({ error, label, ...props }: MockInputProps) => {
  const errorId = error ? `${props.name}-error` : undefined;

  return (
    <>
      <input aria-label={label} aria-describedby={errorId} {...props} />
      {error && (
        <p id={errorId} role="alert">
          {error}
        </p>
      )}
    </>
  );
};

const FormWrapper = ({ children }: { children: ReactNode }) => {
  const methods = useForm({ defaultValues: { user: { email: "" } } });
  return <FormProvider {...methods}>{children}</FormProvider>;
};

describe("Field", () => {
  it("registers the input with react-hook-form", () => {
    render(
      <FormWrapper>
        <Field input={MockInput} name="user.email" label="Email" />
      </FormWrapper>,
    );

    expect(screen.getByRole("textbox", { name: "Email" })).toHaveAttribute(
      "name",
      "user.email",
    );
  });

  it("displays nested error message from react-hook-form", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { user: { email: "" } } });
      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <Field input={MockInput} name="user.email" label="Email" />
          <button
            onClick={() => setError("user.email", { message: "Invalid email" })}
          >
            Trigger error
          </button>
        </FormProvider>
      );
    };

    render(<WrapperWithError />);

    const user = userEvent.setup();
    await user.click(screen.getByRole("button", { name: "Trigger error" }));

    await expect(
      screen.findByText("Invalid email"),
    ).resolves.toBeInTheDocument();
  });
});
