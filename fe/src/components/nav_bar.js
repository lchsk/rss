var m = require("mithril");
const getLink = require("./link");
const LoginComponent = require("./login");

function loginInfo() {
  if (LoginComponent.signedIn === true) {
    return "Signed in";
  }

  return getLink(".nav-link", "/login", "Login");
}

const NavBar = {
  view: () => (
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="collapse navbar-collapse" id="navbarToggler">
        <a class="navbar-brand" href="#">
          <img src="data/text2011.png" alt="rss"/>
        </a>
        <div class="mr-auto mt-2 mt-lg-0"></div>
        <ul class="navbar-nav mt-2 mt-lg-0">
          <li class="nav-item">{loginInfo()}</li>
        </ul>
      </div>
    </nav>
  )
};

module.exports = NavBar;
