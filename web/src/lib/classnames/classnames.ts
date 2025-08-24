export const classNames = (
  ...classes: (string | undefined | { [key: string]: boolean })[]
) => {
  return classes
    .filter((c) => c !== undefined)
    .map((c) =>
      typeof c === "string"
        ? c
        : Object.entries(c)
            .filter(([_, value]) => value)
            .map(([key]) => key)
            .join(" ")
    )
    .join(" ");
};
