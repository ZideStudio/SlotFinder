import { userEvent } from "@vitest/browser/context";
import { type ComponentProps, type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
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
  it("registers the input with react-hook-form", async () => {
    await render(
      <FormWrapper>
        <Field input={MockInput} name="user.email" label="Email" />
      </FormWrapper>,
    );

    await expect
      .element(page.getByRole("textbox", { name: "Email" }))
      .toHaveAttribute("name", "user.email");
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

    await render(<WrapperWithError />);

    await userEvent.click(page.getByRole("button", { name: "Trigger error" }));

    await expect.element(page.getByText("Invalid email")).toBeInTheDocument();
  });
});
