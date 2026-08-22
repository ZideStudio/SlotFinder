import { urlToFileList } from "../urlToFileList";

describe("urlToFileList", () => {
  // eslint-disable-next-line @typescript-eslint/init-declarations
  let fetchSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    fetchSpy = vi.spyOn(global, "fetch");
  });

  afterEach(() => {
    fetchSpy.mockRestore();
  });

  it.each([
    ["image/jpeg", "avatar", "avatar.jpg"],
    ["image/png", "profile", "profile.png"],
    ["image/webp", "image", "image.webp"],
  ])(
    "should create file with correct extension for %s",
    async (mimeType, filename, expectedName) => {
      const blob = new Blob(["image data"], { type: mimeType });
      fetchSpy.mockResolvedValueOnce({
        ok: true,
        blob: () => Promise.resolve(blob),
      } as Response);

      const result = await urlToFileList("http://example.com/image", filename);

      expect(result[0].name).toBe(expectedName);
      expect(result[0].type).toBe(mimeType);
      expect(fetchSpy).toHaveBeenCalledWith("http://example.com/image");
    },
  );

  it("should default to jpg when blob type is missing", async () => {
    const blob = new Blob(["image data"]);
    fetchSpy.mockResolvedValueOnce({
      ok: true,
      blob: () => Promise.resolve(blob),
    } as Response);

    const result = await urlToFileList("http://example.com/image", "photo");

    expect(result[0].name).toBe("photo.jpg");
    expect(result[0].type).toBe("image/jpeg");
  });

  it("should preserve filename with extension", async () => {
    const blob = new Blob(["image data"], { type: "image/png" });
    fetchSpy.mockResolvedValueOnce({
      ok: true,
      blob: () => Promise.resolve(blob),
    } as Response);

    const result = await urlToFileList(
      "http://example.com/image",
      "custom.gif",
    );

    expect(result[0].name).toBe("custom.gif");
  });

  it("should throw when fetch fails", async () => {
    fetchSpy.mockResolvedValueOnce({ ok: false } as Response);

    await expect(urlToFileList("http://example.com/404")).rejects.toThrow(
      "Impossible de récupérer l'image",
    );
  });

  it("should throw when fetch throws network error", async () => {
    fetchSpy.mockRejectedValueOnce(new Error("Network error"));

    await expect(urlToFileList("http://example.com/image")).rejects.toThrow(
      "Network error",
    );
  });
});
