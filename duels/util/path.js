function trimTrailingSlashes(
  /** @type {string} */ path
) {
  return path.replace(/\/*$/, "");
}

export function pathsEqual(
  /** @type {string} */ path,
  /** @type {string} */ other
) {
  return trimTrailingSlashes(path) === trimTrailingSlashes(other);
}
