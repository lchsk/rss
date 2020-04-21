var m = require("mithril");
var { getErrorMessage, getSingleError } = require("./error");

const User = require("./user");
const getLoadingView = require("./loading");

var Login = {
  current: {},
  setError: function(error) {
    this.current.error = error;
  },
  submit: function() {
    return m
      .request({
        method: "POST",
        url: "http://localhost:8000/api/authentication",
        data: Login.current,
        withCredentials: true
      })
      .then(result => {
        this.setError("");
        m.route.set("/index");
      })
      .catch(e => {
        this.setError(getSingleError(e.message));
      });
  }
};

var LoginComponent = {
  oninit: node => {
    User.load();
  },
  view: node => {
    if (User.authState === User.AuthState.SIGNED_IN) {
      m.route.set("/index");
    } else if (User.authState === User.AuthState.UNKNOWN) {
      return m("div", getLoadingView());
    } else if (User.authState === User.AuthState.SIGNED_OUT) {
      return m(
        "form.form-signin",
        {
          onsubmit: e => {
            e.preventDefault();
            Login.submit();
          }
        },
        [
          m("div#login-error", Login.current.error),
          m("input[type=email][placeholder=Email] .form-control", {
            oninput: m.withAttr("value", value => {
              Login.current.email = value;
            })
          }),
          m("input[type=password][placeholder=Password] .form-control", {
            oninput: m.withAttr("value", value => {
              Login.current.password = value;
            })
          }),
          m("button[type=submit] .btn .btn-lg .btn-primary", "Sign in")
        ]
      );
    }
  }
};

module.exports = LoginComponent;
