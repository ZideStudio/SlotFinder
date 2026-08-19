import { userEvent } from "@vitest/browser/context";
import { type ReactNode } from "react";
import { FormProvider, useForm } from "react-hook-form";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { ColorField } from "../ColorField";

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

describe("ColorField", () => {
  it("should render color input with correct label and name attribute", async () => {
    await render(
      <FormWrapper>
        <ColorField name="color" label="Color" description="Choose a color" />
      </FormWrapper>,
    );

    await expect.element(page.getByLabelText("Color")).toBeInTheDocument();
    await expect
      .element(page.getByLabelText("Color"))
      .toHaveAttribute("name", "color");
  });

  it("displays the error message from form state when validation fails", async () => {
    const WrapperWithError = () => {
      const methods = useForm({ defaultValues: { color: "" } });

      const { setError } = methods;

      return (
        <FormProvider {...methods}>
          <ColorField name="color" label="Color" description="Choose a color" />
          <button
            onClick={() =>
              setError("color", { message: "This field is required" })
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
    const input = page.getByLabelText("Color");

    await expect.element(errorMessage).toBeInTheDocument();
    await expect.element(input).toHaveAttribute("aria-invalid", "true");
    await expect
      .element(input)
      .toHaveAttribute("aria-describedby", (await errorMessage.element()).id);
  });

  it("updates the color value on selection", async () => {
    await render(
      <FormWrapper>
        <ColorField name="color" label="Color" description="Choose a color" />
      </FormWrapper>,
    );

    const input = page.getByLabelText("Color");
    await userEvent.fill(input, "#ff0000");

    await expect.element(input).toHaveValue("#ff0000");
  });

  it("should reflect defaultValues in the color input", async () => {
    await render(
      <FormWrapper defaultValues={{ color: "#ff0000" }}>
        <ColorField name="color" label="Color" description="Choose a color" />
      </FormWrapper>,
    );

    await expect.element(page.getByLabelText("Color")).toHaveValue("#ff0000");
  });

  it("should update the color input after reset", async () => {
    const WrapperWithReset = () => {
      const methods = useForm({ defaultValues: { color: "#ff0000" } });
      return (
        <FormProvider {...methods}>
          <ColorField name="color" label="Color" description="Choose a color" />
          <button onClick={() => methods.reset({ color: "#00ff00" })}>
            Reset
          </button>
        </FormProvider>
      );
    };

    await render(<WrapperWithReset />);

    await userEvent.click(page.getByRole("button", { name: "Reset" }));

    await expect.element(page.getByLabelText("Color")).toHaveValue("#00ff00");
  });
});
