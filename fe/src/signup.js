var m = require("mithril");
var {getErrorMessage, getSingleError} = require("./error");
const User = require("./user");
const getLoadingView = require("./loading")

var SignUp = {
  current: {},
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    let password1 = SignUp.current.password1;
    let password2 = SignUp.current.password2;

    if (password1 !== password2) {
      this.setError(getErrorMessage("passwords_are_not_the_same"));
    } else {
      this.setError("");
    }

    if (SignUp.current.error === "") {
      return m
        .request({
          method: "POST",
          url: "http://localhost:8000/api/users",
          data: {
            email: SignUp.current.email,
            username: SignUp.current.username,
            password: SignUp.current.password1
          }
        })
        .then(result => {
          this.setError("");
          m.route.set("/login");
        })
        .catch(e => {
          this.setError(getSingleError(e.message));
        });
    }
  }
};

var SignUpComponent = {
  oninit: node => {
    SignUp.setError("");
    User.load();
  },
  view: node => {
    if (User.authState === User.AuthState.SIGNED_IN) {
      m.route.set("/index");
    } else if (User.authState === User.AuthState.UNKNOWN) {
      return m("div", getLoadingView());
    } else if (User.authState === User.AuthState.SIGNED_OUT) {
      return m(
        "form",
        {
          onsubmit: e => {
            e.preventDefault();
            SignUp.submit();
          }
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
          m("button[type=submit]", "Sign up")
        ]
      );
    }
  }
};

module.exports = SignUpComponent;
