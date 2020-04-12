import m from "mithril";



const errorCodes = {
  // API errors
  invalid_username_len: "Invalid username length",
  invalid_password_len: "Invalid password length",
  invalid_email_len: "Invalid email length",
  invalid_email: "Email is not valid",
  invalid_input_format: "Invalid input format",
  db_error: "Internal error",

  // UI errors
  passwords_are_not_the_same: "Passwords are not the same",
};

function getErrorMessage(errorCode) {
  if (errorCodes[errorCode]) {
    return errorCodes[errorCode];
  }

  return 'Unknown error, please try again';
}

function getSingleError(error) {
  let resp;

  try {
    resp = JSON.parse(error);
  } catch (e) {
    return getErrorMessage('errMsgUnknown');
  }

  if (resp.errors) {
    if (resp.errors.length > 0) {
      const errorCode = resp.errors[0].error_code;

      return getErrorMessage(errorCode);
    }
  }

  return getErrorMessage('errMsgUnknown');
}

var SignUp = {
  current: {},
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    let password1 = SignUp.current.password1;
    let password2 = SignUp.current.password2;

    if (password1 !== password2) {
      this.setError(getErrorMessage('passwords_are_not_the_same'));
    } else {
      this.setError('');
    }

    if (SignUp.current.error === '') {
    return m.request({
      method: "POST",
      url: "http://localhost:8000/api/users",
      data: {
        email: SignUp.current.email,
        username: SignUp.current.username,
        password: SignUp.current.password1,
      },
    })
      .then(result => {
        m.route.set('/login');
      }).catch(e => {
        this.setError(getSingleError(e.message));
      })
    }
  },
};

var SignUpComponent = {
  oninit: (node) => {
    SignUp.setError("")
  },
  view: (node) => {
    return m(
      "form",
      {
        onsubmit: (e) => {
          e.preventDefault();
          SignUp.submit();
        },
      },
      [
        m("div#signup-error", SignUp.current.error),
        m("input[type=text][placeholder=Email]", {
          oninput: m.withAttr("value", value => {
            SignUp.current.email = value;
          })
        }),
        m("input[type=text][placeholder=Username]", {
          oninput: m.withAttr("value", value => {
            SignUp.current.username = value;
          })
        }),
        m("input[type=password][placeholder=Password]", {
          oninput: m.withAttr("value", value => {
            SignUp.current.password1 = value;
          })
        }),
        m("input[type=password][placeholder=Repeat password]", {
          oninput: m.withAttr("value", value => {
            SignUp.current.password2 = value;
          })
        }),
        m("button[type=submit]", "Sign up"),
      ],
    );
  },
};

export var SignUpComponent;
