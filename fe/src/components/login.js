var m = require("mithril");

const User = require("../actions/user");
const getLoadingView = require("./loading");

const Login = {
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
        "form.form-major",
        {
          onsubmit: e => {
            e.preventDefault();
            // Login.submit();
          }
        },
        [
          // m("div#login-error", Login.current.error),
          m("input[type=email][placeholder=Email] .form-control .together-top", {
            oninput: m.withAttr("value", value => {
              // Login.current.email = value;
            })
          }),
          m("input[type=password][placeholder=Password] .form-control .together-bottom", {
            oninput: m.withAttr("value", value => {
              // Login.current.password = value;
            })
          }),
          m("button[type=submit] .btn .btn-lg .btn-primary", "Sign in")
        ]
      );
    }
  }
};

module.exports = Login;
