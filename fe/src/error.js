const errorCodes = {
  // API errors
  invalid_username_len: "Invalid username length",
  invalid_password_len: "Invalid password length",
  invalid_email_len: "Invalid email length",
  invalid_email: "Email is not valid",
  invalid_input_format: "Invalid input format",

  // TODO: Add more specific errors
  db_error: "Invalid",

  // UI errors
  passwords_are_not_the_same: "Passwords are not the same"
};

function getErrorMessage(errorCode) {
  if (errorCodes[errorCode]) {
    return errorCodes[errorCode];
  }

  return "Unknown error, please try again";
}

function getSingleError(error) {
  let resp;

  try {
    resp = JSON.parse(error);
  } catch (e) {
    return getErrorMessage("errMsgUnknown");
  }

  if (resp.errors) {
    if (resp.errors.length > 0) {
      const errorCode = resp.errors[0].error_code;

      return getErrorMessage(errorCode);
    }
  }

  return getErrorMessage("errMsgUnknown");
}

module.exports = {
  getErrorMessage: getErrorMessage,
  getSingleError: getSingleError,
};
