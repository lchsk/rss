const m = require("mithril");
const { getErrorMessage, getSingleError } = require("../common/error");
const LoginComponent = require("./login");
const getLink = require("./link");
const Config = require("../config");

const SignUp = {
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
  submitForm: e => {
    e.preventDefault();

    const usernameInput = document.getElementById("form-username");
    const emailInput = document.getElementById("form-email");
    const passwordInput1 = document.getElementById("form-password-1");
    const passwordInput2 = document.getElementById("form-password-2");

    SignUp.current.username = usernameInput.value;
    SignUp.current.email = emailInput.value;
    SignUp.current.password1 = passwordInput1.value;
    SignUp.current.password2 = passwordInput2.value;

    SignUp.submit();
  },
  view: node => {
    if (LoginComponent.signedIn === true) {
      m.route.set("/");
      return;
    }

    if (LoginComponent.signedIn === false) {
      return (
        <form class="form-major" onsubmit={SignUpComponent.submitForm}>
          <div className="text-center mb-3">
            <img src="../assets/text2011.png" alt="rss" />
          </div>
          <div class="signup-error">{SignUp.current.error}</div>
          <input
            id="form-email"
            type="email"
            placeholder="Email"
            className="form-control together-top"
          />
          <input
            id="form-username"
            type="text"
            placeholder="Username"
            className="form-control together-both"
          />
          <input
            id="form-password-1"
            type="password"
            placeholder="Password"
            className="form-control together-both"
          />
          <input
            id="form-password-2"
            type="password"
            placeholder="Repeat password"
            className="form-control together-bottom"
          />
          <div className="text-center">
            <button type="submit" className="btn btn-primary">
              Sign up
            </button>
            <div className="text-center mt-3">
              {getLink(
                ".mb-0",
                "/login",
                "Sign in if you already have an account"
              )}
            </div>
          </div>
          <div class="text-muted small mt-4">
            By continuing, you agree to{" "}
            <a href="/terms.html">Terms of Service</a> and{" "}
            <a href="/privacy.html">Privacy Policy</a>.
          </div>
        </form>
      );
    }
  }
};

module.exports = SignUpComponent;
