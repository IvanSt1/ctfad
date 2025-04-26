import { errors } from '@vinejs/vine';

export /** @returns {Response} */ function jsonResponse(
  /** @type {Response} */ response
) {
  response.headers.set('Content-Type', 'application/json');
  return response;
}

export /** @returns {Response} */ function htmlResponse(
  /** @type {Response} */ response
) {
  response.headers.set('Content-Type', 'text/html')
  return response
}

export /** @returns {Response} */ function errorResponse(
  /** @type {string} */ message,
  /** @type {any} */ details,
  /** @type {number} */ statusCode
) {
  return jsonResponse(
    new Response(
      JSON.stringify({
        msg: message,
        details: details,
      }),
      {
        status: statusCode
      }
    )
  );
}

export /** @returns {Response | null} */ function requireAuth(
  /** @type {RequestCtx} */ ctx
) {
  if (ctx.user !== null) {
    return null;
  }

  return errorResponse(
    "You must be logged in to perform this action",
    null,
    403
  );
}

/** @returns {Response} */ function handleValidationError(err) {
  if (err instanceof errors.E_VALIDATION_ERROR) {
    return errorResponse(
      "Validation failed",
      err.messages,
      400
    )
  }

  console.error(err);

  return errorResponse(
    "Internal server error",
    err.messages,
    500
  );
}

async function validateOrError(
  /** @type {any} */ data,
  /** @type {import('@vinejs/vine').VineValidator<T, T>} */ validator,
) {
  try {
    const validated = await validator.validate(data);
    return { data: validated, errResp: null };
  } catch (err) {
    return { data: null, errResp: handleValidationError(err) };
  }
}

export async function validateFormDataOrEror(
  /** @type {Request} */ request,
  /** @type {import('@vinejs/vine').VineValidator<T, T>} */ validator
) {
  const rawData = (
    request.body === null
      ? {}
      : Object.fromEntries(await request.formData())
  );
  return await validateOrError(rawData, validator);
}

export async function validateQueryParamsOrError(
  /** @type {Request} */ request,
  /** @type {import('@vinejs/vine').VineValidator<T, T>} */ validator
) {
  const url = new URL(request.url);
  const params = Object.fromEntries(url.searchParams.entries());
  return await validateOrError(params, validator)
}

