/**
 * Converts FormData to a plain object, handling array fields
 */
export function formDataToObject(formData: FormData, arrayFields: string[] = []): Record<string, unknown> {
  const data: Record<string, unknown> = {};

  for (const key of formData.keys()) {
    if (arrayFields.includes(key)) {
      data[key] = formData.getAll(key);
    } else {
      data[key] = formData.get(key);
    }
  }

  return data;
}

/**
 * Appends values to FormData, handling arrays
 */
export function appendToFormData(
  formData: FormData,
  data: Record<string, unknown>,
  fieldKey: string,
  formDataKey?: string
): void {
  const value = data[fieldKey];
  const key = formDataKey || fieldKey;

  if (Array.isArray(value)) {
    value.forEach((item) => {
      formData.append(key, String(item));
    });
  } else if (value !== undefined && value !== null) {
    formData.append(key, String(value));
  }
}
