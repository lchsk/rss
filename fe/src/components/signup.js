var m = require("mithril");
var { getErrorMessage, getSingleError } = require("../common/error");
const User = require("../actions/user");
const getLoadingView = require("./loading");
const LoginComponent = require("./login");

const Config = require("../config");

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
          url: Config.api_url + "/users",
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
            SignUp.submit();
          }
        },
        [
          m("div#signup-error", SignUp.current.error),
          m("input[type=text][placeholder=Email] .form-control", {
            oninput: m.withAttr("value", value => {
              SignUp.current.email = value;
            })
          }),
          m("input[type=text][placeholder=Username] .form-control", {
            oninput: m.withAttr("value", value => {
              SignUp.current.username = value;
            })
          }),
          m("input[type=password][placeholder=Password] .form-control", {
            oninput: m.withAttr("value", value => {
              SignUp.current.password1 = value;
            })
          }),
          m("input[type=password][placeholder=Repeat password] .form-control", {
            oninput: m.withAttr("value", value => {
              SignUp.current.password2 = value;
            })
          }),
          m("button[type=submit] .btn .btn-lg .btn-primary", "Sign up")
        ]
      );
    }
  }
};

module.exports = SignUpComponent;
