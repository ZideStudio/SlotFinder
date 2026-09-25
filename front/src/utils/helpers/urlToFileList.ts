const getMimeTypeExtension = (mimeType: string): string => {
  const mimeToExtension: Record<string, string> = {
    "image/jpeg": "jpg",
    "image/jpg": "jpg",
    "image/png": "png",
    "image/webp": "webp",
  };

  const extension = mimeToExtension[mimeType.toLowerCase()];
  return extension || "jpg";
};

export const urlToFileList = async (
  url: string,
  fileName = "image.jpg",
): Promise<FileList> => {
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error("Impossible de récupérer l'image");
  }

  const blob = await response.blob();
  const blobType = blob.type || "image/jpeg";

  let finalFileName = fileName;
  const extension = getMimeTypeExtension(blobType);

  if (fileName === "image.jpg") {
    finalFileName = `image.${extension}`;
  } else if (fileName.includes(".")) {
    finalFileName = fileName;
  } else {
    finalFileName = `${fileName}.${extension}`;
  }

  const file = new File([blob], finalFileName, {
    type: blobType,
  });

  const dataTransfer = new DataTransfer();
  dataTransfer.items.add(file);

  return dataTransfer.files;
};
