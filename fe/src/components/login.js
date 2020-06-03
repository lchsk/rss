var m = require("mithril");

const User = require("../actions/user");
const Login = require("../actions/login");
const getLoadingView = require("./loading");
const Config = require("../config");

const LoginComponent = {
  oninit: node => {
    m.request({
      method: "GET",
      url: Config.api_url + "/users/current",
      withCredentials: true,
      responseType: "json",
      extract: function(xhr, options) {
        if (xhr.status === 401) {
          LoginComponent.signedIn = false;
          return false;
        }

        LoginComponent.signedIn = true;
        return true;
      }
    }).then(function(signedIn) {
      console.log("signed-in", signedIn);
    });
  },
  view: node => {
    console.log("LoginComponent.signedIn", LoginComponent.signedIn);

    if (LoginComponent.signedIn === true) {
      m.route.set("/");
      return;
    }

    if (LoginComponent.signedIn === false) {
      return m(
        "form.form-major",
        {
          onsubmit: e => {
            e.preventDefault();
            Login.submit();
          }
        },
        [
          m("div#login-error", Login.current.error),
          m(
            "input[type=email][placeholder=Email] .form-control .together-top",
            {
              oninput: m.withAttr("value", value => {
                Login.current.email = value;
              })
            }
          ),
          m(
            "input[type=password][placeholder=Password] .form-control .together-bottom",
            {
              oninput: m.withAttr("value", value => {
                Login.current.password = value;
              })
            }
          ),
          m("button[type=submit] .btn .btn-lg .btn-primary", "Sign in")
        ]
      );
    }
  }
};

module.exports = LoginComponent;
