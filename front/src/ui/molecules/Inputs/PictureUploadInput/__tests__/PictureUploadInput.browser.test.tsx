import { userEvent } from "@vitest/browser/context";
import { render } from "vitest-browser-react";
import { page } from "vitest/browser";
import { PictureUploadInput } from "../PictureUploadInput";

describe("PictureUploadInput", () => {
  beforeAll(() => {
    URL.createObjectURL = vi
      .spyOn(URL, "createObjectURL")
      .mockReturnValue("blob:http://localhost/fake-url");
    URL.revokeObjectURL = vi
      .spyOn(URL, "revokeObjectURL")
      .mockReturnValue(undefined);
  });

  afterAll(() => {
    vi.restoreAllMocks();
  });

  it("should render the component with label", async () => {
    await render(<PictureUploadInput label="Test Label" name="test-input" />);

    await expect.element(page.getByText("Test Label")).toBeInTheDocument();
  });

  it("should apply custom className", async () => {
    await render(
      <PictureUploadInput
        label="Test Label"
        name="test-input"
        className="custom-class"
      />,
    );

    const label = page.getByText("Test Label");
    await expect.element(label).toBeInTheDocument();
    const container = (await label.element()).closest(
      ".ds-picture-upload-input",
    );
    expect(container).toHaveClass("custom-class");
  });

  it("should render error message when error prop is provided", async () => {
    await render(
      <PictureUploadInput
        label="Test Label"
        name="test-input"
        error="This is an error message"
      />,
    );

    await expect
      .element(page.getByText("This is an error message"))
      .toBeInTheDocument();
  });

  it("should render image preview when a valid image file is selected", async () => {
    await render(<PictureUploadInput label="Test Label" name="test-input" />);

    const input = page.getByLabelText("Test Label");
    const file = new File(["test"], "test-image.png", { type: "image/png" });

    await userEvent.upload(input, file);

    await expect.element(page.getByAltText("Preview")).toBeInTheDocument();
    await expect
      .element(page.getByAltText("Preview"))
      .toHaveAttribute("src", "blob:http://localhost/fake-url");
  });

  it("should not render image preview when a non-image file is selected", async () => {
    await render(<PictureUploadInput label="Test Label" name="test-input" />);

    const input = page.getByLabelText("Test Label");
    const file = new File(["dummy content"], "document.pdf", {
      type: "application/pdf",
    });

    await userEvent.upload(input, file);

    await expect
      .element(page.getByAltText("Preview"))
      .toHaveAttribute("hidden");
  });
});
